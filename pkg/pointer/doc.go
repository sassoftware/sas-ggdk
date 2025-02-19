// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package pointer provides a mechanism for creating a pointer from expressions
whose address cannot be directly taken. For example, literal strings. &"literal"
is not valid but pointer.Ptr("literal") will return the address of a string
variable that contains "literal".
*/
package pointer
