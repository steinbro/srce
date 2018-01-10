package srce

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// Helper methods for comparing trees
func (n Node) toString() string {
	var out []string
	n._toString(&out, 0)
	return strings.Join(out, "\n")
}

func (n Node) _toString(out *[]string, indent int) {
	*out = append(*out, fmt.Sprintf("%s%s", strings.Repeat("  ", indent), n.name))
	// Sort each level alphabetically (to make output deterministic)
	names := make([]string, 0)
	for k, _ := range n.children {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, c := range names {
		n.children[c]._toString(out, indent+1)
	}
}

var treeTests = []struct {
	description string
	in          []string
	out         string
}{
	{"one file",
		[]string{"foo"},
		"root\n  foo"},
	{"one file (repeated)",
		[]string{"foo", "foo"},
		"root\n  foo"},
	{"one file (two levels)",
		[]string{"foo/bar"},
		"root\n  foo\n    bar"},
	{"one file (three levels)",
		[]string{"foo/bar/baz"},
		"root\n  foo\n    bar\n      baz"},
	{"one file (separate parent)",
		[]string{"foo", "foo/bar"},
		"root\n  foo\n    bar"},
	{"three files (two folders)",
		[]string{"foo/bar", "foo", "qux/baz"},
		"root\n  foo\n    bar\n  qux\n    baz"},
}

func TestTree(t *testing.T) {
	for _, tt := range treeTests {
		tree := makeTree()
		for _, path := range tt.in {
			tree.add(path, "")
		}
		if result := tree.toString(); result != tt.out {
			t.Errorf("%s\n[expected]\n%s\n[got]\n%s",
				tt.description, tt.out, result)
		}
	}
}

func TestTreeHash(t *testing.T) {
	desired_hash := "deadbeef"
	tree := makeTree()
	tree.add("foo/bar", desired_hash)

	actual := tree.children["foo"].children["bar"].sha1
	if actual != desired_hash {
		t.Errorf("%s should be %s", actual, desired_hash)
	}

	otherHash := tree.children["foo"].sha1
	if otherHash == desired_hash {
		t.Errorf("%s should not be %s", actual, desired_hash)
	}
}
