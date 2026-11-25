/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload_test

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	"chainguard.dev/sdk/harness/upload"
	harnessv1alpha2 "chainguard.dev/sdk/proto/chainguard/platform/harness/v1alpha2"
	"chainguard.dev/sdk/uploads"
)

// plainSizeFor returns the plaintext size that seals to exactly envelope
// bytes under the fake organization's key.
func plainSizeFor(t *testing.T, envelope int64) int {
	t.Helper()
	for n := envelope; n > 0; n-- {
		got, err := uploads.SealedSize(&orgKey().PublicKey, n)
		if err != nil {
			t.Fatalf("SealedSize: %v", err)
		}
		if got == envelope {
			return int(n)
		}
		if got < envelope {
			break
		}
	}
	t.Fatalf("no plaintext seals to %d bytes", envelope)
	return 0
}

// checkParts checks every write against WriteUpload's rules and against
// sending each byte after the first part only once storage asked for it:
// each write starts where storage's count stood, which for writes without a
// short ack is where the last one ended.
func checkParts(t *testing.T, writes []write, envelope, partSize int64) {
	t.Helper()
	if len(writes) == 0 {
		t.Fatal("no WriteUpload calls")
	}
	for i, w := range writes {
		if w.offset%alignment != 0 {
			t.Errorf("write %d: offset %d is not a multiple of 256 KiB", i, w.offset)
		}
		if w.size > partSize {
			t.Errorf("write %d: %d bytes, more than the part size %d", i, w.size, partSize)
		}
		if w.offset+w.size < envelope && w.size%alignment != 0 {
			t.Errorf("write %d: %d bytes ends before the envelope and is not a multiple of 256 KiB", i, w.size)
		}
		if i > 0 && w.offset > writes[i-1].offset+writes[i-1].size {
			t.Errorf("write %d: starts at %d, past the previous part's end %d", i, w.offset, writes[i-1].offset+writes[i-1].size)
		}
	}
	if w := writes[len(writes)-1]; w.offset+w.size != envelope {
		t.Errorf("last write ends at %d, want the envelope's end %d", w.offset+w.size, envelope)
	}
}

// wholeParts returns the writes of an envelope sent once in whole parts of
// size part, the last one shorter.
func wholeParts(envelope, part int64) []write {
	writes := make([]write, 0, (envelope+part-1)/part)
	for offset := int64(0); offset < envelope; offset += part {
		writes = append(writes, write{offset: offset, size: min(part, envelope-offset)})
	}
	return writes
}

// checkReady checks the upload completed READY, which the fake grants only
// to a whole object matching CompleteUpload's CRC32C, and that storage holds
// the sealed bundle.
func checkReady(t *testing.T, h *partHarness, got *harnessv1alpha2.Upload, data []byte) {
	t.Helper()
	if got.GetStatus() != harnessv1alpha2.UploadStatus_UPLOAD_STATUS_READY {
		t.Errorf("status: got = %v, want = READY", got.GetStatus())
	}
	if !bytes.Equal(open(t, h.parts.stored(t)), data) {
		t.Error("the stored object does not open to the bundle")
	}
}

// checkUploadErr checks Upload failed with an error containing want that
// names the upload.
func checkUploadErr(t *testing.T, h *partHarness, err error, want string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("Upload: got = %v, want an error containing %q", err, want)
	}
	if id := h.parts.only(t).upload.GetId(); !strings.Contains(err.Error(), id) {
		t.Errorf("error %q does not name upload %s", err, id)
	}
}

func TestPartUploadSendsAlignedParts(t *testing.T) {
	for _, tc := range []struct {
		name string
		// size returns the plaintext size.
		size       func(*testing.T) int
		advertised int64
		// partSize is PartClient.PartSize.
		partSize int64
		// wantPart is the part size the client should use.
		wantPart int64
	}{{
		name:       "one byte",
		size:       func(*testing.T) int { return 1 },
		advertised: serverMaxPartSize,
		wantPart:   serverMaxPartSize,
	}, {
		name:       "several parts of the smallest size",
		size:       func(*testing.T) int { return 3*alignment + 100 },
		advertised: alignment,
		wantPart:   alignment,
	}, {
		// The last part is a whole part, and is held until sealing ends.
		name:       "an envelope of exactly two parts",
		size:       func(t *testing.T) int { return plainSizeFor(t, 4*alignment) },
		advertised: 2 * alignment,
		wantPart:   2 * alignment,
	}, {
		name:       "max_part_size of 8 MiB",
		size:       func(*testing.T) int { return serverMaxPartSize + 100 },
		advertised: serverMaxPartSize,
		wantPart:   serverMaxPartSize,
	}, {
		name:       "max_part_size above 8 MiB is capped",
		size:       func(*testing.T) int { return serverMaxPartSize + 100 },
		advertised: 2 * serverMaxPartSize,
		wantPart:   serverMaxPartSize,
	}, {
		name:       "PartSize below max_part_size",
		size:       func(*testing.T) int { return 3*alignment + 100 },
		advertised: serverMaxPartSize,
		partSize:   alignment,
		wantPart:   alignment,
	}, {
		name:       "PartSize above max_part_size is clamped",
		size:       func(*testing.T) int { return 3*alignment + 100 },
		advertised: alignment,
		partSize:   4 * alignment,
		wantPart:   alignment,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = tc.advertised
			data := randomBytes(t, tc.size(t))
			b, o := memBundle(data, data)

			got, err := (&upload.PartClient{Uploads: h.client, PartSize: tc.partSize}).Upload(t.Context(), org(), b)
			if err != nil {
				t.Fatalf("Upload: %v", err)
			}
			checkReady(t, h, got, data)
			writes, completed := h.parts.calls()
			want := wholeParts(h.parts.only(t).upload.GetEnvelopeSize(), tc.wantPart)
			if diff := cmp.Diff(want, writes, cmp.AllowUnexported(write{})); diff != "" {
				t.Errorf("writes (-want, +got):\n%s", diff)
			}
			if completed != 1 {
				t.Errorf("CompleteUpload calls: got = %d, want = 1", completed)
			}
			if calls, closes := o.calls.Load(), o.closes.Load(); calls != 1 || closes != 1 {
				t.Errorf("bundle opened %d and closed %d times, want 1 and 1", calls, closes)
			}
		})
	}
}

// TestPartUploadResumesFromCommittedSize covers storage keeping less of a
// part than it was sent: the next part starts at committed_size, the object
// still holds every byte once, and each request's Data still holds the bytes
// stored at its offset after Upload returns, as gRPC requires of a message
// once sent.
func TestPartUploadResumesFromCommittedSize(t *testing.T) {
	for _, tc := range []struct {
		name string
		keep func(call int, n int64) int64
	}{{
		name: "every whole part acked one alignment short",
		keep: func(_ int, n int64) int64 {
			if n == 2*alignment {
				return n - alignment
			}
			return n
		},
	}, {
		name: "an unaligned count, rounded down by the server",
		keep: func(call int, n int64) int64 {
			if call == 1 {
				return alignment + 1
			}
			return n
		},
	}, {
		name: "the last part acked short",
		keep: func(_ int, n int64) int64 {
			if n < alignment {
				return n
			}
			return n - 1
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = 2 * alignment
			h.parts.keep = tc.keep
			// The last part is longer than 256 KiB, so it can be acked short.
			data := randomBytes(t, 5*alignment+alignment/2)
			b, _ := memBundle(data, data)
			rec := &recorder{UploadsClient: h.client}

			got, err := (&upload.PartClient{Uploads: rec}).Upload(t.Context(), org(), b)
			if err != nil {
				t.Fatalf("Upload: %v", err)
			}
			checkReady(t, h, got, data)
			writes, _ := h.parts.calls()
			checkParts(t, writes, h.parts.only(t).upload.GetEnvelopeSize(), 2*alignment)
			if len(h.slept()) != 0 {
				t.Errorf("waited %v, want no backoff for a short ack", h.slept())
			}
			stored := h.parts.stored(t)
			for i, req := range rec.sent {
				if want := stored[req.GetOffset() : req.GetOffset()+int64(len(req.GetData()))]; !bytes.Equal(req.GetData(), want) {
					t.Errorf("request %d at offset %d: Data changed after it was sent", i, req.GetOffset())
				}
			}
		})
	}
}

// recorder keeps each WriteUpload request, as a stats handler or
// interceptor may, each committed_size answered, and each CompleteUpload
// request. Upload calls it on its own goroutine, so a test reads it unlocked
// once Upload returns.
type recorder struct {
	harnessv1alpha2.UploadsClient
	sent      []*harnessv1alpha2.WriteUploadRequest
	acked     []int64
	completed []*harnessv1alpha2.CompleteUploadRequest
}

func (r *recorder) CompleteUpload(ctx context.Context, req *harnessv1alpha2.CompleteUploadRequest, opts ...grpc.CallOption) (*harnessv1alpha2.Upload, error) {
	r.completed = append(r.completed, req)
	return r.UploadsClient.CompleteUpload(ctx, req, opts...)
}

func (r *recorder) WriteUpload(ctx context.Context, req *harnessv1alpha2.WriteUploadRequest, opts ...grpc.CallOption) (*harnessv1alpha2.WriteUploadResponse, error) {
	r.sent = append(r.sent, req)
	resp, err := r.UploadsClient.WriteUpload(ctx, req, opts...)
	if err == nil {
		r.acked = append(r.acked, resp.GetCommittedSize())
	}
	return resp, err
}

// TestPartUploadRetriesAFailedPart covers a part whose request or response
// was lost: the part is sent again from committed_size after one backoff.
func TestPartUploadRetriesAFailedPart(t *testing.T) {
	unavailable := status.Error(codes.Unavailable, "the storage write failed")
	for _, tc := range []struct {
		name   string
		setup  func(*partHarness)
		stored bool
	}{{
		name: "Unavailable before the part is stored",
		setup: func(h *partHarness) {
			h.parts.before = func(_ context.Context, call int) error { return failOn(call, 2, unavailable) }
		},
	}, {
		name: "Unavailable after the part is stored",
		setup: func(h *partHarness) {
			h.parts.after = func(_ context.Context, call int) error { return failOn(call, 2, unavailable) }
		},
		stored: true,
	}, {
		name: "DeadlineExceeded from the server",
		setup: func(h *partHarness) {
			h.parts.after = func(_ context.Context, call int) error {
				return failOn(call, 2, status.Error(codes.DeadlineExceeded, "the storage write timed out"))
			}
		},
		stored: true,
	}, {
		name: "the connection drops before the part is stored",
		setup: func(h *partHarness) {
			h.parts.before = func(ctx context.Context, call int) error { return dropOn(ctx, h, call, 2) }
		},
	}, {
		name: "the connection drops after the part is stored",
		setup: func(h *partHarness) {
			h.parts.after = func(ctx context.Context, call int) error { return dropOn(ctx, h, call, 2) }
		},
		stored: true,
	}, {
		name: "storage keeps none of the part",
		setup: func(h *partHarness) {
			h.parts.keep = func(call int, n int64) int64 {
				if call == 2 {
					return 0
				}
				return n
			}
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			tc.setup(h)
			data := randomBytes(t, 3*alignment+100)
			b, _ := memBundle(data, data)

			got, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
			if err != nil {
				t.Fatalf("Upload: %v", err)
			}
			checkReady(t, h, got, data)
			writes, _ := h.parts.calls()
			// The second part is sent twice, from the same offset: storage
			// either kept it and skips the repeat, or didn't.
			want := wholeParts(h.parts.only(t).upload.GetEnvelopeSize(), alignment)
			want = slices.Insert(want, 2, want[1])
			if diff := cmp.Diff(want, writes, cmp.AllowUnexported(write{})); diff != "" {
				t.Errorf("writes (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff([]time.Duration{time.Second}, h.slept()); diff != "" {
				t.Errorf("waits (-want, +got):\n%s", diff)
			}
		})
	}
}

// failOn returns err for the want'th call.
func failOn(call, want int, err error) error {
	if call == want {
		return err
	}
	return nil
}

// dropOn drops the client's connection during the want'th call, and returns
// once the server has seen the call end.
func dropOn(ctx context.Context, h *partHarness, call, want int) error {
	if call != want {
		return nil
	}
	h.drop()
	<-ctx.Done()
	return ctx.Err()
}

// TestPartUploadGivesUp covers a part storage never takes: Upload stops
// after 8 sends, with a bounded backoff, and never completes.
func TestPartUploadGivesUp(t *testing.T) {
	for _, tc := range []struct {
		name     string
		setup    func(*partHarness)
		wantCode codes.Code
		wantErr  string
	}{{
		name: "Unavailable on every attempt",
		setup: func(h *partHarness) {
			h.parts.before = func(context.Context, int) error { return status.Error(codes.Unavailable, "the storage write failed") }
		},
		wantCode: codes.Unavailable,
		wantErr:  "gave up after 8 attempts",
	}, {
		name:     "storage keeps none on every attempt",
		setup:    func(h *partHarness) { h.parts.keep = func(int, int64) int64 { return 0 } },
		wantCode: codes.Unknown,
		wantErr:  "gave up after 8 attempts: storage kept none of the part",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			tc.setup(h)
			data := randomBytes(t, 100)
			b, _ := memBundle(data, data)

			_, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
			checkUploadErr(t, h, err, tc.wantErr)
			if got := status.Code(err); got != tc.wantCode {
				t.Errorf("status.Code: got = %v, want = %v", got, tc.wantCode)
			}
			writes, completed := h.parts.calls()
			if len(writes) != 8 {
				t.Errorf("WriteUpload calls: got = %d, want = 8", len(writes))
			}
			if completed != 0 {
				t.Errorf("CompleteUpload calls: got = %d, want = 0", completed)
			}
			want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 30 * time.Second, 30 * time.Second}
			if diff := cmp.Diff(want, h.slept()); diff != "" {
				t.Errorf("waits (-want, +got):\n%s", diff)
			}
		})
	}
}

// TestPartUploadStopsOnARefusal covers a failure a retry can't mend: Upload
// returns it after one attempt, with its code.
func TestPartUploadStopsOnARefusal(t *testing.T) {
	for _, code := range []codes.Code{codes.InvalidArgument, codes.FailedPrecondition, codes.Internal, codes.PermissionDenied} {
		t.Run(code.String(), func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.before = func(context.Context, int) error { return status.Error(code, "refused") }
			data := randomBytes(t, 100)
			b, _ := memBundle(data, data)

			_, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
			if got := status.Code(err); got != code {
				t.Fatalf("Upload: got code %v (%v), want %v", got, err, code)
			}
			if writes, completed := h.parts.calls(); len(writes) != 1 || completed != 0 {
				t.Errorf("WriteUpload and CompleteUpload calls: got = %d and %d, want = 1 and 0", len(writes), completed)
			}
			if len(h.slept()) != 0 {
				t.Errorf("waited %v, want no backoff", h.slept())
			}
		})
	}
}

// TestPartUploadReportsProgress pins Progress to what storage acknowledged:
// the first call is 0 before any part is sent, each later call is a
// committed_size, the count only grows, the last call is the whole
// envelope, and every call's total is the envelope's size.
func TestPartUploadReportsProgress(t *testing.T) {
	h := newPartHarness(t)
	h.parts.maxPartSize = 2 * alignment
	// One short ack and one lost response, so Progress sees neither a part
	// end storage didn't keep nor a repeat.
	h.parts.keep = func(call int, n int64) int64 {
		if call == 1 {
			return n - alignment
		}
		return n
	}
	h.parts.after = func(_ context.Context, call int) error {
		if call == 3 {
			return status.Error(codes.Unavailable, "the response was lost")
		}
		return nil
	}
	data := randomBytes(t, 5*alignment)
	b, _ := memBundle(data, data)
	rec := &recorder{UploadsClient: h.client}
	// writesBefore is how many parts had been sent at each Progress call.
	var sent, totals, writesBefore []int64

	_, err := (&upload.PartClient{Uploads: rec, Progress: func(n, total int64) {
		sent = append(sent, n)
		totals = append(totals, total)
		writesBefore = append(writesBefore, int64(len(rec.sent)))
	}}).Upload(t.Context(), org(), b)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	envelope := h.parts.only(t).upload.GetEnvelopeSize()
	if diff := cmp.Diff(slices.Repeat([]int64{envelope}, len(sent)), totals); diff != "" {
		t.Errorf("Progress totals (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(append([]int64{0}, slices.Compact(slices.Clone(rec.acked))...), sent); diff != "" {
		t.Errorf("Progress against the committed_size answers (-want, +got):\n%s", diff)
	}
	if len(writesBefore) == 0 || writesBefore[0] != 0 {
		t.Errorf("parts sent before the first Progress call: got = %v, want none", writesBefore)
	}
}

// TestPartUploadReportsOnePart pins that a bundle that fits in one part is
// reported twice, at 0 and at the whole envelope, so a progress bar moves
// rather than first appearing full.
func TestPartUploadReportsOnePart(t *testing.T) {
	h := newPartHarness(t)
	data := randomBytes(t, alignment)
	b, _ := memBundle(data, data)
	type report struct{ Sent, Total int64 }
	var got []report

	_, err := (&upload.PartClient{Uploads: h.client, Progress: func(n, total int64) {
		got = append(got, report{n, total})
	}}).Upload(t.Context(), org(), b)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	envelope := h.parts.only(t).upload.GetEnvelopeSize()
	if diff := cmp.Diff([]report{{0, envelope}, {envelope, envelope}}, got); diff != "" {
		t.Errorf("Progress calls (-want, +got):\n%s", diff)
	}
}

// TestPartUploadStopsWhenCancelled covers ctx ending at the first progress
// report, between parts, after the last, during a part, and during a
// backoff: Upload returns ctx's error, reads and sends nothing more, and
// never completes.
func TestPartUploadStopsWhenCancelled(t *testing.T) {
	for _, tc := range []struct {
		name string
		// setup arranges for cancel to be called once the upload is under
		// way, and returns the Progress callback to use.
		setup func(*testing.T, *partHarness, context.CancelFunc) func(sent, total int64)
		// size is the plaintext size.
		size       int
		wantWrites int
	}{{
		name: "at the first report",
		setup: func(_ *testing.T, _ *partHarness, cancel context.CancelFunc) func(int64, int64) {
			return func(int64, int64) { cancel() }
		},
		size:       3 * alignment,
		wantWrites: 0,
	}, {
		name: "between parts",
		setup: func(_ *testing.T, _ *partHarness, cancel context.CancelFunc) func(int64, int64) {
			return func(sent, _ int64) {
				if sent > 0 {
					cancel()
				}
			}
		},
		size:       3 * alignment,
		wantWrites: 1,
	}, {
		// The part that fills at the first sealed chunk's end leaves too
		// little for another part, so the sealer would go on to read the
		// next chunk.
		name: "between parts, with the sealer's write used up",
		setup: func(_ *testing.T, _ *partHarness, cancel context.CancelFunc) func(int64, int64) {
			return func(sent, _ int64) {
				if sent == uploads.StreamChunkSize {
					cancel()
				}
			}
		},
		size:       2 * uploads.StreamChunkSize,
		wantWrites: uploads.StreamChunkSize / alignment,
	}, {
		name: "after the last part",
		setup: func(_ *testing.T, _ *partHarness, cancel context.CancelFunc) func(int64, int64) {
			return func(sent, total int64) {
				if sent == total {
					cancel()
				}
			}
		},
		size:       3 * alignment,
		wantWrites: 4,
	}, {
		name: "during a part",
		setup: func(_ *testing.T, h *partHarness, cancel context.CancelFunc) func(int64, int64) {
			h.parts.before = func(ctx context.Context, call int) error {
				if call == 2 {
					cancel()
					<-ctx.Done()
				}
				return nil
			}
			return nil
		},
		size:       3 * alignment,
		wantWrites: 2,
	}, {
		name: "during a backoff",
		setup: func(t *testing.T, h *partHarness, cancel context.CancelFunc) func(int64, int64) {
			h.parts.before = func(context.Context, int) error { return status.Error(codes.Unavailable, "the storage write failed") }
			upload.ReplaceSleep(t, func(ctx context.Context, _ time.Duration) error {
				cancel()
				return ctx.Err()
			})
			return nil
		},
		size:       3 * alignment,
		wantWrites: 1,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			data := randomBytes(t, tc.size)
			b, o := memBundle(data, data)
			// readAtCancel is how much of the bundle Upload had read when
			// ctx was cancelled.
			var readAtCancel atomic.Int64
			progress := tc.setup(t, h, func() {
				readAtCancel.Store(o.read.Load())
				cancel()
			})

			_, err := (&upload.PartClient{Uploads: h.client, Progress: progress}).Upload(ctx, org(), b)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Upload: got = %v, want context.Canceled", err)
			}
			writes, completed := h.parts.calls()
			if len(writes) != tc.wantWrites || completed != 0 {
				t.Errorf("WriteUpload and CompleteUpload calls: got = %d and %d, want = %d and 0", len(writes), completed, tc.wantWrites)
			}
			// A cancel at the first report comes before the bundle is opened.
			if calls, closes := o.calls.Load(), o.closes.Load(); calls > 1 || closes != calls {
				t.Errorf("bundle opened %d and closed %d times, want each open closed once", calls, closes)
			}
			if extra := o.read.Load() - readAtCancel.Load(); extra != 0 {
				t.Errorf("read %d bytes of the bundle after ctx was cancelled, want none", extra)
			}
		})
	}
}

// TestSleep covers the real wait between attempts: it lasts its delay, and
// ends early with ctx's error when ctx is cancelled.
func TestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		if err := upload.Sleep(t.Context(), time.Minute); err != nil {
			t.Fatalf("Sleep: %v", err)
		}
		if got := time.Since(start); got != time.Minute {
			t.Errorf("slept %v, want 1m", got)
		}

		ctx, cancel := context.WithCancel(t.Context())
		time.AfterFunc(time.Second, cancel)
		start = time.Now()
		if err := upload.Sleep(ctx, time.Hour); !errors.Is(err, context.Canceled) {
			t.Fatalf("Sleep: got = %v, want context.Canceled", err)
		}
		if got := time.Since(start); got != time.Second {
			t.Errorf("slept %v, want 1s", got)
		}
	})
}

// TestPartUploadRefusesAServerOutsideItsContract covers answers
// uploads.proto rules out. Upload stops rather than send bytes it no longer
// holds or at an offset the server would refuse.
func TestPartUploadRefusesAServerOutsideItsContract(t *testing.T) {
	for _, tc := range []struct {
		name string
		// create edits CreateUpload's response, and write the call'th
		// WriteUpload's.
		create  func(*harnessv1alpha2.CreateUploadResponse)
		write   func(call int, resp *harnessv1alpha2.WriteUploadResponse)
		wantErr string
	}{{
		name:    "no max_part_size",
		create:  func(r *harnessv1alpha2.CreateUploadResponse) { r.MaxPartSize = 0 },
		wantErr: "max_part_size 0 is not a positive multiple of 262144",
	}, {
		name:    "max_part_size below 256 KiB",
		create:  func(r *harnessv1alpha2.CreateUploadResponse) { r.MaxPartSize = alignment / 2 },
		wantErr: "max_part_size 131072 is not a positive multiple",
	}, {
		name:    "max_part_size not a multiple of 256 KiB",
		create:  func(r *harnessv1alpha2.CreateUploadResponse) { r.MaxPartSize = alignment + 1 },
		wantErr: "max_part_size 262145 is not a positive multiple",
	}, {
		name:    "committed_size past the part",
		write:   func(_ int, r *harnessv1alpha2.WriteUploadResponse) { r.CommittedSize += alignment },
		wantErr: "storage reports 524288 bytes stored, outside the 0 to 262144",
	}, {
		name: "committed_size going back",
		write: func(call int, r *harnessv1alpha2.WriteUploadResponse) {
			if call == 2 {
				r.CommittedSize = 0
			}
		},
		wantErr: "storage reports 0 bytes stored, outside the 262144 to 524288",
	}, {
		name:    "committed_size not aligned before the end",
		write:   func(_ int, r *harnessv1alpha2.WriteUploadResponse) { r.CommittedSize-- },
		wantErr: "storage reports 262143 bytes stored, not a multiple of 262144",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			data := randomBytes(t, 3*alignment)
			b, _ := memBundle(data, data)
			client := &editingParts{UploadsClient: h.client, create: tc.create, write: tc.write}

			_, err := (&upload.PartClient{Uploads: client}).Upload(t.Context(), org(), b)
			checkUploadErr(t, h, err, tc.wantErr)
			if _, completed := h.parts.calls(); completed != 0 {
				t.Errorf("CompleteUpload calls: got = %d, want = 0", completed)
			}
		})
	}
}

// editingParts edits the responses of a v1alpha2 Uploads client.
type editingParts struct {
	harnessv1alpha2.UploadsClient
	create func(*harnessv1alpha2.CreateUploadResponse)
	write  func(int, *harnessv1alpha2.WriteUploadResponse)
	writes int
}

func (c *editingParts) CreateUpload(ctx context.Context, req *harnessv1alpha2.CreateUploadRequest, opts ...grpc.CallOption) (*harnessv1alpha2.CreateUploadResponse, error) {
	resp, err := c.UploadsClient.CreateUpload(ctx, req, opts...)
	if err == nil && c.create != nil {
		c.create(resp)
	}
	return resp, err
}

func (c *editingParts) WriteUpload(ctx context.Context, req *harnessv1alpha2.WriteUploadRequest, opts ...grpc.CallOption) (*harnessv1alpha2.WriteUploadResponse, error) {
	resp, err := c.UploadsClient.WriteUpload(ctx, req, opts...)
	c.writes++
	if err == nil && c.write != nil {
		c.write(c.writes, resp)
	}
	return resp, err
}

// TestPartUploadReportsAChangedBundle covers a bundle whose bytes differ
// from what was described. The part that would finalize the object is never
// sent, so storage never holds a whole object for the server to vouch for,
// and the upload is never completed.
func TestPartUploadReportsAChangedBundle(t *testing.T) {
	const prefix = "the bundle changed since it was described: "
	small := randomBytes(t, 100)
	large := randomBytes(t, 3*alignment+7)
	other := bytes.Clone(large)
	other[len(other)-1] ^= 1
	// whole seals to exactly four parts, so its last part fills while
	// sealing is still under way.
	whole := randomBytes(t, plainSizeFor(t, 4*alignment))
	wholeOther := bytes.Clone(whole)
	wholeOther[len(wholeOther)-1] ^= 1
	for _, tc := range []struct {
		name         string
		data, served []byte
		wantErr      string
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
		name:    "several parts longer",
		data:    small,
		served:  large,
		wantErr: prefix + "yielded more than the declared 100 bytes",
	}, {
		name:    "other bytes of the same length, after several parts",
		data:    large,
		served:  other,
		wantErr: prefix + "yielded digest " + digest(other) + ", declared " + digest(large),
	}, {
		// The declared envelope's last part is full and held when the
		// extra sealed bytes arrive.
		name:    "one byte longer, ending on a part boundary",
		data:    whole,
		served:  append(bytes.Clone(whole), 'x'),
		wantErr: prefix + "yielded more than the declared " + strconv.Itoa(len(whole)) + " bytes",
	}, {
		name:    "other bytes of the same length, ending on a part boundary",
		data:    whole,
		served:  wholeOther,
		wantErr: prefix + "yielded digest " + digest(wholeOther) + ", declared " + digest(whole),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			b, o := memBundle(tc.data, tc.served)

			_, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
			checkUploadErr(t, h, err, tc.wantErr)
			if stored, envelope := int64(len(h.parts.stored(t))), h.parts.only(t).upload.GetEnvelopeSize(); stored >= envelope {
				t.Errorf("stored %d of %d bytes, want the object never finalized", stored, envelope)
			}
			if _, completed := h.parts.calls(); completed != 0 {
				t.Errorf("CompleteUpload calls: got = %d, want = 0", completed)
			}
			if got, limit := o.read.Load(), int64(len(tc.data))+1; got > limit {
				t.Errorf("read %d bytes of the bundle, want at most %d", got, limit)
			}
		})
	}
}

// TestPartUploadCompletion covers completing the upload: a mismatch is a
// *PartRejectedError, and a lost answer to the last part doesn't stop the
// upload completing.
func TestPartUploadCompletion(t *testing.T) {
	t.Run("storage holds other bytes", func(t *testing.T) {
		h := newPartHarness(t)
		h.parts.mutate = func(b []byte) []byte { b[len(b)-1] ^= 1; return b }
		data := randomBytes(t, 100)
		b, _ := memBundle(data, data)

		got, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
		if got != nil {
			t.Errorf("Upload returned %v, want nil", got)
		}
		if _, ok := errors.AsType[*upload.PartRejectedError](err); !ok {
			t.Fatalf("Upload: got = %v, want a *PartRejectedError", err)
		}
		for _, want := range []string{h.parts.only(t).upload.GetId(), "UPLOAD_STATUS_REJECTED", "UPLOAD_REJECTION_REASON_UPLOAD_MISMATCH"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not name %q", err, want)
			}
		}
	})
	// The server records the object finalized before it answers, so a retry
	// of the lost answer gets the whole size and completion succeeds.
	t.Run("the last part's response lost", func(t *testing.T) {
		h := newPartHarness(t)
		h.parts.after = func(_ context.Context, call int) error {
			return failOn(call, 1, status.Error(codes.Unavailable, "the response was lost"))
		}
		data := randomBytes(t, 100)
		b, _ := memBundle(data, data)

		got, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
		if err != nil {
			t.Fatalf("Upload: %v", err)
		}
		checkReady(t, h, got, data)
		writes, completed := h.parts.calls()
		last := wholeParts(h.parts.only(t).upload.GetEnvelopeSize(), serverMaxPartSize)[0]
		if diff := cmp.Diff([]write{last, last}, writes, cmp.AllowUnexported(write{})); diff != "" {
			t.Errorf("writes (-want, +got):\n%s", diff)
		}
		if completed != 1 {
			t.Errorf("CompleteUpload calls: got = %d, want = 1", completed)
		}
	})
}

// TestPartUploadRetriesCompletion covers a CompleteUpload whose request or
// response was lost: the same request is sent again after one backoff, no
// part is sent again, and the upload ends READY. Completion is idempotent,
// so a repeat after the server closed the upload returns it unchanged.
func TestPartUploadRetriesCompletion(t *testing.T) {
	unavailable := status.Error(codes.Unavailable, "the completion failed")
	deadline := status.Error(codes.DeadlineExceeded, "the completion timed out")
	for _, tc := range []struct {
		name  string
		setup func(*partHarness)
	}{{
		name: "Unavailable before the upload is closed",
		setup: func(h *partHarness) {
			h.parts.beforeComplete = func(_ context.Context, call int) error { return failOn(call, 1, unavailable) }
		},
	}, {
		name: "Unavailable after the upload is closed",
		setup: func(h *partHarness) {
			h.parts.afterComplete = func(_ context.Context, call int) error { return failOn(call, 1, unavailable) }
		},
	}, {
		name: "DeadlineExceeded before the upload is closed",
		setup: func(h *partHarness) {
			h.parts.beforeComplete = func(_ context.Context, call int) error { return failOn(call, 1, deadline) }
		},
	}, {
		name: "DeadlineExceeded after the upload is closed",
		setup: func(h *partHarness) {
			h.parts.afterComplete = func(_ context.Context, call int) error { return failOn(call, 1, deadline) }
		},
	}, {
		name: "the connection drops before the upload is closed",
		setup: func(h *partHarness) {
			h.parts.beforeComplete = func(ctx context.Context, call int) error { return dropOn(ctx, h, call, 1) }
		},
	}, {
		name: "the connection drops after the upload is closed",
		setup: func(h *partHarness) {
			h.parts.afterComplete = func(ctx context.Context, call int) error { return dropOn(ctx, h, call, 1) }
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.maxPartSize = alignment
			tc.setup(h)
			data := randomBytes(t, 3*alignment+100)
			b, _ := memBundle(data, data)
			rec := &recorder{UploadsClient: h.client}

			got, err := (&upload.PartClient{Uploads: rec}).Upload(t.Context(), org(), b)
			if err != nil {
				t.Fatalf("Upload: %v", err)
			}
			checkReady(t, h, got, data)
			if diff := cmp.Diff(h.parts.only(t).upload, got, protocmp.Transform()); diff != "" {
				t.Errorf("upload (-want, +got):\n%s", diff)
			}
			writes, completed := h.parts.calls()
			if diff := cmp.Diff(wholeParts(h.parts.only(t).upload.GetEnvelopeSize(), alignment), writes, cmp.AllowUnexported(write{})); diff != "" {
				t.Errorf("writes (-want, +got):\n%s", diff)
			}
			if completed != 2 || len(rec.completed) != 2 {
				t.Fatalf("CompleteUpload calls: got = %d, %d received, want = 2", len(rec.completed), completed)
			}
			if diff := cmp.Diff(rec.completed[0], rec.completed[1], protocmp.Transform()); diff != "" {
				t.Errorf("repeated CompleteUpload request (-first, +repeat):\n%s", diff)
			}
			if diff := cmp.Diff([]time.Duration{time.Second}, h.slept()); diff != "" {
				t.Errorf("waits (-want, +got):\n%s", diff)
			}
		})
	}
}

// TestPartUploadCompletionGivesUp covers a CompleteUpload that never gets
// through: Upload stops after 8 attempts, with the parts' backoff, and keeps
// the last failure's code.
func TestPartUploadCompletionGivesUp(t *testing.T) {
	h := newPartHarness(t)
	h.parts.beforeComplete = func(context.Context, int) error { return status.Error(codes.Unavailable, "the completion failed") }
	data := randomBytes(t, 100)
	b, _ := memBundle(data, data)

	_, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
	checkUploadErr(t, h, err, "gave up after 8 attempts")
	if got := status.Code(err); got != codes.Unavailable {
		t.Errorf("status.Code: got = %v, want = %v", got, codes.Unavailable)
	}
	writes, completed := h.parts.calls()
	if len(writes) != 1 || completed != 8 {
		t.Errorf("WriteUpload and CompleteUpload calls: got = %d and %d, want = 1 and 8", len(writes), completed)
	}
	want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 30 * time.Second, 30 * time.Second}
	if diff := cmp.Diff(want, h.slept()); diff != "" {
		t.Errorf("waits (-want, +got):\n%s", diff)
	}
}

// TestPartUploadCompletionStopsOnARefusal covers a CompleteUpload failure a
// retry can't mend: Upload returns it after one attempt, with its code.
func TestPartUploadCompletionStopsOnARefusal(t *testing.T) {
	for _, code := range []codes.Code{codes.InvalidArgument, codes.FailedPrecondition, codes.NotFound, codes.Internal, codes.PermissionDenied} {
		t.Run(code.String(), func(t *testing.T) {
			h := newPartHarness(t)
			h.parts.beforeComplete = func(context.Context, int) error { return status.Error(code, "refused") }
			data := randomBytes(t, 100)
			b, _ := memBundle(data, data)

			_, err := (&upload.PartClient{Uploads: h.client}).Upload(t.Context(), org(), b)
			if got := status.Code(err); got != code {
				t.Fatalf("Upload: got code %v (%v), want %v", got, err, code)
			}
			if _, completed := h.parts.calls(); completed != 1 {
				t.Errorf("CompleteUpload calls: got = %d, want = 1", completed)
			}
			if len(h.slept()) != 0 {
				t.Errorf("waited %v, want no backoff", h.slept())
			}
		})
	}
}

// TestPartUploadCompletionStopsWhenCancelled covers ctx ending during
// CompleteUpload and during its backoff: Upload returns ctx's error and
// doesn't ask again.
func TestPartUploadCompletionStopsWhenCancelled(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*testing.T, *partHarness, context.CancelFunc)
	}{{
		name: "during CompleteUpload",
		setup: func(_ *testing.T, h *partHarness, cancel context.CancelFunc) {
			h.parts.beforeComplete = func(ctx context.Context, _ int) error {
				cancel()
				<-ctx.Done()
				return ctx.Err()
			}
		},
	}, {
		name: "during a backoff",
		setup: func(t *testing.T, h *partHarness, cancel context.CancelFunc) {
			h.parts.beforeComplete = func(context.Context, int) error { return status.Error(codes.Unavailable, "the completion failed") }
			upload.ReplaceSleep(t, func(ctx context.Context, _ time.Duration) error {
				cancel()
				return ctx.Err()
			})
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			tc.setup(t, h, cancel)
			data := randomBytes(t, 100)
			b, _ := memBundle(data, data)

			_, err := (&upload.PartClient{Uploads: h.client}).Upload(ctx, org(), b)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Upload: got = %v, want context.Canceled", err)
			}
			if _, completed := h.parts.calls(); completed != 1 {
				t.Errorf("CompleteUpload calls: got = %d, want = 1", completed)
			}
		})
	}
}

// TestPartUploadRefusesBeforeAnyRPC covers what Upload checks of its own
// inputs, before it asks the server for anything. checkBundle's other cases
// are in TestUploadRefusesBeforeAnyRPC.
func TestPartUploadRefusesBeforeAnyRPC(t *testing.T) {
	data := []byte("a\n")
	for _, tc := range []struct {
		name     string
		noAPI    bool
		bundle   func() *upload.Bundle
		partSize int64
		wantErr  string
	}{{
		name:    "no Uploads client",
		noAPI:   true,
		bundle:  func() *upload.Bundle { b, _ := memBundle(data, data); return b },
		wantErr: "no Uploads client",
	}, {
		name:    "nil bundle",
		bundle:  func() *upload.Bundle { return nil },
		wantErr: "nil bundle",
	}, {
		name:    "empty bundle",
		bundle:  func() *upload.Bundle { b, _ := memBundle(nil, nil); return b },
		wantErr: "must not be empty",
	}, {
		name:     "negative PartSize",
		bundle:   func() *upload.Bundle { b, _ := memBundle(data, data); return b },
		partSize: -alignment,
		wantErr:  "PartSize -262144 is not a positive multiple of 262144",
	}, {
		name:     "PartSize below 256 KiB",
		bundle:   func() *upload.Bundle { b, _ := memBundle(data, data); return b },
		partSize: alignment / 2,
		wantErr:  "PartSize 131072 is not a positive multiple of 262144",
	}, {
		name:     "PartSize not a multiple of 256 KiB",
		bundle:   func() *upload.Bundle { b, _ := memBundle(data, data); return b },
		partSize: alignment + 1,
		wantErr:  "PartSize 262145 is not a positive multiple of 262144",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newPartHarness(t)
			c := &upload.PartClient{Uploads: h.client, PartSize: tc.partSize}
			if tc.noAPI {
				c.Uploads = nil
			}

			_, err := c.Upload(t.Context(), org(), tc.bundle())
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Upload: got = %v, want an error containing %q", err, tc.wantErr)
			}
			h.parts.mu.Lock()
			defer h.parts.mu.Unlock()
			if h.parts.keys != 0 || len(h.parts.records) != 0 {
				t.Errorf("GetOrgPublicKey and CreateUpload calls: got = %d and %d, want none", h.parts.keys, len(h.parts.records))
			}
		})
	}
}
