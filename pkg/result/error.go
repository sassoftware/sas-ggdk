// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

import "fmt"

// Error returns a Result encapsulating the given error.
func Error[T any](e error) Result[T] {
	return &err[T]{err: e}
}

// err implements Result for errors.
type err[T any] struct {
	err error
}

// Error returns the encapsulated error.
func (e *err[T]) Error() error {
	return e.err
}

// IsError returns true for all errors.
func (e *err[T]) IsError() bool {
	return true
}

// MustGet panics for all errors.
func (e *err[T]) MustGet() T {
	panic(e.err) // Error has no value
}

// OrElse returns the provided default value for all errors.
func (e *err[T]) OrElse(d T) T {
	return d
}

// OrElseGet returns the result of calling the provided function for all errors.
func (e *err[T]) OrElseGet(f Getter[T]) T {
	return f()
}

// String returns a representation of the error.
func (e *err[T]) String() string {
	return fmt.Sprintf("{Error: %#v}", e.err)
}
