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

type named interface {
	GetName() string
}

type instance struct{}

func (i *instance) GetName() string { return "my name" }

func Test_As(t *testing.T) {
	m := maybe.Just(&instance{})
	a := maybe.As[named](m)
	require.True(t, a.IsJust())
	require.Equal(t, "my name", a.MustGet().GetName())
}

func Test_As_NotSupported(t *testing.T) {
	m := maybe.Just("string")
	a := maybe.As[int](m)
	require.False(t, a.IsJust())
}

func Test_As_FromNothing(t *testing.T) {
	m := maybe.Nothing[string]()
	a := maybe.As[int](m)
	require.False(t, a.IsJust())
}

type closer struct {
	called bool
}

func (c *closer) Close() error {
	c.called = true
	return nil
}

func Test_CloseCloser(t *testing.T) {
	c := &closer{}
	closerMaybe := maybe.Just(c)
	err := maybe.Close(closerMaybe)
	require.NoError(t, err)
	require.True(t, c.called)
}

func Test_CloseNonCloser(t *testing.T) {
	closerMaybe := maybe.Just(99)
	err := maybe.Close(closerMaybe)
	require.NoError(t, err)
}

func Test_CloseError(t *testing.T) {
	closerMaybe := maybe.Nothing[int]()
	err := maybe.Close(closerMaybe)
	require.NoError(t, err)
}
