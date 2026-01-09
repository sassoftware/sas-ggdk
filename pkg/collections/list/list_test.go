// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package list_test

import (
	"slices"
	"sort"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/collections/list"
	"github.com/sassoftware/sas-ggdk/pkg/collections/set"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_Map(t *testing.T) {
	colors := list.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(value string) result.Result[int] {
		return result.Ok(len(value))
	}
	actual := list.Map[string, int](mapper, colors)
	require.NoError(t, actual.Error())
	expected := list.New[int](colors.Len())
	for _, color := range colors.ToSlice() {
		length := mapper(color)
		require.NoError(t, length.Error())
		expected.Add(length.MustGet())
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_raw_new(t *testing.T) {
	instance := new(list.List[string])
	require.NotNil(t, instance)
	instance.Add(`red`)
	require.True(t, instance.Contains(`red`))
}

func Test_New(t *testing.T) {
	numbers := list.New[int](2)
	size := numbers.Len()
	require.Equal(t, 0, size)
}

func Test_NewFrom(t *testing.T) {
	expected := []int{1, 2, 3}
	numbers := list.NewFrom[int](expected...)
	require.NotNil(t, numbers)
	actualSize := numbers.Len()
	expectedSize := len(expected)
	require.Equal(t, expectedSize, actualSize)
	actual := numbers.ToSlice()
	sort.Ints(actual)
	require.Equal(t, expected, actual)
}

func Test_NewFromCollection(t *testing.T) {
	source := set.NewFrom(1, 2, 3)
	numbers := list.NewFromCollection[int](source)
	require.NotNil(t, numbers)
	actualSize := numbers.Len()
	expectedSize := source.Len()
	require.Equal(t, expectedSize, actualSize)
	actual := numbers.ToSlice()
	sort.Ints(actual)
	expected := source.ToSlice()
	sort.Ints(expected)
	require.Equal(t, expected, actual)
}

func Test_NewWithAccessor(t *testing.T) {
	instance, accessor := list.NewWithAccessor[int](5)
	require.NotNil(t, instance)
	require.NotNil(t, accessor)
	instance.Add(20, 10, 50)
	actual := accessor()
	expected := []int{20, 10, 50} // nolint: prealloc
	require.Equal(t, expected, actual)
	instance.Add(70)
	actual = accessor()
	expected = append(expected, 70)
	require.Equal(t, expected, actual)
	actual = append(actual, 80)
	require.True(t, slices.Contains(actual, 80))
	require.False(t, instance.Contains(80))
}

func Test_Map_fail(t *testing.T) {
	colors := list.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(_ string) result.Result[int] {
		return result.Error[int](errors.New(`fail`))
	}
	actual := list.Map[string, int](mapper, colors)
	require.Error(t, actual.Error())
}

func Test_List_Add_Contains(t *testing.T) {
	numbers := list.New[int](2)
	state := numbers.Contains(10)
	require.False(t, state)
	r := numbers.Add(10)
	require.Same(t, numbers, r)
	size := numbers.Len()
	require.Equal(t, 1, size)
	state = numbers.Contains(10)
	require.True(t, state)
	numbers.Add(10)
	size = numbers.Len()
	require.Equal(t, 2, size)
}

func Test_List_Detect(t *testing.T) {
	numbers := list.NewFrom[int](10, 20, 30, 40)
	filter := func(item int) result.Result[bool] {
		return result.Ok(item > 10 && item < 40)
	}
	actual := numbers.Detect(filter)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	item := actual.MustGet().MustGet()
	state := item == 20 || item == 30
	require.True(t, state)
	numbers.Remove(item)
	actual = numbers.Detect(filter)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	item = actual.MustGet().MustGet()
	state = item == 20 || item == 30
	require.True(t, state)
	numbers.Remove(item)
	actual = numbers.Detect(filter)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_List_Remove_Contains(t *testing.T) {
	numbers := list.New[int](5)
	size := numbers.Len()
	require.Equal(t, 0, size)
	r := numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 0, size)
	require.Same(t, numbers, r)
	numbers.Add(10).
		Add(10).
		Add(20)
	size = numbers.Len()
	require.Equal(t, 3, size)
	state := numbers.Contains(10)
	require.True(t, state)
	numbers.Remove(10)
	state = numbers.Contains(10)
	require.True(t, state)
	size = numbers.Len()
	require.Equal(t, 2, size)
	numbers.Remove(10)
	state = numbers.Contains(10)
	require.False(t, state)
	size = numbers.Len()
	require.Equal(t, 1, size)
	state = numbers.Contains(20)
	require.True(t, state)
	numbers.Remove(20)
	state = numbers.Contains(20)
	require.False(t, state)
	size = numbers.Len()
	require.Equal(t, 0, size)
}

func Test_List_Len(t *testing.T) {
	numbers := list.New[int](5)
	size := numbers.Len()
	require.Equal(t, 0, size)
	numbers.Add(10)
	size = numbers.Len()
	require.Equal(t, 1, size)
	numbers.Add(10)
	size = numbers.Len()
	require.Equal(t, 2, size)
	numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 1, size)
	numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 0, size)
}

func Test_List_Select(t *testing.T) {
	numbers := list.NewFrom[int](1, 2, 3, 4, 5, 2)
	isEvenFilter := func(value int) result.Result[bool] {
		return result.Ok(value%2 == 0)
	}
	evenNumbers := numbers.Select(isEvenFilter)
	require.NoError(t, evenNumbers.Error())
	actual := evenNumbers.MustGet().ToSlice()
	expected := []int{2, 4, 2}
	require.Equal(t, expected, actual)
}

func Test_List_Select_error(t *testing.T) {
	numbers := list.NewFrom[int](1, 2)
	expectedMessage := `failing SELECT filter `
	itemCount := 1
	isEvenFilter := func(_ int) result.Result[bool] {
		if itemCount == 1 {
			itemCount++
			return result.Ok(true) // Ensure we have some data in the target.
		}
		return result.Error[bool](errors.New(expectedMessage))
	}
	evenNumbers := numbers.Select(isEvenFilter)
	require.ErrorContains(t, evenNumbers.Error(), expectedMessage)
}

func Test_List_ToSlice(t *testing.T) {
	numbers := list.New[int](7)
	actual := numbers.ToSlice()
	require.NotNil(t, actual)
	require.Empty(t, actual)
	numbers.Add(1, 2, 3)
	actual = numbers.ToSlice()
	expected := []int{1, 2, 3}
	require.Equal(t, expected, actual)
}

func Test_List_String(t *testing.T) {
	numbers := list.NewFrom[int](1, 2, 3)
	actual := numbers.String()
	expected := `[1 2 3]`
	require.Equal(t, expected, actual)
}

func Test_List_ToCollection(t *testing.T) {
	numbers := list.New[string](0)
	instance := numbers.ToCollection()
	require.NotNil(t, instance)
}
