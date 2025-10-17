/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"net/url"
)

// ErrInvalidHelmRepoURL describes invalid Helm repository URLs
var ErrInvalidHelmRepoURL = fmt.Errorf("helm repository URL must be a valid URL")

// ValidateHelmRepoURL validates that a Helm repository URL is valid
func ValidateHelmRepoURL(helmRepoURL string) error {
	if helmRepoURL == "" {
		return fmt.Errorf("helm repository URL cannot be empty")
	}

	// Parse the URL
	_, err := url.ParseRequestURI(helmRepoURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	return nil
}
