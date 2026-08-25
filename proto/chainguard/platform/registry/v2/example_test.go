/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2_test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2"
)

func ExampleNewClientsFromConnection() {
	clients := registry.NewClientsFromConnection(nil)
	fmt.Println(clients != nil)
	// Output: true
}

func ExampleClients_ReposService() {
	clients := registry.NewClientsFromConnection(nil)
	fmt.Println(clients.ReposService() != nil)
	// Output: true
}

func ExampleClients_TagsService() {
	clients := registry.NewClientsFromConnection(nil)
	fmt.Println(clients.TagsService() != nil)
	// Output: true
}

func ExampleClients_ListReposIter() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListReposIter(context.Background(), &registry.ListReposRequest{})
	fmt.Println("iterator created")
	// Output: iterator created
}

func ExampleClients_ListReposAll() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListReposAll
	fmt.Println("func available")
	// Output: func available
}

func ExampleClients_ListCatalogImagesIter() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListCatalogImagesIter(context.Background(), &registry.ListCatalogImagesRequest{})
	fmt.Println("iterator created")
	// Output: iterator created
}

func ExampleClients_ListCatalogImagesAll() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListCatalogImagesAll
	fmt.Println("func available")
	// Output: func available
}

func ExampleClients_ListTagsIter() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListTagsIter(context.Background(), &registry.ListTagsRequest{})
	fmt.Println("iterator created")
	// Output: iterator created
}

func ExampleClients_ListTagsAll() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.ListTagsAll
	fmt.Println("func available")
	// Output: func available
}

func ExampleClients_ImagesService() {
	clients := registry.NewClientsFromConnection(nil)
	fmt.Println(clients.ImagesService() != nil)
	// Output: true
}

func ExampleClients_Close() {
	clients := registry.NewClientsFromConnection(nil)
	_ = clients.Close()
	fmt.Println("closed")
	// Output: closed
}
