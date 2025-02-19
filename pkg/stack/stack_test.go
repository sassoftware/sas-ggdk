// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stack_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/stack"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	instance := stack.New[string](0)
	require.NotNil(t, instance)
}

func Test_Peek(t *testing.T) {
	instance := stack.New[string](2)
	actual, err := instance.Peek()
	require.Error(t, err)
	require.Empty(t, actual)
	expected := `one`
	instance.Push(expected)
	actual, err = instance.Peek()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	expected = `two`
	instance.Push(expected)
	actual, err = instance.Peek()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func Test_Push(t *testing.T) {
	instance := stack.New[string](2)
	expected := `one`
	result := instance.Push(expected)
	require.Same(t, instance, result)
	expectedSize := 1
	actualSize := instance.Size()
	require.Equal(t, expectedSize, actualSize)
	actual, err := instance.Peek()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	expected = `two`
	result = instance.Push(expected)
	require.Same(t, instance, result)
	expectedSize = 2
	actualSize = instance.Size()
	require.Equal(t, expectedSize, actualSize)
	actual, err = instance.Peek()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func Test_Pop(t *testing.T) {
	instance := stack.New[int](3)
	instance.Push(1).
		Push(2).
		Push(3)
	expected := []int{
		3, 2, 1,
	}
	for i := 0; instance.Size() != 0; i++ {
		actual, err := instance.Pop()
		require.NoError(t, err)
		require.Equal(t, expected[i], actual)
	}
	size := instance.Size()
	require.Equal(t, 0, size)
	actual, err := instance.Pop()
	require.Error(t, err)
	require.Equal(t, 0, actual)
}

func Test_ToSlice(t *testing.T) {
	instance := stack.New[int](3)
	actual := instance.ToSlice()
	require.NotNil(t, actual)
	require.Empty(t, actual)
	instance.Push(1).
		Push(2).
		Push(3)
	actual = instance.ToSlice()
	expectedSize := instance.Size()
	actualSize := len(actual)
	require.Equal(t, expectedSize, actualSize)
	expected := []int{3, 2, 1}
	require.Equal(t, expected, actual)
}
