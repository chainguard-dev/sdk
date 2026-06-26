/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateHelmRepoURL(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Expect bool
	}{
		// Valid URLs
		"valid https URL":          {"https://prometheus-community.github.io/helm-charts", true},
		"valid https with path":    {"https://argoproj.github.io/argo-helm", true},
		"valid oci URL":            {"oci://registry.example.com/charts", true},
		"localhost http allowed":   {"http://localhost:8080/helm", true},
		"localhost 127.0.0.1 http": {"http://127.0.0.1:8080/helm", true},
		"localhost ::1 http":       {"http://[::1]:8080/helm", true},
		"with port":                {"https://example.com:443/helm", true},
		"with query params":        {"https://example.com/helm?version=1.0", true},

		// Invalid: empty or no scheme
		"empty":       {"", false},
		"no scheme":   {"example.com/helm", false},
		"invalid URL": {"not-a-url", false},

		// Invalid: disallowed schemes (SSRF / security)
		"file scheme rejected":       {"file:///local/path/to/charts", false},
		"ftp scheme rejected":        {"ftp://example.com/helm", false},
		"javascript scheme rejected": {"javascript:alert(1)", false},
		"data scheme rejected":       {"data:text/plain,hello", false},

		// Invalid: http for non-localhost
		"http non-localhost rejected": {"http://example.com/helm", false},

		// Invalid: empty host
		"no host https rejected": {"https://", false},
		"no host oci rejected":   {"oci://", false},

		// Invalid: userinfo in URL
		"userinfo rejected": {"https://user:pass@example.com/helm", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateHelmRepoURL(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected ValidateHelmRepoURL(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
