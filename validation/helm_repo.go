/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"net/netip"
	"net/url"
)

// ErrInvalidHelmRepoURL describes invalid Helm repository URLs
var ErrInvalidHelmRepoURL = fmt.Errorf("helm repository URL must be a valid URL")

// allowedHelmSchemes is the set of URL schemes permitted for Helm repository URLs.
// Only https and oci are allowed in general; http is permitted only for localhost.
var allowedHelmSchemes = map[string]struct{}{
	"https": {},
	"oci":   {},
	"http":  {},
}

// isLoopbackHost reports whether host is a loopback address.
// It accepts the "localhost" hostname and any IP address for which
// netip.Addr.IsLoopback returns true (e.g. 127.0.0.1, ::1).
func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	return addr.IsLoopback()
}

// ValidateHelmRepoURL validates that a Helm repository URL is valid.
// It restricts URLs to https and oci schemes (http is allowed only for localhost),
// requires a non-empty host, and rejects userinfo (credentials in the URL).
func ValidateHelmRepoURL(helmRepoURL string) error {
	if helmRepoURL == "" {
		return fmt.Errorf("helm repository URL cannot be empty")
	}

	// Parse the URL
	u, err := url.ParseRequestURI(helmRepoURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Restrict to allowed schemes only
	if _, ok := allowedHelmSchemes[u.Scheme]; !ok {
		return fmt.Errorf("helm repository URL scheme %q is not allowed; must be one of: https, oci, http (localhost only)", u.Scheme)
	}

	// Require a non-empty host
	if u.Hostname() == "" {
		return fmt.Errorf("helm repository URL must have a non-empty host")
	}

	// Reject userinfo (credentials embedded in the URL)
	if u.User != nil {
		return fmt.Errorf("helm repository URL must not contain userinfo (credentials)")
	}

	// http is only allowed for localhost
	if u.Scheme == "http" {
		host := u.Hostname()
		if !isLoopbackHost(host) {
			return fmt.Errorf("http scheme is only allowed for localhost; use https for remote hosts")
		}
	}

	return nil
}
