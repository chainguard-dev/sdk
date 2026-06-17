/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms_test

import (
	"fmt"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
)

func ExampleDocumentMetadata() {
	doc := terms.DocumentMetadata("guardener-tos.v1")
	fmt.Println(doc.Label)
	fmt.Println(doc.URL)
	// Output:
	// Terms of Service
	// https://www.chainguard.dev/legal/guardener
}
