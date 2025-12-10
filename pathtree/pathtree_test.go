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

type testPath struct {
	path  string
	value string
	label string
}

func TestTree_Add(t *testing.T) {
	tests := []struct {
		name    string
		paths   []testPath
		wantErr bool
	}{{
		name: "one root",
		paths: []testPath{
			{"/root", "end", "end of the root"},
		},
		wantErr: false,
	}, {
		name: "one root, single long path",
		paths: []testPath{
			{"/root/sub1/sub2", "end", "long path from root"},
		},
		wantErr: false,
	}, {
		name: "one root, multiple long paths (all named)",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub3", "end2", "second path through sub1"},
			{"/root/sub1", "end3", "in the middle"},
		},
		wantErr: false,
	}, {
		name: "one root, multiple long paths (some unnamed)",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub3", "end2", "second path through sub1"},
		},
		wantErr: false,
	}, {
		name: "one root, add existing error",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub2", "end2", "second path through sub1"},
		},
		wantErr: true,
	}, {
		name: "multiple roots",
		paths: []testPath{
			{"/root1/a", "a end", "end of a from root1"},
			{"/root2/a", "a end", "end of a from root2"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := New()
			errors := make([]error, 0, len(tt.paths))
			for _, p := range tt.paths {
				if err := tree.Add(p.path, p.value, p.label); err != nil {
					errors = append(errors, err)
				}
			}

			if len(errors) > 0 && !tt.wantErr || len(errors) == 0 && tt.wantErr {
				t.Errorf("Add() returned error = %v, wantErr = %v", len(errors) > 0, tt.wantErr)
			}
		})
	}
}

func TestTree_Get(t *testing.T) {
	tests := []struct {
		name    string
		paths   []testPath
		search  string
		want    string
		wantErr bool
	}{{
		name: "one root",
		paths: []testPath{
			{"/root", "end", "end of the root"},
		},
		search:  "/root",
		want:    "end",
		wantErr: false,
	}, {
		name: "one root, single long path",
		paths: []testPath{
			{"/root/sub1/sub2", "end", "long path from root"},
		},
		search:  "/root/sub1/sub2",
		want:    "end",
		wantErr: false,
	}, {
		name: "one root, multiple long paths, get middle named",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub3", "end2", "second path through sub1"},
			{"/root/sub1", "end3", "in the middle"},
		},
		search:  "/root/sub1",
		want:    "end3",
		wantErr: false,
	}, {
		name: "one root, multiple long paths, get middle unnamed",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub3", "end2", "second path through sub1"},
		},
		search:  "/root/sub1",
		want:    "root/sub1",
		wantErr: false,
	}, {
		name: "one root, key not found",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "end of path"},
		},
		search:  "/root/sub1/sub3",
		wantErr: true,
	}, {
		name: "multiple roots",
		paths: []testPath{
			{"/root1/sub1/sub2", "end112", "end of path from root1"},
			{"/root2/sub1/sub2", "end212", "end of path from root2"},
		},
		search:  "/root2/sub1/sub2",
		want:    "end212",
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tree := New()
			for _, tp := range tt.paths {
				_ = tree.Add(tp.path, tp.value, tp.label)
			}

			// Act
			got, err := tree.Get(tt.search)

			// Assert
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("diff (-want / +got):\n%v", diff)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error: %v; wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

const oneRootKey string = `root`
const oneRootValue string = `[the root]`
const oneRootKeyValueLabel string = `root [the root] root node`
const oneRootSinglePathKey string = `
root
└ a
  └ b`
const oneRootSinglePathValue string = `
[the root]
└ [y]
  └ [x]`
const oneRootSinglePathKeyValue string = `
root [the root]
└ a [y]
  └ b [x]`
const oneRootSinglePathValueLabel string = `
[the root] root node
└ [y] first child
  └ [x] child two`
const oneRootSinglePathKeyValueLabel string = `
root [the root] root node
└ a [y] first child
  └ b [x] child two`
const oneRootMultipleGrandchildrenKey string = `
root
└ a
  ├ b
  └ c`
const oneRootMultipleGrandchildrenValue string = `
[the root]
└ [z]
  ├ [x]
  └ [y]`
const oneRootMultipleGrandchildrenLabel string = `
root node
└ first child
  ├ child two
  └ end child`
const oneRootMultipleGrandchildrenKeyValue string = `
root [the root]
└ a [z]
  ├ b [y]
  └ c [x]`
const oneRootMultipleGrandchildrenValueLabel string = `
[the root] root node
└ [z] first child
  ├ [x] end child
  └ [y] child two`
const oneRootMultipleGrandchildrenKeyValueLabel string = `
root [the root] root node
└ a [z] first child
  ├ b [y] child two
  └ c [x] end child`
const oneRootMultipleGrandchildrenUnValue string = `
root
└ a
  ├ [x]
  └ [y]`
const oneRootMultipleGrandchildrenUnLabel string = `
root
└ a
  ├ child two
  └ end child`
const oneRootMultipleGrandchildrenUnKeyValue string = `
root
└ a
  ├ b [y]
  └ c [x]`
const oneRootMultipleGrandchildrenUnValueLabel string = `
root
└ a
  ├ [x] end child
  └ [y] child two`
const oneRootMultipleGrandchildrenUnKeyValueLabel string = `
root
└ a
  ├ b [y] child two
  └ c [x] end child`
const multipleRootsKey string = `
root1
└ a
root2
└ a`
const multipleRootsValue string = `
[da root2]
└ [x]
[the root1]
└ [y]`
const multipleRootsLabel string = `
first root
└ first child
second root
└ second child`
const multipleRootsKeyValue string = `
root1 [the root1]
└ a [y]
root2 [da root2]
└ a [x]`
const multipleRootsValueLabel string = `
[da root2] second root
└ [x] second child
[the root1] first root
└ [y] first child`
const multipleRootsKeyValueLabel string = `
root1 [the root1] first root
└ a [y] first child
root2 [da root2] second root
└ a [x] second child`
const jaggedKey string = `
root
├ a
│ └ b
└ c`
const jaggedValue string = `
[the root]
├ [x]
└ [z]
  └ [y]`
const jaggedLabel string = `
root node
├ end child
└ first child
  └ child two`
const jaggedKeyValue string = `
root [the root]
├ a [z]
│ └ b [y]
└ c [x]`
const jaggedValueLabel string = `
[the root] root node
├ [x] end child
└ [z] first child
  └ [y] child two`
const jaggedKeyValueLabel string = `
root [the root] root node
├ a [z] first child
│ └ b [y] child two
└ c [x] end child`

func TestTree_Fprint(t *testing.T) {
	tests := []struct {
		name  string
		paths []testPath
		want  []formattedWant
	}{{
		name: "one root",
		paths: []testPath{
			{"/root", "the root", "root node"},
		},
		want: []formattedWant{{
			po:   Key,
			want: oneRootKey,
		}, {
			po:   Value,
			want: oneRootValue,
		}, {
			po:   KeyValueLabel,
			want: oneRootKeyValueLabel,
		}},
	}, {
		name: "one root, single long path",
		paths: []testPath{
			{"/root", "the root", "root node"},
			{"/root/a", "y", "first child"},
			{"/root/a/b", "x", "child two"},
		},
		want: []formattedWant{{
			po:   Key,
			want: oneRootSinglePathKey,
		}, {
			po:   Value,
			want: oneRootSinglePathValue,
		}, {
			po:   KeyValue,
			want: oneRootSinglePathKeyValue,
		}, {
			po:   ValueLabel,
			want: oneRootSinglePathValueLabel,
		}, {
			po:   KeyValueLabel,
			want: oneRootSinglePathKeyValueLabel,
		}},
	}, {
		name: "one root, multiple grandchildren (all named)",
		paths: []testPath{
			{"/root", "the root", "root node"},
			{"/root/a", "z", "first child"},
			{"/root/a/b", "y", "child two"},
			{"/root/a/c", "x", "end child"},
		},
		want: []formattedWant{{
			po:   Key,
			want: oneRootMultipleGrandchildrenKey,
		}, {
			po:   Value,
			want: oneRootMultipleGrandchildrenValue,
		}, {
			po:   Label,
			want: oneRootMultipleGrandchildrenLabel,
		}, {
			po:   KeyValue,
			want: oneRootMultipleGrandchildrenKeyValue,
		}, {
			po:   ValueLabel,
			want: oneRootMultipleGrandchildrenValueLabel,
		}, {
			po:   KeyValueLabel,
			want: oneRootMultipleGrandchildrenKeyValueLabel,
		}},
	}, {
		name: "one root, multiple grandchildren (some unnamed)",
		paths: []testPath{
			{"/root/a/b", "y", "child two"},
			{"/root/a/c", "x", "end child"},
		},
		want: []formattedWant{{
			po:   Value,
			want: oneRootMultipleGrandchildrenUnValue,
		}, {
			po:   Label,
			want: oneRootMultipleGrandchildrenUnLabel,
		}, {
			po:   KeyValue,
			want: oneRootMultipleGrandchildrenUnKeyValue,
		}, {
			po:   ValueLabel,
			want: oneRootMultipleGrandchildrenUnValueLabel,
		}, {
			po:   KeyValueLabel,
			want: oneRootMultipleGrandchildrenUnKeyValueLabel,
		}},
	}, {
		name: "multiple roots",
		paths: []testPath{
			{"/root1", "the root1", "first root"},
			{"/root2", "da root2", "second root"},
			{"/root1/a", "y", "first child"},
			{"/root2/a", "x", "second child"},
		},
		want: []formattedWant{{
			po:   Key,
			want: multipleRootsKey,
		}, {
			po:   Value,
			want: multipleRootsValue,
		}, {
			po:   Label,
			want: multipleRootsLabel,
		}, {
			po:   KeyValue,
			want: multipleRootsKeyValue,
		}, {
			po:   ValueLabel,
			want: multipleRootsValueLabel,
		}, {
			po:   KeyValueLabel,
			want: multipleRootsKeyValueLabel,
		}},
	}, {
		name: "jagged",
		paths: []testPath{
			{"/root", "the root", "root node"},
			{"/root/a", "z", "first child"},
			{"/root/a/b", "y", "child two"},
			{"/root/c", "x", "end child"},
		},
		want: []formattedWant{{
			po:   Key,
			want: jaggedKey,
		}, {
			po:   Value,
			want: jaggedValue,
		}, {
			po:   Label,
			want: jaggedLabel,
		}, {
			po:   KeyValue,
			want: jaggedKeyValue,
		}, {
			po:   ValueLabel,
			want: jaggedValueLabel,
		}, {
			po:   KeyValueLabel,
			want: jaggedKeyValueLabel,
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tw := range tt.want {
				t.Run(tw.name(tt.name), func(t *testing.T) {
					tree := New()
					for _, p := range tt.paths {
						_ = tree.Add(p.path, p.value, p.label)
					}
					tree.PrintOption = tw.po

					out := &bytes.Buffer{}
					tree.Fprint(out)
					if diff := cmp.Diff(strings.TrimSpace(tw.want), strings.TrimSpace(out.String())); diff != "" {
						t.Errorf("diff (-want / +got):\n%v", diff)
					}
				})
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {
	tests := []struct {
		name    string
		paths   []testPath
		search  string
		want    string
		wantErr bool
	}{{
		name: "one root",
		paths: []testPath{
			{"/root", "end", "end of the root"},
		},
		search:  "/root",
		want:    "end",
		wantErr: false,
	}, {
		name: "one root, single long path",
		paths: []testPath{
			{"/root/sub1/sub2", "end", "long path from root"},
		},
		search:  "/root/sub1/sub2",
		want:    "end",
		wantErr: false,
	}, {
		name: "one root, multiple children",
		paths: []testPath{
			{"/root/a", "z", "long path from root"},
			{"/root/b", "y", "long path from root"},
			{"/root/c", "x", "long path from root"},
		},
		search:  "/root/b",
		want:    "y",
		wantErr: false,
	}, {
		name: "one root, multiple long paths, delete middle",
		paths: []testPath{
			{"/root/sub1/sub2", "end1", "first path through sub1"},
			{"/root/sub1/sub3", "end2", "second path through sub1"},
			{"/root/sub1", "end3", "in the middle"},
		},
		search:  "/root/sub1",
		wantErr: true,
	}, {
		name: "multiple roots, delete root",
		paths: []testPath{
			{"/root1/sub1", "end", "end of path from root1"},
			{"/root2", "second root", "end of path from root2"},
		},
		search:  "/root2",
		want:    "second root",
		wantErr: false,
	}, {
		name: "leaf not found",
		paths: []testPath{
			{"/root/a/b", "end", ""},
		},
		search:  "/root/a/c",
		wantErr: true,
	}, {
		name: "branch not found",
		paths: []testPath{
			{"/root/a/b", "end", ""},
		},
		search:  "/root/c/b",
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tree := New()
			for _, tp := range tt.paths {
				_ = tree.Add(tp.path, tp.value, tp.label)
			}

			// Act
			got, err := tree.Delete(tt.search)

			// Assert
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("diff (-want / +got):\n%v", diff)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error: %v; wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestTree_Update(t *testing.T) {
	tests := []struct {
		name    string
		paths   []testPath
		update  testPath
		wantErr bool
	}{{
		name: "one root",
		paths: []testPath{
			{"/root", "end", "end of the root"},
		},
		update:  testPath{"/root", "new end", "new label"},
		wantErr: false,
	}, {
		name: "one root, single long path",
		paths: []testPath{
			{"/root/a/b", "x", "label for x"},
		},
		update:  testPath{"/root/a/b", "new x", "new label"},
		wantErr: false,
	}, {
		name: "one root, multiple children",
		paths: []testPath{
			{"/root/a", "x", "label for x"},
			{"/root/b", "y", "label for y"},
			{"/root/c", "z", "label for y"},
		},
		update:  testPath{"/root/b", "new y", "new label"},
		wantErr: false,
	}, {
		name: "one root, duplicate leaf key part",
		paths: []testPath{
			{"/root/a/c", "x", "c through a"},
			{"/root/b/c", "y", "c through b"},
		},
		update:  testPath{"/root/b/c", "new y", "new label"},
		wantErr: false,
	}, {
		name: "multiple roots",
		paths: []testPath{
			{"/root1/a", "x", "end of a from root1"},
			{"/root2/a", "y", "end of a from root2"},
		},
		update:  testPath{"/root1/a", "new x", "new label"},
		wantErr: false,
	}, {
		name: "key doesn't exist - leaf node",
		paths: []testPath{
			{"/root/a/b", "x", "label for x"},
			{"/root/a/c", "y", "label for y"},
		},
		update:  testPath{"/root/a/d", "fail", "fail"},
		wantErr: true,
	}, {
		name: "key doesn't exist - middle node",
		paths: []testPath{
			{"/root/a/b", "x", "label for x"},
			{"/root/a/c", "y", "label for y"},
		},
		update:  testPath{"/root/x/b", "fail", "fail"},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tree := New()
			for _, p := range tt.paths {
				_ = tree.Add(p.path, p.value, p.label)
			}

			// Act
			if err := tree.Update(tt.update.path, tt.update.value, tt.update.label); (err != nil) != tt.wantErr {
				t.Errorf("Update() error: %v; wantErr: %v", err, tt.wantErr)
			}

			// Assert
			got, _ := tree.Get(tt.update.path)
			if diff := cmp.Diff(tt.update.value, got); diff != "" && !tt.wantErr {
				t.Errorf("diff (-want / +got):\n%v", diff)
			}
		})
	}
}

func TestTree_Flatten(t *testing.T) {
	type want struct {
		nodes []*Node
		po    PrintOption
	}
	tests := []struct {
		name  string
		paths []testPath
		want  []want
	}{{
		name:  "empty",
		paths: []testPath{},
		want: []want{{
			nodes: []*Node{},
		}},
	}, {
		name: "one root, two gens",
		paths: []testPath{
			{"/root/a/b", "y", ""},
			{"/root/a/c", "x", ""},
			{"/root/a", "z", ""},
			{"/root/d", "w", ""},
		},
		want: []want{{
			po: Key,
			nodes: []*Node{
				{Key: "root"},
				{Key: "root/a", Value: "z"},
				{Key: "root/a/b", Value: "y"},
				{Key: "root/a/c", Value: "x"},
				{Key: "root/d", Value: "w"},
			},
		}, {
			po: Value,
			nodes: []*Node{
				{Key: "root"},
				{Key: "root/d", Value: "w"},
				{Key: "root/a", Value: "z"},
				{Key: "root/a/c", Value: "x"},
				{Key: "root/a/b", Value: "y"},
			},
		}},
	}, {
		name: "multi roots",
		paths: []testPath{
			{"/root1", "z", ""},
			{"/root2", "y", ""},
			{"/root3", "x", ""},
		},
		want: []want{{
			po: Key,
			nodes: []*Node{
				{Key: "root1", Value: "z"},
				{Key: "root2", Value: "y"},
				{Key: "root3", Value: "x"},
			},
		}, {
			po: Value,
			nodes: []*Node{
				{Key: "root3", Value: "x"},
				{Key: "root2", Value: "y"},
				{Key: "root1", Value: "z"},
			},
		}},
	}, {
		name: "multi roots, two gens",
		paths: []testPath{
			{"/root1", "k1", ""},
			{"/root1/a/x", "n1", ""},
			{"/root1/a/y", "m1", ""},
			{"/root1/a", "t1", ""},
			{"/root1/b", "s1", ""},
			{"/root2/a/x", "n2", ""},
			{"/root2/a/y", "m2", ""},
			{"/root2/a", "t2", ""},
			{"/root2/b", "s2", ""},
			{"/root2", "j2", ""},
		},
		want: []want{{
			po: Key,
			nodes: []*Node{
				{Key: "root1", Value: "k1"},
				{Key: "root1/a", Value: "t1"},
				{Key: "root1/a/x", Value: "n1"},
				{Key: "root1/a/y", Value: "m1"},
				{Key: "root1/b", Value: "s1"},
				{Key: "root2", Value: "j2"},
				{Key: "root2/a", Value: "t2"},
				{Key: "root2/a/x", Value: "n2"},
				{Key: "root2/a/y", Value: "m2"},
				{Key: "root2/b", Value: "s2"},
			},
		}, {
			po: Value,
			nodes: []*Node{
				{Key: "root2", Value: "j2"},
				{Key: "root2/b", Value: "s2"},
				{Key: "root2/a", Value: "t2"},
				{Key: "root2/a/y", Value: "m2"},
				{Key: "root2/a/x", Value: "n2"},
				{Key: "root1", Value: "k1"},
				{Key: "root1/b", Value: "s1"},
				{Key: "root1/a", Value: "t1"},
				{Key: "root1/a/y", Value: "m1"},
				{Key: "root1/a/x", Value: "n1"},
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tw := range tt.want {
				t.Run(fmt.Sprintf("%s [%s]", tt.name, tw.po), func(t *testing.T) {
					tree := New()
					for _, p := range tt.paths {
						_ = tree.Add(p.path, p.value, p.label)
					}
					tree.PrintOption = tw.po
					nodes := tree.Flatten()

					if diff := cmp.Diff(tw.nodes, nodes, cmp.AllowUnexported(Node{})); diff != "" {
						t.Errorf("diff (-want / +got):\n%v", diff)
					}
				})
			}
		})
	}
}
