/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chainguard.dev/sdk/harness/upload"
	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
	"chainguard.dev/sdk/uidp"
	"chainguard.dev/sdk/uploads"
)

func randomBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("rand: %v", err)
	}
	return b
}

func digest(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func org() string { return string(uidp.NewUIDP("")) }

// opens records how a Bundle's Open was used.
type opens struct {
	calls, closes atomic.Int32
	read          atomic.Int64
}

type trackedReader struct {
	r io.Reader
	o *opens
}

func (r trackedReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.o.read.Add(int64(n))
	return n, err
}

func (r trackedReader) Close() error {
	r.o.closes.Add(1)
	return nil
}

// memBundle describes data truthfully and has Open yield served instead,
// so a test can make the bytes Upload reads differ from the declared ones.
func memBundle(data, served []byte) (*upload.Bundle, *opens) {
	o := &opens{}
	return &upload.Bundle{
		Size:   int64(len(data)),
		SHA256: digest(data),
		Open: func() (io.ReadCloser, error) {
			o.calls.Add(1)
			return trackedReader{r: bytes.NewReader(served), o: o}, nil
		},
	}, o
}

func newClient(h *harness) *upload.Client {
	return &upload.Client{Uploads: h.client, HTTP: h.store.srv.Client()}
}

// open decrypts a stored object with the fake organization's key.
func open(t *testing.T, sealed []byte) []byte {
	t.Helper()
	r, err := uploads.OpenStream(bytes.NewReader(sealed), func(version string, wrapped []byte) ([]byte, error) {
		if version != "1" {
			t.Errorf("key version: got = %q, want = %q", version, "1")
		}
		return rsa.DecryptOAEP(sha256.New(), nil, orgKey(), wrapped, nil)
	})
	if err != nil {
		t.Fatalf("OpenStream: %v", err)
	}
	plain, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading the sealed stream: %v", err)
	}
	return plain
}

// TestUploadDeclaresWhatStorageHolds pins the two values CompleteUpload
// checks: the envelope size declared at CreateUpload and the CRC32C sent at
// CompleteUpload must both describe the stored ciphertext, across the chunk
// boundaries of the sealed stream.
func TestUploadDeclaresWhatStorageHolds(t *testing.T) {
	for _, tc := range []struct {
		name string
		size int
	}{
		{name: "one byte", size: 1},
		{name: "small", size: 100},
		{name: "one chunk", size: uploads.StreamChunkSize},
		{name: "several chunks", size: 2*uploads.StreamChunkSize + 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			data := randomBytes(t, tc.size)
			b, o := memBundle(data, data)

			got, err := newClient(h).Upload(t.Context(), org(), b)
			if err != nil {
				t.Fatalf("Upload: %v", err)
			}
			rec := h.uploads.only(t)
			stored := h.store.object(rec.object)
			if want := int64(len(stored)); rec.upload.GetEnvelopeSize() != want {
				t.Errorf("declared envelope_size: got = %d, want = %d (stored)", rec.upload.GetEnvelopeSize(), want)
			}
			if _, _, completed := h.uploads.counts(); completed != 1 {
				t.Errorf("CompleteUpload calls: got = %d, want = 1", completed)
			}
			if got, want := h.uploads.lastCompletion(t).GetEnvelopeCrc32C(), int64(crc32.Checksum(stored, castagnoli)); got != want {
				t.Errorf("envelope_crc32c: got = %d, want = %d (stored)", got, want)
			}
			if got.GetStatus() != harnesspb.UploadStatus_UPLOAD_STATUS_READY {
				t.Errorf("status: got = %v, want = READY", got.GetStatus())
			}
			if calls, closes := o.calls.Load(), o.closes.Load(); calls != 1 || closes != 1 {
				t.Errorf("bundle opened %d and closed %d times, want 1 and 1", calls, closes)
			}
		})
	}
}

// TestUploadOpensToTheDeclaredBundle checks what a hunt checks after
// decryption: the plaintext is the bundle, with the declared size and digest.
func TestUploadOpensToTheDeclaredBundle(t *testing.T) {
	h := newHarness(t)
	data := randomBytes(t, uploads.StreamChunkSize+11)
	path := filepath.Join(t.TempDir(), "hunt.bundle")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	b, err := upload.DescribeFile(path)
	if err != nil {
		t.Fatalf("DescribeFile: %v", err)
	}

	if _, err := newClient(h).Upload(t.Context(), org(), b); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	rec := h.uploads.only(t)
	plain := open(t, h.store.object(rec.object))
	if !bytes.Equal(plain, data) {
		t.Fatalf("opened %d bytes that differ from the %d-byte bundle", len(plain), len(data))
	}
	if got, want := digest(plain), rec.upload.GetBundleSha256(); got != want {
		t.Errorf("bundle_sha256: got = %s, want = %s", got, want)
	}
	if got, want := int64(len(plain)), rec.upload.GetBundleSize(); got != want {
		t.Errorf("bundle_size: got = %d, want = %d", got, want)
	}
}

// TestUploadReportsAMismatchAsRejected covers storage holding other bytes
// than were sent: the server rejects the upload, and Upload returns it.
func TestUploadReportsAMismatchAsRejected(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func([]byte) []byte
	}{{
		name:   "storage keeps fewer bytes",
		mutate: func(b []byte) []byte { return b[:len(b)-1] },
	}, {
		name:   "storage flips a byte",
		mutate: func(b []byte) []byte { b[len(b)-1] ^= 1; return b },
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			h.store.mutate = tc.mutate
			data := randomBytes(t, 100)
			b, _ := memBundle(data, data)

			got, err := newClient(h).Upload(t.Context(), org(), b)
			if got != nil {
				t.Errorf("Upload returned %v, want nil", got)
			}
			rejected, ok := errors.AsType[*upload.RejectedError](err)
			if !ok {
				t.Fatalf("Upload: got = %v, want a *RejectedError", err)
			}
			if got, want := rejected.Upload.GetStatus(), harnesspb.UploadStatus_UPLOAD_STATUS_REJECTED; got != want {
				t.Errorf("status: got = %v, want = %v", got, want)
			}
			if got, want := rejected.Upload.GetRejectionReason(), harnesspb.UploadRejectionReason_UPLOAD_REJECTION_REASON_UPLOAD_MISMATCH; got != want {
				t.Errorf("rejection reason: got = %v, want = %v", got, want)
			}
			for _, want := range []string{h.uploads.only(t).upload.GetId(), "UPLOAD_STATUS_REJECTED", "UPLOAD_REJECTION_REASON_UPLOAD_MISMATCH"} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not name %q", err, want)
				}
			}
		})
	}
}

// TestUploadRetriesAnAbortedCompletion pins the single retry: the server
// answers Aborted without changing the upload, so one retry readies an
// object already stored, and a second Aborted is returned.
func TestUploadRetriesAnAbortedCompletion(t *testing.T) {
	t.Run("aborted once", func(t *testing.T) {
		h := newHarness(t)
		h.uploads.abort = 1
		data := randomBytes(t, 100)
		b, _ := memBundle(data, data)

		got, err := newClient(h).Upload(t.Context(), org(), b)
		if err != nil {
			t.Fatalf("Upload: %v", err)
		}
		if got.GetStatus() != harnesspb.UploadStatus_UPLOAD_STATUS_READY {
			t.Errorf("status: got = %v, want = READY", got.GetStatus())
		}
		if _, _, completed := h.uploads.counts(); completed != 2 {
			t.Errorf("CompleteUpload calls: got = %d, want = 2", completed)
		}
	})
	t.Run("aborted twice", func(t *testing.T) {
		h := newHarness(t)
		h.uploads.abort = 2
		data := randomBytes(t, 100)
		b, _ := memBundle(data, data)

		_, err := newClient(h).Upload(t.Context(), org(), b)
		if got := status.Code(err); got != codes.Aborted {
			t.Fatalf("Upload: got code %v (%v), want Aborted", got, err)
		}
		if id := h.uploads.only(t).upload.GetId(); !strings.Contains(err.Error(), id) {
			t.Errorf("error %q does not name upload %s", err, id)
		}
		if _, _, completed := h.uploads.counts(); completed != 2 {
			t.Errorf("CompleteUpload calls: got = %d, want = 2", completed)
		}
	})
}

// TestUploadRefusesBeforeAnyRPC covers what Upload checks of its own
// inputs, before it asks the server for anything.
func TestUploadRefusesBeforeAnyRPC(t *testing.T) {
	data := []byte("a\n")
	valid := func() *upload.Bundle { b, _ := memBundle(data, data); return b }
	for _, tc := range []struct {
		name    string
		noAPI   bool
		bundle  func() *upload.Bundle
		wantErr string
	}{{
		name:    "no Uploads client",
		noAPI:   true,
		bundle:  valid,
		wantErr: "no Uploads client",
	}, {
		name:    "nil bundle",
		bundle:  func() *upload.Bundle { return nil },
		wantErr: "nil bundle",
	}, {
		name:    "no Open",
		bundle:  func() *upload.Bundle { b := valid(); b.Open = nil; return b },
		wantErr: "no Open",
	}, {
		name:    "empty bundle",
		bundle:  func() *upload.Bundle { b, _ := memBundle(nil, nil); return b },
		wantErr: "must not be empty",
	}, {
		name:    "negative size",
		bundle:  func() *upload.Bundle { b := valid(); b.Size = -1; return b },
		wantErr: "must not be empty",
	}, {
		name:    "uppercase digest",
		bundle:  func() *upload.Bundle { b := valid(); b.SHA256 = strings.ToUpper(b.SHA256); return b },
		wantErr: "64 lowercase hex digits",
	}, {
		name:    "bare hex digest",
		bundle:  func() *upload.Bundle { b := valid(); b.SHA256 = strings.TrimPrefix(b.SHA256, "sha256:"); return b },
		wantErr: "64 lowercase hex digits",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			c := newClient(h)
			if tc.noAPI {
				c.Uploads = nil
			}

			_, err := c.Upload(t.Context(), org(), tc.bundle())
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			if keys, _, _ := h.uploads.counts(); keys != 0 {
				t.Errorf("GetOrgPublicKey calls: got = %d, want = 0", keys)
			}
		})
	}
}

func TestUploadFailures(t *testing.T) {
	data := randomBytes(t, 100)
	for _, tc := range []struct {
		name   string
		setup  func(*testing.T, *harness)
		served []byte
		open   func() (io.ReadCloser, error)
		// http, when set, returns the client that carries the PUT.
		http         func(*harness) *http.Client
		wantErr      string
		wantCreate   bool
		wantComplete bool
	}{{
		name:       "storage refuses",
		setup:      func(_ *testing.T, h *harness) { h.store.refuse = http.StatusForbidden },
		wantErr:    "storage refused the upload: 403",
		wantCreate: true,
	}, {
		name: "key is not RSA-3072",
		setup: func(t *testing.T, h *harness) {
			k, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("GenerateKey: %v", err)
			}
			h.uploads.key.Pem = publicPEM(t, &k.PublicKey)
		},
		wantErr: "not RSA-3072",
	}, {
		name:    "key is not PEM",
		setup:   func(_ *testing.T, h *harness) { h.uploads.key.Pem = "not a key" },
		wantErr: "not PEM",
	}, {
		name:    "key names another algorithm",
		setup:   func(_ *testing.T, h *harness) { h.uploads.key.Algorithm = "RSA_DECRYPT_OAEP_4096_SHA256" },
		wantErr: "encryption algorithm",
	}, {
		name:    "key version is not a positive integer",
		setup:   func(_ *testing.T, h *harness) { h.uploads.key.KeyVersion = "0" },
		wantErr: "key version",
	}, {
		name:       "bundle does not open",
		open:       func() (io.ReadCloser, error) { return nil, errors.New("no such bundle") },
		wantErr:    "no such bundle",
		wantCreate: true,
	}, {
		// The default client does not trust the fake's certificate, which
		// shows a nil HTTP reaches storage through it.
		name:       "nil HTTP is the default client",
		http:       func(*harness) *http.Client { return nil },
		wantErr:    "certificate",
		wantCreate: true,
	}, {
		name:       "nil Transport is the default transport",
		http:       func(*harness) *http.Client { return &http.Client{} },
		wantErr:    "certificate",
		wantCreate: true,
	}, {
		name:       "transport returns neither a response nor an error",
		http:       func(*harness) *http.Client { return &http.Client{Transport: nilTransport{}} },
		wantErr:    "returned a nil *Response with a nil error",
		wantCreate: true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			if tc.setup != nil {
				tc.setup(t, h)
			}
			served := data
			if tc.served != nil {
				served = tc.served
			}
			b, _ := memBundle(data, served)
			if tc.open != nil {
				b.Open = tc.open
			}
			c := newClient(h)
			if tc.http != nil {
				c.HTTP = tc.http(h)
			}

			_, err := c.Upload(t.Context(), org(), b)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			_, created, completed := h.uploads.counts()
			if got := created > 0; got != tc.wantCreate {
				t.Errorf("CreateUpload called: got = %t, want = %t", got, tc.wantCreate)
			}
			if got := completed > 0; got != tc.wantComplete {
				t.Errorf("CompleteUpload called: got = %t, want = %t", got, tc.wantComplete)
			}
			if tc.wantCreate {
				if id := h.uploads.only(t).upload.GetId(); !strings.Contains(err.Error(), id) {
					t.Errorf("error %q does not name upload %s", err, id)
				}
			}
		})
	}
}

// TestUploadReportsAChangedBundle covers a bundle whose bytes differ from
// what was described. A changed length also fails the PUT, over a transport
// that enforces Content-Length, but the error names the bundle either way,
// and the upload is never completed.
func TestUploadReportsAChangedBundle(t *testing.T) {
	const prefix = "the bundle changed since it was described: "
	small := randomBytes(t, 100)
	large := randomBytes(t, 2*uploads.StreamChunkSize+7)
	other := bytes.Clone(small)
	other[0] ^= 1
	draining := func(*harness) *http.Client { return &http.Client{Transport: drainingTransport{}} }
	for _, tc := range []struct {
		name         string
		data, served []byte
		// http, when set, returns the client that carries the PUT.
		http    func(*harness) *http.Client
		wantErr string
	}{{
		name:    "one byte shorter",
		data:    small,
		served:  small[:len(small)-1],
		wantErr: prefix + "yielded 99 bytes, declared 100",
	}, {
		name:    "one byte longer",
		data:    small,
		served:  append(bytes.Clone(small), 'x'),
		wantErr: prefix + "yielded more than the declared 100 bytes",
	}, {
		name:    "chunks shorter",
		data:    large,
		served:  large[:uploads.StreamChunkSize],
		wantErr: prefix + "yielded 1048576 bytes, declared 2097159",
	}, {
		name:    "chunks longer",
		data:    small,
		served:  large,
		wantErr: prefix + "yielded more than the declared 100 bytes",
	}, {
		name:    "other bytes of the same length",
		data:    small,
		served:  other,
		wantErr: prefix + "yielded digest " + digest(other) + ", declared " + digest(small),
	}, {
		// The PUT succeeds here, so only Upload's own count stops it.
		name:    "shorter, over a transport that ignores Content-Length",
		data:    small,
		served:  small[:len(small)-1],
		http:    draining,
		wantErr: prefix + "yielded 99 bytes, declared 100",
	}, {
		name:    "longer, over a transport that ignores Content-Length",
		data:    small,
		served:  append(bytes.Clone(small), 'x'),
		http:    draining,
		wantErr: prefix + "yielded more than the declared 100 bytes",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			b, o := memBundle(tc.data, tc.served)
			c := newClient(h)
			if tc.http != nil {
				c.HTTP = tc.http(h)
			}

			_, err := c.Upload(t.Context(), org(), b)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			if id := h.uploads.only(t).upload.GetId(); !strings.Contains(err.Error(), id) {
				t.Errorf("error %q does not name upload %s", err, id)
			}
			if _, _, completed := h.uploads.counts(); completed != 0 {
				t.Errorf("CompleteUpload calls: got = %d, want = 0", completed)
			}
			// One byte past the declared size is enough to tell it grew.
			if got, limit := o.read.Load(), int64(len(tc.data))+1; got > limit {
				t.Errorf("read %d bytes of the bundle, want at most %d", got, limit)
			}
		})
	}
}

// TestPutRefusesURLs covers what the client checks of the signed URL itself
// before a byte is sent.
func TestPutRefusesURLs(t *testing.T) {
	for _, tc := range []struct {
		name    string
		edit    func(*harnesspb.UploadUrl)
		wantErr string
	}{{
		name:    "plain http",
		edit:    func(u *harnesspb.UploadUrl) { u.Url = strings.Replace(u.GetUrl(), "https://", "http://", 1) },
		wantErr: "not an https URL",
	}, {
		name:    "no host",
		edit:    func(u *harnesspb.UploadUrl) { u.Url = "https:///object" },
		wantErr: "not an https URL",
	}, {
		name:    "signed for another size",
		edit:    func(u *harnesspb.UploadUrl) { u.Headers["Content-Length"] = "1" },
		wantErr: "signed for 1 bytes",
	}, {
		name: "signed for another size under a lowercase name",
		edit: func(u *harnesspb.UploadUrl) {
			delete(u.Headers, "Content-Length")
			u.Headers["content-length"] = "1"
		},
		wantErr: "signed for 1 bytes",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			data := []byte("a\n")
			b, o := memBundle(data, data)
			client := &editingClient{UploadsClient: h.client, edit: tc.edit}

			_, err := (&upload.Client{Uploads: client, HTTP: h.store.srv.Client()}).Upload(t.Context(), org(), b)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			if obj := h.store.object(h.uploads.only(t).object); obj != nil {
				t.Errorf("stored %d bytes, want nothing sent", len(obj))
			}
			if calls := o.calls.Load(); calls != 0 {
				t.Errorf("bundle opened %d times, want none", calls)
			}
		})
	}
}

// drainingTransport reads each request body to its end and answers 200,
// whatever its length.
type drainingTransport struct{}

func (drainingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := io.Copy(io.Discard, req.Body)
	if err := errors.Join(err, req.Body.Close()); err != nil {
		return nil, err
	}
	return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Body: http.NoBody, Request: req}, nil
}

// TestUploadLeavesTheCallersClientAlone checks that the redirect refusal is
// made on a copy of the caller's client, even when storage redirects.
func TestUploadLeavesTheCallersClientAlone(t *testing.T) {
	h := newHarness(t)
	h.store.redirect = "https://%zz"
	hc := h.store.srv.Client()
	transport := hc.Transport
	data := []byte("a\n")
	b, _ := memBundle(data, data)

	if _, err := (&upload.Client{Uploads: h.client, HTTP: hc}).Upload(t.Context(), org(), b); err == nil {
		t.Fatal("Upload: got nil error, want the redirect refused")
	}
	if hc.Transport != transport {
		t.Errorf("Transport: got = %T, want = %T", hc.Transport, transport)
	}
	if hc.CheckRedirect != nil {
		t.Error("CheckRedirect: got a function, want nil")
	}
}

// nilTransport breaks the RoundTripper contract by returning neither a
// response nor an error.
type nilTransport struct{}

func (nilTransport) RoundTrip(*http.Request) (*http.Response, error) { return nil, nil }

// failingTransport fails every request as a dropped connection does.
type failingTransport struct{}

func (failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("connection reset by peer")
}

// TestUploadKeepsTheSignedURLOutOfErrors checks that no part of the signed
// URL's query, a live upload credential, reaches an error the caller prints.
func TestUploadKeepsTheSignedURLOutOfErrors(t *testing.T) {
	signed := url.Values{
		"X-Goog-Algorithm":     {"GOOG4-RSA-SHA256"},
		"X-Goog-Credential":    {"uploader@example-project.iam.gserviceaccount.com/20261018/auto/storage/goog4_request"},
		"X-Goog-Date":          {"20261018T120000Z"},
		"X-Goog-Expires":       {"900"},
		"X-Goog-SignedHeaders": {"content-length;host;x-goog-if-generation-match"},
		"X-Goog-Signature":     {"5ea1ed0ff1ce5ea1ed0ff1ce"},
	}
	for _, tc := range []struct {
		name string
		// setup returns the client that carries the PUT.
		setup   func(*testing.T, *harness) *http.Client
		wantErr string
	}{{
		name:    "transport fails",
		setup:   func(*testing.T, *harness) *http.Client { return &http.Client{Transport: failingTransport{}} },
		wantErr: "connection reset by peer",
	}, {
		name:    "storage echoes the query",
		setup:   func(_ *testing.T, h *harness) *http.Client { return h.store.srv.Client() },
		wantErr: "storage refused the upload: 403 Forbidden: SignatureDoesNotMatch",
	}, {
		name: "storage redirects",
		setup: func(t *testing.T, h *harness) *http.Client {
			var followed atomic.Int32
			elsewhere := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { followed.Add(1) }))
			t.Cleanup(func() {
				elsewhere.Close()
				if n := followed.Load(); n != 0 {
					t.Errorf("redirect followed %d times, want none", n)
				}
			})
			h.store.redirect = elsewhere.URL
			return h.store.srv.Client()
		},
		wantErr: "storage refused the upload: 302 Found",
	}, {
		// net/http parses a redirect's Location before CheckRedirect runs,
		// and its parse error quotes the whole Location.
		name: "storage redirects to a malformed Location",
		setup: func(_ *testing.T, h *harness) *http.Client {
			h.store.redirect = "https://%zz"
			return h.store.srv.Client()
		},
		wantErr: "storage refused the upload: 302 Found",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			client := &editingClient{UploadsClient: h.client, edit: func(u *harnesspb.UploadUrl) {
				parsed, err := url.Parse(u.GetUrl())
				if err != nil {
					t.Fatalf("Parse: %v", err)
				}
				parsed.RawQuery = signed.Encode()
				u.Url = parsed.String()
			}}
			data := []byte("a\n")
			b, _ := memBundle(data, data)

			_, err := (&upload.Client{Uploads: client, HTTP: tc.setup(t, h)}).Upload(t.Context(), org(), b)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			for _, k := range leakedQuery(err.Error(), signed) {
				t.Errorf("Upload error carries %s: %v", k, err)
			}
		})
	}
}

func TestDescribeFile(t *testing.T) {
	data := randomBytes(t, 3*uploads.StreamChunkSize/2)
	path := filepath.Join(t.TempDir(), "hunt.bundle")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	b, err := upload.DescribeFile(path)
	if err != nil {
		t.Fatalf("DescribeFile: %v", err)
	}
	if got, want := b.Size, int64(len(data)); got != want {
		t.Errorf("Size: got = %d, want = %d", got, want)
	}
	if got, want := b.SHA256, digest(data); got != want {
		t.Errorf("SHA256: got = %s, want = %s", got, want)
	}
	// Open reopens the file from its first byte each time.
	for i := range 2 {
		rc, err := b.Open()
		if err != nil {
			t.Fatalf("Open %d: %v", i, err)
		}
		got, err := io.ReadAll(rc)
		if err := errors.Join(err, rc.Close()); err != nil {
			t.Fatalf("reading Open %d: %v", i, err)
		}
		if !bytes.Equal(got, data) {
			t.Errorf("Open %d yielded %d bytes that differ from the file", i, len(got))
		}
	}

	if err := os.Remove(path); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	rc, err := b.Open()
	if err == nil || rc != nil {
		t.Errorf("Open of a removed file: got (%v, %v), want (nil, error)", rc, err)
	}
	if _, err := upload.DescribeFile(path); err == nil {
		t.Error("DescribeFile of a missing file: got nil error, want an error")
	}
}

// leakedQuery returns the keys of q that msg shows a sign of: a key=value
// pair, raw or query-escaped, or a value long enough to be distinctive. A
// short value alone, such as an expiry of 900 seconds, also turns up by
// chance in a random upload id or a port.
func leakedQuery(msg string, q url.Values) []string {
	const distinctive = 16
	var leaked []string
	for k, vs := range q {
		for _, v := range vs {
			signs := []string{k + "=" + v, url.QueryEscape(k) + "=" + url.QueryEscape(v)}
			if len(v) >= distinctive {
				signs = append(signs, v, url.QueryEscape(v))
			}
			if slices.ContainsFunc(signs, func(sign string) bool { return strings.Contains(msg, sign) }) {
				leaked = append(leaked, k)
			}
		}
	}
	return leaked
}

// editingClient edits the signed URL CreateUpload returns.
type editingClient struct {
	harnesspb.UploadsClient
	edit func(*harnesspb.UploadUrl)
}

func (c *editingClient) CreateUpload(ctx context.Context, req *harnesspb.CreateUploadRequest, opts ...grpc.CallOption) (*harnesspb.CreateUploadResponse, error) {
	resp, err := c.UploadsClient.CreateUpload(ctx, req, opts...)
	if err == nil {
		c.edit(resp.GetUploadUrl())
	}
	return resp, err
}
