/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sigstore/policy-controller/pkg/apis/glob"
	"github.com/sigstore/policy-controller/pkg/apis/policy/v1alpha1"
	"github.com/sigstore/policy-controller/pkg/apis/policy/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/yaml"
)

// Parse decodes a provided YAML document containing zero or more objects into
// a collection of unstructured.Unstructured objects.
func Parse(_ context.Context, document string) ([]*unstructured.Unstructured, error) {
	docs := strings.Split(document, "\n---\n")

	objs := make([]*unstructured.Unstructured, 0, len(docs))
	for i, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(doc), &obj); err != nil {
			return nil, fmt.Errorf("decoding object[%d]: %w", i, err)
		}
		if obj.GetAPIVersion() == "" {
			return nil, apis.ErrMissingField("apiVersion").ViaIndex(i)
		}
		if obj.GetName() == "" {
			return nil, apis.ErrMissingField("metadata.name").ViaIndex(i)
		}
		objs = append(objs, &obj)
	}
	return objs, nil
}

// ParseClusterImagePolicies returns ClusterImagePolicy objects found in the
// policy document.
func ParseClusterImagePolicies(ctx context.Context, document string) (cips []*v1beta1.ClusterImagePolicy, warns error, err error) {
	if warns, err = Validate(ctx, document); err != nil {
		return nil, warns, err
	}

	cips, err = parseClusterImagePolicies(ctx, document)
	if err != nil {
		return nil, warns, err
	}

	return cips, warns, nil
}

// UnsafeParseClusterImagePolicies returns ClusterImagePolicy objects found in the
// policy document, without validating if the policy is valid.
func UnsafeParseClusterImagePolicies(ctx context.Context, document string) (cips []*v1beta1.ClusterImagePolicy, err error) {
	return parseClusterImagePolicies(ctx, document)
}

// CoversImage parses the given policy document, and checks if the target image
// matches any of the image globs included in the policy.
func CoversImage(ctx context.Context, document, target string) (bool, error) {
	cips, _, err := ParseClusterImagePolicies(ctx, document)
	if err != nil {
		return false, err
	}
	if len(cips) != 1 {
		return false, fmt.Errorf("document must contain exactly one ClusterImagePolicy (%d found)", len(cips))
	}

	for _, image := range cips[0].Spec.Images {
		ok, err := glob.Match(image.Glob, target)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func parseClusterImagePolicies(ctx context.Context, document string) (cips []*v1beta1.ClusterImagePolicy, err error) {
	ol, err := Parse(ctx, document)
	if err != nil {
		return nil, err
	}

	cips = make([]*v1beta1.ClusterImagePolicy, 0)
	for _, obj := range ol {
		gv, err := schema.ParseGroupVersion(obj.GetAPIVersion())
		if err != nil {
			// Practically unstructured.Unstructured won't let this happen.
			return nil, fmt.Errorf("error parsing apiVersion of: %w", err)
		}

		cip := &v1beta1.ClusterImagePolicy{}

		switch gv.WithKind(obj.GetKind()) {
		case v1alpha1.SchemeGroupVersion.WithKind("ClusterImagePolicy"):
			v1a1 := &v1alpha1.ClusterImagePolicy{}
			if err := convert(obj, v1a1); err != nil {
				return nil, err
			}
			if err := v1a1.ConvertTo(ctx, cip); err != nil {
				return nil, err
			}

		case v1beta1.SchemeGroupVersion.WithKind("ClusterImagePolicy"):
			// This is allowed, but we should convert things.
			if err := convert(obj, cip); err != nil {
				return nil, err
			}

		default:
			continue
		}

		cips = append(cips, cip)
	}
	return cips, nil
}

func convert(from interface{}, to runtime.Object) error {
	bs, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Marshal() = %w", err)
	}
	if err := json.Unmarshal(bs, to); err != nil {
		return fmt.Errorf("Unmarshal() = %w", err)
	}
	return nil
}
