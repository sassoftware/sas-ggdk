// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

const value = "value"

func Test_Map_Ok(t *testing.T) {
	instance := result.Ok(value)
	instance = result.Map(func(a string) (string, error) {
		return "mapped " + a, nil
	}, instance)
	validateOkResult(t, instance, "mapped value", "")
}

func Test_Map_Err(t *testing.T) {
	instance := result.Error[string](errors.New("failed"))
	called := false
	instance = result.Map(func(a string) (string, error) {
		called = true
		return "mapped " + a, nil
	}, instance)
	require.False(t, called)
	validateErrResult(t, instance, "else value")
}

func append2(a, b string) (string, error) {
	return a + b, nil
}
func Test_Map2_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	actual := result.Map2(append2, instance1, instance2)
	validateOkResult(t, actual, "value1value2", "")
}

func Test_Map2_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	actual := result.Map2(append2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_Map2_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	actual := result.Map2(append2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func append3(a, b, c string) (string, error) {
	return a + b + c, nil
}
func Test_Map3_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.Map3(append3, instance1, instance2, instance3)
	validateOkResult(t, actual, "value1value2value3", "")
}

func Test_Map3_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.Map3(append3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_Map3_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	instance3 := result.Ok("value3")
	actual := result.Map3(append3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_Map3_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error"))
	actual := result.Map3(append3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func append4(a, b, c, d string) (string, error) {
	return a + b + c + d, nil
}

func Test_Map4_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.Map4(append4, instance1, instance2, instance3, instance4)
	validateOkResult(t, actual, "value1value2value3value4", "")
}

func Test_Map4_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.Map4(append4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error1")
}

func Test_Map4_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.Map4(append4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error2")
}

func Test_Map4_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	actual := result.Map4(append4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error3")
}

func Test_Map4_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	actual := result.Map4(append4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error4")
}

func append5(a, b, c, d, e string) (string, error) {
	return a + b + c + d + e, nil
}

func Test_Map5_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateOkResult(t, actual, "value1value2value3value4value5", "")
}

func Test_Map5_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error1")
}

func Test_Map5_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error2")
}

func Test_Map5_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error3")
}

func Test_Map5_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	instance5 := result.Ok("value5")
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error4")
}

func Test_Map5_Error5(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Error[string](errors.New("error4"))
	actual := result.Map5(append5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error4")
}

func Test_MapNoError_Ok(t *testing.T) {
	instance := result.Ok(value)
	instance = result.MapNoError(func(a string) string {
		return "mapped " + a
	}, instance)
	validateOkResult(t, instance, "mapped value", "")
}

func Test_MapNoError_Err(t *testing.T) {
	instance := result.Error[string](errors.New("failed"))
	called := false
	instance = result.MapNoError(func(a string) string {
		called = true
		return "mapped " + a
	}, instance)
	require.False(t, called)
	validateErrResult(t, instance, "else value")
}

func appendNoError2(a, b string) string {
	return a + b
}
func Test_MapNoError2_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	actual := result.MapNoError2(appendNoError2, instance1, instance2)
	validateOkResult(t, actual, "value1value2", "")
}

func Test_MapNoError2_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	actual := result.MapNoError2(appendNoError2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_MapNoError2_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	actual := result.MapNoError2(appendNoError2, instance1, instance2)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func appendNoError3(a, b, c string) string {
	return a + b + c
}
func Test_MapNoError3_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.MapNoError3(appendNoError3, instance1, instance2, instance3)
	validateOkResult(t, actual, "value1value2value3", "")
}

func Test_MapNoError3_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	actual := result.MapNoError3(appendNoError3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_MapNoError3_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error"))
	instance3 := result.Ok("value3")
	actual := result.MapNoError3(appendNoError3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func Test_MapNoError3_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error"))
	actual := result.MapNoError3(appendNoError3, instance1, instance2, instance3)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error")
}

func appendNoError4(a, b, c, d string) string {
	return a + b + c + d
}

func Test_MapNoError4_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.MapNoError4(appendNoError4, instance1, instance2, instance3, instance4)
	validateOkResult(t, actual, "value1value2value3value4", "")
}

func Test_MapNoError4_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.MapNoError4(appendNoError4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error1")
}

func Test_MapNoError4_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	actual := result.MapNoError4(appendNoError4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error2")
}

func Test_MapNoError4_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	actual := result.MapNoError4(appendNoError4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error3")
}

func Test_MapNoError4_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	actual := result.MapNoError4(appendNoError4, instance1, instance2, instance3, instance4)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error4")
}

func appendNoError5(a, b, c, d, e string) string {
	return a + b + c + d + e
}

func Test_MapNoError5_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateOkResult(t, actual, "value1value2value3value4value5", "")
}

func Test_MapNoError5_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error1")
}

func Test_MapNoError5_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error2")
}

func Test_MapNoError5_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error3")
}

func Test_MapNoError5_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	instance5 := result.Ok("value5")
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error4")
}

func Test_MapNoError5_Error5(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Error[string](errors.New("error5"))
	actual := result.MapNoError5(appendNoError5, instance1, instance2, instance3, instance4, instance5)
	validateErrResult(t, actual, "failed")
	require.Equal(t, actual.Error().Error(), "error5")
}

func Test_MapError_Ok(t *testing.T) {
	instance := result.Ok(value)
	err := result.MapErrorOnly(func(a string) error {
		return nil
	}, instance)
	require.NoError(t, err)
}

func Test_MapError_on_Ok_error(t *testing.T) {
	instance := result.Ok(value)
	err := result.MapErrorOnly(func(a string) error {
		return errors.New("failed")
	}, instance)
	require.ErrorContains(t, err, "failed")
}

func Test_MapError_Err(t *testing.T) {
	instance := result.Error[string](errors.New("failed"))
	called := false
	err := result.MapErrorOnly(func(a string) error {
		called = true
		return errors.New("inner error")
	}, instance)
	require.False(t, called)
	require.ErrorContains(t, err, "failed")
}

func appendError2(a, b string) error {
	return errors.New("inner error")
}
func Test_MapError2_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	err := result.MapErrorOnly2(appendError2, instance1, instance2)
	require.ErrorContains(t, err, "inner error")
}

func Test_MapError2_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	err := result.MapErrorOnly2(appendError2, instance1, instance2)
	require.ErrorContains(t, err, "error1")
}

func Test_MapError2_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	err := result.MapErrorOnly2(appendError2, instance1, instance2)
	require.ErrorContains(t, err, "error2")
}

func appendError3(a, b, c string) error {
	return errors.New("inner error")
}
func Test_MapError3_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	err := result.MapErrorOnly3(appendError3, instance1, instance2, instance3)
	require.ErrorContains(t, err, "inner error")
}

func Test_MapError3_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	err := result.MapErrorOnly3(appendError3, instance1, instance2, instance3)
	require.ErrorContains(t, err, "error1")
}

func Test_MapError3_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	err := result.MapErrorOnly3(appendError3, instance1, instance2, instance3)
	require.ErrorContains(t, err, "error2")
}

func Test_MapError3_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	err := result.MapErrorOnly3(appendError3, instance1, instance2, instance3)
	require.ErrorContains(t, err, "error3")
}

func appendError4(a, b, c, d string) error {
	return errors.New("inner error")
}

func Test_MapError4_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	err := result.MapErrorOnly4(appendError4, instance1, instance2, instance3, instance4)
	require.ErrorContains(t, err, "inner error")
}

func Test_MapError4_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	err := result.MapErrorOnly4(appendError4, instance1, instance2, instance3, instance4)
	require.ErrorContains(t, err, "error1")
}

func Test_MapError4_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	err := result.MapErrorOnly4(appendError4, instance1, instance2, instance3, instance4)
	require.ErrorContains(t, err, "error2")
}

func Test_MapError4_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	err := result.MapErrorOnly4(appendError4, instance1, instance2, instance3, instance4)
	require.ErrorContains(t, err, "error3")
}

func Test_MapError4_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	err := result.MapErrorOnly4(appendError4, instance1, instance2, instance3, instance4)
	require.ErrorContains(t, err, "error4")
}

func appendError5(a, b, c, d, e string) error {
	return errors.New("inner error")
}

func Test_MapError5_Ok(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "inner error")
}

func Test_MapError5_Error1(t *testing.T) {
	instance1 := result.Error[string](errors.New("error1"))
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "error1")
}

func Test_MapError5_Error2(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Error[string](errors.New("error2"))
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "error2")
}

func Test_MapError5_Error3(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Error[string](errors.New("error3"))
	instance4 := result.Ok("value4")
	instance5 := result.Ok("value5")
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "error3")
}

func Test_MapError5_Error4(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Error[string](errors.New("error4"))
	instance5 := result.Ok("value5")
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "error4")
}

func Test_MapError5_Error5(t *testing.T) {
	instance1 := result.Ok("value1")
	instance2 := result.Ok("value2")
	instance3 := result.Ok("value3")
	instance4 := result.Ok("value5")
	instance5 := result.Error[string](errors.New("error5"))
	err := result.MapErrorOnly5(appendError5, instance1, instance2, instance3, instance4, instance5)
	require.ErrorContains(t, err, "error5")
}
