// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe_test

import (
	"fmt"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/stretchr/testify/require"
)

func Test_New_Nothing(t *testing.T) {
	instance := maybe.Nothing[bool]()
	validateNothingMaybe(t, instance, true)
}

func Test_Map_Noting(t *testing.T) {
	instance := maybe.Nothing[string]()
	called := false
	instance = maybe.Map(func(a string) string {
		called = true
		return "mapped " + a
	}, instance)
	require.False(t, called)
	validateNothingMaybe(t, instance, "else value")
}

func Test_FlatMap_Nothing(t *testing.T) {
	instance := maybe.Nothing[int]()
	called := false
	mapped := maybe.FlatMap(func(_ int) maybe.Maybe[string] {
		called = true
		return maybe.Just("ok")
	}, instance)
	require.False(t, called)
	validateNothingMaybe(t, mapped, "else")
}

func Test_Nothing_String(t *testing.T) {
	instance := maybe.Nothing[int]()
	require.Equal(t, "{Nothing}", fmt.Sprintf("%v", instance))
}

func Test_Nothing_Is_Just(t *testing.T) {
	instance := maybe.Nothing[int]()
	require.False(t, maybe.IsJust(instance))
}
