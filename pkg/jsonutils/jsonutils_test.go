// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package jsonutils_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/jsonutils"
	"github.com/stretchr/testify/require"
)

const (
	jsonFilename = `person.json`
	personName   = `John Smith`
	personAge    = 75
)

type person struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Spouse *person `json:"spouse,omitempty"`
}

func Test_PrintJSONOn(t *testing.T) {
	expected, err := readTestdata(jsonFilename)
	require.NoError(t, err)
	instance := &person{
		Name: personName,
		Age:  personAge,
	}
	buffer := bytes.NewBuffer(nil)
	err = jsonutils.PrintJSONOn(instance, buffer)
	require.NoError(t, err)
	actual := buffer.Bytes()
	require.Equal(t, bytes.TrimSpace(expected), actual)
}

func Test_PrintJSONOn_fail(t *testing.T) {
	instance := &person{
		Name: personName,
		Age:  personAge,
	}
	// You cannot be your own spouse!
	instance.Spouse = instance
	buffer := bytes.NewBuffer(nil)
	err := jsonutils.PrintJSONOn(instance, buffer)
	require.EqualError(t, err, `json: unsupported value: encountered a cycle via *jsonutils_test.person`)
	actual := buffer.Bytes()
	require.Empty(t, actual)
}

func Test_ToJSON(t *testing.T) {
	bites, err := readTestdata(jsonFilename)
	require.NoError(t, err)
	expected := string(bites)
	instance := new(person)
	err = json.Unmarshal(bites, instance)
	require.NoError(t, err)
	actual := jsonutils.ToJSON(instance)
	require.NoError(t, actual.Error())
	require.Equal(t, strings.TrimSpace(expected), actual.MustGet())
}

func Test_ToJSON_fail(t *testing.T) {
	instance := &person{
		Name: personName,
		Age:  personAge,
	}
	// You cannot be your own spouse!
	instance.Spouse = instance
	actual := jsonutils.ToJSON(instance)
	expected := `json: unsupported value: encountered a cycle via *jsonutils_test.person`
	require.ErrorContains(t, actual.Error(), expected)
}

func Test_UnmarshalFromReader(t *testing.T) {
	reader, err := toTestdataFile(jsonFilename)
	require.NoError(t, err)
	r := jsonutils.UnmarshalFromReader[person](reader)
	require.False(t, r.IsError())
	require.Equal(t, personName, r.MustGet().Name)
	require.Equal(t, personAge, r.MustGet().Age)
}

func Test_UnmarshalFromReaderInto(t *testing.T) {
	reader, err := toTestdataFile(jsonFilename)
	require.NoError(t, err)
	instance := new(person)
	err = jsonutils.UnmarshalFromReaderInto(reader, instance)
	require.NoError(t, err)
	require.Equal(t, personName, instance.Name)
	require.Equal(t, personAge, instance.Age)
}

type failingReader struct{}

func (reader *failingReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New(`failed to read`)
}

func Test_UnmarshalFromReader_failingReader(t *testing.T) {
	reader := new(failingReader)
	r := jsonutils.UnmarshalFromReader[int](reader)
	require.True(t, r.IsError())
	require.Error(t, r.Error())
}

func Test_UnmarshalFromReader_failingUnmarshal(t *testing.T) {
	reader := strings.NewReader(`not JSON`)
	r := jsonutils.UnmarshalFromReader[int](reader)
	require.True(t, r.IsError())
	require.Error(t, r.Error())
}

func Test_UnmarshalFromReaderInto_failingReader(t *testing.T) {
	reader := new(failingReader)
	err := jsonutils.UnmarshalFromReaderInto[int](reader, nil)
	require.Error(t, err)
}

func Test_UnmarshalFromReaderInto_failingUnmarshal(t *testing.T) {
	reader := strings.NewReader(`not JSON`)
	err := jsonutils.UnmarshalFromReaderInto[int](reader, nil)
	require.Error(t, err)
}

func Test_LoadAs(t *testing.T) {
	expected := person{
		Name: "John Smith",
		Age:  75,
	}
	actual := jsonutils.LoadAs[person](toTestdataFilename(jsonFilename))
	require.False(t, actual.IsError())
	require.Equal(t, expected, actual.MustGet())
}

func Test_LoadAs_ptr(t *testing.T) {
	expectedPtr := &person{
		Name: "John Smith",
		Age:  75,
	}
	actualPtr := jsonutils.LoadAs[*person](toTestdataFilename(jsonFilename))
	require.False(t, actualPtr.IsError())
	require.Equal(t, expectedPtr, actualPtr.MustGet())
}

func Test_LoadAs_fail(t *testing.T) {
	actual := jsonutils.LoadAs[person](toTestdataFilename("missing.json"))
	require.True(t, actual.IsError())
}

func Test_LoadWith(t *testing.T) {
	expected := person{
		Name: "John Smith",
		Age:  75,
	}
	readF := func(_ string) ([]byte, error) {
		return []byte(`{"name":"John Smith","age":75}`), nil
	}
	actual := jsonutils.LoadWith[person](readF, "/tmp/person.json")
	require.False(t, actual.IsError())
	require.Equal(t, expected, actual.MustGet())
}

func Test_LoadWith_ptr(t *testing.T) {
	expectedPtr := &person{
		Name: "John Smith",
		Age:  75,
	}
	readF := func(_ string) ([]byte, error) {
		return []byte(`{"name":"John Smith","age":75}`), nil
	}
	actualPtr := jsonutils.LoadWith[*person](readF, "/tmp/person.json")
	require.False(t, actualPtr.IsError())
	require.Equal(t, expectedPtr, actualPtr.MustGet())
}

func Test_LoadWith_fail(t *testing.T) {
	readF := func(_ string) ([]byte, error) {
		return []byte{}, errors.New("failed")
	}
	actual := jsonutils.LoadWith[person](readF, "/tmp/missing.json")
	require.True(t, actual.IsError())
}

func Test_Save(t *testing.T) {
	instance := person{
		Name: "John Smith",
		Age:  75,
	}
	tmpdir := t.TempDir()
	path := filepath.Join(tmpdir, "person.json")
	err := jsonutils.Save(instance, path, 0700)
	require.NoError(t, err)
	defer func() { require.NoError(t, os.RemoveAll(tmpdir)) }()
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func Test_SaveWith(t *testing.T) {
	instance := person{
		Name: "John Smith",
		Age:  75,
	}
	called := false
	saveF := func(_ string, _ []byte, _ os.FileMode) error {
		called = true
		return nil
	}
	err := jsonutils.SaveWith(saveF, instance, "/tmp/person.json", 0700)
	require.NoError(t, err)
	require.True(t, called)
}

func Test_SaveWith_failed(t *testing.T) {
	instance := person{
		Name: "John Smith",
		Age:  75,
	}
	saveF := func(_ string, _ []byte, _ os.FileMode) error {
		return errors.New("failed")
	}
	err := jsonutils.SaveWith(saveF, instance, "/tmp/person.json", 0700)
	require.Error(t, err)
}

func Test_UnmarshalAs(t *testing.T) {
	expectedPtr := person{
		Name: "John Smith",
		Age:  75,
	}
	content, err := readTestdata(jsonFilename)
	require.NoError(t, err)
	actualPtr := jsonutils.UnmarshalAs[person](content)
	require.False(t, actualPtr.IsError())
	require.Equal(t, expectedPtr, actualPtr.MustGet())
}

func Test_UnmarshalAs_ptr(t *testing.T) {
	expectedPtr := &person{
		Name: "John Smith",
		Age:  75,
	}
	content, err := readTestdata(jsonFilename)
	require.NoError(t, err)
	actualPtr := jsonutils.UnmarshalAs[*person](content)
	require.False(t, actualPtr.IsError())
	require.Equal(t, expectedPtr, actualPtr.MustGet())
}

func Test_UnmarshalAs_fail(t *testing.T) {
	actual := jsonutils.UnmarshalAs[person]([]byte(`{"invalid": "json`))
	require.True(t, actual.IsError())
}

func readTestdata(elements ...string) ([]byte, error) {
	filename := toTestdataFilename(elements...)
	path := filepath.Clean(filename)
	return os.ReadFile(path)
}

func toTestdataFile(element string) (*os.File, error) {
	filename := toTestdataFilename(element)
	path := filepath.Clean(filename)
	return os.Open(path)
}

func toTestdataFilename(elements ...string) string {
	s := []string{
		`testdata`,
	}
	s = append(s, elements...)
	return filepath.Join(s...)
}
