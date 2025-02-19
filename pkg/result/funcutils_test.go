// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_CallFuncOk(t *testing.T) {
	a := 0
	f := func() {
		a = 1
	}
	result.CallFunc(result.Ok(f))
	require.Equal(t, 1, a)
}

func Test_CallFuncErr(t *testing.T) {
	a := 0
	f := func() {
		a = 1
	}
	result.CallFunc(result.New(f, errors.New("failed")))
	require.Equal(t, 0, a)
}

func Test_CallErrorOnlyFuncOk(t *testing.T) {
	a := 0
	f := func() error {
		a = 1
		return nil
	}
	err := result.CallErrorOnlyFunc(result.Ok(f))
	require.NoError(t, err)
	require.Equal(t, 1, a)
}

func Test_CallErrorOnlyFuncErrResult(t *testing.T) {
	a := 0
	f := func() error {
		a = 1
		return nil
	}
	err := result.CallErrorOnlyFunc(result.New(f, errors.New("failed")))
	require.ErrorContains(t, err, "failed")
	require.Equal(t, 0, a)
}

func Test_CallErrorOnlyFuncErrFunc(t *testing.T) {
	a := 0
	f := func() error {
		a = 1
		return errors.New("failed")
	}
	err := result.CallErrorOnlyFunc(result.Ok(f))
	require.ErrorContains(t, err, "failed")
	require.Equal(t, 1, a)
}

func Test_CallFlatFuncOk(t *testing.T) {
	a := 0
	f := func() result.Result[int] {
		a = 1
		return result.Ok(a + 1)
	}
	actual := result.CallFlatFunc(result.Ok(f))
	require.NoError(t, actual.Error())
	require.Equal(t, 1, a)
	require.Equal(t, 2, actual.MustGet())
}

func Test_CallFlatFuncErrResult(t *testing.T) {
	a := 0
	f := func() result.Result[int] {
		a = 1
		return result.Ok(a + 1)
	}
	actual := result.CallFlatFunc(result.New(f, errors.New("failed")))
	require.ErrorContains(t, actual.Error(), "failed")
	require.Equal(t, 0, a)
}

func Test_CallFaltFuncErrFunc(t *testing.T) {
	a := 0
	f := func() result.Result[int] {
		a = 1
		return result.Error[int](errors.New("failed"))
	}
	actual := result.CallFlatFunc(result.Ok(f))
	require.ErrorContains(t, actual.Error(), "failed")
	require.Equal(t, 1, a)
}
