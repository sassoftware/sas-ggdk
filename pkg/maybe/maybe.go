// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe

import (
	"io"
)

// Getter defines a function that gets the default value for a Maybe that
// encapsulates an absent value.
type Getter[T any] func() T

// Maybe encapsulates an optional value that may or may not exist.
type Maybe[T any] interface {
	// IsJust returns true if the Maybe encapsulates a value and false if the
	// Maybe encapsulates nothing.
	IsJust() bool
	// MustGet returns the encapsulated value if the Maybe encapsulates a value
	// and panics if the Maybe encapsulates nothing.
	MustGet() T
	// OrElse returns the encapsulated value if the Maybe encapsulates a value
	// and returns the given default if the Maybe encapsulates nothing.
	OrElse(elseValue T) T
	// OrElseGet returns the encapsulated value if the Maybe encapsulates a value
	// and returns the result of calling the given Getter function if the Maybe
	// encapsulates nothing.
	OrElseGet(getter Getter[T]) T
}

// As encapsulates the value in the given Maybe in a new Maybe of the target
// type. If the given Maybe encapsulates nothing, or the encapsulated value does
// not implement the requested type, then a Nothing Maybe will be returned of
// the target type.
func As[T, S any](src Maybe[S]) Maybe[T] {
	if !src.IsJust() {
		return Nothing[T]()
	}
	var v any = src.MustGet()
	val, ok := v.(T)
	if !ok {
		return Nothing[T]()
	}
	return Just(val)
}

// Close calls close on the encapsulated value if that value implements
// io.Closer. It is a no-op if the value encapsulates nothing or if the
// encasuplated value does not implement io.Closer.
func Close[T any](src Maybe[T]) error {
	if !src.IsJust() {
		return nil
	}
	v := any(src.MustGet())
	closer, ok := v.(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}

// MapperNoError defines a function with one argument of any type that returns a
// single value.
type MapperNoError[T, R any] func(
	T,
) R

// Map applies the given MapperNoError to the Maybe. If the Maybe encapsulates a
// value then that value is passed to the Mapper and a Maybe of the return type
// is returned. If the Maybe encapsulates nothing then the nothing is propagated
// to a Maybe of the return type.
func Map[T, R any](f MapperNoError[T, R], m Maybe[T]) Maybe[R] {
	if m.IsJust() {
		return Just(f(m.MustGet()))
	}
	return Nothing[R]()
}

// FlatMapper defines a function with one argument that maps from the argument's
// type to a Maybe of that type.
type FlatMapper[T, R any] func(T) Maybe[R]

// FlatMap applies the given FlatMapper to the Maybe. If the Maybe encapsulates
// a value then that value is passed to the Mapper and the Mapper's return value
// is returned. If the Maybe encapsulates an absent value then Maybe
// encapsulating an absent value of the return type of the Mapper is returned.
func FlatMap[T, R any](f FlatMapper[T, R], m Maybe[T]) Maybe[R] {
	if m.IsJust() {
		return f(m.MustGet())
	}
	return Nothing[R]()
}

// IsJust returns true if the given maybe.Maybe is a maybe.Just. This function
// exists so that it can be used in a result.Map* call.
func IsJust[T any](value Maybe[T]) bool {
	return value.IsJust()
}
