// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

var isEvenFilter filters.Filter[int] = func(value int) result.Result[bool] {
	return result.Ok(value%2 == 0)
}

func Test_ApplyFilter(t *testing.T) {
	actual := filters.ApplyFilter(isEvenFilter, 2)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
	actual = filters.ApplyFilter(isEvenFilter, 1)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet())
	actual = filters.ApplyFilter(nil, 2)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
}

func Test_Filter_And(t *testing.T) {
	filter := isEvenFilter.And(newIsGreaterThanFilter(5))
	actual := filter(10)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
	actual = filter(2)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet())
	actual = filter(11)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet())
}

func Test_Filter_And_error(t *testing.T) {
	expectedMessage := `failed AND this filter`
	filter := newFailingFilter[int](expectedMessage).And(isEvenFilter)
	actual := filter(10)
	require.ErrorContains(t, actual.Error(), expectedMessage)
	expectedMessage = `failed AND that filter`
	filter = isEvenFilter.And(newFailingFilter[int](expectedMessage))
	actual = filter(10)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Filter_Or(t *testing.T) {
	filter := isEvenFilter.
		Or(newIsGreaterThanFilter(5))
	actual := filter(2)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
	actual = filter(3)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet())
	actual = filter(7)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
}

func Test_Filter_Or_error(t *testing.T) {
	expectedMessage := `failed OR this filter`
	filter := newFailingFilter[int](expectedMessage).
		Or(isEvenFilter)
	actual := filter(2)
	require.ErrorContains(t, actual.Error(), expectedMessage)
	expectedMessage = `failed OR that filter`
	filter = isEvenFilter.
		Or(newFailingFilter[int](expectedMessage))
	actual = filter(3)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

func Test_Filter_Not(t *testing.T) {
	filter := isEvenFilter.
		Not()
	actual := filter(1)
	require.NoError(t, actual.Error())
	require.True(t, actual.MustGet())
	actual = filter(2)
	require.NoError(t, actual.Error())
	require.False(t, actual.MustGet())
}

func Test_Filter_Not_error(t *testing.T) {
	expectedMessage := `failed NOT filter`
	filter := newFailingFilter[int](expectedMessage).
		Not()
	actual := filter(1)
	require.ErrorContains(t, actual.Error(), expectedMessage)
}

//nolint:funlen
func Test_Filter_table(t *testing.T) {
	type testCaseData[T any] struct {
		filter   filters.Filter[T]
		input    T
		expected bool
		err      bool
	}
	table := map[string]*testCaseData[int]{
		`filter match`: {
			filter:   isEvenFilter,
			input:    10,
			expected: true,
		},
		`filter match failure`: {
			filter:   isEvenFilter,
			input:    11,
			expected: false,
		},
		`filter AND match`: {
			filter: isEvenFilter.
				And(newIsGreaterThanFilter(10)),
			input:    20,
			expected: true,
		},
		`this filter AND failure`: {
			filter: isEvenFilter.
				And(newIsGreaterThanFilter(10)),
			input:    21,
			expected: false,
		},
		`that filter AND failure`: {
			filter: isEvenFilter.
				And(newIsGreaterThanFilter(10)),
			input:    8,
			expected: false,
		},
		`this filter OR match`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10)),
			input:    10,
			expected: true,
		},
		`that filter OR match`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10)),
			input:    11,
			expected: true,
		},
		`filter OR failure`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10)),
			input:    9,
			expected: false,
		},
		`filter NOT match`: {
			filter: isEvenFilter.
				Not(),
			input:    9,
			expected: true,
		},
		`filter NOT failure`: {
			filter: isEvenFilter.
				Not(),
			input:    10,
			expected: false,
		},
		`triple AND filter match`: {
			filter: isEvenFilter.
				And(newIsGreaterThanFilter(10).
					And(newIsGreaterThanFilter(100).Not())),
			input:    50,
			expected: true,
		},
		`triple AND filter failure`: {
			filter: isEvenFilter.
				And(newIsGreaterThanFilter(10).
					And(newIsGreaterThanFilter(100).Not())),
			input:    101,
			expected: false,
		},
		`triple OR filter match 12`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10).Not()).
				Or(newIsGreaterThanFilter(100)),
			input:    12,
			expected: true,
		},
		`triple OR filter match 101`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10).Not()).
				Or(newIsGreaterThanFilter(100)),
			input:    101,
			expected: true,
		},
		`triple OR filter failure`: {
			filter: isEvenFilter.
				Or(newIsGreaterThanFilter(10).Not()).
				Or(newIsGreaterThanFilter(100)),
			input:    49,
			expected: false,
		},
		`filter error`: {
			filter:   newFailingFilter[int](`fail`),
			input:    10,
			expected: false,
			err:      true,
		},
	}
	for name, value := range table {
		subtestFunc := func(t *testing.T) {
			actual := filters.ApplyFilter(value.filter, value.input)
			if value.err {
				require.Error(t, actual.Error())
			} else {
				require.NoError(t, actual.Error())
				require.Equal(t, value.expected, actual.MustGet())
			}
		}
		t.Run(name, subtestFunc)
	}
}

func Test_MatchAllFilter(t *testing.T) {
	expected := result.Ok(true)
	testIntValues(t, filters.MatchAll[int], expected)
	testStringValues(t, filters.MatchAll[string], expected)
}

func Test_MatchNoneFilter(t *testing.T) {
	expected := result.Ok(false)
	testIntValues(t, filters.MatchNone[int], expected)
	testStringValues(t, filters.MatchNone[string], expected)
}

func testIntValues(t *testing.T, fn filters.Filter[int], expected result.Result[bool]) {
	// Test the zero value
	var checkVal int
	actual := fn(checkVal)
	require.Equal(t, expected, actual, "Tested value: %d", checkVal)
	// Test 1000 random integers
	for range 1000 {
		checkVal = rand.Int() //nolint:gosec // Weak random generation is ok when creating test data.
		actual = fn(checkVal)
		require.Equal(t, expected, actual, "Tested value: %d", checkVal)
	}
}

func testStringValues(t *testing.T, fn filters.Filter[string], expected result.Result[bool]) {
	// Test the zero value
	var checkVal string
	actual := fn(checkVal)
	require.Equal(t, expected, actual, "Tested value: %s", checkVal)
	// Test 1000 random strings
	for range 1000 {
		checkVal = strconv.Itoa(rand.Int()) //nolint:gosec // Weak random generation is ok when creating test data.
		actual = fn(checkVal)
		require.Equal(t, expected, actual, "Tested value: %s", checkVal)
	}
}

func newFailingFilter[T any](message string) filters.Filter[T] {
	return func(_ T) result.Result[bool] {
		return result.Error[bool](errors.New(message))
	}
}

func newIsGreaterThanFilter(thatValue int) filters.Filter[int] {
	return func(thisValue int) result.Result[bool] {
		return result.Ok(thisValue > thatValue)
	}
}
