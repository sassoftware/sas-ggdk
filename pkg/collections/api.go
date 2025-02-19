// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package collections

import (
	"cmp"

	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// Collection defines the API for all collection types, regardless of implementation.
type Collection[T comparable] interface {
	// Add the given items to the receiver. Returns the receiver as a
	// collections.Collection.
	Add(item ...T) Collection[T]
	// Contains returns true if the given item is contained in the receiver,
	// otherwise false.
	Contains(item T) bool
	// Detect queries the receiver for an item matching the given filter. Returns the
	// detected item, and true when found, otherwise the zero value and false;
	// an error is returned if the filter fails.
	Detect(detect filters.Filter[T]) result.Result[maybe.Maybe[T]]
	// Len returns the length of the receiver; an empty receiver has a length of
	// zero.
	Len() int
	// Remove the given item from the receiver. If the item does not exist the
	// receiver is unchanged.Returns the receiver as a collections.Collection.
	Remove(item T) Collection[T]
	// Select returns the receiver's items that match the given filter; an error
	// is returned if the filter fails. The returned collection is the same
	// species as the receiver, and never nil.
	Select(filter filters.Filter[T]) result.Result[Collection[T]]
	// ToSlice returns a new slice of items contained in the receiver. Changes to the
	// slice do not affect the contents of the receiver.
	ToSlice() []T
}

// OrderedCollection defines the API for all constraints.Ordered collection
// types, regardless of implementation. OrderedCollection is composed of the
// Collection interface.
type OrderedCollection[T cmp.Ordered] interface {
	// Collection is a composed type.
	Collection[T]
	// First returns the first item in the receiver.
	First() result.Result[T]
	// Get returns the item in the receiver at the given index.
	Get(index int) result.Result[T]
	// Index returns the index of the given item in the receiver, or -1 if the
	// given item is not contained in the receiver.
	Index(item T) int
	// Insert at the given index the given item.  Returns an error if index is out of
	// bounds.
	Insert(index int, item T) error
	// Largest returns the lexicographical largest item contained in the
	// receiver.
	Largest() result.Result[T]
	// Last returns the last item in the receiver.
	Last() result.Result[T]
	// Set the given item at the given index in the receiver, returning the item
	// that previously occupied the given index, or the zero value, if the index
	// is previously unoccupied.
	Set(index int, item T) result.Result[T]
	// Smallest returns the lexicographical smallest item contained in the receiver.
	Smallest() result.Result[T]
	// SortAscending sorts the content of the receiver in ascending lexical order.
	SortAscending() OrderedCollection[T]
	// SortDescending sorts the content of the receiver in descending lexical order.
	SortDescending() OrderedCollection[T]
}
