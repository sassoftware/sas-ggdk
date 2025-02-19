// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sliceutils

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/condition"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/folders"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// AnySliceToSlice returns a typed slice from the specified 'any' slice;
// additionally, returns false when the specified 'any' slice cannot be converted
// into a typed slice, otherwise true.
func AnySliceToSlice[T any](value []any) maybe.Maybe[[]T] {
	if condition.IsNil(value) {
		return maybe.Nothing[[]T]()
	}
	slice := make([]T, len(value))
	for i, each := range value {
		target, ok := each.(T)
		if !ok {
			return maybe.Nothing[[]T]()
		}
		slice[i] = target
	}
	return maybe.Just(slice)
}

// AnyToAnySlice returns an 'any' slice from the specified 'any' value;
// additionally, returns false when the specified 'any' value cannot be converted
// to an 'any' slice, otherwise true.
func AnyToAnySlice(value any) maybe.Maybe[[]any] {
	if condition.IsNil(value) {
		return maybe.Nothing[[]any]()
	}
	slice, ok := value.([]any)
	if !ok {
		return maybe.Nothing[[]any]()
	}
	return maybe.Just(slice)
}

// AssertContains returns an error if the given slice of values does not contain
// the give value, otherwise it returns nil. The returned error describes the
// missing values.
func AssertContains[T comparable](target []T, value T) error {
	if slices.Contains(target, value) {
		return nil
	}
	// We could call stringutils.ToQuoted here but that would create a cycle as
	// stringutils uses sliceutils. I believe sliceutils to be the more
	// foundational so it should reimplement this mapper.Mapper function rather
	// than stringutils reimplementing MapNoError.
	quotedValuesResult := Map(
		func(value T) result.Result[string] {
			return result.Ok(fmt.Sprintf("%v", value))
		},
		target,
	)
	// The mapper above does not ever return an error, so we can call MustGet
	// without checking IsError first.
	quotedValues := quotedValuesResult.MustGet()
	sort.Strings(quotedValues)
	csv := strings.Join(quotedValues, `, `)
	return errors.New(`must be one of %s`, csv)
}

// CollectErrors returns a multi-error for all errors in the given slice of
// result.Result.
func CollectErrors[T any](root error, elements ...result.Result[T]) error {
	errs := FoldNoError(func(errs []error, e result.Result[T]) []error {
		if e.IsError() {
			errs = append(errs, e.Error())
		}
		return errs
	}, []error{}, elements)
	if len(errs) != 0 {
		return errors.WrapAll(errs, root.Error())
	}
	return nil
}

// Detect returns the first value found in the given slice of values that matches
// the given filter. Returns a maybe.Nothing if no value is found. Returns a
// result.Error if the filter fails.
func Detect[T any](
	detect filters.Filter[T],
	values []T,
) result.Result[maybe.Maybe[T]] {
	if detect == nil {
		return result.Ok(maybe.Nothing[T]())
	}
	// This could be implemented as Fold but Fold does not know how to stop
	// iterating. This implementation exits as soon as the first matching element
	// is detected.
	for index, value := range values {
		detected := detect(value)
		if detected.IsError() {
			return result.Error[maybe.Maybe[T]](
				errors.Wrap(
					detected.Error(),
					`detector filter failure; value[%d]=%#v`,
					index,
					value,
				),
			)
		}
		if detected.MustGet() {
			return result.Ok(maybe.Just(value))
		}
	}
	return result.Ok(maybe.Nothing[T]())
}

// DetectNoError returns the first value found in the given slice of values that
// matches the given filter. Returns a maybe.Nothing if no value is found.
func DetectNoError[T any](
	detect filters.FilterNoError[T],
	values []T,
) maybe.Maybe[T] {
	if detect == nil {
		return maybe.Nothing[T]()
	}
	// This could be implemented as Fold but Fold does not know how to stop
	// iterating. This implementation exits as soon as the first matching element
	// is detected.
	for _, value := range values {
		detected := detect(value)
		if detected {
			return maybe.Just(value)
		}
	}
	return maybe.Nothing[T]()
}

// DetectNoErrorResult calls DetectNoError on the slice encapsulated in the
// given result. If the given result encapsulates an error, that error is
// returned in a result of the return type.
func DetectNoErrorResult[T any](
	filter filters.FilterNoError[T],
	values result.Result[[]T],
) result.Result[maybe.Maybe[T]] {
	return result.MapNoError(
		func(values []T) maybe.Maybe[T] {
			return DetectNoError(filter, values)
		},
		values,
	)
}

// DetectResult returns the first value found in the given slice of values that
// matches the given filter. Returns a maybe.Nothing if no value is found.
// Returns a result.Error if the filter fails.
func DetectResult[T any](
	filter filters.Filter[T],
	source result.Result[[]T],
) result.Result[maybe.Maybe[T]] {
	return result.FlatMap2(Detect[T], result.Ok(filter), source)
}

// Disjoint returns the disjoint set of the given two slices; never nil.
func Disjoint[T comparable](left []T, right []T) result.Result[[]T] {
	leftOnly := Select(
		filters.NewSliceContainsFilter(right).Not(),
		left,
	)
	rightOnly := Select(
		filters.NewSliceContainsFilter(left).Not(),
		right,
	)
	return UniqueUnion(leftOnly.MustGet(), rightOnly.MustGet())
}

// FirstError returns the first error, if any, in the given slice of
// result.Result.
func FirstError[T any](elements ...result.Result[T]) error {
	errMaybe := DetectNoError(func(e result.Result[T]) bool {
		return e.IsError()
	}, elements)
	// If a failing result.Result was found, return its error.
	if errMaybe.IsJust() {
		return errMaybe.MustGet().Error()
	}
	return nil
}

// Fold passes the first element of the given slice and the given initial value
// to the given folder. The next element and the value in the result returned
// from the last call to the given folder is passed to the folder until all
// elements have been passed or an error is returned from the folder.
func Fold[T, S any](
	folder folders.Folder[T, S],
	initial T,
	source []S,
) result.Result[T] {
	accumulator := initial
	for _, v := range source {
		res := folder(accumulator, v)
		if res.IsError() {
			return result.Error[T](res.Error())
		}
		accumulator = res.MustGet()
	}
	return result.Ok(accumulator)
}

// FoldNoError passes the first element of the given slice and the given initial
// value to the given folder. The next element and the result of the last call to
// the given folder is passed to the folder until all elements have been passed.
func FoldNoError[T, S any](
	folder folders.FolderNoError[T, S],
	initial T,
	source []S,
) T {
	accumulator := initial
	for _, v := range source {
		accumulator = folder(accumulator, v)
	}
	return accumulator
}

// FoldNoErrorResult calls Fold on the slice encapsulated in the given result. If the
// given result encapsulates an error, that error is returned in a result of the
// return type.
func FoldNoErrorResult[T, S any](
	folder folders.FolderNoError[T, S],
	initial T,
	source result.Result[[]S],
) result.Result[T] {
	return result.MapNoError(
		func(values []S) T {
			return FoldNoError(folder, initial, values)
		},
		source,
	)
}

// FoldResult calls Fold on the slice encapsulated in the given result. If the
// given result encapsulates an error, that error is returned in a result of the
// return type.
func FoldResult[T, S any](
	folder folders.Folder[T, S],
	initial T,
	source result.Result[[]S],
) result.Result[T] {
	return result.FlatMap(
		func(values []S) result.Result[T] {
			return Fold(folder, initial, values)
		},
		source,
	)
}

// Head returns the first element of the given slice.
func Head[T any](source []T) maybe.Maybe[T] {
	if len(source) == 0 {
		return maybe.Nothing[T]()
	}
	return maybe.Just(source[0])
}

// Intersection returns the intersection of the given two slices; never nil.
func Intersection[T comparable](left []T, right []T) result.Result[[]T] {
	return Select(
		filters.NewSliceContainsFilter[T](right),
		left,
	)
}

// LenAll returns the accumulated lengths of each of the given values; never
// negative.
func LenAll[T any](values ...[]T) int {
	lenAll := Fold(lenAllFolder[T], 0, values)
	return lenAll.MustGet()
}

// Map the given source slice, using the given mapper, into a result of a new
// slice. Returns a result of the new slice containing the mapped values, never
// nil.
func Map[T, S any](
	mapper result.FlatMapper[S, T],
	source []S,
) result.Result[[]T] {
	target := make([]T, 0, len(source))
	for _, v := range source {
		mapped := mapper(v)
		if mapped.IsError() {
			return result.Error[[]T](mapped.Error())
		}
		target = append(target, mapped.MustGet())
	}
	return result.Ok(target)
}

// MapperNoError defines a function with one argument of any type that returns a
// single value.
type MapperNoError[T, R any] func(
	T,
) R

// MapNoError the given source slice, using the given mapper, into a new slice.
// Returns the new slice containing the mapped values, never nil.
func MapNoError[T, S any](
	mapper MapperNoError[S, T],
	source []S,
) []T {
	target := make([]T, 0, len(source))
	for _, v := range source {
		value := mapper(v)
		target = append(target, value)
	}
	return target
}

// MapNoErrorResult calls MapNoError on the slice encapsulated in the given
// result. If the given result encapsulates an error, that error is returned in a
// result of the return type.
func MapNoErrorResult[T, S any](
	mapper MapperNoError[S, T],
	source result.Result[[]S],
) result.Result[[]T] {
	return result.FlatMap(
		func(values []S) result.Result[[]T] {
			return result.Ok(MapNoError(mapper, values))
		},
		source,
	)
}

// MapResult calls Map on the slice encapsulated in the given result. If the
// given result encapsulates an error, that error is returned in a result of the
// return type.
func MapResult[T, S any](
	mapper result.FlatMapper[S, T],
	source result.Result[[]S],
) result.Result[[]T] {
	return result.FlatMap(
		func(values []S) result.Result[[]T] {
			return Map(mapper, values)
		},
		source,
	)
}

// Prepend the given values on the target. Returns a new slice, never nil.
func Prepend[T any](target []T, values ...T) []T {
	return slices.Insert(target, 0, values...)
}

// Remove from the given slice the given value; returns a new slice, never nil.
func Remove[T comparable](target []T, values ...T) []T {
	for _, value := range values {
		i := slices.Index(target, value)
		if i == -1 {
			continue
		}
		target = slices.Delete(target, i, i+1)
	}
	return target
}

// Reverse the given slice. Returns a new slice; never nil.
func Reverse[T any](values []T) []T {
	size := len(values)
	res := make([]T, size)
	copy(res, values)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

// Select returns a slice the values containing the values from the given slice
// that match the given filter; never nil.
func Select[T any](
	filter filters.Filter[T],
	source []T,
) result.Result[[]T] {
	target := make([]T, 0, len(source))
	if filter == nil {
		return result.Ok(target)
	}
	return Fold(folders.NewMatchesFolder(filter), target, source)
}

// SelectNoError returns a slice the values containing the values from the given slice
// that match the given filter; never nil.
func SelectNoError[T any](
	filter filters.FilterNoError[T],
	source []T,
) []T {
	target := make([]T, 0, len(source))
	if filter == nil {
		return target
	}
	return FoldNoError(folders.NewMatchesFolderNoError(filter), target, source)
}

// SelectNoErrorResult returns a slice the values containing the values from the
// given slice that match the given filter; never nil.
func SelectNoErrorResult[T any](
	filter filters.FilterNoError[T],
	source result.Result[[]T],
) result.Result[[]T] {
	return result.MapNoError(
		func(values []T) []T {
			return SelectNoError(filter, values)
		},
		source,
	)
}

// SelectResult returns a slice the values containing the values from the given slice
// that match the given filter; never nil.
func SelectResult[T any](
	filter filters.Filter[T],
	source result.Result[[]T],
) result.Result[[]T] {
	return result.FlatMap2(Select[T], result.Ok(filter), source)
}

// SelectUsingSubsetInterfaceFilter mirrors Select, with support for a filter
// function that acts against an interface of the type to be filtered. If the
// type to be filtered is the interface itself, then just use Select.
func SelectUsingSubsetInterfaceFilter[T any, S any](filter filters.Filter[S], source []T) result.Result[[]T] {
	target := make([]T, 0, len(source))
	for _, element := range source {
		asAnS, ok := any(element).(S)
		if !ok {
			err := errors.New("the filter function does not support the element to be filtered")
			return result.Error[[]T](err)
		}
		matchesResult := filter(asAnS)
		if matchesResult.IsError() {
			return result.Error[[]T](matchesResult.Error())
		}
		matches := matchesResult.MustGet()
		if matches {
			target = append(target, element)
		}
	}
	return result.Ok(target)
}

// SelectUsingSubsetInterfaceFilterNoError mirrors Select, with support for a
// filter function that acts against an interface of the type to be filtered. If
// the type to be filtered is the interface itself, then just use Select.
func SelectUsingSubsetInterfaceFilterNoError[T any, S any](
	filter filters.FilterNoError[S],
	source []T,
) result.Result[[]T] {
	target := make([]T, 0, len(source))
	for _, element := range source {
		asAnS, ok := any(element).(S)
		if !ok {
			err := errors.New("the filter function does not support the element to be filtered")
			return result.Error[[]T](err)
		}
		matches := filter(asAnS)
		if matches {
			target = append(target, element)
		}
	}
	return result.Ok(target)
}

// Tail returns all the elements of the given slice except for the first one.
func Tail[T any](source []T) maybe.Maybe[[]T] {
	if len(source) == 0 {
		return maybe.Nothing[[]T]()
	}
	return maybe.Just(source[1:])
}

func ToSlice[T any](values ...T) []T {
	return values
}

// Union returns a slice containing the union of the given slices; never nil.
func Union[T comparable](values ...[]T) result.Result[[]T] {
	size := LenAll(values...)
	accumulator := make([]T, 0, size)
	return Fold(unionFolder[T], accumulator, values)
}

// UniqueUnion returns a slice containing the union of the unique elements of
// the given slices; never nil.
func UniqueUnion[T comparable](values ...[]T) result.Result[[]T] {
	size := LenAll(values...)
	accumulator := make([]T, 0, size)
	return Fold(uniqueUnionFolder[T], accumulator, values)
}

func lenAllFolder[T any](accumulator int, value []T) result.Result[int] {
	return result.Ok(len(value) + accumulator)
}

func unionFolder[T any](accumulator []T, value []T) result.Result[[]T] {
	// This could be a library folder but that would separate it from
	// uniqueUnionFolder.
	accumulator = append(accumulator, value...)
	return result.Ok[[]T](accumulator)
}

func uniqueUnionFolder[T comparable](accumulator []T, value []T) result.Result[[]T] {
	// This could be a library folder but would create a circular dependency on
	// Select.
	uniqueValues := Select(filters.NewSliceContainsFilter(accumulator).Not(), value)
	return result.MapNoError(
		func(values []T) []T { return append(accumulator, values...) },
		uniqueValues,
	)
}
