/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import argos "chainguard.dev/sdk/proto/platform/argos/v1"

var _ argos.Clients = (*MockArgosClients)(nil)

type MockArgosClients struct {
	ArgosDocumentsClient MockArgosDocumentsClient
	ArgosOSVClient       MockArgosOSVClient

	OnClose error
}

func (m MockArgosClients) ArgosDocuments() argos.ArgosDocumentsClient {
	return &m.ArgosDocumentsClient
}

func (m MockArgosClients) ArgosOSV() argos.ArgosOSVClient {
	return &m.ArgosOSVClient
}

func (m MockArgosClients) Close() error {
	return m.OnClose
}
