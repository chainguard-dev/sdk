/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrChartLockNotFound is returned when no chart-lock attestation exists.
var ErrChartLockNotFound = errors.New("chart-lock attestation not found")

// ParseChartLockAttestation parses a chart-lock predicate from a DSSE-wrapped
// in-toto attestation payload. Returns ErrChartLockNotFound if the payload is
// not a chart-lock attestation.
//
// For an already-unwrapped in-toto statement (e.g. one surfaced by a verifier
// that decodes the DSSE envelope), use ParseChartLockAttestationFromPayload.
func ParseChartLockAttestation(payload []byte) (*Lock, error) {
	var envelope struct {
		Payload []byte `json:"payload"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, fmt.Errorf("parsing DSSE envelope: %w", err)
	}

	return ParseChartLockAttestationFromPayload(envelope.Payload)
}

// ParseChartLockAttestationFromPayload parses a chart-lock predicate from an
// in-toto attestation statement — the DSSE envelope's payload, already
// unwrapped. Returns ErrChartLockNotFound if the statement is not a
// chart-lock attestation.
func ParseChartLockAttestationFromPayload(statement []byte) (*Lock, error) {
	var stmt struct {
		PredicateType string `json:"predicateType"`
		Predicate     Lock   `json:"predicate"`
	}
	if err := json.Unmarshal(statement, &stmt); err != nil {
		return nil, fmt.Errorf("parsing in-toto statement: %w", err)
	}

	if stmt.PredicateType != PredicateType {
		return nil, ErrChartLockNotFound
	}

	return &stmt.Predicate, nil
}
