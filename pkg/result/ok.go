// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

import "fmt"

// Ok returns a Result encapsulating the given value.
func Ok[T any](value T) Result[T] {
	return &ok[T]{value: value}
}

// ok implements Result for values.
type ok[T any] struct {
	value T
}

// Error returns nil for all values.
func (o *ok[T]) Error() error {
	return nil
}

// IsError returns false for all values.
func (o *ok[T]) IsError() bool {
	return false
}

// MustGet returns the encapsulated value.
func (o *ok[T]) MustGet() T {
	return o.get()
}

// OrElse returns the encapsulated value. The default value is unused.
func (o *ok[T]) OrElse(_ T) T {
	return o.get()
}

// OrElseGet returns the encapsulated value. The default value getter is unused.
func (o *ok[T]) OrElseGet(_ Getter[T]) T {
	return o.get()
}

// get returns the encapsulated value. This is not a Result interface method.
func (o *ok[T]) get() T {
	return o.value
}

func (o *ok[T]) String() string {
	return fmt.Sprintf("{Ok: %#v}", o.value)
}
