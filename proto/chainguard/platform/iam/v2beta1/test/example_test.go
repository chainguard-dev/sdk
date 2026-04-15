/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	iamtest "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1/test"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

func ExampleMockClients() {
	mock := &iamtest.MockClients{
		GroupsServiceClient: iamtest.MockGroupsServiceClient{
			OnGetGroup: []test.On[*iam.GetGroupRequest, *iam.Group]{{
				Given:  &iam.GetGroupRequest{Uid: "abc123"},
				Result: &iam.Group{Uid: "abc123", Name: "my-org"},
			}},
		},
	}

	group, err := mock.GroupsService().GetGroup(context.Background(), &iam.GetGroupRequest{Uid: "abc123"})
	if err != nil {
		panic(err)
	}
	fmt.Println(group.Name)
	// Output: my-org
}
