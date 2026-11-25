/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package upload

import (
	"context"
	"testing"
	"time"
)

// Sleep is the wait between attempts, for a test of its own.
var Sleep = sleep

// ReplaceSleep makes the wait between attempts f for the rest of t.
func ReplaceSleep(t *testing.T, f func(context.Context, time.Duration) error) {
	t.Helper()
	old := sleep
	sleep = f
	t.Cleanup(func() { sleep = old })
}
