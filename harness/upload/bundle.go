/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Bundle is a plaintext harness bundle to upload.
type Bundle struct {
	// Size is the bundle's length in bytes.
	Size int64
	// SHA256 is the bundle's digest, "sha256:" and 64 lowercase hex digits.
	SHA256 string
	// Open returns a new reader over the bundle from its first byte. Upload
	// calls it once and closes what it returns. Upload returns only once
	// that reader does, so cancelling ctx does not interrupt a Read that
	// blocks.
	Open func() (io.ReadCloser, error)
}

// DescribeFile reads the file at path once and returns it as a Bundle
// whose Open reopens path.
func DescribeFile(path string) (*Bundle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening the bundle: %w", err)
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("reading the bundle: %w", err)
	}
	return &Bundle{
		Size:   n,
		SHA256: "sha256:" + hex.EncodeToString(h.Sum(nil)),
		Open: func() (io.ReadCloser, error) {
			// Returning os.Open's results directly would turn a nil
			// *os.File into a non-nil io.ReadCloser.
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			return f, nil
		},
	}, nil
}
