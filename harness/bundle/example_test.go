/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"chainguard.dev/sdk/harness/bundle"
)

// content is one source file an entry reads from. Intake opens a file on disk
// or a git blob; the writer only needs it twice, once to digest and once to
// write.
func content(body string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(body)), nil }
}

func ExampleWriter_Write() {
	ctx := context.Background()
	src := map[string]string{
		"main.go":   "package main\n\nfunc main() {}\n",
		"go.mod":    "module example.com/widgets\n\ngo 1.27.0\n",
		"README.md": "# widgets\n",
	}
	entries := make([]bundle.Entry, 0, len(src))
	for path, body := range src {
		entries = append(entries, bundle.Entry{
			Path: path,
			Type: bundle.TypeFile,
			Mode: 0o644,
			Size: int64(len(body)),
			Open: content(body),
		})
	}

	var buf bytes.Buffer
	m, err := (&bundle.Writer{}).Write(ctx, &buf, bundle.Manifest{
		Kind:   bundle.KindFolder,
		Client: "chainctl v0.1.2",
	}, entries)
	if err != nil {
		fmt.Println(err)
		return
	}
	// The manifest is filled in from what was packed, and the bytes do not
	// depend on the order the entries were handed over in.
	fmt.Println(m.Format, m.Kind, m.Files, m.Bytes)
	fmt.Println(m.TreeSHA256)

	// Output:
	// harness-source-v1 folder 3 77
	// 9ce04431168e90a1a01f1281360022033f7510b5aa6b259400441d12ad8e6e56
}

// Intake walks a source tree and keeps what the format can carry, warning about
// the rest. It asks Normalize rather than restating the rules, so what it drops
// and what the writer would refuse stay the same set.
func ExampleNormalize() {
	candidates := []bundle.Entry{
		{Path: "cmd/run.sh", Type: bundle.TypeFile, Mode: 0o700, Open: content("")},
		{Path: "docs/entry.md", Type: bundle.TypeSymlink, LinkTarget: "../README.md"},
		{Path: "etc/passwd", Type: bundle.TypeSymlink, LinkTarget: "../../../etc/passwd"},
		{Path: "var/run.sock", Type: bundle.EntryType("socket")},
	}
	keep := make([]bundle.Entry, 0, len(candidates))
	for _, c := range candidates {
		e, err := bundle.Normalize(c)
		if err != nil {
			// Intake warns and moves on; the writer, handed the same entry,
			// would refuse the whole bundle.
			r, _ := errors.AsType[*bundle.Error](err)
			fmt.Printf("dropped %s: %s\n", c.Path, r.Reason)
			continue
		}
		keep = append(keep, e)
	}
	for _, e := range keep {
		fmt.Printf("kept %s %04o\n", e.Path, e.Mode.Perm())
	}

	// Output:
	// dropped etc/passwd: symlink-escape
	// dropped var/run.sock: unsupported-entry-type
	// kept cmd/run.sh 0755
	// kept docs/entry.md 0777
}

func ExampleValidate() {
	ctx := context.Background()
	var buf bytes.Buffer
	if _, err := (&bundle.Writer{}).Write(ctx, &buf, bundle.Manifest{
		Kind:   bundle.KindFolder,
		Client: "chainctl v0.1.2",
	}, []bundle.Entry{{
		Path: "main.go",
		Type: bundle.TypeFile,
		Mode: 0o644,
		Size: 29,
		Open: content("package main\n\nfunc main() {}\n"),
	}}); err != nil {
		fmt.Println(err)
		return
	}

	// The server runs Validate over every bundle, whatever the client checked.
	m, err := bundle.Validate(ctx, bytes.NewReader(buf.Bytes()))
	fmt.Println(m.Kind, m.Files, err)

	// A rejection carries a stable reason to surface, not just a message.
	_, err = bundle.Validate(ctx, strings.NewReader("not a bundle"))
	if e, ok := errors.AsType[*bundle.Error](err); ok {
		fmt.Println(e.Reason)
	}

	// Output:
	// folder 1 <nil>
	// malformed-bundle
}

func ExampleReader_Read() {
	ctx := context.Background()
	var buf bytes.Buffer
	if _, err := (&bundle.Writer{}).Write(ctx, &buf, bundle.Manifest{
		Kind:   bundle.KindFiles,
		Client: "chainctl v0.1.2",
	}, []bundle.Entry{{
		Path: "notes.txt", Type: bundle.TypeFile, Mode: 0o600, Size: 6, Open: content("hello\n"),
	}, {
		Path: "run.sh", Type: bundle.TypeFile, Mode: 0o755, Size: 12, Open: content("echo hello\n\n"),
	}}); err != nil {
		fmt.Println(err)
		return
	}

	// Entries arrive sorted, with modes normalized. Nothing read here is
	// trustworthy until Read returns: tree_sha256 is only confirmed at the end,
	// so a consumer that writes must stage and commit on success.
	_, err := (&bundle.Reader{}).Read(ctx, &buf, func(e bundle.Entry, body io.Reader) error {
		b, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		fmt.Printf("%s %04o %d\n", e.Path, e.Mode.Perm(), len(b))
		return nil
	})
	fmt.Println(err)

	// Output:
	// notes.txt 0644 6
	// run.sh 0755 12
	// <nil>
}
