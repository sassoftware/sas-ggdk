// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stack

import (
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
)

// New creates a stack of the given size.
func New[T any](size int) *Stack[T] {
	return &Stack[T]{
		content: make([]T, 0, size),
	}
}

// Stack models a stack of values that can be pushed, peeked, and popped. This
// type is not thread-safe.
type Stack[T any] struct {
	content []T
}

// Peek returns the value of the item on the top of the stack; additionally, the
// zero value and an error is returned if the stack is empty.
func (stack *Stack[T]) Peek() (item T, err error) {
	err = stack.ensureNotEmpty()
	if err != nil {
		return
	}
	index := len(stack.content) - 1
	item = stack.content[index]
	return
}

// Pop removes and returns the item on the top of the stack; additionally, the
// zero value and an error is returned if the stack is empty.
func (stack *Stack[T]) Pop() (item T, err error) {
	err = stack.ensureNotEmpty()
	if err != nil {
		return
	}
	index := len(stack.content) - 1
	item = stack.content[index]
	stack.content = stack.content[0:index]
	return
}

// Push adds the given item to the top of the stack; the stack is returned; never
// nil
func (stack *Stack[T]) Push(item T) *Stack[T] {
	stack.content = append(stack.content, item)
	return stack
}

// Size returns the size of the stack; never negative.
func (stack *Stack[T]) Size() int {
	return len(stack.content)
}

// ToSlice returns a slice containing the content of the stack. The first item in
// the slice is the top of the stack, and the last time in the slice is the
// bottom of the stack. Changes to the returned slice do not change the state of
// the stack. Returns an empty slice if the stack is empty; never nil.
func (stack *Stack[T]) ToSlice() []T {
	return sliceutils.Reverse(stack.content)
}

func (stack *Stack[T]) ensureNotEmpty() error {
	if stack.Size() != 0 {
		return nil
	}
	return errors.New(`%T is empty`, stack)
}
