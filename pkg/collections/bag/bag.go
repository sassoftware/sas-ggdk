// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package bag

import (
	"fmt"

	"github.com/sassoftware/sas-ggdk/pkg/collections"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
)

// Map the given Bag using the given mapper, returning a Bag of mapped
// items, never nil.
func Map[S comparable, T comparable](
	mapper result.FlatMapper[S, T],
	source *Bag[S],
) result.Result[*Bag[T]] {
	slice := sliceutils.Map(
		mapper,
		source.ToSlice(),
	)
	return result.MapNoError(newFromSlice[T], slice)
}

// New creates a Bag of the given size.
func New[T comparable](size int) *Bag[T] {
	return &Bag[T]{
		content: make(map[T]int, size),
	}
}

// NewFrom creates a Bag containing the given items.
func NewFrom[T comparable](items ...T) *Bag[T] {
	size := calculateNewSize(items)
	instance := New[T](size)
	instance.Add(items...)
	return instance
}

// NewFromCollection creates a Bag containing the items in the given Collection.
func NewFromCollection[T comparable](source collections.Collection[T]) *Bag[T] {
	slice := source.ToSlice()
	return NewFrom[T](slice...)
}

// NewWithAccessor creates a Bag of the given size, along with an accessor
// function. The accessor function is intended for use by composite types that
// need direct access to the Bag's implementation; the accessor function should
// not be shared.
func NewWithAccessor[T comparable](size int) (*Bag[T], func() map[T]int) {
	instance := New[T](size)
	return instance, instance.getContent
}

// Bag defines the structure of a bag collection, an implementation of the
// collections.Collection interface. A Bag is a hashed collection of comparable
// items; duplicates are allowed. This type is not thread-safe.
type Bag[T comparable] struct {
	content map[T]int
}

// Add the given items to the receiver. Returns the receiver as a
// collections.Collection.
func (bag *Bag[T]) Add(items ...T) collections.Collection[T] {
	content := bag.getContent()
	for _, item := range items {
		count := content[item]
		content[item] = count + 1
	}
	return bag
}

// Contains returns true if the given item is contained in the receiver,
// otherwise false.
func (bag *Bag[T]) Contains(item T) bool {
	content := bag.getContent()
	_, exists := content[item]
	return exists
}

// Detect queries the receiver for an item matching the given filter. Returns the
// detected item, and true when found, otherwise the zero value and false;
// an error is returned if the filter fails.
func (bag *Bag[T]) Detect(filter filters.Filter[T]) result.Result[maybe.Maybe[T]] {
	slice := bag.ToSlice()
	return sliceutils.Detect(filter, slice)
}

// Len returns the length of the receiver; an empty receiver has a length of
// zero.
func (bag *Bag[T]) Len() int {
	res := 0
	content := bag.getContent()
	for _, count := range content {
		res += count
	}
	return res
}

// Remove the given item from the receiver. If the item does not exist the
// receiver is unchanged.Returns the receiver as a collections.Collection.
func (bag *Bag[T]) Remove(item T) collections.Collection[T] {
	content := bag.getContent()
	count := content[item]
	if count == 0 {
		return bag
	}
	count--
	if count == 0 {
		delete(content, item)
		return bag
	}
	content[item] = count
	return bag
}

// Select returns the receiver's items that match the given filter; an error is
// returned if the filter fails. The returned collection is the same species as
// the receiver, and never nil.
func (bag *Bag[T]) Select(filter filters.Filter[T]) result.Result[collections.Collection[T]] {
	values := bag.ToSlice()
	selection := sliceutils.Select[T](filter, values)
	return result.MapNoError(newCollectionFromSlice[T], selection)
}

// String returns a human-readable representation of the receiver.
func (bag *Bag[T]) String() string {
	slice := bag.ToSlice()
	return fmt.Sprintf(`%v`, slice)
}

// ToCollection returns the receiver as a collections.Collection.
func (bag *Bag[T]) ToCollection() collections.Collection[T] {
	return bag
}

// ToSlice returns a new slice of items contained in the receiver. Changes to the
// slice do not affect the contents of the receiver.
func (bag *Bag[T]) ToSlice() []T {
	size := bag.Len()
	results := make([]T, 0, size)
	content := bag.getContent()
	for key, count := range content {
		for i := 0; i < count; i++ {
			results = append(results, key)
		}
	}
	return results
}

func (bag *Bag[T]) getContent() map[T]int {
	if bag.content == nil {
		bag.content = make(map[T]int, 5)
	}
	return bag.content
}

func calculateNewSize[T any](items []T) int {
	count := float32(len(items))
	return int(count * 2.5)
}

func newFromSlice[T comparable](values []T) *Bag[T] {
	return NewFrom(values...)
}

func newCollectionFromSlice[T comparable](values []T) collections.Collection[T] {
	return NewFrom[T](values...)
}
