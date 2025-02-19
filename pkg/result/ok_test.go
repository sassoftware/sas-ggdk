// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"fmt"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_New_Ok_True(t *testing.T) {
	value := 10
	instance := result.Ok(value == 10)
	validateOkResult(t, instance, true, false)
}

func Test_New_Ok_False(t *testing.T) {
	value := 10
	instance := result.Ok(value != 10)
	validateOkResult(t, instance, false, true)
}

func Test_New_Ok_String(t *testing.T) {
	value := "a string value"
	instance := result.Ok(value)
	validateOkResult(t, instance, value, "ignored")
}

func Test_New_Ok_Int(t *testing.T) {
	value := 10
	instance := result.Ok(value)
	validateOkResult(t, instance, 10, 0)
}

func Test_Ok_String(t *testing.T) {
	instance1 := result.Ok(10)
	require.Equal(t, "{Ok: 10}", fmt.Sprintf("%v", instance1))

	instance2 := result.Ok([]int{1, 2, 3})
	require.Equal(t, "{Ok: []int{1, 2, 3}}", fmt.Sprintf("%v", instance2))

	instance3 := result.Ok[[]string](nil)
	require.Equal(t, "{Ok: []string(nil)}", fmt.Sprintf("%v", instance3))
}
