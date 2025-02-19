// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maputils_test

import (
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/maputils"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func Test_AssertValuesNotNil(t *testing.T) {
	m := map[string]*int32{
		`small`:  proto.Int32(10),
		`medium`: proto.Int32(20),
		`large`:  proto.Int32(30),
	}
	err := maputils.AssertValuesNotNil(m)
	require.NoError(t, err)
	m[`xl`] = nil
	err = maputils.AssertValuesNotNil(m)
	require.Error(t, err)
	message := err.Error()
	state := strings.HasSuffix(message, `missing required fields: "xl"`)
	require.True(t, state)
	m[`xxl`] = nil
	err = maputils.AssertValuesNotNil(m)
	require.Error(t, err)
	message = err.Error()
	state = strings.HasSuffix(message, `missing required fields: "xl", "xxl"`)
	require.True(t, state)
}

func Test_DeleteKeys(t *testing.T) {
	one := `one`
	two := `two`
	three := `three`
	four := `four`
	m := map[string]int{
		one:   1,
		two:   2,
		three: 3,
		four:  4,
	}
	require.Len(t, m, 4)
	maputils.DeleteKeys(m, two)
	_, exists := m[two]
	require.False(t, exists)
	require.Len(t, m, 3)
	maputils.DeleteKeys(m, one, three)
	_, exists = m[one]
	require.False(t, exists)
	_, exists = m[three]
	require.False(t, exists)
	require.Len(t, m, 1)
	maputils.DeleteKeys(m, one)
	require.Len(t, m, 1)
	_, exists = m[four]
	require.True(t, exists)
}

func Test_Map(t *testing.T) {
	m := map[string]string{
		`rabbit`: `Bugs`,
		`pig`:    `Porky`,
		`horse`:  `Shadowfax`,
	}
	mapper := func(animal string) result.Result[int] {
		name := m[animal]
		return result.Ok(len(name))
	}
	actual := maputils.Map(mapper, m)
	require.NoError(t, actual.Error())
	require.NotNil(t, actual.MustGet())
	require.Equal(t, 4, actual.MustGet()[`rabbit`])
	require.Equal(t, 5, actual.MustGet()[`pig`])
	require.Equal(t, 9, actual.MustGet()[`horse`])
}

func Test_Map_error(t *testing.T) {
	m := map[string]string{
		`key`: `value`,
	}
	failedMap := `failed MAP`
	mapper := func(value string) result.Result[int] {
		return result.Error[int](errors.New(failedMap))
	}
	actual := maputils.Map(mapper, m)
	require.ErrorContains(t, actual.Error(), failedMap)
}

func Test_Map_withSelect(t *testing.T) {
	m := map[string]string{
		`rabbit`: `Bugs`,
		`pig`:    `Porky`,
		`horse`:  `Shadowfax`,
	}
	filter := func(animal string) result.Result[bool] {
		return result.Ok(animal != `pig`)
	}
	mapper := func(animal string) result.Result[int] {
		name := m[animal]
		return result.Ok(len(name))
	}
	filtered := maputils.Select(m, filter)
	actual := maputils.MapResult(filtered, mapper)
	require.NoError(t, actual.Error())
	require.NotNil(t, actual.MustGet())
	require.Equal(t, 4, actual.MustGet()[`rabbit`])
	require.Equal(t, 9, actual.MustGet()[`horse`])
}

func Test_Merge(t *testing.T) {
	m1 := map[string]string{
		`dog`: `Rover`,
		`cat`: `Luna`,
	}
	m2 := map[string]string{
		`rabbit`: `Bugs`,
		`dog`:    `Shep`,
		`pig`:    `Porky`,
	}
	m3 := map[string]string{
		`horse`: `Shadowfax`,
	}
	actual := maputils.Merge(m1, m2, m3)
	expected := map[string]string{
		`cat`:    `Luna`,
		`rabbit`: `Bugs`,
		`dog`:    `Shep`,
		`pig`:    `Porky`,
		`horse`:  `Shadowfax`,
	}
	require.Equal(t, expected, actual)
}

func Test_Fold(t *testing.T) {
	// id -> name
	m := map[int]string{
		61: `luna`,
		32: `bugs`,
		35: `shep`,
		74: `porky`,
		50: `shadowfax`,
	}
	transformer := func(accumulator map[string]int, id int) result.Result[map[string]int] {
		animalName := m[id]
		for _, i32 := range animalName { // Each letter.
			letter := string(i32)
			accumulator[letter]++
		}
		return result.Ok(accumulator)
	}
	actual := maputils.Fold(
		transformer,
		make(map[string]int, 26),
		m,
	)
	require.NoError(t, actual.Error())
	expected := map[string]int{
		`a`: 3, `b`: 1, `d`: 1, `e`: 1,
		`f`: 1, `g`: 1, `h`: 2, `k`: 1,
		`l`: 1, `n`: 1, `o`: 2, `p`: 2,
		`r`: 1, `s`: 3, `u`: 2, `w`: 1,
		`x`: 1, `y`: 1}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_fail(t *testing.T) {
	m := map[string]string{
		`key`: `value`,
	}
	expectedMessage := `failed REDUCE`
	transformer := func(target map[string]int, source string) result.Result[map[string]int] {
		return result.Error[map[string]int](errors.New(expectedMessage))
	}
	actual := maputils.Fold(
		transformer,
		make(map[string]int, 0),
		m,
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Fold_withSelect_fail(t *testing.T) {
	m := map[string]string{
		`key`: `value`,
	}
	expectedMessage := `failed REDUCE filter`
	filter := func(animal string) result.Result[bool] {
		return result.Error[bool](errors.New(expectedMessage))
	}
	transformer := func(target map[string]int, animal string) result.Result[map[string]int] {
		return result.Ok(target)
	}
	filtered := maputils.Select(m, filter)
	actual := maputils.FoldResult(
		transformer,
		make(map[string]int, 26),
		filtered,
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Fold_withSelect_nil(t *testing.T) {
	// animal -> name
	m := map[string]string{
		`cat`:    `luna`,
		`rabbit`: `bugs`,
		`dog`:    `shep`,
		`pig`:    `porky`,
		`horse`:  `shadowfax`,
	}
	transformer := func(accumulator map[string]int, animal string) result.Result[map[string]int] {
		for _, i32 := range m[animal] {
			letter := string(i32)
			accumulator[letter]++
		}
		return result.Ok(accumulator)
	}
	filtered := maputils.Select(m, nil)
	actual := maputils.FoldResult(
		transformer,
		make(map[string]int, 26),
		filtered,
	)
	require.NoError(t, actual.Error())
	expected := map[string]int{}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_withFilter(t *testing.T) {
	// animal -> name
	m := map[string]string{
		`cat`:    `luna`,
		`rabbit`: `bugs`,
		`dog`:    `shep`,
		`pig`:    `porky`,
		`horse`:  `shadowfax`,
	}
	filter := func(animal string) result.Result[bool] {
		return result.Ok(animal == `cat` || animal == `dog`)
	}
	transformer := func(accumulator map[string]int, animal string) result.Result[map[string]int] {
		for _, i32 := range m[animal] {
			letter := string(i32)
			accumulator[letter]++
		}
		return result.Ok(accumulator)
	}
	filtered := maputils.Select(m, filter)
	actual := maputils.FoldResult(
		transformer,
		make(map[string]int, 26),
		filtered,
	)
	require.NoError(t, actual.Error())
	expected := map[string]int{
		`a`: 1, `e`: 1, `h`: 1, `l`: 1,
		`n`: 1, `p`: 1, `s`: 1, `u`: 1,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_withFilter_fail(t *testing.T) {
	m := map[string]string{
		`key`: `value`,
	}
	expectedMessage := `failed REDUCE filter`
	filter := func(animal string) result.Result[bool] {
		return result.Error[bool](errors.New(expectedMessage))
	}
	transformer := func(target map[string]int, animal string) result.Result[map[string]int] {
		return result.Ok(target)
	}
	filtered := maputils.Select(m, filter)
	actual := maputils.FoldResult(
		transformer,
		make(map[string]int, 26),
		filtered,
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Select(t *testing.T) {
	source := map[string]int{
		`red`:   3,
		`green`: 5,
		`blue`:  4,
	}
	filter := func(key string) result.Result[bool] {
		return result.Ok(source[key] > 3)
	}
	actual := maputils.Select(source, filter)
	require.NoError(t, actual.Error())
	expected := map[string]int{
		`green`: 5,
		`blue`:  4,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Select_nil(t *testing.T) {
	source := map[string]int{
		`red`:   3,
		`green`: 5,
		`blue`:  4,
	}
	actual := maputils.Select(source, nil)
	require.NoError(t, actual.Error())
	require.NotNil(t, actual.MustGet())
	require.Empty(t, actual.MustGet())
}

func Test_ToAnyMap(t *testing.T) {
	expected := map[any]any{
		`key1`: `value1`,
		`key2`: `value2`,
	}
	actual := maputils.ToAnyMap(expected)
	require.NotNil(t, actual)
	require.Len(t, actual, 2)
}

func Test_ToAnyMap_nil(t *testing.T) {
	actual := maputils.ToAnyMap(nil)
	require.Nil(t, actual)
	var m map[string]string
	actual = maputils.ToAnyMap(m)
	require.Nil(t, actual)
}

func Test_ToAnyMap_fail(t *testing.T) {
	var expected any
	actual := maputils.ToAnyMap(expected)
	require.Nil(t, actual)
	expected = `Hello`
	actual = maputils.ToAnyMap(expected)
	require.Nil(t, actual)
}

func Test_ToBool(t *testing.T) {
	m := make(map[string]any)
	key := `true`
	m[key] = true
	state, err := maputils.ToBool(m, key)
	require.NoError(t, err)
	require.True(t, state)
	key = `false`
	m[key] = false
	state, err = maputils.ToBool(m, key)
	require.NoError(t, err)
	require.False(t, state)
	key = `nope`
	m[key] = `not a bool`
	state, err = maputils.ToBool(m, key)
	require.Error(t, err)
	require.False(t, state)
	key = `absent`
	state, err = maputils.ToBool(m, key)
	require.Error(t, err)
	require.False(t, state)
}

func Test_Get(t *testing.T) {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	v := maputils.Get(m, "one")
	require.True(t, v.IsJust())
	require.Equal(t, 1, v.MustGet())
	v = maputils.Get(m, "four")
	require.False(t, v.IsJust())
	require.Equal(t, 4, v.OrElse(4))
}
