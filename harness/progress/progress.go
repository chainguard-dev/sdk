/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package progress

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Update is an operation's count so far in Phase. Total is Bytes' final
// value, or 0 where it isn't known; Files is 0 where files aren't counted.
type Update struct {
	Phase               string
	Bytes, Total, Files int64
}

// Printer writes progress to w: on a terminal one line rewritten in place
// about once a second, and otherwise a plain line every 30 seconds. Quiet
// drops the progress but not Infof's lines. It is not safe for concurrent use.
type Printer struct {
	Quiet bool

	w   io.Writer
	tty bool
	now func() time.Time

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
	}
	return p
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
		p.phaseStart = now
	}
	p.cur, p.dirty = u, true
	every := 30 * time.Second
	if p.tty {
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
	var line string
	switch {
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
