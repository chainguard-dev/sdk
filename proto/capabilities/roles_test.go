/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortCaps(t *testing.T) {
	tests := []struct {
		name      string
		capsLists [][]Capability
		want      []Capability
	}{{
		name:      "empty",
		capsLists: [][]Capability{},
		want:      []Capability{},
	}, {
		name: "single slice",
		capsLists: [][]Capability{
			{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_REPO_LIST, Capability_CAP_IAM_GROUPS_LIST},
		},
		want: []Capability{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_REPO_LIST},
	}, {
		name: "multiple slices",
		capsLists: [][]Capability{
			{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_REPO_LIST, Capability_CAP_IAM_GROUPS_LIST},
			{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_APK_LIST, Capability_CAP_IAM_ROLES_LIST, Capability_CAP_REPO_LIST},
		},
		want: []Capability{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_IAM_ROLES_LIST, Capability_CAP_REPO_LIST, Capability_CAP_APK_LIST},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SortCaps(test.capsLists...)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("(-want, +got) = %s", diff)
			}
		})
	}
}
