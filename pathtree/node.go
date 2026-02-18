/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package pathtree

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Node struct {
	Key      string
	Value    string
	Label    string
	children map[string]*Node
}

func newNode(key, value, label string) *Node {
	return &Node{
		Key:      key,
		Value:    value,
		Label:    label,
		children: make(map[string]*Node, 2), // TODO: configurable expected number of children?
	}
}

func (n *Node) flatten(po PrintOption) []*Node {
	if n == nil {
		return []*Node{}
	}
	nodes := make([]*Node, 0, len(n.children)+1) // TODO: choose the cap more intelligently
	nodes = append(nodes, shallow(n))
	kids := sortNodes(n.children, po)

	for _, kid := range kids {
		nodes = append(nodes, kid.flatten(po)...)
	}
	return nodes
}

func (n *Node) formatNodeInfo(po PrintOption) string {
	// Current output format: "Key [Value] Label"
	// Collect relevant info for this Node, based on po and what exists
	var info string
	key := n.Key[strings.LastIndex(n.Key, "/")+1:]
	if po%Key == 0 {
		info = key + " "
	}
	if po%Value == 0 && n.Value != "" {
		info = fmt.Sprintf("%s[%s] ", info, n.Value)
	}
	if po%Label == 0 && n.Label != "" {
		info = fmt.Sprintf("%s%s", info, n.Label)
	}

	// In case of invalid PrintOption or missing Value/Label info
	if info == "" {
		info = key
	}

	return strings.TrimSpace(info)
}

func (n *Node) print(w io.Writer, prefix string, po PrintOption) {
	info := n.formatNodeInfo(po)

	// Print this Node
	_, _ = fmt.Fprintf(w, "%s%s\n", prefix, info)

	// Adjust prefix for the next generation line(s), if necessary
	// "├ " is replaced with "│ "
	// "└ " is replaced with "  "
	if p, ok := strings.CutSuffix(prefix, teeSpace); ok {
		prefix = p + pipeSpace
	} else if p, ok := strings.CutSuffix(prefix, elbowSpace); ok {
		prefix = p + spaceSpace
	}

	children := sortNodes(n.children, po)
	for i, child := range children {
		if i == len(children)-1 {
			// Last child
			child.print(w, prefix+elbowSpace, po)
		} else {
			// Not the last children
			child.print(w, prefix+teeSpace, po)
		}
	}
}

// lessThan this *Node is less than the other *Node, based on the given PrintOption.
// Composite types are evaluated in the order listed in their name.
// For example: KeyValue is evaluated by Key, then Value.
func (n *Node) lessThan(other *Node, po PrintOption) bool {
	switch po {
	case Key:
		return n.Key < other.Key
	case Value:
		// If both Values are empty, order by Key
		if n.Value == "" && other.Value == "" {
			return n.lessThan(other, Key)
		}
		return n.Value < other.Value
	case Label:
		// If both Labels are empty, order by Key
		if n.Label == "" && other.Label == "" {
			return n.lessThan(other, Key)
		}
		return n.Label < other.Label
	case KeyValue:
		if n.Key != other.Key {
			return n.lessThan(other, Key)
		}
		return n.lessThan(other, Value)
	case KeyLabel:
		if n.Key != other.Key {
			return n.lessThan(other, Key)
		}
		return n.lessThan(other, Label)
	case ValueLabel:
		if n.Value != other.Value {
			return n.lessThan(other, Value)
		}
		return n.lessThan(other, Label)
	case KeyValueLabel:
		if n.Key != other.Key {
			return n.lessThan(other, Key)
		}
		if n.Value != other.Value {
			return n.lessThan(other, Value)
		}
		return n.lessThan(other, Label)
	default:
		// For invalid values of po, consider all fields.
		return n.lessThan(other, KeyValueLabel)
	}
}

// shallow creates a shallow (childless) copy of n
func shallow(n *Node) *Node {
	return &Node{
		Key:   n.Key,
		Value: n.Value,
		Label: n.Label,
	}
}

// Keeping sortedKeys around for now in case it is useful
// nolint
func sortedKeys(m map[string]*Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortNodes(m map[string]*Node, po PrintOption) []*Node {
	nodes := make([]*Node, 0, len(m))
	for k := range m {
		nodes = append(nodes, m[k])
	}

	// Sort based on the given PrintOption. Composite types are sorted in order listed.
	// For example:
	//	ValueLabel will sort by Value -> Label.
	//	KeyValue will sort by Key -> Value
	//  etc...
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].lessThan(nodes[j], po)
	})
	return nodes
}
