/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	v1 "chainguard.dev/sdk/helm/v1"
)

// ExamplePredicateType demonstrates the chart-lock attestation predicate type.
func ExamplePredicateType() {
	fmt.Println(v1.PredicateType)
	// Output:
	// https://chainguard.dev/attestation/chart-lock/v1
}

// ExampleLock demonstrates constructing a Lock value.
func ExampleLock() {
	lock := &v1.Lock{
		Chart: &v1.Chart{
			Package: "nginx",
			Ref:     "cgr.dev/chainguard/nginx-chart:1.0.0",
		},
	}
	fmt.Println(lock.Chart.Package)
	// Output:
	// nginx
}
