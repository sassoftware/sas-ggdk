// SPDX-FileCopyrightText: 2022, SAS Institute Inc., Cary, NC, USA. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package result

// Mapper defines a function with one argument of any type that returns a single
// value and an error.
type Mapper[T, R any] func(
	T,
) (R, error)

// Mapper2 defines a function with two arguments of any type that returns a
// single value and an error.
type Mapper2[T1, T2, R any] func(
	T1, T2,
) (R, error)

// Mapper3 defines a function with three arguments of any type that returns a
// single value and an error.
type Mapper3[T1, T2, T3, R any] func(
	T1, T2, T3,
) (R, error)

// Mapper4 defines a function with four arguments of any type that returns a
// single value and an error.
type Mapper4[T1, T2, T3, T4, R any] func(
	T1, T2, T3, T4,
) (R, error)

// Mapper5 defines a function with five arguments of any type that returns a
// single value and an error.
type Mapper5[T1, T2, T3, T4, T5, R any] func(
	T1, T2, T3, T4, T5,
) (R, error)

// MapperNoError defines a function with one argument of any type that returns a
// single value.
type MapperNoError[T, R any] func(
	T,
) R

// MapperNoError2 defines a function with two arguments of any type that returns
// a single value.
type MapperNoError2[T1, T2, R any] func(
	T1, T2,
) R

// MapperNoError3 defines a function with three arguments of any type that returns
// a single value.
type MapperNoError3[T1, T2, T3, R any] func(
	T1, T2, T3,
) R

// MapperNoError4 defines a function with four arguments of any type that returns
// a single value.
type MapperNoError4[T1, T2, T3, T4, R any] func(
	T1, T2, T3, T4,
) R

// MapperNoError5 defines a function with five arguments of any type that returns
// a single value.
type MapperNoError5[T1, T2, T3, T4, T5, R any] func(
	T1, T2, T3, T4, T5,
) R

// MapperErrorOnly defines a function with one argument of any type that returns
// an error.
type MapperErrorOnly[T any] func(
	T,
) error

// MapperErrorOnly2 defines a function with two arguments of any type that
// returns an error.
type MapperErrorOnly2[T1, T2 any] func(
	T1, T2,
) error

// MapperErrorOnly3 defines a function with three arguments of any type that
// returns an error.
type MapperErrorOnly3[T1, T2, T3 any] func(
	T1, T2, T3,
) error

// MapperErrorOnly4 defines a function with four arguments of any type that
// returns an error.
type MapperErrorOnly4[T1, T2, T3, T4 any] func(
	T1, T2, T3, T4,
) error

// MapperErrorOnly5 defines a function with five arguments of any type that
// returns an error.
type MapperErrorOnly5[T1, T2, T3, T4, T5 any] func(
	T1, T2, T3, T4, T5,
) error
