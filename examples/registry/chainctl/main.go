/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"chainguard.dev/sdk/auth"
	"chainguard.dev/sdk/auth/ggcr"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func main() {
	ctx := context.Background()
	ts := auth.NewChainctlTokenSource(ctx, auth.WithAudience("cgr.dev"))

	desc, err := remote.Get(name.MustParseReference("cgr.dev/chainguard/static"), remote.WithAuthFromKeychain(ggcr.TokenSourceKeychain(ts)))
	if err != nil {
		log.Fatalf("error getting reference: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(desc)
}
