// SPDX-FileCopyrightText:  2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package streamutils

import (
	"os"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_FailedStdoutPipe(t *testing.T) {
	createPipeFunc = func() (r fileI, w fileI, err error) {
		return nil, nil, errors.New("pipe failed")
	}
	defer func() { createPipeFunc = createPipe }()
	f := func() result.Result[bool] {
		return result.Ok(true)
	}
	res, stdout, stderr := CaptureStdStreams(f)
	require.ErrorContains(t, res.Error(), "pipe failed")
	require.ErrorContains(t, stdout.Error(), "pipe failed")
	require.ErrorContains(t, stderr.Error(), "pipe failed")
}

func Test_FailedStderrPipe(t *testing.T) {
	count := 0
	createPipeFunc = func() (r fileI, w fileI, err error) {
		if count == 0 {
			count++
			return createPipe()
		}
		return nil, nil, errors.New("pipe failed")
	}
	defer func() { createPipeFunc = createPipe }()
	f := func() result.Result[bool] {
		return result.Ok(true)
	}
	res, stdout, stderr := CaptureStdStreams(f)
	require.ErrorContains(t, res.Error(), "pipe failed")
	require.ErrorContains(t, stdout.Error(), "pipe failed")
	require.ErrorContains(t, stderr.Error(), "pipe failed")
}

func Test_FailedClose(t *testing.T) {
	count := 0
	createPipeFunc = func() (fileI, fileI, error) {
		msg := "stdout stream failed"
		if count > 0 {
			msg = "stderr stream failed"
		}
		count++
		r, w, err := os.Pipe()
		require.NoError(t, err)
		return &fileFailCLose{File: r, msg: msg},
			&fileFailCLose{File: w, msg: msg},
			nil
	}
	defer func() { createPipeFunc = createPipe }()
	f := func() result.Result[bool] {
		return result.Ok(true)
	}
	res, stdout, stderr := CaptureStdStreams(f)
	require.NoError(t, res.Error())
	require.ErrorContains(t, stdout.Error(), "stdout stream failed")
	require.ErrorContains(t, stderr.Error(), "stderr stream failed")
}

type fileFailCLose struct {
	*os.File
	msg string
}

func (f *fileFailCLose) GetFile() *os.File {
	return f.File
}
func (f *fileFailCLose) Close() error {
	_ = f.File.Close()
	return errors.New(f.msg)
}
