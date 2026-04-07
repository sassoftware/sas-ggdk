// SPDX-FileCopyrightText:  2026, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package httputils_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/httputils"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	resp := httputils.Do(server.Client(), req)
	_ = resp.MustGet().Body.Close()
	require.False(t, resp.IsError(), "Expected successful response")
	require.Equal(t, http.StatusOK, resp.MustGet().StatusCode, "Expected code to match") //nolint:bodyclose
}

func TestDoAsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d := httputils.DoAsJSON[data](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
}

func TestDoAsJSONFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`I'm a teapot`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d := httputils.DoAsJSON[data](server.Client(), req)
	require.True(t, d.IsError(), "Expected failure response")
	expectedPort := getPort(server.URL)
	expectedMessage := fmt.Sprintf("GET http://127.0.0.1:%s failed\n Caused by:\n\t* 418 I'm a teapot", expectedPort)
	require.ErrorContains(t, d.Error(), expectedMessage)
}

func TestDoAsJSONSetResponseHeadersPointer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d := httputils.DoAsJSON[*dataSetResponseHeadersPointerReceiver](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
	require.Len(t, d.MustGet().headers, 3)
	require.Equal(t, "28", d.MustGet().headers.Get("Content-Length"))
}

func TestDoAsJSONSetResponseHeadersValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d := httputils.DoAsJSON[dataSetResponseHeadersValueReceiver](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
	require.Len(t, d.MustGet().headers, 3)
	require.Equal(t, "28", d.MustGet().headers.Get("Content-Length"))
}

func TestDoAsJSONDumpEnabled(t *testing.T) {
	buffer := &strings.Builder{}
	httputils.SetDebug("TRUE")
	httputils.SetDebugOutput(buffer)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d := httputils.DoAsJSON[data](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
	expectedPort := getPort(server.URL)
	expectedDump := fmt.Sprintf("> GET http://127.0.0.1:%s\n< 200 http://127.0.0.1:%s\n", expectedPort, expectedPort)
	require.Equal(t, expectedDump, buffer.String(), "Expected dump output to match")

	// Ensure dumping can be disabled
	buffer = &strings.Builder{}
	httputils.SetDebug("FALSE")
	httputils.SetDebugOutput(buffer)
	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	d = httputils.DoAsJSON[data](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
	require.Empty(t, buffer.String(), "Expected dump to be empty")
}

func TestDoAsJSONDumpFull(t *testing.T) {
	buffer := &strings.Builder{}
	httputils.SetDebug("FULL")
	httputils.SetDebugOutput(buffer)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		server.URL,
		strings.NewReader("request body"),
	)
	require.NoError(t, err)
	d := httputils.DoAsJSON[data](server.Client(), req)
	require.False(t, d.IsError(), "Expected successful response")
	require.Equal(t, "Hello, World!", d.MustGet().Message, "Expected message to match")
	expectedDump := fmt.Sprintf(`=============================================
> SANITIZED
> GET / HTTP/1.1
> Host: 127.0.0.1:%s
> User-Agent: Go-http-client/1.1
> Content-Length: 12
> Accept-Encoding: gzip
>
> request body
---------------------------------------------
---------------------------------------------
< SANITIZED
< HTTP/1.1 200 OK
< Content-Length: 28
< Content-Type: text/plain; charset=utf-8
< Date: SANITIZED
<
< {"message": "Hello, World!"}
=============================================
`, getPort(server.URL))
	timestampSanitizer := regexp.MustCompile(`([><]) \d{4}-\d{2}-\d{2}.*\n`)
	dateSanitizer := regexp.MustCompile(`([><] Date:) .*\n`)
	emptyLineSanitizer := regexp.MustCompile(`([><]) *\n`)
	actualDump := timestampSanitizer.ReplaceAllString(buffer.String(), "${1} SANITIZED\n")
	actualDump = dateSanitizer.ReplaceAllString(actualDump, "${1} SANITIZED\n")
	actualDump = emptyLineSanitizer.ReplaceAllString(actualDump, "${1}\n")
	require.Equal(t, expectedDump, actualDump, "Expected dump output to match")
}

func TestDoDumpFullFailure(t *testing.T) {
	buffer := &strings.Builder{}
	httputils.SetDebug("FULL")
	httputils.SetDebugOutput(buffer)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, &failedReader{})
	require.NoError(t, err)
	d := httputils.DoAsJSON[data](server.Client(), req)
	require.True(t, d.IsError(), "Expected failure response")
	expectedDump := `=============================================
> SANITIZED
> read error
---------------------------------------------
`
	timestampSanitizer := regexp.MustCompile(`([><]) \d{4}-\d{2}-\d{2}.*\n`)
	actualDump := timestampSanitizer.ReplaceAllString(buffer.String(), "${1} SANITIZED\n")
	require.Equal(t, expectedDump, actualDump, "Expected dump output to match")
}

func TestDoAsString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	s := httputils.DoAsString(server.Client(), req)
	require.False(t, s.IsError(), "Expected successful response")
	require.JSONEq(t, `{"message": "Hello, World!"}`, s.MustGet(), "Expected message to match")
}

func TestDoNoResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	err = httputils.DoNoResponse(server.Client(), req)
	require.NoError(t, err)
}

func TestDoNoResponseFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`Forbidden`))
	}))
	defer server.Close()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	err = httputils.DoNoResponse(server.Client(), req)
	expectedPort := getPort(server.URL)
	expectedMessage := fmt.Sprintf("GET http://127.0.0.1:%s failed\n Caused by:\n\t* 403 Forbidden", expectedPort)
	require.ErrorContains(t, err, expectedMessage)
}

type failedReader struct{}

func (f *failedReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}

type data struct {
	Message string `json:"message"`
}

type dataSetResponseHeadersPointerReceiver struct {
	Message string `json:"message"`
	headers http.Header
}

func (d *dataSetResponseHeadersPointerReceiver) SetResponseHeaders(
	headers http.Header,
) *dataSetResponseHeadersPointerReceiver {
	d.headers = headers
	return d
}

type dataSetResponseHeadersValueReceiver struct {
	Message string `json:"message"`
	headers http.Header
}

func (d dataSetResponseHeadersValueReceiver) SetResponseHeaders(
	headers http.Header,
) dataSetResponseHeadersValueReceiver {
	d.headers = headers
	return d
}

func getPort(_url string) string {
	parsedURL, err := url.Parse(_url)
	if err != nil {
		return ""
	}
	return parsedURL.Port()
}
