/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package pathtree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

func splitKey(key string, sep string) []string {
	parts := make([]string, 0, strings.Count(key, sep)+1)
	for part := range strings.SplitSeq(key, sep) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// PrintOption determines what information about each Node is printed.
// The values of Key, Value, and Label were chosen to be relatively prime
// to make determining what combination of identifiers to print purely mathematical
// Actual output may vary depending on available identifiers for each Node.
type PrintOption int

const (
	Key           PrintOption = 2
	Value         PrintOption = 3
	Label         PrintOption = 5
	KeyValue      PrintOption = 6
	KeyLabel      PrintOption = 10
	ValueLabel    PrintOption = 15
	KeyValueLabel PrintOption = 30
)

func (po PrintOption) String() string {
	switch po {
	case Key:
		return "Key"
	case Value:
		return "Value"
	case Label:
		return "Label"
	case KeyValue:
		return "KeyValue"
	case KeyLabel:
		return "KeyLabel"
	case ValueLabel:
		return "ValueLabel"
	case KeyValueLabel:
		return "KeyValueLabel"
	default:
		return "<unknown>"
	}
}

type Tree struct {
	PrintOption PrintOption
	root        *Node
	keySep      string // TODO: make separator configurable
	count       int
	prefix      string
}

// New returns a pointer to an initialized Tree
func New() *Tree {
	return &Tree{
		PrintOption: ValueLabel,
		root:        newNode("/", "/", ""),
		keySep:      "/",
	}
}

// Add adds the given Value to the tree at the given Key path. If interior nodes between the root
// and leaf do not exist, empty nodes without a Value are created along the path. Add cannot be used to
// overwrite a Value that already exists at the given Key path.
func (t *Tree) Add(key, value, label string) error {
	if t == nil || t.root == nil {
		return errors.New("tree is uninitialized; call New() to create empty Tree")
	}

	parts := splitKey(key, t.keySep)
	if len(parts) == 0 {
		return errors.New("key path cannot be empty or only /")
	}

	n := t.root
	path := ""
	for _, k := range parts {
		fullKey := strings.TrimLeft(fmt.Sprintf("%s/%s", path, k), "/") // Remove a leading / if present
		next, ok := n.children[k]
		// add nodes along path if needed
		if !ok {
			next = newNode(fullKey, "", "")
			n.children[k] = next
			t.count++
		}
		n = next
		path = fullKey
	}
	// We've reached the leaf, if there is a Value it must have already existed
	if n.Value != "" {
		return fmt.Errorf("value \"%s\" already exists at key \"%s\"", value, key)
	}
	// Remove newlines from the value and label, they mess with output
	n.Value = strings.ReplaceAll(value, "\n", " ")
	n.Label = strings.ReplaceAll(label, "\n", "")
	t.count++
	return nil
}

// Contains returns true if Key is found in the tree, even if the Value is unset.
func (t *Tree) Contains(key string) bool {
	if t == nil || t.root == nil {
		return false
	}

	parts := splitKey(key, t.keySep)
	n := t.root
	for _, k := range parts {
		next, ok := n.children[k]
		if !ok {
			return false
		}
		n = next
	}

	// TODO: return false if the Value is unset? (i.e. this was an intermediate Node created along a path with no Value)
	return true
}

// Delete deletes the Node identified by the given Key. Currently, only deleting leaves is supported.
func (t *Tree) Delete(key string) (string, error) {
	if t == nil || t.root == nil {
		return "", errors.New("tree is uninitialized; call New() to create empty Tree")
	}

	parts := splitKey(key, t.keySep)
	n := t.root
	p := t.root
	for _, k := range parts {
		next, ok := n.children[k]
		if !ok {
			return "", fmt.Errorf("key \"%s\" does not exist in the tree", key)
		}
		p = n
		n = next
	}

	if len(n.children) > 0 {
		return "", fmt.Errorf("cannot delete key \"%s\" because it has %d children", key, len(n.children))
	}
	delete(p.children, parts[len(parts)-1])
	t.count--
	return n.Value, nil
}

func (t *Tree) Flatten() []*Node {
	if t == nil || t.root == nil {
		return []*Node{}
	}

	nodes := make([]*Node, 0, t.count)
	roots := sortNodes(t.root.children, t.PrintOption)
	for _, root := range roots {
		nodes = append(nodes, root.flatten(t.PrintOption)...)
	}
	return nodes
}

func (t *Tree) SetPrefix(prefix string) {
	t.prefix = prefix
}

// Fprint prints the tree to the writer.
func (t *Tree) Fprint(w io.Writer) {
	if t == nil || t.root == nil {
		return
	}

	// Check for invalid PrintOption Value; default to Key
	po := t.PrintOption
	if po%Key != 0 && po%Value != 0 && po%Label != 0 {
		po = Key
	}

	roots := sortNodes(t.root.children, po)
	for _, r := range roots {
		r.print(w, t.prefix, po)
	}
}

// String returns the tree as a string.
func (t *Tree) String() string {
	b := &bytes.Buffer{}
	t.Fprint(b)
	return b.String()
}

// Get returns the Value in the Tree stored at the given Key path. If no Value was stored but the
// Key exists, return the local portion of the Key path.
func (t *Tree) Get(key string) (string, error) {
	if t == nil || t.root == nil {
		return "", errors.New("tree is uninitialized; call New() to create empty Tree")
	}

	parts := splitKey(key, t.keySep)
	n := t.root
	for _, k := range parts {
		next, ok := n.children[k]
		if !ok {
			return "", fmt.Errorf("key \"%s\" not found in \"%s\"", k, n.Key)
		}
		n = next
	}

	// TODO: If no Value is saved at this location, should we return an error instead?
	if n.Value == "" {
		return n.Key, nil
	}
	return n.Value, nil
}

// Update updates the Node at the given Key location, if it exists. Returns an error if any ancestor
// Node does not exist (since that implies the Node at Key does not exist yet).
func (t *Tree) Update(key, value, label string) error {
	if t == nil || t.root == nil {
		return errors.New("tree is uninitialized; call New() to create empty Tree")
	}

	parts := splitKey(key, t.keySep)
	if len(parts) == 0 {
		return errors.New("key path cannot be empty or only /")
	}

	n := t.root
	for _, k := range parts {
		next, ok := n.children[k]
		// Any missing nodes in the path imply this Key hasn't been added yet
		if !ok {
			return fmt.Errorf("key \"%s\" does not exist in the tree", key)
		}
		n = next
	}

	n.Value = value
	n.Label = label
	return nil
}
