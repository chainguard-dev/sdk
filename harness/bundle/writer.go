/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import (
	"archive/tar"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/klauspost/compress/zstd"
)

// epoch is the mtime every entry carries: metadata is zeroed, so a bundle does
// not leak when a customer's files were touched.
var epoch = time.Unix(0, 0)

// Writer packs a bundle. The zero value is the one production callers use; its
// Limits default to [DefaultLimits].
type Writer struct {
	// Limits bounds the bundle being written. A zero field takes its default,
	// so the zero Writer enforces exactly [DefaultLimits]. Only tests set it,
	// to drive a guard at a small size.
	Limits Limits
}

// Write packs entries into a harness-source-v1 bundle on dst and returns the
// manifest it wrote.
//
// The caller supplies Kind, Commit, Dirty, IncludeIgnored and Client; Write
// computes Format, Files, Bytes and TreeSHA256 from the entries so the manifest
// cannot disagree with what was packed. Entries are sorted by path in byte
// order and normalized here, so one tree packs to one tar whatever order intake
// enumerated it in. The compressed bytes are not pinned that way: the encoder's
// version and level decide them, which is why the manifest carries tree_sha256
// and intake declares the bundle's SHA-256. See the package doc.
//
// Write rejects invalid entries and unsafe symlink chains. [Normalize] lets
// intake check individual entries; chain containment requires the whole tree.
//
// Write streams, so dst holds a partial bundle whenever Write returns an error
// — including the [ReasonSourceChanged] it raises when a file changes under the
// pack. The caller discards dst on any error; nothing it holds is a bundle.
//
// A cancelled ctx stops the pack between entries and is returned as it is, not
// as an [Error]: nothing was wrong with the bundle.
func (w *Writer) Write(ctx context.Context, dst io.Writer, m Manifest, entries []Entry) (*Manifest, error) {
	lim := w.Limits.withDefaults()
	pr, pw := io.Pipe()
	checked := make(chan error, 1)
	go func() {
		_, err := (&Reader{Limits: lim}).Read(ctx, pr, nil)
		_ = pr.CloseWithError(err)
		checked <- err
	}()
	// The ratio is a property of decoded prefixes, not just the final sizes.
	// Validate the actual compressed stream without buffering the bundle or
	// duplicating the reader's policy in the encoder.
	result, err := write(ctx, io.MultiWriter(dst, pw), m, entries, lim)
	_ = pw.CloseWithError(err)
	validationErr := <-checked
	if validationErr != nil && (err == nil || errors.Is(err, validationErr)) {
		return nil, validationErr
	}
	return result, err
}

func write(ctx context.Context, dst io.Writer, m Manifest, entries []Entry, lim Limits) (*Manifest, error) {
	m.Format = Format
	// What the caller supplied is checked before a byte of the tree is read: a
	// manifest the writer could never emit should not cost a digest pass over a
	// monorepo first.
	if err := m.checkSource(); err != nil {
		return nil, err
	}
	sorted, err := normalize(entries, lim)
	if err != nil {
		return nil, err
	}

	// First pass: digest content and build the tree hash, so the manifest the
	// second pass writes first already describes every entry behind it.
	th := newTreeHasher(sha256.New())
	payloads := make([]string, len(sorted))
	var files int
	var content, decoded int64
	for i, e := range sorted {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		switch e.Type {
		case TypeFile:
			if payloads[i], err = digestEntry(e); err != nil {
				return nil, err
			}
			files++
			content += e.Size
		case TypeSymlink:
			payloads[i] = e.LinkTarget
		case TypeDir:
			// A directory's tree-hash payload is the empty string.
		}
		th.add(e, payloads[i])
		decoded += tarBytes(e.Size)
		if decoded > lim.MaxDecodedBytes {
			return nil, rejectf(ReasonStreamTooLarge, e.Path, "tar stream exceeds %d bytes", lim.MaxDecodedBytes)
		}
	}

	m.Files = files
	m.Bytes = content
	m.TreeSHA256 = th.sum()
	if err := m.validate(); err != nil {
		return nil, err
	}
	doc, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, rejectf(ReasonInvalidManifest, ManifestPath, "encoding manifest: %w", err)
	}
	doc = append(doc, '\n')
	if int64(len(doc)) > maxManifestBytes {
		return nil, rejectf(ReasonInvalidManifest, ManifestPath, "manifest is %d bytes, limit is %d", len(doc), maxManifestBytes)
	}

	out := &countingWriter{w: dst, what: "compressed bundle", reason: ReasonBundleTooLarge, limit: lim.MaxCompressedBytes}
	// The format pins the codec, not the ratio, so the encoder keeps its default
	// level and a later zstd may compress the same tar differently. Concurrency
	// is pinned anyway, so one build's output does not also depend on how many
	// cores the machine that packed it had.
	zw, err := zstd.NewWriter(out, zstd.WithEncoderConcurrency(1))
	if err != nil {
		return nil, rejectf(ReasonMalformed, "", "creating zstd writer: %w", err)
	}
	// The first pass's running total is a lower bound — it has no manifest
	// entry, no extended header for a path too long for ustar, and no
	// end-of-archive marker — so the limit is also applied to the tar bytes
	// themselves. Without this the writer would accept trees the reader then
	// refuses, which for 250 000 long paths is a quarter of a gigabyte of
	// extended headers the estimate never sees.
	raw := &countingWriter{w: zw, what: "tar stream", reason: ReasonStreamTooLarge, limit: lim.MaxDecodedBytes}
	tw := tar.NewWriter(raw)
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     ManifestPath,
		Size:     int64(len(doc)),
		Mode:     int64(modeFile),
		ModTime:  epoch,
		Format:   tar.FormatPAX,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, wrapWrite(err, raw, out)
	}
	if _, err := tw.Write(doc); err != nil {
		return nil, wrapWrite(err, raw, out)
	}
	for i, e := range sorted {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := writeEntry(tw, e, payloads[i]); err != nil {
			return nil, wrapWrite(err, raw, out)
		}
	}
	if err := tw.Close(); err != nil {
		return nil, wrapWrite(err, raw, out)
	}
	if err := zw.Close(); err != nil {
		return nil, wrapWrite(err, raw, out)
	}
	if err := raw.err.Load(); err != nil {
		return nil, err
	}
	if err := out.err.Load(); err != nil {
		return nil, err
	}
	return &m, nil
}

// normalize returns the entries in bundle order, each checked against the
// format's rules and its mode normalized. The caller's slice is left alone.
func normalize(entries []Entry, lim Limits) ([]Entry, error) {
	// The manifest is an entry too, and counts against the limit.
	if len(entries)+1 > lim.MaxEntries {
		return nil, rejectf(ReasonTooManyEntries, "", "%d entries, limit is %d", len(entries)+1, lim.MaxEntries)
	}
	out := slices.Clone(entries)
	// Sorting on Entry.Path, not on the name the entry takes in the tar: a
	// directory is stored with a trailing slash, and "cmd/" sorts after
	// "cmd.go" while "cmd" sorts before it. Path is the key the reader
	// re-checks and the tree hash covers.
	slices.SortFunc(out, func(a, b Entry) int { return strings.Compare(a.Path, b.Path) })
	var tree entryTree
	for i := range out {
		e, err := normalizeEntry(out[i], lim)
		if err != nil {
			return nil, err
		}
		if i > 0 && out[i-1].Path == e.Path {
			return nil, rejectf(ReasonDuplicatePath, e.Path, "path appears twice")
		}
		if err := tree.add(e); err != nil {
			return nil, err
		}
		out[i] = e
	}
	if err := tree.checkLinks(); err != nil {
		return nil, err
	}
	return out, nil
}

// digestEntry reads a file entry's content to its hex SHA-256, the payload the
// tree hash covers.
func digestEntry(e Entry) (string, error) {
	rc, err := e.Open()
	if err != nil {
		return "", rejectf(ReasonMalformed, e.Path, "opening content: %w", err)
	}
	defer rc.Close()
	h := sha256.New()
	n, err := io.Copy(h, io.LimitReader(rc, e.Size+1))
	if err != nil {
		return "", rejectf(ReasonMalformed, e.Path, "reading content: %w", err)
	}
	if n != e.Size {
		return "", rejectf(ReasonSourceChanged, e.Path, "content is %d bytes, entry declares %d", n, e.Size)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// writeEntry writes one entry's header and, for a file, its content. It
// re-digests the content it writes: the digest in the manifest ahead of it came
// from a separate read, so a source that changed in between is caught here
// rather than packed into a bundle whose tree hash is already wrong.
func writeEntry(tw *tar.Writer, e Entry, payload string) error {
	hdr := &tar.Header{
		Name:    e.Path,
		Mode:    int64(e.Mode.Perm()),
		ModTime: epoch,
		Format:  tar.FormatPAX,
	}
	switch e.Type {
	case TypeFile:
		hdr.Typeflag, hdr.Size = tar.TypeReg, e.Size
	case TypeDir:
		// A directory is stored with a trailing slash, the convention every
		// reader expects; the slash is not part of the entry's path.
		hdr.Typeflag, hdr.Name = tar.TypeDir, e.Path+"/"
	case TypeSymlink:
		hdr.Typeflag, hdr.Linkname = tar.TypeSymlink, e.LinkTarget
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if e.Type != TypeFile {
		return nil
	}
	rc, err := e.Open()
	if err != nil {
		return rejectf(ReasonMalformed, e.Path, "opening content: %w", err)
	}
	defer rc.Close()
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(tw, h), io.LimitReader(rc, e.Size+1))
	if err != nil {
		return rejectf(ReasonMalformed, e.Path, "writing content: %w", err)
	}
	if n != e.Size || hex.EncodeToString(h.Sum(nil)) != payload {
		return rejectf(ReasonSourceChanged, e.Path, "content changed while packing")
	}
	return nil
}

// tarBytes is the lower bound on what one entry costs in the decoded stream: a
// 512-byte header plus its content padded to a whole block. A path too long for
// a ustar header costs an extended header on top, which only the tar writer
// knows about.
func tarBytes(size int64) int64 {
	return 512 + (size+511)/512*512
}

// countingWriter bounds one of the writer's two streams as it is produced,
// giving the client the same rejection the server would reach on the way back
// in. The limit is applied before the bytes are handed on, so the stream never
// grows past it.
type countingWriter struct {
	w      io.Writer
	what   string
	reason Reason
	limit  int64

	n   atomic.Int64
	err atomic.Pointer[Error]
}

func (c *countingWriter) Write(p []byte) (int, error) {
	if err := c.err.Load(); err != nil {
		return 0, err
	}
	if c.n.Load()+int64(len(p)) > c.limit {
		over := rejectf(c.reason, "", "%s exceeds %d bytes", c.what, c.limit)
		c.err.Store(over)
		return 0, over
	}
	n, err := c.w.Write(p)
	c.n.Add(int64(n))
	return n, err
}

// wrapWrite reports the limit the writer hit, if one stopped it, rather than
// the wrapped plumbing error the archive layers saw. The counters are given
// innermost first, so the rejection returned is the one raised nearest the tar
// writer rather than the plumbing failure it caused further out.
func wrapWrite(err error, counters ...*countingWriter) error {
	for _, c := range counters {
		if over := c.err.Load(); over != nil {
			return over
		}
	}
	if e, ok := errors.AsType[*Error](err); ok {
		return e
	}
	return rejectf(ReasonMalformed, "", "writing bundle: %w", err)
}
