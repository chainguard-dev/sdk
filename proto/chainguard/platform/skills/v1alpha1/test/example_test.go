/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
	skillstest "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1/test"
)

// ExampleMockSkillsClients demonstrates constructing a mock skills
// catalog client.
func ExampleMockSkillsClients() {
	mock := skillstest.MockSkillsClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}

// ExampleMockHardenClient polls a harden operation through the mock bundle.
func ExampleMockHardenClient() {
	const (
		group = "720909c9f5279097d847ad02a2f24ba8f59de36a"
		name  = "harden/" + group + "/0123abcd"
	)
	mock := skillstest.MockSkillsClients{
		HardenClient: skillstest.MockHardenClient{
			OnGetHardenOperation: []skillstest.HardenOnGetOperation{{
				Given:     &skills.GetHardenOperationRequest{Group: group, Name: name},
				Operation: &longrunningpb.Operation{Name: name, Done: true},
			}},
		},
	}
	op, err := mock.Harden().GetHardenOperation(context.Background(), &skills.GetHardenOperationRequest{Group: group, Name: name})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(op.GetDone())
	// Output:
	// true
}

// ExampleMockHardenClient_list lists a group's harden operations through the mock bundle.
func ExampleMockHardenClient_list() {
	const (
		group = "720909c9f5279097d847ad02a2f24ba8f59de36a"
		name  = "harden/" + group + "/0123abcd"
	)
	mock := skillstest.MockSkillsClients{
		HardenClient: skillstest.MockHardenClient{
			OnListHardenOperations: []skillstest.HardenOnListOperations{{
				Given: &skills.ListHardenOperationsRequest{Group: group},
				List: &skills.ListHardenOperationsResponse{
					Operations: []*longrunningpb.Operation{{Name: name}},
				},
			}},
		},
	}
	resp, err := mock.Harden().ListHardenOperations(context.Background(), &skills.ListHardenOperationsRequest{Group: group})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, op := range resp.GetOperations() {
		fmt.Println(op.GetName())
	}
	// Output:
	// harden/720909c9f5279097d847ad02a2f24ba8f59de36a/0123abcd
}
