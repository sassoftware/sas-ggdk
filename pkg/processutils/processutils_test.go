// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package processutils_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/processutils"
	"github.com/stretchr/testify/require"
)

func Test_Whoami(t *testing.T) {
	value := processutils.ProcessName()
	require.False(t, value.IsError())
}
