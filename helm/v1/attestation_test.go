/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"encoding/base64"
	"errors"
	"testing"
)

const chartLockStatement = `{
	"_type": "https://in-toto.io/Statement/v0.1",
	"predicateType": "` + PredicateType + `",
	"predicate": {
		"chart": {"package": "nginx", "ref": "cgr.dev/chainguard/charts/nginx@sha256:abc"},
		"images": {"refs": {"nginx": {"repoName": "nginx", "tag": "latest", "digest": "sha256:def"}}}
	}
}`

func TestParseChartLockAttestationFromPayload(t *testing.T) {
	t.Run("parses a chart-lock statement", func(t *testing.T) {
		lock, err := ParseChartLockAttestationFromPayload([]byte(chartLockStatement))
		if err != nil {
			t.Fatalf("ParseChartLockAttestationFromPayload() = %v", err)
		}
		if got, want := lock.Chart.Package, "nginx"; got != want {
			t.Errorf("Chart.Package: got = %q, want = %q", got, want)
		}
		if got := len(lock.Images.Refs); got != 1 {
			t.Errorf("Images.Refs: got = %d refs, want = 1", got)
		}
	})

	t.Run("other predicate type is not found", func(t *testing.T) {
		statement := []byte(`{"predicateType": "https://example.com/other", "predicate": {}}`)
		if _, err := ParseChartLockAttestationFromPayload(statement); !errors.Is(err, ErrChartLockNotFound) {
			t.Fatalf("ParseChartLockAttestationFromPayload() = %v, want ErrChartLockNotFound", err)
		}
	})

	t.Run("malformed statement errors", func(t *testing.T) {
		if _, err := ParseChartLockAttestationFromPayload([]byte("not json")); err == nil {
			t.Fatal("ParseChartLockAttestationFromPayload() = nil, want error")
		}
	})
}

func TestParseChartLockAttestation(t *testing.T) {
	t.Run("parses a DSSE-wrapped chart-lock", func(t *testing.T) {
		envelope := []byte(`{"payloadType":"application/vnd.in-toto+json","payload":"` +
			base64.StdEncoding.EncodeToString([]byte(chartLockStatement)) + `","signatures":[]}`)
		lock, err := ParseChartLockAttestation(envelope)
		if err != nil {
			t.Fatalf("ParseChartLockAttestation() = %v", err)
		}
		if got, want := lock.Chart.Package, "nginx"; got != want {
			t.Errorf("Chart.Package: got = %q, want = %q", got, want)
		}
	})

	t.Run("malformed envelope errors", func(t *testing.T) {
		if _, err := ParseChartLockAttestation([]byte("not json")); err == nil {
			t.Fatal("ParseChartLockAttestation() = nil, want error")
		}
	})
}
