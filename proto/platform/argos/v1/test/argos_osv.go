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

var _ argos.ArgosOSVClient = (*MockArgosOSVClient)(nil)

type MockArgosOSVClient struct {
	argos.ArgosOSVClient

	OnQuery      []ArgosOSVOnQuery
	OnQueryBatch []ArgosOSVOnQueryBatch
	OnGetVuln    []ArgosOSVOnGetVuln
	OnDump       []ArgosOSVOnDump
}

type ArgosOSVOnQuery struct {
	Given    *argos.OSVQueryRequest
	Response *argos.OSVQueryResponse
	Error    error
}

type ArgosOSVOnQueryBatch struct {
	Given    *argos.OSVQueryBatchRequest
	Response *argos.OSVQueryBatchResponse
	Error    error
}

type ArgosOSVOnGetVuln struct {
	Given  *argos.GetOSVRequest
	Record *argos.OSVRecord
	Error  error
}

type ArgosOSVOnDump struct {
	Given    *argos.DumpOSVRequest
	Messages []*argos.DumpOSVResponse
	Error    error
}

func (m MockArgosOSVClient) Query(_ context.Context, given *argos.OSVQueryRequest, _ ...grpc.CallOption) (*argos.OSVQueryResponse, error) {
	for _, o := range m.OnQuery {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosOSVClient) QueryBatch(_ context.Context, given *argos.OSVQueryBatchRequest, _ ...grpc.CallOption) (*argos.OSVQueryBatchResponse, error) {
	for _, o := range m.OnQueryBatch {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosOSVClient) GetVuln(_ context.Context, given *argos.GetOSVRequest, _ ...grpc.CallOption) (*argos.OSVRecord, error) {
	for _, o := range m.OnGetVuln {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Record, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosOSVClient) Dump(_ context.Context, given *argos.DumpOSVRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[argos.DumpOSVResponse], error) {
	for _, o := range m.OnDump {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &dumpStream{msgs: o.Messages, err: o.Error}, nil
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// dumpStream is a minimal ServerStreamingClient backed by a slice of messages.
type dumpStream struct {
	grpc.ClientStream
	msgs []*argos.DumpOSVResponse
	err  error
}

func (s *dumpStream) Recv() (*argos.DumpOSVResponse, error) {
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

func (s *dumpStream) Header() (metadata.MD, error) { return nil, nil }
func (s *dumpStream) Trailer() metadata.MD         { return nil }
func (s *dumpStream) CloseSend() error             { return nil }
func (s *dumpStream) Context() context.Context     { return context.Background() }
func (s *dumpStream) SendMsg(any) error            { return nil }
func (s *dumpStream) RecvMsg(any) error            { return nil }
