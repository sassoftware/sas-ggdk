// SPDX-FileCopyrightText:  2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package streamutils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/streamutils"
	"github.com/stretchr/testify/require"
)

func Test_CaptureStdStreams(t *testing.T) {
	f := func() result.Result[bool] {
		fmt.Println("stdout content")
		_, _ = fmt.Fprintln(os.Stderr, "stderr content")
		return result.Ok(true)
	}
	res, stdout, stderr := streamutils.CaptureStdStreams(f)
	require.NoError(t, res.Error())
	require.NoError(t, stdout.Error())
	require.NoError(t, stderr.Error())
	require.True(t, res.MustGet())
	require.Equal(t, "stdout content\n", string(stdout.MustGet()))
	require.Equal(t, "stderr content\n", string(stderr.MustGet()))
}
