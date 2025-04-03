// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maputils

import (
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/condition"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/folders"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/sassoftware/sas-ggdk/pkg/stringutils"
)

// AssertValuesNotNil returns an error if the given map contains a key with a
// nil value, otherwise nil.
func AssertValuesNotNil[M ~map[K]V, K comparable, V any](source M) error {
	ok, missing := valuesNotNil(source)
	if ok {
		return nil
	}
	quotedMissing := stringutils.ToQuoted(missing...)
	sort.Strings(quotedMissing)
	csv := strings.Join(quotedMissing, `, `)
	message := `missing required fields: ` + csv
	return errors.New(message)
}

// DeleteKeys deletes the given keys from the given map.
func DeleteKeys[M ~map[K]V, K comparable, V any](source M, keys ...K) {
	for _, key := range keys {
		delete(source, key)
	}
}

// Fold passes a slice of the keys of the given map and the given initial value
// and folder to sliceutils.Fold.
func Fold[M ~map[K]V, K comparable, V, T any](
	folder folders.Folder[T, K],
	initial T,
	source M,
) result.Result[T] {
	keys := make([]K, 0, len(source))
	keys = slices.AppendSeq(keys, maps.Keys(source))
	return sliceutils.Fold(folder, initial, keys)
}

// FoldResult calls Fold on the map encapsulated in the given result. If the
// given result encapsulates an error, that error is returned in a result of the
// return type.
func FoldResult[M ~map[K]V, K comparable, V, T any](
	folder folders.Folder[T, K],
	initial T,
	source result.Result[M],
) result.Result[T] {
	return result.FlatMap(
		func(m M) result.Result[T] {
			return Fold(folder, initial, m)
		},
		source,
	)
}

// Get calls the two value map index and converts to a Maybe.
func Get[M ~map[K]V, K comparable, V any](source M, key K) maybe.Maybe[V] {
	v, ok := source[key]
	if ok {
		return maybe.Just(v)
	}
	return maybe.Nothing[V]()
}

// LenAll returns the accumulated lengths of the values of the given maps.
func LenAll[M ~map[K]V, K comparable, V any](sources ...M) int {
	size := 0
	for _, source := range sources {
		size += len(source)
	}
	return size
}

// Map the given source map of type, using the given mapper, into a new map. The keys of
// the source map are used to store the mapped values in the new map. If a filter
// is specified, only values whose keys match the filter are mapped into the new
// map. Returns the new map containing the mapped values, never nil;
// additionally, an error is returned if either the mapper or the filter fails.
func Map[M ~map[K]V, K comparable, V, T any](
	mapper result.FlatMapper[K, T],
	source M,
) result.Result[map[K]T] {
	folder := func(accumulator map[K]T, key K) result.Result[map[K]T] {
		value := mapper(key)
		if value.IsError() {
			return result.Error[map[K]T](value.Error())
		}
		accumulator[key] = value.MustGet()
		return result.Ok(accumulator)
	}
	keys := make([]K, 0, len(source))
	keys = slices.AppendSeq(keys, maps.Keys(source))
	return sliceutils.Fold(
		folder,
		make(map[K]T, len(source)),
		keys,
	)
}

// MapResult calls Map on the map encapsulated in the given result. If the given
// result encapsulates an error, that error is returned in a result of the return
// type.
func MapResult[M ~map[K]V, K comparable, V, T any](
	source result.Result[M],
	mapper result.FlatMapper[K, T],
) result.Result[map[K]T] {
	return result.FlatMap(
		func(m M) result.Result[map[K]T] {
			return Map(mapper, m)
		},
		source,
	)
}

// Merge returns a map containing all the keys and values of the given maps; when duplicate
// keys exist the value of the "last-man-in" wins; never nil.
func Merge[M ~map[K]V, K comparable, V any](sources ...M) M {
	folder := func(accumulator M, source M) result.Result[M] {
		for key, value := range source {
			accumulator[key] = value
		}
		return result.Ok(accumulator)
	}
	size := LenAll(sources...)
	target := sliceutils.Fold(folder, make(M, size), sources)
	return target.MustGet()
}

// Select returns a map containing the values for the keys from the given map
// that match the given filter; never nil.
func Select[M ~map[K]V, K comparable, V any](
	source M,
	filter filters.Filter[K],
) result.Result[M] {
	if filter == nil {
		return result.Ok(make(M, 0))
	}
	folder := func(accumulator M, key K) result.Result[M] {
		match := filter(key)
		if match.IsError() {
			return result.Error[M](match.Error())
		}
		if match.MustGet() {
			accumulator[key] = source[key]
		}
		return result.Ok(accumulator)
	}
	keys := make([]K, 0, len(source))
	keys = slices.AppendSeq(keys, maps.Keys(source))
	return sliceutils.Fold(folder, make(M, len(source)), keys)
}

// ToAnyMap attempts to type-convert the given value to an "any map", otherwise
// returns a "nil map".
func ToAnyMap(value any) map[any]any {
	if condition.IsNil(value) {
		return nil
	}
	m, ok := value.(map[any]any)
	if !ok {
		return nil
	}
	return m
}

// ToBool returns the bool value for the given key in the given map. If the
// value is not a bool an error is returned.
func ToBool[M ~map[string]any](
	source M,
	key string,
) (bool, error) {
	value, ok := source[key]
	if !ok {
		return false, errors.New(`unable to read the "%s" key`, key)
	}
	b, ok := value.(bool)
	if !ok {
		return false, errors.New(`unable to parse the "%s" key's value "%s"`, key, value)
	}
	return b, nil
}

// valuesNotNil returns true (and nil) if the given map does not contain a key
// with a nil value, otherwise false (and the keys that had nil values).
func valuesNotNil[M ~map[K]V, K comparable, V any](source M) (bool, []K) {
	nilKeys := make([]K, 0, len(source))
	for key, value := range source {
		if condition.IsNil(value) {
			nilKeys = append(nilKeys, key)
		}
	}
	if len(nilKeys) == 0 {
		return true, nil
	}
	return false, nilKeys
}
