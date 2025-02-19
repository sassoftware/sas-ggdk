// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	serverError                = `Server error`
	fileNotFound               = `The file was not found!`
	fileNotFoundTemplate       = `The file '%s' was not found!`
	autoExecBat                = `autoexec.bat`
	applicationFailure         = `Application Failure`
	applicationFailureTemplate = `Application Failure: %d`
)

func Test_Message(t *testing.T) {
	instance := errors.Message(fileNotFound)
	require.NotNil(t, instance)
	requireHasSuffix(t, instance, fileNotFound)
}

func Test_Message_withParameters(t *testing.T) {
	template := fileNotFoundTemplate
	filename := autoExecBat
	expected := fmt.Sprintf(template, filename)
	instance := errors.Message(template, filename)
	require.NotNil(t, instance)
	requireHasSuffix(t, instance, expected)
}

func Test_New(t *testing.T) {
	instance := errors.New(fileNotFound)
	require.NotNil(t, instance)
	require.Error(t, instance)
	actual := instance.Error()
	requireHasSuffix(t, actual, fileNotFound)
}

func Test_New_withParameters(t *testing.T) {
	template := fileNotFoundTemplate
	filename := autoExecBat
	instance := errors.New(template, filename)
	require.NotNil(t, instance)
	require.Error(t, instance)
	suffix := fmt.Sprintf(template, filename)
	actual := instance.Error()
	requireHasSuffix(t, actual, suffix)
}

func Test_Panic(t *testing.T) {
	require.Panics(t, func() {
		errors.Panic(applicationFailure)
	})
}

func Test_ProjectRoot(t *testing.T) {
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	require.NoError(t, err)
	base := filepath.Base(string(root))
	base = strings.TrimSpace(base)
	v := errors.Verbose
	r := errors.ProjectRoot
	errors.Verbose = true
	errors.ProjectRoot = "/" + base + "/"
	defer func() { errors.Verbose = v; errors.ProjectRoot = r }()
	err = errors.New(applicationFailure)
	require.ErrorContains(t, err, applicationFailure)
	require.NotContains(t, err.Error(), base)
}

func Test_Unwrap(t *testing.T) {
	causeMessage := fileNotFound
	cause := errors.New(causeMessage)
	message := applicationFailure
	err := errors.Wrap(cause, message)
	instances := errors.Unwrap(err)
	require.NotNil(t, instances)
	require.Len(t, instances, 2)

	actualErr := instances[0]
	require.NotNil(t, actualErr)
	actualMessage := actualErr.Error()
	requireHasSuffix(t, actualMessage, applicationFailure)

	actualErr = instances[1]
	require.NotNil(t, actualErr)
	actualMessage = actualErr.Error()
	requireHasSuffix(t, actualMessage, fileNotFound)

	instances = errors.Unwrap(nil)
	require.NotNil(t, instances)
	require.Empty(t, instances)
}

func Test_Wrap(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	wrappedMessage := fileNotFound
	cause := errors.New(wrappedMessage)
	message := applicationFailure
	instance := errors.Wrap(cause, message)
	require.NotNil(t, instance)
	require.Error(t, instance)
	actual := instance.Error()
	exp := fmt.Sprintf(`.+\.go:\d+: %s\r?\n Caused by:\r?\n\t\*.+\.go:\d+: %s`, applicationFailure, fileNotFound)
	require.Regexp(t, exp, actual)
}

func Test_Wrap_nil(t *testing.T) {
	instance := errors.Wrap(nil, `error`)
	require.ErrorContains(t, instance, `error`)
	require.Equal(t, 1, len(errors.Unwrap(instance)))
}

func Test_Wrap_twice(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	const error1 = `Error 1`
	const error2 = `Error 2`
	const error3 = `Error 3`
	e1 := errors.New(error1)
	e2 := errors.Wrap(e1, error2)
	e3 := errors.Wrap(e2, error3)
	actual := e3.Error()
	exp := fmt.Sprintf(`.+\.go:\d+: %s\r?\n Caused by:\r?\n\t\* .+:\d+: %s\r?\n\t\* .+:\d+: %s`, error3, error2, error1)
	require.Regexp(t, exp, actual)
}

func Test_Wrap_withParameters(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	wrappedMessage := fileNotFound
	wrappedErr := errors.New(wrappedMessage)
	template := applicationFailureTemplate
	rc := -1
	instance := errors.Wrap(wrappedErr, template, rc)
	require.NotNil(t, instance)
	require.Error(t, instance)
	actual := instance.Error()
	applicationFailureMsg := fmt.Sprintf(template, rc)
	exp := fmt.Sprintf(`.+\.go:\d+: %s\r?\n Caused by:\r?\n\t\*.+\.go:\d+: %s`, applicationFailureMsg, fileNotFound)
	require.Regexp(t, exp, actual)
}

func Test_WrapAll(t *testing.T) {
	wrappedMessage1 := fileNotFound
	cause1 := errors.New(wrappedMessage1)
	wrappedMessage2 := serverError
	cause2 := errors.New(wrappedMessage2)
	message := applicationFailure
	instance := errors.WrapAll([]error{cause1, cause2}, message)
	require.NotNil(t, instance)
	require.Error(t, instance)
	require.ErrorContains(t, instance, fileNotFound)
	require.ErrorContains(t, instance, serverError)
	require.ErrorContains(t, instance, applicationFailure)
}

func Test_WrapAll_nil(t *testing.T) {
	instance := errors.WrapAll(nil, applicationFailure)
	require.NotNil(t, instance)
	require.Error(t, instance)
	require.ErrorContains(t, instance, applicationFailure)
	require.Equal(t, 1, len(errors.Unwrap(instance)))
}

func Test_WrapMessage(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	causeMessage := fileNotFound
	cause := errors.New(causeMessage)
	message := applicationFailure
	instance := errors.WrapMessage(cause, message)
	require.NotNil(t, instance)
	exp := fmt.Sprintf(`.+\.go:\d+: %s\r?\n Caused by:\r?\n\t\*.+\.go:\d+: %s`, applicationFailure, fileNotFound)
	require.Regexp(t, exp, instance)
}

func Test_WrapMessage_nil(t *testing.T) {
	message := applicationFailure
	instance := errors.WrapMessage(nil, message)
	require.NotNil(t, instance)
	require.Contains(t, instance, applicationFailure)
}

func Test_WrapMessage_withParameter(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	cause := errors.New(fileNotFound)
	rc := -1
	instance := errors.WrapMessage(cause, applicationFailureTemplate, rc)
	require.NotNil(t, instance)
	exp := fmt.Sprintf(`.+\.go:\d+: %s: %d\r?\n Caused by:\r?\n\t\*.+\.go:\d+: %s`, applicationFailure, rc, fileNotFound)
	require.Regexp(t, exp, instance)
}

func Test_NilErrors(t *testing.T) {
	e := errors.New("test")
	m := e.(*multierror.Error)
	m.Errors = nil
	require.Equal(t, "", m.Error())
}

func Test_ToStrings(t *testing.T) {
	v := errors.Verbose
	errors.Verbose = true
	defer func() { errors.Verbose = v }()
	const error1 = `Error 1`
	const error2 = `Error 2`
	const error3 = `Error 3`
	e1 := errors.New(error1)
	e2 := errors.New(error2)
	e3 := errors.New(error3)
	actual := errors.ToStrings(e1, e2, e3)
	exp1 := `.+\.go:\d+: Error 1`
	exp2 := `.+\.go:\d+: Error 2`
	exp3 := `.+\.go:\d+: Error 3`
	require.Regexp(t, exp1, actual[0])
	require.Regexp(t, exp2, actual[1])
	require.Regexp(t, exp3, actual[2])
}

func requireHasSuffix(t *testing.T, value string, suffix string) {
	state := strings.HasSuffix(value, suffix)
	require.True(t, state, `the value "%s" does not have the suffix "%s"`, value, suffix)
}
