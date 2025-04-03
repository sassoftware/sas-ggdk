// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stringutils

import (
	"fmt"
	"strconv"

	"github.com/sassoftware/sas-ggdk/pkg/condition"
	"github.com/sassoftware/sas-ggdk/pkg/result"
	"github.com/sassoftware/sas-ggdk/pkg/sliceutils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DefaultToString returns the %v formatting of the given value.
func DefaultToString[T any](value T) result.Result[string] {
	return result.Ok(fmt.Sprintf(`%v`, value))
}

// AnyToString returns nil if the given value is nil or not a string. Otherwise,
// casts the value to string and returns a pointer to that string.
func AnyToString(value any) *string {
	if condition.IsNil(value) {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return nil
	}
	return &s
}

// ToQuoted returns a slice of strings that is the result of adding double quotes
// to the %v formatting of each element of the slice.
func ToQuoted[T any](values ...T) []string {
	quoter := result.MakeFlatMapperNoError(func(value T) string {
		return fmt.Sprintf(`"%v"`, value)
	})
	quotedValues := ToStringsWith(values, quoter)
	// Because the quoter never returns an error, we don't need to check before
	// calling MustGet().
	return quotedValues.MustGet()
}

// ToStrings ...
func ToStrings[T any](values ...T) result.Result[[]string] {
	return ToStringsWith(values, nil)
}

// ToStringsWith ...
func ToStringsWith[T any](values []T, toString result.FlatMapper[T, string]) result.Result[[]string] {
	if toString == nil {
		toString = DefaultToString[T]
	}
	return sliceutils.Map(toString, values)
}

// ToTitle returns the given string in English title case.
func ToTitle(value string) string {
	caser := cases.Title(language.English)
	return caser.String(value)
}

// AsBool returns the boolean value represented by the string as per
// https://golang.org/pkg/strconv/#ParseBool. If the value cannot be parsed,
// then the default answer is returned.
func AsBool(value string, defaultValue bool) bool {
	m := result.MakeFlatMapper(strconv.ParseBool)
	return m(value).OrElse(defaultValue)
}

// AsInt returns the boolean value represented by the string as per
// https://golang.org/pkg/strconv/#Atoi. If the value cannot be parsed, then
// the default answer is returned.
func AsInt(value string, defaultValue int) int {
	m := result.MakeFlatMapper(strconv.Atoi)
	return m(value).OrElse(defaultValue)
}

// AsInt64 returns the given string in the given base and bitSize as an int64.
// If the conversion fails, the given defaultValue is returned.
func AsInt64(value string, base, bitSize int, defaultValue int64) int64 {
	m := result.MakeFlatMapper3(strconv.ParseInt)
	return m(value, base, bitSize).OrElse(defaultValue)
}

// AsFloat returns the boolean value represented by the string as per
// https://golang.org/pkg/strconv/#ParseFloat. If the value cannot be parsed,
// then the default answer is returned.
func AsFloat(value string, bitSize int, defaultValue float64) float64 {
	m := result.MakeFlatMapper2(strconv.ParseFloat)
	return m(value, bitSize).OrElse(defaultValue)
}

// AsComplex returns the boolean value represented by the string as per
// https://golang.org/pkg/strconv/#ParseComplex. If the value cannot be parsed,
// then the default answer is returned.
func AsComplex(value string, bitSize int, defaultValue complex128) complex128 {
	m := result.MakeFlatMapper2(strconv.ParseComplex)
	return m(value, bitSize).OrElse(defaultValue)
}

// AsUint returns the boolean value represented by the string as per
// https://golang.org/pkg/strconv/#ParseUint. If the value cannot be parsed,
// then the default answer is returned.
func AsUint(value string, base, bitSize int, defaultValue uint64) uint64 {
	m := result.MakeFlatMapper3(strconv.ParseUint)
	return m(value, base, bitSize).OrElse(defaultValue)
}
