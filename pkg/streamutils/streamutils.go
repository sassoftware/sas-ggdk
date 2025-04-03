// SPDX-FileCopyrightText:  2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package streamutils

import (
	"bytes"
	"io"
	"os"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// createPipeFunc is a hook for testing.
// nolint:revive
var createPipeFunc func() (fileI, fileI, error) = createPipe

// CaptureStdStreams captures the stdout and stderr of a function that returns a
// result.
// nolint:nakedret
func CaptureStdStreams[T any](fn func() result.Result[T]) (res result.Result[T], stdout, stderr result.Result[[]byte]) {
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()
	outReader, outWriter, err := createPipeFunc()
	if err != nil {
		res = result.Error[T](err)
		stdout = result.Error[[]byte](err)
		stderr = result.Error[[]byte](err)
		return
	}
	errReader, errWriter, err := createPipeFunc()
	if err != nil {
		res = result.Error[T](err)
		stdout = result.Error[[]byte](err)
		stderr = result.Error[[]byte](err)
		return
	}
	stdoutCh := make(chan result.Result[[]byte])
	stderrCh := make(chan result.Result[[]byte])
	go func() {
		buf := bytes.Buffer{}
		_, err := io.Copy(&buf, outReader.GetFile())
		stdoutCh <- result.New(buf.Bytes(), err)
	}()
	go func() {
		buf := bytes.Buffer{}
		_, err := io.Copy(&buf, errReader.GetFile())
		stderrCh <- result.New(buf.Bytes(), err)
	}()
	os.Stdout = outWriter.GetFile()
	os.Stderr = errWriter.GetFile()
	res = fn()
	err = outWriter.Close()
	if err != nil {
		stdout = result.Error[[]byte](err)
	} else {
		stdout = <-stdoutCh
	}
	err = errWriter.Close()
	if err != nil {
		stderr = result.Error[[]byte](err)
	} else {
		stderr = <-stderrCh
	}
	return
}

// fileI is a wrapper around os.File so that testware can simulate a file that
// fails to close.
type fileI interface {
	GetFile() *os.File
	Close() error
}

// file is an implementation of fileI that delegates to os for all
// functionality.
type file struct {
	*os.File
}

// GetFile returns the wrapped os.File.
func (f *file) GetFile() *os.File {
	return f.File
}

// createPipe returns wrapped os.Files from os.Pipe.
func createPipe() (fileI, fileI, error) {
	r, w, err := os.Pipe()
	return &file{File: r}, &file{File: w}, err
}
