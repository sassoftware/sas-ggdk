// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/stretchr/testify/require"
)

func validateJust[T any](t *testing.T, m maybe.Maybe[T], expected, elseValue T) {
	require.True(t, m.IsJust())
	require.Equal(t, expected, m.OrElse(elseValue))
	require.Equal(t, expected, m.MustGet())
}

func validateNothingMaybe[T any](t *testing.T, m maybe.Maybe[T], elseValue T) {
	require.False(t, m.IsJust())
	require.Equal(t, elseValue, m.OrElse(elseValue))
	require.Panics(t, func() { m.MustGet() })
}

func Test_Or_Else_Get(t *testing.T) {
	instance := maybe.Just(10)
	value := instance.OrElseGet(func() int { return 100 })
	require.Equal(t, 10, value)

	instance = maybe.Nothing[int]()
	value = instance.OrElseGet(func() int { return 100 })
	require.Equal(t, 100, value)
}

func Test_Just_Is_Just(t *testing.T) {
	instance := maybe.Just(10)
	require.True(t, maybe.IsJust(instance))
}
