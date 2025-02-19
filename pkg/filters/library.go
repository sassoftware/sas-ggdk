// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters

import (
	"os"
	"slices"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// NewIsDirectoryFilter returns a new Filter function that returns true when the
// given os.DirEntry represents a directory, otherwise false. This filter never
// fails.
func NewIsDirectoryFilter() Filter[os.DirEntry] {
	return func(entry os.DirEntry) result.Result[bool] {
		return result.Ok(entry.IsDir())
	}
}

// NewIsFileFilter returns a filter that returns true when the given os.DirEntry
// represents a file, otherwise false. This filter never fails.
func NewIsFileFilter() Filter[os.DirEntry] {
	return NewIsDirectoryFilter().Not()
}

// NewIsEmptyMapFilter returns a filter that returns true if the given map is
// empty, otherwise false. This filter never fails.
func NewIsEmptyMapFilter[M ~map[K]V, K comparable, V any]() Filter[M] {
	return func(each M) result.Result[bool] {
		return result.Ok(len(each) == 0)
	}
}

// NewIsEmptySliceFilter returns a filter that returns true if the given slice is
// empty, otherwise false. This filter never fails.
func NewIsEmptySliceFilter[T any]() Filter[[]T] {
	return func(each []T) result.Result[bool] {
		return result.Ok(len(each) == 0)
	}
}

// NewIsEmptyStringFilter returns a filter that returns true if the given string
// is empty, otherwise false. This filter never fails.
func NewIsEmptyStringFilter() Filter[string] {
	return func(each string) result.Result[bool] {
		return result.Ok(len(each) == 0)
	}
}

// NewIsEqualFilter returns a filter that returns true if two comparable items
// are equal, otherwise false. This filter never fails.
func NewIsEqualFilter[T comparable](value T) Filter[T] {
	return func(each T) result.Result[bool] {
		return result.Ok(each == value)
	}
}

// NewMapContainsKeyFilter returns a filter that returns true if a map contains a given key.
func NewMapContainsKeyFilter[M ~map[K]V, K comparable, V any](target M) Filter[K] {
	return func(key K) result.Result[bool] {
		_, exists := target[key]
		return result.Ok(exists)
	}
}

// NewSliceContainsFilter returns true if a value is contained in a slice of
// values.
// NOTE: If this filter ever returns an error then stringutils.Disjoint will
// panic.
func NewSliceContainsFilter[T comparable](values []T) Filter[T] {
	return func(value T) result.Result[bool] {
		return result.Ok(slices.Contains(values, value))
	}
}
