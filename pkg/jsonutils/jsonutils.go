// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package jsonutils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// PrintJSONOn pretty-prints the given data as JSON on the given writer. If the
// given data is not valid JSON an error is returned.
func PrintJSONOn[T any](data T, writer io.Writer) error {
	bites := result.New(json.MarshalIndent(data, ``, `    `))
	return result.MapErrorOnly(func(b []byte) error {
		_, err := fmt.Fprintf(writer, "%s", b)
		return err
	}, bites)
}

// ToJSON returns a pretty-printed JSON string for the given data. If the given
// data is not valid JSON an error is returned. If you are marshaling a large
// data structure and are concerned about the performance of iteratively growing
// the destination strings.Builder then create your own builder, grow it to your
// expected size, and then call PrintJSONOn. If you want the bytes then call the
// following.
//
//	result.New(json.MarshalIndent(data, ``, `    `))
func ToJSON[T any](data T) result.Result[string] {
	bites := result.New(json.MarshalIndent(data, ``, `    `))
	return result.MapNoError(func(b []byte) string { return string(b) }, bites)
}

// UnmarshalFromReader unmarshals the JSON in the given reader and populates a
// new instance of T. If the data in the reader is not valid JSON an error is
// returned.
func UnmarshalFromReader[T any](reader io.Reader) result.Result[T] {
	bites := result.New(io.ReadAll(reader))
	return result.FlatMap(UnmarshalAs[T], bites)
}

// UnmarshalAs ummarshals the given content into a result of a new instance of T.
func UnmarshalAs[T any](content []byte) result.Result[T] {
	var t T
	err := json.Unmarshal(content, &t)
	return result.New(t, err)
}

// UnmarshalFromReaderInto unmarshals the JSON in the given reader and populates
// the given instance. If the data in the reader is not valid JSON an error is
// returned.
func UnmarshalFromReaderInto[T any](reader io.Reader, value *T) error {
	bites := result.New(io.ReadAll(reader))
	return result.MapErrorOnly(func(b []byte) error {
		return json.Unmarshal(b, value)
	}, bites)
}

// LoadAs reads the content from the given path and ummarshals it into a result
// of a new instance of T.
func LoadAs[T any](path string) result.Result[T] {
	return LoadWith[T](os.ReadFile, path)
}

// LoadWith passes the given path to the given read function to get the content.
// That content is then marshaled into a result of a new instance of T. This is
// useful when using a file system abstraction like
// https://github.com/spf13/afero.
//
//	fs := afero.Afero{Fs: afero.NewMemMapFs()}
//	err := jsonutils.LoadWith[string](fs.ReadFile, "/tmp/person.json")
func LoadWith[T any](
	read func(string) ([]byte, error),
	path string,
) result.Result[T] {
	content := result.New(read(path))
	return result.FlatMap(UnmarshalAs[T], content)
}

// Save marshals the given T and writes it to a file at the given path with the
// given permissions. The file is truncated if it exists.
func Save[T any](value T, path string, perm os.FileMode) error {
	content := result.New(json.Marshal(value))
	writeF := func(b []byte) error {
		return os.WriteFile(path, b, perm)
	}
	return result.MapErrorOnly(writeF, content)
}

// SaveWith marshals the given T and calls the given writeFunc with the
// resulting content, given path, and given permissions.This is
// useful when using a file system abstraction like
// https://github.com/spf13/afero.
//
//	fs := afero.Afero{Fs: afero.NewMemMapFs()}
//	err := jsonutils.SaveWith(fs.WriteFile, instance, "/tmp/person.json", 0700)
func SaveWith[T any](
	writeFunc func(string, []byte, os.FileMode) error,
	value T,
	path string,
	perm os.FileMode,
) error {
	content := result.New(json.Marshal(value))
	writeF := func(b []byte) error {
		return writeFunc(path, b, perm)
	}
	return result.MapErrorOnly(writeF, content)
}
