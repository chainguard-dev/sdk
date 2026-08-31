/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/testing/protocmp"

	argos "chainguard.dev/sdk/proto/platform/argos/v1"
)

var _ argos.ArgosVEXClient = (*MockArgosVEXClient)(nil)

type MockArgosVEXClient struct {
	argos.ArgosVEXClient

	OnGetDocument []ArgosVEXOnGetDocument
	OnDump        []ArgosVEXOnDump
}

type ArgosVEXOnGetDocument struct {
	Given    *argos.GetVEXDocumentRequest
	Document *argos.VEXDocument
	Error    error
}

type ArgosVEXOnDump struct {
	Given    *argos.DumpVEXRequest
	Messages []*argos.DumpVEXResponse
	Error    error
}

func (m MockArgosVEXClient) GetDocument(_ context.Context, given *argos.GetVEXDocumentRequest, _ ...grpc.CallOption) (*argos.VEXDocument, error) {
	for _, o := range m.OnGetDocument {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Document, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosVEXClient) Dump(_ context.Context, given *argos.DumpVEXRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[argos.DumpVEXResponse], error) {
	for _, o := range m.OnDump {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &vexDumpStream{msgs: o.Messages, err: o.Error}, nil
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// vexDumpStream is a minimal ServerStreamingClient backed by a slice of messages.
type vexDumpStream struct {
	grpc.ClientStream
	msgs []*argos.DumpVEXResponse
	err  error
}

func (s *vexDumpStream) Recv() (*argos.DumpVEXResponse, error) {
	if len(s.msgs) == 0 {
		if s.err != nil {
			return nil, s.err
		}
		return nil, io.EOF
	}
	msg := s.msgs[0]
	s.msgs = s.msgs[1:]
	return msg, nil
}

func (s *vexDumpStream) Header() (metadata.MD, error) { return nil, nil }
func (s *vexDumpStream) Trailer() metadata.MD         { return nil }
func (s *vexDumpStream) CloseSend() error             { return nil }
func (s *vexDumpStream) Context() context.Context     { return context.Background() }
func (s *vexDumpStream) SendMsg(any) error            { return nil }
func (s *vexDumpStream) RecvMsg(any) error            { return nil }
