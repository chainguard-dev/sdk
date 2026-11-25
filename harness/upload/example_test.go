/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"context"
	"errors"
	"fmt"
	"os"

	"chainguard.dev/sdk/harness/progress"
	"chainguard.dev/sdk/harness/upload"
	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
	harnessv1alpha2 "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha2"
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

// ExamplePartClient_Upload uploads a packed bundle on disk to an
// organization in parts, printing progress to stderr. uploadsClient is a
// harness v1alpha2 Uploads API client, and org is the organization's UIDP.
func ExamplePartClient_Upload() {
	var (
		ctx           = context.Background()
		uploadsClient harnessv1alpha2.UploadsClient
		org           string
	)

	b, err := upload.DescribeFile("hunt.bundle")
	if err != nil {
		fmt.Println(err)
		return
	}
	pr := progress.New(os.Stderr)
	u, err := (&upload.PartClient{
		Uploads:  uploadsClient,
		Progress: func(sent, total int64) { pr.Report(progress.Update{Phase: "uploading", Bytes: sent, Total: total}) },
	}).Upload(ctx, org, b)
	pr.Done()
	if rejected, ok := errors.AsType[*upload.PartRejectedError](err); ok {
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
