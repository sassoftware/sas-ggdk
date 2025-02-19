// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters

import (
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// Filter defines a function type for filtering values. The function returns true
// when the value matches the filter, otherwise false; additionally, an error is
// returned when the filter fails.
type Filter[T any] func(value T) result.Result[bool]

// And creates a Filter that ANDs the receiver and the given filter.
func (thisFilter Filter[T]) And(thatFilter Filter[T]) Filter[T] {
	return func(value T) result.Result[bool] {
		match := ApplyFilter(thisFilter, value)
		match = result.FlatMap(func(b bool) result.Result[bool] {
			if !b {
				return result.Ok(b)
			}
			return ApplyFilter(thatFilter, value)
		}, match)
		context := func(err error) error { return errors.Wrap(err, `filter AND failure`) }
		return result.ErrorMap[bool](context, match)
	}
}

// Not creates a Filter that negates the receiver.
func (thisFilter Filter[T]) Not() Filter[T] {
	return func(value T) result.Result[bool] {
		match := ApplyFilter(thisFilter, value)
		match = result.MapNoError(func(value bool) bool {
			return !value
		}, match)
		context := func(err error) error { return errors.Wrap(err, `filter NOT failure`) }
		return result.ErrorMap[bool](context, match)
	}
}

// Or creates a Filter that ORs the receiver and the given filter.
func (thisFilter Filter[T]) Or(thatFilter Filter[T]) Filter[T] {
	return func(value T) result.Result[bool] {
		match := ApplyFilter(thisFilter, value)
		match = result.FlatMap(func(b bool) result.Result[bool] {
			if b {
				return result.Ok(b)
			}
			return ApplyFilter(thatFilter, value)
		}, match)
		context := func(err error) error { return errors.Wrap(err, `filter OR failure`) }
		return result.ErrorMap[bool](context, match)
	}
}

// ApplyFilter is a helper function that evaluates the given filter and returns
// the result; when the given filter is nil true is returned.
func ApplyFilter[T any](filter Filter[T], value T) result.Result[bool] {
	if filter == nil {
		return result.Ok[bool](true)
	}
	return filter(value)
}

// MatchAll is a filter that always returns true.
func MatchAll[T any](_ T) result.Result[bool] {
	return result.Ok(true)
}

// MatchNone is a filter that always returns false.
func MatchNone[T any](_ T) result.Result[bool] {
	return result.Ok(false)
}
