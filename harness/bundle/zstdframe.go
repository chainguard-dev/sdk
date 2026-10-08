/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import "io"

// A zstd stream is a sequence of frames, and a decoder silently skips the
// skippable kind: it reads the magic, reads the declared length, and drops
// that many bytes. Those bytes are never decoded, so neither the trailer check
// nor tree_sha256 ever sees them, and they still count toward the compressed
// size the ratio guard divides by. That is the channel the 32 KiB trailer
// allowance exists to close, except hundreds of megabytes wide.
//
// A bundle is one frame, so the guard says exactly that: one standard frame,
// nothing before it and nothing after. It walks the frame's own structure
// rather than scanning for the skippable magic, because those four bytes occur
// in compressed data by chance often enough to matter over a 512 MiB cap.
const (
	zstdMagic         = 0xFD2FB528
	zstdSkippableMask = 0xFFFFFFF0
	zstdSkippableMag  = 0x184D2A50
)

type guardState int

const (
	wantMagic guardState = iota
	wantDescriptor
	wantBlockHeader
	skipBytes
	wantChecksum
	wantEOF
)

// zstdFrameGuard passes bytes through to the decoder while tracking where it
// is in the frame, so a second frame of any kind is refused at its first byte.
type zstdFrameGuard struct {
	r io.Reader

	state guardState
	// buf accumulates one header field group, and the widest group the guard
	// reads is 4 bytes; runs it only counts past are skipped, never buffered.
	buf      [4]byte
	have     int
	need     int
	remain   int64      // uninspected bytes still to pass through
	after    guardState // where remain running out lands
	checksum bool
	err      *Error
}

func (g *zstdFrameGuard) Read(p []byte) (int, error) {
	if g.err != nil {
		return 0, g.err
	}
	n, err := g.r.Read(p)
	if n > 0 {
		if e := g.walk(p[:n]); e != nil {
			g.err = e
			return 0, e
		}
	}
	return n, err
}

// walk advances the state machine over bytes already handed to the decoder.
// Rejecting after the fact is sound because the decoder cannot act on a frame
// it has not finished reading, and Read returns the error before the caller
// sees the entry that frame would have carried.
func (g *zstdFrameGuard) walk(b []byte) *Error {
	for len(b) > 0 {
		if g.state == skipBytes {
			take := min(int64(len(b)), g.remain)
			b = b[take:]
			g.remain -= take
			if g.remain == 0 {
				g.state, g.need, g.have = g.after, 0, 0
			}
			continue
		}
		if g.state == wantEOF {
			return rejectf(ReasonMalformed, "", "trailing data after the bundle's zstd frame: a bundle is one frame, and bytes outside it are never decoded")
		}
		if g.need == 0 {
			switch g.state {
			case wantMagic:
				g.need = 4
			case wantDescriptor:
				g.need = 1
			case wantBlockHeader:
				g.need = 3
			case wantChecksum:
				g.need = 4
			}
			g.have = 0
		}
		take := min(len(b), g.need-g.have)
		copy(g.buf[g.have:], b[:take])
		g.have += take
		b = b[take:]
		if g.have < g.need {
			return nil
		}
		if e := g.field(); e != nil {
			return e
		}
	}
	return nil
}

// field consumes one completed header field group and chooses the next state.
func (g *zstdFrameGuard) field() *Error {
	v := uint32(g.buf[0]) | uint32(g.buf[1])<<8 | uint32(g.buf[2])<<16 | uint32(g.buf[3])<<24
	switch g.state {
	case wantMagic:
		if v&zstdSkippableMask == zstdSkippableMag {
			return rejectf(ReasonMalformed, "", "bundle carries a zstd skippable frame: its bytes are never decoded, so no entry declares them and tree_sha256 does not cover them")
		}
		if v != zstdMagic {
			return rejectf(ReasonMalformed, "", "bundle does not start with a zstd frame")
		}
		g.state, g.need = wantDescriptor, 0

	case wantDescriptor:
		d := g.buf[0]
		fcsFlag := d >> 6
		singleSegment := d&0x20 != 0
		g.checksum = d&0x04 != 0
		didFlag := d & 0x03

		rest := 0
		if !singleSegment {
			rest++ // window descriptor
		}
		switch didFlag {
		case 1:
			rest++
		case 2:
			rest += 2
		case 3:
			rest += 4
		}
		switch fcsFlag {
		case 0:
			if singleSegment {
				rest++
			}
		case 1:
			rest += 2
		case 2:
			rest += 4
		case 3:
			rest += 8
		}
		if rest == 0 {
			g.state, g.need = wantBlockHeader, 0
		} else {
			// The remainder reaches 13 bytes, wider than buf, and nothing in
			// it is read: the guard only has to land on the first block.
			g.state, g.remain, g.after = skipBytes, int64(rest), wantBlockHeader
		}

	case wantBlockHeader:
		h := uint32(g.buf[0]) | uint32(g.buf[1])<<8 | uint32(g.buf[2])<<16
		last := h&1 != 0
		btype := (h >> 1) & 3
		size := int64(h >> 3)
		if btype == 3 {
			return rejectf(ReasonMalformed, "", "bundle has a reserved zstd block type")
		}
		if btype == 1 { // RLE: one byte on the wire, size is the expansion
			size = 1
		}
		// Where the frame lands once its last block has passed through.
		end := wantEOF
		if g.checksum {
			end = wantChecksum
		}
		switch {
		case size > 0:
			g.state, g.remain, g.after = skipBytes, size, wantBlockHeader
			if last {
				g.after = end
			}
		case last:
			g.state, g.need, g.have = end, 0, 0
		default:
			g.state, g.need = wantBlockHeader, 0
		}

	case wantChecksum:
		g.state, g.need = wantEOF, 0
	}
	return nil
}
