// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"fmt"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

func Test_New_Error_Bool(t *testing.T) {
	instance := result.Error[bool](errors.New("failed"))
	validateErrResult(t, instance, true)
}

func Test_New_Error_String(t *testing.T) {
	instance := result.Error[string](errors.New("failed"))
	validateErrResult(t, instance, "a different value")
}

func Test_New_Error_Int(t *testing.T) {
	instance := result.Error[int](errors.New("failed"))
	validateErrResult(t, instance, 10)
}

func Test_Error_String(t *testing.T) {
	instance1 := result.Error[int](fmt.Errorf("an int error"))
	require.Equal(t, "{Error: &errors.errorString{s:\"an int error\"}}", fmt.Sprintf("%v", instance1))

	instance2 := result.Error[[]int](fmt.Errorf("a slice of int error"))
	require.Equal(t, "{Error: &errors.errorString{s:\"a slice of int error\"}}", fmt.Sprintf("%v", instance2))

	instance3 := result.Error[[]string](nil)
	require.Equal(t, "{Error: <nil>}", fmt.Sprintf("%v", instance3))
}
