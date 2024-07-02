/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	notifications "chainguard.dev/sdk/proto/platform/notifications/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ notifications.Clients = (*MockNotificationsClients)(nil)

type MockNotificationsClients struct {
	OnClose error

	NotificationsClient MockNotificationsClient
}

func (m MockNotificationsClients) Notifications() notifications.NotificationsClient {
	return &m.NotificationsClient
}

func (m MockNotificationsClients) Close() error {
	return m.OnClose
}

var _ notifications.NotificationsClient = (*MockNotificationsClient)(nil)

type MockNotificationsClient struct {
	notifications.NotificationsClient

	OnListNotifications []NotificationsOnList
}

type NotificationsOnList struct {
	Given *notifications.NotificationsFilter
	List  *notifications.NotificationsList
	Error error
}

func (m MockNotificationsClient) List(_ context.Context, given *notifications.NotificationsFilter, _ ...grpc.CallOption) (*notifications.NotificationsList, error) { //nolint: revive
	for _, o := range m.OnListNotifications {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
