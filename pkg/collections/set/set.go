// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package set

import (
	"github.com/sassoftware/sas-ggdk/pkg/collections"
	"github.com/sassoftware/sas-ggdk/pkg/collections/bag"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
)

// Map the given Set using the given mapper and filter, returning a Set of mapped
// items, never nil; an error is returned if either the mapper or the filter
// fails.
func Map[S comparable, T comparable](
	mapper result.FlatMapper[S, T],
	source *Set[S],
) result.Result[*Set[T]] {
	slice := sliceutils.Map(
		mapper,
		source.ToSlice(),
	)
	return result.MapNoError(newFromSlice[T], slice)
}

// New creates a Set of the given size.
func New[T comparable](size int) *Set[T] {
	bg, accessor := bag.NewWithAccessor[T](size)
	return &Set[T]{
		Bag:      bg,
		accessor: accessor,
	}
}

// NewFrom creates a Set containing the given items.
func NewFrom[T comparable](items ...T) *Set[T] {
	size := calculateNewSize(items)
	instance := New[T](size)
	instance.Add(items...)
	return instance
}

// NewFromCollection creates a Set containing the items in the given Collection.
func NewFromCollection[T comparable](source collections.Collection[T]) *Set[T] {
	slice := source.ToSlice()
	return NewFrom[T](slice...)
}

// NewWithAccessor creates a Set of the given size, along with an accessor
// function. The accessor function is intended for use by composite types that
// need direct access to the Set's implementation; the accessor function should
// not be shared.
func NewWithAccessor[T comparable](size int) (*Set[T], func() map[T]int) {
	bg, accessor := bag.NewWithAccessor[T](size)
	instance := &Set[T]{
		Bag:      bg,
		accessor: accessor,
	}
	return instance, instance.accessor
}

// Set defines the structure of a set collection, an implementation of the
// collections.Collection interface. A Set is a hashed collection of comparable
// items; duplicates are not allowed. This type is not thread-safe.
type Set[T comparable] struct {
	*bag.Bag[T]
	accessor func() map[T]int
}

// Add the given items to the receiver. Returns the receiver as a
// collections.Collection.
// collections.Collection.
func (set *Set[T]) Add(items ...T) collections.Collection[T] {
	bg := set.getBag()
	for _, item := range items {
		exists := set.Contains(item)
		if !exists {
			bg.Add(item) // Call contained Bag's Add function.
		}
	}
	return set
}

// Contains returns true if the given item is contained in the receiver,
// otherwise false.
func (set *Set[T]) Contains(item T) bool {
	bg := set.getBag()
	return bg.Contains(item)
}

// Detect queries the receiver for an item matching the given filter. Returns the
// detected item, and true when found, otherwise the zero value and false;
// an error is returned if the filter fails.
func (set *Set[T]) Detect(filter filters.Filter[T]) result.Result[maybe.Maybe[T]] {
	bg := set.getBag()
	return bg.Detect(filter)
}

// Len returns the length of the receiver; an empty receiver has a length of
// zero.
func (set *Set[T]) Len() int {
	bg := set.getBag()
	return bg.Len()
}

// Remove the given item from the receiver. If the item does not exist the
// receiver is unchanged.Returns the receiver as a collections.Collection.
func (set *Set[T]) Remove(item T) collections.Collection[T] {
	bg := set.getBag()
	bg.Remove(item)
	return set
}

// Select returns the receiver's items that match the given filter; an error is
// returned if the filter fails. The returned collection is the same species as
// the receiver, and never nil.
func (set *Set[T]) Select(filter filters.Filter[T]) result.Result[collections.Collection[T]] {
	bg := set.getBag()
	selection := bg.Select(filter)
	return result.MapNoError(newCollectionFromCollection[T], selection)
}

// String returns a human-readable representation of the receiver.
func (set *Set[T]) String() string {
	bg := set.getBag()
	return bg.String()
}

// ToCollection returns the receiver as a collections.Collection.
func (set *Set[T]) ToCollection() collections.Collection[T] {
	return set
}

// ToSlice returns a new slice of items contained in the receiver. Changes to the
// slice do not affect the contents of the receiver.
func (set *Set[T]) ToSlice() []T {
	bg := set.getBag()
	return bg.ToSlice()
}

func (set *Set[T]) getBag() *bag.Bag[T] {
	if set.Bag == nil {
		set.Bag, set.accessor = bag.NewWithAccessor[T](5)
	}
	return set.Bag
}

func calculateNewSize[T any](items []T) int {
	count := float32(len(items))
	return int(count * 2.5)
}

func newFromSlice[T comparable](values []T) *Set[T] {
	return NewFrom(values...)
}

func newCollectionFromCollection[T comparable](source collections.Collection[T]) collections.Collection[T] {
	return NewFromCollection[T](source)
}
