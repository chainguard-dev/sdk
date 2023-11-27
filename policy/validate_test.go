/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package policy

import (
	"context"
	"errors"
	"testing"

	policycontrollerconfig "github.com/sigstore/policy-controller/pkg/config"
	"knative.dev/pkg/apis"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name                  string
		doc                   string
		wantWarns             error
		wantErr               error
		allowEmptyAuthorities bool
	}{{
		name: "good single object",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
  authorities:
  - keyless:
      identities:
      -  issuer: https://issuer.example.com
         subject: foo@example.com
      url: https://fulcio.sigstore.dev
`,
		wantErr: nil,
	}, {
		name: "good CIP and Secret",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
  authorities:
  - keyless:
      identities:
      -  issuer: https://issuer.example.com
         subject: foo@example.com
      url: https://fulcio.sigstore.dev
---
apiVersion: v1
kind: Secret
metadata:
  name: foo
  namespace: cosign-system
stringData:
  foo: bar
`,
		wantErr: nil,
	}, {
		name: "bad secret namespace",
		doc: `
apiVersion: v1
kind: Secret
metadata:
  name: foo
  namespace: something-system
stringData:
  foo: bar
`,
		wantErr: errors.New(`invalid value: something-system: [0].metadata.namespace`),
	}, {
		name: "bad image policy",
		doc: `
apiVersion: policy.sigstore.dev/v1alpha1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
  authorities:
  - key: {}
`,
		wantErr: apis.ErrMissingOneOf("data", "kms", "secretref").ViaField("key").ViaFieldIndex("authorities", 0).ViaField("spec"),
	}, {
		name:    "empty document",
		doc:     ``,
		wantErr: ErrEmptyDocument,
	}, {
		name: "object missing kind",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
# Missing: kind: ClusterImagePolicy
metadata:
  name: blah
spec: {}
`,
		wantErr: errors.New(`decoding object[0]: error unmarshaling JSON: while decoding JSON: Object 'Kind' is missing in '{"apiVersion":"policy.sigstore.dev/v1beta1","metadata":{"name":"blah"},"spec":{}}'`),
	}, {
		name: "unknown field",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  asdf: dfsadf
`,
		wantErr: errors.New(`unable to unmarshal: json: unknown field "asdf"`),
	}, {
		name: "unknown type",
		doc: `
apiVersion: unknown.dev/v1
kind: OtherPolicy
metadata:
  name: blah
spec: {}
`,
		wantErr: errors.New(`unknown type: unknown.dev/v1, Kind=OtherPolicy`),
	}, {
		name: "error - missing field",
		doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
  authorities:
  - keyless:
      url: https://fulcio.sigstore.dev
`,
		wantWarns: nil,
		wantErr:   errors.New("missing field(s): spec.authorities[0].keyless.identities"),
	},
		{
			name: "admit - missing authorities",
			doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
`,
			wantErr:               nil,
			allowEmptyAuthorities: true,
		}, {
			name: "deny - missing authorities",
			doc: `
apiVersion: policy.sigstore.dev/v1beta1
kind: ClusterImagePolicy
metadata:
  name: blah
spec:
  images:
  - glob: '*'
`,
			wantErr: errors.New("missing field(s): spec.authorities"),
		}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testContext := context.Background()
			if test.allowEmptyAuthorities {
				testContext = policycontrollerconfig.ToContext(testContext, &policycontrollerconfig.PolicyControllerConfig{FailOnEmptyAuthorities: false})
			}
			gotWarns, gotErr := Validate(testContext, test.doc)
			if (gotErr != nil) != (test.wantErr != nil) {
				t.Fatalf("Parse() = %v, wanted %v", gotErr, test.wantErr)
			}
			if (gotWarns != nil) != (test.wantWarns != nil) {
				t.Fatalf("Parse() = %v, wanted %v", gotWarns, test.wantWarns)
			}
			if gotErr != nil && gotErr.Error() != test.wantErr.Error() {
				t.Fatalf("Parse() = %v, wanted %v", gotErr, test.wantErr)
			}
			if gotWarns != nil && gotWarns.Error() != test.wantWarns.Error() {
				t.Fatalf("Parse() = %v, wanted %v", gotWarns, test.wantWarns)
			}
		})
	}
}
