/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	EnforcerOptionWebhookFailOpen          = "webhook_fail_open"
	EnforcerOptionEnableCIPCache           = "enable_cip_cache"
	EnforcerOptionNamespaceEnforcementMode = "namespace_enforcement_mode"

	NamespaceEnforcementModeOptOut = "opt-out"
	NamespaceEnforcementModeOptIn  = "opt-in"
)

var (
	// ValidNamespaceEnforcementModeOpts are valid webhook label selection options for a managed cluster
	ValidNamespaceEnforcementModeOpts = sets.NewString(NamespaceEnforcementModeOptIn, NamespaceEnforcementModeOptOut)

	// ValidEnforcerOptions are valid field names of enforcer options available
	ValidEnforcerOptions = sets.NewString(EnforcerOptionWebhookFailOpen, EnforcerOptionEnableCIPCache, EnforcerOptionNamespaceEnforcementMode)
)
