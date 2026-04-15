/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package iter_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"chainguard.dev/sdk/proto/chainguard/platform/iter"
	pb "chainguard.dev/sdk/proto/chainguard/platform/test"
	"github.com/chainguard-dev/clog/slogtest"
)

func TestAll(t *testing.T) {
	for _, tt := range []struct {
		name    string
		seq     func(yield func(string, error) bool)
		want    []string
		wantErr string
	}{{
		name: "collects all items",
		seq: func(yield func(string, error) bool) {
			_ = yield("a", nil) && yield("b", nil) && yield("c", nil)
		},
		want: []string{"a", "b", "c"},
	}, {
		name: "empty iterator",
		seq:  func(_ func(string, error) bool) {},
		want: nil,
	}, {
		name: "stops on error",
		seq: func(yield func(string, error) bool) {
			_ = yield("a", nil) && yield("", errors.New("boom"))
		},
		wantErr: "boom",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iter.All(tt.seq)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("error: got = nil, wanted = %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error: got = %q, wanted = %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("length: got = %d, wanted = %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("item[%d]: got = %q, wanted = %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	ctx := slogtest.Context(t)

	t.Run("single page", func(t *testing.T) {
		items := []string{fmt.Sprintf("item-%d", rand.Int64()), fmt.Sprintf("item-%d", rand.Int64())}
		got, err := iter.All(iter.List(ctx, "things", func(pageToken string) ([]string, string, error) {
			if pageToken != "" {
				t.Fatalf("unexpected page token: %q", pageToken)
			}
			return items, "", nil
		}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != len(items) {
			t.Fatalf("length: got = %d, wanted = %d", len(got), len(items))
		}
		for i := range got {
			if got[i] != items[i] {
				t.Errorf("item[%d]: got = %q, wanted = %q", i, got[i], items[i])
			}
		}
	})

	t.Run("multiple pages", func(t *testing.T) {
		pages := map[string]struct {
			items     []string
			nextToken string
		}{
			"":      {items: []string{"a", "b"}, nextToken: "page2"},
			"page2": {items: []string{"c", "d"}, nextToken: "page3"},
			"page3": {items: []string{"e"}, nextToken: ""},
		}
		got, err := iter.All(iter.List(ctx, "things", func(pageToken string) ([]string, string, error) {
			page, ok := pages[pageToken]
			if !ok {
				t.Fatalf("unexpected page token: %q", pageToken)
			}
			return page.items, page.nextToken, nil
		}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"a", "b", "c", "d", "e"}
		if len(got) != len(want) {
			t.Fatalf("length: got = %d, wanted = %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("item[%d]: got = %q, wanted = %q", i, got[i], want[i])
			}
		}
	})

	t.Run("fetch error wraps resource name", func(t *testing.T) {
		for _, err := range iter.List(ctx, "widgets", func(_ string) ([]string, string, error) {
			return nil, "", errors.New("connection refused")
		}) {
			if err == nil {
				t.Fatal("error: got = nil, wanted error")
			}
			want := `listing "widgets": connection refused`
			if err.Error() != want {
				t.Fatalf("error: got = %q, wanted = %q", err.Error(), want)
			}
			return
		}
		t.Fatal("iterator yielded no items")
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		for _, err := range iter.List(ctx, "things", func(_ string) ([]string, string, error) {
			t.Fatal("fetch should not be called after context cancellation")
			return nil, "", nil
		}) {
			if err == nil {
				t.Fatal("error: got = nil, wanted context.Canceled")
			}
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("error: got = %v, wanted context.Canceled", err)
			}
			return
		}
		t.Fatal("iterator yielded no items")
	})

	t.Run("empty page", func(t *testing.T) {
		got, err := iter.All(iter.List(ctx, "things", func(_ string) ([]string, string, error) {
			return nil, "", nil
		}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("length: got = %d, wanted = 0", len(got))
		}
	})

	t.Run("break stops iteration", func(t *testing.T) {
		fetchCount := 0
		var got []string
		for item, err := range iter.List(ctx, "things", func(_ string) ([]string, string, error) {
			fetchCount++
			return []string{"a", "b"}, "next", nil
		}) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got = append(got, item)
			if len(got) == 1 {
				break
			}
		}
		if len(got) != 1 {
			t.Fatalf("length: got = %d, wanted = 1", len(got))
		}
		if fetchCount != 1 {
			t.Fatalf("fetch count: got = %d, wanted = 1", fetchCount)
		}
	})
}

func TestPaginate(t *testing.T) {
	ctx := slogtest.Context(t)

	// fakeListExemplars simulates a paginated list RPC backed by a slice.
	fakeListExemplars := func(all []*pb.Exemplar) func(context.Context, *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
		return func(_ context.Context, r *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			offset := 0
			if r.GetPageToken() != "" {
				fmt.Sscanf(r.GetPageToken(), "%d", &offset)
			}
			end := offset + int(r.GetPageSize())
			if end > len(all) {
				end = len(all)
			}
			nextToken := ""
			if end < len(all) {
				nextToken = fmt.Sprintf("%d", end)
			}
			return all[offset:end], nextToken, nil
		}
	}

	t.Run("paginates across multiple pages", func(t *testing.T) {
		all := make([]*pb.Exemplar, 5)
		for i := range all {
			all[i] = &pb.Exemplar{Uid: fmt.Sprintf("uid-%d", rand.Int64()), Name: fmt.Sprintf("name-%d", i)}
		}
		got, err := iter.All(iter.Paginate(ctx, &pb.ListExemplarsRequest{PageSize: 2}, "exemplars", fakeListExemplars(all)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != len(all) {
			t.Fatalf("length: got = %d, wanted = %d", len(got), len(all))
		}
		for i := range got {
			if got[i].GetUid() != all[i].GetUid() {
				t.Errorf("item[%d].uid: got = %q, wanted = %q", i, got[i].GetUid(), all[i].GetUid())
			}
		}
	})

	t.Run("nil request uses defaults", func(t *testing.T) {
		all := make([]*pb.Exemplar, 3)
		for i := range all {
			all[i] = &pb.Exemplar{Uid: fmt.Sprintf("uid-%d", rand.Int64())}
		}
		var receivedPageSize int32
		seq := iter.Paginate(ctx, (*pb.ListExemplarsRequest)(nil), "exemplars", func(_ context.Context, r *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			receivedPageSize = r.GetPageSize()
			return all, "", nil
		})
		got, err := iter.All(seq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != len(all) {
			t.Fatalf("length: got = %d, wanted = %d", len(got), len(all))
		}
		if receivedPageSize != iter.DefaultPageSize {
			t.Errorf("page size: got = %d, wanted = %d", receivedPageSize, iter.DefaultPageSize)
		}
	})

	t.Run("does not mutate caller request", func(t *testing.T) {
		req := &pb.ListExemplarsRequest{Parent: fmt.Sprintf("parent-%d", rand.Int64()), PageSize: 2}
		origParent := req.GetParent()
		seq := iter.Paginate(ctx, req, "exemplars", func(_ context.Context, _ *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			return []*pb.Exemplar{{Uid: "1"}}, "", nil
		})
		if _, err := iter.All(seq); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.GetParent() != origParent {
			t.Errorf("parent mutated: got = %q, wanted = %q", req.GetParent(), origParent)
		}
		if req.GetPageToken() != "" {
			t.Errorf("page_token mutated: got = %q, wanted = %q", req.GetPageToken(), "")
		}
	})

	t.Run("defaults page size when zero", func(t *testing.T) {
		var receivedPageSize int32
		seq := iter.Paginate(ctx, &pb.ListExemplarsRequest{}, "exemplars", func(_ context.Context, r *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			receivedPageSize = r.GetPageSize()
			return nil, "", nil
		})
		if _, err := iter.All(seq); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedPageSize != iter.DefaultPageSize {
			t.Errorf("page size: got = %d, wanted = %d", receivedPageSize, iter.DefaultPageSize)
		}
	})

	t.Run("preserves custom page size", func(t *testing.T) {
		var receivedPageSize int32
		wantSize := int32(10)
		seq := iter.Paginate(ctx, &pb.ListExemplarsRequest{PageSize: wantSize}, "exemplars", func(_ context.Context, r *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			receivedPageSize = r.GetPageSize()
			return nil, "", nil
		})
		if _, err := iter.All(seq); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedPageSize != wantSize {
			t.Errorf("page size: got = %d, wanted = %d", receivedPageSize, wantSize)
		}
	})

	t.Run("fetch error propagates", func(t *testing.T) {
		wantErr := fmt.Sprintf("rpc-error-%d", rand.Int64())
		seq := iter.Paginate(ctx, &pb.ListExemplarsRequest{}, "exemplars", func(_ context.Context, _ *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			return nil, "", errors.New(wantErr)
		})
		_, err := iter.All(seq)
		if err == nil {
			t.Fatalf("error: got = nil, wanted = %q", wantErr)
		}
		want := fmt.Sprintf(`listing "exemplars": %s`, wantErr)
		if err.Error() != want {
			t.Fatalf("error: got = %q, wanted = %q", err.Error(), want)
		}
	})

	t.Run("context cancellation stops pagination", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		fetchCount := 0
		seq := iter.Paginate(ctx, &pb.ListExemplarsRequest{PageSize: 1}, "exemplars", func(_ context.Context, _ *pb.ListExemplarsRequest) ([]*pb.Exemplar, string, error) {
			fetchCount++
			cancel()
			return []*pb.Exemplar{{Uid: "1"}}, "next", nil
		})
		_, err := iter.All(seq)
		if err == nil {
			t.Fatal("error: got = nil, wanted context.Canceled")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error: got = %v, wanted context.Canceled", err)
		}
		if fetchCount != 1 {
			t.Errorf("fetch count: got = %d, wanted = 1", fetchCount)
		}
	})
}
