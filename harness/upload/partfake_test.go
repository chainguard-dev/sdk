/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"context"
	"hash/crc32"
	"math"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"

	"chainguard.dev/sdk/harness/upload"
	harnessv1alpha2 "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha2"
	"chainguard.dev/sdk/uidp"
)

// The fake below is the v1alpha2 Uploads service as uploads.proto's RPC
// comments define it: WriteUpload's size and offset rules and their codes,
// bytes below committed_size skipped, storage free to keep less of a part
// than it was sent, the last part finalizing the object, and a
// CompleteUpload that refuses an object still missing bytes, otherwise
// compares its size and CRC32C, and returns a closed upload unchanged. Keep
// it in step with uploads.proto.

const (
	alignment = 256 << 10
	// serverMaxPartSize is the max_part_size the server advertises.
	serverMaxPartSize = 8 << 20
	// serverRecvLimit is the server's gRPC receive limit.
	serverRecvLimit = 10 << 20
)

// write is one WriteUpload the server received.
type write struct {
	offset, size int64
}

// partRecord is one upload as the fake holds it.
type partRecord struct {
	upload *harnessv1alpha2.Upload
	// stored is the sealed bytes storage holds from the start.
	stored []byte
}

// partsServer is the v1alpha2 Uploads service.
type partsServer struct {
	harnessv1alpha2.UnimplementedUploadsServer

	key         *harnessv1alpha2.OrgPublicKey
	maxPartSize int64

	// before, when set, is called with each WriteUpload's 1-based number
	// before the part is stored. A non-nil error is returned and nothing is
	// stored.
	before func(ctx context.Context, call int) error
	// keep, when set, returns how many of a part's n new bytes storage
	// keeps. The fake rounds it down to 256 KiB unless it ends the object,
	// as the server answers committed_size.
	keep func(call int, n int64) int64
	// after, when set, is called with each WriteUpload's number once the
	// part is stored. A non-nil error is returned in place of the response.
	after func(ctx context.Context, call int) error
	// beforeComplete and afterComplete are before and after for
	// CompleteUpload, called around closing the upload.
	beforeComplete, afterComplete func(ctx context.Context, call int) error

	mu        sync.Mutex
	records   map[string]*partRecord
	writes    []write
	completes int
	keys      int
	// mutate, when set, changes the stored object when it is finalized.
	mutate func([]byte) []byte
}

func (p *partsServer) GetOrgPublicKey(_ context.Context, req *harnessv1alpha2.GetOrgPublicKeyRequest) (*harnessv1alpha2.OrgPublicKey, error) {
	p.mu.Lock()
	p.keys++
	p.mu.Unlock()
	if !uidp.Valid(req.GetParentId()) || !uidp.InRoot(req.GetParentId()) {
		return nil, status.Error(codes.InvalidArgument, "parent_id is not an organization UIDP")
	}
	return p.key, nil
}

func (p *partsServer) CreateUpload(_ context.Context, req *harnessv1alpha2.CreateUploadRequest) (*harnessv1alpha2.CreateUploadResponse, error) {
	switch {
	case !uidp.Valid(req.GetParentId()):
		return nil, status.Error(codes.InvalidArgument, "group is not a UIDP")
	case req.GetEnvelopeSize() <= 0 || req.GetEnvelopeSize() > maxEnvelopeSize:
		return nil, status.Errorf(codes.InvalidArgument, "envelope_size must be between 1 and %d", maxEnvelopeSize)
	case req.GetBundleSize() <= 0 || req.GetBundleSize() > maxBundleSize:
		return nil, status.Errorf(codes.InvalidArgument, "bundle_size must be between 1 and %d", maxBundleSize)
	case !bundleDigest.MatchString(req.GetBundleSha256()):
		return nil, status.Error(codes.InvalidArgument, `bundle_sha256 must be "sha256:" and 64 lowercase hex digits`)
	}
	rec := &partRecord{upload: &harnessv1alpha2.Upload{
		Id:           string(uidp.UIDP(req.GetParentId()).NewChild()),
		Status:       harnessv1alpha2.UploadStatus_UPLOAD_STATUS_PENDING,
		EnvelopeSize: req.GetEnvelopeSize(),
		BundleSize:   req.GetBundleSize(),
		BundleSha256: req.GetBundleSha256(),
	}}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.records == nil {
		p.records = map[string]*partRecord{}
	}
	p.records[rec.upload.GetId()] = rec
	return &harnessv1alpha2.CreateUploadResponse{Upload: rec.upload, MaxPartSize: p.maxPartSize}, nil
}

// checkID applies what every per-upload RPC checks of parent_id and id.
func checkID(parent, id string) error {
	switch {
	case !uidp.Valid(parent):
		return status.Error(codes.InvalidArgument, "group is not a UIDP")
	case !uidp.Valid(id) || uidp.InRoot(id):
		return status.Error(codes.InvalidArgument, "id is not an upload UIDP")
	case !uidp.IsAncestor(parent, id):
		return status.Error(codes.NotFound, "upload not found")
	}
	return nil
}

func (p *partsServer) WriteUpload(ctx context.Context, req *harnessv1alpha2.WriteUploadRequest) (*harnessv1alpha2.WriteUploadResponse, error) {
	if err := checkID(req.GetParentId(), req.GetId()); err != nil {
		return nil, err
	}
	data, offset := req.GetData(), req.GetOffset()
	p.mu.Lock()
	p.writes = append(p.writes, write{offset: offset, size: int64(len(data))})
	call := len(p.writes)
	p.mu.Unlock()

	if p.before != nil {
		if err := p.before(ctx, call); err != nil {
			return nil, err
		}
	}
	committed, err := p.store(call, req.GetId(), offset, data)
	if err != nil {
		return nil, err
	}
	if p.after != nil {
		if err := p.after(ctx, call); err != nil {
			return nil, err
		}
	}
	return &harnessv1alpha2.WriteUploadResponse{CommittedSize: committed}, nil
}

// store applies WriteUpload's rules to one part and returns committed_size.
func (p *partsServer) store(call int, id string, offset int64, data []byte) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	rec, ok := p.records[id]
	if !ok {
		return 0, status.Error(codes.NotFound, "upload not found")
	}
	if rec.upload.GetStatus() != harnessv1alpha2.UploadStatus_UPLOAD_STATUS_PENDING {
		return 0, status.Error(codes.FailedPrecondition, "the upload no longer accepts bytes")
	}
	size, end := rec.upload.GetEnvelopeSize(), offset+int64(len(data))
	committed := int64(len(rec.stored))
	switch {
	case len(data) == 0:
		return 0, status.Error(codes.InvalidArgument, "the part is empty")
	case int64(len(data)) > p.maxPartSize:
		return 0, status.Error(codes.InvalidArgument, "the part is larger than max_part_size")
	case offset < 0 || offset%alignment != 0:
		return 0, status.Error(codes.InvalidArgument, "offset is not a multiple of 256 KiB")
	case end < size && len(data)%alignment != 0:
		return 0, status.Error(codes.InvalidArgument, "a part before the last is not a multiple of 256 KiB")
	case end > size:
		return 0, status.Error(codes.InvalidArgument, "the part extends past envelope_size")
	case offset > committed:
		return 0, status.Error(codes.FailedPrecondition, "the part starts past committed_size")
	}
	if end <= committed {
		return committed, nil
	}
	fresh := data[committed-offset:]
	n := int64(len(fresh))
	if p.keep != nil {
		n = max(0, min(p.keep(call, n), n))
	}
	if committed+n < size {
		n -= (committed + n) % alignment
	}
	rec.stored = append(rec.stored, fresh[:n]...)
	if int64(len(rec.stored)) == size && p.mutate != nil {
		rec.stored = p.mutate(rec.stored)
	}
	return int64(len(rec.stored)), nil
}

func (p *partsServer) CompleteUpload(ctx context.Context, req *harnessv1alpha2.CompleteUploadRequest) (*harnessv1alpha2.Upload, error) {
	if err := checkID(req.GetParentId(), req.GetId()); err != nil {
		return nil, err
	}
	if req.GetEnvelopeCrc32C() < 0 || req.GetEnvelopeCrc32C() > math.MaxUint32 {
		return nil, status.Error(codes.InvalidArgument, "envelope_crc32c is not a CRC32C")
	}
	p.mu.Lock()
	p.completes++
	call := p.completes
	p.mu.Unlock()

	if p.beforeComplete != nil {
		if err := p.beforeComplete(ctx, call); err != nil {
			return nil, err
		}
	}
	done, err := p.closeUpload(req)
	if err != nil {
		return nil, err
	}
	if p.afterComplete != nil {
		if err := p.afterComplete(ctx, call); err != nil {
			return nil, err
		}
	}
	return done, nil
}

// closeUpload applies CompleteUpload's rules and returns a copy of the
// upload, so gRPC marshals it outside the lock.
func (p *partsServer) closeUpload(req *harnessv1alpha2.CompleteUploadRequest) (*harnessv1alpha2.Upload, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	rec, ok := p.records[req.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "upload not found")
	}
	if rec.upload.GetStatus() == harnessv1alpha2.UploadStatus_UPLOAD_STATUS_PENDING {
		if int64(len(rec.stored)) < rec.upload.GetEnvelopeSize() {
			return nil, status.Error(codes.FailedPrecondition, "the sealed bundle has not all arrived")
		}
		rec.upload.Status = harnessv1alpha2.UploadStatus_UPLOAD_STATUS_READY
		if int64(len(rec.stored)) != rec.upload.GetEnvelopeSize() || int64(crc32.Checksum(rec.stored, castagnoli)) != req.GetEnvelopeCrc32C() {
			rec.upload.Status = harnessv1alpha2.UploadStatus_UPLOAD_STATUS_REJECTED
			rec.upload.RejectionReason = harnessv1alpha2.UploadRejectionReason_UPLOAD_REJECTION_REASON_UPLOAD_MISMATCH
		}
	}
	return proto.CloneOf(rec.upload), nil
}

// only returns the single upload recorded.
func (p *partsServer) only(t *testing.T) *partRecord {
	t.Helper()
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.records) != 1 {
		t.Fatalf("recorded %d uploads, want 1", len(p.records))
	}
	for _, r := range p.records {
		return r
	}
	return nil
}

// stored returns a copy of what storage holds of the single upload.
func (p *partsServer) stored(t *testing.T) []byte {
	t.Helper()
	rec := p.only(t)
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]byte(nil), rec.stored...)
}

// calls returns the WriteUpload calls and the number of CompleteUpload calls.
func (p *partsServer) calls() ([]write, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]write(nil), p.writes...), p.completes
}

// partHarness is a fake v1alpha2 Uploads API reached over gRPC, with a
// connection the test can drop.
type partHarness struct {
	parts  *partsServer
	client harnessv1alpha2.UploadsClient

	mu    sync.Mutex
	conns []net.Conn
	// waits are the delays between attempts Upload asked for.
	waits []time.Duration
}

func newPartHarness(t *testing.T) *partHarness {
	t.Helper()
	h := &partHarness{parts: &partsServer{
		maxPartSize: serverMaxPartSize,
		key: &harnessv1alpha2.OrgPublicKey{
			Pem:        publicPEM(t, &orgKey().PublicKey),
			Algorithm:  "RSA_DECRYPT_OAEP_3072_SHA256",
			KeyVersion: "1",
		},
	}}
	upload.ReplaceSleep(t, func(ctx context.Context, d time.Duration) error {
		h.mu.Lock()
		h.waits = append(h.waits, d)
		h.mu.Unlock()
		return ctx.Err()
	})

	lis := bufconn.Listen(1 << 20)
	gs := grpc.NewServer(grpc.MaxRecvMsgSize(serverRecvLimit))
	harnessv1alpha2.RegisterUploadsServer(gs, h.parts)
	go func() { _ = gs.Serve(lis) }()
	t.Cleanup(gs.Stop)
	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			c, err := lis.DialContext(ctx)
			if err != nil {
				return nil, err
			}
			h.mu.Lock()
			h.conns = append(h.conns, c)
			h.mu.Unlock()
			return c, nil
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	h.client = harnessv1alpha2.NewUploadsClient(conn)
	return h
}

// drop closes every connection the client has dialed, as a network that
// fails mid-request does. gRPC dials again for the next call.
func (h *partHarness) drop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.conns {
		_ = c.Close()
	}
	h.conns = nil
}

// slept returns the delays Upload waited between attempts.
func (h *partHarness) slept() []time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]time.Duration(nil), h.waits...)
}
