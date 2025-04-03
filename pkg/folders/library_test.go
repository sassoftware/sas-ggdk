// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package folders_test

import (
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/folders"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/stretchr/testify/require"
)

func Test_NewCountingFolder(t *testing.T) {
	// For each element in source, increment step by 1 starting at 0
	source := []string{`red`, `blue`, `green`}
	step := 1
	folder := folders.NewCountingFolder[string](step)
	actual := sliceutils.Fold(folder, 0, source)
	require.NoError(t, actual.Error())
	expected := len(source)
	require.Equal(t, expected, actual.MustGet())
	// For each element in source, increment step by 10 starting at 100
	step = 10
	folder = folders.NewCountingFolder[string](step)
	actual = sliceutils.Fold(folder, 100, source)
	require.NoError(t, actual.Error())
	expected = 100 + len(source)*step
	require.Equal(t, expected, actual.MustGet())
}

func Test_NewMapFolder(t *testing.T) {
	// For each Color in the slice, set the key to be the upper case of the color
	// name and the value to be the length of the color name in an empty
	// map[string]int.
	type color struct {
		name string
	}
	source := []*color{
		{`red`},
		{`blue`},
		{`green`},
	}
	toKey := func(each *color) result.Result[string] {
		return result.Ok(strings.ToUpper(each.name))
	}
	toValue := func(each *color) result.Result[int] {
		return result.Ok(len(each.name))
	}
	folder := folders.NewMapFolder[*color, string, int](toKey, toValue)
	size := len(source)
	actual := sliceutils.Fold(
		folder,
		make(map[string]int, size),
		source,
	)
	require.NoError(t, actual.Error())
	expected := make(map[string]int, size)
	for _, each := range source {
		key := strings.ToUpper(each.name)
		expected[key] = len(each.name)
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_NewMapFolder_toKey_error(t *testing.T) {
	// Test error handling when the toKey function returns an error.
	expectedMessage := `failed key creation`
	toKey := func(_ string) result.Result[string] {
		return result.Error[string](errors.New(expectedMessage))
	}
	toValue := func(value string) result.Result[int] {
		return result.Ok(len(value))
	}
	folder := folders.NewMapFolder[string, string, int](toKey, toValue)
	actual := sliceutils.Fold(
		folder,
		make(map[string]int),
		[]string{`red`, `blue`, `green`},
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_NewMapFolder_toValue_error(t *testing.T) {
	// Test error handling when the toValue function returns an error.
	expectedMessage := `failed value creation`
	toKey := result.Ok[string]
	toValue := func(_ string) result.Result[int] {
		return result.Error[int](errors.New(expectedMessage))
	}
	folder := folders.NewMapFolder[string, string, int](toKey, toValue)
	actual := sliceutils.Fold(
		folder,
		make(map[string]int),
		[]string{`red`, `blue`, `green`},
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_NewSliceFolder(t *testing.T) {
	// Append each element of primary to an empty slice.
	folder := folders.NewSliceFolder[string]()
	primary := []string{`red`, `blue`, `green`}
	secondary := []string{`cyan`, `magenta`, `yellow`}
	size := sliceutils.LenAll(primary, secondary)
	initial := make([]string, 0, size)
	actual := sliceutils.Fold(folder, initial, primary)
	require.NoError(t, actual.Error())
	require.Equal(t, primary, actual.MustGet())
	// Append each element of secondary to the slice created by the last call to Fold.
	actual = sliceutils.Fold(folder, actual.MustGet(), secondary)
	require.NoError(t, actual.Error())
	expected := make([]string, 0, size)
	expected = append(expected, primary...)
	expected = append(expected, secondary...)
	require.Equal(t, expected, actual.MustGet())
}

func Test_LargestFolder(t *testing.T) {
	primary := []int{1, 5, 4, 3, 2}
	actual := sliceutils.Fold(folders.LargestFolder[int], primary[0], primary)
	require.NoError(t, actual.Error())
	require.Equal(t, 5, actual.MustGet())
}

func Test_SmallestFolder(t *testing.T) {
	primary := []int{5, 4, 3, 1, 2}
	actual := sliceutils.Fold(folders.SmallestFolder[int], primary[0], primary)
	require.NoError(t, actual.Error())
	require.Equal(t, 1, actual.MustGet())
}

func Test_NewMatchesFolder(t *testing.T) {
	primary := []int{5, 4, 3, 1, 2}
	filter := func(i int) result.Result[bool] {
		if i < 4 {
			return result.Ok(true)
		}
		return result.Ok(false)
	}
	folder := folders.NewMatchesFolder(filter)
	actual := sliceutils.Fold(folder, make([]int, 0), primary)
	require.NoError(t, actual.Error())
	require.Equal(t, []int{3, 1, 2}, actual.MustGet())
}

func Test_NewMatchesFolder_Error(t *testing.T) {
	primary := []int{5, 4, 3, 1, 2}
	filter := func(i int) result.Result[bool] {
		if i < 4 {
			return result.Ok(true)
		}
		return result.Error[bool](errors.New("failed"))
	}
	folder := folders.NewMatchesFolder(filter)
	actual := sliceutils.Fold(folder, make([]int, 0), primary)
	require.ErrorContains(t, actual.Error(), "failed")
}

func Test_NewMatchesFolderNoError(t *testing.T) {
	primary := []int{5, 4, 3, 1, 2}
	filter := func(i int) bool {
		return i < 4
	}
	folder := folders.NewMatchesFolderNoError(filter)
	actual := sliceutils.FoldNoError(folder, make([]int, 0), primary)
	require.Equal(t, []int{3, 1, 2}, actual)
}
