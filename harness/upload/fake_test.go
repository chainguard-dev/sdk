/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"hash/crc32"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"chainguard.dev/sdk/harness/bundle"
	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
	"chainguard.dev/sdk/uidp"
)

// The fakes below mirror the harness Uploads server: the same request
// validation, a URL pinned to one object, its length and
// x-goog-if-generation-match: 0, and a CompleteUpload that compares the
// stored object's size and CRC32C with the declared size and the client's
// checksum, aborting without a change when the server's store does. Keep
// them in step with the server.

var (
	// orgKey is generated once: RSA-3072 generation is slow under -race.
	orgKey = sync.OnceValue(func() *rsa.PrivateKey {
		k, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			panic(err)
		}
		return k
	})

	castagnoli   = crc32.MakeTable(crc32.Castagnoli)
	bundleDigest = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	// The server's bounds.
	maxBundleSize   = bundle.DefaultLimits().MaxCompressedBytes
	maxEnvelopeSize = maxBundleSize + 1<<20
)

func publicPEM(t *testing.T, pub *rsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("MarshalPKIXPublicKey: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// storage is the intake bucket behind the signed URLs.
type storage struct {
	srv *httptest.Server

	mu      sync.Mutex
	objects map[string][]byte
	// mutate, when set, changes what is stored, as a corrupting transport
	// would.
	mutate func([]byte) []byte
	// refuse, when non-zero, answers every PUT with that status.
	refuse int
	// redirect, when set, answers every request with a 302 to the same path
	// and query on that origin.
	redirect string
}

func newStorage(t *testing.T) *storage {
	s := &storage{objects: map[string][]byte{}}
	s.srv = httptest.NewTLSServer(http.HandlerFunc(s.serve))
	t.Cleanup(s.srv.Close)
	return s
}

// sign returns a URL for object pinned to size, and the headers it covers.
func (s *storage) sign(object string, size int64) *harnesspb.UploadUrl {
	return &harnesspb.UploadUrl{
		Url: s.srv.URL + "/" + object + "?X-Goog-Signature=" + strconv.FormatInt(size, 10),
		Headers: map[string]string{
			"Content-Length":             strconv.FormatInt(size, 10),
			"x-goog-if-generation-match": "0",
		},
	}
}

// serve enforces what a V4 signature over the object, Content-Length and
// x-goog-if-generation-match: 0 enforces in GCS.
func (s *storage) serve(w http.ResponseWriter, r *http.Request) {
	if s.redirect != "" {
		http.Redirect(w, r, s.redirect+r.URL.RequestURI(), http.StatusFound)
		return
	}
	if r.Method != http.MethodPut {
		storageError(w, r, http.StatusForbidden, "SignatureDoesNotMatch")
		return
	}
	signed, err := strconv.ParseInt(r.URL.Query().Get("X-Goog-Signature"), 10, 64)
	if err != nil || r.ContentLength != signed || r.Header.Get("Content-Length") != "" && r.Header.Get("Content-Length") != strconv.FormatInt(signed, 10) {
		storageError(w, r, http.StatusForbidden, "SignatureDoesNotMatch")
		return
	}
	if r.Header.Get("x-goog-if-generation-match") != "0" {
		storageError(w, r, http.StatusForbidden, "SignatureDoesNotMatch")
		return
	}
	if s.refuse != 0 {
		storageError(w, r, s.refuse, "AccessDenied")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, signed+1))
	if err != nil || int64(len(body)) != signed {
		storageError(w, r, http.StatusBadRequest, "IncompleteBody")
		return
	}
	object := strings.TrimPrefix(r.URL.Path, "/")
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.objects[object]; ok {
		storageError(w, r, http.StatusPreconditionFailed, "PreconditionFailed")
		return
	}
	if s.mutate != nil {
		body = s.mutate(body)
	}
	s.objects[object] = body
}

// storageError writes a GCS XML error. Each one echoes the query but the
// signature, as GCS does for a bad signature, so every path is leak-checked.
func storageError(w http.ResponseWriter, r *http.Request, code int, gcsCode string) {
	q := r.URL.Query()
	q.Del("X-Goog-Signature")
	body, err := xml.Marshal(struct {
		XMLName          xml.Name `xml:"Error"`
		Code             string
		CanonicalRequest string
	}{
		Code:             gcsCode,
		CanonicalRequest: r.Method + "\n" + r.URL.Path + "\n" + q.Encode(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

// stat is what the server reads of a stored object: size and CRC32C, never
// the bytes.
func (s *storage) stat(object string) (int64, uint32, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.objects[object]
	return int64(len(b)), crc32.Checksum(b, castagnoli), ok
}

func (s *storage) object(object string) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.objects[object]
}

// record is one upload as the fake datastore holds it.
type record struct {
	upload *harnesspb.Upload
	object string
}

// uploadsServer is the Uploads service.
type uploadsServer struct {
	harnesspb.UnimplementedUploadsServer

	store *storage
	// key is what GetOrgPublicKey returns. Tests set it before Upload runs.
	key *harnesspb.OrgPublicKey

	mu          sync.Mutex
	keyRequests int
	records     map[string]*record
	completed   []*harnesspb.CompleteUploadRequest
	// abort is how many more CompleteUpload calls on a pending upload end
	// Aborted, as the server answers when its store aborts the write.
	abort int
}

func (u *uploadsServer) GetOrgPublicKey(_ context.Context, req *harnesspb.GetOrgPublicKeyRequest) (*harnesspb.OrgPublicKey, error) {
	u.mu.Lock()
	u.keyRequests++
	u.mu.Unlock()
	if !uidp.Valid(req.GetParentId()) || !uidp.InRoot(req.GetParentId()) {
		return nil, status.Error(codes.InvalidArgument, "parent_id is not an organization UIDP")
	}
	return u.key, nil
}

func (u *uploadsServer) CreateUpload(_ context.Context, req *harnesspb.CreateUploadRequest) (*harnesspb.CreateUploadResponse, error) {
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
	rec := &record{
		object: "bundles/" + string(uidp.NewUID()),
		upload: &harnesspb.Upload{
			Id:           string(uidp.UIDP(req.GetParentId()).NewChild()),
			Status:       harnesspb.UploadStatus_UPLOAD_STATUS_PENDING,
			EnvelopeSize: req.GetEnvelopeSize(),
			BundleSize:   req.GetBundleSize(),
			BundleSha256: req.GetBundleSha256(),
		},
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.records == nil {
		u.records = map[string]*record{}
	}
	u.records[rec.upload.GetId()] = rec
	return &harnesspb.CreateUploadResponse{Upload: rec.upload, UploadUrl: u.store.sign(rec.object, req.GetEnvelopeSize())}, nil
}

func (u *uploadsServer) CompleteUpload(_ context.Context, req *harnesspb.CompleteUploadRequest) (*harnesspb.Upload, error) {
	switch {
	case !uidp.Valid(req.GetParentId()):
		return nil, status.Error(codes.InvalidArgument, "group is not a UIDP")
	case !uidp.Valid(req.GetId()) || uidp.InRoot(req.GetId()):
		return nil, status.Error(codes.InvalidArgument, "id is not an upload UIDP")
	case !uidp.IsAncestor(req.GetParentId(), req.GetId()):
		return nil, status.Error(codes.NotFound, "upload not found")
	case req.GetEnvelopeCrc32C() < 0 || req.GetEnvelopeCrc32C() > math.MaxUint32:
		return nil, status.Error(codes.InvalidArgument, "envelope_crc32c is not a CRC32C")
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.completed = append(u.completed, req)
	rec, ok := u.records[req.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "upload not found")
	}
	if rec.upload.GetStatus() == harnesspb.UploadStatus_UPLOAD_STATUS_PENDING {
		size, stored, ok := u.store.stat(rec.object)
		if !ok {
			return nil, status.Error(codes.FailedPrecondition, "the sealed bundle has not arrived")
		}
		if u.abort > 0 {
			u.abort--
			return nil, status.Error(codes.Aborted, "the upload could not be confirmed; retry")
		}
		rec.upload.Status = harnesspb.UploadStatus_UPLOAD_STATUS_READY
		if size != rec.upload.GetEnvelopeSize() || int64(stored) != req.GetEnvelopeCrc32C() {
			rec.upload.Status = harnesspb.UploadStatus_UPLOAD_STATUS_REJECTED
			rec.upload.RejectionReason = harnesspb.UploadRejectionReason_UPLOAD_REJECTION_REASON_UPLOAD_MISMATCH
		}
	}
	return rec.upload, nil
}

// counts returns how many GetOrgPublicKey, CreateUpload and CompleteUpload
// calls the server has seen.
func (u *uploadsServer) counts() (keys, created, completed int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.keyRequests, len(u.records), len(u.completed)
}

// lastCompletion returns the last CompleteUpload request.
func (u *uploadsServer) lastCompletion(t *testing.T) *harnesspb.CompleteUploadRequest {
	t.Helper()
	u.mu.Lock()
	defer u.mu.Unlock()
	if len(u.completed) == 0 {
		t.Fatal("CompleteUpload was never called")
	}
	return u.completed[len(u.completed)-1]
}

// only returns the single upload recorded.
func (u *uploadsServer) only(t *testing.T) *record {
	t.Helper()
	u.mu.Lock()
	defer u.mu.Unlock()
	if len(u.records) != 1 {
		t.Fatalf("recorded %d uploads, want 1", len(u.records))
	}
	for _, r := range u.records {
		return r
	}
	return nil
}

// harness is a fake harness Uploads API and intake bucket.
type harness struct {
	store   *storage
	uploads *uploadsServer
	client  harnesspb.UploadsClient
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	h := &harness{store: newStorage(t)}
	h.uploads = &uploadsServer{store: h.store, key: &harnesspb.OrgPublicKey{
		Pem:        publicPEM(t, &orgKey().PublicKey),
		Algorithm:  "RSA_DECRYPT_OAEP_3072_SHA256",
		KeyVersion: "1",
	}}

	lis := bufconn.Listen(1 << 20)
	gs := grpc.NewServer()
	harnesspb.RegisterUploadsServer(gs, h.uploads)
	go func() { _ = gs.Serve(lis) }()
	t.Cleanup(gs.Stop)
	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	h.client = harnesspb.NewUploadsClient(conn)
	return h
}
