// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package errors // nolint revive

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/sassoftware/sas-ggdk/pkg/condition"
)

const (
	indentation = "\n\t "
	newline     = "\n"
	space       = ` `
	star        = `*`
	tab         = "\t"
)

// Verbose is used to control verbose output. When true, error messages will
// include the caller's location.
var Verbose = false

// ProjectRoot is used to trim paths in verbose output. Anything in the path up
// to, and including ProjectRoot will be removed from the front of the output
// paths.
var ProjectRoot = ``

// New returns an error formatted according to the given template and
// parameters.
func New(template string, params ...any) error {
	template = addCallerLocation(template)
	err := fmt.Errorf(template, params...)
	e := multierror.Append(nil, err)
	e.ErrorFormat = printErrors
	return e
}

// Wrap returns an error formatted according to the given template and
// parameters with the given error as its cause.
func Wrap(cause error, template string, params ...any) error {
	template = addCallerLocation(template)
	err := fmt.Errorf(template, params...)
	if condition.IsNil(cause) {
		return err
	}
	e := multierror.Append(err, cause)
	e.ErrorFormat = printErrors
	return e
}

// WrapAll returns an error formatted according to the given template and
// parameters with the given errors as its cause.
func WrapAll(causes []error, template string, params ...any) error {
	template = addCallerLocation(template)
	err := fmt.Errorf(template, params...)
	if condition.IsNil(causes) {
		return err
	}
	e := multierror.Append(err, causes...)
	e.ErrorFormat = printErrors
	return e
}

// Message returns a message formatted according to the given template and
// parameters.
func Message(template string, params ...any) string {
	template = addCallerLocation(template)
	return fmt.Sprintf(template, params...)
}

// Panic panics with a message formatted the way an error returned from this
// package would be.
func Panic(template string, params ...any) {
	template = addCallerLocation(template)
	message := fmt.Sprintf(template, params...)
	panic(message)
}

// Unwrap takes an error, which can be a single error or a multi-error, and
// unwraps it into a flat slice of errors.
func Unwrap(err error) []error {
	var result = make([]error, 0, 3)
	if err != nil {
		all := &multierror.Error{}
		ok := errors.As(err, &all)
		if ok {
			// Collect each error in the multierror...
			result = append(result, all.Errors...)
		} else {
			// Not a multierror...
			result = append(result, err)
		}
	}
	return result
}

// WrapMessage returns a message formatted according to the given template and
// parameters with the given error as its cause.
func WrapMessage(cause error, template string, params ...any) string {
	template = addCallerLocation(template)
	err := fmt.Errorf(template, params...)
	if condition.IsNil(cause) {
		return err.Error()
	}
	e := multierror.Append(err, cause)
	e.ErrorFormat = printErrors
	return e.Error()
}

// ToStrings takes a slice of errors and converts it to a slice of strings.
func ToStrings(errors ...error) []string {
	strs := make([]string, 0, len(errors))
	for _, err := range errors {
		strs = append(strs, err.Error())
	}
	return strs
}

func addCallerLocation(template string) string {
	if !Verbose {
		return template
	}
	_, path, line, ok := runtime.Caller(2)
	if ok {
		filename := path
		if len(ProjectRoot) > 0 {
			index := strings.LastIndex(path, ProjectRoot)
			if index != -1 {
				index += len(ProjectRoot)
				filename = path[index:]
			}
		}
		template = fmt.Sprintf(`%s:%d: %s`, filename, line, template)
	}
	return template
}

func printErrors(errors []error) string {
	size := len(errors)
	if size == 0 {
		return ""
	}
	writer := new(strings.Builder)
	capacity := size * 75
	writer.Grow(capacity)
	printErrorOn(errors[0], writer)
	if size > 1 {
		printCausesOn(errors[1:], writer)
	}
	result := writer.String()
	return result
}

func printCauseOn(err error, writer io.StringWriter) {
	printOn(newline, writer)
	printOn(tab, writer)
	printOn(star, writer)
	printOn(space, writer)
	msg := err.Error()
	msg = strings.ReplaceAll(msg, newline, indentation)
	printOn(msg, writer)
}

func printCausesOn(errors []error, writer io.StringWriter) {
	printOn(newline, writer)
	printOn(space, writer)
	printOn(`Caused by:`, writer)
	for _, err := range errors {
		printCauseOn(err, writer)
	}
}

func printErrorOn(err error, writer io.StringWriter) {
	msg := err.Error()
	printOn(msg, writer)
}

func printOn(value string, writer io.StringWriter) {
	_, _ = writer.WriteString(value)
}
