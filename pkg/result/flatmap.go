// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

// FlatMap applies the given FlatMapper to the Result. If the Result
// encapsulates a value then that value is passed to the FlatMapper and the
// FlatMapper's return value is returned. If the Result encapsulates an error
// then the error is propagated to a Result of the return type.
func FlatMap[T, R any](
	f FlatMapper[T, R],
	r Result[T],
) Result[R] {
	if r.IsError() {
		return Error[R](r.Error())
	}
	return f(r.MustGet())
}

// FlatMap2 applies the given FlatMapper2 to the given Results. If both Results
// encapsulate a value then those values are passed to the FlatMapper2 and the
// FlatMapper2's return value is returned. If either Result encapsulates an
// error then the error of the first argument that has an error (from left to
// right) is propagated to a Result of the return type.
func FlatMap2[T1, T2, R any](
	f FlatMapper2[T1, T2, R],
	r1 Result[T1],
	r2 Result[T2],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	return f(r1.MustGet(), r2.MustGet())
}

// FlatMap3 applies the given FlatMapper3 to the given Results. If all Results
// encapsulate a value then those values are passed to the FlatMapper3 and the
// FlatMapper3's return value is returned. If any Result encapsulates an error
// then the error of the first argument that has an error (from left to right)
// is propagated to a Result of the return type.
func FlatMap3[T1, T2, T3, R any](
	f FlatMapper3[T1, T2, T3, R],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	if r3.IsError() {
		return Error[R](r3.Error())
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet())
}

// FlatMap4 applies the given FlatMapper4 to the given Results. If all Results
// encapsulate a value then those values are passed to the FlatMapper4 and the
// FlatMapper4's return value is returned. If any Result encapsulates an error
// then the error of the first argument that has an error (from left to right)
// is propagated to a Result of the return type.
func FlatMap4[T1, T2, T3, T4, R any](
	f FlatMapper4[T1, T2, T3, T4, R],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
	r4 Result[T4],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	if r3.IsError() {
		return Error[R](r3.Error())
	}
	if r4.IsError() {
		return Error[R](r4.Error())
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet())
}

// FlatMap5 applies the given FlatMapper5 to the given Results. If all Results
// encapsulate a value then those values are passed to the FlatMapper5 and the
// FlatMapper5's return value is returned. If any Result encapsulates an error
// then the error of the first argument that has an error (from left to right)
// is propagated to a Result of the return type.
func FlatMap5[T1, T2, T3, T4, T5, R any](
	f FlatMapper5[T1, T2, T3, T4, T5, R],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
	r4 Result[T4],
	r5 Result[T5],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	if r3.IsError() {
		return Error[R](r3.Error())
	}
	if r4.IsError() {
		return Error[R](r4.Error())
	}
	if r5.IsError() {
		return Error[R](r5.Error())
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet(), r5.MustGet())
}
