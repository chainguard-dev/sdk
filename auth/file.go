/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc/credentials"
	"knative.dev/pkg/logging"
)

type fileAuth struct {
	file   string
	secure bool

	m           sync.Mutex
	lastUpdated time.Time
	cache       []byte
}

var _ credentials.PerRPCCredentials = (*fileAuth)(nil)

// GetRequestMetadata implements credentials.PerRPCCredentials
func (fa *fileAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) { //nolint: revive
	fa.m.Lock()
	defer fa.m.Unlock()

	// According to https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	//    The kubelet proactively rotates the token if it is older than 80% of
	//    its total TTL, or if the token is older than 24 hours.
	// We use the minimum lifetime (10m) so the nearest a token should get to
	// expiration before it is refreshed is 2 minutes.  Use 1 minute to give
	// us wiggle room, but we should have 2 minutes to validate the token.
	if time.Since(fa.lastUpdated) < time.Minute {
		logging.FromContext(ctx).Infof("Using cached token, last refreshed %v", fa.lastUpdated)
		return map[string]string{
			"Authorization": string(fa.cache),
		}, nil
	}

	b, err := os.ReadFile(fa.file)
	if err != nil {
		return nil, err
	}
	fa.cache = b
	fa.lastUpdated = time.Now()
	logging.FromContext(ctx).Info("Using fresh token.")
	return map[string]string{
		"Authorization": string(b),
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (fa *fileAuth) RequireTransportSecurity() bool {
	return fa.secure
}
