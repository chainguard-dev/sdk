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
		"valid https URL":       {"https://prometheus-community.github.io/helm-charts", true},
		"valid https with path": {"https://argoproj.github.io/argo-helm", true},
		"valid http URL":        {"http://example.com/helm", true},
		"valid oci URL":         {"oci://registry.example.com/charts", true},
		"valid file URL":        {"file:///local/path/to/charts", true},
		"valid ftp URL":         {"ftp://example.com/helm", true},
		"empty":                 {"", false},
		"no scheme":             {"example.com/helm", false},
		"invalid URL":           {"not-a-url", false},
		"localhost allowed":     {"https://localhost:8080/helm", true},
		"with port":             {"https://example.com:443/helm", true},
		"with query params":     {"https://example.com/helm?version=1.0", true},
		"no host allowed":       {"https://", true},
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
