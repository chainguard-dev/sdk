/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package images

import (
	"fmt"
	"strings"
)

// token represents either a literal string or a RefField reference.
type token interface {
	token()
}

// literal is a literal string segment in a token list.
type literal string

func (literal) token() {}

// RefField represents a field of an image reference (registry, repo, tag, etc.).
type RefField string

func (RefField) token() {}

const (
	Registry     RefField = "registry"
	Repo         RefField = "repo"
	RegistryRepo RefField = "registry_repo"
	Tag          RefField = "tag"
	Digest       RefField = "digest"
	PseudoTag    RefField = "pseudo_tag"
	Ref          RefField = "ref"
)

var knownFields = map[string]RefField{
	"registry":      Registry,
	"repo":          Repo,
	"registry_repo": RegistryRepo,
	"tag":           Tag,
	"digest":        Digest,
	"pseudo_tag":    PseudoTag,
	"ref":           Ref,
}

// TokenList is a sequence of tokens from lexing a value string.
type TokenList []token

// Map applies fn to each RefField, preserves literals, returns parts.
func (t TokenList) Map(fn func(RefField) any) []any {
	parts := make([]any, len(t))
	for i, tok := range t {
		switch v := tok.(type) {
		case literal:
			parts[i] = string(v)
		case RefField:
			parts[i] = fn(v)
		}
	}
	return parts
}

func lex(s string) (TokenList, error) {
	var tokens TokenList
	for len(s) > 0 {
		idx := strings.Index(s, "${")
		if idx == -1 {
			if s != "" {
				tokens = append(tokens, literal(s))
			}
			break
		}
		if idx > 0 {
			tokens = append(tokens, literal(s[:idx]))
		}
		s = s[idx+2:]
		end := strings.Index(s, "}")
		if end == -1 {
			return nil, fmt.Errorf("unclosed variable marker")
		}
		varName := s[:end]
		v, ok := knownFields[varName]
		if !ok {
			return nil, fmt.Errorf("unknown field %q", varName)
		}
		tokens = append(tokens, v)
		s = s[end+1:]
	}
	return tokens, nil
}
