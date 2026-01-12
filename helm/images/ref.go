/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package images

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// OCIRef holds the components of an OCI reference.
// Use NewRef to construct - it parses and validates the reference string.
type OCIRef struct {
	Registry     string // Registry host, e.g., "cgr.dev"
	Repo         string // Repository path without registry, e.g., "chainguard/nginx"
	Tag          string // OCI tag, e.g., "latest"
	Digest       string // OCI digest, e.g., "sha256:abc123..."
	RegistryRepo string // Combined registry/repo, e.g., "cgr.dev/chainguard/nginx"
	PseudoTag    string // Tag with digest, or "unused@digest" if no tag
	FullRef      string // Full reference, e.g., "cgr.dev/chainguard/nginx:latest@sha256:abc123..."
}

// NewRef parses and validates an OCI reference string.
// NOTE: This is relatively custom because we try to preserve the tag in
// "pseduo_tag", which ggcr trims (for good reason).
func NewRef(reference string) (OCIRef, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return OCIRef{}, err
	}

	registry := ref.Context().RegistryStr()
	repo := ref.Context().RepositoryStr()
	registryRepo := ref.Context().Name()
	var tag, digest string

	if t, ok := ref.(name.Tag); ok {
		tag = t.TagStr()
	}
	if d, ok := ref.(name.Digest); ok {
		digest = d.DigestStr()
		// Recover tag from tag@digest format - ParseReference loses it.
		beforeAt, _, _ := strings.Cut(reference, "@")
		if tagRef, err := name.ParseReference(beforeAt); err == nil {
			if t, ok := tagRef.(name.Tag); ok {
				extracted := t.TagStr()
				// Only use if explicit (not implicit "latest")
				if strings.HasSuffix(beforeAt, ":"+extracted) {
					tag = extracted
				}
			}
		}
	}

	if digest == "" {
		return OCIRef{}, fmt.Errorf("reference %q must include a digest", reference)
	}

	var pseudoTag string
	if tag == "" {
		pseudoTag = "unused@" + digest
	} else {
		pseudoTag = tag + "@" + digest
	}

	fullRef := registryRepo
	if tag != "" {
		fullRef += ":" + tag
	}
	if digest != "" {
		fullRef += "@" + digest
	}

	return OCIRef{
		Registry:     registry,
		Repo:         repo,
		Tag:          tag,
		Digest:       digest,
		RegistryRepo: registryRepo,
		PseudoTag:    pseudoTag,
		FullRef:      fullRef,
	}, nil
}

// Resolve returns a WalkFunc that resolves tokens to strings using refs.
func Resolve(refs map[string]OCIRef) WalkFunc {
	return func(imageID string, tokens TokenList) (any, error) {
		ref, ok := refs[imageID]
		if !ok {
			return nil, fmt.Errorf("no ref found for image %q", imageID)
		}

		var sb strings.Builder
		for _, tok := range tokens {
			switch v := tok.(type) {
			case literal:
				sb.WriteString(string(v))
			case RefField:
				var val string
				switch v {
				case Registry:
					val = ref.Registry
				case Repo:
					val = ref.Repo
				case RegistryRepo:
					val = ref.RegistryRepo
				case Tag:
					val = ref.Tag
				case Digest:
					val = ref.Digest
				case PseudoTag:
					val = ref.PseudoTag
				case Ref:
					val = ref.FullRef
				}
				if val == "" {
					return nil, fmt.Errorf("empty value for field %q: use ${pseudo_tag} for digest-only refs", v)
				}
				sb.WriteString(val)
			}
		}
		return sb.String(), nil
	}
}
