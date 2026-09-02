/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

// capturingVulnReportsClient records the page sizes of the requests it
// receives and returns one empty terminal page.
type capturingVulnReportsClient struct {
	VulnReportsServiceClient
	pageSizes []int32
}

func (c *capturingVulnReportsClient) ListVulnCountReports(_ context.Context, req *ListVulnCountReportsRequest, _ ...grpc.CallOption) (*ListVulnCountReportsResponse, error) {
	c.pageSizes = append(c.pageSizes, req.GetPageSize())
	return &ListVulnCountReportsResponse{}, nil
}

// Test_VulnCountReports_ListIter_DefaultsPageSizeToServerDefault pins the
// guard against silent truncation: the shared iterator stops after its page
// cap, so an unset page_size must become the server's tuned default rather
// than the iterator's generic 50.
func Test_VulnCountReports_ListIter_DefaultsPageSizeToServerDefault(t *testing.T) {
	tests := []struct {
		name      string
		requested int32
		want      int32
	}{{
		name:      "unset becomes the server default",
		requested: 0,
		want:      vulnCountReportsIterPageSize,
	}, {
		name:      "explicit page size is respected",
		requested: 7,
		want:      7,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture := &capturingVulnReportsClient{}
			c := &clients{vulnReportsService: capture}
			req := &ListVulnCountReportsRequest{PageSize: tt.requested}
			if _, err := c.ListVulnCountReportsAll(t.Context(), req); err != nil {
				t.Fatalf("ListVulnCountReportsAll() error: %v", err)
			}
			if len(capture.pageSizes) != 1 || capture.pageSizes[0] != tt.want {
				t.Errorf("page sizes sent: got = %v, want = [%d]", capture.pageSizes, tt.want)
			}
			if req.GetPageSize() != tt.requested {
				t.Errorf("caller's request mutated: got = %d, want = %d", req.GetPageSize(), tt.requested)
			}
		})
	}
}
