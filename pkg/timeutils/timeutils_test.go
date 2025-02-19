// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package timeutils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sassoftware/sas-ggdk/pkg/timeutils"
	"github.com/stretchr/testify/require"
)

const (
	militaryTimeLayout = `15:04:05`
)

func Test_GetTimestamp(t *testing.T) {
	// The time format is `dd mmm yy hh:mm zzz`, for example `23 Feb 22 09:56 EST`.
	timestamp := timeutils.GetTimestamp()
	ts, err := time.Parse(time.RFC822, timestamp)
	require.NoError(t, err)
	require.NotNil(t, ts)
}

func Test_NewMilitaryTime(t *testing.T) {
	mt := timeutils.NewMilitaryTime()
	require.Len(t, mt, 3)
}

func Test_NewMilitaryTimeFrom(t *testing.T) {
	actual := time.Now()
	mt := timeutils.NewMilitaryTimeFrom(actual)
	require.Len(t, mt, 3)
	value := fmt.Sprintf(`%s:%s:%s`,
		mt.Hours(),
		mt.Minutes(),
		mt.Seconds(),
	)
	expected, err := time.Parse(militaryTimeLayout, value)
	require.NoError(t, err)
	require.Equal(t, expected.Hour(), actual.Hour())
	require.Equal(t, expected.Minute(), actual.Minute())
	require.Equal(t, expected.Second(), actual.Second())
}
