// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/maybe"
	"github.com/stretchr/testify/require"
)

func Test_CallFuncJust(t *testing.T) {
	a := 0
	f := func() {
		a = 1
	}
	maybe.CallFunc(maybe.Just(f))
	require.Equal(t, 1, a)
}

func Test_CallFuncNothing(t *testing.T) {
	maybe.CallFunc(maybe.Nothing[func()]())
}
