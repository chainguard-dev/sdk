/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	cgtest "chainguard.dev/sdk/proto/chainguard/platform/test"
)

// ExampleMockIter demonstrates converting a slice into an iterator for mock clients.
func ExampleMockIter() {
	items := []string{"a", "b", "c"}
	seq := cgtest.MockIter(items, nil)
	var got []string
	for item, err := range seq {
		if err != nil {
			break
		}
		got = append(got, item)
	}
	fmt.Println(got)
	// Output:
	// [a b c]
}
