// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/stretchr/testify/require"
)

var isEvenFilterNoError filters.FilterNoError[int] = func(value int) bool {
	return value%2 == 0
}

func Test_ApplyFilterNoError(t *testing.T) {
	actual := filters.ApplyFilterNoError(isEvenFilterNoError, 2)
	require.True(t, actual)
	actual = filters.ApplyFilterNoError(isEvenFilterNoError, 1)
	require.False(t, actual)
	actual = filters.ApplyFilterNoError(nil, 2)
	require.True(t, actual)
}

func Test_FilterNoError_And(t *testing.T) {
	filter := isEvenFilterNoError.And(newIsGreaterThanFilterNoError(5))
	actual := filter(10)
	require.True(t, actual)
	actual = filter(2)
	require.False(t, actual)
	actual = filter(11)
	require.False(t, actual)
}

func Test_FilterNoError_Or(t *testing.T) {
	filter := isEvenFilterNoError.
		Or(newIsGreaterThanFilterNoError(5))
	actual := filter(2)
	require.True(t, actual)
	actual = filter(3)
	require.False(t, actual)
	actual = filter(7)
	require.True(t, actual)
}

func Test_FilterNoError_Not(t *testing.T) {
	filter := isEvenFilterNoError.
		Not()
	actual := filter(1)
	require.True(t, actual)
	actual = filter(2)
	require.False(t, actual)
}

//nolint:funlen
func Test_FilterNoError_table(t *testing.T) {
	type testCaseData[T any] struct {
		filter   filters.FilterNoError[T]
		input    T
		expected bool
	}
	table := map[string]*testCaseData[int]{
		`filter match`: {
			filter:   isEvenFilterNoError,
			input:    10,
			expected: true,
		},
		`filter match failure`: {
			filter:   isEvenFilterNoError,
			input:    11,
			expected: false,
		},
		`filter AND match`: {
			filter: isEvenFilterNoError.
				And(newIsGreaterThanFilterNoError(10)),
			input:    20,
			expected: true,
		},
		`this filter AND failure`: {
			filter: isEvenFilterNoError.
				And(newIsGreaterThanFilterNoError(10)),
			input:    21,
			expected: false,
		},
		`that filter AND failure`: {
			filter: isEvenFilterNoError.
				And(newIsGreaterThanFilterNoError(10)),
			input:    8,
			expected: false,
		},
		`this filter OR match`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10)),
			input:    10,
			expected: true,
		},
		`that filter OR match`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10)),
			input:    11,
			expected: true,
		},
		`filter OR failure`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10)),
			input:    9,
			expected: false,
		},
		`filter NOT match`: {
			filter: isEvenFilterNoError.
				Not(),
			input:    9,
			expected: true,
		},
		`filter NOT failure`: {
			filter: isEvenFilterNoError.
				Not(),
			input:    10,
			expected: false,
		},
		`triple AND filter match`: {
			filter: isEvenFilterNoError.
				And(newIsGreaterThanFilterNoError(10).
					And(newIsGreaterThanFilterNoError(100).Not())),
			input:    50,
			expected: true,
		},
		`triple AND filter failure`: {
			filter: isEvenFilterNoError.
				And(newIsGreaterThanFilterNoError(10).
					And(newIsGreaterThanFilterNoError(100).Not())),
			input:    101,
			expected: false,
		},
		`triple OR filter match 12`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10).Not()).
				Or(newIsGreaterThanFilterNoError(100)),
			input:    12,
			expected: true,
		},
		`triple OR filter match 101`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10).Not()).
				Or(newIsGreaterThanFilterNoError(100)),
			input:    101,
			expected: true,
		},
		`triple OR filter failure`: {
			filter: isEvenFilterNoError.
				Or(newIsGreaterThanFilterNoError(10).Not()).
				Or(newIsGreaterThanFilterNoError(100)),
			input:    49,
			expected: false,
		},
	}
	for name, value := range table {
		subtestFunc := func(t *testing.T) {
			actual := filters.ApplyFilterNoError(value.filter, value.input)
			require.Equal(t, value.expected, actual)
		}
		t.Run(name, subtestFunc)
	}
}

func Test_MatchAllFilterNoError(t *testing.T) {
	testIntValuesNoError(t, filters.MatchAllNoError[int], true)
	testStringValuesNoError(t, filters.MatchAllNoError[string], true)
}

func Test_MatchNoneFilterNoError(t *testing.T) {
	testIntValuesNoError(t, filters.MatchNoneNoError[int], false)
	testStringValuesNoError(t, filters.MatchNoneNoError[string], false)
}

func testIntValuesNoError(t *testing.T, fn filters.FilterNoError[int], expected bool) {
	// Test the zero value
	var checkVal int
	actual := fn(checkVal)
	require.Equal(t, expected, actual, "Tested value: %d", checkVal)
	// Test 1000 random integers
	for range 1000 {
		checkVal = rand.Int() //nolint:gosec // Weak random generation is ok when generating test data.
		actual = fn(checkVal)
		require.Equal(t, expected, actual, "Tested value: %d", checkVal)
	}
}

func testStringValuesNoError(t *testing.T, fn filters.FilterNoError[string], expected bool) {
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

func newIsGreaterThanFilterNoError(thatValue int) filters.FilterNoError[int] {
	return func(thisValue int) bool {
		return thisValue > thatValue
	}
}
