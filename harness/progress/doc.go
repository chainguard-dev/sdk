/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package progress renders a long operation's progress for a command's
// stderr: a phase with its byte and file counts, rate and elapsed time, as a
// line rewritten in place on a terminal or a plain line now and then
// elsewhere, between one-line summaries the command prints. It knows nothing
// of what the operation is, so any CLI can use it.
package progress
