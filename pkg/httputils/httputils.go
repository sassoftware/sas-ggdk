// SPDX-FileCopyrightText:  2026, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package httputils (see doc.go for package documentation). The bodyclose
// linter cannot follow the response body through the result chain to ensure it
// is closed, so it is disabled for this file.
// nolint:bodyclose
package httputils

import (
	"io"
	"net/http"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/jsonutils"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// Do executes the HTTP request using the given client and returns a result
// containing the response without further processing. StatusCodes greater than
// 299 are converted to errors with the response body in the error message. The
// caller is responsible for closing the response body.
func Do(c *http.Client, req *http.Request) result.Result[*http.Response] {
	return do[*http.Response](c, req, result.Ok)
}

// DoAsJSON executes the HTTP request using the given client and returns a
// Result containing the unmarshaled response of type T. StatusCodes greater
// than 299 are converted to errors with the response body in the error message.
func DoAsJSON[T any](c *http.Client, req *http.Request) result.Result[T] {
	return do(c, req, unmarshalResponseBodyFromJSON[T])
}

// DoAsString executes the HTTP request using the given client and returns a
// Result containing the response body as a string. StatusCodes greater than 299
// are converted to errors with the response body in the error message.
func DoAsString(c *http.Client, req *http.Request) result.Result[string] {
	return do[string](c, req, getResponseBodyAsString)
}

// DoNoResponse is deprecated and will be removed. Use DoNoResponseBody instead.
func DoNoResponse(c *http.Client, req *http.Request) error {
	return DoNoResponseBody(c, req)
}

// DoNoResponseBody executes the HTTP request using the given client and returns
// an error if the request fails (status code > 299) or nil if it succeeds. The
// response body is closed before returning.
func DoNoResponseBody(c *http.Client, req *http.Request) error {
	res := Do(c, req)
	return result.MapErrorOnly(func(res *http.Response) error {
		return res.Body.Close()
	}, res)
}

// do performs the given HTTP request with the given HTTP client and returns a
// Result containing the response as processed by the given unmarshal function.
// StatusCodes greater than 299 are converted to errors with the response body
// in the error message. The unmarshal function is responsible for closing the
// response body.
func do[T any](c *http.Client, req *http.Request, unmarshal func(*http.Response) result.Result[T]) result.Result[T] {
	// The doc says "It is an error to set this field in an HTTP client request"
	// so we ensure it is empty.
	req.RequestURI = ""
	req = dumpRequest(req)

	// Perform the request and create a Result. Not a SSRF issue because we are
	// performing the request the caller constructed.
	res := result.New(c.Do(req)) //nolint:gosec
	res = result.MapNoError(dumpResponse, res)
	res = result.FlatMap(checkResponseError, res)
	// Map the Result into an unmarshaled instance of T.
	unmarshaledResponse := result.FlatMap(unmarshal, res)

	// Add context if the response indicates failure.
	unmarshaledResponse = result.ErrorMap(addRequestContextToError(req), unmarshaledResponse)

	// Provide response headers to the caller if T implements the
	// ResponseHeaderSetter interface.
	return result.MapNoError2(setResponseHeaders, unmarshaledResponse, res)
}

// checkResponseError checks the HTTP response for errors. If the status code is
// greater than 299, an error is returned with the response body as the ed error
// message.
func checkResponseError(r *http.Response) result.Result[*http.Response] {
	if r.StatusCode > 299 {
		buffer := strings.Builder{}
		_, err := io.Copy(&buffer, r.Body)
		msg := result.New(buffer.String(), err)
		_ = r.Body.Close()
		return result.FlatMap(func(s string) result.Result[*http.Response] {
			return result.Error[*http.Response](errors.New("%d %s", r.StatusCode, s))
		}, msg)
	}
	return result.Ok(r)
}

// addRequestContextToError returns a function appropriate for use with
// results.ErrorMap that adds the HTTP method and URL to an error message for
// additional context.
func addRequestContextToError(req *http.Request) func(error) error {
	return func(err error) error {
		if err != nil {
			return errors.Wrap(err, "%s %s failed", req.Method, req.URL.String())
		}
		return err
	}
}

// ResponseHeaderSetter defines an interface for types that can have response
// headers set on them. This is used to provide response headers to the caller
// if the type returned by DoAsJSON implements this interface. The
// SetResponseHeaders method returns a T so that the method can have a value
// receiver. If the receiver is a pointer, then the type used with DoAsJson
// should be a pointer type.
type ResponseHeaderSetter[T any] interface {
	SetResponseHeaders(http.Header) T //nolint:inamedparam
}

// setResponseHeaders checks if T implements the ResponseHeaderSetter interface
// and, if so, calls the SetResponseHeaders method with the response headers
// from the given http.Response.
func setResponseHeaders[T any](t T, r *http.Response) T {
	setter, ok := any(t).(ResponseHeaderSetter[T])
	if ok {
		return setter.SetResponseHeaders(r.Header)
	}
	return t
}

// unmarshalResponseBodyFromJSON returns a Result of the response body
// unmarshaled into a instance of type T.
func unmarshalResponseBodyFromJSON[T any](r *http.Response) result.Result[T] {
	defer func() {
		_ = r.Body.Close()
	}()
	return jsonutils.UnmarshalFromReader[T](r.Body)
}

// getResponseBodyAsString returns the response body as a string as a Result without further
// processing.
func getResponseBodyAsString(r *http.Response) result.Result[string] {
	defer func() {
		_ = r.Body.Close()
	}()
	buffer := strings.Builder{}
	_, err := io.Copy(&buffer, r.Body)
	return result.New(buffer.String(), err)
}
