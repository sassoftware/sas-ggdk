// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stringutils_test

import (
	"fmt"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/stringutils"
	"github.com/stretchr/testify/require"
)

func Test_AnyToString(t *testing.T) {
	expected := `Hello`
	actual := stringutils.AnyToString(expected)
	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func Test_AnyToString_fail(t *testing.T) {
	var expected any
	actual := stringutils.AnyToString(expected)
	require.Nil(t, actual)
	expected = 10
	actual = stringutils.AnyToString(expected)
	require.Nil(t, actual)
}

func Test_ToQuoted_empty(t *testing.T) {
	actual := stringutils.ToQuoted[string]()
	require.Empty(t, actual)
}

func Test_ToQuoted_int(t *testing.T) {
	actual := stringutils.ToQuoted(10, 20, 30)
	expected := []string{
		`"10"`, `"20"`, `"30"`,
	}
	require.Equal(t, expected, actual)
}

func Test_ToQuoted_string(t *testing.T) {
	actual := stringutils.ToQuoted(`red`, `blue`, `green`)
	expected := []string{
		`"red"`, `"blue"`, `"green"`,
	}
	require.Equal(t, expected, actual)
}

func Test_ToStrings(t *testing.T) {
	actual := stringutils.ToStrings(10, 20, 30)
	require.NoError(t, actual.Error())
	expected := []string{
		`10`, `20`, `30`,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_ToStringsWith(t *testing.T) {
	numbers := []int{
		10, 20, 30,
	}
	toHexString := func(value int) result.Result[string] {
		return result.Ok(fmt.Sprintf(`%X`, value))
	}
	actual := stringutils.ToStringsWith(numbers, toHexString)
	require.NoError(t, actual.Error())
	expected := []string{
		`A`, `14`, `1E`,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_ToStringsWith_failure(t *testing.T) {
	numbers := []int{
		10, 20, 30,
	}
	expectedMessage := `failed to create string`
	toHexString := func(value int) result.Result[string] {
		return result.Error[string](errors.New(expectedMessage))
	}
	actual := stringutils.ToStringsWith(numbers, toHexString)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_ToStringWith_nil(t *testing.T) {
	numbers := []int{
		10, 20, 30,
	}
	actual := stringutils.ToStringsWith(numbers, nil)
	require.NoError(t, actual.Error())
	expected := []string{
		`10`, `20`, `30`,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_ToTitle(t *testing.T) {
	actual := stringutils.ToTitle(``)
	require.Empty(t, actual)
	expected := `Hello World`
	actual = stringutils.ToTitle(`hello world`)
	require.Equal(t, expected, actual)
	actual = stringutils.ToTitle(`HELLO WORLD`)
	require.Equal(t, expected, actual)
	expected = "Hello\nWorld"
	actual = stringutils.ToTitle("hello\nworld")
	require.Equal(t, expected, actual)
}

func Test_AsBool(t *testing.T) {
	var input string
	actual := stringutils.AsBool(input, false)
	require.False(t, actual)
	actual = stringutils.AsBool(input, true)
	require.True(t, actual)
	input = `true`
	actual = stringutils.AsBool(input, false)
	require.True(t, actual)
	input = `false`
	actual = stringutils.AsBool(input, true)
	require.False(t, actual)
	input = `xyz`
	actual = stringutils.AsBool(input, false)
	require.False(t, actual)
	actual = stringutils.AsBool(input, true)
	require.True(t, actual)
}

func Test_AsInt(t *testing.T) {
	var input string
	actual := stringutils.AsInt(input, 0)
	require.Equal(t, 0, actual)
	input = `1`
	actual = stringutils.AsInt(input, 0)
	require.Equal(t, 1, actual)
	input = `invalid`
	actual = stringutils.AsInt(input, 0)
	require.Equal(t, 0, actual)
}

func Test_AsInt64(t *testing.T) {
	var input string
	actual := stringutils.AsInt64(input, 10, 33, 0)
	require.Equal(t, int64(0), actual)
	input = `1`
	actual = stringutils.AsInt64(input, 10, 32, 0)
	require.Equal(t, int64(1), actual)
	input = `invalid`
	actual = stringutils.AsInt64(input, 10, 32, 0)
	require.Equal(t, int64(0), actual)
}

func Test_AsFloat(t *testing.T) {
	var input string
	actual := stringutils.AsFloat(input, 32, 0)
	require.Equal(t, 0., actual)
	input = `1.5`
	actual = stringutils.AsFloat(input, 32, 0)
	require.Equal(t, 1.5, actual)
	input = `invalid`
	actual = stringutils.AsFloat(input, 32, 0)
	require.Equal(t, 0., actual)
}

func Test_AsComplex(t *testing.T) {
	var input string
	actual := stringutils.AsComplex(input, 32, 0)
	require.Equal(t, complex128(0), actual)
	input = `1.2`
	actual = stringutils.AsComplex(input, 32, 0)
	require.Equal(t, complex128(1.2), actual)
	input = `invalid`
	actual = stringutils.AsComplex(input, 32, 0)
	require.Equal(t, complex128(0), actual)
}

func Test_AsUint(t *testing.T) {
	var input string
	actual := stringutils.AsUint(input, 10, 32, 0)
	require.Equal(t, uint64(0), actual)
	input = `1`
	actual = stringutils.AsUint(input, 10, 32, 0)
	require.Equal(t, uint64(1), actual)
	input = `invalid`
	actual = stringutils.AsUint(input, 10, 32, 0)
	require.Equal(t, uint64(0), actual)
}
