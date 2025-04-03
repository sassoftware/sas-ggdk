// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package di_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/di"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

type aCloser struct {
	closed bool
}

func (a *aCloser) Close() error {
	if a.closed == true {
		return errors.New("closed twice")
	}
	a.closed = true
	return nil
}

func (a *aCloser) String() string {
	return "aCloser"
}

type bCloser struct {
	a      *aCloser
	closed bool
}

func (b *bCloser) Close() error {
	// Note: Do not Close a because it was plumbed with Get. The framework
	// will call Close.
	if b.closed == true {
		return errors.New("closed twice")
	}
	b.closed = true
	return nil
}

type cStruct struct{}

type failingCloser struct{}

func (f *failingCloser) Close() error {
	return errors.New("failed")
}

func TestSimple(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	a := di.Get[*aCloser]("aCloser")
	require.NoError(t, a.Error())
	require.NotNil(t, a)
	err = stopfn()
	require.NoError(t, err)
	require.True(t, a.MustGet().closed)
}

func TestComposite(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	aFactory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	bFactory := func() result.Result[any] {
		a := di.Get[*aCloser]("aCloser")
		if a.Error() != nil {
			return result.Error[any](a.Error())
		}
		return result.Ok[any](&bCloser{
			a: a.MustGet(),
		})
	}
	di.RegisterLazySingletonFactory("aCloser", aFactory)
	di.RegisterLazySingletonFactory("bCloser", bFactory)
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	b := di.Get[*bCloser]("bCloser")
	require.NoError(t, b.Error())
	require.NotNil(t, b.MustGet())
	a := di.Get[*aCloser]("aCloser")
	require.NoError(t, a.Error())
	require.NotNil(t, a.MustGet())
	err = stopfn()
	require.NoError(t, err)
	require.True(t, b.MustGet().closed)
	require.True(t, b.MustGet().a.closed)
	require.Equal(t, a.MustGet(), b.MustGet().a)
}

func TestGetWithoutStart(t *testing.T) {
	a := di.Get[*aCloser]("aCloser")
	require.Error(t, a.Error())
}

func TestNotFound(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	_, err := di.Start()
	require.NoError(t, err)
	a := di.Get[*aCloser]("aCloser")
	require.Error(t, a.Error())
}

func TestReplaceNotAllowed(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	di.RegisterLazySingletonFactory("aCloser", factory)
	stopfn, err := di.Start()
	require.Error(t, err)
	require.Nil(t, stopfn)
}

func TestReplaceAllowed(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	di.RegisterLazySingletonFactory("aCloser", factory)
	_, err := di.StartAllowReplaced()
	require.NoError(t, err)
}

func TestMultipleStart(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	_, err := di.Start()
	require.NoError(t, err)
	_, err = di.StartAllowReplaced()
	require.Error(t, err)
	_, err = di.Start()
	require.Error(t, err)
	_, err = di.StartAllowReplaced()
	require.Error(t, err)
}

func TestMultipleStop(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	err = stopfn()
	require.NoError(t, err)
	err = stopfn()
	require.Error(t, err)
}

func TestMultipleClose(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	aAsCloser := di.Get[io.Closer]("aCloser")
	require.NoError(t, aAsCloser.Error())
	require.NotNil(t, aAsCloser.MustGet())
	err = aAsCloser.MustGet().Close()
	require.NoError(t, err)
	err = stopfn()
	require.Error(t, err)
}

func TestGetByInterface(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	_, err := di.Start()
	require.NoError(t, err)
	a := di.Get[*aCloser]("aCloser")
	require.NoError(t, a.Error())
	require.NotNil(t, a)
	aAsStringer := di.Get[fmt.Stringer]("aCloser")
	require.NoError(t, aAsStringer.Error())
	require.NotNil(t, aAsStringer.MustGet())
	s := aAsStringer.MustGet().String()
	require.Equal(t, "aCloser", s)
}

func TestGetByInterfaceNotFound(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&aCloser{})
	}
	di.RegisterLazySingletonFactory("aCloser", factory)
	_, err := di.Start()
	require.NoError(t, err)
	a := di.Get[io.Reader]("aCloser")
	require.Error(t, a.Error())
}

func TestFailingCloser(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&failingCloser{})
	}
	di.RegisterLazySingletonFactory("failingCloser", factory)
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	f := di.Get[io.Closer]("failingCloser")
	require.NoError(t, f.Error())
	require.NotNil(t, f.MustGet())
	err = stopfn()
	require.Error(t, err)
}

func TestNonCloser(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Ok[any](&cStruct{})
	}
	di.RegisterLazySingletonFactory("cStruct", factory)
	stopfn, err := di.Start()
	require.NoError(t, err)
	require.NotNil(t, stopfn)
	c := di.Get[*cStruct]("cStruct")
	require.NoError(t, c.Error())
	require.NotNil(t, c.MustGet())
	err = stopfn()
	require.NoError(t, err)
}

func TestFailingFactory(t *testing.T) {
	defer func() { err := di.Reset(); require.NoError(t, err) }()
	factory := func() result.Result[any] {
		return result.Error[any](errors.New("failed"))
	}
	di.RegisterLazySingletonFactory("fail", factory)
	_, err := di.Start()
	require.NoError(t, err)
	a := di.Get[any]("fail")
	require.Error(t, a.Error())
}
