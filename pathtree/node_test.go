/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package pathtree

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type formattedWant struct {
	po   PrintOption
	want string
}

func (tw formattedWant) name(name string) string {
	return fmt.Sprintf("%s [%s]", name, tw.po)
}

const noKidsKey string = `root`
const noKidsValue string = `[the root]`
const noKidsKeyValue string = `root [the root]`
const noKidsValueLabel string = `[the root] node root`
const pathOfDescendantsKey string = `
root
└ a
  └ b
    └ c`
const pathOfDescendantsValue string = `
[the root]
└ [z]
  └ [y]
    └ [x]`
const pathOfDescendantsLabel string = `
node root
└ child
  └ grandchild
    └ great grandchild`
const pathOfDescendantsKeyValue string = `
root [the root]
└ a [z]
  └ b [y]
    └ c [x]`
const pathOfDescendantsValueLabel string = `
[the root] node root
└ [z] child
  └ [y] grandchild
    └ [x] great grandchild`
const multipleKidsKeyValueLabel string = `
root [the root] node root
├ a [z] first child
├ b [y] child two
└ c [x] end child`
const multipleKidsKey string = `
root
├ a
├ b
└ c`
const multipleKidsValue string = `
[the root]
├ [x]
├ [y]
└ [z]`
const multipleKidsLabel string = `
node root
├ child two
├ end child
└ first child`
const multipleKidsKeyValue string = `
root [the root]
├ a [z]
├ b [y]
└ c [x]`
const multipleKidsValueLabel string = `
[the root] node root
├ [x] end child
├ [y] child two
└ [z] first child`

func Test_node_print(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want []formattedWant
	}{{
		name: "no children",
		node: newNode("root", "the root", "node root"),
		want: []formattedWant{{
			po:   Key,
			want: noKidsKey,
		}, {
			po:   Value,
			want: noKidsValue,
		}, {
			po:   KeyValue,
			want: noKidsKeyValue,
		}, {
			po:   ValueLabel,
			want: noKidsValueLabel,
		}},
	}, {
		name: "path of descendants",
		node: &Node{
			Key:   "root",
			Value: "the root",
			Label: "node root",
			children: map[string]*Node{
				"a": {
					Key:   "a",
					Value: "z",
					Label: "child",
					children: map[string]*Node{
						"b": {
							Key:   "b",
							Value: "y",
							Label: "grandchild",
							children: map[string]*Node{
								"c": {
									Key:   "c",
									Value: "x",
									Label: "great grandchild",
								},
							},
						},
					},
				},
			},
		},
		want: []formattedWant{{
			po:   Key,
			want: pathOfDescendantsKey,
		}, {
			po:   Value,
			want: pathOfDescendantsValue,
		}, {
			po:   Label,
			want: pathOfDescendantsLabel,
		}, {
			po:   KeyValue,
			want: pathOfDescendantsKeyValue,
		}, {
			po:   ValueLabel,
			want: pathOfDescendantsValueLabel,
		}, {
			// Intentionally invalid PrintOption Value
			po:   1,
			want: pathOfDescendantsKey,
		}},
	}, {
		name: "multiple kids",
		node: &Node{
			Key:   "root",
			Value: "the root",
			Label: "node root",
			children: map[string]*Node{
				"a": {
					Key:   "a",
					Value: "z",
					Label: "first child",
				},
				"b": {
					Key:   "b",
					Value: "y",
					Label: "child two",
				},
				"c": {
					Key:   "c",
					Value: "x",
					Label: "end child",
				},
			},
		},
		want: []formattedWant{{
			po:   KeyValueLabel,
			want: multipleKidsKeyValueLabel,
		}, {
			po:   Key,
			want: multipleKidsKey,
		}, {
			po:   Value,
			want: multipleKidsValue,
		}, {
			po:   Label,
			want: multipleKidsLabel,
		}, {
			po:   KeyValue,
			want: multipleKidsKeyValue,
		}, {
			po:   ValueLabel,
			want: multipleKidsValueLabel,
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tw := range tt.want {
				t.Run(tw.name(tt.name), func(t *testing.T) {
					out := &bytes.Buffer{}
					tt.node.print(out, "", tw.po)
					if diff := cmp.Diff(strings.TrimSpace(tw.want), strings.TrimSpace(out.String())); diff != "" {
						t.Errorf("diff (-want / +got):\n%v", diff)
					}
				})
			}
		})
	}
}

func Test_node_flatten(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want []*Node
	}{{
		name: "nil node",
		node: nil,
		want: []*Node{},
	}, {
		name: "no kids",
		node: &Node{Key: "a", Value: "a"},
		want: []*Node{
			{Key: "a", Value: "a"},
		},
	}, {
		name: "one gen",
		node: &Node{
			Key:   "a",
			Value: "a",
			children: map[string]*Node{
				"b": {Key: "b", Value: "b"},
				"c": {Key: "c", Value: "c"},
			},
		},
		want: []*Node{
			{Key: "a", Value: "a"},
			{Key: "b", Value: "b"},
			{Key: "c", Value: "c"},
		},
	}, {
		name: "multi gen",
		node: &Node{
			Key:   "a",
			Value: "a",
			children: map[string]*Node{
				"b": {
					Key: "b", Value: "b",
					children: map[string]*Node{
						"c": {Key: "c", Value: "c"},
					},
				},
				"d": {Key: "d", Value: "d"},
			},
		},
		want: []*Node{
			{Key: "a", Value: "a"},
			{Key: "b", Value: "b"},
			{Key: "c", Value: "c"},
			{Key: "d", Value: "d"},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			po := KeyValue
			flat := tt.node.flatten(po)

			if diff := cmp.Diff(tt.want, flat, cmp.AllowUnexported(Node{})); diff != "" {
				t.Errorf("diff (-want / +got):\n%v", diff)
			}
		})
	}
}
