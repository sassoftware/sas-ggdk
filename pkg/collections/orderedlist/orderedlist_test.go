// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package orderedlist_test

import (
	"slices"
	"sort"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/collections/orderedlist"
	"github.com/sassoftware/sas-ggdk/pkg/collections/set"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/stretchr/testify/require"
)

func Test_Map(t *testing.T) {
	colors := orderedlist.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(value string) result.Result[int] {
		return result.Ok(len(value))
	}
	collection := orderedlist.Map[string, int](mapper, colors)
	require.NoError(t, collection.Error())
	actual := collection.MustGet().ToSlice()
	expected := make([]int, 0, colors.Len())
	for _, color := range colors.ToSlice() {
		length := mapper(color)
		require.NoError(t, length.Error())
		expected = append(expected, length.MustGet())
	}
	require.Equal(t, expected, actual)
}

func Test_Map_fail(t *testing.T) {
	colors := orderedlist.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(value string) result.Result[int] {
		return result.Error[int](errors.New(`fail`))
	}
	actual := orderedlist.Map[string, int](mapper, colors)
	require.Error(t, actual.Error())
}

func Test_raw_new(t *testing.T) {
	instance := new(orderedlist.OrderedList[string])
	require.NotNil(t, instance)
	expected := `red`
	instance.Add(expected)
	actual := instance.Get(0)
	require.NoError(t, actual.Error())
	require.Equal(t, expected, actual.MustGet())
}

func Test_NewFromCollection(t *testing.T) {
	source := set.NewFrom(1, 2, 3)
	numbers := orderedlist.NewFromCollection[int](source)
	require.NotNil(t, numbers)
	actualSize := numbers.Len()
	expectedSize := source.Len()
	require.Equal(t, expectedSize, actualSize)
	numbers.SortAscending()
	actual := numbers.ToSlice()
	expected := source.ToSlice()
	sort.Ints(expected)
	require.Equal(t, expected, actual)
}

func Test_NewWithAccessor(t *testing.T) {
	instance, accessor := orderedlist.NewWithAccessor[int](5)
	require.NotNil(t, instance)
	require.NotNil(t, accessor)
	instance.Add(20, 10, 50)
	actual := accessor()
	expected := []int{20, 10, 50}
	require.Equal(t, expected, actual)
	instance.Add(70)
	actual = accessor()
	expected = append(expected, 70)
	require.Equal(t, expected, actual)
	actual = append(actual, 80)
	require.True(t, slices.Contains(actual, 80))
	require.False(t, instance.Contains(80))
}

func Test_OrderedList_Add(t *testing.T) {
	instance := orderedlist.New[int](10)
	instance.Add(10).Add(20).Add(20)
}

func Test_OrderedList_Detect(t *testing.T) {
	numbers := orderedlist.NewFrom[int](10, 20, 30, 40)
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

func Test_OrderedList_First(t *testing.T) {
	source := []string{`red`, `blue`, `green`}
	colors := orderedlist.NewFrom(source...)
	actual := colors.First()
	require.NoError(t, actual.Error())
	expected := source[0]
	require.Equal(t, expected, actual.MustGet())
}

func Test_OrderedList_First_fail(t *testing.T) {
	instance := orderedlist.New[int](0)
	actual := instance.First()
	require.Error(t, actual.Error())
}

func Test_OrderedList_Get(t *testing.T) {
	slice := []string{`red`, `blue`, `green`}
	colors := orderedlist.NewFrom(slice...)
	for i, expected := range slice {
		actual := colors.Get(i)
		require.NoError(t, actual.Error())
		require.Equal(t, expected, actual.MustGet())
	}
}

func Test_OrderedList_Get_fail(t *testing.T) {
	colors := orderedlist.NewFrom(`red`, `blue`)
	actual := colors.Get(-1)
	require.ErrorContains(t, actual.Error(), `the index -1 is out of bounds`)
	actual = colors.Get(3)
	require.ErrorContains(t, actual.Error(), `the index 3 is out of bounds`)
}

func Test_OrderedList_Index(t *testing.T) {
	slice := []string{`red`, `blue`, `green`}
	colors := orderedlist.NewFrom(slice...)
	for expected, color := range slice {
		actual := colors.Index(color)
		require.Equal(t, expected, actual)
	}
}

func Test_OrderedList_Index_fail(t *testing.T) {
	colors := orderedlist.NewFrom(`red`, `blue`)
	actual := colors.Index(`green`)
	expected := -1
	require.Equal(t, expected, actual)
}

func Test_OrderedList_Insert(t *testing.T) {
	numbers := orderedlist.NewFrom(`two`, `four`)
	err := numbers.Insert(0, `one`)
	require.NoError(t, err)
	actual := numbers.ToSlice()
	expected := sliceutils.ToSlice(`one`, `two`, `four`)
	require.Equal(t, expected, actual)
	err = numbers.Insert(2, `three`)
	require.NoError(t, err)
	actual = numbers.ToSlice()
	expected = sliceutils.ToSlice(`one`, `two`, `three`, `four`)
	require.Equal(t, expected, actual)
}

func Test_OrderedList_Insert_fail(t *testing.T) {
	numbers := orderedlist.New[string](3)
	err := numbers.Insert(0, `one`)
	require.Error(t, err)
	numbers = orderedlist.NewFrom(`one`, `two`)
	err = numbers.Insert(-1, `minus one`)
	require.Error(t, err)
	err = numbers.Insert(2, `three`)
	require.Error(t, err)
}

func Test_OrderedList_Largest(t *testing.T) {
	numbers := orderedlist.NewFrom(10, 23, 1, 100, 0)
	actual := numbers.Largest()
	require.NoError(t, actual.Error())
	expected := 100
	require.Equal(t, expected, actual.MustGet())
}

func Test_OrderedList_Last(t *testing.T) {
	source := []string{`red`, `blue`, `green`}
	colors := orderedlist.NewFrom(source...)
	actual := colors.Last()
	require.NoError(t, actual.Error())
	expected := source[len(source)-1]
	require.Equal(t, expected, actual.MustGet())
}

func Test_OrderedList_Last_fail(t *testing.T) {
	instance := orderedlist.New[int](0)
	actual := instance.Last()
	require.Error(t, actual.Error())
}

func Test_OrderedList_Largest_fail(t *testing.T) {
	numbers := orderedlist.New[int](0)
	actual := numbers.Largest()
	require.Error(t, actual.Error())
}

func Test_OrderedList_Remove_Contains(t *testing.T) {
	numbers := orderedlist.New[int](5)
	size := numbers.Len()
	require.Equal(t, 0, size)
	res := numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 0, size)
	require.Same(t, numbers, res)
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

func Test_OrderedList_Select(t *testing.T) {
	numbers := orderedlist.NewFrom[int](1, 2, 3, 4, 5, 2)
	isEvenFilter := func(value int) result.Result[bool] {
		return result.Ok(value%2 == 0)
	}
	evenNumbers := numbers.Select(isEvenFilter)
	require.NoError(t, evenNumbers.Error())
	actual := evenNumbers.MustGet().ToSlice()
	expected := []int{2, 4, 2}
	require.Equal(t, expected, actual)
}

func Test_OrderedList_Set(t *testing.T) {
	colors := orderedlist.NewFrom(`red`, `blue`, `green`)
	actual := colors.Set(1, `purple`)
	require.NoError(t, actual.Error())
	expected := `blue`
	require.Equal(t, expected, actual.MustGet())
	actual = colors.Get(1)
	require.NoError(t, actual.Error())
	expected = `purple`
	require.Equal(t, expected, actual.MustGet())
}

func Test_OrderedList_Set_fail(t *testing.T) {
	colors := orderedlist.NewFrom(`red`, `blue`)
	actual := colors.Set(-1, `green`)
	require.ErrorContains(t, actual.Error(), `the index -1 is out of bounds`)
	actual = colors.Set(3, `green`)
	require.ErrorContains(t, actual.Error(), `the index 3 is out of bounds`)
}

func Test_OrderedList_Smallest(t *testing.T) {
	values := orderedlist.NewFrom(`H`, `A`, `Z`)
	actual := values.Smallest()
	require.NoError(t, actual.Error())
	expected := `A`
	require.Equal(t, expected, actual.MustGet())
}

func Test_OrderedList_Smallest_fail(t *testing.T) {
	numbers := orderedlist.New[string](0)
	actual := numbers.Smallest()
	require.Error(t, actual.Error())
}

func Test_OrderedList_SortAscending(t *testing.T) {
	numbers := orderedlist.NewFrom(10, 5, 80, 2)
	numbers.SortAscending()
	actual := numbers.ToSlice()
	expected := []int{
		2, 5, 10, 80,
	}
	require.Equal(t, expected, actual)
}

func Test_OrderedList_SortDescending(t *testing.T) {
	numbers := orderedlist.NewFrom(`file-0011`, `file-0001`, `file-1234`)
	numbers.SortDescending()
	actual := numbers.ToSlice()
	expected := []string{
		`file-1234`, `file-0011`, `file-0001`,
	}
	require.Equal(t, expected, actual)
}

func Test_OrderedList_String(t *testing.T) {
	numbers := orderedlist.NewFrom[int](1, 2, 3)
	actual := numbers.String()
	expected := `[1 2 3]`
	require.Equal(t, expected, actual)
}

func Test_OrderedList_ToCollection(t *testing.T) {
	numbers := orderedlist.New[string](0)
	instance := numbers.ToCollection()
	require.NotNil(t, instance)
}

func Test_OrderedList_ToOrderedCollection(t *testing.T) {
	numbers := orderedlist.New[int](0)
	instance := numbers.ToOrderedCollection()
	require.NotNil(t, instance)
}
