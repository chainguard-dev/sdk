/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/klauspost/compress/zstd"
)

// maxDecoderWindow bounds the zstd window the decoder will allocate for a
// bundle, so a frame header alone cannot cost hundreds of megabytes before a
// single entry has been read. It is far above what the default encoder level
// produces.
const maxDecoderWindow = 64 << 20

// Reader streams a bundle. The zero value is the one production callers use;
// its Limits default to [DefaultLimits].
type Reader struct {
	// Limits bounds the stream. A zero field takes its default, so the zero
	// Reader enforces exactly [DefaultLimits]. Only tests set it, to drive a
	// guard at a small size.
	Limits Limits
}

// Validate streams the bundle in src under [DefaultLimits], enforcing every
// format rule, and returns its manifest. It is the server's entry point, run
// over every bundle regardless of what the client checked.
func Validate(ctx context.Context, src io.Reader) (*Manifest, error) {
	return (&Reader{}).Read(ctx, src, nil)
}

// Read streams the bundle in src, enforcing every format rule and limit as the
// bytes arrive, and returns its manifest.
//
// If fn is non-nil it is called for each entry after the manifest, in tar
// order, with a reader over a file's content that is valid only until Read
// moves on; unread content is drained. A non-nil error from fn stops the stream
// and is returned unchanged.
//
// Entries reach fn before the stream has been fully checked: tree_sha256,
// counts and symlink chains can only be confirmed at the end. A caller that
// materializes a bundle must confine every filesystem operation through an
// os.Root opened on its staging directory, then commit only once Read returns
// nil. Staging alone does not contain a link that escapes before validation.
//
// A cancelled ctx stops the stream between entries and is returned as it is,
// not as an [Error]: nothing was wrong with the bundle.
func (r *Reader) Read(ctx context.Context, src io.Reader, fn func(Entry, io.Reader) error) (*Manifest, error) {
	lim := r.Limits.withDefaults()
	in := &countingReader{r: src, limit: lim.MaxCompressedBytes}
	// The guard sits between the size count and the decoder, so a frame the
	// decoder would silently skip is refused instead (zstdframe.go).
	zr, err := zstd.NewReader(&zstdFrameGuard{r: in},
		zstd.WithDecoderConcurrency(1),
		zstd.WithDecoderLowmem(true),
		zstd.WithDecoderMaxWindow(maxDecoderWindow))
	if err != nil {
		return nil, rejectf(ReasonMalformed, "", "creating zstd reader: %w", err)
	}
	defer zr.Close()
	// Every byte the tar reader pulls — headers, skipped payloads and content
	// alike — passes through dec, so the guards bite while the stream is still
	// being decoded rather than once it has been.
	dec := &decodedReader{r: zr, in: in, lim: lim}
	tr := tar.NewReader(dec)

	var (
		m        *Manifest
		th       = newTreeHasher(sha256.New())
		previous string
		entries  int
		files    int
		content  int64
		tree     entryTree
	)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, dec.wrap(err)
		}
		entries++
		if entries > lim.MaxEntries {
			return nil, rejectf(ReasonTooManyEntries, "", "more than %d entries", lim.MaxEntries)
		}
		if entries == 1 {
			if m, err = readManifest(tr, hdr, lim); err != nil {
				return nil, err
			}
			continue
		}
		e, err := entryFromHeader(hdr, lim)
		if err != nil {
			return nil, err
		}
		switch {
		case e.Path == previous:
			return nil, rejectf(ReasonDuplicatePath, e.Path, "path appears twice")
		case e.Path < previous:
			return nil, rejectf(ReasonUnsortedEntry, e.Path, "entry sorts before %q", previous)
		}
		previous = e.Path
		if err := tree.add(e); err != nil {
			return nil, err
		}
		payload, err := consume(tr, e, fn)
		if err != nil {
			return nil, err
		}
		th.add(e, payload)
		if e.Type == TypeFile {
			files++
			content += e.Size
		}
	}
	if err := checkTrailer(dec, lim); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, rejectf(ReasonInvalidManifest, "", "bundle is empty; %q must be its first entry", ManifestPath)
	}
	if err := tree.checkLinks(); err != nil {
		return nil, err
	}
	if got := th.sum(); got != m.TreeSHA256 {
		return nil, rejectf(ReasonTreeHashMismatch, "", "entries hash to %s, manifest declares %s", got, m.TreeSHA256)
	}
	if files != m.Files || content != m.Bytes {
		return nil, rejectf(ReasonInvalidManifest, "", "entries are %d files of %d bytes, manifest declares %d of %d", files, content, m.Files, m.Bytes)
	}
	return m, nil
}

// checkTrailer reads what follows the tar's end-of-archive marker. The tree
// hash covers entries, so anything past the marker would otherwise ride along
// unread and unhashed; refusing it outright would reject a writer that blocks
// its output, which GNU tar does by default. Zero bytes within the allowance
// are what is left: a channel that carries only its own length.
//
// The bytes are read through dec, so they are counted against the decoded
// stream and the ratio like any other.
func checkTrailer(dec *decodedReader, lim Limits) error {
	var buf [4096]byte
	var n int64
	for {
		read, err := dec.Read(buf[:])
		for _, b := range buf[:read] {
			if b != 0 {
				return rejectf(ReasonTrailingBytes, "", "non-zero byte %#02x after the end-of-archive marker", b)
			}
		}
		if n += int64(read); n > lim.MaxTrailingBytes {
			return rejectf(ReasonTrailingBytes, "", "more than %d bytes of padding after the end-of-archive marker", lim.MaxTrailingBytes)
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return dec.wrap(err)
		}
	}
}

// readManifest reads the reserved first entry.
func readManifest(tr *tar.Reader, hdr *tar.Header, lim Limits) (*Manifest, error) {
	if hdr.Name != ManifestPath {
		return nil, rejectf(ReasonInvalidManifest, hdr.Name, "first entry must be %q", ManifestPath)
	}
	if hdr.Typeflag != tar.TypeReg {
		return nil, rejectf(ReasonInvalidManifest, hdr.Name, "manifest is not a regular file")
	}
	if hdr.Mode != int64(modeFile) {
		return nil, rejectf(ReasonMetadata, hdr.Name, "manifest mode is %04o, want 0644", hdr.Mode)
	}
	// entryFromHeader refuses a link target on any non-symlink entry, but the
	// manifest never passes through it. Without this a manifest header carries
	// a ustar linkname, or up to a megabyte through a PAX linkpath record,
	// that no entry declares and the tree hash does not cover.
	if hdr.Linkname != "" {
		return nil, rejectf(ReasonMetadata, hdr.Name, "manifest carries a link target")
	}
	if err := checkHeader(hdr, lim); err != nil {
		return nil, err
	}
	if hdr.Size > maxManifestBytes {
		return nil, rejectf(ReasonInvalidManifest, hdr.Name, "manifest is %d bytes, limit is %d", hdr.Size, maxManifestBytes)
	}
	doc, err := io.ReadAll(io.LimitReader(tr, maxManifestBytes+1))
	if err != nil {
		return nil, wrapRead(err, hdr.Name, "reading manifest")
	}
	var m Manifest
	d := json.NewDecoder(bytes.NewReader(doc))
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
		return nil, rejectf(ReasonInvalidManifest, hdr.Name, "decoding manifest: %w", err)
	}
	if err := d.Decode(new(json.RawMessage)); err != io.EOF {
		return nil, rejectf(ReasonInvalidManifest, hdr.Name, "manifest has trailing content")
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// entryFromHeader turns a tar header into an entry, rejecting anything the
// format does not permit.
func entryFromHeader(hdr *tar.Header, lim Limits) (Entry, error) {
	e := Entry{Path: hdr.Name}
	// The mode is matched against the normalized values rather than converted
	// from the header, so only a mode the format permits ever reaches an entry.
	switch hdr.Typeflag {
	case tar.TypeReg:
		e.Type, e.Size = TypeFile, hdr.Size
		switch hdr.Mode {
		case int64(modeFile):
			e.Mode = modeFile
		case int64(modeExec):
			e.Mode = modeExec
		default:
			return Entry{}, rejectf(ReasonMetadata, hdr.Name, "file mode is %04o, want 0644 or 0755", hdr.Mode)
		}
	case tar.TypeDir:
		name, ok := strings.CutSuffix(hdr.Name, "/")
		if !ok {
			return Entry{}, rejectf(ReasonUnsafePath, hdr.Name, "directory name has no trailing slash")
		}
		if hdr.Mode != int64(modeDir) {
			return Entry{}, rejectf(ReasonMetadata, hdr.Name, "directory mode is %04o, want 0755", hdr.Mode)
		}
		e.Type, e.Mode, e.Path = TypeDir, modeDir, name
	case tar.TypeSymlink:
		if hdr.Mode != int64(modeLink) {
			return Entry{}, rejectf(ReasonMetadata, hdr.Name, "symlink mode is %04o, want 0777", hdr.Mode)
		}
		e.Type, e.Mode, e.LinkTarget = TypeSymlink, modeLink, hdr.Linkname
	case tar.TypeLink:
		return Entry{}, rejectf(ReasonEntryType, hdr.Name, "hardlinks are not permitted")
	case tar.TypeChar, tar.TypeBlock:
		return Entry{}, rejectf(ReasonEntryType, hdr.Name, "device files are not permitted")
	case tar.TypeFifo:
		return Entry{}, rejectf(ReasonEntryType, hdr.Name, "FIFOs are not permitted")
	default:
		return Entry{}, rejectf(ReasonEntryType, hdr.Name, "entry type %q is not permitted", string(hdr.Typeflag))
	}
	if err := checkPath(e.Path, lim); err != nil {
		return Entry{}, err
	}
	if e.Path == ManifestPath {
		return Entry{}, rejectf(ReasonReservedPath, e.Path, "source tree contains the reserved manifest path")
	}
	switch e.Type {
	case TypeFile:
		if e.Size < 0 || e.Size > lim.MaxFileBytes {
			return Entry{}, rejectf(ReasonFileTooLarge, e.Path, "file is %d bytes, limit is %d", e.Size, lim.MaxFileBytes)
		}
	case TypeSymlink:
		if err := checkLinkTarget(e.Path, e.LinkTarget, lim); err != nil {
			return Entry{}, err
		}
	}
	if e.Type != TypeSymlink && hdr.Linkname != "" {
		return Entry{}, rejectf(ReasonMetadata, e.Path, "non-symlink entry has link target %q", hdr.Linkname)
	}
	if e.Type != TypeFile && hdr.Size != 0 {
		return Entry{}, rejectf(ReasonMetadata, e.Path, "%s entry declares %d content bytes", e.Type, hdr.Size)
	}
	if err := checkHeader(hdr, lim); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// checkHeader enforces the metadata every entry zeroes.
func checkHeader(hdr *tar.Header, lim Limits) error {
	// A v7 or GNU header can carry shapes this format does not describe —
	// GNU sparse entries above all — so only the two formats the writer emits
	// are read.
	if hdr.Format&(tar.FormatUSTAR|tar.FormatPAX) == 0 {
		return rejectf(ReasonMalformed, hdr.Name, "unsupported tar format %v", hdr.Format)
	}
	for k, v := range hdr.PAXRecords {
		if k != "path" && k != "linkpath" {
			return rejectf(ReasonMetadata, hdr.Name, "unexpected PAX record %q=%q", k, v)
		}
	}
	if !hdr.ModTime.Equal(epoch) {
		return rejectf(ReasonMetadata, hdr.Name, "mtime is %s, want the Unix epoch", hdr.ModTime.UTC().Format(time.RFC3339))
	}
	if !hdr.AccessTime.IsZero() || !hdr.ChangeTime.IsZero() {
		return rejectf(ReasonMetadata, hdr.Name, "entry carries an access or change time")
	}
	if hdr.Uid != 0 || hdr.Gid != 0 {
		return rejectf(ReasonMetadata, hdr.Name, "uid/gid are %d/%d, want 0/0", hdr.Uid, hdr.Gid)
	}
	if hdr.Uname != "" || hdr.Gname != "" {
		return rejectf(ReasonMetadata, hdr.Name, "uname/gname are %q/%q, want empty", hdr.Uname, hdr.Gname)
	}
	if hdr.Devmajor != 0 || hdr.Devminor != 0 {
		return rejectf(ReasonMetadata, hdr.Name, "entry carries device numbers")
	}
	name := hdr.Name
	if hdr.Typeflag == tar.TypeDir {
		name = strings.TrimSuffix(name, "/")
	}
	if len(name) > lim.MaxPathBytes {
		return rejectf(ReasonPathTooLong, name, "path is %d bytes, limit is %d", len(name), lim.MaxPathBytes)
	}
	return nil
}

// consume hands a file's content to fn, drains whatever fn left, and returns
// the entry's tree-hash payload.
func consume(tr *tar.Reader, e Entry, fn func(Entry, io.Reader) error) (string, error) {
	if e.Type != TypeFile {
		if fn != nil {
			if err := fn(e, bytes.NewReader(nil)); err != nil {
				return "", err
			}
		}
		if e.Type == TypeSymlink {
			return e.LinkTarget, nil
		}
		return "", nil
	}
	h := sha256.New()
	if fn != nil {
		if err := fn(e, io.TeeReader(tr, h)); err != nil {
			return "", err
		}
	}
	if _, err := io.Copy(h, tr); err != nil {
		return "", wrapRead(err, e.Path, "reading content")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// wrapRead keeps the rejection a limit already raised: the stream guards fire
// inside a read the tar layer is driving, and that rejection is the one the
// caller has to see, not "malformed archive".
func wrapRead(err error, p, doing string) error {
	if over, ok := errors.AsType[*Error](err); ok {
		return over
	}
	return rejectf(ReasonMalformed, p, doing+": %w", err)
}

// countingReader bounds the compressed bundle as it is consumed.
type countingReader struct {
	r     io.Reader
	limit int64

	n   atomic.Int64
	err atomic.Pointer[Error]
}

func (c *countingReader) Read(p []byte) (int, error) {
	if err := c.err.Load(); err != nil {
		return 0, err
	}
	n, err := c.r.Read(p)
	if total := c.n.Add(int64(n)); total > c.limit {
		over := rejectf(ReasonBundleTooLarge, "", "compressed bundle exceeds %d bytes", c.limit)
		c.err.Store(over)
		return 0, over
	}
	return n, err
}

// decodedReader bounds the decoded tar stream and its ratio against the
// compressed bytes read so far. Both guards are applied per read, so a bomb is
// refused while it is inflating rather than once it has inflated.
type decodedReader struct {
	r   io.Reader
	in  *countingReader
	lim Limits

	n   atomic.Int64
	err atomic.Pointer[Error]
}

func (d *decodedReader) Read(p []byte) (int, error) {
	if err := d.err.Load(); err != nil {
		return 0, err
	}
	n, err := d.r.Read(p)
	out := d.n.Add(int64(n))
	if out > d.lim.MaxDecodedBytes {
		return 0, d.fail(rejectf(ReasonStreamTooLarge, "", "decoded stream exceeds %d bytes", d.lim.MaxDecodedBytes))
	}
	if in := d.in.n.Load(); out > d.lim.RatioFloorBytes && out > in*d.lim.MaxRatio {
		return 0, d.fail(rejectf(ReasonDecompressRatio, "", "%d decoded bytes from %d compressed exceeds %d:1", out, in, d.lim.MaxRatio))
	}
	return n, err
}

// fail records a rejection so it survives whatever the decoder wraps it in.
func (d *decodedReader) fail(err *Error) error {
	d.err.Store(err)
	return err
}

// wrap reports the limit the stream hit, if one stopped it, rather than the
// decoder or tar error that surfaced in its place.
func (d *decodedReader) wrap(err error) error {
	if over := d.err.Load(); over != nil {
		return over
	}
	if over := d.in.err.Load(); over != nil {
		return over
	}
	if e, ok := errors.AsType[*Error](err); ok {
		return e
	}
	return rejectf(ReasonMalformed, "", "reading bundle: %w", err)
}
