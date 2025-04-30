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
	ViewerCaps = sortCaps(append(append([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_LIST,

		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_LIST,
		Capability_CAP_IAM_GROUP_INVITES_LIST,
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_IAM_ROLE_BINDINGS_LIST,
		Capability_CAP_IAM_ROLES_LIST,
		Capability_CAP_IAM_IDENTITY_LIST,
		Capability_CAP_IAM_IDENTITY_PROVIDERS_LIST,

		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,

		Capability_CAP_VERSION_LIST,

		Capability_CAP_VULN_REPORT_LIST,

		Capability_CAP_BUILD_REPORT_LIST,

		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,

		Capability_CAP_REGISTRY_ENTITLEMENTS_LIST,
	},
		// Viewers can also list repos and tags, and pull images.
		RegistryPullCaps...), APKPullCaps...))

	// editorCaps can modify state, but not grant roles/permissions.
	EditorCaps = sortCaps(append([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_CREATE,
		Capability_CAP_EVENTS_SUBSCRIPTION_DELETE,
		Capability_CAP_EVENTS_SUBSCRIPTION_UPDATE,
		Capability_CAP_REPO_UPDATE,
	}, ViewerCaps...))

	// ownerCaps includes all capabilities possible by a user.
	OwnerCaps = sortCaps(append(append(append(append(append([]Capability{
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_CREATE,
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_DELETE,
		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_UPDATE,

		Capability_CAP_IAM_GROUP_INVITES_CREATE,
		Capability_CAP_IAM_GROUP_INVITES_DELETE,

		Capability_CAP_IAM_GROUPS_CREATE,
		Capability_CAP_IAM_GROUPS_DELETE,
		Capability_CAP_IAM_GROUPS_UPDATE,

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

		Capability_CAP_VULN_CREATE,
		Capability_CAP_VULN_REPORT_CREATE,

		Capability_CAP_LIBRARIES_ENTITLEMENTS_CREATE,
		Capability_CAP_LIBRARIES_ENTITLEMENTS_DELETE,
	}, EditorCaps...),
		// Owners can also push and delete images, subject to the identity allowlist.
		RegistryPushCaps...),
		APKPushCaps...),
		// Owners can pull artifacts from ecosystem libraries and grant this role to others in their org.
		// NB: The org must also be entitled to the ecosystem to pull artifacts.
		LibrariesJavaPullCaps...), LibrariesPythonPullCaps...))

	RegistryPullCaps = sortCaps([]Capability{
		Capability_CAP_IAM_GROUPS_LIST,

		Capability_CAP_REPO_LIST,
		Capability_CAP_MANIFEST_LIST,
		Capability_CAP_TAG_LIST,
		Capability_CAP_MANIFEST_METADATA_LIST,

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

		// To create nested groups as needed on push.
		Capability_CAP_IAM_GROUPS_CREATE,
	}, RegistryPullCaps...))

	// PullTokenCreatorCaps is the minimal set of capabilities to create a pull token.
	PullTokenCreatorCaps = sortCaps([]Capability{
		// To create the token (identity + rolebinding)
		Capability_CAP_IAM_ROLE_BINDINGS_CREATE,
		Capability_CAP_IAM_IDENTITY_CREATE,
		// To validate the role of the token.
		Capability_CAP_IAM_ROLES_LIST,
		// To validate the parent group the token will attach to.
		Capability_CAP_IAM_GROUPS_LIST,
	})

	RegistryPullTokenCreatorCaps = sortCaps(append(append(PullTokenCreatorCaps, RegistryPullCaps...), APKPullCaps...))

	APKPullCaps = sortCaps([]Capability{
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_APK_LIST,
	})

	APKPushCaps = sortCaps(append([]Capability{
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_APK_CREATE,
		Capability_CAP_APK_DELETE,
	}, APKPullCaps...))

	LibrariesJavaPullCaps = sortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_JAVA_LIST,
	})

	LibrariesPythonPullCaps = sortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_PYTHON_LIST,
	})
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
