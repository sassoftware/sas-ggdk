// SPDX-FileCopyrightText: 2022, SAS Institute Inc., Cary, NC, USA. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

// FlatMapper defines a function with one argument that maps from the
// argument's type to a Result.
type FlatMapper[T, R any] func(
	T,
) Result[R]

// FlatMapper2 defines a function with one argument that maps from the
// argument's type to a Result.
type FlatMapper2[T1, T2, R any] func(
	T1, T2,
) Result[R]

// FlatMapper3 defines a function with one argument that maps from the
// argument's type to a Result.
type FlatMapper3[T1, T2, T3, R any] func(
	T1, T2, T3,
) Result[R]

// FlatMapper4 defines a function with one argument that maps from the
// argument's type to a Result.
type FlatMapper4[T1, T2, T3, T4, R any] func(
	T1, T2, T3, T4,
) Result[R]

// FlatMapper5 defines a function with one argument that maps from the
// argument's type to a Result.
type FlatMapper5[T1, T2, T3, T4, T5, R any] func(
	T1, T2, T3, T4, T5,
) Result[R]

// MakeFlatMapper takes a Mapper and returns a FlatMapper that can be used with
// FlatMap.
func MakeFlatMapper[T, R any](
	f Mapper[T, R],
) FlatMapper[T, R] {
	return func(t T) Result[R] {
		return New(f(t))
	}
}

// MakeFlatMapper2 takes a Mapper2 and returns a FlatMapper2 that can be used
// with FlatMap2.
func MakeFlatMapper2[T1, T2, R any](
	f Mapper2[T1, T2, R],
) FlatMapper2[T1, T2, R] {
	return func(t1 T1, t2 T2) Result[R] {
		return New(f(t1, t2))
	}
}

// MakeFlatMapper3 takes a Mapper3 and returns a FlatMapper3 that can be used
// with FlatMap3.
func MakeFlatMapper3[T1, T2, T3, R any](
	f Mapper3[T1, T2, T3, R],
) FlatMapper3[T1, T2, T3, R] {
	return func(t1 T1, t2 T2, t3 T3) Result[R] {
		return New(f(t1, t2, t3))
	}
}

// MakeFlatMapper4 takes a Mapper4 and returns a FlatMapper4 that can be used
// with FlatMap4.
func MakeFlatMapper4[T1, T2, T3, T4, R any](
	f Mapper4[T1, T2, T3, T4, R],
) FlatMapper4[T1, T2, T3, T4, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) Result[R] {
		return New(f(t1, t2, t3, t4))
	}
}

// MakeFlatMapper5 takes a Mapper5 and returns a FlatMapper5 that can be used
// with FlatMap5.
func MakeFlatMapper5[T1, T2, T3, T4, T5, R any](
	f Mapper5[T1, T2, T3, T4, T5, R],
) FlatMapper5[T1, T2, T3, T4, T5, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Result[R] {
		return New(f(t1, t2, t3, t4, t5))
	}
}

// MakeFlatMapperNoError takes a MapperNoError and returns a FlatMapper that can
// be used with FlatMap.
func MakeFlatMapperNoError[T, R any](
	f MapperNoError[T, R],
) FlatMapper[T, R] {
	return func(t T) Result[R] {
		return Ok(f(t))
	}
}

// MakeFlatMapperNoError2 takes a MapperNoError2 and returns a FlatMapper2 that
// can be used with FlatMap2.
func MakeFlatMapperNoError2[T1, T2, R any](
	f MapperNoError2[T1, T2, R],
) FlatMapper2[T1, T2, R] {
	return func(t1 T1, t2 T2) Result[R] {
		return Ok(f(t1, t2))
	}
}

// MakeFlatMapperNoError3 takes a MapperNoErrors3 and returns a FlatMapper3 that
// can be used with FlatMap3.
func MakeFlatMapperNoError3[T1, T2, T3, R any](
	f MapperNoError3[T1, T2, T3, R],
) FlatMapper3[T1, T2, T3, R] {
	return func(t1 T1, t2 T2, t3 T3) Result[R] {
		return Ok(f(t1, t2, t3))
	}
}

// MakeFlatMapperNoError4 takes a MapperNoError4 and returns a FlatMapper4 that
// can be used with FlatMap4.
func MakeFlatMapperNoError4[T1, T2, T3, T4, R any](
	f MapperNoError4[T1, T2, T3, T4, R],
) FlatMapper4[T1, T2, T3, T4, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) Result[R] {
		return Ok(f(t1, t2, t3, t4))
	}
}

// MakeFlatMapperNoError5 takes a MapperNoError5 and returns a FlatMapper5 that
// can be used with FlatMap5.
func MakeFlatMapperNoError5[T1, T2, T3, T4, T5, R any](
	f MapperNoError5[T1, T2, T3, T4, T5, R],
) FlatMapper5[T1, T2, T3, T4, T5, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Result[R] {
		return Ok(f(t1, t2, t3, t4, t5))
	}
}
