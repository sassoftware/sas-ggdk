// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package maybe

// CallFunc is a helper function that will call a function with no return value
// that is encapsulated in a Maybe. This is useful when deferring a call to a
// function that was returned in a Maybe.
func CallFunc(functionResult Maybe[func()]) {
	if !functionResult.IsJust() {
		return
	}
	fn := functionResult.MustGet()
	fn()
}
