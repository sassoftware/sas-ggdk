// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package collections

import (
	"cmp"

	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// Add the given items to the collection. Returns the collection as a
// collections.Collection.
func Add[T comparable](c Collection[T], item T) Collection[T] {
	return c.Add(item)
}

// Contains returns true if the given item is contained in the collection,
// otherwise false.
func Contains[T comparable](c Collection[T], item T) bool {
	return c.Contains(item)
}

// Detect queries the collection for an item matching the given filter. Returns
// the detected item, and true when found, otherwise the zero value and false; an
// error is returned if the filter fails.
func Detect[T comparable](c Collection[T], detect filters.Filter[T]) result.Result[maybe.Maybe[T]] {
	return c.Detect(detect)
}

// Len returns the length of the collection; an empty collection has a length of
// zero.
func Len[T comparable](c Collection[T]) int {
	return c.Len()
}

// Remove the given item from the collection. If the item does not exist the
// collection is unchanged.Returns the collection as a collections.Collection.
func Remove[T comparable](c Collection[T], item T) Collection[T] {
	return c.Remove(item)
}

// Select returns the collection's items that match the given filter; an error is
// returned if the filter fails. The returned collection is the same species as
// the collection, and never nil.
func Select[T comparable](c Collection[T], filter filters.Filter[T]) result.Result[Collection[T]] {
	return c.Select(filter)
}

// ToSlice returns a new slice of items contained in the collection. Changes to
// the slice do not affect the contents of the collection.
func ToSlice[T comparable](c Collection[T]) []T {
	return c.ToSlice()
}

// First returns the first item in the collection.
func First[T cmp.Ordered](c OrderedCollection[T]) result.Result[T] {
	return c.First()
}

// Get returns the item in the collection at the given index.
func Get[T cmp.Ordered](c OrderedCollection[T], index int) result.Result[T] {
	return c.Get(index)
}

// Index returns the index of the given item in the collection, or -1 if the
// given item is not contained in the collection.
func Index[T cmp.Ordered](c OrderedCollection[T], item T) int {
	return c.Index(item)
}

// Insert at the given index the given item. Returns the modified collection.
func Insert[T cmp.Ordered](c OrderedCollection[T], index int, item T) result.Result[OrderedCollection[T]] {
	return result.New(c, c.Insert(index, item))
}

// Largest returns the lexicographical largest item contained in the collection.
func Largest[T cmp.Ordered](c OrderedCollection[T]) result.Result[T] {
	return c.Largest()
}

// Last returns the last item in the collection.
func Last[T cmp.Ordered](c OrderedCollection[T]) result.Result[T] {
	return c.Last()
}

// Set the given item at the given index in the collection, returning the item
// that previously occupied the given index, or the zero value, if the index is
// previously unoccupied.
func Set[T cmp.Ordered](c OrderedCollection[T], index int, item T) result.Result[T] {
	return c.Set(index, item)
}

// Smallest returns the lexicographical smallest item contained in the
// collection.
func Smallest[T cmp.Ordered](c OrderedCollection[T]) result.Result[T] {
	return c.Smallest()
}

// SortAscending sorts the content of the collection in ascending lexical order.
func SortAscending[T cmp.Ordered](c OrderedCollection[T]) OrderedCollection[T] {
	return c.SortAscending()
}

// SortDescending sorts the content of the collection in descending lexical
// order.
func SortDescending[T cmp.Ordered](c OrderedCollection[T]) OrderedCollection[T] {
	return c.SortDescending()
}
