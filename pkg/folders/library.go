// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package folders

import (
	"cmp"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// NewCountingFolder returns a new Folder function that increments the
// accumulator value by the given step amount for each element folded over. The
// elements are ignored.
func NewCountingFolder[S any](step int) Folder[int, S] {
	return func(accumulator int, _ S) result.Result[int] {
		return result.Ok(accumulator + step)
	}
}

// NewSliceFolder returns a new Folder function that appends each element folded
// over to the accumulator slice. This is similar to calling append(accumulator,
// slice...) except that this folder can be used in a call to sliceutils.Fold*.
func NewSliceFolder[S any]() Folder[[]S, S] {
	return func(accumulator []S, source S) result.Result[[]S] {
		return result.Ok(append(accumulator, source))
	}
}

// NewMapFolder returns a new Folder function that takes two functions, applies
// the toKey function and toValue function to the element being folded over and
// sets the resulting key to be the resulting value in the provided accumulator
// map.
func NewMapFolder[S any, K comparable, V any](
	toKey ToKey[S, K],
	toValue ToValue[S, V],
) Folder[map[K]V, S] {
	return func(accumulator map[K]V, source S) result.Result[map[K]V] {
		key := toKey(source)
		if key.IsError() {
			return result.Error[map[K]V](errors.Wrap(key.Error(), `transformer "to key" failure`))
		}
		value := toValue(source)
		if value.IsError() {
			return result.Error[map[K]V](errors.Wrap(value.Error(), `transformer "to value" failure`))
		}
		accumulator[key.MustGet()] = value.MustGet()
		return result.Ok(accumulator)
	}
}

// LargestFolder returns the largest of a given value and an accumulator value.
func LargestFolder[T cmp.Ordered](accumulator T, value T) result.Result[T] {
	if value > accumulator {
		return result.Ok(value)
	}
	return result.Ok(accumulator)
}

// SmallestFolder returns the smallest of a given value and an accumulator value.
func SmallestFolder[T cmp.Ordered](accumulator T, value T) result.Result[T] {
	if value < accumulator {
		return result.Ok(value)
	}
	return result.Ok(accumulator)
}

// NewMatchesFolder returns a new Folder function that applies a Filter function to a
// given value and appends it to a provided accumulator slice for a match.
func NewMatchesFolder[T any](filter filters.Filter[T]) Folder[[]T, T] {
	return func(accumulator []T, value T) result.Result[[]T] {
		match := filter(value)
		if match.IsError() {
			return result.Error[[]T](match.Error())
		}
		if match.MustGet() {
			accumulator = append(accumulator, value)
		}
		return result.Ok(accumulator)
	}
}

// NewMatchesFolderNoError returns a new Folder function that applies a Filter function to a
// given value and appends it to a provided accumulator slice for a match.
func NewMatchesFolderNoError[T any](filter filters.FilterNoError[T]) FolderNoError[[]T, T] {
	return func(accumulator []T, value T) []T {
		match := filter(value)
		if match {
			accumulator = append(accumulator, value)
		}
		return accumulator
	}
}
