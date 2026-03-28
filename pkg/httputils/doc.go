// SPDX-FileCopyrightText:  2026, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package httputils defines utilities for making HTTP calls. Includes functions
for returning the entire response as a result.Result[*http.Response] or just the
body as a result.Result[string] or the result.Result[T] where the response body
is unmarashaled from JSON into a T.

In addition to the response body, if the type T implements the
ResponseHeaderSetter interface then the response headers will be set on the T
before it is returned.

StatusCodes greater than 299 are converted to errors.

Debug output of the request and response can be enabled by the GGDK_HTTP_DEBUG
environment variable. Valid values are anything that strconv.ParseBool accepts
as well as the value "full". A value of true enables basic debug output, while a
value of "full" enables more detailed debug output including the request and
response headers and bodies. By default, this output is written to stderr that
functions are provided to change this destination to any io.Writer.

Example:

	type data struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
		headers http.Header
	}
	func (d *data) SetResponseHeaders(headers http.Header) {
		d.headers = headers
	}
	func main() {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		if err != nil {
			log.Fatal(err)
		}
		d := httputils.DoAsJSON[*data](http.DefaultClient, req)
	}
*/
package httputils
