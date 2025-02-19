// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package lazy defines functions and types for implementing lazy initialization.

Without this package, each struct that wants to lazily create the value of a
field must manage a getter function, a creator function, and state to know when
the creator function has been called. Something like the following.

	type Thing struct {
		expensiveStruct *ExpensiveStruct
	}

	func New() *Thing {
		return &Thing{}
	}

	func (t *Thing) createExpensiveStruct() *ExpensiveStruct {
		return &ExpensiveStruct{
			//fields
		}
	}

	func (t *Thing) getExpensiveStruct() *ExpensiveStruct {
		if t.expensiveStruct != nil {
			t.expensiveStruct = t.createExpensiveStruct()
		}
		return t.expensiveStruct
	}

	func (t *Thing) ExportedFunction() ReturnType {
		expensiveStruct = t.getExpensiveStruct()
		return expensiveStruct.ExportedFunction()
	}

With this package, there is no need for each struct to re-implement this logic.
Also, because the creator holds the instance, it is impossible to accidentally
access the field directly.

	type Thing struct {
		getExpensiveStruct memoize.Creator[*ExpensiveStruct]
	}

	func New() *Thing {
		t := &Thing{}
		t.getExpensiveStruct = memoize.MakeGetter(t.createExpensiveStruct)
		return t
	}

	func (t *Thing) createExpensiveStruct() *ExpensiveStruct {
		return &ExpensiveStruct{
			//fields
		}
	}

	func (t *Thing) ExportedFunction() ReturnType {
		expensiveStruct = t.getExpensiveStruct()
		return expensiveStruct.ExportedFunction()
	}
*/
package lazy
