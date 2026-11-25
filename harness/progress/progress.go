/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// Update is an operation's count so far in Phase. Total is Bytes' final
// value, or 0 where it isn't known; Files is 0 where files aren't counted.
type Update struct {
	Phase               string
	Bytes, Total, Files int64
}

// Printer writes progress to w: on a terminal one line rewritten in place,
// and otherwise a plain line every 30 seconds. Where the terminal can show
// it, the line is a bar that Linky the octopus swims along. Quiet drops the
// progress but not Infof's lines. It is not safe for concurrent use.
type Printer struct {
	Quiet bool

	w   io.Writer
	tty bool
	// bar draws the terminal's line as the bar, color adds its escapes, and
	// cols is the terminal's width.
	bar, color bool
	cols       func() int
	// frame counts the phase's bar lines, and reserve is the widest text
	// beside the bar so far, so the bar doesn't grow and shrink with it.
	frame, reserve int
	// drawn is whether the last line was the bar.
	drawn bool
	now   func() time.Time

	cur                     Update
	started, pending, dirty bool
	start, phaseStart, last time.Time
}

// New returns a Printer to w, rewriting its line only if w is a terminal.
func New(w io.Writer) *Printer {
	p := &Printer{w: w, now: time.Now}
	if f, ok := w.(*os.File); ok {
		fi, err := f.Stat()
		p.tty = err == nil && fi.Mode()&os.ModeCharDevice != 0
		if p.tty {
			p.bar, p.color = style(os.Getenv)
			p.cols = func() int {
				// An error leaves 0 columns, which is no room for the bar.
				cols, _, _ := term.GetSize(int(f.Fd()))
				return cols
			}
		}
	}
	return p
}

// style reports whether the environment's terminal can draw the bar, which
// needs UTF-8, and whether to color it (https://no-color.org).
func style(getenv func(string) string) (bar, color bool) {
	locale := getenv("LC_ALL")
	if locale == "" {
		locale = getenv("LC_CTYPE")
	}
	if locale == "" {
		locale = getenv("LANG")
	}
	locale = strings.ToUpper(locale)
	bar = getenv("TERM") != "dumb" && (strings.Contains(locale, "UTF-8") || strings.Contains(locale, "UTF8"))
	return bar, bar && getenv("NO_COLOR") == ""
}

// Report records u, printing it once the interval since the last line has
// passed, or at once on a terminal when the phase changes.
func (p *Printer) Report(u Update) {
	if p.Quiet {
		return
	}
	now := p.now()
	if !p.started {
		p.start, p.last, p.started = now, now, true
	}
	changed := u.Phase != p.cur.Phase
	if changed {
		p.phaseStart, p.frame, p.reserve = now, 0, 0
	}
	p.cur, p.dirty = u, true
	every := 30 * time.Second
	switch {
	case p.drawn:
		// Often enough that Linky's wake moves. Not for the text a narrow
		// terminal gets instead: that line can wrap, and each redraw of a
		// wrapped line leaves a line behind.
		every = 200 * time.Millisecond
	case p.tty:
		every = time.Second
	}
	if (p.tty && changed) || now.Sub(p.last) >= every {
		p.print(now)
	}
}

// Done prints the last update if no line shows it yet, and ends a terminal's
// line so later output doesn't overwrite it.
func (p *Printer) Done() {
	if p.dirty {
		p.print(p.now())
	}
	p.started = false
	p.endLine()
}

// Infof ends any progress line on a terminal, then prints as [fmt.Fprintf].
func (p *Printer) Infof(format string, a ...any) {
	p.endLine()
	fmt.Fprintf(p.w, format, a...)
}

func (p *Printer) endLine() {
	if p.pending {
		fmt.Fprintln(p.w)
		p.pending = false
	}
}

func (p *Printer) print(now time.Time) {
	p.last, p.dirty = now, false
	elapsed := now.Sub(p.start).Truncate(time.Second)
	var rate int64
	if s := now.Sub(p.phaseStart).Seconds(); s > 0 {
		rate = int64(float64(p.cur.Bytes) / s)
	}
	line := p.barLine(rate, elapsed)
	p.drawn = line != ""
	switch {
	case p.drawn:
	case p.cur.Total > 0:
		line = fmt.Sprintf("%s %s / %s (%d%%), %s/s, %s elapsed", p.cur.Phase, Size(p.cur.Bytes), Size(p.cur.Total),
			p.cur.Bytes*100/p.cur.Total, Size(rate), elapsed)
	case p.cur.Files > 0:
		line = fmt.Sprintf("%s %d files, %s, %s/s, %s elapsed", p.cur.Phase, p.cur.Files, Size(p.cur.Bytes), Size(rate), elapsed)
	case p.cur.Bytes > 0:
		line = fmt.Sprintf("%s %s, %s/s, %s elapsed", p.cur.Phase, Size(p.cur.Bytes), Size(rate), elapsed)
	default:
		line = fmt.Sprintf("%s, %s elapsed", p.cur.Phase, elapsed)
	}
	if p.tty {
		// Clearing to the end of the line drops a longer line's tail.
		fmt.Fprintf(p.w, "\r%s\x1b[K", line)
		p.pending = true
	} else {
		fmt.Fprintln(p.w, line)
	}
}

const (
	maxBar = 26
	minBar = 8
	// linkyCells is how many terminal cells the octopus takes.
	linkyCells = 2
	// maxEta is the longest time left the bar shows. A few bytes into a
	// large total the estimate is longer still, and its width would stay
	// reserved for the rest of the phase.
	maxEta = 100 * time.Hour
)

// barLine returns the terminal's line with a bar sized to the terminal, or ""
// where there is no bar or no room for one: a line wider than the terminal
// wraps, and a wrapped line can't be rewritten in place.
func (p *Printer) barLine(rate int64, elapsed time.Duration) string {
	if !p.bar {
		return ""
	}
	u := p.cur
	phase := fmt.Sprintf("%-11s", u.Phase)
	var stats string
	// What's left, or how long it has run where that isn't known.
	clock := elapsed.String()
	switch {
	case u.Total > 0:
		stats = fmt.Sprintf("%3d%%  %s / %s  %s/s  ", u.Bytes*100/u.Total, Size(u.Bytes), Size(u.Total), Size(rate))
		if rate > 0 && u.Bytes < u.Total && (u.Total-u.Bytes)/rate < int64(maxEta/time.Second) {
			clock = "eta " + (time.Duration((u.Total-u.Bytes)/rate) * time.Second).String()
		}
	case u.Files > 0:
		stats = fmt.Sprintf("%d files  %s  %s/s  ", u.Files, Size(u.Bytes), Size(rate))
	case u.Bytes > 0:
		stats = fmt.Sprintf("%s  %s/s  ", Size(u.Bytes), Size(rate))
	}
	// One column stays free: a line that fills the last one wraps on some
	// terminals.
	p.reserve = max(p.reserve, len(stats)+len(clock))
	width := min(maxBar, p.cols()-1-len(phase)-len("   ")-p.reserve)
	if width < minBar {
		return ""
	}

	// Linky sits at one of track+1 places, with water behind and dots ahead.
	track := width - linkyCells
	var at, wake int
	if u.Total > 0 {
		at = int(min(u.Bytes, u.Total) * int64(track) / u.Total)
		wake = at
	} else {
		// With no total to fill toward, Linky swims there and back.
		at = p.frame % (2 * track)
		wake = min(at, 3)
		if at > track {
			at = 2*track - at
			wake = 0
		}
	}
	p.frame++
	var water strings.Builder
	for i := at - wake; i < at; i++ {
		if (i+p.frame)%3 == 0 {
			water.WriteString("≈")
		} else {
			water.WriteString("~")
		}
	}
	return p.paint("1", phase) + " " +
		p.paint("2", strings.Repeat("·", at-wake)) + p.paint("36", water.String()) + "🐙" + p.paint("2", strings.Repeat("·", track-at)) +
		"  " + stats + p.paint("2", clock)
}

// paint wraps s in the SGR attribute code where the Printer colors its line.
func (p *Printer) paint(code, s string) string {
	if !p.color || s == "" {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

// Size formats n bytes in binary units.
func Size(n int64) string {
	if n < 1<<10 {
		return fmt.Sprintf("%d B", n)
	}
	v, unit := float64(n)/(1<<10), 0
	for v >= 1<<10 && unit < 3 {
		v, unit = v/(1<<10), unit+1
	}
	return fmt.Sprintf("%.1f %s", v, [...]string{"KiB", "MiB", "GiB", "TiB"}[unit])
}
