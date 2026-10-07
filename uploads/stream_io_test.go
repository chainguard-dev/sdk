/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"chainguard.dev/sdk/uploads"
	"github.com/google/go-cmp/cmp"
)

func TestOpenStreamHTTPBody(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, size := range []int{0, 32, uploads.StreamChunkSize, uploads.StreamChunkSize + 32} {
		for _, truncated := range []bool{false, true} {
			t.Run(fmt.Sprintf("size=%d/truncated=%t", size, truncated), func(t *testing.T) {
				want := bytesOfLen(size)
				sealed := sealStream(t, &priv.PublicKey, want, 7919)
				contentLength := len(sealed)
				if truncated {
					contentLength++
				}
				header := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", contentLength)
				resp, err := http.ReadResponse(bufio.NewReader(io.MultiReader(strings.NewReader(header), bytes.NewReader(sealed))), nil)
				if err != nil {
					t.Fatalf("read HTTP response: %v", err)
				}
				defer resp.Body.Close()
				r, err := uploads.OpenStream(resp.Body, rsaStreamUnwrap(priv))
				if err != nil {
					t.Fatalf("open: %v", err)
				}
				got, err := io.ReadAll(r)
				if truncated {
					if !errors.Is(err, io.ErrUnexpectedEOF) {
						t.Errorf("truncated HTTP body: got err = %v, want io.ErrUnexpectedEOF", err)
					}
					// Earlier authenticated chunks may be released, but the final
					// chunk must not be released after a transport failure.
					prefix := max(0, (size-1)/uploads.StreamChunkSize) * uploads.StreamChunkSize
					if !bytes.Equal(want[:prefix], got) {
						diff := cmp.Diff(want[:prefix], got)
						t.Errorf("plaintext before transport failure (-want, +got):\n%s", diff)
					}
					if n, err := r.Read(make([]byte, 32)); n != 0 || !errors.Is(err, io.ErrUnexpectedEOF) {
						t.Errorf("read after transport failure: got (%d, %v), want (0, io.ErrUnexpectedEOF)", n, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("read: %v", err)
				}
				if !bytes.Equal(want, got) {
					diff := cmp.Diff(want, got)
					t.Errorf("plaintext (-want, +got):\n%s", diff)
				}
			})
		}
	}
}

func TestOpenStreamReaderVariants(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, size := range []int{0, 32, uploads.StreamChunkSize, uploads.StreamChunkSize + 32} {
		for _, tc := range []struct {
			name string
			wrap func(io.Reader) io.Reader
		}{
			{"half", iotest.HalfReader},
			{"data with EOF", iotest.DataErrReader},
		} {
			t.Run(fmt.Sprintf("size=%d/%s", size, tc.name), func(t *testing.T) {
				want := bytesOfLen(size)
				sealed := sealStream(t, &priv.PublicKey, want, 7919)
				r, err := uploads.OpenStream(tc.wrap(bytes.NewReader(sealed)), rsaStreamUnwrap(priv))
				if err != nil {
					t.Fatalf("open: %v", err)
				}
				got, err := io.ReadAll(r)
				if err != nil {
					t.Fatalf("read: %v", err)
				}
				if !bytes.Equal(want, got) {
					diff := cmp.Diff(want, got)
					t.Errorf("plaintext (-want, +got):\n%s", diff)
				}
			})
		}
	}
}
