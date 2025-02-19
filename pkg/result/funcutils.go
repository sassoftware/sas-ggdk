// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

// CallErrorOnlyFunc is a helper function that will call an error returning
// function that is encapsulated in a Result. This is useful when deferring a
// call to a function that was returned in a Result. If the Result contains an
// error, that error is returned. Otherwise, the value returned by the
// encapsulated function is returned.
func CallErrorOnlyFunc(functionResult Result[func() error]) error {
	if functionResult.IsError() {
		return functionResult.Error()
	}
	fn := functionResult.MustGet()
	return fn()
}

// CallFlatFunc is a helper function that will call a Result returning function
// that is encapsulated in a Result. This is useful when deferring a call to a
// function that was returned in a Result.
func CallFlatFunc[T any](functionResult Result[func() Result[T]]) Result[T] {
	if functionResult.IsError() {
		return Error[T](functionResult.Error())
	}
	fn := functionResult.MustGet()
	return fn()
}

// CallFunc is a helper function that will call a function with no return value
// that is encapsulated in a Result. This is useful when deferring a call to a
// function that was returned in a Result.
func CallFunc(functionResult Result[func()]) {
	if functionResult.IsError() {
		return
	}
	fn := functionResult.MustGet()
	fn()
}
