// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe

// Nothing returns a Maybe encapsulating an absent value.
func Nothing[T any]() Maybe[T] {
	return nothing[T]{}
}

// nothing implements Maybe for an absent value.
type nothing[T any] struct {
}

// IsJust returns false for all absent values.
func (n nothing[T]) IsJust() bool {
	return false
}

// MustGet panics for all absent values.
func (n nothing[T]) MustGet() T {
	panic("nothing does not have a value")
}

// OrElse returns the provided default value for all absent values.
func (n nothing[T]) OrElse(d T) T {
	return d
}

// OrElseGet returns the result of calling the provided function for all absent
// values.
func (n nothing[T]) OrElseGet(f Getter[T]) T {
	return f()
}

// String returns a representation of the value held in this Nothing
func (n nothing[T]) String() string {
	return "{Nothing}"
}
