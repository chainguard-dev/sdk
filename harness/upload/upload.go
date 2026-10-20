/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload

import (
	"cmp"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
	"chainguard.dev/sdk/uploads"
)

// orgKeyBits is the modulus size of RSA_DECRYPT_OAEP_3072_SHA256, the only
// algorithm SealStream accepts.
const orgKeyBits = 3072

var (
	// bundleDigest is the form CreateUpload accepts for bundle_sha256.
	bundleDigest = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	// storageCode matches a GCS error code, such as SignatureDoesNotMatch.
	storageCode = regexp.MustCompile(`^[A-Za-z]{1,64}$`)
)

// Client seals bundles and uploads them through the harness Uploads API.
type Client struct {
	// Uploads is the harness Uploads API.
	Uploads harnesspb.UploadsClient
	// HTTP carries the sealed bundle to the signed URL. Nil means
	// http.DefaultClient. Upload uses a copy that never follows redirects.
	HTTP *http.Client
}

// RejectedError is returned by Upload when CompleteUpload ends in any status
// other than READY.
type RejectedError struct {
	// Upload is the upload as CompleteUpload returned it.
	Upload *harnesspb.Upload
}

func (e *RejectedError) Error() string {
	return fmt.Sprintf("upload %s was not accepted: %s %s", e.Upload.GetId(), e.Upload.GetStatus(), e.Upload.GetRejectionReason())
}

// Upload seals b to org's harness key, uploads it to a signed URL, and
// completes the upload. It returns the upload once the server has marked it
// READY.
//
// Upload re-hashes the plaintext as it seals, and does not complete an
// upload whose bytes differ from b's Size or SHA256. Every error after
// CreateUpload names the upload, and gRPC errors stay visible to
// status.Code.
func (c *Client) Upload(ctx context.Context, org string, b *Bundle) (*harnesspb.Upload, error) {
	switch {
	case c.Uploads == nil:
		return nil, errors.New("no Uploads client")
	case b == nil:
		return nil, errors.New("nil bundle")
	case b.Open == nil:
		return nil, errors.New("the bundle has no Open function")
	case b.Size <= 0:
		return nil, fmt.Errorf("the bundle is %d bytes; it must not be empty", b.Size)
	case !bundleDigest.MatchString(b.SHA256):
		return nil, fmt.Errorf(`the bundle digest %q is not "sha256:" and 64 lowercase hex digits`, b.SHA256)
	}

	key, err := c.Uploads.GetOrgPublicKey(ctx, &harnesspb.GetOrgPublicKeyRequest{ParentId: org})
	if err != nil {
		return nil, fmt.Errorf("getting the organization's harness key: %w", err)
	}
	s, err := newSealer(b, key)
	if err != nil {
		return nil, err
	}
	size, err := uploads.SealedSize(s.pub, b.Size)
	if err != nil {
		return nil, fmt.Errorf("sizing the sealed bundle: %w", err)
	}
	created, err := c.Uploads.CreateUpload(ctx, &harnesspb.CreateUploadRequest{
		ParentId:     org,
		EnvelopeSize: size,
		BundleSize:   b.Size,
		BundleSha256: b.SHA256,
	})
	if err != nil {
		return nil, fmt.Errorf("creating the upload: %w", err)
	}
	id := created.GetUpload().GetId()
	crc, err := c.put(ctx, created.GetUploadUrl(), size, s)
	if err != nil {
		return nil, fmt.Errorf("upload %s: %w", id, err)
	}
	done, err := c.complete(ctx, &harnesspb.CompleteUploadRequest{
		ParentId:       org,
		Id:             id,
		EnvelopeCrc32C: int64(crc),
	})
	if err != nil {
		return nil, fmt.Errorf("completing upload %s: %w", id, err)
	}
	if done.GetStatus() != harnesspb.UploadStatus_UPLOAD_STATUS_READY {
		return nil, &RejectedError{Upload: done}
	}
	return done, nil
}

// complete calls CompleteUpload, retrying once on Aborted: the server aborts
// without changing the upload, and CompleteUpload is idempotent, so a retry
// re-checks the object already stored instead of making the caller upload it
// again.
func (c *Client) complete(ctx context.Context, req *harnesspb.CompleteUploadRequest) (*harnesspb.Upload, error) {
	done, err := c.Uploads.CompleteUpload(ctx, req)
	if status.Code(err) == codes.Aborted {
		done, err = c.Uploads.CompleteUpload(ctx, req)
	}
	return done, err
}

// sealer seals one bundle to one organization key.
type sealer struct {
	bundle     *Bundle
	pub        *rsa.PublicKey
	algorithm  string
	keyVersion string
}

func newSealer(b *Bundle, key *harnesspb.OrgPublicKey) (*sealer, error) {
	block, _ := pem.Decode([]byte(key.GetPem()))
	if block == nil {
		return nil, errors.New("the organization's harness key is not PEM")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing the organization's harness key: %w", err)
	}
	pub, ok := parsed.(*rsa.PublicKey)
	if !ok || pub.N.BitLen() != orgKeyBits {
		return nil, fmt.Errorf("the organization's harness key is not RSA-%d", orgKeyBits)
	}
	// Sealing starts only once the PUT is under way, after CreateUpload.
	// Sealing a header to nowhere first refuses an algorithm or key version
	// SealStream would reject before an upload is created for it.
	if _, err := uploads.SealStream(io.Discard, pub, key.GetAlgorithm(), key.GetKeyVersion()); err != nil {
		return nil, fmt.Errorf("sealing to the organization's harness key: %w", err)
	}
	return &sealer{bundle: b, pub: pub, algorithm: key.GetAlgorithm(), keyVersion: key.GetKeyVersion()}, nil
}

// sealed is what one pass of sealer.seal saw and wrote, as far as it got.
type sealed struct {
	// crc and n are the CRC32C and length of the bytes the destination
	// accepted.
	crc uint32
	n   int64
	// plainN and plainSHA256 describe the plaintext read from the bundle,
	// which is all of it only when readAll is set.
	plainN      int64
	plainSHA256 string
	readAll     bool
}

// seal reads the bundle once and writes it sealed to w. It reports what it
// read and wrote even when it fails.
func (s *sealer) seal(w io.Writer) (sealed, error) {
	rc, err := s.bundle.Open()
	if err != nil {
		return sealed{}, fmt.Errorf("opening the bundle: %w", err)
	}
	defer rc.Close()
	counted := &counter{w: w}
	crc := uploads.NewCRC32CWriter(counted)
	h := sha256.New()
	plain := &counter{w: h}
	got := func(readAll bool) sealed {
		return sealed{
			crc:         crc.Sum32(),
			n:           counted.n,
			plainN:      plain.n,
			plainSHA256: "sha256:" + hex.EncodeToString(h.Sum(nil)),
			readAll:     readAll,
		}
	}
	sw, err := uploads.SealStream(crc, s.pub, s.algorithm, s.keyVersion)
	if err != nil {
		return got(false), fmt.Errorf("sealing the bundle: %w", err)
	}
	// One byte past Size is enough to tell that the bundle grew, and stops a
	// bundle that grew without bound from being read to its end.
	if _, err := io.Copy(sw, io.TeeReader(io.LimitReader(rc, s.bundle.Size+1), plain)); err != nil {
		return got(false), fmt.Errorf("sealing the bundle: %w", err)
	}
	if err := sw.Close(); err != nil {
		return got(true), fmt.Errorf("sealing the bundle: %w", err)
	}
	return got(true), nil
}

// changed reports a plaintext that differs from the bundle as described. It
// returns nil when the plaintext matches, or when a read cut short leaves
// that unknown.
func changed(b *Bundle, got sealed) error {
	const prefix = "the bundle changed since it was described"
	switch {
	case got.plainN > b.Size:
		return fmt.Errorf("%s: yielded more than the declared %d bytes", prefix, b.Size)
	case !got.readAll:
		return nil
	case got.plainN != b.Size:
		return fmt.Errorf("%s: yielded %d bytes, declared %d", prefix, got.plainN, b.Size)
	case got.plainSHA256 != b.SHA256:
		return fmt.Errorf("%s: yielded digest %s, declared %s", prefix, got.plainSHA256, b.SHA256)
	}
	return nil
}

// put seals the bundle into a single PUT to dst and returns the CRC32C of
// what it sent. The URL's query is a live credential, so no error carries it.
func (c *Client) put(ctx context.Context, dst *harnesspb.UploadUrl, size int64, s *sealer) (uint32, error) {
	target, err := url.Parse(dst.GetUrl())
	if err != nil || target.Scheme != "https" || target.Host == "" {
		return 0, errors.New("the signed upload URL is not an https URL")
	}
	pr, pw := io.Pipe()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, target.String(), pr)
	if err != nil {
		return 0, fmt.Errorf("building the upload request: %w", err)
	}
	// Go sends Content-Length from the request, never from the header map.
	req.ContentLength = size
	for k, v := range dst.GetHeaders() {
		if http.CanonicalHeaderKey(k) == "Content-Length" {
			if v != strconv.FormatInt(size, 10) {
				return 0, fmt.Errorf("the upload URL is signed for %s bytes, not %d", v, size)
			}
			continue
		}
		req.Header.Set(k, v)
	}

	type result struct {
		sealed
		err error
	}
	done := make(chan result, 1)
	go func() {
		got, err := s.seal(pw)
		_ = pw.CloseWithError(err)
		done <- result{sealed: got, err: err}
	}()
	// A redirect is a refusal: following it would replay the query elsewhere,
	// and turns the PUT into a GET whose 200 looks like success. The client
	// is a copy, so the caller's is left as it was.
	client := *cmp.Or(c.HTTP, http.DefaultClient)
	client.Transport = noLocation{rt: cmp.Or(client.Transport, http.DefaultTransport)}
	// Unreachable while noLocation strips every Location, but it keeps a
	// redirect from being followed if a 3xx ever arrives with one.
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := client.Do(req)
	// Unblock the sealer if storage stopped reading early.
	_ = pr.Close()
	got := <-done
	if err == nil {
		defer resp.Body.Close()
	}
	// A bundle that changed length also fails the PUT itself, as a body
	// whose length is not the signed one. Naming the bundle is the useful
	// error in that case, and after a 2xx it keeps the server from vouching
	// for a bundle other than the one CreateUpload declared.
	if err := changed(s.bundle, got.sealed); err != nil {
		return 0, err
	}
	if err != nil {
		// url.Error quotes the whole URL; keep only its cause.
		if ue, ok := errors.AsType[*url.Error](err); ok {
			err = ue.Err
		}
		return 0, fmt.Errorf("uploading the sealed bundle to %s: %w", target.Host, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// A GCS error body can echo the canonical query, so report its code only.
		var body struct{ Code string }
		_ = xml.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&body)
		if !storageCode.MatchString(body.Code) {
			return 0, fmt.Errorf("storage refused the upload: %s", resp.Status)
		}
		return 0, fmt.Errorf("storage refused the upload: %s: %s", resp.Status, body.Code)
	}
	switch {
	case got.err != nil:
		return 0, got.err
	case got.n != size:
		// The plaintext matched, so SealedSize and SealStream disagree.
		return 0, fmt.Errorf("sealed %d bytes, declared %d", got.n, size)
	}
	return got.crc, nil
}

// noLocation strips the Location header from 3xx responses. http.Client
// parses a redirect's Location before CheckRedirect can refuse it, and the
// error for a malformed one quotes the whole Location, which storage may
// have built from the signed query. Without a Location, http.Client returns
// the 3xx as it is.
type noLocation struct{ rt http.RoundTripper }

func (t noLocation) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.rt.RoundTrip(req)
	// A nil response is the wrapped transport's contract violation, which
	// http.Client reports as an error once it reaches it.
	if err == nil && resp != nil && resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode < http.StatusBadRequest {
		resp.Header.Del("Location")
	}
	return resp, err
}

// counter counts the bytes its writer accepts.
type counter struct {
	w io.Writer
	n int64
}

func (c *counter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += int64(n)
	return n, err
}
