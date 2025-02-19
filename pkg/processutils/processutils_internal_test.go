// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package processutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExecutableFailed(t *testing.T) {
	executableFunc = func() (string, error) {
		return "", fmt.Errorf("executable failed")
	}
	defer func() { executableFunc = os.Executable }()
	value := ProcessName()
	require.NotEmpty(t, value)
}
