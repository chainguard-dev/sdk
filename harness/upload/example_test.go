/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"context"
	"errors"
	"fmt"

	"chainguard.dev/sdk/harness/upload"
	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
)

// ExampleClient_Upload uploads a packed bundle on disk to an organization.
// uploadsClient is a harness Uploads API client, and org is the
// organization's UIDP.
func ExampleClient_Upload() {
	var (
		ctx           = context.Background()
		uploadsClient harnesspb.UploadsClient
		org           string
	)

	b, err := upload.DescribeFile("hunt.bundle")
	if err != nil {
		fmt.Println(err)
		return
	}
	u, err := (&upload.Client{Uploads: uploadsClient}).Upload(ctx, org, b)
	if rejected, ok := errors.AsType[*upload.RejectedError](err); ok {
		// The server holds the upload but refused it, and says why.
		fmt.Println(rejected.Upload.GetRejectionReason())
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u.GetId())
}
