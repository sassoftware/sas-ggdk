// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters

// FilterNoError defines a function type for filtering values. The function
// returns true when the value matches the filter, otherwise false.
type FilterNoError[T any] func(value T) bool

// And creates a FilterNoError that ANDs the receiver and the given filter.
func (thisFilter FilterNoError[T]) And(thatFilter FilterNoError[T]) FilterNoError[T] {
	return func(value T) bool {
		return ApplyFilterNoError(thisFilter, value) &&
			ApplyFilterNoError(thatFilter, value)
	}
}

// Not creates a FilterNoError that negates the receiver.
func (thisFilter FilterNoError[T]) Not() FilterNoError[T] {
	return func(value T) bool {
		return !ApplyFilterNoError(thisFilter, value)
	}
}

// Or creates a FilterNoError that ORs the receiver and the given filter.
func (thisFilter FilterNoError[T]) Or(thatFilter FilterNoError[T]) FilterNoError[T] {
	return func(value T) bool {
		return ApplyFilterNoError(thisFilter, value) ||
			ApplyFilterNoError(thatFilter, value)
	}
}

// ApplyFilterNoError is a helper function that evaluates the given filter and
// returns the result; when the given filter is nil true is returned.
func ApplyFilterNoError[T any](filter FilterNoError[T], value T) bool {
	if filter == nil {
		return true
	}
	return filter(value)
}

// MatchAllNoError is a filter that always returns true.
func MatchAllNoError[T any](_ T) bool {
	return true
}

// MatchNoneNoError is a filter that always returns false.
func MatchNoneNoError[T any](_ T) bool {
	return false
}
