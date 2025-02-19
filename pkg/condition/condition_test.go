// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package condition_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/condition"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"github.com/stretchr/testify/require"
)

type Parent interface {
	GetChildren() *Children
}

type Person struct {
	name     string
	children *Children
}

func (p *Person) GetChildren() *Children {
	return p.children
}

type Children struct {
	family []*Person
}

func newPerson(name string, childNames ...string) result.Result[*Person] {
	instance := &Person{
		name:     name,
		children: new(Children),
	}
	count := len(childNames)
	if count > 0 {
		mapper := func(name string) result.Result[*Person] {
			return newPerson(name)
		}
		family := sliceutils.Map(mapper, childNames)
		if family.IsError() {
			return result.Error[*Person](family.Error())
		}
		instance.children.family = family.MustGet()
	}
	return result.Ok(instance)
}

func Test_IsNil_nil(t *testing.T) {
	state := condition.IsNil(nil)
	require.True(t, state)
}

func Test_IsNil_bool(t *testing.T) {
	var b bool
	require.NotNil(t, b)
	state := condition.IsNil(b)
	require.False(t, state)
	b = true
	require.NotNil(t, b)
	state = condition.IsNil(b)
	require.False(t, state)
}

func Test_IsNil_byte(t *testing.T) {
	var b byte
	require.NotNil(t, b)
	state := condition.IsNil(b)
	require.False(t, state)
	b = byte('a')
	require.NotNil(t, b)
	state = condition.IsNil(b)
	require.False(t, state)
}

func Test_IsNil_complex(t *testing.T) {
	var cpx complex128
	require.NotNil(t, cpx)
	state := condition.IsNil(cpx)
	require.False(t, state)
	cpx = complex(1.0, 2.0)
	require.NotNil(t, cpx)
	state = condition.IsNil(cpx)
	require.False(t, state)
}

func Test_IsNil_float(t *testing.T) {
	var f float32
	require.NotNil(t, f)
	state := condition.IsNil(f)
	require.False(t, state)
	f = 3.14
	require.NotNil(t, f)
	state = condition.IsNil(f)
	require.False(t, state)
}

func Test_IsNil_func(t *testing.T) {
	var f func(*testing.T)
	require.Nil(t, f)
	state := condition.IsNil(f)
	require.True(t, state)
	f = Test_IsNil_func
	require.NotNil(t, f)
	state = condition.IsNil(f)
	require.False(t, state)
}

func Test_IsNil_int(t *testing.T) {
	var i int
	require.NotNil(t, i)
	state := condition.IsNil(i)
	require.False(t, state)
	i = 0
	require.NotNil(t, i)
	state = condition.IsNil(i)
	require.False(t, state)
}

func Test_IsNil_interface_nil(t *testing.T) {
	f := func(parent Parent) {
		if parent == nil {
			require.Fail(t, `Unexpected equality with nil`)
		}
		require.Nil(t, parent)
		f2 := func() {
			defer func() {
				recovered := recover()
				err, ok := recovered.(error)
				require.True(t, ok)
				require.Error(t, err)
				require.EqualError(t, err, `runtime error: invalid memory address or nil pointer dereference`)
			}()
			parent.GetChildren()
		}
		f2()
		state := condition.IsNil(parent)
		require.True(t, state)
	}
	var person *Person // Zero value of nil.
	f(person)
}

func Test_IsNil_interface_notNil(t *testing.T) {
	f := func(parent Parent) {
		if parent == nil {
			require.Fail(t, `Unexpected equality with nil`)
		}
		require.NotNil(t, parent)
		children := parent.GetChildren()
		require.NotNil(t, children)
		require.Len(t, children.family, 1)
		state := condition.IsNil(parent)
		require.False(t, state)
	}
	person := newPerson(`John`, `Mary`)
	require.NoError(t, person.Error())
	f(person.MustGet())
}

func Test_IsNil_map(t *testing.T) {
	var m map[string]string
	require.Nil(t, m)
	state := condition.IsNil(m)
	require.True(t, state)
	m = make(map[string]string)
	require.NotNil(t, m)
	state = condition.IsNil(m)
	require.False(t, state)
}

func Test_IsNil_map_readAndWriteWhenNil(t *testing.T) {
	var m1 map[string]string
	require.Nil(t, m1)
	x := m1[`x`]
	require.Empty(t, x)
	f := func() {
		defer func() {
			recovered := recover()
			err, ok := recovered.(error)
			require.True(t, ok)
			require.EqualError(t, err, `assignment to entry in nil map`)
		}()
		m1[`y`] = `abc`
	}
	f()
	state := condition.IsNil(m1)
	require.True(t, state)
}

func Test_IsNil_slice(t *testing.T) {
	var s1 []int
	require.Nil(t, s1)
	state := condition.IsNil(s1)
	require.True(t, state)
	s2 := make([]int, 0, 5)
	require.NotNil(t, s2)
	state = condition.IsNil(s2)
	require.False(t, state)
}

func Test_IsNil_string(t *testing.T) {
	var s string
	require.NotNil(t, s)
	state := condition.IsNil(s)
	require.False(t, state)
	s = ``
	require.NotNil(t, s)
	state = condition.IsNil(s)
	require.False(t, state)
}
