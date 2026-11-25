/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	harnesspb "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha1"
	harnessv1alpha2 "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha2"
	"chainguard.dev/sdk/uploads"
)

const (
	// partAlignment is the multiple WriteUpload requires of every offset and
	// of every part's size but the last.
	partAlignment = 256 << 10
	// maxPartSize is the largest part PartClient sends, whatever the server
	// advertises. It bounds the memory a part holds while it waits for an
	// ack, and keeps a request under the 10 MiB a server receives by default.
	maxPartSize = 8 << 20
	// maxAttempts is how many times one part is sent without storage
	// keeping any more of it, or CompleteUpload is called, before Upload
	// gives up.
	maxAttempts = 8
	// firstRetryDelay doubles between attempts up to lastRetryDelay, so a
	// part or a completion rides out about two minutes of outage.
	firstRetryDelay = time.Second
	lastRetryDelay  = 30 * time.Second
)

// sleep waits for d or until ctx is done. Tests replace it.
var sleep = func(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// PartClient seals bundles and uploads them in parts through the harness
// v1alpha2 Uploads API. It is not safe for concurrent use.
type PartClient struct {
	// Uploads is the harness v1alpha2 Uploads API.
	Uploads harnessv1alpha2.UploadsClient
	// Progress, when set, is called on Upload's goroutine with the sealed
	// bytes storage holds and the sealed bundle's size: first with sent 0
	// once the upload exists, then each time that count grows. Its last call
	// has sent equal to total.
	Progress func(sent, total int64)
	// PartSize, when set, is the most one part carries, for a link too slow
	// to send a whole max_part_size part within a request's time limit. It
	// must be a positive multiple of 256 KiB. Zero, or a value over
	// max_part_size or 8 MiB, means the smaller of those two.
	PartSize int64
}

// PartRejectedError is returned by PartClient.Upload when CompleteUpload
// ends in any status other than READY.
type PartRejectedError struct {
	// Upload is the upload as CompleteUpload returned it.
	Upload *harnessv1alpha2.Upload
}

func (e *PartRejectedError) Error() string {
	return fmt.Sprintf("upload %s was not accepted: %s %s", e.Upload.GetId(), e.Upload.GetStatus(), e.Upload.GetRejectionReason())
}

// Upload seals b to org's harness key, writes it in parts, and completes
// the upload. It returns the upload once the server has marked it READY.
//
// A part is kept until storage acknowledges holding it, and a part that
// fails with Unavailable or DeadlineExceeded, or that storage keeps none of,
// is sent again after a backoff. CompleteUpload is repeated the same way on
// those two codes, since completing a closed upload returns it unchanged. The
// sealed bytes are randomized, so a
// failure Upload doesn't retry, or one that outlasts its retries, ends the
// upload: a caller starts again with a new Upload, and the abandoned upload
// is left pending for the server to reclaim.
//
// Upload re-hashes the plaintext as it seals, and holds back the part that
// would finalize the stored object until the bytes it read match b's Size
// and SHA256. Every error after CreateUpload names the upload, and gRPC
// errors stay visible to status.Code.
func (c *PartClient) Upload(ctx context.Context, org string, b *Bundle) (*harnessv1alpha2.Upload, error) {
	if c.Uploads == nil {
		return nil, errors.New("no Uploads client")
	}
	if err := checkBundle(b); err != nil {
		return nil, err
	}
	if c.PartSize < 0 || c.PartSize%partAlignment != 0 {
		return nil, fmt.Errorf("PartSize %d is not a positive multiple of %d", c.PartSize, partAlignment)
	}

	key, err := c.Uploads.GetOrgPublicKey(ctx, &harnessv1alpha2.GetOrgPublicKeyRequest{ParentId: org})
	if err != nil {
		return nil, fmt.Errorf("getting the organization's harness key: %w", err)
	}
	s, err := newSealer(b, &harnesspb.OrgPublicKey{Pem: key.GetPem(), Algorithm: key.GetAlgorithm(), KeyVersion: key.GetKeyVersion()})
	if err != nil {
		return nil, err
	}
	size, err := uploads.SealedSize(s.pub, b.Size)
	if err != nil {
		return nil, fmt.Errorf("sizing the sealed bundle: %w", err)
	}
	created, err := c.Uploads.CreateUpload(ctx, &harnessv1alpha2.CreateUploadRequest{
		ParentId:     org,
		EnvelopeSize: size,
		BundleSize:   b.Size,
		BundleSha256: b.SHA256,
	})
	if err != nil {
		return nil, fmt.Errorf("creating the upload: %w", err)
	}
	id := created.GetUpload().GetId()
	advertised := created.GetMaxPartSize()
	if advertised < partAlignment || advertised%partAlignment != 0 {
		return nil, fmt.Errorf("upload %s: the server's max_part_size %d is not a positive multiple of %d", id, advertised, partAlignment)
	}
	w := &partWriter{
		ctx:      ctx,
		uploads:  c.Uploads,
		org:      org,
		id:       id,
		size:     size,
		partSize: int(min(cmp.Or(c.PartSize, maxPartSize), advertised, maxPartSize)),
		progress: c.Progress,
	}
	if c.Progress != nil {
		// Without this, nothing is reported until the first part lands, so
		// a one-part bundle is seen only at 100%.
		c.Progress(0, size)
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("upload %s: %w", id, err)
		}
	}
	crc, err := w.write(s)
	if err != nil {
		return nil, fmt.Errorf("upload %s: %w", id, err)
	}
	done, err := complete(ctx, c.Uploads, &harnessv1alpha2.CompleteUploadRequest{
		ParentId:       org,
		Id:             id,
		EnvelopeCrc32C: int64(crc),
	})
	if err != nil {
		return nil, fmt.Errorf("completing upload %s: %w", id, err)
	}
	if done.GetStatus() != harnessv1alpha2.UploadStatus_UPLOAD_STATUS_READY {
		return nil, &PartRejectedError{Upload: done}
	}
	return done, nil
}

// partWriter is the io.Writer the sealer writes to. It holds the sealed bytes
// from committed, the count storage holds, and sends them a part at a time
// as it fills, all but the last. The last part finalizes the stored object,
// so write sends it only once the sealer has checked the plaintext.
type partWriter struct {
	ctx      context.Context
	uploads  harnessv1alpha2.UploadsClient
	org, id  string
	size     int64
	partSize int
	progress func(sent, total int64)

	// committed is the sealed bytes storage holds, and buf the ones after
	// them, at most one part.
	committed int64
	buf       []byte
	// err is the first failure, kept apart from the sealer's wrapping.
	err error
}

// write seals into w and sends the last part, returning the CRC32C of the
// sealed bundle.
func (w *partWriter) write(s *sealer) (uint32, error) {
	w.buf = make([]byte, 0, min(int64(w.partSize), w.size))
	got, err := s.seal(w)
	if err := changed(s.bundle, got); err != nil {
		return 0, err
	}
	switch {
	case w.err != nil && errors.Is(err, w.err):
		// The sealer only relayed the write's failure, which says more
		// unwrapped.
		return 0, w.err
	case err != nil:
		return 0, err
	case got.n != w.size:
		// The plaintext matched, so SealedSize and SealStream disagree.
		return 0, fmt.Errorf("sealed %d bytes, declared %d", got.n, w.size)
	}
	for w.committed < w.size {
		if err := w.sendBuf(); err != nil {
			return 0, err
		}
	}
	return got.crc, nil
}

func (w *partWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if w.committed+int64(len(w.buf))+int64(len(p)) > w.size {
		w.err = fmt.Errorf("sealed more than the declared %d bytes", w.size)
		return 0, w.err
	}
	n := 0
	for len(p) > 0 {
		k := min(len(p), w.partSize-len(w.buf))
		w.buf = append(w.buf, p[:k]...)
		p, n = p[k:], n+k
		if len(w.buf) == w.partSize && w.committed+int64(w.partSize) < w.size {
			if w.err = w.sendBuf(); w.err != nil {
				return n, w.err
			}
		}
	}
	return n, nil
}

// sendBuf sends buf, then drops what storage holds of it.
func (w *partWriter) sendBuf() error {
	committed, err := w.send(w.committed, w.buf)
	if err != nil {
		return err
	}
	if committed == w.size {
		// No part follows, so none needs a buffer.
		w.buf = nil
	} else {
		// The next part starts at committed. It goes in a new buffer: gRPC
		// forbids changing a sent message, and a stats handler or
		// interceptor may still read the part after WriteUpload returns.
		w.buf = append(make([]byte, 0, cap(w.buf)), w.buf[committed-w.committed:]...)
	}
	w.committed = committed
	if w.progress != nil {
		w.progress(committed, w.size)
	}
	// Progress can cancel ctx. Stopping here keeps the sealer from reading
	// on until the next part fills, and keeps CompleteUpload from reporting
	// a cancel after the last part only as a Canceled status.
	return w.ctx.Err()
}

// send writes part at offset until storage holds more of it, and returns
// the committed_size that says how much.
func (w *partWriter) send(offset int64, part []byte) (int64, error) {
	end := offset + int64(len(part))
	var last error
	for attempt := range maxAttempts {
		if err := backoff(w.ctx, attempt); err != nil {
			return 0, fmt.Errorf("writing the part at offset %d: %w", offset, err)
		}
		resp, err := w.uploads.WriteUpload(w.ctx, &harnessv1alpha2.WriteUploadRequest{
			ParentId: w.org,
			Id:       w.id,
			Offset:   offset,
			Data:     part,
		})
		if err != nil {
			// A cancelled or expired ctx explains the failure better than
			// the status gRPC made of it.
			if ctxErr := w.ctx.Err(); ctxErr != nil {
				err = ctxErr
			}
			if !retryable(err) {
				return 0, fmt.Errorf("writing the part at offset %d: %w", offset, err)
			}
			last = err
			continue
		}
		committed := resp.GetCommittedSize()
		switch {
		case committed < offset || committed > end:
			return 0, fmt.Errorf("writing the part at offset %d: storage reports %d bytes stored, outside the %d to %d it could hold", offset, committed, offset, end)
		case committed < w.size && committed%partAlignment != 0:
			return 0, fmt.Errorf("writing the part at offset %d: storage reports %d bytes stored, not a multiple of %d", offset, committed, partAlignment)
		case committed == offset:
			last = errors.New("storage kept none of the part")
			continue
		}
		return committed, nil
	}
	return 0, fmt.Errorf("writing the part at offset %d: gave up after %d attempts: %w", offset, maxAttempts, last)
}

// complete closes the upload, repeating req after a failure that may have
// lost it or its answer. The server acks the whole bundle only once it has
// recorded the object finalized, so CompleteUpload finds nothing missing, and
// any refusal it gives would be repeated.
func complete(ctx context.Context, uploads harnessv1alpha2.UploadsClient, req *harnessv1alpha2.CompleteUploadRequest) (*harnessv1alpha2.Upload, error) {
	var last error
	for attempt := range maxAttempts {
		if err := backoff(ctx, attempt); err != nil {
			return nil, err
		}
		done, err := uploads.CompleteUpload(ctx, req)
		if err == nil {
			return done, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
		if !retryable(err) {
			return nil, err
		}
		last = err
	}
	return nil, fmt.Errorf("gave up after %d attempts: %w", maxAttempts, last)
}

// backoff waits before the attempt'th try of a request, counting from zero,
// and returns ctx's error if ctx ends first.
func backoff(ctx context.Context, attempt int) error {
	if attempt == 0 {
		return nil
	}
	return sleep(ctx, min(firstRetryDelay<<(attempt-1), lastRetryDelay))
}

// retryable reports a failure after which the request may or may not have
// taken effect: the connection or a hop behind it failed. A WriteUpload retry
// learns committed_size either way, and a CompleteUpload retry gets the
// closed upload. The caller has already replaced a failure of ctx's own with
// ctx's error, which isn't retried.
func retryable(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}
