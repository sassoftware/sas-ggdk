// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package set_test

import (
	"sort"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/collections/list"
	"github.com/sassoftware/sas-ggdk/pkg/collections/set"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/stretchr/testify/require"
)

func Test_Map(t *testing.T) {
	colors := set.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(value string) result.Result[int] {
		return result.Ok(len(value))
	}
	actual := set.Map[string, int](mapper, colors)
	require.NoError(t, actual.Error())
	actualSlice := actual.MustGet().ToSlice()
	sort.Ints(actualSlice)
	expected := set.New[int](colors.Len())
	for _, color := range colors.ToSlice() {
		length := mapper(color)
		require.NoError(t, length.Error())
		expected.Add(length.MustGet())
	}
	expectedSlice := expected.ToSlice()
	sort.Ints(expectedSlice)
	require.Equal(t, expectedSlice, actualSlice)
}

func Test_Map_fail(t *testing.T) {
	colors := set.NewFrom[string](`red`, `green`, `blue`)
	mapper := func(_ string) result.Result[int] {
		return result.Error[int](errors.New(`fail`))
	}
	actual := set.Map[string, int](mapper, colors)
	require.Error(t, actual.Error())
}

func Test_raw_new(t *testing.T) {
	instance := new(set.Set[string])
	require.NotNil(t, instance)
	instance.Add(`red`)
	require.True(t, instance.Contains(`red`))
}

func Test_New(t *testing.T) {
	numbers := set.New[int](5)
	require.NotNil(t, numbers)
	size := numbers.Len()
	require.Equal(t, 0, size)
}

func Test_NewFrom(t *testing.T) {
	expected := []int{1, 2, 3}
	numbers := set.NewFrom[int](expected...)
	require.NotNil(t, numbers)
	actualSize := numbers.Len()
	expectedSize := len(expected)
	require.Equal(t, expectedSize, actualSize)
	actual := numbers.ToSlice()
	sort.Ints(actual)
	require.Equal(t, expected, actual)
}

func Test_NewFromCollection(t *testing.T) {
	source := list.NewFrom(1, 2, 3, 2)
	numbers := set.NewFromCollection[int](source)
	require.NotNil(t, numbers)
	actualSize := numbers.Len()
	expectedSize := 3
	require.Equal(t, expectedSize, actualSize)
	actual := numbers.ToSlice()
	sort.Ints(actual)
	expected := []int{1, 2, 3}
	require.Equal(t, expected, actual)
}

func Test_NewWithAccessor(t *testing.T) {
	words, access := set.NewWithAccessor[string](5)
	words.
		Add(`the`).
		Add(`day`).
		Add(`of`).
		Add(`the`).
		Add(`triffids`)
	actual := access()
	expected := map[string]int{
		`the`:      1,
		`day`:      1,
		`of`:       1,
		`triffids`: 1,
	}
	require.Equal(t, expected, actual)
}

func Test_Set_Add_Contains(t *testing.T) {
	numbers := set.New[int](5)
	state := numbers.Contains(10)
	require.False(t, state)
	res := numbers.Add(10)
	require.Same(t, numbers, res)
	size := numbers.Len()
	require.Equal(t, 1, size)
	state = numbers.Contains(10)
	require.True(t, state)
	numbers.Add(11)
	size = numbers.Len()
	require.Equal(t, 2, size)
	state = numbers.Contains(11)
	require.True(t, state)
	numbers.Add(11)
	size = numbers.Len()
	require.Equal(t, 2, size)
}

func Test_Set_Detect(t *testing.T) {
	numbers := set.NewFrom[int](10, 20, 30, 40)
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

func Test_Set_Len(t *testing.T) {
	numbers := set.New[int](5)
	size := numbers.Len()
	require.Equal(t, 0, size)
	numbers.Add(10)
	size = numbers.Len()
	require.Equal(t, 1, size)
	numbers.Add(10)
	size = numbers.Len()
	require.Equal(t, 1, size)
	numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 0, size)
	numbers.Remove(10)
	size = numbers.Len()
	require.Equal(t, 0, size)
}

func Test_Set_Remove_Contains(t *testing.T) {
	numbers := set.New[int](5)
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
	require.Equal(t, 2, size)
	state := numbers.Contains(10)
	require.True(t, state)
	numbers.Remove(10)
	state = numbers.Contains(10)
	require.False(t, state)
	size = numbers.Len()
	require.Equal(t, 1, size)
	numbers.Remove(20)
	state = numbers.Contains(20)
	require.False(t, state)
	size = numbers.Len()
	require.Equal(t, 0, size)
}

func Test_Set_Select(t *testing.T) {
	numbers := set.NewFrom[int](1, 2, 3, 4, 5)
	isEvenFilter := func(value int) result.Result[bool] {
		return result.Ok(value%2 == 0)
	}
	evenNumbers := numbers.Select(isEvenFilter)
	require.NoError(t, evenNumbers.Error())
	actual := evenNumbers.MustGet().ToSlice()
	sort.Ints(actual)
	expected := []int{2, 4}
	require.Equal(t, expected, actual)
}

func Test_Set_Select_error(t *testing.T) {
	numbers := set.NewFrom[int](1, 2)
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

func Test_Set_ToSlice(t *testing.T) {
	numbers := set.New[int](7)
	actual := numbers.ToSlice()
	require.NotNil(t, actual)
	require.Empty(t, actual)
	numbers.Add(1, 2, 3)
	actual = numbers.ToSlice()
	sort.Ints(actual)
	expected := []int{1, 2, 3}
	require.Equal(t, expected, actual)
}

func Test_Set_String(t *testing.T) {
	numbers := set.New[int](3)
	numbers.Add(1, 2, 3)
	target := numbers.String()
	possibilities := []string{
		`[1 2 3]`,
		`[1 3 2]`,
		`[2 1 3]`,
		`[2 3 1]`,
		`[3 1 2]`,
		`[3 2 1]`,
	}
	actual := sliceutils.Detect(
		filters.NewIsEqualFilter(target),
		possibilities,
	)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
}

func Test_Set_ToCollection(t *testing.T) {
	numbers := set.New[int](0)
	instance := numbers.ToCollection()
	require.NotNil(t, instance)
}
