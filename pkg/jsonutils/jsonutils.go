// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package jsonutils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// PrintJSONOn pretty-prints the given data as JSON on the given writer. If the
// given data is not valid JSON an error is returned.
func PrintJSONOn(data any, writer io.Writer) error {
	bites, err := json.MarshalIndent(data, ``, `    `)
	if err != nil {
		return err
	}
	content := string(bites)
	_, err = fmt.Fprintf(writer, "%s\n", content)
	return err
}

// ToJSON returns a pretty-printed JSON string for the given data. If the given
// data is not valid JSON an error is returned. If you are marshaling a large
// data structure and are concerned about the performance of iteratively growing
// the destination strings.Builder then create your own builder, grow it to your
// expected size, and then call PrintJSONOn.
func ToJSON(data any) result.Result[string] {
	writer := new(strings.Builder)
	err := PrintJSONOn(data, writer)
	if err != nil {
		return result.Error[string](err)
	}
	content := writer.String()
	return result.Ok(content)
}

// UnmarshalFromReader inflates the JSON in the given reader and populates the
// given instance. If the data in the reader is not valid JSON an error is
// returned.
func UnmarshalFromReader(reader io.Reader, instance any) error {
	bites, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bites, instance)
	if err != nil {
		return err
	}
	return nil
}

// LoadAs reads the content from the given path and ummarshals it into a result
// of a new instance of T.
func LoadAs[T any](path string) result.Result[T] {
	content := result.New(os.ReadFile(path))
	return result.FlatMap(UnmarshalAs[T], content)
}

// UnmarshalAs ummarshals the given content into a result of anew instance of T.
func UnmarshalAs[T any](content []byte) result.Result[T] {
	var t T
	err := json.Unmarshal(content, &t)
	return result.New(t, err)
}
