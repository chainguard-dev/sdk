/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"context"
	"fmt"

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
	skillstest "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1/test"
	common "chainguard.dev/sdk/proto/platform/common/v1"
)

// This example lists source filter values across two pages.
// The mock supplies responses so the example runs without a live API.
func ExampleSkillsClient_ListSkillSources() {
	const group = "720909c9f5279097d847ad02a2f24ba8f59de36a"
	scope := &common.UIDPFilter{ChildrenOf: group}
	client := &skillstest.MockCatalogClient{
		OnListSkillSources: []skillstest.SkillsOnListSkillSources{{
			Given: &skills.ListSkillSourcesRequest{Uidp: scope, PageSize: 1, OrderBy: "source asc"},
			List: &skills.ListSkillSourcesResponse{
				Items:         []*skills.Source{{Source: "anthropics"}},
				NextPageToken: "opaque-source-page-token",
			},
		}, {
			Given: &skills.ListSkillSourcesRequest{Uidp: scope, PageSize: 1, OrderBy: "source asc", PageToken: "opaque-source-page-token"},
			List: &skills.ListSkillSourcesResponse{
				Items: []*skills.Source{{Source: "openai"}},
			},
		}},
	}

	// Name the org explicitly. The API also applies the caller's authorized scope.
	req := &skills.ListSkillSourcesRequest{Uidp: scope, PageSize: 1, OrderBy: "source asc"}
	for {
		page, err := client.ListSkillSources(context.Background(), req)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, source := range page.GetItems() {
			// Pass the slug to ListSkillsRequest.Source.
			fmt.Println(source.GetSource())
		}
		if page.GetNextPageToken() == "" {
			break
		}
		// Treat tokens as opaque and preserve the scope and ordering.
		req.PageToken = page.GetNextPageToken()
	}
	// Output:
	// anthropics
	// openai
}
