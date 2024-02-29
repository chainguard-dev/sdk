/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	admissionv1 "k8s.io/api/admission/v1"
)

// ReviewEventType is the cloudevents event type when an admission controller
// is invoked.
const ReviewEventType = "dev.chainguard.admission.v1"

// ReviewBody is the body of the Chainguard event Occurrence when the event type
// is ReviewEventType.
type ReviewBody admissionv1.AdmissionReview
