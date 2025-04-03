// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe

import "fmt"

// Just returns a Maybe encapsulating a value.
func Just[T any](value T) Maybe[T] {
	return &just[T]{value: value}
}

// just implements Maybe for values.
type just[T any] struct {
	value T
}

// IsJust returns true for all values.
func (j *just[T]) IsJust() bool {
	return true
}

// MustGet returns the encapsulated value.
func (j *just[T]) MustGet() T {
	return j.get()
}

// OrElse returns the encapsulated value. The default value is unused.
func (j *just[T]) OrElse(_ T) T {
	return j.get()
}

// OrElseGet returns the encapsulated value. The default value getter is unused.
func (j *just[T]) OrElseGet(_ Getter[T]) T {
	return j.get()
}

// get returns the encapsulated value. This is not a Maybe interface method.
func (j *just[T]) get() T {
	return j.value
}

// String returns a representation of the value held in this Just.
func (j *just[T]) String() string {
	return fmt.Sprintf("{Just: %#v}", j.value)
}
