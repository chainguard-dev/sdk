/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	authv1 "chainguard.dev/sdk/proto/platform/auth/v1"
)

// ExampleSession_CloudEventsSubject demonstrates the CloudEventsSubject method.
func ExampleSession_CloudEventsSubject() {
	s := &authv1.Session{}
	subject := s.CloudEventsSubject()
	fmt.Println(subject == "")
	// Output:
	// true
}
