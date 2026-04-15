/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/proto/chainguard/platform/test"
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1"
	vulntest "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1/test"
)

func ExampleMockClients() {
	mock := &vulntest.MockClients{
		AdvisoriesServiceClient: vulntest.MockAdvisoriesServiceClient{
			OnListAdvisories: []test.On[*vuln.ListAdvisoriesRequest, *vuln.ListAdvisoriesResponse]{{
				Given: &vuln.ListAdvisoriesRequest{},
				Result: &vuln.ListAdvisoriesResponse{
					Advisories: []*vuln.Advisory{{AdvisoryId: "CGA-2026-1234"}},
				},
			}},
		},
	}

	resp, err := mock.AdvisoriesService().ListAdvisories(context.Background(), &vuln.ListAdvisoriesRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Advisories[0].AdvisoryId)
	// Output: CGA-2026-1234
}
