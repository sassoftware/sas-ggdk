// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

// Map applies the given Mapper to the Result. If the Result encapsulates a
// value then that value is passed to the Mapper and a Result of the return type
// is returned. If the Result encapsulates an error then the error is propagated
// to a Result of the return type.
func Map[T, R any](
	f Mapper[T, R],
	r Result[T],
) Result[R] {
	if r.IsError() {
		return Error[R](r.Error())
	}
	return New(f(r.MustGet()))
}

// Map2 applies the given Mapper2 to the given Results. If both Results
// encapsulate values then those values are passed to the Mapper2 and a Result
// of the return type is returned. If either Result encapsulates an error then
// the error of the first argument that has an error (from left to right) is
// propagated to a Result of the return type.
func Map2[T1, T2, R any](
	f Mapper2[T1, T2, R],
	r1 Result[T1],
	r2 Result[T2],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	return New(f(r1.MustGet(), r2.MustGet()))
}

// Map3 applies the given Mapper3 to the given Results. If all Results
// encapsulate values then those values are passed to the Mapper3 and a Result
// of the return type is returned. If either Result encapsulates an error then
// the error of the first argument that has an error (from left to right) is
// propagated to a Result of the return type.
func Map3[T1, T2, T3, R any](
	f Mapper3[T1, T2, T3, R],
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
	return New(f(r1.MustGet(), r2.MustGet(), r3.MustGet()))
}

// Map4 applies the given Mapper4 to the given Results. If all Results
// encapsulate values then those values are passed to the Mapper4 and a Result
// of the return type is returned. If either Result encapsulates an error then
// the error of the first argument that has an error (from left to right) is
// propagated to a Result of the return type.
func Map4[T1, T2, T3, T4, R any](
	f Mapper4[T1, T2, T3, T4, R],
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
	return New(f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet()))
}

// Map5 applies the given Mapper5 to the given Results. If all Results
// encapsulate values then those values are passed to the Mapper5 and a Result
// of the return type is returned. If either Result encapsulates an error then
// the error of the first argument that has an error (from left to right) is
// propagated to a Result of the return type.
func Map5[T1, T2, T3, T4, T5, R any](
	f Mapper5[T1, T2, T3, T4, T5, R],
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
	return New(f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet(), r5.MustGet()))
}

// MapNoError applies the given MapperNoError to the Result. If the Result
// encapsulates a value then that value is passed to the MapperNoError and a
// Result of the return type is returned. If the Result encapsulates an error
// then the error is propagated to a Result of the return type.
func MapNoError[T, R any](
	f MapperNoError[T, R],
	r Result[T],
) Result[R] {
	if r.IsError() {
		return Error[R](r.Error())
	}
	return Ok[R](f(r.MustGet()))
}

// MapNoError2 applies the given MapperNoError2 to the given Results. If both
// Results encapsulate values then those values are passed to the MapperNoError2
// and a Result of the return type is returned. If either Result encapsulates an
// error then the error of the first argument that has an error (from left to
// right) is propagated to a Result of the return type.
func MapNoError2[T1, T2, R any](
	f MapperNoError2[T1, T2, R],
	r1 Result[T1],
	r2 Result[T2],
) Result[R] {
	if r1.IsError() {
		return Error[R](r1.Error())
	}
	if r2.IsError() {
		return Error[R](r2.Error())
	}
	return Ok[R](f(r1.MustGet(), r2.MustGet()))
}

// MapNoError3 applies the given MapperNoError3 to the given Results. If all
// Results encapsulate values then those values are passed to the MapperNoError3
// and a Result of the return type is returned. If either Result encapsulates an
// error then the error of the first argument that has an error (from left to
// right) is propagated to a Result of the return type.
func MapNoError3[T1, T2, T3, R any](
	f MapperNoError3[T1, T2, T3, R],
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
	return Ok[R](f(r1.MustGet(), r2.MustGet(), r3.MustGet()))
}

// MapNoError4 applies the given MapperNoError4 to the given Results. If all
// Results encapsulate values then those values are passed to the MapperNoError4
// and a Result of the return type is returned. If either Result encapsulates an
// error then the error of the first argument that has an error (from left to
// right) is propagated to a Result of the return type.
func MapNoError4[T1, T2, T3, T4, R any](
	f MapperNoError4[T1, T2, T3, T4, R],
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
	return Ok[R](f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet()))
}

// MapNoError5 applies the given MapperNoError5 to the given Results. If all
// Results encapsulate values then those values are passed to the MapperNoError5
// and a Result of the return type is returned. If either Result encapsulates an
// error then the error of the first argument that has an error (from left to
// right) is propagated to a Result of the return type.
func MapNoError5[T1, T2, T3, T4, T5, R any](
	f MapperNoError5[T1, T2, T3, T4, T5, R],
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
	return Ok[R](f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet(), r5.MustGet()))
}

// MapErrorOnly applies the given MapperErrorOnly to the Result. If the Result
// encapsulates a value then that value is passed to the MapperErrorOnly and the
// error is returned. If the Result encapsulates an error then that error is
// returned.
func MapErrorOnly[T any](
	f MapperErrorOnly[T],
	r Result[T],
) error {
	if r.IsError() {
		return r.Error()
	}
	return f(r.MustGet())
}

// MapErrorOnly2 applies the given MapperErrorOnly2 to the Result. If the Result
// encapsulates a value then that value is passed to the MapperErrorOnly2 and
// the error is returned. If either Result encapsulates an error then the error
// of the first argument that has an error (from left to right) is returned.
func MapErrorOnly2[T1, T2 any](
	f MapperErrorOnly2[T1, T2],
	r1 Result[T1],
	r2 Result[T2],
) error {
	if r1.IsError() {
		return r1.Error()
	}
	if r2.IsError() {
		return r2.Error()
	}
	return f(r1.MustGet(), r2.MustGet())
}

// MapErrorOnly3 applies the given MapperErrorOnly3 to the Result. If the Result
// encapsulates a value then that value is passed to the MapperErrorOnly3 and
// the error is returned. If any Result encapsulates an error then the error of
// the first argument that has an error (from left to right) is returned.
func MapErrorOnly3[T1, T2, T3 any](
	f MapperErrorOnly3[T1, T2, T3],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
) error {
	if r1.IsError() {
		return r1.Error()
	}
	if r2.IsError() {
		return r2.Error()
	}
	if r3.IsError() {
		return r3.Error()
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet())
}

// MapErrorOnly4 applies the given MapperErrorOnly4 to the Result. If the Result
// encapsulates a value then that value is passed to the MapperErrorOnly4 and
// the error is returned. If any Result encapsulates an error then the error of
// the first argument that has an error (from left to right) is returned.
func MapErrorOnly4[T1, T2, T3, T4 any](
	f MapperErrorOnly4[T1, T2, T3, T4],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
	r4 Result[T4],
) error {
	if r1.IsError() {
		return r1.Error()
	}
	if r2.IsError() {
		return r2.Error()
	}
	if r3.IsError() {
		return r3.Error()
	}
	if r4.IsError() {
		return r4.Error()
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet())
}

// MapErrorOnly5 applies the given MapperErrorOnly5 to the Result. If the Result
// encapsulates a value then that value is passed to the MapperErrorOnly5 and
// the error is returned. If any Result encapsulates an error then the error of
// the first argument that has an error (from left to right) is returned.
func MapErrorOnly5[T1, T2, T3, T4, T5 any](
	f MapperErrorOnly5[T1, T2, T3, T4, T5],
	r1 Result[T1],
	r2 Result[T2],
	r3 Result[T3],
	r4 Result[T4],
	r5 Result[T5],
) error {
	if r1.IsError() {
		return r1.Error()
	}
	if r2.IsError() {
		return r2.Error()
	}
	if r3.IsError() {
		return r3.Error()
	}
	if r4.IsError() {
		return r4.Error()
	}
	if r5.IsError() {
		return r5.Error()
	}
	return f(r1.MustGet(), r2.MustGet(), r3.MustGet(), r4.MustGet(), r5.MustGet())
}
