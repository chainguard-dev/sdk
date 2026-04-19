/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package civil implements types for civil time: a time-zone-independent
// representation of time following the proleptic Gregorian calendar with
// exactly 24-hour days, 60-minute hours, and 60-second minutes.
//
// Because they lack location information, these types do not represent
// unique moments or intervals of time. Use time.Time for that purpose.
package civil
