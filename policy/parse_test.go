/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package policy

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"knative.dev/pkg/apis"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		doc     string
		want    []*unstructured.Unstructured
		wantErr error
	}{{
		name: "good single object",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec: {}
`,
		want: []*unstructured.Unstructured{{
			Object: map[string]interface{}{
				"apiVersion": "policy.sigstore.dev/v1beta1",
				"kind":       "ClusterImagePolicy",
				"metadata": map[string]interface{}{
					"name": "blah",
				},
				"spec": map[string]interface{}{},
			},
		}},
	}, {
		name: "good multi-object",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec: {}
---
---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: foo
spec: {}
---

---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: bar
spec: {}
`,
		want: []*unstructured.Unstructured{{
			Object: map[string]interface{}{
				"apiVersion": "policy.sigstore.dev/v1beta1",
				"kind":       "ClusterImagePolicy",
				"metadata": map[string]interface{}{
					"name": "blah",
				},
				"spec": map[string]interface{}{},
			},
		}, {
			Object: map[string]interface{}{
				"apiVersion": "policy.sigstore.dev/v1beta1",
				"kind":       "ClusterImagePolicy",
				"metadata": map[string]interface{}{
					"name": "foo",
				},
				"spec": map[string]interface{}{},
			},
		}, {
			Object: map[string]interface{}{
				"apiVersion": "policy.sigstore.dev/v1beta1",
				"kind":       "ClusterImagePolicy",
				"metadata": map[string]interface{}{
					"name": "bar",
				},
				"spec": map[string]interface{}{},
			},
		}},
	}, {
		name: "bad missing apiVersion",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec: {}
---
# Missing: apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: foo
spec: {}
---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: bar
spec: {}
`,
		wantErr: apis.ErrMissingField("[1].apiVersion"),
	}, {
		name: "bad missing kind",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec: {}
---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: foo
spec: {}
---
apiVersion: policy.sigstore.dev/v1beta1
# Missing: kind: ClusterImagePolicy
metadata:
  name: bar
spec: {}
`,
		wantErr: errors.New(`decoding object[2]: error unmarshaling JSON: while decoding JSON: Object 'Kind' is missing in '{"apiVersion":"policy.sigstore.dev/v1beta1","metadata":{"name":"bar"},"spec":{}}'`),
	}, {
		name: "bad missing apiVersion",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  # Missing: name: blah
sp dec: {}
---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: foo
spec: {}
---
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: bar
spec: {}
`,
		wantErr: apis.ErrMissingField("[0].metadata.name"),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Parse(context.Background(), test.doc)

			switch {
			case (gotErr != nil) != (test.wantErr != nil):
				t.Fatalf("Parse() = %v, wanted %v", gotErr, test.wantErr)
			case gotErr != nil && gotErr.Error() != test.wantErr.Error():
				t.Fatalf("Parse() = %v, wanted %v", gotErr, test.wantErr)
			case gotErr != nil:
				return // This was an error test.
			}

			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("Parse (-got, +want) = %s", diff)
			}
		})
	}
}
