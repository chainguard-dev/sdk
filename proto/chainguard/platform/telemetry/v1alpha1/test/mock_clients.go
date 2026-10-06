/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import telemetry "chainguard.dev/sdk/proto/chainguard/platform/telemetry/v1alpha1"

var _ telemetry.Clients = (*MockTelemetryClients)(nil)

type MockTelemetryClients struct {
	TelemetryClient MockTelemetryServiceClient

	OnClose error
}

func (m MockTelemetryClients) Telemetry() telemetry.TelemetryServiceClient {
	return &m.TelemetryClient
}

func (m MockTelemetryClients) Close() error {
	return m.OnClose
}
