/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package civil_test

import (
	"fmt"
	"time"

	"chainguard.dev/sdk/civil"
)

// ExampleDateOf demonstrates converting a time.Time to a civil Date.
func ExampleDateOf() {
	t := time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)
	d := civil.DateOf(t)
	fmt.Println(d.Year)
	fmt.Println(int(d.Month))
	fmt.Println(d.Day)
	// Output:
	// 2024
	// 3
	// 15
}

// ExampleParseDate demonstrates parsing a date string.
func ExampleParseDate() {
	d, err := civil.ParseDate("2024-03-15")
	fmt.Println(err)
	fmt.Println(d.Year)
	// Output:
	// <nil>
	// 2024
}
