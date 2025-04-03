// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package collections_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/collections"
	"github.com/sassoftware/sas-ggdk/pkg/collections/orderedlist"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_Add(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1))
	s = result.MapNoError2(collections.Add[int], s, result.Ok(2))
	require.NoError(t, s.Error())
	require.Equal(t, []int{1, 2}, s.MustGet().ToSlice())
}

func Test_Test_Contains(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1))
	present := result.MapNoError2(collections.Contains[int], s, result.Ok(1))
	require.NoError(t, present.Error())
	require.True(t, present.MustGet())
	absent := result.MapNoError2(collections.Contains[int], s, result.Ok(1))
	require.NoError(t, absent.Error())
	require.True(t, absent.MustGet())
}

func Test_Test_Detect(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1))
	detectPresent := func(v int) result.Result[bool] {
		return result.Ok(v == 1)
	}
	detectAbsent := func(v int) result.Result[bool] {
		return result.Ok(v == 2)
	}
	detectError := func(_ int) result.Result[bool] {
		return result.Error[bool](errors.New("failed"))
	}
	present := result.FlatMap2(collections.Detect[int], s, result.Ok(filters.Filter[int](detectPresent)))
	require.NoError(t, present.Error())
	require.True(t, present.MustGet().IsJust())
	require.Equal(t, 1, present.MustGet().MustGet())
	absent := result.FlatMap2(collections.Detect[int], s, result.Ok(filters.Filter[int](detectAbsent)))
	require.NoError(t, absent.Error())
	require.False(t, absent.MustGet().IsJust())
	err := result.FlatMap2(collections.Detect[int], s, result.Ok(filters.Filter[int](detectError)))
	require.ErrorContains(t, err.Error(), "failed")
}

func Test_Len(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1, 2, 3))
	r := result.MapNoError(collections.Len[int], s)
	require.NoError(t, r.Error())
	require.Equal(t, 3, r.MustGet())
}

func Test_Remove(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1, 2, 3))
	r := result.MapNoError2(collections.Remove[int], s, result.Ok(2))
	require.NoError(t, r.Error())
	require.Equal(t, []int{1, 3}, r.MustGet().ToSlice())
}

func Test_Select(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1, 2, 3))
	selectPresent := func(v int) result.Result[bool] {
		return result.Ok(v == 1 || v == 2)
	}
	selectAbsent := func(v int) result.Result[bool] {
		return result.Ok(v == 100)
	}
	selectError := func(_ int) result.Result[bool] {
		return result.Error[bool](errors.New("failed"))
	}
	present := result.FlatMap2(collections.Select[int], s, result.Ok(filters.Filter[int](selectPresent)))
	require.NoError(t, present.Error())
	require.True(t, present.MustGet().Contains(1))
	require.True(t, present.MustGet().Contains(2))
	require.False(t, present.MustGet().Contains(3))
	absent := result.FlatMap2(collections.Select[int], s, result.Ok(filters.Filter[int](selectAbsent)))
	require.NoError(t, absent.Error())
	require.Equal(t, 0, absent.MustGet().Len())
	err := result.FlatMap2(collections.Select[int], s, result.Ok(filters.Filter[int](selectError)))
	require.ErrorContains(t, err.Error(), "failed")
}

func Test_ToSlice(t *testing.T) {
	s := result.Ok[collections.Collection[int]](orderedlist.NewFrom(1, 2, 3))
	sl := result.MapNoError(collections.ToSlice[int], s)
	require.NoError(t, sl.Error())
	require.Equal(t, []int{1, 2, 3}, sl.MustGet())
}

func Test_First(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	first := result.FlatMap(collections.First[int], s)
	require.NoError(t, first.Error())
	require.Equal(t, 1, first.MustGet())
}

func Test_Get(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	first := result.FlatMap2(collections.Get[int], s, result.Ok(2))
	require.NoError(t, first.Error())
	require.Equal(t, 3, first.MustGet())
	err := result.FlatMap2(collections.Get[int], s, result.Ok(3))
	require.ErrorContains(t, err.Error(), `index 3 is out of bounds`)
}

func Test_Index(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	first := result.MapNoError2(collections.Index[int], s, result.Ok(2))
	require.NoError(t, first.Error())
	require.Equal(t, 1, first.MustGet())
	err := result.FlatMap2(collections.Get[int], s, result.Ok(3))
	require.ErrorContains(t, err.Error(), `index 3 is out of bounds`)
}

func Test_Insert(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	s = result.FlatMap3(collections.Insert[int], s, result.Ok(2), result.Ok(4))
	require.NoError(t, s.Error())
	require.Equal(t, []int{1, 2, 4, 3}, s.MustGet().ToSlice())
	s = result.FlatMap3(collections.Insert[int], s, result.Ok(9), result.Ok(4))
	require.ErrorContains(t, s.Error(), `index 9 is out of bounds`)
}

func Test_Largest(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	largest := result.FlatMap(collections.Largest[int], s)
	require.NoError(t, largest.Error())
	require.Equal(t, 3, largest.MustGet())
}

func Test_Last(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	last := result.FlatMap(collections.Last[int], s)
	require.NoError(t, last.Error())
	require.Equal(t, 3, last.MustGet())
}

func Test_Set(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	replaced := result.FlatMap3(collections.Set[int], s, result.Ok(2), result.Ok(4))
	require.NoError(t, replaced.Error())
	require.Equal(t, 3, replaced.MustGet())
	replaced = result.FlatMap3(collections.Set[int], s, result.Ok(9), result.Ok(4))
	require.ErrorContains(t, replaced.Error(), `index 9 is out of bounds`)
}

func Test_Smallest(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 2, 3))
	smallest := result.FlatMap(collections.Smallest[int], s)
	require.NoError(t, smallest.Error())
	require.Equal(t, 1, smallest.MustGet())
}

func Test_SortAscending(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 3, 2))
	sorted := result.MapNoError(collections.SortAscending[int], s)
	require.NoError(t, sorted.Error())
	require.Equal(t, []int{1, 2, 3}, sorted.MustGet().ToSlice())
}

func Test_SortDescending(t *testing.T) {
	s := result.Ok[collections.OrderedCollection[int]](orderedlist.NewFrom(1, 3, 2))
	sorted := result.MapNoError(collections.SortDescending[int], s)
	require.NoError(t, sorted.Error())
	require.Equal(t, []int{3, 2, 1}, sorted.MustGet().ToSlice())
}
