// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_FlatMap_Ok(t *testing.T) {
	instance := result.Ok(10)
	mapped := result.FlatMap(func(i int) result.Result[string] {
		return result.Ok(fmt.Sprintf("%v", i))
	}, instance)
	validateOkResult(t, mapped, "10", "failed")
}

func Test_FlatMap_Ok_To_Failed(t *testing.T) {
	instance := result.Ok(10)
	mapped := result.FlatMap(func(i int) result.Result[string] {
		return result.Error[string](errors.New("failed"))
	}, instance)
	validateErrResult(t, mapped, "else")
}

func Test_FlatMap_Failed(t *testing.T) {
	instance := result.Error[int](errors.New("failed"))
	called := false
	mapped := result.FlatMap(func(i int) result.Result[string] {
		called = true
		return result.Ok("ok")
	}, instance)
	require.False(t, called)
	validateErrResult(t, mapped, "else")
}

func Test_FlatMapWithErrorFn(t *testing.T) {
	mapper := result.MakeFlatMapper(strconv.Atoi)
	instance := result.Ok("123")
	mapped := result.FlatMap(mapper, instance)
	validateOkResult(t, mapped, 123, 0)

	instance = result.Ok("not a number")
	mapped = result.FlatMap(mapper, instance)
	validateErrResult(t, mapped, 0)
}

func append2WithError(a, b string) (string, error) {
	return a + b, nil
}

var flatAppend2 = result.MakeFlatMapper2(append2WithError)

func Test_FlatMap2_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	actual := result.FlatMap2(flatAppend2, instance1, instance2)
	validateOkResult(t, actual, "value1value2", "")
}

func Test_FlatMap2_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	actual := result.FlatMap2(flatAppend2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap2_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	actual := result.FlatMap2(flatAppend2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func append3WithError(a, b, c string) (string, error) {
	return a + b + c, nil
}

var flatAppend3 = result.MakeFlatMapper3(append3WithError)

func Test_FlatMap3_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.FlatMap3(flatAppend3, instance1, instance2, instance3)
	validateOkResult(t, actual, "value1value2value3", "")
}

func Test_FlatMap3_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.FlatMap3(flatAppend3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap3_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	instance3 := result.Ok("value3")
	actual := result.FlatMap3(flatAppend3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap3_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error"))
	actual := result.FlatMap3(flatAppend3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func append4WithError(a, b, c, d string) (string, error) {
	return a + b + c + d, nil
}

var flatAppend4 = result.MakeFlatMapper4(append4WithError)

func Test_FlatMap4_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.FlatMap4(flatAppend4, instance1, instance2, instance3, instance4)
	validateOkResult(t, actual, "value1value2value3value4", "")
}

func Test_FlatMap4_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.FlatMap4(flatAppend4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap4_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.FlatMap4(flatAppend4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap4_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error"))
	instance4 := result.Ok("value4")
	actual := result.FlatMap4(flatAppend4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap4_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error"))
	actual := result.FlatMap4(flatAppend4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func append5WithError(a, b, c, d, e string) (string, error) {
	return a + b + c + d + e, nil
}

var flatAppend5 = result.MakeFlatMapper5(append5WithError)

func Test_FlatMap5_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateOkResult(t, actual, "value1value2value3value4value5", "")
}

func Test_FlatMap5_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap5_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap5_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error"))
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap5_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error"))
	instance5 := result.Ok("value5")
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_FlatMap5_Error5(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Error[string](errors.New("error"))
	actual := result.FlatMap5(flatAppend5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}
