// SPDX-FileCopyrightText:  2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer

// Ptr returns a pointer to the given value.
func Ptr[T any](o T) *T {
	return &o
}

// UnPtr returns either the value pointed to by p or the zero value of T.
func UnPtr[T any](p *T) T {
	if p != nil {
		return *p
	}
	return *new(T)
}
