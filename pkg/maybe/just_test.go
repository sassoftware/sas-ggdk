// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/stretchr/testify/require"
)

func Test_New_Just_true(t *testing.T) {
	value := 10
	instance := maybe.Just(value == 10)
	validateJust(t, instance, true, false)
}

func Test_New_Just_false(t *testing.T) {
	value := 10
	instance := maybe.Just(value != 10)
	validateJust(t, instance, false, true)
}

func Test_New_Just_string(t *testing.T) {
	value := "a string value"
	instance := maybe.Just(value)
	validateJust(t, instance, value, "ignored")
}

func Test_New_Just_int(t *testing.T) {
	value := 10
	instance := maybe.Just(value)
	validateJust(t, instance, 10, 0)
}

func Test_Map_Just(t *testing.T) {
	value := "value"
	instance := maybe.Just(value)
	instance = maybe.Map(func(a string) string {
		return "mapped " + a
	}, instance)
	validateJust(t, instance, "mapped value", "")
}

func Test_FlatMap_Just(t *testing.T) {
	instance := maybe.Just(10)
	mapped := maybe.FlatMap(func(i int) maybe.Maybe[string] {
		return maybe.Just(strconv.Itoa(i))
	}, instance)
	validateJust(t, mapped, "10", "failed")
}

func Test_FlatMap_Just_To_Nothing(t *testing.T) {
	instance := maybe.Just(10)
	mapped := maybe.FlatMap(func(_ int) maybe.Maybe[string] {
		return maybe.Nothing[string]()
	}, instance)
	validateNothingMaybe(t, mapped, "else")
}

func Test_Just_String(t *testing.T) {
	instance1 := maybe.Just(10)
	require.Equal(t, "{Just: 10}", fmt.Sprintf("%v", instance1))

	instance2 := maybe.Just([]int{1, 2, 3})
	require.Equal(t, "{Just: []int{1, 2, 3}}", fmt.Sprintf("%v", instance2))

	instance3 := maybe.Just[[]string](nil)
	require.Equal(t, "{Just: []string(nil)}", fmt.Sprintf("%v", instance3))
}
