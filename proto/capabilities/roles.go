/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	// viewerCaps are read-only capabilities that do not affect state.
	ViewerCaps = sortCaps(append([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_LIST,

		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_LIST,
		Capability_CAP_IAM_GROUP_INVITES_LIST,
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_IAM_ROLE_BINDINGS_LIST,
		Capability_CAP_IAM_ROLES_LIST,
		Capability_CAP_IAM_POLICY_LIST,
		Capability_CAP_IAM_IDENTITY_LIST,
		Capability_CAP_IAM_IDENTITY_PROVIDERS_LIST,

		Capability_CAP_TENANT_CLUSTERS_DISCOVER,
		Capability_CAP_TENANT_CLUSTERS_LIST,
		Capability_CAP_TENANT_NAMESPACES_LIST,
		Capability_CAP_TENANT_NODES_LIST,
		Capability_CAP_TENANT_RECORDS_LIST,
		Capability_CAP_TENANT_RECORD_CONTEXTS_LIST,
		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_RECORD_POLICY_RESULTS_LIST,
		Capability_CAP_TENANT_RISKS_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,
		Capability_CAP_TENANT_WORKLOADS_LIST,

		Capability_CAP_SIGSTORE_LIST,
	},
		// Viewers can also list repos and tags, and pull images.
		RegistryPullCaps...))

	// editorCaps can modify state, but not grant roles/permissions.
	EditorCaps = sortCaps(append([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_CREATE,
		Capability_CAP_EVENTS_SUBSCRIPTION_DELETE,
		Capability_CAP_EVENTS_SUBSCRIPTION_UPDATE,

		Capability_CAP_TENANT_CLUSTERS_CREATE,
		Capability_CAP_TENANT_CLUSTERS_UPDATE,
		Capability_CAP_TENANT_CLUSTERS_DELETE,

		Capability_CAP_SIGSTORE_CERTIFICATE_CREATE,
		Capability_CAP_SIGSTORE_CREATE,
		Capability_CAP_SIGSTORE_DELETE,
		Capability_CAP_SIGSTORE_UPDATE,
	}, ViewerCaps...))

	// ownerCaps includes all capabilities possible by a user.
	OwnerCaps = sortCaps(append([]Capability{
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_CREATE,
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_DELETE,
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_UPDATE,

		Capability_CAP_IAM_GROUP_INVITES_CREATE,
		Capability_CAP_IAM_GROUP_INVITES_DELETE,

		Capability_CAP_IAM_GROUPS_CREATE,
		Capability_CAP_IAM_GROUPS_DELETE,
		Capability_CAP_IAM_GROUPS_UPDATE,

		Capability_CAP_IAM_POLICY_CREATE,
		Capability_CAP_IAM_POLICY_UPDATE,
		Capability_CAP_IAM_POLICY_DELETE,

		Capability_CAP_IAM_IDENTITY_CREATE,
		Capability_CAP_IAM_IDENTITY_DELETE,
		Capability_CAP_IAM_IDENTITY_UPDATE,

		Capability_CAP_IAM_IDENTITY_PROVIDERS_CREATE,
		Capability_CAP_IAM_IDENTITY_PROVIDERS_DELETE,
		Capability_CAP_IAM_IDENTITY_PROVIDERS_UPDATE,

		Capability_CAP_IAM_ROLE_BINDINGS_CREATE,
		Capability_CAP_IAM_ROLE_BINDINGS_DELETE,
		Capability_CAP_IAM_ROLE_BINDINGS_UPDATE,

		Capability_CAP_IAM_ROLES_CREATE,
		Capability_CAP_IAM_ROLES_DELETE,
		Capability_CAP_IAM_ROLES_UPDATE,

		// Add gulfstream capability to owner so owners can rolebind
		// identities to the gulfstream role.
		Capability_CAP_GULFSTREAM,
	}, append(EditorCaps,
		// Owners can also push and delete images, subject to the identity allowlist.
		RegistryPushCaps...)...))

	RegistryPullCaps = sortCaps([]Capability{
		Capability_CAP_IAM_GROUPS_LIST,

		Capability_CAP_REPO_LIST,
		Capability_CAP_MANIFEST_LIST,
		Capability_CAP_TAG_LIST,

		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,
	})

	RegistryPushCaps = sortCaps(append([]Capability{
		Capability_CAP_REPO_CREATE,
		Capability_CAP_REPO_UPDATE,
		Capability_CAP_REPO_DELETE,

		Capability_CAP_MANIFEST_CREATE,
		Capability_CAP_MANIFEST_UPDATE,
		Capability_CAP_MANIFEST_DELETE,

		Capability_CAP_TAG_CREATE,
		Capability_CAP_TAG_UPDATE,
		Capability_CAP_TAG_DELETE,
	}, RegistryPullCaps...))

	RegistryPullTokenCreatorCaps = sortCaps(append([]Capability{
		// Minimal set of capabilities to create a registry pull token.
		Capability_CAP_IAM_ROLE_BINDINGS_CREATE,
		Capability_CAP_IAM_IDENTITY_CREATE,

		Capability_CAP_IAM_ROLES_LIST,
	}, RegistryPullCaps...))

	SigningViewerCaps = sortCaps([]Capability{
		Capability_CAP_SIGSTORE_LIST,
	})

	SigningCertRequesterCaps = sortCaps(append([]Capability{
		Capability_CAP_SIGSTORE_CERTIFICATE_CREATE,
	}, SigningViewerCaps...))

	SigningEditorCaps = sortCaps(append([]Capability{
		Capability_CAP_SIGSTORE_CREATE,
		Capability_CAP_SIGSTORE_DELETE,
		Capability_CAP_SIGSTORE_UPDATE,
	}, SigningCertRequesterCaps...))
)

func sortCaps(caps []Capability) []Capability {
	uniq := map[Capability]struct{}{}
	for _, c := range caps {
		uniq[c] = struct{}{}
	}
	out := maps.Keys(uniq)
	slices.Sort(out)
	return out
}
