// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package result defines interfaces and structs for encapsulating a value or an
error.

Without Result, code must continually handle errors in the middle of business
logic. With Result, the algorithm is more clear and error handling is
abstracted.

Three constructors are provided for creating a Result; Ok which creates a Result
encapsulating a value, Error which creates a Result encapsulating an error, and
New which creates a Result from a value and an error. New is useful for
converting return values from functions that return both a value and an error
into a Result. See those functions for more information.

Usage:

Without Result:

	func getPostCount() (int, error) {
		resp, err := http.Get("https://api.project.com/posts")
		if err != nil {
			return 0, errors.New("failed to fetch posts: cause %v", e)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		posts := []*Post{}
		err = json.Unmarshal(body, &posts)
		if err != nil {
			return 0, err
		}
		return len(posts), nil
	}

Using Result:

	func getPostCount() result.Result[int] {
		resp := fetch("https://api.project.com/posts")
		body := result.FlatMap(readAll, resp)
		posts := result.FlatMap(marshal, body)
		return result.MapNoError(count, posts)
	}

	func fetch(url string) result.Result[*http.Response] {
		resp := result.New(http.Get(url))
		return result.ErrorMap(addErrorContext, resp)
	}

	func addErrorContext(err error) error {
		return errors.New("failed to fetch posts: cause %v", e)
	}

	func readAll(r *http.Response) result.Result[[]byte] {
		defer r.Body.Close()
		return result.New(io.ReadAll(r.Body))
	}

	func marshal(content []byte) result.Result[[]*Post] {
		posts := []*Post{}
		err := json.Unmarshal(content, &posts)
		return result.New(posts, err)
	}

	func count(posts []*Post) int {
		return len(posts)
	}
*/
package result
