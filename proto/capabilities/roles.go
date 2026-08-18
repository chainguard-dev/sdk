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
	// ConsoleViewerCaps are read-only API capabilities that do not affect state,
	// and cannot pull artifacts from registries.
	ConsoleViewerCaps = SortCaps([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_LIST,

		Capability_CAP_IAM_ACCOUNT_ASSOCIATIONS_LIST,
		Capability_CAP_IAM_GROUP_INVITES_LIST,
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_IAM_ROLE_BINDINGS_LIST,
		Capability_CAP_IAM_ROLES_LIST,
		Capability_CAP_IAM_IDENTITY_LIST,
		Capability_CAP_IAM_IDENTITY_PROVIDERS_LIST,

		Capability_CAP_REPO_LIST,
		Capability_CAP_MANIFEST_LIST,
		Capability_CAP_REFERRERS_LIST,
		Capability_CAP_TAG_LIST,
		Capability_CAP_MANIFEST_METADATA_LIST,
		Capability_CAP_APK_LIST,

		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,

		Capability_CAP_VERSION_LIST,

		Capability_CAP_VULN_REPORT_LIST,

		Capability_CAP_BUILD_REPORT_LIST,

		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_ARTIFACTS_LIST,
		Capability_CAP_LIBRARIES_CACHE_LIST,

		// Reading an org's library build requests changes nothing, so every member
		// can see what has been asked for and where it stands. One capability covers
		// all of those reads — the groups, their items, and the per-library view —
		// matching how a Get reuses its List capability elsewhere in this file.
		// Creating and submitting requests are separate, in the editor and owner sets.
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_LIST,

		Capability_CAP_SKILLS_ENTITLEMENTS_LIST,
		Capability_CAP_SKILLS_LIST,

		Capability_CAP_REGISTRY_ENTITLEMENTS_LIST,
		Capability_CAP_REGISTRY_DEPLOYMENTS_LIST,

		Capability_CAP_REGISTRY_SETTINGS_LIST,

		Capability_CAP_REGISTRY_OVERLAYS_LIST,

		Capability_CAP_POLICIES_POLICY_LIST,
		Capability_CAP_POLICIES_BINDING_LIST,
		Capability_CAP_POLICIES_DECISION_LIST,
		Capability_CAP_POLICIES_OVERRIDE_LIST,

		Capability_CAP_LIBRARIES_POLICY_LIST,
		Capability_CAP_LIBRARIES_POLICY_BINDING_LIST,
		Capability_CAP_LIBRARIES_POLICY_BLOCK_EVENT_LIST,

		Capability_CAP_LIBRARIES_AWS_MARKETPLACE_SUBSCRIPTIONS_LIST,

		Capability_CAP_PACKAGES_ENTITLEMENTS_LIST,
		Capability_CAP_ADVISORIES_LIST,

		Capability_CAP_ACTIONS_LIST,

		Capability_CAP_GUARDENER_ENTITLEMENT_LIST,
		Capability_CAP_GUARDENER_SCAN_LIST,
		Capability_CAP_GUARDENER_SCAN_GET,
	})

	// ViewerCaps are read-only capabilities that do not affect state,
	// but are allowed to pull from container and apk registries.
	ViewerCaps = SortCaps(ConsoleViewerCaps,
		// Any additional capabilities needed to specifically pull from
		// the container or apk registries.
		RegistryPullCaps, APKPullCaps,
		// Viewers can also run guardener (DFC) sessions.
		GuardenerUserCaps)

	// EditorCaps can modify state, but not grant roles/permissions.
	EditorCaps = SortCaps([]Capability{
		Capability_CAP_EVENTS_SUBSCRIPTION_CREATE,
		Capability_CAP_EVENTS_SUBSCRIPTION_DELETE,
		Capability_CAP_EVENTS_SUBSCRIPTION_UPDATE,

		// Building up a library request — creating a draft, renaming it, editing
		// which items it carries, discarding it — is ordinary member work and
		// commits Chainguard to nothing. Only submitting spends build capacity, and
		// that capability lives in OwnerCaps alone. delete here reaches drafts only;
		// removing a submitted group needs the admin deleteSubmitted capability.
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_CREATE,
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_UPDATE,
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_DELETE,
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_ITEMS_UPDATE,
	}, RegistryEditorCaps, ViewerCaps)

	// SkillsPublishCaps is the capability set required to publish skill artifacts
	// to skills.cgr.dev. Included in OwnerCaps so org owners can push after
	// provisioning a skills entitlement. Can also be granted explicitly to a
	// service principal or non-owner member via chainctl iam role-bindings create.
	//
	// CAP_TERMS_LIST lets chainctl's pre-flight read the org's accepted-docs
	// list so it knows whether to prompt with the TUI before remote.Write.
	//
	// The skills.{write,delete} catalog-write caps are internal_only and
	// deliberately NOT here: this bundle is customer-grantable, and the catalog
	// is written by the Chainguard publish pipeline via a CAP_INTERNAL role, not
	// by customer skills publishers.
	//
	// The skills.harden.{read,write,cancel} caps are customer-facing: they let a
	// skills-entitled publisher submit, track, and cancel harden jobs. They do
	// NOT write the hardened output — publishing the hardened artifact and
	// recording the job result is done by internal identities (skills.write is
	// internal_only; the record-result RPC is gated on CAP_INTERNAL), so this
	// bundle stays customer-safe.
	SkillsPublishCaps = SortCaps([]Capability{
		Capability_CAP_SKILLS_PUBLISH,
		Capability_CAP_SKILLS_ENTITLEMENTS_LIST,
		Capability_CAP_SKILLS_HARDEN_READ,
		Capability_CAP_SKILLS_HARDEN_WRITE,
		Capability_CAP_SKILLS_HARDEN_CANCEL,
		Capability_CAP_TERMS_LIST,
	})

	// OwnerCaps includes all capabilities possible by a user.
	OwnerCaps = SortCaps([]Capability{
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
		Capability_CAP_IAM_SCIM_MANAGE,

		Capability_CAP_IAM_ROLE_BINDINGS_CREATE,
		Capability_CAP_IAM_ROLE_BINDINGS_DELETE,
		Capability_CAP_IAM_ROLE_BINDINGS_UPDATE,

		Capability_CAP_IAM_ROLES_CREATE,
		Capability_CAP_IAM_ROLES_DELETE,
		Capability_CAP_IAM_ROLES_UPDATE,

		Capability_CAP_VULN_CREATE,
		Capability_CAP_VULN_REPORT_CREATE,

		Capability_CAP_REGISTRY_ENTITLEMENTS_IMAGES_ADD,
		Capability_CAP_REGISTRY_ENTITLEMENTS_IMAGES_REMOVE,
		Capability_CAP_REGISTRY_ENTITLEMENTS_IMAGES_SWAP,

		Capability_CAP_LIBRARIES_ENTITLEMENTS_CREATE,
		Capability_CAP_LIBRARIES_ENTITLEMENTS_DELETE,

		Capability_CAP_SKILLS_ENTITLEMENTS_CREATE,
		Capability_CAP_SKILLS_ENTITLEMENTS_DELETE,

		Capability_CAP_LIBRARIES_CACHE_INVALIDATE,
		Capability_CAP_REPO_CHECK_POLICIES,

		Capability_CAP_POLICIES_POLICY_CREATE,
		Capability_CAP_POLICIES_POLICY_UPDATE,
		Capability_CAP_POLICIES_POLICY_DELETE,
		Capability_CAP_POLICIES_BINDING_CREATE,
		Capability_CAP_POLICIES_BINDING_UPDATE,
		Capability_CAP_POLICIES_BINDING_DELETE,
		Capability_CAP_POLICIES_OVERRIDE_CREATE,
		Capability_CAP_POLICIES_OVERRIDE_DELETE,

		Capability_CAP_LIBRARIES_POLICY_CREATE,
		Capability_CAP_LIBRARIES_POLICY_UPDATE,
		Capability_CAP_LIBRARIES_POLICY_DELETE,
		Capability_CAP_LIBRARIES_POLICY_BINDING_CREATE,
		Capability_CAP_LIBRARIES_POLICY_BINDING_UPDATE,
		Capability_CAP_LIBRARIES_POLICY_BINDING_DELETE,

		Capability_CAP_TERMS_ACCEPT,
		Capability_CAP_TERMS_LIST,

		Capability_CAP_LIBRARIES_AWS_MARKETPLACE_SUBSCRIPTIONS_CREATE,
		Capability_CAP_LIBRARIES_AWS_MARKETPLACE_SUBSCRIPTIONS_UPDATE,

		// Owner-only, per the product requirement. Submitting is the point where a
		// draft stops being a preview and starts consuming Chainguard build
		// capacity, and it carries the CVE-remediation opt-in for the whole group.
		// Members can prepare a request; committing to it is the owner's call.
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_SUBMIT,

		Capability_CAP_MCP_TOOL_CALL,
	}, EditorCaps,
		// Owners can also push and delete OCI, APK Java, JavaScript, Python, and .NET artifacts, subject to the identity allowlist.
		RegistryPushCaps, APKPushCaps, LibrariesJavaPushCaps, LibrariesJavascriptPushCaps, LibrariesPythonPushCaps, LibrariesDotnetPushCaps, LibrariesGoPushCaps,
		// Owners can pull artifacts from ecosystem libraries and grant this role to others in their org.
		// NB: The org must also be entitled to the ecosystem to pull artifacts.
		LibrariesJavaPullCaps, LibrariesPythonPullCaps, LibrariesJavascriptPullCaps, LibrariesDotnetPullCaps, LibrariesGoPullCaps,
		// Owners can publish skill artifacts to skills.cgr.dev, subject to the org
		// having a skills entitlement. SkillsPublishCaps can also be granted
		// independently to a service principal or non-owner via role-bindings.
		SkillsPublishCaps,
		// Owners get the full guardener admin set, so they can self-service
		// link/unlink a GitHub org to the group in addition to running DFC sessions.
		GuardenerAdminCaps,
	)

	RegistryRepoAdminCaps = SortCaps(RegistryEditorCaps)

	RegistryEditorCaps = SortCaps([]Capability{
		Capability_CAP_REPO_CREATE,
		Capability_CAP_REPO_UPDATE,
		Capability_CAP_REPO_DELETE,

		// Tag-scoped Custom Assembly: whoever can edit overlays can also
		// list them.
		Capability_CAP_REGISTRY_OVERLAYS_LIST,
		Capability_CAP_REGISTRY_OVERLAYS_EDIT,
	}, RegistryPullCaps)

	RegistryPullCaps = SortCaps([]Capability{
		Capability_CAP_REPO_BLOBS_GET,

		Capability_CAP_IAM_GROUPS_LIST,

		Capability_CAP_REPO_LIST,
		Capability_CAP_MANIFEST_LIST,
		Capability_CAP_REFERRERS_LIST,
		Capability_CAP_TAG_LIST,
		Capability_CAP_MANIFEST_METADATA_LIST,

		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,
	})

	RegistryPushCaps = SortCaps([]Capability{
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
	}, RegistryPullCaps)

	// PullTokenCreatorCaps is the minimal set of capabilities to create a pull token.
	PullTokenCreatorCaps = SortCaps([]Capability{
		// To create the token (identity + role binding)
		Capability_CAP_IAM_ROLE_BINDINGS_CREATE,
		Capability_CAP_IAM_IDENTITY_CREATE,
		// To validate the role of the token.
		Capability_CAP_IAM_ROLES_LIST,
		// To validate the parent group the token will attach to.
		Capability_CAP_IAM_GROUPS_LIST,
	})

	RegistryPullTokenCreatorCaps = SortCaps(PullTokenCreatorCaps, RegistryPullCaps, APKPullCaps)

	APKPullCaps = SortCaps([]Capability{
		Capability_CAP_APK_BLOBS_GET,
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_APK_LIST,
	})

	APKPushCaps = SortCaps([]Capability{
		Capability_CAP_IAM_GROUPS_LIST,
		Capability_CAP_APK_CREATE,
		Capability_CAP_APK_DELETE,
	}, APKPullCaps)

	AdvisoriesViewerCaps = SortCaps([]Capability{
		Capability_CAP_ADVISORIES_LIST,
	})

	// ActionsViewerCaps is the minimal capability set for reading the
	// Chainguard Actions catalog. It backs the actions.viewer role, granted to
	// the public-puller identity so the catalog is readable without org
	// membership.
	ActionsViewerCaps = SortCaps([]Capability{
		Capability_CAP_ACTIONS_LIST,
	})

	AdvisoriesCreatorCaps = SortCaps([]Capability{
		Capability_CAP_ADVISORIES_CREATE,
		Capability_CAP_ADVISORIES_UPDATE,
	}, AdvisoriesViewerCaps)

	AdvisoriesApproverCaps = SortCaps([]Capability{
		Capability_CAP_ADVISORIES_APPROVE,
	}, AdvisoriesCreatorCaps)

	AdvisoriesAdminCaps = SortCaps([]Capability{
		Capability_CAP_ADVISORIES_DELETE,
	}, AdvisoriesApproverCaps)

	LibrariesJavaPullCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_JAVA_LIST,
	})

	LibrariesJavaPushCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_JAVA_CREATE,
	}, LibrariesJavaPullCaps)

	LibrariesPythonPullCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_PYTHON_LIST,
	})

	LibrariesPythonPushCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_PYTHON_CREATE,
	}, LibrariesPythonPullCaps)

	LibrariesJavascriptPullCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_JAVASCRIPT_LIST,
	})

	LibrariesJavascriptPushCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_JAVASCRIPT_CREATE,
	}, LibrariesJavascriptPullCaps)

	LibrariesDotnetPullCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_DOTNET_LIST,
	})

	LibrariesDotnetPushCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_DOTNET_CREATE,
	}, LibrariesDotnetPullCaps)

	LibrariesGoPullCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_GO_LIST,
	})

	LibrariesGoPushCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_ENTITLEMENTS_LIST,
		Capability_CAP_LIBRARIES_GO_CREATE,
	}, LibrariesGoPullCaps)

	// The rebuilder capabilities are all (internal_only), so every rebuilder
	// bundle must also grant CAP_INTERNAL — the internal-only validation enforces
	// that pairing, which keeps these capabilities off customer-grantable roles.
	LibrariesRebuilderRequestsCreateCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_CREATE,
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_LIST,
		Capability_CAP_INTERNAL,
	})

	LibrariesRebuilderViewerCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_REMEDIATED_ARTIFACTS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_LIST,
		Capability_CAP_LIBRARIES_REBUILDER_NEW_VERSIONS_READ,
		Capability_CAP_INTERNAL,
	})

	// LibrariesRebuilderTrustedCaps is for the trusted SA (e.g. prepare phase) that
	// can fetch build-scoped JWT tokens.
	LibrariesRebuilderTrustedCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_TOKENS_FETCH,
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_READ,
		Capability_CAP_INTERNAL,
	})

	// LibrariesRebuilderUntrustedCaps is for untrusted build/test pods. They combine
	// this Chainguard identity check with a per-build JWT token for defence in depth.
	// The read-only malware-status capability lets the build pipeline (e.g. the
	// prepare phase) gate which package versions enter the build matrix; it is
	// read-only and low-risk, so it lives on the untrusted role rather than
	// requiring a separate trusted service account.
	LibrariesRebuilderUntrustedCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_UNTRUSTED,
		Capability_CAP_LIBRARIES_REBUILDER_MALWARE_STATUS_READ,
		Capability_CAP_INTERNAL,
	})

	LibrariesRebuilderAdminCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_CREATE,
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_LIST,
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_CANCEL,
		Capability_CAP_LIBRARIES_REBUILDER_REQUESTS_GROUP_UPDATE,
		Capability_CAP_LIBRARIES_REBUILDER_ARTIFACTS_INVALIDATE,
		Capability_CAP_LIBRARIES_REBUILDER_EXCLUSIONS_MANAGE,
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_TOKENS_FETCH,
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_UNTRUSTED,
		Capability_CAP_LIBRARIES_REBUILDER_REMEDIATED_ARTIFACTS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_NEW_VERSIONS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_NEW_VERSIONS_REFRESH,
		Capability_CAP_LIBRARIES_REBUILDER_MALWARE_STATUS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_CVE_REMEDIATIONS_LIST,
		Capability_CAP_LIBRARIES_REBUILDER_CVE_REMEDIATIONS_UPDATE,
		// The two escape hatches on the customer-facing v2beta1 surface, both
		// internal_only so no customer role can carry them. refreshCoverage forces a
		// full re-scan of a group, which customers do not get because submitted
		// groups are kept current by the daily sweep and an on-demand full scan is
		// expensive; it exists for support and incident escalation.
		// deleteSubmitted removes a group the customer has already committed to,
		// which their own delete capability deliberately cannot reach.
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_REFRESH_COVERAGE,
		Capability_CAP_LIBRARIES_REQUEST_GROUPS_DELETE_SUBMITTED,
		// source_coordinates.resolve lets a staff operator enqueue on-demand source
		// resolution (DESIGN-rebuilder-api-integration.md §5.1: "exercised by cg and
		// tests until a first driver is chosen"). Admin already holds
		// artifacts.invalidate / exclusions.manage, so enqueueing a resolve is no
		// privilege escalation here. Note source_coordinates.write is deliberately
		// NOT in this bundle — persisting canonical coordinates is reconciler-only
		// (LibrariesRebuilderSourceResolverCaps), and admin must not be able to
		// forge resolved rows.
		Capability_CAP_LIBRARIES_REBUILDER_SOURCE_COORDINATES_RESOLVE,
		// This bundle grants (internal_only) capabilities, so it must also grant
		// CAP_INTERNAL; the internal-only validation enforces that pairing, which
		// keeps these capabilities off customer-grantable roles.
		Capability_CAP_INTERNAL,
	})

	// LibrariesRebuilderSourceResolverCaps is for the source-coordinate reconciler's
	// identity (the source-resolver-reconciler service): it reads build and source data
	// (builds.read) and writes resolved source coordinates back
	// (source_coordinates.write). Least-privilege: only those two, plus CAP_INTERNAL for
	// the internal-only pairing.
	LibrariesRebuilderSourceResolverCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_BUILDS_READ,
		Capability_CAP_LIBRARIES_REBUILDER_SOURCE_COORDINATES_WRITE,
		Capability_CAP_INTERNAL,
	})

	// LibrariesRebuilderResolveCaps is the enqueue-only counterpart to the
	// reconciler's write bundle above: it lets a caller request on-demand source
	// resolution via rebuilder-api's :resolve endpoint (source_coordinates.resolve),
	// and nothing more. It deliberately excludes source_coordinates.write —
	// enqueueing a resolve must not confer the ability to write canonical
	// coordinates (that stays reconciler-only, LibrariesRebuilderSourceResolverCaps).
	// Granting this rather than the broad admin bundle (which also carries resolve)
	// gives staff / libraries engineers a least-privilege way to drive resolution.
	// CAP_INTERNAL is required because source_coordinates.resolve is (internal_only).
	LibrariesRebuilderResolveCaps = SortCaps([]Capability{
		Capability_CAP_LIBRARIES_REBUILDER_SOURCE_COORDINATES_RESOLVE,
		Capability_CAP_INTERNAL,
	})

	// GuardenerUserCaps is the minimum capability set required to run
	// guardener (DFC) sessions against a group. terms.list is required
	// because the DFC server checks the group's terms-of-service
	// acceptance with the caller's credentials before starting a session.
	// Registry pull caps are included so users can fetch the base images
	// the DFC session converts to. association.list is the read-only
	// counterpart to association.manage, so users can see which GitHub orgs
	// are linked to the group without being able to link or unlink them.
	// actions.migrate lets users enqueue an on-demand GitHub Actions
	// migration for a repository in the group. scan.list/scan.get let
	// users read the dependency scans guardener builds for the group's
	// repositories.
	GuardenerUserCaps = SortCaps([]Capability{
		Capability_CAP_GUARDENER_DFC_CONVERT,
		Capability_CAP_GUARDENER_ASSOCIATION_LIST,
		Capability_CAP_GUARDENER_ACTIONS_MIGRATE,
		Capability_CAP_GUARDENER_SCAN_LIST,
		Capability_CAP_GUARDENER_SCAN_GET,
		Capability_CAP_TERMS_LIST,
	}, RegistryPullCaps)

	// GuardenerAdminCaps extends GuardenerUserCaps with terms.accept so
	// admins can record acceptance of the legal documents that the DFC
	// server requires before sessions can start, and association.manage so
	// admins can self-service link/unlink a GitHub org to the group.
	GuardenerAdminCaps = SortCaps([]Capability{
		Capability_CAP_TERMS_ACCEPT,
		Capability_CAP_GUARDENER_ASSOCIATION_MANAGE,
	}, GuardenerUserCaps)

	MCPToolUserCaps = SortCaps([]Capability{
		Capability_CAP_MCP_TOOL_CALL,
	})

	// ArgosOperatorCaps is the capability set for managing argos
	// (client-side-encrypted document) records.
	ArgosOperatorCaps = SortCaps([]Capability{
		Capability_CAP_ARGOS_DOCUMENTS_CREATE,
		Capability_CAP_ARGOS_DOCUMENTS_DELETE,
		Capability_CAP_ARGOS_DOCUMENTS_LIST,
	})

	// ArgosOSVReaderCaps is the capability set for reading the shared,
	// customer-blind Private OSV corpus. It is kept separate from
	// ArgosOperatorCaps (document management) and is granted to participating
	// orgs via the dedicated argos.osv.reader role, never to staff or org
	// owners by default.
	ArgosOSVReaderCaps = SortCaps([]Capability{
		Capability_CAP_ARGOS_OSV_READ,
	})

	// ArgosOSVDumperCaps is the capability set for bulk-exporting the shared
	// Private OSV corpus via the Dump RPC. The Dump RPC requires
	// CAP_ARGOS_OSV_DUMP.
	ArgosOSVDumperCaps = SortCaps([]Capability{
		Capability_CAP_ARGOS_OSV_DUMP,
	})
)

func SortCaps(caps ...[]Capability) []Capability {
	uniq := map[Capability]struct{}{}
	for _, cs := range caps {
		for _, c := range cs {
			uniq[c] = struct{}{}
		}
	}
	out := maps.Keys(uniq)
	slices.Sort(out) // These are sorted by enum value, not string name
	return out
}
