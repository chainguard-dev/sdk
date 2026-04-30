/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"errors"
	"testing"

	"github.com/chainguard-dev/clog/slogtest"

	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1"
	eventstest "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1/test"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

func TestListSubscriptionsIter(t *testing.T) {
	ctx := slogtest.Context(t)

	sub1 := &events.Subscription{Uid: "uid-1", Sink: "https://example.com/hook1"}
	sub2 := &events.Subscription{Uid: "uid-2", Sink: "https://example.com/hook2"}
	sub3 := &events.Subscription{Uid: "uid-3", Sink: "https://example.com/hook3"}

	t.Run("single page", func(t *testing.T) {
		mock := &eventstest.MockClients{
			SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
				T: t,
				OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
					Given: &events.ListSubscriptionsRequest{},
					Result: &events.ListSubscriptionsResponse{
						Subscriptions: []*events.Subscription{sub1, sub2},
					},
				}},
			},
		}

		got, err := mock.ListSubscriptionsAll(ctx, &events.ListSubscriptionsRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("length: got = %d, wanted = 2", len(got))
		}
		if got[0].GetUid() != sub1.GetUid() {
			t.Errorf("item[0].uid: got = %q, wanted = %q", got[0].GetUid(), sub1.GetUid())
		}
		if got[1].GetUid() != sub2.GetUid() {
			t.Errorf("item[1].uid: got = %q, wanted = %q", got[1].GetUid(), sub2.GetUid())
		}
	})

	t.Run("empty result", func(t *testing.T) {
		mock := &eventstest.MockClients{
			SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
				T: t,
				OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
					Given:  &events.ListSubscriptionsRequest{},
					Result: &events.ListSubscriptionsResponse{},
				}},
			},
		}

		got, err := mock.ListSubscriptionsAll(ctx, &events.ListSubscriptionsRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("length: got = %d, wanted = 0", len(got))
		}
	})

	t.Run("error propagates", func(t *testing.T) {
		wantErr := errors.New("rpc error")
		mock := &eventstest.MockClients{
			SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
				T: t,
				OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
					Given: &events.ListSubscriptionsRequest{},
					Error: wantErr,
				}},
			},
		}

		_, err := mock.ListSubscriptionsAll(ctx, &events.ListSubscriptionsRequest{})
		if err == nil {
			t.Fatal("error: got = nil, wanted error")
		}
		if !errors.Is(err, wantErr) {
			t.Errorf("error: got = %v, wanted = %v", err, wantErr)
		}
	})

	t.Run("iter yields items one by one", func(t *testing.T) {
		mock := &eventstest.MockClients{
			SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
				T: t,
				OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
					Given: &events.ListSubscriptionsRequest{},
					Result: &events.ListSubscriptionsResponse{
						Subscriptions: []*events.Subscription{sub1, sub2, sub3},
					},
				}},
			},
		}

		var got []*events.Subscription
		for sub, err := range mock.ListSubscriptionsIter(ctx, &events.ListSubscriptionsRequest{}) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got = append(got, sub)
		}
		if len(got) != 3 {
			t.Fatalf("length: got = %d, wanted = 3", len(got))
		}
		for i, want := range []*events.Subscription{sub1, sub2, sub3} {
			if got[i].GetUid() != want.GetUid() {
				t.Errorf("item[%d].uid: got = %q, wanted = %q", i, got[i].GetUid(), want.GetUid())
			}
		}
	})

	t.Run("iter break stops early", func(t *testing.T) {
		mock := &eventstest.MockClients{
			SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
				T: t,
				OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
					Given: &events.ListSubscriptionsRequest{},
					Result: &events.ListSubscriptionsResponse{
						Subscriptions: []*events.Subscription{sub1, sub2, sub3},
					},
				}},
			},
		}

		count := 0
		for _, err := range mock.ListSubscriptionsIter(ctx, &events.ListSubscriptionsRequest{}) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			count++
			break
		}
		if count != 1 {
			t.Errorf("count: got = %d, wanted = 1", count)
		}
	})
}

func TestListSubscriptionsAll_MockPagination(t *testing.T) {
	ctx := slogtest.Context(t)

	// Simulate two pages of results using the mock.
	page1Req := &events.ListSubscriptionsRequest{PageSize: 2}
	page2Req := &events.ListSubscriptionsRequest{PageSize: 2, PageToken: "page2"}

	sub1 := &events.Subscription{Uid: "uid-1", Sink: "https://example.com/hook1"}
	sub2 := &events.Subscription{Uid: "uid-2", Sink: "https://example.com/hook2"}
	sub3 := &events.Subscription{Uid: "uid-3", Sink: "https://example.com/hook3"}

	mock := &eventstest.MockClients{
		SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
			T: t,
			OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{
				{
					Given: page1Req,
					Result: &events.ListSubscriptionsResponse{
						Subscriptions: []*events.Subscription{sub1, sub2},
						NextPageToken: "page2",
					},
				},
				{
					Given: page2Req,
					Result: &events.ListSubscriptionsResponse{
						Subscriptions: []*events.Subscription{sub3},
					},
				},
			},
		},
	}

	// The mock ListSubscriptionsAll only calls ListSubscriptions once (no real pagination),
	// so verify the mock correctly delegates to the service client for the first page.
	got, err := mock.ListSubscriptionsAll(ctx, page1Req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mock returns only the first page (no pagination in mock).
	if len(got) != 2 {
		t.Fatalf("length: got = %d, wanted = 2", len(got))
	}
	if got[0].GetUid() != sub1.GetUid() {
		t.Errorf("item[0].uid: got = %q, wanted = %q", got[0].GetUid(), sub1.GetUid())
	}
	if got[1].GetUid() != sub2.GetUid() {
		t.Errorf("item[1].uid: got = %q, wanted = %q", got[1].GetUid(), sub2.GetUid())
	}
}
