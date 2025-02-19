// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lazy_test

import (
	"fmt"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/lazy"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

type ExpensiveStruct struct {
}

func (rw *ExpensiveStruct) Dispose() error {
	return nil
}

func (rw *ExpensiveStruct) String() string {
	return "expensive struct"
}

type ThingData struct {
	creationCount int
	disposalCount int
	stringCount   int
}

type Thing struct {
	ThingData
	getExpensiveStruct lazy.Creator[*ExpensiveStruct]
}

func NewThing() *Thing {
	t := &Thing{}
	t.getExpensiveStruct = lazy.MakeGetter(
		t.createExpensiveStruct,
	)
	return t
}

func (t *Thing) createExpensiveStruct() *ExpensiveStruct {
	t.creationCount++
	return &ExpensiveStruct{}
}

func (t *Thing) String() string {
	t.stringCount++
	return fmt.Sprintf("Value: %v", t.getExpensiveStruct())
}

type ThingWithError struct {
	ThingData
	getRepositoryWarehouse     lazy.Creator[result.Result[*ExpensiveStruct]]
	disposeRepositoryWarehouse lazy.Disposal[error]
}

func NewThingWithError() *ThingWithError {
	t := &ThingWithError{}
	t.getRepositoryWarehouse, t.disposeRepositoryWarehouse = lazy.MakeGetterWithDispose(
		t.createExpensiveStruct,
		t.disposeOfExpensiveStruct,
	)
	return t
}

func (t *ThingWithError) createExpensiveStruct() result.Result[*ExpensiveStruct] {
	t.creationCount++
	return result.Ok(&ExpensiveStruct{})
}

func (t *ThingWithError) disposeOfExpensiveStruct(rw result.Result[*ExpensiveStruct]) error {
	t.disposalCount++
	if rw.IsError() == false {
		return rw.MustGet().Dispose()
	}
	return nil
}

func (t *ThingWithError) Dispose() error {
	return t.disposeRepositoryWarehouse()
}
func (t *ThingWithError) String() string {
	t.stringCount++
	value := t.getRepositoryWarehouse()
	if value.IsError() {
		return fmt.Sprintf("Error: %v", value.Error())
	}
	return fmt.Sprintf("Value: %v", value.MustGet())
}

func Test_LazyInit(t *testing.T) {
	thing := NewThing()
	s := thing.String()
	require.Equal(t, "Value: expensive struct", s)
	s = thing.String()
	_ = s
	s = thing.String()
	_ = s
	s = thing.String()
	_ = s
	require.Equal(t, 1, thing.creationCount)
	require.Equal(t, 4, thing.stringCount)
}

func Test_LazyInitWithError(t *testing.T) {
	thing := NewThingWithError()
	s := thing.String()
	require.Equal(t, "Value: expensive struct", s)
	s = thing.String()
	_ = s
	s = thing.String()
	_ = s
	s = thing.String()
	_ = s
	require.Equal(t, 1, thing.creationCount)
	require.Equal(t, 4, thing.stringCount)
	_ = thing.Dispose()
	_ = thing.Dispose()
	_ = thing.Dispose()
	require.Equal(t, 1, thing.disposalCount)
}
