/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package apk provides CloudEvents types for APK package registry operations.

# Overview

This package defines event types and data structures for tracking APK package
push and pull operations in the Chainguard registry. These events are emitted
as CloudEvents and can be consumed by event processing systems for auditing,
analytics, and monitoring purposes.

# Features

  - PushEvent: Captures metadata when APK packages are pushed to the registry
  - PullEvent: Captures metadata when APK packages are pulled from the registry
  - Location tracking: Records approximate geographic location of clients
  - Proxy support: Tracks pulls through virtualapk.cgr.dev proxy with UIDP and hash
  - Error reporting: Includes error details for failed push operations

# Event Types

The package defines two CloudEvents event types:

  - dev.chainguard.apk.push.v1: Emitted when an APK is pushed
  - dev.chainguard.apk.pull.v1: Emitted when an APK is pulled

# Usage

Events are typically consumed from a CloudEvents source and unmarshaled into
the appropriate event type:

	import (
		"encoding/json"
		cloudevents "github.com/cloudevents/sdk-go/v2"
		"chainguard.dev/sdk/events/apk"
	)

	func handleEvent(event cloudevents.Event) error {
		switch event.Type() {
		case apk.PushedEventType:
			var push apk.PushEvent
			if err := event.DataAs(&push); err != nil {
				return err
			}
			// Process push event
			log.Printf("Package pushed: %s/%s-%s", push.Architecture, push.Package, push.Version)

		case apk.PulledEventType:
			var pull apk.PullEvent
			if err := event.DataAs(&pull); err != nil {
				return err
			}
			// Process pull event
			log.Printf("Package pulled: %s", pull.APKPath())
		}
		return nil
	}

# Integration Patterns

Events can be consumed from various sources:

  - Cloud Pub/Sub subscriptions
  - HTTP CloudEvents endpoints
  - Event streaming platforms

Example integration with Cloud Pub/Sub:

	import (
		"cloud.google.com/go/pubsub"
		cloudevents "github.com/cloudevents/sdk-go/v2"
		"chainguard.dev/sdk/events/apk"
	)

	func processPubSubMessage(msg *pubsub.Message) error {
		event := cloudevents.NewEvent()
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			return err
		}

		if event.Type() == apk.PushedEventType {
			var push apk.PushEvent
			if err := event.DataAs(&push); err != nil {
				return err
			}
			// Handle push event
		}
		return nil
	}

# Path Construction

Both PushEvent and PullEvent provide convenience methods for constructing
APK file paths:

  - APKPath(): Returns the full path including repository ID
  - APKBasePath(): Returns the base path without repository ID

These methods follow the standard APK repository layout:
{architecture}/{package}-{version}.apk
*/
package apk
