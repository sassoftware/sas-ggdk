// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/stretchr/testify/require"
)

const failMe = "fail me"

func Test_MakeFlatMapperNoError(t *testing.T) {
	f := func(s string) int {
		return len(s)
	}
	rm := result.MakeFlatMapperNoError(f)
	actual := rm("a string")
	require.NoError(t, actual.Error())
	require.Equal(t, 8, actual.MustGet())
}

func Test_MakeFlatMapperNoError2(t *testing.T) {
	f := func(a, b string) string {
		return a + b
	}
	rm := result.MakeFlatMapperNoError2(f)
	actual := rm("value1", "value2")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2", actual.MustGet())
}

func Test_MakeFlatMapperNoError3(t *testing.T) {
	f := func(a, b, c string) string {
		return a + b + c
	}
	rm := result.MakeFlatMapperNoError3(f)
	actual := rm("value1", "value2", "value3")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3", actual.MustGet())
}

func Test_MakeFlatMapperNoError4(t *testing.T) {
	f := func(a, b, c, d string) string {
		return a + b + c + d
	}
	rm := result.MakeFlatMapperNoError4(f)
	actual := rm("value1", "value2", "value3", "value4")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3value4", actual.MustGet())
}

func Test_MakeFlatMapperNoError5(t *testing.T) {
	f := func(a, b, c, d, e string) string {
		return a + b + c + d + e
	}
	rm := result.MakeFlatMapperNoError5(f)
	actual := rm("value1", "value2", "value3", "value4", "value5")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3value4value5", actual.MustGet())
}

func Test_MakeFlatMapper(t *testing.T) {
	f := func(s string) (int, error) {
		if s == failMe {
			return 0, errors.New("I failed")
		}
		return len(s), nil
	}
	rm := result.MakeFlatMapper(f)
	actual := rm("a string")
	require.NoError(t, actual.Error())
	require.Equal(t, 8, actual.MustGet())
	actual = rm(failMe)
	require.ErrorContains(t, actual.Error(), "I failed")
}

func Test_MakeFlatMapper2(t *testing.T) {
	f := func(a, b string) (string, error) {
		if a == failMe {
			return "", errors.New("I failed")
		}
		return a + b, nil
	}
	rm := result.MakeFlatMapper2(f)
	actual := rm("value1", "value2")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2", actual.MustGet())
	actual = rm(failMe, "")
	require.ErrorContains(t, actual.Error(), "I failed")
}

func Test_MakeFlatMapper3(t *testing.T) {
	f := func(a, b, c string) (string, error) {
		if a == failMe {
			return "", errors.New("I failed")
		}
		return a + b + c, nil
	}
	rm := result.MakeFlatMapper3(f)
	actual := rm("value1", "value2", "value3")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3", actual.MustGet())
	actual = rm(failMe, "", "")
	require.ErrorContains(t, actual.Error(), "I failed")
}

func Test_MakeFlatMapper4(t *testing.T) {
	f := func(a, b, c, d string) (string, error) {
		if a == failMe {
			return "", errors.New("I failed")
		}
		return a + b + c + d, nil
	}
	rm := result.MakeFlatMapper4(f)
	actual := rm("value1", "value2", "value3", "value4")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3value4", actual.MustGet())
	actual = rm(failMe, "", "", "")
	require.ErrorContains(t, actual.Error(), "I failed")
}

func Test_MakeFlatMapper5(t *testing.T) {
	f := func(a, b, c, d, e string) (string, error) {
		if a == failMe {
			return "", errors.New("I failed")
		}
		return a + b + c + d + e, nil
	}
	rm := result.MakeFlatMapper5(f)
	actual := rm("value1", "value2", "value3", "value4", "value5")
	require.NoError(t, actual.Error())
	require.Equal(t, "value1value2value3value4value5", actual.MustGet())
	actual = rm(failMe, "", "", "", "")
	require.ErrorContains(t, actual.Error(), "I failed")
}
