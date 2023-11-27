/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events

const (
	// DeliveryTypeKey is the CloudEvents extension name to filter on for kinds
	// of deliveries.
	DeliveryTypeKey = "chainguarddev1delivery"
	// DeliveryTypeWebhook defines webhook delivery type.
	DeliveryTypeWebhook = "webhook"

	// DeliveryWebhookTargetKey is the CloudEvents extension name to store the
	// target url for webhooks.
	DeliveryWebhookTargetKey = "chainguarddev1webhook"

	// DeliverySubscriptionKey is the CloudEvents extension name to store the
	// subscription id that caused the event.
	DeliverySubscriptionKey = "chainguarddev1subscription"

	// GroupKey is the CloudEvents extension name to store the group associated
	// to the event.
	GroupKey = "group"
	// ClusterKey is the CloudEvents extension name to store the cluster associated
	// to the event.
	ClusterKey = "cluster"
	// ImageKey is the CloudEvents extension name to store the image associated
	// to the event.
	ImageKey = "image"

	// AudienceKey labels an event for its intended audience ["internal", "customer"].
	AudienceKey = "audience"
	// AudienceInternal are events intended for the internal platform.
	AudienceInternal = "internal"
	// AudienceCustomer are events targeting outside the platform.
	AudienceCustomer = "customer"

	// ArrivalTimeKey is the CloudEvents extension name to store the Knative
	// arrival timestamp
	ArrivalTimeKey = "knativearrivaltime"
)
