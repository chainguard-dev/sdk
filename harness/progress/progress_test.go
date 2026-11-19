/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const mib = 1 << 20

// fakeClock is a clock the test moves by hand.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time { return c.t }

func TestPrinterPlain(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	var out bytes.Buffer
	p := &Printer{w: &out, now: clock.now}
	size := int64(100 * mib)

	p.Infof("Opening it\n")
	for _, step := range []struct {
		after time.Duration
		u     Update
	}{
		{0, Update{Phase: "downloading", Total: size}},
		{10 * time.Second, Update{Phase: "downloading", Bytes: 10 * mib, Total: size}},
		{19 * time.Second, Update{Phase: "downloading", Bytes: 29 * mib, Total: size}},
		{time.Second, Update{Phase: "downloading", Bytes: 30 * mib, Total: size}},
		{29 * time.Second, Update{Phase: "verifying", Total: size}},
		{time.Second, Update{Phase: "extracting", Files: 3, Bytes: 60 * mib}},
		{5 * time.Second, Update{Phase: "renaming"}},
	} {
		clock.t = clock.t.Add(step.after)
		p.Report(step.u)
	}
	p.Done()
	p.Infof("Opened it\n")

	want := []string{
		"Opening it",
		"downloading 30.0 MiB / 100.0 MiB (30%), 1.0 MiB/s, 30s elapsed",
		"extracting 3 files, 60.0 MiB, 0 B/s, 1m0s elapsed",
		"renaming, 1m5s elapsed",
		"Opened it",
	}
	if diff := cmp.Diff(want, strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}
}

func TestPrinterTerminal(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	var out bytes.Buffer
	p := &Printer{w: &out, tty: true, now: clock.now}

	p.Infof("Opening it\n")
	p.Report(Update{Phase: "downloading", Total: 4 * mib})
	clock.t = clock.t.Add(500 * time.Millisecond)
	p.Report(Update{Phase: "downloading", Bytes: mib, Total: 4 * mib})
	clock.t = clock.t.Add(500 * time.Millisecond)
	p.Report(Update{Phase: "downloading", Bytes: 2 * mib, Total: 4 * mib})
	clock.t = clock.t.Add(time.Second)
	p.Report(Update{Phase: "writing", Bytes: 512})
	p.Done()
	p.Infof("Opened it\n")

	want := "Opening it\n" +
		"\rdownloading 0 B / 4.0 MiB (0%), 0 B/s, 0s elapsed\x1b[K" +
		"\rdownloading 2.0 MiB / 4.0 MiB (50%), 2.0 MiB/s, 1s elapsed\x1b[K" +
		"\rwriting 512 B, 0 B/s, 2s elapsed\x1b[K\n" +
		"Opened it\n"
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("output (-want, +got):\n%s", diff)
	}
}

// An Infof line mid-progress starts on a fresh line, not over the progress.
func TestPrinterInfofEndsTheTerminalLine(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	var out bytes.Buffer
	p := &Printer{w: &out, tty: true, now: clock.now}

	p.Report(Update{Phase: "renaming"})
	p.Infof("warning: slow\n")

	if want := "\rrenaming, 0s elapsed\x1b[K\nwarning: slow\n"; out.String() != want {
		t.Errorf("output: got %q, want %q", out.String(), want)
	}
}

func TestPrinterQuiet(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	for _, tty := range []bool{false, true} {
		var out bytes.Buffer
		p := &Printer{Quiet: true, w: &out, tty: tty, now: clock.now}

		p.Infof("Opening it\n")
		p.Report(Update{Phase: "downloading", Total: 1})
		clock.t = clock.t.Add(time.Minute)
		p.Report(Update{Phase: "downloading", Bytes: 1, Total: 1})
		p.Done()
		p.Infof("Opened it\n")

		if want := "Opening it\nOpened it\n"; out.String() != want {
			t.Errorf("tty %t: got %q, want %q", tty, out.String(), want)
		}
	}
}

func TestPrinterDoneWithoutReports(t *testing.T) {
	var out bytes.Buffer
	p := &Printer{w: &out, tty: true, now: time.Now}
	p.Done()
	if out.Len() != 0 {
		t.Errorf("output: got %q, want nothing", out.String())
	}
}

func TestNewIsPlainOffATerminal(t *testing.T) {
	if p := New(&bytes.Buffer{}); p.tty {
		t.Error("New: got a terminal printer for a buffer")
	}
}

func TestSize(t *testing.T) {
	for n, want := range map[int64]string{0: "0 B", 1023: "1023 B", 1 << 10: "1.0 KiB", 1717986918: "1.6 GiB", 3 << 40: "3.0 TiB"} {
		if got := Size(n); got != want {
			t.Errorf("Size(%d): got %q, want %q", n, got, want)
		}
	}
}

func TestPrinterDoneDoesNotRepeatAPrintedLine(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	var out bytes.Buffer
	p := &Printer{w: &out, now: clock.now}

	p.Report(Update{Phase: "downloading", Total: 100})
	clock.t = clock.t.Add(30 * time.Second)
	p.Report(Update{Phase: "downloading", Bytes: 100, Total: 100})
	p.Done()

	want := "downloading 100 B / 100 B (100%), 3 B/s, 30s elapsed\n"
	if got := out.String(); got != want {
		t.Errorf("output: got = %q, want = %q", got, want)
	}
}
