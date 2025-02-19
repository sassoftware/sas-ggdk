// SPDX-FileCopyrightText:  2021, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package timeutils

import (
	"time"
)

const (
	militaryTimeLayout = `15:04:05` // Don't change this.
	timeHours          = `hours`
	timeMinutes        = `minutes`
	timeSeconds        = `seconds`
)

// MilitaryTime models the hours, minutes, and seconds in military time (aka "24
// hour time").
type MilitaryTime interface {
	// Hours returns the hours since midnight.
	Hours() string
	// Minutes returns the minutes since midnight.
	Minutes() string
	// Seconds returns the seconds since midnight.
	Seconds() string
}

type militaryTimeMap map[string]string

// NewMilitaryTime returns a new military time.
func NewMilitaryTime() MilitaryTime {
	return NewMilitaryTimeFrom(time.Now())
}

// NewMilitaryTimeFrom returns a new military time based on the given time. See
// https://pkg.go.dev/time#pkg-constants
func NewMilitaryTimeFrom(t time.Time) MilitaryTime {
	formattedTime := t.Format(militaryTimeLayout)
	return militaryTimeMap{
		timeHours:   formattedTime[0:2],
		timeMinutes: formattedTime[3:5],
		timeSeconds: formattedTime[6:8],
	}
}

// Hours returns the hours since midnight.
func (mt militaryTimeMap) Hours() string {
	return mt[timeHours]
}

// Minutes returns the minutes since midnight.
func (mt militaryTimeMap) Minutes() string {
	return mt[timeMinutes]
}

// Seconds returns the seconds since midnight.
func (mt militaryTimeMap) Seconds() string {
	return mt[timeSeconds]
}

// GetTimestamp returns the current time formatted per RFC822. The format is:
// `dd mmm yy hh:mm zzz`, for example `23 Feb 22 09:56 EST`.
// See https://pkg.go.dev/time#pkg-constants
func GetTimestamp() string {
	now := time.Now()
	return now.Format(time.RFC822)
}
