// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package folders

import "github.com/sassoftware/sas-ggdk/pkg/result"

// Folder is a function that can transform a source value and store it in a
// target value. The various Fold functions will call a Folder repeatedly feeding
// the value in the returned result as the accumulator each call.
type Folder[T, S any] func(accumulator T, source S) result.Result[T]

// FolderNoError is a function that can transform a source value and store it in
// a target value. The various Fold functions will call a FolderNoError
// repeatedly feeding the returned value as the accumulator each call.
type FolderNoError[T, S any] func(accumulator T, source S) T

// ToKey is a function that can convert a source value to a key.
type ToKey[S any, K comparable] func(source S) result.Result[K]

// ToValue is a function that can convert a source value to another value.
type ToValue[S, V any] func(value S) result.Result[V]
