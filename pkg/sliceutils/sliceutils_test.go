// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sliceutils_test

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/pointer"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AnySliceToSlice(t *testing.T) {
	value := []any{
		1, 2, 3,
	}
	actual := sliceutils.AnySliceToSlice[int](value)
	require.True(t, actual.IsJust())
	expected := []int{
		1, 2, 3,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_AnySliceToSlice_fail(t *testing.T) {
	value := []any{ // Not strings.
		1, 2, 3,
	}
	actual := sliceutils.AnySliceToSlice[string](value)
	require.False(t, actual.IsJust())
	value = nil // Not a slice of strings.
	actual = sliceutils.AnySliceToSlice[string](value)
	require.False(t, actual.IsJust())
}

func Test_AssertContains(t *testing.T) {
	var values []string
	err := sliceutils.AssertContains(values, `red`)
	require.Error(t, err)
	values = []string{
		`red`,
		`blue`,
		`green`,
	}
	err = sliceutils.AssertContains(values, `blue`)
	require.NoError(t, err)
	err = sliceutils.AssertContains(values, `purple`)
	require.Error(t, err)
	values = make([]string, 0)
	err = sliceutils.AssertContains(values, `purple`)
	require.Error(t, err)
	values = nil
	err = sliceutils.AssertContains(values, `purple`)
	require.Error(t, err)
}

func Test_CollectErrors(t *testing.T) {
	errs := []result.Result[string]{
		result.Ok("ok"),
		result.Error[string](errors.New("failure1")),
		result.Error[string](errors.New("failure2")),
		result.Ok("ok"),
		result.Error[string](errors.New("failure3")),
	}
	err := sliceutils.CollectErrors(errors.New("Failure:"), errs...)
	require.Error(t, err)
	reg := `(?ms)Failure:$.*Caused by:$.*failure1$.*failure2$.*failure3$`
	msg := err.Error()
	require.Regexp(t, reg, msg)
}

func Test_CollectErrors_no_errors(t *testing.T) {
	errs := []result.Result[string]{
		result.Ok("ok"),
		result.Ok("ok"),
	}
	err := sliceutils.CollectErrors(errors.New("Failure:"), errs...)
	require.NoError(t, err)
}

func Test_CollectErrors_no_elements(t *testing.T) {
	errs := []result.Result[string]{}
	err := sliceutils.CollectErrors(errors.New("Failure:"), errs...)
	require.NoError(t, err)
}

func Test_ContainsAll(t *testing.T) {
	haystack := []string{"one", "two", "three", "four", "five"}
	presentNeedles := []string{"one", "three", "four"}
	absentNeedles1 := []string{"one", "two", "three", "four", "five", "six"}
	absentNeedles2 := []string{"one", "six"}
	require.True(t, sliceutils.ContainsAll(haystack, presentNeedles))
	require.False(t, sliceutils.ContainsAll(haystack, absentNeedles1))
	require.False(t, sliceutils.ContainsAll(haystack, absentNeedles2))
}

func Test_Detect_error(t *testing.T) {
	expectedMessage := `failed DETECT filter`
	filter := newFailingFilter[string](expectedMessage)
	actual := sliceutils.Detect(
		filter,
		[]string{`red`, `blue`, `green`},
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

type myType struct {
	value int
}

func Test_Detect_myType_pointer(t *testing.T) {
	data := []*myType{{3}, {7}, {9}}
	filter := func(each *myType) result.Result[bool] {
		return result.Ok(each.value > 5)
	}
	actual := sliceutils.Detect(filter, data)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	require.NotNil(t, actual.MustGet().MustGet())
	require.Same(t, data[1], actual.MustGet().MustGet())
}

func Test_Detect_myType_pointer_notDetected(t *testing.T) {
	data := []*myType{{3}, {7}, {9}}
	filter := func(each *myType) result.Result[bool] {
		return result.Ok(each.value > 15)
	}
	actual := sliceutils.Detect(filter, data)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_Detect_myType_value(t *testing.T) {
	data := []myType{{3}, {7}, {9}}
	filter := func(each myType) result.Result[bool] {
		return result.Ok(each.value > 5)
	}
	actual := sliceutils.Detect(filter, data)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	require.NotNil(t, actual.MustGet().MustGet())
	require.Equal(t, data[1], actual.MustGet().MustGet())
	require.NotSame(t, &data[1], pointer.Ptr(actual.MustGet().MustGet()))
}

func Test_Detect_myType_value_notDetected(t *testing.T) {
	data := []myType{{3}, {7}, {9}}
	filter := func(each myType) result.Result[bool] {
		return result.Ok(each.value > 15)
	}
	actual := sliceutils.Detect(filter, data)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_Detect_nil(t *testing.T) {
	actual := sliceutils.Detect(
		nil,
		[]string{`red`, `blue`, `green`},
	)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_Detect_notDetected(t *testing.T) {
	filter := func(color string) result.Result[bool] {
		return result.Ok(color == `orange`)
	}
	actual := sliceutils.Detect(
		filter,
		[]string{`red`, `blue`, `green`},
	)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_Detect_success(t *testing.T) {
	filter := func(v int) result.Result[bool] {
		return result.Ok(v > 2)
	}
	actual := sliceutils.Detect(
		filter,
		[]int{1, 2, 3, 4},
	)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	expected := 3
	require.Equal(t, expected, actual.MustGet().MustGet())
}

func Test_Detect_zeroValue(t *testing.T) {
	filter := func(v int) result.Result[bool] {
		return result.Ok(v == 0)
	}
	actual := sliceutils.Detect(
		filter,
		[]int{1, 2, 3, 4},
	)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet().IsJust())
}

func Test_DetectNoError(t *testing.T) {
	filter := func(v int) bool {
		return v == 0
	}
	actual := sliceutils.DetectNoError(
		filter,
		[]int{1, 2, 0, 3, 4},
	)
	require.True(t, actual.IsJust())
	require.Equal(t, 0, actual.MustGet())
}

func Test_DetectNoError_noFilter(t *testing.T) {
	actual := sliceutils.DetectNoError(
		nil,
		[]int{1, 2, 3, 4},
	)
	require.False(t, actual.IsJust())
}

func Test_DetectNoError_notDetected(t *testing.T) {
	filter := func(v int) bool {
		return v == 0
	}
	actual := sliceutils.DetectNoError(
		filter,
		[]int{1, 2, 3, 4},
	)
	require.False(t, actual.IsJust())
}

func Test_DetectNoErrorResult(t *testing.T) {
	filter := func(v int) bool {
		return v == 0
	}
	actual := sliceutils.DetectNoErrorResult(
		filter,
		result.Ok([]int{1, 2, 0, 3, 4}),
	)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	require.Equal(t, 0, actual.MustGet().MustGet())
}

func Test_DetectResult(t *testing.T) {
	filter := func(v int) result.Result[bool] {
		return result.Ok(v == 0)
	}
	actual := sliceutils.DetectResult(
		filter,
		result.Ok([]int{1, 2, 0, 3, 4}),
	)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet().IsJust())
	require.Equal(t, 0, actual.MustGet().MustGet())
}

func Test_Disjoint(t *testing.T) {
	left := []int{
		1, 2, 3,
	}
	right := []int{
		2, 3, 4,
	}
	actual := sliceutils.Disjoint(left, right)
	sort.Ints(actual.MustGet())
	expected := []int{
		1, 4,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Disjoint_empty(t *testing.T) {
	left := []int{
		1, 2, 3,
	}
	right := []int{
		1, 2, 3,
	}
	actual := sliceutils.Disjoint(left, right)
	require.NotNil(t, actual.MustGet())
	require.Empty(t, actual.MustGet())
}

func Test_Duplicate(t *testing.T) {
	src := []int{1, 2, 3}
	dst := sliceutils.Duplicate(src)
	require.Equal(t, src, dst)
	require.NotSame(t, &src[0], &dst[0])
}

func Test_FirstError(t *testing.T) {
	errs := []result.Result[string]{
		result.Ok("ok"),
		result.Error[string](errors.New("failure1")),
		result.Error[string](errors.New("failure2")),
		result.Ok("ok"),
		result.Error[string](errors.New("failure3")),
	}
	err := sliceutils.FirstError(errs...)
	require.ErrorContains(t, err, "failure1")
}

func Test_FirstError_noErrors(t *testing.T) {
	errs := []result.Result[string]{
		result.Ok("ok"),
		result.Ok("ok"),
	}
	err := sliceutils.FirstError(errs...)
	require.NoError(t, err)
}

func Test_FirstError_noElements(t *testing.T) {
	errs := []result.Result[string]{}
	err := sliceutils.FirstError(errs...)
	require.NoError(t, err)
}

func Test_Fold_Buffer(t *testing.T) {
	source := []int{10, 20, 30}
	count := 0
	toCommaSeparatedString := func(accumulator *bytes.Buffer, number int) result.Result[*bytes.Buffer] {
		value := strconv.Itoa(number)
		accumulator.WriteString(value)
		if count < len(source)-1 {
			accumulator.WriteString(`, `)
		}
		count++
		return result.Ok(accumulator)
	}
	actual := sliceutils.Fold(
		toCommaSeparatedString,
		bytes.NewBuffer(make([]byte, 0, 10)),
		source,
	)
	require.NoError(t, actual.Error())
	actualString := actual.MustGet().String()
	expectedString := `10, 20, 30`
	require.Equal(t, expectedString, actualString)
}

func Test_Fold_int(t *testing.T) {
	transformer := func(accumulator int, number int) result.Result[int] {
		return result.Ok(accumulator + number)
	}
	actual := sliceutils.Fold(
		transformer,
		0,
		[]int{10, 20, 30},
	)
	require.NoError(t, actual.Error())
	expected := 60
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_int_error(t *testing.T) {
	expected := 0
	expectedMessage := `failed REDUCE`
	transformer := func(total int, number int) result.Result[int] {
		return result.Error[int](errors.New(expectedMessage))
	}
	actual := sliceutils.Fold(
		transformer,
		expected,
		[]int{10, 20, 30},
	)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Fold_pointer(t *testing.T) {
	transformer := func(accumulator int, color *string) result.Result[int] {
		if color != nil { // It's the transformer's responsibility to handle nil.
			accumulator += len(*color)
		}
		return result.Ok(accumulator)
	}
	red := `red`
	blue := `blue`
	green := `green`
	colors := sliceutils.ToSlice(
		&red,
		&blue,
		nil, // Ensure nil is passed to the transformer.
		&green,
	)
	actual := sliceutils.Fold(
		transformer,
		0,
		colors,
	)
	require.NoError(t, actual.Error())
	expected := len(red) + len(blue) + len(green)
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_withSelect(t *testing.T) {
	transformer := func(accumulator int, number int) result.Result[int] {
		return result.Ok(accumulator + number)
	}
	// The filters.Filter type is required here.
	filter1 := filters.Filter[int](func(value int) result.Result[bool] {
		return result.Ok(value >= 15)
	})
	filter2 := func(value int) result.Result[bool] {
		return result.Ok(value <= 55)
	}
	filter := filter1.And(filter2)
	filtered := sliceutils.Select(filter, []int{10, 20, 30, 40, 50, 60})
	actual := sliceutils.FoldResult(transformer, 0, filtered)
	require.NoError(t, actual.Error())
	expected := 140
	require.Equal(t, expected, actual.MustGet())
}

func Test_Fold_withSelect_error(t *testing.T) {
	transformer := func(accumulator int, number int) result.Result[int] {
		return result.Ok(accumulator + number)
	}
	expectedMessage := `failed REDUCE filter`
	filtered := sliceutils.Select(newFailingFilter[int](expectedMessage), []int{10, 20, 30, 40, 50, 60})
	actual := sliceutils.FoldResult(transformer, 0, filtered)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Fold_withSelect_nil(t *testing.T) {
	transformer := func(accumulator int, number int) result.Result[int] {
		return result.Ok(accumulator + number)
	}
	filtered := sliceutils.Select(nil, []int{10, 20, 30, 40, 50, 60})
	actual := sliceutils.FoldResult(transformer, 0, filtered)
	require.NoError(t, actual.Error())
	expected := 0
	require.Equal(t, expected, actual.MustGet())
}

func Test_FoldNoError_int(t *testing.T) {
	transformer := func(accumulator int, number int) int {
		return accumulator + number
	}
	actual := sliceutils.FoldNoError(
		transformer,
		0,
		[]int{10, 20, 30},
	)
	expected := 60
	require.Equal(t, expected, actual)
}

func Test_FoldNoError_result(t *testing.T) {
	transformer := func(accumulator int, number int) int {
		return accumulator + number
	}
	actual := sliceutils.FoldNoErrorResult(
		transformer,
		0,
		result.Ok([]int{10, 20, 30}),
	)
	expected := 60
	require.NoError(t, actual.Error())
	require.Equal(t, expected, actual.MustGet())
}

func Test_Head(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	head := sliceutils.Head(s)
	require.True(t, head.IsJust())
	require.Equal(t, 1, head.MustGet())
	head = sliceutils.Head[int](nil)
	require.False(t, head.IsJust())
	head = sliceutils.Head([]int{})
	require.False(t, head.IsJust())
}

func Test_Intersection(t *testing.T) {
	left := []int{1, 2, 3, 4}
	right := []int{2, 4, 6}
	actual := sliceutils.Intersection(left, right)
	expected := []int{2, 4}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Intersection_empty(t *testing.T) {
	left := []int{1, 2, 3, 4}
	right := []int{5, 6, 7}
	actual := sliceutils.Intersection(left, right)
	require.NotNil(t, actual.MustGet())
	require.Empty(t, actual.MustGet())
}

func Test_LenAll(t *testing.T) {
	breakfast := []string{
		`tea`, `eggs`, `toast`,
	}
	lunch := []string{
		`coke`, `chicken`, `salad`,
	}
	dinner := []string{
		`hamburger`, `salad`, `beer`,
	}
	actual := sliceutils.LenAll(breakfast, lunch, dinner)
	expected := len(breakfast) + len(lunch) + len(dinner)
	require.Equal(t, expected, actual)
	actual = sliceutils.LenAll[string](make([]string, 0))
	expected = 0
	require.Equal(t, expected, actual)
	actual = sliceutils.LenAll[string]()
	expected = 0
	require.Equal(t, expected, actual)
	actual = sliceutils.LenAll[string](nil)
	expected = 0
	require.Equal(t, expected, actual)
}

func Test_Map(t *testing.T) {
	colors := []string{
		`red`,
		`green`,
		`blue`,
	}
	mapper := func(value string) result.Result[int] {
		return result.Ok(len(value))
	}
	actual := sliceutils.Map(mapper, colors)
	require.NoError(t, actual.Error())
	expected := make([]int, 0, len(colors))
	for _, color := range colors {
		length := mapper(color)
		require.NoError(t, length.Error())
		expected = append(expected, length.MustGet())
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_MapNoError(t *testing.T) {
	colors := []string{
		`red`,
		`green`,
		`blue`,
	}
	mapper := func(value string) int {
		return len(value)
	}
	actual := sliceutils.MapNoError(mapper, colors)
	expected := make([]int, 0, len(colors))
	for _, color := range colors {
		length := mapper(color)
		expected = append(expected, length)
	}
	require.Equal(t, expected, actual)
}

func Test_MapNoErrorResult(t *testing.T) {
	colors := []string{
		`red`,
		`green`,
		`blue`,
	}
	mapper := func(value string) int {
		return len(value)
	}
	colorsResult := result.Ok(colors)
	actual := sliceutils.MapNoErrorResult(mapper, colorsResult)
	require.NoError(t, actual.Error())
	expected := make([]int, 0, len(colors))
	for _, color := range colors {
		length := mapper(color)
		expected = append(expected, length)
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Map_error(t *testing.T) {
	colors := []string{
		`red`,
		`green`,
		`blue`,
	}
	failedMap := `failed MAP`
	mapper := func(value string) result.Result[int] {
		return result.Error[int](errors.New(failedMap))
	}
	actual := sliceutils.Map(mapper, colors)
	require.ErrorContains(t, actual.Error(), failedMap)
}

func Test_Map_withSelect(t *testing.T) {
	filter := func(value int) result.Result[bool] {
		return result.Ok(value > 20)
	}
	mapper := func(value int) result.Result[int] {
		return result.Ok(value * 2)
	}
	filtered := sliceutils.Select(filter, []int{10, 20, 30, 40, 50})
	actual := sliceutils.MapResult(mapper, filtered)
	require.NoError(t, actual.Error())
	expected := []int{60, 80, 100}
	require.Equal(t, expected, actual.MustGet())
}

func Test_Map_withSelect_error(t *testing.T) {
	expectedMessage := `failed MAP filter`
	filter := newFailingFilter[int](expectedMessage)
	mapper := func(value int) result.Result[int] {
		return result.Ok(value * 2)
	}
	filtered := sliceutils.Select(filter, []int{10, 20, 30, 40, 50})
	actual := sliceutils.MapResult(mapper, filtered)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Prepend(t *testing.T) {
	foods := []string{
		`tea`, `toast`,
	}
	actual := sliceutils.Prepend(foods, `eggs`)
	expected := make([]string, 0, len(foods)+1)
	expected = append(expected, `eggs`)
	expected = append(expected, foods...)
	require.Equal(t, expected, actual)
	actual = make([]string, 0, 1)
	actual = sliceutils.Prepend(actual, `eggs`)
	expected = []string{
		`eggs`,
	}
	require.Equal(t, expected, actual)
	actual = nil
	actual = sliceutils.Prepend(actual, `eggs`)
	expected = []string{
		`eggs`,
	}
	require.Equal(t, expected, actual)
}

func Test_Remove(t *testing.T) {
	food := `eggs`
	foods := []string{
		`tea`, food, `toast`, food,
	}
	actual := sliceutils.Remove(foods, food)
	expected := []string{
		`tea`, `toast`, food,
	}
	require.Equal(t, expected, actual)
	actual = sliceutils.Remove(actual, food)
	expected = []string{
		`tea`, `toast`,
	}
	require.Equal(t, expected, actual)
	actual = sliceutils.Remove(actual, food)
	require.Equal(t, expected, actual)
	actual = make([]string, 0)
	actual = sliceutils.Remove(actual, food)
	require.NotNil(t, actual)
	require.Empty(t, actual)
	actual = sliceutils.Remove(nil, food)
	require.Nil(t, actual)
}

func Test_Remove_pointer(t *testing.T) {
	cheese := `cheese`
	ham := `ham`
	foods := sliceutils.ToSlice(cheese, ham)
	actual := sliceutils.Remove(foods, cheese)
	require.Equal(t, sliceutils.ToSlice(ham), actual)
}

func Test_Reverse(t *testing.T) {
	scores := []int{10, 20, 30}
	actual := sliceutils.Reverse(scores)
	expected := []int{30, 20, 10}
	require.Equal(t, expected, actual)
	expectedScores := []int{10, 20, 30}
	require.Equal(t, expectedScores, scores)
	emptySlice := sliceutils.Reverse[any](nil)
	require.NotNil(t, emptySlice)
	require.Empty(t, emptySlice)
}

func Test_Select(t *testing.T) {
	values := []string{
		`friendly`, `bright`, `happily`,
		`calmly`, `slowly`, `fast`,
	}
	filter := func(value string) result.Result[bool] {
		return result.Ok(!strings.HasSuffix(value, `ly`))
	}
	expected := []string{
		`bright`,
		`fast`,
	}
	actual := sliceutils.Select(filter, values)
	require.NoError(t, actual.Error())
	require.Equal(t, expected, actual.MustGet())
	filter = func(value string) result.Result[bool] {
		return result.Ok(strings.HasSuffix(value, `ed`))
	}
	actual = sliceutils.Select(filter, values)
	require.NoError(t, actual.Error())
	require.Empty(t, actual.MustGet())
	require.NotNil(t, actual.MustGet())
}

func Test_Select_error(t *testing.T) {
	values := []string{
		`abc`,
	}
	expectedMessage := `failing SELECT filter`
	filter := newFailingFilter[string](expectedMessage)
	actual := sliceutils.Select(filter, values)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Select_nil(t *testing.T) {
	actual := sliceutils.Select(nil, []int{10, 20, 30})
	require.NoError(t, actual.Error())
	require.NotNil(t, actual.MustGet())
	require.Empty(t, actual.MustGet())
}

func Test_SelectNoError(t *testing.T) {
	values := []string{
		`friendly`, `bright`, `happily`,
		`calmly`, `slowly`, `fast`,
	}
	filter := func(value string) bool {
		return !strings.HasSuffix(value, `ly`)
	}
	expected := []string{
		`bright`,
		`fast`,
	}
	actual := sliceutils.SelectNoError(filter, values)
	require.Equal(t, expected, actual)
	filter = func(value string) bool {
		return strings.HasSuffix(value, `ed`)
	}
	actual = sliceutils.SelectNoError(filter, values)
	require.Empty(t, actual)
	require.NotNil(t, actual)
}

func Test_SelectNoError_nil(t *testing.T) {
	actual := sliceutils.SelectNoError(nil, []int{10, 20, 30})
	require.NotNil(t, actual)
	require.Empty(t, actual)
}

func Test_SelectNoErrorResult(t *testing.T) {
	values := []string{
		`friendly`, `bright`, `happily`,
		`calmly`, `slowly`, `fast`,
	}
	filter := func(value string) bool {
		return !strings.HasSuffix(value, `ly`)
	}
	expected := []string{
		`bright`,
		`fast`,
	}
	actual := sliceutils.SelectNoErrorResult(filter, result.Ok(values))
	require.NoError(t, actual.Error())
	require.Equal(t, expected, actual.MustGet())
	filter = func(value string) bool {
		return strings.HasSuffix(value, `ed`)
	}
	actual = sliceutils.SelectNoErrorResult(filter, result.Ok(values))
	require.NoError(t, actual.Error())
	require.Empty(t, actual.MustGet())
	require.NotNil(t, actual.MustGet())
}

func Test_SelectResult(t *testing.T) {
	values := []string{
		`friendly`, `bright`, `happily`,
		`calmly`, `slowly`, `fast`,
	}
	filter := func(value string) result.Result[bool] {
		return result.Ok(!strings.HasSuffix(value, `ly`))
	}
	expected := []string{
		`bright`,
		`fast`,
	}
	actual := sliceutils.SelectResult(filter, result.Ok(values))
	require.NoError(t, actual.Error())
	require.Equal(t, expected, actual.MustGet())
	filter = func(value string) result.Result[bool] {
		return result.Ok(strings.HasSuffix(value, `ed`))
	}
	actual = sliceutils.Select(filter, values)
	require.NoError(t, actual.Error())
	require.Empty(t, actual.MustGet())
	require.NotNil(t, actual.MustGet())
}

// Set up definitions for Test_SelectWithSubtypeFilter
type itemForFiltering struct {
	value int
}

func (i itemForFiltering) Value() int {
	return i.value
}

func (i itemForFiltering) Other() {
}

// subsetInterface contains some, but not all, of what's available in
// itemForFiltering.
type subsetInterface interface {
	Value() int
}

func Test_SelectWithSubtypeFilter(t *testing.T) {
	zeroErrorMsg := "0 throws an error"
	filterFn := func(si subsetInterface) result.Result[bool] {
		if si.Value() == 0 {
			return result.Error[bool](errors.New(zeroErrorMsg))
		}
		answer := si.Value() < 10
		return result.Ok(answer)
	}
	values := []itemForFiltering{
		{6},
		{15},
		{3},
		{72},
		{-5},
		{0},
		{4},
		{12},
	}
	actual := sliceutils.SelectUsingSubsetInterfaceFilter[itemForFiltering, subsetInterface](filterFn, values)
	require.True(t, actual.IsError(), "An error should be thrown")
	assert.ErrorContains(t, actual.Error(), zeroErrorMsg)
	valuesSubset := values[0:4]
	actual = sliceutils.SelectUsingSubsetInterfaceFilter[itemForFiltering, subsetInterface](filterFn, valuesSubset)
	expected := []itemForFiltering{{6}, {3}}
	require.False(t, actual.IsError(), "An error should not have been thrown")
	assert.Equal(t, expected, actual.MustGet())
}

func Test_SelectWithSubtypeFilter_fail(t *testing.T) {
	errorMsg := "the filter function does not support the element to be filtered"
	filterFn := func(si int) result.Result[bool] {
		if si == 0 {
			return result.Error[bool](errors.New(errorMsg))
		}
		answer := si < 10
		return result.Ok(answer)
	}
	values := []itemForFiltering{
		{12},
	}
	actual := sliceutils.SelectUsingSubsetInterfaceFilter[itemForFiltering, int](filterFn, values)
	require.True(t, actual.IsError(), "An error should be thrown")
	assert.ErrorContains(t, actual.Error(), errorMsg)
}

func Test_SelectWithSubtypeFilterNoError(t *testing.T) {
	filterFn := func(si subsetInterface) bool {
		return si.Value() < 10
	}
	values := []itemForFiltering{
		{6},
		{15},
		{3},
		{72},
		{-5},
		{0},
		{4},
		{12},
	}
	actual := sliceutils.SelectUsingSubsetInterfaceFilterNoError[itemForFiltering, subsetInterface](filterFn, values)
	require.NoError(t, actual.Error(), "An error should not have been thrown")
	valuesSubset := values[0:4]
	actual = sliceutils.SelectUsingSubsetInterfaceFilterNoError[itemForFiltering, subsetInterface](filterFn, valuesSubset)
	expected := []itemForFiltering{{6}, {3}}
	require.NoError(t, actual.Error(), "An error should not have been thrown")
	assert.Equal(t, expected, actual.MustGet())
}

func Test_SelectWithSubtypeFilterNoError_fail(t *testing.T) {
	errorMsg := "the filter function does not support the element to be filtered"
	filterFn := func(si int) bool {
		return si < 10
	}
	values := []itemForFiltering{
		{12},
	}
	actual := sliceutils.SelectUsingSubsetInterfaceFilterNoError[itemForFiltering, int](filterFn, values)
	require.True(t, actual.IsError(), "An error should be thrown")
	assert.ErrorContains(t, actual.Error(), errorMsg)
}

func Test_Tail(t *testing.T) {
	tail := sliceutils.Tail([]int{1, 2, 3, 4, 5})
	require.True(t, tail.IsJust())
	require.Equal(t, []int{2, 3, 4, 5}, tail.MustGet())
	tail = sliceutils.Tail([]int{1})
	require.True(t, tail.IsJust())
	require.Equal(t, []int{}, tail.MustGet())
	tail = sliceutils.Tail[int](nil)
	require.False(t, tail.IsJust())
	tail = sliceutils.Tail([]int{})
	require.False(t, tail.IsJust())
}

func Test_ToAnySlice(t *testing.T) {
	var value any = []any{
		1, 2, 3,
	}
	actual := sliceutils.AnyToAnySlice(value)
	require.True(t, actual.IsJust())
	expected := []any{
		1, 2, 3,
	}
	require.Equal(t, expected, actual.MustGet())
}

func Test_ToAnySlice_fail(t *testing.T) {
	actual := sliceutils.AnyToAnySlice(1) // Not a slice.
	require.False(t, actual.IsJust())
	numbers := sliceutils.ToSlice(1, 2, 3)
	actual = sliceutils.AnyToAnySlice(numbers) // Not an any slice.
	require.False(t, actual.IsJust())
	actual = sliceutils.AnyToAnySlice(nil) // Not a slice.
	require.False(t, actual.IsJust())
}

func Test_ToSlice(t *testing.T) {
	actual := sliceutils.ToSlice(1, 2, 3)
	expected := []int{1, 2, 3}
	require.Equal(t, expected, actual)
	actual = sliceutils.ToSlice(1)
	expected = []int{1}
	require.Equal(t, expected, actual)
	actual = sliceutils.ToSlice[int]()
	require.Nil(t, actual)
}

func Test_Union(t *testing.T) {
	breakfast := []string{
		`coffee`, `eggs`, `toast`,
	}
	lunch := []string{
		`coffee`, `chicken`, `salad`,
	}
	dinner := []string{
		`hamburger`, `salad`, `beer`,
	}
	actual := sliceutils.Union(breakfast, lunch, dinner)
	expected := []string{
		`beer`, `chicken`, `coffee`,
		`coffee`, `eggs`, `hamburger`,
		`salad`, `salad`, `toast`,
	}
	require.NoError(t, actual.Error())
	actualStrings := actual.MustGet()
	sort.Strings(actualStrings)
	require.Equal(t, expected, actualStrings)
}

func Test_UniqueUnion(t *testing.T) {
	breakfast := []string{
		`coffee`, `eggs`, `toast`,
	}
	lunch := []string{
		`coffee`, `chicken`, `salad`,
	}
	dinner := []string{
		`hamburger`, `salad`, `beer`,
	}
	actual := sliceutils.UniqueUnion(breakfast, lunch, dinner)
	expected := []string{
		`beer`, `chicken`, `coffee`,
		`eggs`, `hamburger`, `salad`,
		`toast`,
	}
	require.NoError(t, actual.Error())
	actualStrings := actual.MustGet()
	sort.Strings(actualStrings)
	require.Equal(t, expected, actualStrings)
}

func newFailingFilter[T any](message string) filters.Filter[T] {
	return func(value T) result.Result[bool] {
		return result.Error[bool](errors.New(message))
	}
}
