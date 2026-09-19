/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"chainguard.dev/sdk/auth"
	"chainguard.dev/sdk/auth/ggcr"
	"github.com/chainguard-dev/clog"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const (
	sub = "720909c9f5279097d847ad02a2f24ba8f59de36a/a033a6fabe0bfa0d"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	// We need a valid token to do the exchange - this happens to be convinient.
	ts := auth.NewChainctlTokenSource(ctx)

	desc, err := remote.Get(name.MustParseReference("cgr.dev/chainguard/static"), remote.WithAuthFromKeychain(ggcr.Keychain(sub, ts)))
	if err != nil {
		clog.FatalContextf(ctx, "error getting reference: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(desc)
}
