/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestSkillsPublishCaps pins the bundle's membership so a regression that
// drops CAP_TERMS_LIST (which chainctl pre-flight needs to read acceptance
// state) won't silently leave skills publishers unable to push.
func TestSkillsPublishCaps(t *testing.T) {
	required := map[Capability]struct{}{
		Capability_CAP_SKILLS_PUBLISH:           {},
		Capability_CAP_SKILLS_ENTITLEMENTS_LIST: {},
		Capability_CAP_TERMS_LIST:               {},
	}
	got := make(map[Capability]struct{}, len(SkillsPublishCaps))
	for _, c := range SkillsPublishCaps {
		got[c] = struct{}{}
	}
	for c := range required {
		if _, ok := got[c]; !ok {
			t.Errorf("SkillsPublishCaps missing %v", c)
		}
	}
}

// TestGuardenerAdminCaps pins the admin bundle's membership so a regression
// can't silently drop the capability that gates self-service GitHub org
// linking, and confirms admins retain the full user capability set.
func TestGuardenerAdminCaps(t *testing.T) {
	required := map[Capability]struct{}{
		Capability_CAP_GUARDENER_ASSOCIATION_MANAGE: {},
		Capability_CAP_TERMS_ACCEPT:                 {},
		Capability_CAP_GUARDENER_DFC_CONVERT:        {},
	}
	got := make(map[Capability]struct{}, len(GuardenerAdminCaps))
	for _, c := range GuardenerAdminCaps {
		got[c] = struct{}{}
	}
	for c := range required {
		if _, ok := got[c]; !ok {
			t.Errorf("GuardenerAdminCaps missing %v", c)
		}
	}
}

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
