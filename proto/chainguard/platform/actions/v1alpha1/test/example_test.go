/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	actions "chainguard.dev/sdk/proto/chainguard/platform/actions/v1alpha1"
	actionstest "chainguard.dev/sdk/proto/chainguard/platform/actions/v1alpha1/test"
)

func ExampleMockActionsClient() {
	mock := &actionstest.MockActionsClient{
		OnListActions: []actionstest.ActionsOnListActions{{
			Given: &actions.ListActionsRequest{},
			Got:   &actions.ListActionsResponse{},
		}},
	}

	resp, err := mock.ListActions(context.Background(), &actions.ListActionsRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp != nil)
	// Output: true
}
