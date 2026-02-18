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
func ParseChartLockAttestation(payload []byte) (*Lock, error) {
	var envelope struct {
		Payload []byte `json:"payload"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, fmt.Errorf("parsing DSSE envelope: %w", err)
	}

	var statement struct {
		PredicateType string `json:"predicateType"`
		Predicate     Lock   `json:"predicate"`
	}
	if err := json.Unmarshal(envelope.Payload, &statement); err != nil {
		return nil, fmt.Errorf("parsing in-toto statement: %w", err)
	}

	if statement.PredicateType != PredicateType {
		return nil, ErrChartLockNotFound
	}

	return &statement.Predicate, nil
}
