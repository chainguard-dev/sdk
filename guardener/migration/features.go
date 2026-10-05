/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package migration

import (
	"slices"
	"strings"

	cappb "chainguard.dev/sdk/proto/capabilities"
	guardpb "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"
)

// Definition identifies a migration worker and its required capability.
type Definition struct {
	Feature     guardpb.MigrationFeature
	QueueSuffix string
	Capability  cappb.Capability
	Name        string
}

var features = []Definition{
	{Feature: guardpb.MigrationFeature_MIGRATION_FEATURE_ACTIONS, QueueSuffix: "/tree/HEAD/.github/workflows", Capability: cappb.Capability_CAP_GUARDENER_ACTIONS_MIGRATE, Name: "Actions"},
	{Feature: guardpb.MigrationFeature_MIGRATION_FEATURE_IMAGES, QueueSuffix: "/tree/HEAD/.chainguard/images.yaml", Capability: cappb.Capability_CAP_GUARDENER_IMAGES_MIGRATE, Name: "Images"},
}

// All returns the registered features in their display order.
func All() []Definition { return slices.Clone(features) }

// Lookup returns the definition for a registered feature.
func Lookup(feature guardpb.MigrationFeature) (Definition, bool) {
	for _, definition := range features {
		if definition.Feature == feature {
			return definition, true
		}
	}
	return Definition{}, false
}

// Key returns the repository workqueue key shared by enqueueing and reporting.
func (d Definition) Key(repoURL string) string {
	if _, rest, ok := strings.Cut(repoURL, "://"); ok {
		repoURL = rest
	}
	return strings.TrimRight(repoURL, "/") + d.QueueSuffix
}
