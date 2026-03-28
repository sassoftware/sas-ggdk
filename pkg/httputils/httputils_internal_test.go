// SPDX-FileCopyrightText:  2026, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package httputils

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDumpResponseWithFailedBodyRead(t *testing.T) {
	resp := http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: http.StatusInternalServerError,
		Body:       errReadCloser{},
	}
	buffer := &strings.Builder{}
	SetDebug("FULL")
	SetDebugOutput(buffer)
	_ = dumpResponse(&resp) //nolint:bodyclose
	expected := `---------------------------------------------
< SANITIZED
< errReadCloser
=============================================
`
	timestampSanitizer := regexp.MustCompile(`([><]) \d{4}-\d{2}-\d{2}.*\n`)
	actualDump := timestampSanitizer.ReplaceAllString(buffer.String(), "${1} SANITIZED\n")
	require.Equal(t, expected, actualDump)
}

func TestAddRequestContextToError(t *testing.T) {
	err := addRequestContextToError(nil)(nil)
	require.NoError(t, err)
}

type errReadCloser struct {
}

func (e errReadCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New("errReadCloser")
}

func (e errReadCloser) Close() error {
	return nil
}
