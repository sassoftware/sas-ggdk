// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package condition

import (
	"reflect"
)

// IsNil returns true if its parameter has a nil value, otherwise false.
func IsNil(i any) bool {
	if i == nil {
		return true
	}
	value := reflect.ValueOf(i)
	kind := value.Kind()
	switch kind {
	case
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return value.IsNil()
	}
	return false
}
