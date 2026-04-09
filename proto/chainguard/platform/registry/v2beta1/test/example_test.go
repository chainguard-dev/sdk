/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1"
	registrytest "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1/test"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

func ExampleMockClients() {
	mock := &registrytest.MockClients{
		ReposServiceClient: registrytest.MockReposServiceClient{
			OnGetRepo: []test.On[*registry.GetRepoRequest, *registry.Repo]{{
				Given:  &registry.GetRepoRequest{Uid: "abc123"},
				Result: &registry.Repo{Uid: "abc123", Name: "my-repo"},
			}},
		},
	}

	repo, err := mock.ReposService().GetRepo(context.Background(), &registry.GetRepoRequest{Uid: "abc123"})
	if err != nil {
		panic(err)
	}
	fmt.Println(repo.Name)
	// Output: my-repo
}
