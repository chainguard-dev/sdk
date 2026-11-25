/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload

import (
	"context"
	"runtime"
	"testing"

	"google.golang.org/grpc"

	harnessv1alpha2 "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha2"
)

// fullAcks is an Uploads API whose storage keeps every part whole.
type fullAcks struct {
	harnessv1alpha2.UploadsClient
}

func (fullAcks) WriteUpload(_ context.Context, req *harnessv1alpha2.WriteUploadRequest, _ ...grpc.CallOption) (*harnessv1alpha2.WriteUploadResponse, error) {
	return &harnessv1alpha2.WriteUploadResponse{CommittedSize: req.GetOffset() + int64(len(req.GetData()))}, nil
}

// TestSendBufAllocatesNoPartAfterTheLast covers the ack that reports the
// whole bundle stored: no part follows it, so sendBuf allocates no buffer
// for one.
func TestSendBufAllocatesNoPartAfterTheLast(t *testing.T) {
	w := &partWriter{
		ctx:      t.Context(),
		uploads:  fullAcks{},
		size:     maxPartSize,
		partSize: maxPartSize,
		buf:      make([]byte, maxPartSize),
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err := w.sendBuf()
	runtime.ReadMemStats(&after)
	if err != nil {
		t.Fatalf("sendBuf: %v", err)
	}
	if n := after.TotalAlloc - before.TotalAlloc; n >= maxPartSize/2 {
		t.Errorf("sendBuf allocated %d bytes after the last part, want far fewer than a %d-byte part", n, maxPartSize)
	}
}
