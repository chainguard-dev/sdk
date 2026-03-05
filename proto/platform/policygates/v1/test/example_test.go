/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	policygates "chainguard.dev/sdk/proto/platform/policygates/v1"
	"chainguard.dev/sdk/proto/platform/policygates/v1/test"
)

// ExampleMockPolicyGatesClients demonstrates how to use the mock clients
// for testing code that depends on the PolicyGates API.
func ExampleMockPolicyGatesClients() {
	ctx := context.Background()

	// Create a mock with configured responses
	mock := &test.MockPolicyGatesClients{
		PoliciesOnClient: test.MockPoliciesClient{
			OnListPolicies: []test.OnListPolicies{{
				Given: &policygates.PolicyFilter{},
				List: &policygates.PolicyList{
					Items: []*policygates.Policy{{
						Id:   "policy-1",
						Name: "example-policy",
					}},
				},
			}},
		},
	}

	// Use the mock in your code
	policies := mock.Policies()
	list, err := policies.ListPolicies(ctx, &policygates.PolicyFilter{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Found %d policies\n", len(list.Items))
	// Output: Found 1 policies
}

// ExampleMockPoliciesClient_CreatePolicy demonstrates mocking policy creation.
func ExampleMockPoliciesClient_CreatePolicy() {
	ctx := context.Background()

	mock := &test.MockPoliciesClient{
		OnCreatePolicy: []test.OnCreatePolicy{{
			Given: &policygates.CreatePolicyRequest{
				Policy: &policygates.Policy{
					Name: "new-policy",
				},
			},
			Created: &policygates.Policy{
				Id:   "policy-123",
				Name: "new-policy",
			},
		}},
	}

	created, err := mock.CreatePolicy(ctx, &policygates.CreatePolicyRequest{
		Policy: &policygates.Policy{
			Name: "new-policy",
		},
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Created policy: %s\n", created.Id)
	// Output: Created policy: policy-123
}

// ExampleMockPoliciesClient_UpdatePolicy demonstrates mocking policy updates.
func ExampleMockPoliciesClient_UpdatePolicy() {
	ctx := context.Background()

	mock := &test.MockPoliciesClient{
		OnUpdatePolicy: []test.OnUpdatePolicy{{
			Given: &policygates.Policy{
				Id:   "policy-1",
				Name: "updated-name",
			},
			Updated: &policygates.Policy{
				Id:   "policy-1",
				Name: "updated-name",
			},
		}},
	}

	updated, err := mock.UpdatePolicy(ctx, &policygates.Policy{
		Id:   "policy-1",
		Name: "updated-name",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Updated policy: %s\n", updated.Name)
	// Output: Updated policy: updated-name
}

// ExampleMockPoliciesClient_ListPolicies demonstrates mocking policy listing.
func ExampleMockPoliciesClient_ListPolicies() {
	ctx := context.Background()

	mock := &test.MockPoliciesClient{
		OnListPolicies: []test.OnListPolicies{{
			Given: &policygates.PolicyFilter{},
			List: &policygates.PolicyList{
				Items: []*policygates.Policy{{
					Id:   "policy-1",
					Name: "first-policy",
				}, {
					Id:   "policy-2",
					Name: "second-policy",
				}},
			},
		}},
	}

	list, err := mock.ListPolicies(ctx, &policygates.PolicyFilter{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Found %d policies\n", len(list.Items))
	// Output: Found 2 policies
}

// ExampleMockPoliciesClient_DeletePolicy demonstrates mocking policy deletion.
func ExampleMockPoliciesClient_DeletePolicy() {
	ctx := context.Background()

	mock := &test.MockPoliciesClient{
		OnDeletePolicy: []test.OnDeletePolicy{{
			Given: &policygates.DeletePolicyRequest{
				Id: "policy-1",
			},
		}},
	}

	_, err := mock.DeletePolicy(ctx, &policygates.DeletePolicyRequest{
		Id: "policy-1",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println("Policy deleted successfully")
	// Output: Policy deleted successfully
}

// ExampleMockBindingsClient_CreateBinding demonstrates mocking binding creation.
func ExampleMockBindingsClient_CreateBinding() {
	ctx := context.Background()

	mock := &test.MockBindingsClient{
		OnCreateBinding: []test.OnCreateBinding{{
			Given: &policygates.CreateBindingRequest{
				Binding: &policygates.Binding{
					Policy: "policy-1",
				},
			},
			Created: &policygates.Binding{
				Id:     "binding-123",
				Policy: "policy-1",
			},
		}},
	}

	created, err := mock.CreateBinding(ctx, &policygates.CreateBindingRequest{
		Binding: &policygates.Binding{
			Policy: "policy-1",
		},
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Created binding: %s\n", created.Id)
	// Output: Created binding: binding-123
}

// ExampleMockBindingsClient_UpdateBinding demonstrates mocking binding updates.
func ExampleMockBindingsClient_UpdateBinding() {
	ctx := context.Background()

	mock := &test.MockBindingsClient{
		OnUpdateBinding: []test.OnUpdateBinding{{
			Given: &policygates.Binding{
				Id:     "binding-1",
				Policy: "policy-2",
			},
			Updated: &policygates.Binding{
				Id:     "binding-1",
				Policy: "policy-2",
			},
		}},
	}

	updated, err := mock.UpdateBinding(ctx, &policygates.Binding{
		Id:     "binding-1",
		Policy: "policy-2",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Updated binding policy: %s\n", updated.Policy)
	// Output: Updated binding policy: policy-2
}

// ExampleMockBindingsClient_ListBindings demonstrates mocking binding listing.
func ExampleMockBindingsClient_ListBindings() {
	ctx := context.Background()

	mock := &test.MockBindingsClient{
		OnListBindings: []test.OnListBindings{{
			Given: &policygates.BindingFilter{},
			List: &policygates.BindingList{
				Items: []*policygates.Binding{{
					Id:     "binding-1",
					Policy: "policy-1",
				}, {
					Id:     "binding-2",
					Policy: "policy-2",
				}},
			},
		}},
	}

	list, err := mock.ListBindings(ctx, &policygates.BindingFilter{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Found %d bindings\n", len(list.Items))
	// Output: Found 2 bindings
}

// ExampleMockBindingsClient_DeleteBinding demonstrates mocking binding deletion.
func ExampleMockBindingsClient_DeleteBinding() {
	ctx := context.Background()

	mock := &test.MockBindingsClient{
		OnDeleteBinding: []test.OnDeleteBinding{{
			Given: &policygates.DeleteBindingRequest{
				Id: "binding-1",
			},
		}},
	}

	_, err := mock.DeleteBinding(ctx, &policygates.DeleteBindingRequest{
		Id: "binding-1",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println("Binding deleted successfully")
	// Output: Binding deleted successfully
}
