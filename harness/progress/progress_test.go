/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package progress

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

// barPrinter is a Printer that draws the bar on a terminal cols wide, with
// its output and its clock.
func barPrinter(cols int) (*Printer, *bytes.Buffer, *fakeClock) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	var out bytes.Buffer
	return &Printer{w: &out, tty: true, bar: true, cols: func() int { return cols }, now: clock.now}, &out, clock
}

// frames splits a terminal's output into the lines it drew in place.
func frames(out string) []string {
	var lines []string
	for l := range strings.SplitSeq(out, "\r") {
		if l = strings.TrimSuffix(strings.TrimSuffix(l, "\n"), "\x1b[K"); l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func TestPrinterBar(t *testing.T) {
	p, out, clock := barPrinter(80)

	for _, step := range []struct {
		after time.Duration
		u     Update
	}{
		{0, Update{Phase: "uploading", Total: 4 * mib}},
		// Sooner than the bar redraws.
		{100 * time.Millisecond, Update{Phase: "uploading", Bytes: mib / 2, Total: 4 * mib}},
		{900 * time.Millisecond, Update{Phase: "uploading", Bytes: mib, Total: 4 * mib}},
		{time.Second, Update{Phase: "uploading", Bytes: 3 * mib, Total: 4 * mib}},
		{time.Second, Update{Phase: "uploading", Bytes: 4 * mib, Total: 4 * mib}},
		// Past its total, Linky stays at the end of the bar.
		{time.Second, Update{Phase: "uploading", Bytes: 5 * mib, Total: 4 * mib}},
		{time.Second, Update{Phase: "building", Files: 7, Bytes: mib}},
		{time.Second, Update{Phase: "building", Files: 9, Bytes: 2 * mib}},
		{time.Second, Update{Phase: "writing", Bytes: 512}},
		{time.Second, Update{Phase: "renaming"}},
	} {
		clock.t = clock.t.Add(step.after)
		p.Report(step.u)
	}
	p.Done()

	want := []string{
		"uploading   🐙························    0%  0 B / 4.0 MiB  0 B/s  0s",
		// The bar gives up the columns the longer numbers need.
		"uploading   ~≈~~≈🐙················   25%  1.0 MiB / 4.0 MiB  1.0 MiB/s  eta 3s",
		"uploading   ≈~~≈~~≈~~≈~~≈~~🐙······   75%  3.0 MiB / 4.0 MiB  1.5 MiB/s  eta 0s",
		"uploading   ~~≈~~≈~~≈~~≈~~≈~~≈~~≈🐙  100%  4.0 MiB / 4.0 MiB  1.3 MiB/s  3s",
		"uploading   ~≈~~≈~~≈~~≈~~≈~~≈~~≈~🐙  125%  5.0 MiB / 4.0 MiB  1.2 MiB/s  4s",
		"building    🐙························  7 files  1.0 MiB  0 B/s  5s",
		"building    ~🐙·······················  9 files  2.0 MiB  2.0 MiB/s  6s",
		"writing     🐙························  512 B  0 B/s  7s",
		"renaming    🐙························  8s",
	}
	if diff := cmp.Diff(want, frames(out.String())); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}
}

func TestPrinterBarColor(t *testing.T) {
	p, out, clock := barPrinter(80)
	p.color = true

	p.Report(Update{Phase: "uploading", Total: 4 * mib})
	clock.t = clock.t.Add(time.Second)
	p.Report(Update{Phase: "uploading", Bytes: mib, Total: 4 * mib})

	want := []string{
		"\x1b[1muploading  \x1b[0m 🐙\x1b[2m························\x1b[0m    0%  0 B / 4.0 MiB  0 B/s  \x1b[2m0s\x1b[0m",
		"\x1b[1muploading  \x1b[0m \x1b[36m~≈~~≈\x1b[0m🐙\x1b[2m················\x1b[0m   25%  1.0 MiB / 4.0 MiB  1.0 MiB/s  \x1b[2meta 3s\x1b[0m",
	}
	if diff := cmp.Diff(want, frames(out.String())); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}
}

// With no total, Linky swims to the end of the bar and back, trailing a wake
// on the way out.
func TestPrinterBarWithoutATotalSwimsThereAndBack(t *testing.T) {
	// Room for the shortest bar and no more: Linky and six places to swim.
	cols := len("renaming    ") + minBar + len("  0s") + 1
	p, out, clock := barPrinter(cols)

	for range 14 {
		p.Report(Update{Phase: "renaming"})
		clock.t = clock.t.Add(200 * time.Millisecond)
	}

	want := []string{
		"renaming    🐙······  0s",
		"renaming    ~🐙·····  0s",
		"renaming    ≈~🐙····  0s",
		"renaming    ~~≈🐙···  0s",
		"renaming    ·≈~~🐙··  0s",
		"renaming    ··~≈~🐙·  1s",
		"renaming    ···~~≈🐙  1s",
		"renaming    ·····🐙·  1s",
		"renaming    ····🐙··  1s",
		"renaming    ···🐙···  1s",
		"renaming    ··🐙····  2s",
		"renaming    ·🐙·····  2s",
		"renaming    🐙······  2s",
		"renaming    ~🐙·····  2s",
	}
	if diff := cmp.Diff(want, frames(out.String())); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}

	// One column fewer leaves no room for the shortest bar.
	p, out, _ = barPrinter(cols - 1)
	p.Report(Update{Phase: "renaming"})
	if diff := cmp.Diff([]string{"renaming, 0s elapsed"}, frames(out.String())); diff != "" {
		t.Errorf("lines one column narrower (-want, +got):\n%s", diff)
	}
}

// A line wider than the terminal wraps, and then can't be rewritten in place.
func TestPrinterBarFitsTheTerminal(t *testing.T) {
	updates := []Update{
		{Phase: "downloading", Bytes: 1717986918, Total: 3 << 40},
		{Phase: "extracting", Files: 41872, Bytes: 5 << 30},
		{Phase: "uploading", Bytes: 512},
		{Phase: "renaming"},
	}
	for cols := range 140 {
		for _, u := range updates {
			p, out, clock := barPrinter(cols)
			p.Report(u)
			clock.t = clock.t.Add(90 * time.Second)
			p.Report(u)

			for _, l := range frames(out.String()) {
				if !strings.Contains(l, "🐙") {
					continue
				}
				// Linky takes two cells and every other rune one.
				if cells := utf8.RuneCountInString(l) + 1; cells >= cols {
					t.Errorf("%d columns, %s: line takes %d cells: %q", cols, u.Phase, cells, l)
				}
			}
		}
	}
}

func TestPrinterBarFallsBackToTextOnANarrowTerminal(t *testing.T) {
	p, out, clock := barPrinter(40)

	// The text line is wider than this terminal, so each redraw leaves a
	// line behind: it keeps the text line's interval, not the bar's.
	for range 11 {
		p.Report(Update{Phase: "uploading", Bytes: mib, Total: 4 * mib})
		clock.t = clock.t.Add(200 * time.Millisecond)
	}

	want := []string{
		"uploading 1.0 MiB / 4.0 MiB (25%), 0 B/s, 0s elapsed",
		"uploading 1.0 MiB / 4.0 MiB (25%), 1.0 MiB/s, 1s elapsed",
		"uploading 1.0 MiB / 4.0 MiB (25%), 512.0 KiB/s, 2s elapsed",
	}
	if diff := cmp.Diff(want, frames(out.String())); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}
}

// A few bytes into a large total, the time left is longer than anyone waits
// and too wide to keep room for.
func TestPrinterBarShowsNoDistantEta(t *testing.T) {
	p, out, clock := barPrinter(80)
	total := int64(100 << 30)

	p.Report(Update{Phase: "extracting", Total: total})
	clock.t = clock.t.Add(time.Second)
	p.Report(Update{Phase: "extracting", Bytes: 10, Total: total})
	clock.t = clock.t.Add(time.Second)
	p.Report(Update{Phase: "extracting", Bytes: total / 2, Total: total})

	want := []string{
		"extracting  🐙························    0%  0 B / 100.0 GiB  0 B/s  0s",
		"extracting  🐙························    0%  10 B / 100.0 GiB  10 B/s  1s",
		"extracting  ≈~~≈~~≈~🐙·········   50%  50.0 GiB / 100.0 GiB  25.0 GiB/s  eta 2s",
	}
	if diff := cmp.Diff(want, frames(out.String())); diff != "" {
		t.Errorf("lines (-want, +got):\n%s", diff)
	}
}

func TestStyle(t *testing.T) {
	for _, tc := range []struct {
		name       string
		env        map[string]string
		bar, color bool
	}{
		{name: "UTF-8 locale", env: map[string]string{"LANG": "en_US.UTF-8", "TERM": "xterm-256color"}, bar: true, color: true},
		{name: "lowercase utf8", env: map[string]string{"LANG": "C.utf8"}, bar: true, color: true},
		{name: "NO_COLOR", env: map[string]string{"LANG": "en_US.UTF-8", "NO_COLOR": "1"}, bar: true},
		{name: "no locale", env: map[string]string{"TERM": "xterm-256color"}},
		{name: "LC_ALL overrides LANG", env: map[string]string{"LC_ALL": "C", "LANG": "en_US.UTF-8"}},
		{name: "LC_CTYPE overrides LANG", env: map[string]string{"LC_CTYPE": "en_US.UTF-8", "LANG": "C"}, bar: true, color: true},
		{name: "dumb terminal", env: map[string]string{"LANG": "en_US.UTF-8", "TERM": "dumb"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bar, color := style(func(k string) string { return tc.env[k] })
			if bar != tc.bar || color != tc.color {
				t.Errorf("style: got = (%t, %t), want = (%t, %t)", bar, color, tc.bar, tc.color)
			}
		})
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

// New draws the bar only on a terminal, sized to it.
func TestNew(t *testing.T) {
	// A locale that would draw the bar, so only the writer decides.
	t.Setenv("LC_ALL", "en_US.UTF-8")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	file, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer file.Close()
	// A character device, which New takes for a terminal. Its size can't be
	// read, which leaves no room for the bar.
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("OpenFile(%s): %v", os.DevNull, err)
	}
	defer null.Close()

	for _, tc := range []struct {
		name string
		w    io.Writer
		want Printer
	}{
		{name: "buffer", w: &bytes.Buffer{}},
		{name: "file", w: file},
		{name: "character device", w: null, want: Printer{tty: true, bar: true, color: true, cols: func() int { return 0 }}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.w)
			tc.want.w = tc.w
			if diff := cmp.Diff(tc.want, *p,
				cmp.AllowUnexported(Printer{}),
				cmpopts.IgnoreFields(Printer{}, "now"),
				cmp.Comparer(func(a, b io.Writer) bool { return a == b }),
				// Widths compare by the columns they report.
				cmp.Comparer(func(a, b func() int) bool {
					if a == nil || b == nil {
						return a == nil && b == nil
					}
					return a() == b()
				}),
			); diff != "" {
				t.Errorf("New (-want, +got):\n%s", diff)
			}
		})
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
