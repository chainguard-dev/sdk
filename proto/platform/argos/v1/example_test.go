/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"errors"
	"fmt"

	argosv1 "chainguard.dev/sdk/proto/platform/argos/v1"
)

// ExampleNewClientsFromConnection demonstrates constructing an Argos client
// from an existing gRPC connection.
func ExampleNewClientsFromConnection() {
	cl := argosv1.NewClientsFromConnection(nil)
	fmt.Println(cl != nil)
	// Output:
	// true
}

// ExampleValidateRejectionReason shows that a caller stores the trimmed
// reason, and that a reason carrying a newline is refused rather than
// rewritten.
func ExampleValidateRejectionReason() {
	reason, err := argosv1.ValidateRejectionReason("  Taint tracking misses deserialization sinks.\n")
	fmt.Printf("%q %v\n", reason, err)

	_, err = argosv1.ValidateRejectionReason("first line\nsecond line")
	fmt.Println(errors.Is(err, argosv1.ErrRejectionReasonMalformed))
	// Output:
	// "Taint tracking misses deserialization sinks." <nil>
	// true
}
