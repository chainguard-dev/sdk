/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package iter_test

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Example demonstrates basic usage of the List iterator with a simple
// in-memory data source.
func Example() {
	ctx := context.Background()

	// Simulate a paginated API with in-memory data
	allItems := []string{"item1", "item2", "item3", "item4", "item5"}
	pageSize := 2

	// Create an iterator using List
	for item, err := range iter.List(ctx, "items", func(pageToken string) ([]string, string, error) {
		// Parse page token to determine offset
		offset := 0
		if pageToken != "" {
			fmt.Sscanf(pageToken, "%d", &offset)
		}

		// Return a page of items
		end := offset + pageSize
		if end > len(allItems) {
			end = len(allItems)
		}

		items := allItems[offset:end]

		// Calculate next page token
		nextToken := ""
		if end < len(allItems) {
			nextToken = fmt.Sprintf("%d", end)
		}

		return items, nextToken, nil
	}) {
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Println(item)
	}

	// Output:
	// item1
	// item2
	// item3
	// item4
	// item5
}

// ExampleList demonstrates using the List iterator with a custom struct type.
func ExampleList() {
	ctx := context.Background()

	type User struct {
		ID   int
		Name string
	}

	// Simulate paginated user data
	allUsers := []User{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}
	pageSize := 2

	for user, err := range iter.List(ctx, "users", func(pageToken string) ([]User, string, error) {
		offset := 0
		if pageToken != "" {
			fmt.Sscanf(pageToken, "%d", &offset)
		}

		end := offset + pageSize
		if end > len(allUsers) {
			end = len(allUsers)
		}

		users := allUsers[offset:end]
		nextToken := ""
		if end < len(allUsers) {
			nextToken = fmt.Sprintf("%d", end)
		}

		return users, nextToken, nil
	}) {
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Printf("User %d: %s\n", user.ID, user.Name)
	}

	// Output:
	// User 1: Alice
	// User 2: Bob
	// User 3: Charlie
}

// ExampleList_withDefaultPageSize demonstrates using the DefaultPageSize constant
// when making paginated API requests.
func ExampleList_withDefaultPageSize() {
	ctx := context.Background()

	type Item struct {
		Value string
	}

	// Simulate an API client that accepts a page size parameter
	fetchPage := func(_ string, _ int) ([]Item, string, error) {
		// In a real implementation, this would call an API
		// For this example, return empty results
		return []Item{}, "", nil
	}

	// Use DefaultPageSize when creating the iterator
	for item, err := range iter.List(ctx, "items", func(pageToken string) ([]Item, string, error) {
		return fetchPage(pageToken, iter.DefaultPageSize)
	}) {
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Println(item.Value)
	}
}

// ExampleList_contextCancellation demonstrates how the iterator respects
// context cancellation.
func ExampleList_contextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel after processing 2 items
	count := 0

	for item, err := range iter.List(ctx, "items", func(pageToken string) ([]string, string, error) {
		// Return items in pages
		if pageToken == "" {
			return []string{"item1", "item2"}, "page2", nil
		}
		return []string{"item3", "item4"}, "", nil
	}) {
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Println(item)
		count++
		if count == 2 {
			cancel()
		}
	}

	// Output:
	// item1
	// item2
	// error: context canceled
}
