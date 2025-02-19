// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

import (
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
)

// Getter defines a function that gets the default value for a Result that
// encapsulates an error.
type Getter[T any] func() T

// Result encapsulates a value or an error. An instance implementing this
// interface will either contain a value or an error. Not both.
type Result[T any] interface {
	// Error returns nil if the Result encapsulates an error or the encapsulated
	// error if the Result encapsulates an error.
	Error() error
	// IsError returns false if the Result encapsulates a value and true if the
	// Result encapsulates an error.
	IsError() bool
	// MustGet returns the encapsulated value if the Result encapsulates a value
	// and panics if the Result encapsulates an error.
	MustGet() T
	// OrElse returns the encapsulated value if the Result encapsulates a value
	// and returns the given default if the Result encapsulates an error.
	OrElse(d T) T
	// OrElseGet returns the encapsulated value if the Result encapsulates a
	// value and returns the result of calling the given Getter function if the
	// Result encapsulates an error.
	OrElseGet(Getter[T]) T
}

// New returns a Result encapsulating the given value and error. If the error is
// nil the Result will be an Ok encapsulating the value. Otherwise, it will be an
// Error encapsulate the error.
func New[T any](value T, err error) Result[T] {
	if err != nil {
		return Error[T](err)
	}
	return Ok(value)
}

// As encapsulates the value in the given Result in a new Result of the target
// type. If the given Result encapsulates an error then an Error Result will be
// returned of the target type encapsulating the original error.
func As[T, S any](src Result[S]) Result[T] {
	if src.IsError() {
		return Error[T](src.Error())
	}
	var v any = src.MustGet()
	val, ok := v.(T)
	if !ok {
		err := errors.New(`requested interface not implemented by value`)
		return Error[T](err)
	}
	return Ok(val)
}

// ErrorMapper defines a function that takes an error and returns an error.
// This can be used to convert between errors or add additional context to an
// error.
type ErrorMapper func(error) error

// ErrorMap applies the given ErrorMapper to the Result. If the Result
// encapsulates a value then the Result is returned unchanged. If the Result
// encapsulates an error then the error is passed to the ErrorMapper and
// the returned error is propagated to a Result of the return type.
func ErrorMap[T any](f ErrorMapper, r Result[T]) Result[T] {
	if r.IsError() {
		return Error[T](f(r.Error()))
	}
	return r
}

// FromMaybe converts the given maybe.Maybe to a Result. If the maybe.Maybe is
// maybe.Nothing then the Result will have the given error. Otherwise, the Result
// will have the value of the maybe.Just.
func FromMaybe[T any](value maybe.Maybe[T], err error) Result[T] {
	if value.IsJust() {
		return Ok(value.MustGet())
	}
	return Error[T](err)
}
