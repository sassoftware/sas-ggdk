// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"slices"

	"github.com/sassoftware/sas-ggdk/pkg/collections"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
)

// Map the given List using the given mapper and filter, returning a List of
// mapped items, never nil; an error is returned if either the mapper or the
// filter fails.
func Map[S comparable, T comparable](
	mapper result.FlatMapper[S, T],
	source *List[S],
) result.Result[*List[T]] {
	slice := sliceutils.Map(
		mapper,
		source.getContent(),
	)
	return result.MapNoError(newFromSlice[T], slice)
}

// New creates a List of the given size.
func New[T comparable](size int) *List[T] {
	return &List[T]{
		content: make([]T, 0, size),
	}
}

// NewFrom creates a List containing the given items.
func NewFrom[T comparable](items ...T) *List[T] {
	size := calculateNewSize(items)
	instance := New[T](size)
	instance.Add(items...)
	return instance
}

// NewFromCollection creates a List containing the items in the given Collection.
func NewFromCollection[T comparable](source collections.Collection[T]) *List[T] {
	slice := source.ToSlice()
	return NewFrom[T](slice...)
}

// NewWithAccessor creates a List of the given size, along with an accessor
// function. The accessor function is intended for use by composite types that
// need direct access to the List's implementation; the accessor function should
// not be shared.
func NewWithAccessor[T comparable](size int) (*List[T], func() []T) {
	instance := New[T](size)
	return instance, instance.getContent
}

// List defines the structure of a list collection, an implementation of the
// collections.Collection interface. A List is a slice of comparable items;
// duplicates are allowed. This type is not thread-safe.
type List[T comparable] struct {
	content []T
}

// Add the given items to the receiver. Returns the receiver as a
// collections.Collection.
func (list *List[T]) Add(item ...T) collections.Collection[T] {
	content := list.getContent()
	content = append(content, item...)
	list.content = content
	return list
}

// Contains returns true if the given item is contained in the receiver,
// otherwise false.
func (list *List[T]) Contains(item T) bool {
	content := list.getContent()
	return slices.Contains(content, item)
}

// Detect queries the receiver for an item matching the given filter. Returns the
// detected item, and true when found, otherwise the zero value and false; an
// error is returned if the filter fails.
func (list *List[T]) Detect(filter filters.Filter[T]) result.Result[maybe.Maybe[T]] {
	content := list.getContent()
	return sliceutils.Detect(filter, content)
}

// Len returns the length of the receiver; an empty receiver has a length of
// zero.
func (list *List[T]) Len() int {
	content := list.getContent()
	return len(content)
}

// Remove the given item from the receiver. If the item does not exist the
// receiver is unchanged.Returns the receiver as a collections.Collection.
func (list *List[T]) Remove(item T) collections.Collection[T] {
	content := list.getContent()
	list.content = sliceutils.Remove(content, item)
	return list
}

// Select returns the receiver's items that match the given filter; an error is
// returned if the filter fails. The returned collection is the same species as
// the receiver, and never nil.
func (list *List[T]) Select(filter filters.Filter[T]) result.Result[collections.Collection[T]] {
	content := list.getContent()
	selection := sliceutils.Select[T](filter, content)
	return result.MapNoError(newCollectionFromSlice[T], selection)
}

// String returns a human-readable representation of the receiver.
func (list *List[T]) String() string {
	content := list.getContent()
	return fmt.Sprintf(`%v`, content)
}

// ToCollection returns the receiver as a collections.Collection.
func (list *List[T]) ToCollection() collections.Collection[T] {
	return list
}

// ToSlice returns a new slice of items contained in the receiver. Changes to the
// slice do not affect the contents of the receiver.
func (list *List[T]) ToSlice() []T {
	content := list.getContent()
	dest := make([]T, len(content))
	copy(dest, content)
	return dest
}

func (list *List[T]) getContent() []T {
	if list.content == nil {
		list.content = make([]T, 5)
	}
	return list.content
}

func calculateNewSize[T any](items []T) int {
	count := float32(len(items))
	return int(count * 1.5)
}

func newFromSlice[T comparable](values []T) *List[T] {
	return NewFrom(values...)
}

func newCollectionFromSlice[T comparable](values []T) collections.Collection[T] {
	return NewFrom(values...)
}
