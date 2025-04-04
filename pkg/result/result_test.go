// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func validateOkResult[T any](t *testing.T, r result.Result[T], expected, elseValue T) {
	require.NoError(t, r.Error())
	require.False(t, r.IsError())
	require.Equal(t, expected, r.MustGet())
	require.Equal(t, expected, r.OrElse(elseValue))
}

func validateErrResult[T any](t *testing.T, r result.Result[T], elseValue T) {
	require.Error(t, r.Error())
	require.True(t, r.IsError())
	require.Panics(t, func() { r.MustGet() })
	require.Equal(t, elseValue, r.OrElse(elseValue))
}

func Test_New(t *testing.T) {
	f := func() (int, error) { return 1, nil }
	r := result.New(f())
	validateOkResult(t, r, 1, 0)
	f = func() (int, error) { return 0, errors.New("failed") }
	r = result.New(f())
	validateErrResult(t, r, 10)
}

type named interface {
	GetName() string
}

type instance struct{}

func (i *instance) GetName() string { return "my name" }

func Test_As(t *testing.T) {
	r := result.Ok(&instance{})
	a := result.As[named](r)
	require.NoError(t, a.Error())
	require.Equal(t, "my name", a.MustGet().GetName())
}

func Test_As_Error(t *testing.T) {
	r := result.Ok("string")
	a := result.As[int](r)
	validateErrResult(t, a, 10)
	require.ErrorContains(t, a.Error(), "requested interface not implemented by value")
}

func Test_As_FromError(t *testing.T) {
	r := result.Error[string](errors.New("failed"))
	a := result.As[int](r)
	validateErrResult(t, a, 10)
	require.ErrorContains(t, a.Error(), "failed")
}

func Test_Or_Else_Get(t *testing.T) {
	instance := result.Ok(10)
	value := instance.OrElseGet(func() int { return 100 })
	require.Equal(t, 10, value)

	instance = result.Error[int](errors.New("failed"))
	value = instance.OrElseGet(func() int { return 100 })
	require.Equal(t, 100, value)
}

func Test_ErrorMap_Of_Ok(t *testing.T) {
	instance := result.Ok(10)
	called := false
	mapped := result.ErrorMap(func(e error) error {
		called = true
		return errors.New("failed. caused by: %v", e)
	}, instance)
	validateOkResult(t, mapped, 10, 0)
	require.False(t, called)
}

func Test_ErrorMap_Of_Error(t *testing.T) {
	instance := result.Error[int](errors.New("failed"))
	mapped := result.ErrorMap(func(e error) error {
		return errors.New("failed. caused by: %v", e)
	}, instance)
	validateErrResult(t, mapped, 0)
	require.Equal(t, "failed. caused by: failed", mapped.Error().Error())
}

func Test_FromMaybe(t *testing.T) {
	mOk := maybe.Just(1)
	rOk := result.FromMaybe(mOk, errors.New("failed"))
	require.NoError(t, rOk.Error())
	require.Equal(t, 1, rOk.MustGet())
	mErr := maybe.Nothing[int]()
	rErr := result.FromMaybe(mErr, errors.New("failed"))
	require.Error(t, rErr.Error())
	require.ErrorContains(t, rErr.Error(), "failed")
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
	closerResult := result.Ok(c)
	err := result.Close(closerResult)
	require.NoError(t, err)
	require.True(t, c.called)
}

func Test_CloseNonCloser(t *testing.T) {
	closerResult := result.Ok(99)
	err := result.Close(closerResult)
	require.NoError(t, err)
}

func Test_CloseError(t *testing.T) {
	closerResult := result.Error[int](errors.New("failed"))
	err := result.Close(closerResult)
	require.NoError(t, err)
}
