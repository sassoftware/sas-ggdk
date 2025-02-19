// SPDX-FileCopyrightText:  2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer_test

import (
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/pointer"
	"github.com/stretchr/testify/require"
)

func TestPtr(t *testing.T) {
	literal := "literal"
	expected := &literal
	actual := pointer.Ptr(literal)
	require.Equal(t, expected, actual)
}

func TestUnPtr(t *testing.T) {
	{
		expected := "value"
		actual := pointer.UnPtr(&expected)
		require.Equal(t, expected, actual)
	}
	{
		var expected string
		actual := pointer.UnPtr(&expected)
		require.Equal(t, expected, actual)
		actual = pointer.UnPtr[string](nil)
		require.Equal(t, expected, actual)
	}
}
