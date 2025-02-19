// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lazy

// Creator defines a function provided by a caller that returns any type that
// will be used to create a Memoized return value.
type Creator[T any] func() T

// Disposer defines a function that takes any type and returns any type that
// will be used to dispose of a Memoized value.
type Disposer[T, R any] func(t T) R

// Disposal defines a function that is returned to the caller from the Memoize
// creation functions. When this function is called it will call the Disposer
// provided by the user if, and only if, the Memoized value has been created and
// this function has not been called before.
type Disposal[R any] func() R

// MakeGetter returns a function that calls the provided Creator the first time it
// is called. All subsequent calls to the MakeGetter return the initial value.
func MakeGetter[T any](create Creator[T]) Creator[T] {
	var value T
	created := false
	return func() T {
		if !created {
			created = true
			value = create() // Save T
		}
		return value
	}
}

// MakeGetterWithDispose returns a function that calls the provided Creator the
// first time it is called. All subsequent calls to the Creator return
// the initial value.
func MakeGetterWithDispose[T, R any](
	create Creator[T],
	dispose Disposer[T, R],
) (Creator[T], Disposal[R]) {
	var value T
	var disposeResult R
	created := false
	disposed := false
	creationFunc := func() T {
		if !created {
			created = true
			value = create() // Save T
		}
		return value
	}
	disposalFunc := func() R {
		if created && !disposed {
			disposed = true
			disposeResult = dispose(value) // Dispose of saved T
		}
		return disposeResult
	}
	return creationFunc, disposalFunc
}
