/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
)

var _ tenant.AttestationsClient = (*MockAttestationsClientt)(nil)

type MockAttestationsClientt struct {
	OnList []AttestationsOnList
}

type AttestationsOnList struct {
	Given *tenant.AttestationFilter
	List  *tenant.AttestationList
	Error error
}

func (m MockAttestationsClientt) List(_ context.Context, given *tenant.AttestationFilter, _ ...grpc.CallOption) (*tenant.AttestationList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
