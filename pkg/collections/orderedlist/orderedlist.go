// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package orderedlist

import (
	"cmp"
	"slices"
	"sort"

	"github.com/sassoftware/sas-ggdk/pkg/collections"
	"github.com/sassoftware/sas-ggdk/pkg/collections/list"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/folders"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
)

// Map the given OrderedList using the given mapper and filter, returning an
// OrderedList of mapped items, never nil; an error is returned if either the
// mapper or the filter fails.
func Map[S cmp.Ordered, T cmp.Ordered](
	mapper result.FlatMapper[S, T],
	source *OrderedList[S],
) result.Result[*OrderedList[T]] {
	slice := sliceutils.Map(
		mapper,
		source.accessor(),
	)
	return result.MapNoError(newFromSlice[T], slice)
}

// OrderedList defines the structure of an ordered list collection, an
// implementation of the collections.OrderedCollection interface. A List is a
// slice of constraint.Ordered items; duplicates are allowed. This type is not
// thread-safe.
type OrderedList[T cmp.Ordered] struct {
	*list.List[T]
	accessor func() []T
}

// New creates a OrderedList of the given size.
func New[T cmp.Ordered](size int) *OrderedList[T] {
	instance, _ := NewWithAccessor[T](size)
	return instance
}

// NewFrom creates a OrderedList containing the given items.
func NewFrom[T cmp.Ordered](items ...T) *OrderedList[T] {
	size := calculateNewSize(items)
	instance := New[T](size)
	instance.Add(items...)
	return instance
}

// NewFromCollection creates a OrderedList containing the items in the given
// Collection.
func NewFromCollection[T cmp.Ordered](source collections.Collection[T]) *OrderedList[T] {
	slice := source.ToSlice()
	return NewFrom[T](slice...)
}

// NewWithAccessor creates a OrderedList of the given size, along with an
// accessor function. The accessor function is intended for use by composite
// types that need direct access to the OrderedList's implementation; the
// accessor function should not be shared.
func NewWithAccessor[T cmp.Ordered](size int) (*OrderedList[T], func() []T) {
	lst, accessor := list.NewWithAccessor[T](size)
	instance := &OrderedList[T]{
		List:     lst,
		accessor: accessor,
	}
	return instance, instance.accessor
}

// Add the given items to the receiver. Returns the receiver as a
// collections.Collection.
func (odrList *OrderedList[T]) Add(item ...T) collections.Collection[T] {
	lst := odrList.getList()
	lst.Add(item...)
	return odrList
}

// Contains returns true if the given item is contained in the receiver,
// otherwise false.
func (odrList *OrderedList[T]) Contains(item T) bool {
	lst := odrList.getList()
	return lst.Contains(item)
}

// Detect queries the receiver for an item matching the given filter. Returns the
// detected item, and true when found, otherwise the zero value and false;
// an error is returned if the filter fails.
func (odrList *OrderedList[T]) Detect(filter filters.Filter[T]) result.Result[maybe.Maybe[T]] {
	lst := odrList.getList()
	return lst.Detect(filter)
}

// First returns the first item in the receiver; an error is returned if the
// receiver is empty.
func (odrList *OrderedList[T]) First() result.Result[T] {
	return odrList.Get(0)
}

// Get returns the item in the receiver at the given index; an error is returned
// if the index is invalid.
func (odrList *OrderedList[T]) Get(index int) result.Result[T] {
	err := odrList.assertBounds(index)
	if err != nil {
		return result.Error[T](err)
	}
	content := odrList.accessor()
	item := content[index]
	return result.Ok(item)
}

// Index returns the index of the given item in the receiver, or -1 if the given
// item is not contained in the receiver.
func (odrList *OrderedList[T]) Index(item T) int {
	content := odrList.accessor()
	return slices.Index(content, item)
}

// Insert at the given index the given item. Returns an error if index is out of
// bounds.
func (odrList *OrderedList[T]) Insert(index int, item T) error {
	err := odrList.assertBounds(index)
	if err != nil {
		return err
	}
	content := odrList.accessor()
	content = slices.Insert(content, index, item)
	size := calculateNewSize(content)
	odrList.List, odrList.accessor = list.NewWithAccessor[T](size)
	odrList.List.Add(content...)
	return nil
}

// Largest returns the lexicographical largest item contained in the receiver.
func (odrList *OrderedList[T]) Largest() result.Result[T] {
	if odrList.Len() == 0 {
		return result.Error[T](errors.New(`the ordered list is empty`))
	}
	content := odrList.accessor()
	return sliceutils.Fold(folders.LargestFolder[T], content[0], content)
}

// Last returns the last item in the receiver; an error is returned if the
// receiver is empty.
func (odrList *OrderedList[T]) Last() result.Result[T] {
	index := odrList.Len() - 1
	return odrList.Get(index)
}

// Len returns the length of the receiver; an empty receiver has a length of
// zero.
func (odrList *OrderedList[T]) Len() int {
	lst := odrList.getList()
	return lst.Len()
}

// Remove the given item from the receiver. If the item does not exist the
// receiver is unchanged.Returns the receiver as a collections.Collection.
func (odrList *OrderedList[T]) Remove(item T) collections.Collection[T] {
	lst := odrList.getList()
	lst.Remove(item)
	return odrList
}

// Set the given item at the given index in the receiver, returning the item that
// previously occupied the given index, or the zero value, if the index is
// previously unoccupied; an error is returned if the index is invalid.
func (odrList *OrderedList[T]) Set(index int, item T) result.Result[T] {
	previous := odrList.Get(index)
	if !previous.IsError() {
		content := odrList.accessor()
		content[index] = item
	}
	return previous
}

// Smallest returns the lexicographical smallest item contained in the receiver;
// additional, an error is returned if the receiver is empty.
func (odrList *OrderedList[T]) Smallest() result.Result[T] {
	if odrList.Len() == 0 {
		return result.Error[T](errors.New(`the ordered list is empty`))
	}
	content := odrList.accessor()
	return sliceutils.Fold(folders.SmallestFolder[T], content[0], content)
}

// SortAscending sorts the content of the receiver in ascending lexical order.
func (odrList *OrderedList[T]) SortAscending() collections.OrderedCollection[T] {
	content := odrList.accessor()
	comparator := func(i, j int) bool {
		return content[i] < content[j]
	}
	sort.Slice(content, comparator)
	return odrList
}

// SortDescending sorts the content of the receiver in descending lexical order.
func (odrList *OrderedList[T]) SortDescending() collections.OrderedCollection[T] {
	content := odrList.accessor()
	comparator := func(i, j int) bool {
		return content[i] > content[j]
	}
	sort.Slice(content, comparator)
	return odrList
}

// Select returns the receiver's items that match the given filter; an error is
// returned if the filter fails. The returned collection is the same species as
// the receiver, and never nil.
func (odrList *OrderedList[T]) Select(filter filters.Filter[T]) result.Result[collections.Collection[T]] {
	lst := odrList.getList()
	selection := lst.Select(filter)
	return result.MapNoError(newCollectionFromCollection[T], selection)
}

// String returns a human-readable representation of the receiver.
func (odrList *OrderedList[T]) String() string {
	lst := odrList.getList()
	return lst.String()
}

// ToCollection returns the receiver as a collections.Collection.
func (odrList *OrderedList[T]) ToCollection() collections.Collection[T] {
	return odrList
}

// ToOrderedCollection returns the receiver as a collections.OrderedCollection.
func (odrList *OrderedList[T]) ToOrderedCollection() collections.OrderedCollection[T] {
	return odrList
}

// ToSlice returns a new slice of items contained in the receiver. Changes to the
// slice do not affect the contents of the receiver.
func (odrList *OrderedList[T]) ToSlice() []T {
	lst := odrList.getList()
	return lst.ToSlice()
}

func (odrList *OrderedList[T]) assertBounds(index int) error {
	if index >= 0 && index <= odrList.Len()-1 {
		return nil
	}
	return errors.New(`the index %d is out of bounds`, index)
}

func (odrList *OrderedList[T]) getList() *list.List[T] {
	if odrList.List == nil {
		odrList.List, odrList.accessor = list.NewWithAccessor[T](5)
	}
	return odrList.List
}

func calculateNewSize[T any](items []T) int {
	count := float32(len(items))
	return int(count * 1.5)
}

func newFromSlice[T cmp.Ordered](values []T) *OrderedList[T] {
	return NewFrom(values...)
}

func newCollectionFromCollection[T cmp.Ordered](source collections.Collection[T]) collections.Collection[T] {
	return NewFromCollection[T](source)
}
