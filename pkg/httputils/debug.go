// SPDX-FileCopyrightText:  2026, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package httputils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

// dumpRequest prints the HTTP request and response to standard error based on
// the configured debug level. In enabled mode, only the method, URL, and status
// code are printed. In full mode, the entire request (including headers and
// body) is dumped.
// gosec flags the Fprintf calls as potential XSS vulnerabilities even though
// that analysis is intended to be limited to http.ResponseWriter calls. This
// false positive is fixed in https://github.com/securego/gosec/issues/1548.
func dumpRequest(req *http.Request) *http.Request {
	switch debugLevel {
	case enabled:
		_, _ = fmt.Fprintf(debugOutput, "> %s %s\n", req.Method, req.URL.String()) //nolint:gosec
	case full:
		var requestDumpString string
		requestDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			requestDumpString = err.Error()
		} else {
			requestDumpString = strings.ReplaceAll(string(requestDump), "\r\n", "\n")
			requestDumpString = strings.ReplaceAll(requestDumpString, "\n", "\n> ")
		}
		_, _ = fmt.Fprintf(debugOutput, "=============================================\n")
		_, _ = fmt.Fprintf(debugOutput, "> %v\n", time.Now())
		_, _ = fmt.Fprintf(debugOutput, "> %s\n", requestDumpString) //nolint:gosec
		_, _ = fmt.Fprintf(debugOutput, "---------------------------------------------\n")
	default:
	}
	return req
}

// dumpResponse prints the HTTP response to standard error based on the
// configured debug level. In enabled mode, only the status code and URL are
// printed. In full mode, the entire response (including headers and body) is
// dumped.
// gosec flags the Fprintf calls as potential XSS vulnerabilities even though
// that analysis is intended to be limited to http.ResponseWriter calls. This
// false positive is fixed in https://github.com/securego/gosec/issues/1548.
func dumpResponse(res *http.Response) *http.Response {
	switch debugLevel {
	case enabled:
		_, _ = fmt.Fprintf(debugOutput, "< %d %s\n", res.StatusCode, res.Request.URL.String()) //nolint:gosec
	case full:
		var responseDumpString string
		responseDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			responseDumpString = err.Error()
		} else {
			responseDumpString = strings.ReplaceAll(string(responseDump), "\r\n", "\n")
			responseDumpString = strings.ReplaceAll(responseDumpString, "\n", "\n< ")
		}
		_, _ = fmt.Fprintf(debugOutput, "---------------------------------------------\n")
		_, _ = fmt.Fprintf(debugOutput, "< %v\n", time.Now())
		_, _ = fmt.Fprintf(debugOutput, "< %s\n", responseDumpString)
		_, _ = fmt.Fprintf(debugOutput, "=============================================\n")
	default:
	}
	return res
}

// httpDebugLevel represents the level of HTTP debugging to perform. It can be
// set via the HTTP_DEBUG environment variable, which can be "full", or a value
// that strconv can parse (1, t, T, TRUE, true, True, 0, f, F, FALSE, false,
// False).
type httpDebugLevel int

// debugLevel is the configured level of HTTP debugging to perform. It is set
// during package initialization based on the HTTP_DEBUG environment variable.
// It can be changed at runtime via the SetDebug function.
var debugLevel httpDebugLevel

// debugOutput is the output destination for HTTP debug information. It defaults
// to os.Stderr and can be changed via the SetDebugOutput function.
var debugOutput io.Writer = os.Stderr

const (
	disabled httpDebugLevel = iota
	enabled
	full
)

// SetDebug sets the debug level based on the given string. It accepts "disabled", "enabled",
// "full", or any value that strconv.ParseBool can parse (1, t, T, TRUE, true,
// True, 0, f, F, FALSE, false, False). If the input cannot be parsed, the debug
// level defaults to disabled.
func SetDebug(level string) {
	b, err := strconv.ParseBool(level)
	if err == nil {
		switch b {
		case true:
			debugLevel = enabled
		case false:
			debugLevel = disabled
		}
		return
	}
	// "full" is the only non-boolean value allowed.
	switch strings.ToLower(level) {
	case "full":
		debugLevel = full
	default:
		debugLevel = disabled
	}
}

// SetDebugOutput sets the output destination for HTTP debug information.
func SetDebugOutput(w io.Writer) {
	debugOutput = w
}

const httpDebugEnv = "GGDK_HTTP_DEBUG"

// init initializes the debugLevel variable based on the HTTP_DEBUG environment
// variable.
// nolint:gochecknoinits
func init() {
	SetDebug(os.Getenv(httpDebugEnv))
}
