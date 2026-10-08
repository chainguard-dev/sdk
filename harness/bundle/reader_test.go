/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"archive/tar"
	"bytes"
	"cmp"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/klauspost/compress/zstd"

	"chainguard.dev/sdk/harness/bundle"
)

// zeroHash is a syntactically valid tree_sha256 that matches no tree, for rows
// rejected before the tree hash is reached.
var zeroHash = strings.Repeat("0", 64)

// TestValidateRejects is the hostile-tar table. Every row is a bundle a
// well-behaved writer would never produce, and the server refuses it whatever
// the client checked.
//
// One rule the writer enforces is not reachable from here: a NUL terminates a
// name in a ustar header and archive/tar refuses to encode the PAX record that
// could carry one, so no row can hand the reader a NUL-bearing path. The rule
// is covered writer-side in TestWriteRejects.
func TestValidateRejects(t *testing.T) {
	for _, tc := range []struct {
		name   string
		limits bundle.Limits
		// m is the manifest to write; nil writes a valid one.
		m *bundle.Manifest
		// write adds the entries after the manifest.
		write func(testing.TB, *tar.Writer)
		// build replaces the whole bundle, for input that is not one.
		build func(testing.TB) []byte
		want  bundle.Reason
	}{{
		name:  "path traversal",
		write: entries(file("../escape.txt", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "traversal mid-path",
		write: entries(file("a/../../b.txt", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "absolute path",
		write: entries(file("/etc/shadow", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "dot segment",
		write: entries(file("a/./b.txt", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "empty segment",
		write: entries(file("a//b.txt", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "invalid UTF-8 path",
		write: entries(file("a/\xff\xfe.txt", 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		// Legal POSIX bytes the reader still refuses: a Windows extractor
		// splits on them and lands outside the tree.
		name:  "backslash in a path",
		write: entries(file(`a\..\..\b.txt`, 0o644, "x")),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "backslash in a link target",
		write: entries(link("a/evil", tar.TypeSymlink, `..\..\etc\shadow`)),
		want:  bundle.ReasonSymlinkEscape,
	}, {
		name:  "duplicate path",
		write: entries(file("a.txt", 0o644, "one"), file("a.txt", 0o644, "two")),
		want:  bundle.ReasonDuplicatePath,
	}, {
		name:  "unsorted entries",
		write: entries(file("b.txt", 0o644, "two"), file("a.txt", 0o644, "one")),
		want:  bundle.ReasonUnsortedEntry,
	}, {
		name:  "hardlink",
		write: entries(link("hard.txt", tar.TypeLink, "a.txt")),
		want:  bundle.ReasonEntryType,
	}, {
		name:  "character device",
		write: entries(device("dev/tty", tar.TypeChar)),
		want:  bundle.ReasonEntryType,
	}, {
		name:  "block device",
		write: entries(device("dev/sda", tar.TypeBlock)),
		want:  bundle.ReasonEntryType,
	}, {
		name:  "FIFO",
		write: entries(special("var/pipe", tar.TypeFifo, 0o644)),
		want:  bundle.ReasonEntryType,
	}, {
		name:  "symlink escaping the tree",
		write: entries(link("a/evil", tar.TypeSymlink, "../../etc/shadow")),
		want:  bundle.ReasonSymlinkEscape,
	}, {
		name:  "absolute symlink target",
		write: entries(link("evil", tar.TypeSymlink, "/etc/shadow")),
		want:  bundle.ReasonSymlinkEscape,
	}, {
		name:  "reserved manifest path in the source tree",
		write: entries(file(bundle.ManifestPath, 0o644, "{}")),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "directory without a trailing slash",
		write: entries(special("src", tar.TypeDir, 0o755)),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "world-writable file mode",
		write: entries(file("a.txt", 0o666, "x")),
		want:  bundle.ReasonMetadata,
	}, {
		name:  "setuid file mode",
		write: entries(file("a.txt", 0o4755, "x")),
		want:  bundle.ReasonMetadata,
	}, {
		name:  "directory with a file's mode",
		write: entries(special("src/", tar.TypeDir, 0o644)),
		want:  bundle.ReasonMetadata,
	}, {
		name:  "non-epoch mtime",
		write: entries(tweak(file("a.txt", 0o644, "x"), func(h *tar.Header) { h.ModTime = time.Date(2026, time.October, 5, 12, 0, 0, 0, time.UTC) })),
		want:  bundle.ReasonMetadata,
	}, {
		name:  "non-zero uid",
		write: entries(tweak(file("a.txt", 0o644, "x"), func(h *tar.Header) { h.Uid, h.Gid = 501, 20 })),
		want:  bundle.ReasonMetadata,
	}, {
		name:  "uname set",
		write: entries(tweak(file("a.txt", 0o644, "x"), func(h *tar.Header) { h.Uname = "mallory" })),
		want:  bundle.ReasonMetadata,
	}, {
		name: "extra PAX record",
		write: entries(tweak(file("a.txt", 0o644, "x"), func(h *tar.Header) {
			h.PAXRecords = map[string]string{"SCHILY.xattr.security.selinux": "unconfined_u"}
		})),
		want: bundle.ReasonMetadata,
	}, {
		name:  "symlink declaring content",
		write: entries(tweak(link("a/evil", tar.TypeSymlink, "b"), func(h *tar.Header) { h.Size = 16 })),
		want:  bundle.ReasonMetadata,
	}, {
		name:   "path too long",
		limits: bundle.Limits{MaxPathBytes: 24},
		write:  entries(file(strings.Repeat("a", 32)+".txt", 0o644, "x")),
		want:   bundle.ReasonPathTooLong,
	}, {
		name:   "path too deep",
		limits: bundle.Limits{MaxPathDepth: 3},
		write:  entries(file("a/b/c/d.txt", 0o644, "x")),
		want:   bundle.ReasonTooDeep,
	}, {
		name:   "file too large",
		limits: bundle.Limits{MaxFileBytes: 8},
		write:  entries(file("big.bin", 0o644, strings.Repeat("x", 64))),
		want:   bundle.ReasonFileTooLarge,
	}, {
		name:   "too many entries",
		limits: bundle.Limits{MaxEntries: 2},
		write:  entries(file("a.txt", 0o644, "1"), file("b.txt", 0o644, "2")),
		want:   bundle.ReasonTooManyEntries,
	}, {
		name:   "decoded stream too large",
		limits: bundle.Limits{MaxDecodedBytes: 4096, MaxFileBytes: 1 << 20},
		write:  entries(file("big.bin", 0o644, strings.Repeat("x", 128<<10))),
		want:   bundle.ReasonStreamTooLarge,
	}, {
		name:   "decompression ratio",
		limits: bundle.Limits{MaxRatio: 4, RatioFloorBytes: 4096, MaxFileBytes: 1 << 20, MaxDecodedBytes: 1 << 20},
		write:  entries(file("big.bin", 0o644, strings.Repeat("x", 512<<10))),
		want:   bundle.ReasonDecompressRatio,
	}, {
		name:   "compressed bundle too large",
		limits: bundle.Limits{MaxCompressedBytes: 16},
		write:  entries(file("a.txt", 0o644, "x")),
		want:   bundle.ReasonBundleTooLarge,
	}, {
		name: "no manifest",
		build: func(t testing.TB) []byte {
			return rawTar(t, entries(file("a.txt", 0o644, "x")))
		},
		want: bundle.ReasonInvalidManifest,
	}, {
		name:  "empty bundle",
		build: func(t testing.TB) []byte { return rawTar(t, nil) },
		want:  bundle.ReasonInvalidManifest,
	}, {
		name: "manifest is not the first entry",
		build: func(t testing.TB) []byte {
			return rawTar(t, func(t testing.TB, tw *tar.Writer) {
				entries(file("a.txt", 0o644, "x"))(t, tw)
				writeManifest(t, tw, nil)
			})
		},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "manifest with an unknown field",
		build: func(t testing.TB) []byte {
			return rawTar(t, entries(raw(manifestHeader(), `{"format":"harness-source-v1","kind":"folder","dirty":false,"files":0,"bytes":0,"tree_sha256":"`+zeroHash+`","client":"chainctl v0.1.2","repo_url":"https://github.com/acme/widgets"}`)))
		},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "manifest that is not JSON",
		build: func(t testing.TB) []byte {
			return rawTar(t, entries(raw(manifestHeader(), "harness-source-v1\n")))
		},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "manifest with an unknown kind",
		m:    &bundle.Manifest{Kind: "tarball", Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "manifest with a bad commit",
		m:    &bundle.Manifest{Kind: bundle.KindGit, Commit: "HEAD", Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name:  "tree hash does not match the entries",
		write: entries(file("a.txt", 0o644, "x")),
		want:  bundle.ReasonTreeHashMismatch,
	}, {
		name: "file count does not match the entries",
		m: &bundle.Manifest{
			Kind: bundle.KindFolder, Client: "chainctl v0.1.2", Files: 7, Bytes: 1,
			TreeSHA256: specTreeHash([4]string{"a.txt", "0644", "file", sha256hex([]byte("x"))}),
		},
		write: entries(file("a.txt", 0o644, "x")),
		want:  bundle.ReasonInvalidManifest,
	}, {
		name: "byte count does not match the entries",
		m: &bundle.Manifest{
			Kind: bundle.KindFolder, Client: "chainctl v0.1.2", Files: 1, Bytes: 4096,
			TreeSHA256: specTreeHash([4]string{"a.txt", "0644", "file", sha256hex([]byte("x"))}),
		},
		write: entries(file("a.txt", 0o644, "x")),
		want:  bundle.ReasonInvalidManifest,
	}, {
		name:  "not a bundle",
		build: func(testing.TB) []byte { return []byte("this is not a zstd stream, it is a letter to the agent") },
		want:  bundle.ReasonMalformed,
	}, {
		name: "truncated bundle",
		build: func(t testing.TB) []byte {
			whole := rawBundle(t, nil, entries(file("a.txt", 0o644, "x")))
			return whole[:len(whole)/2]
		},
		want: bundle.ReasonMalformed,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			data := rawBundle(t, tc.m, tc.write)
			if tc.build != nil {
				data = tc.build(t)
			}
			r := &bundle.Reader{Limits: tc.limits}
			if _, err := r.Read(t.Context(), bytes.NewReader(data), nil); !isReason(err, tc.want) {
				t.Errorf("Read: got = %v, want reason %q", err, tc.want)
			}
		})
	}
}

// TestValidateRefusesBombBeforeFullDecode pins the guard that matters most: a
// decompression bomb is refused while it inflates, not once it has. This bundle
// decodes to 32 MiB of zeros followed by a megabyte of incompressible bytes, so
// nearly all of its compressed size sits at the far end — a validator that
// checked sizes after decoding would have read all of it.
func TestValidateRefusesBombBeforeFullDecode(t *testing.T) {
	body := make([]byte, 32<<20+(1<<20))
	rand.Read(body[32<<20:])
	bomb := rawBundle(t, nil, func(t testing.TB, tw *tar.Writer) {
		h := tarHeader("bomb.bin", tar.TypeReg, 0o644)
		h.Size = int64(len(body))
		if err := tw.WriteHeader(h); err != nil {
			t.Fatalf("writing header: %v", err)
		}
		if _, err := tw.Write(body); err != nil {
			t.Fatalf("writing body: %v", err)
		}
	})

	src := &countingSource{r: bytes.NewReader(bomb)}
	r := &bundle.Reader{Limits: bundle.Limits{MaxDecodedBytes: 1 << 20, MaxFileBytes: 64 << 20}}
	if _, err := r.Read(t.Context(), src, nil); !isReason(err, bundle.ReasonStreamTooLarge) {
		t.Fatalf("Read: got = %v, want reason %q", err, bundle.ReasonStreamTooLarge)
	}
	if src.n > int64(len(bomb))/2 {
		t.Errorf("read %d of %d compressed bytes before refusing the bomb; the guard ran after the decode, not during it", src.n, len(bomb))
	}
}

// TestValidateRatioSurvivesAnIncompressiblePrefix pins what the ratio guard
// measures. It is the whole stream so far, not the chunk in hand, so a bundle
// that opens with incompressible bytes — holding the running ratio at roughly
// 1:1 well past RatioFloorBytes — does not buy the inflation that follows a
// free pass.
func TestValidateRatioSurvivesAnIncompressiblePrefix(t *testing.T) {
	prefix := make([]byte, 256<<10)
	rand.Read(prefix)
	r := &bundle.Reader{Limits: bundle.Limits{
		MaxRatio:        4,
		RatioFloorBytes: 4096,
		MaxFileBytes:    16 << 20,
		MaxDecodedBytes: 16 << 20,
	}}

	// The prefix alone is well past the floor and well under the ratio, so it
	// reaches the end of the stream and is refused for its tree hash. Without
	// this the row below would prove nothing about where the guard fired.
	calm := rawBundle(t, nil, entries(file("a-incompressible.bin", 0o644, string(prefix))))
	if _, err := r.Read(t.Context(), bytes.NewReader(calm), nil); !isReason(err, bundle.ReasonTreeHashMismatch) {
		t.Fatalf("Read of the prefix alone: got = %v, want reason %q", err, bundle.ReasonTreeHashMismatch)
	}

	bomb := rawBundle(t, nil, entries(
		file("a-incompressible.bin", 0o644, string(prefix)),
		file("b-zeros.bin", 0o644, strings.Repeat("\x00", 8<<20)),
	))
	if _, err := r.Read(t.Context(), bytes.NewReader(bomb), nil); !isReason(err, bundle.ReasonDecompressRatio) {
		t.Errorf("Read: got = %v, want reason %q", err, bundle.ReasonDecompressRatio)
	}
}

// TestReservedKindIsRejectedAsDeferred pins what a reader of a bundle from a
// client that shipped a repository source ahead of the format is told: the kind
// is spoken for and has no packer yet, which is a different thing from a name
// nobody recognizes.
func TestReservedKindIsRejectedAsDeferred(t *testing.T) {
	m := &bundle.Manifest{Kind: bundle.KindSource, Commit: strings.Repeat("a", 40), Client: "chainctl v0.1.2"}
	_, err := bundle.Validate(t.Context(), bytes.NewReader(rawBundle(t, m, nil)))
	if !isReason(err, bundle.ReasonInvalidManifest) {
		t.Fatalf("Validate: got = %v, want reason %q", err, bundle.ReasonInvalidManifest)
	}
	if !strings.Contains(err.Error(), "deferred") {
		t.Errorf("rejection does not say the kind is deferred: %v", err)
	}
}

// TestReadAcceptsZeroPadding covers the allowance that keeps a GNU-style
// blocked writer conforming: tar pads its output to a whole record, so a reader
// that refused everything after the end-of-archive marker would reject a
// faithful reimplementation of this format.
func TestReadAcceptsZeroPadding(t *testing.T) {
	// A whole 20-block GNU record past the marker, and the allowance itself.
	for _, pad := range []int{0, 512, 10240 - 1024, 32 << 10} {
		t.Run(strconv.Itoa(pad), func(t *testing.T) {
			data := repad(t, padding(pad))
			m, err := bundle.Validate(t.Context(), bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if m.TreeSHA256 != goldenTreeSHA256 {
				t.Errorf("tree_sha256: got = %s, want = %s", m.TreeSHA256, goldenTreeSHA256)
			}
		})
	}
}

// TestReadRejectsTrailingBytes closes the smuggling channel the padding
// allowance leaves open. Bytes after the marker are read but no entry declares
// them and tree_sha256 does not cover them, so the only ones permitted are
// zeros, and only so many.
func TestReadRejectsTrailingBytes(t *testing.T) {
	for _, tc := range []struct {
		name string
		pad  []byte
	}{{
		name: "payload straight after the marker",
		pad:  []byte("exfiltrated: the second stage goes here"),
	}, {
		name: "payload hidden behind conforming padding",
		pad:  append(padding(10240-1024), "a note for the agent"...),
	}, {
		name: "one non-zero byte at the end of the padding",
		pad:  append(padding(4096), 1),
	}, {
		name: "padding past the allowance",
		pad:  padding(32<<10 + 512),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			data := repad(t, tc.pad)
			if _, err := bundle.Validate(t.Context(), bytes.NewReader(data)); !isReason(err, bundle.ReasonTrailingBytes) {
				t.Errorf("Validate: got = %v, want reason %q", err, bundle.ReasonTrailingBytes)
			}
		})
	}
}

// TestTrailingBytesCountAgainstTheStreamLimits pins that padding is not a way
// past a limit. It is read through the same guards as the entries, so a bundle
// cannot hold its tar under MaxDecodedBytes and then inflate past it, or past
// the ratio, on zeros the trailing allowance alone would have permitted.
func TestTrailingBytesCountAgainstTheStreamLimits(t *testing.T) {
	// The golden tar is 5632 bytes, so these limits are met only by what
	// follows it.
	t.Run("decoded stream", func(t *testing.T) {
		data := repad(t, padding(16<<10))
		r := &bundle.Reader{Limits: bundle.Limits{MaxDecodedBytes: 8 << 10, MaxTrailingBytes: 1 << 20}}
		if _, err := r.Read(t.Context(), bytes.NewReader(data), nil); !isReason(err, bundle.ReasonStreamTooLarge) {
			t.Errorf("Read: got = %v, want reason %q", err, bundle.ReasonStreamTooLarge)
		}
	})
	t.Run("ratio", func(t *testing.T) {
		data := repad(t, padding(1<<20))
		r := &bundle.Reader{Limits: bundle.Limits{
			MaxRatio:         4,
			RatioFloorBytes:  4096,
			MaxDecodedBytes:  16 << 20,
			MaxTrailingBytes: 8 << 20,
		}}
		if _, err := r.Read(t.Context(), bytes.NewReader(data), nil); !isReason(err, bundle.ReasonDecompressRatio) {
			t.Errorf("Read: got = %v, want reason %q", err, bundle.ReasonDecompressRatio)
		}
	})
}

func padding(n int) []byte { return make([]byte, n) }

// repad repacks the golden bundle with pad after the tar's end-of-archive
// marker, inside the same zstd frame — where a blocked writer's padding, and
// anything hiding in it, arrives.
func repad(t testing.TB, pad []byte) []byte {
	t.Helper()
	data, _ := pack(t, goldenManifest(), goldenEntries())
	var buf bytes.Buffer
	zw, err := zstd.NewWriter(&buf)
	if err != nil {
		t.Fatalf("zstd writer: %v", err)
	}
	if _, err := zw.Write(append(untar(t, data), pad...)); err != nil {
		t.Fatalf("writing padded tar: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("closing zstd: %v", err)
	}
	return buf.Bytes()
}

func TestReadStopsOnCallbackError(t *testing.T) {
	data, _ := pack(t, goldenManifest(), goldenEntries())
	sentinel := errors.New("no room on disk")
	var seen int
	_, err := (&bundle.Reader{}).Read(t.Context(), bytes.NewReader(data), func(bundle.Entry, io.Reader) error {
		seen++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("Read: got = %v, want = %v", err, sentinel)
	}
	if seen != 1 {
		t.Errorf("entries visited: got = %d, want = 1", seen)
	}
}

func TestReadHonoursContext(t *testing.T) {
	data, _ := pack(t, goldenManifest(), goldenEntries())
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := bundle.Validate(ctx, bytes.NewReader(data)); !errors.Is(err, context.Canceled) {
		t.Errorf("Validate: got = %v, want = %v", err, context.Canceled)
	}
}

func TestManifestTrailingContent(t *testing.T) {
	m := goldenManifest()
	m.Format, m.TreeSHA256 = bundle.Format, specTreeHash()
	doc, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		trailer string
		wantErr bool
	}{
		{name: "EOF"},
		{name: "whitespace", trailer: " \t\r\n"},
		{name: "closing object", trailer: "}", wantErr: true},
		{name: "closing array", trailer: "]", wantErr: true},
		{name: "hidden payload", trailer: "}\ntrailing bytes", wantErr: true},
		{name: "second object", trailer: "{}", wantErr: true},
		{name: "second scalar", trailer: "null", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := rawTar(t, entries(raw(manifestHeader(), string(doc)+tc.trailer)))
			got, err := bundle.Validate(t.Context(), bytes.NewReader(data))
			if tc.wantErr {
				if !isReason(err, bundle.ReasonInvalidManifest) || got != nil {
					t.Errorf("Validate: got = %v, %v, want = nil, reason %q", got, err, bundle.ReasonInvalidManifest)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if diff := gocmp.Diff(&m, got); diff != "" {
				t.Errorf("manifest (-want +got):\n%s", diff)
			}
		})
	}
}

func FuzzValidate(f *testing.F) {
	valid, _ := pack(f, goldenManifest(), goldenEntries())
	f.Add(valid)
	f.Add(rawBundle(f, nil, entries(file("a.txt", 0o644, "x"))))
	f.Add(rawTar(f, nil))
	f.Add(repad(f, padding(512)))
	f.Add(repad(f, []byte("bytes riding past the end-of-archive marker")))
	f.Add([]byte("not a bundle"))
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		// Small limits keep a case cheap; they are the same guards.
		r := &bundle.Reader{Limits: bundle.Limits{
			MaxCompressedBytes: 1 << 20,
			MaxDecodedBytes:    8 << 20,
			MaxFileBytes:       1 << 20,
			MaxEntries:         64,
			MaxPathBytes:       256,
			MaxPathDepth:       8,
			MaxRatio:           100,
			RatioFloorBytes:    64 << 10,
		}}
		m, err := r.Read(t.Context(), bytes.NewReader(data), nil)
		switch {
		case err == nil && m == nil:
			t.Fatal("Read returned no manifest and no error")
		case err == nil:
			return
		}
		// Every rejection is an *Error carrying a reason the caller can
		// surface; a bare error is a path that escaped the format's vocabulary.
		e, ok := errors.AsType[*bundle.Error](err)
		if !ok {
			t.Fatalf("rejection is not a *bundle.Error: %v", err)
		}
		if e.Reason == "" {
			t.Fatalf("rejection carries no reason: %v", err)
		}
	})
}

// rawEntry is a tar entry a test writes by hand, bypassing the writer's rules.
type rawEntry struct {
	h    *tar.Header
	body string
}

func tarHeader(name string, typeflag byte, mode int64) *tar.Header {
	return &tar.Header{
		Typeflag: typeflag,
		Name:     name,
		Mode:     mode,
		ModTime:  time.Unix(0, 0),
		Format:   tar.FormatPAX,
	}
}

func manifestHeader() *tar.Header { return tarHeader(bundle.ManifestPath, tar.TypeReg, 0o644) }

func file(name string, mode int64, body string) rawEntry {
	h := tarHeader(name, tar.TypeReg, mode)
	h.Size = int64(len(body))
	return rawEntry{h: h, body: body}
}

func special(name string, typeflag byte, mode int64) rawEntry {
	return rawEntry{h: tarHeader(name, typeflag, mode)}
}

func device(name string, typeflag byte) rawEntry {
	e := special(name, typeflag, 0o644)
	e.h.Devmajor, e.h.Devminor = 5, 1
	return e
}

// dir names the entry the way the format stores it, with the trailing slash
// the reader requires; Entry.Path carries the name without it.
func dir(name string) rawEntry { return special(name+"/", tar.TypeDir, 0o755) }

func link(name string, typeflag byte, target string) rawEntry {
	e := special(name, typeflag, 0o777)
	e.h.Linkname = target
	return e
}

// raw is an entry whose header the caller built outright.
func raw(h *tar.Header, body string) rawEntry {
	h.Size = int64(len(body))
	return rawEntry{h: h, body: body}
}

func tweak(e rawEntry, f func(*tar.Header)) rawEntry {
	f(e.h)
	return e
}

// entries writes the given entries in the order they are given.
func entries(es ...rawEntry) func(testing.TB, *tar.Writer) {
	return func(t testing.TB, tw *tar.Writer) {
		t.Helper()
		for _, e := range es {
			if err := tw.WriteHeader(e.h); err != nil {
				t.Fatalf("writing header %q: %v", e.h.Name, err)
			}
			if e.body != "" {
				if _, err := io.WriteString(tw, e.body); err != nil {
					t.Fatalf("writing body of %q: %v", e.h.Name, err)
				}
			}
		}
	}
}

// rawTar builds a zstd-compressed tar with no manifest of its own, so a test
// can put in it what the writer would never emit.
func rawTar(t testing.TB, write func(testing.TB, *tar.Writer)) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw, err := zstd.NewWriter(&buf)
	if err != nil {
		t.Fatalf("zstd writer: %v", err)
	}
	tw := tar.NewWriter(zw)
	if write != nil {
		write(t, tw)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("closing zstd: %v", err)
	}
	return buf.Bytes()
}

// rawBundle is rawTar with a manifest first. A nil manifest is a valid one over
// a tree hash that matches nothing.
func rawBundle(t testing.TB, m *bundle.Manifest, write func(testing.TB, *tar.Writer)) []byte {
	t.Helper()
	return rawTar(t, func(t testing.TB, tw *tar.Writer) {
		writeManifest(t, tw, m)
		if write != nil {
			write(t, tw)
		}
	})
}

func writeManifest(t testing.TB, tw *tar.Writer, m *bundle.Manifest) {
	t.Helper()
	valid := bundle.Manifest{Kind: bundle.KindFolder, Client: "chainctl v0.1.2"}
	if m != nil {
		valid = *m
	}
	valid.Format = cmp.Or(valid.Format, bundle.Format)
	valid.TreeSHA256 = cmp.Or(valid.TreeSHA256, zeroHash)
	doc, err := json.MarshalIndent(valid, "", "  ")
	if err != nil {
		t.Fatalf("encoding manifest: %v", err)
	}
	entries(raw(manifestHeader(), string(doc)+"\n"))(t, tw)
}

// specTreeHash computes tree_sha256 straight from the definition in the package
// doc — path, mode, type, payload, each length-prefixed — independently of the
// package's own implementation.
func specTreeHash(records ...[4]string) string {
	h := sha256.New()
	field := func(s string) { fmt.Fprintf(h, "%d:%s\n", len(s), s) }
	field("harness-source-v1-tree")
	for _, r := range records {
		for _, f := range r {
			field(f)
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// countingSource records how much of a bundle a reader actually consumed.
type countingSource struct {
	r io.Reader
	n int64
}

func (c *countingSource) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}
