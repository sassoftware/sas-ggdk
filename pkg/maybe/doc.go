// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package maybe defines interfaces and structs for encapsulating a value (a Just)
or an absent value (a Nothing).

Without Maybe, code must rely on a convention like nil or zero value to indicate
that a value is not supplied. With a Maybe, the intention becomes more clear.

Two constructors are provided for creating a Maybe; Just which creates a Maybe
encapsulating a value, and Nothing which creates a Maybe encapsulating the
absence of a value. See those functions for more information.
*/
package maybe
