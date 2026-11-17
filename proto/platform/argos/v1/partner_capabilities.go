/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

// EmptyPartnerCapabilitiesVersionToken is the version_token of a partner
// holding no capabilities, including one that never saved. A caller that
// reads NotFound sends it as expected_version_token so its first save is
// still checked.
const EmptyPartnerCapabilitiesVersionToken = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" //nolint:gosec // G101: a digest, not a credential

// MaxPartnerCapabilities is the most capabilities one write may send.
const MaxPartnerCapabilities = 100
