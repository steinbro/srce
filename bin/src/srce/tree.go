package srce

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

type Node struct {
	name     string
	sha1     Hash
	children map[string]*Node
}

func newNode(name string) *Node {
	return &Node{
		name: name, sha1: timestampedHash(name), children: make(map[string]*Node)}
}

func makeTree() *Node {
	return newNode("root")
}

func (n *Node) add(path string, sha1 Hash) {
	pathComponents := strings.Split(path, "/")
	this := n
	for _, thing := range pathComponents {
		if _, found := this.children[thing]; !found {
			this.children[thing] = newNode(thing)
		}
		this = this.children[thing]
	}
	// Provided hash applies to last path component (the file)
	this.sha1 = sha1
}

// Find hash for path by walking tree.
func (n *Node) get(path string) (Hash, error) {
	this := n
	var ok bool
	for _, pathComponent := range strings.Split(path, "/") {
		if this, ok = this.children[pathComponent]; !ok {
			return Hash(""), fmt.Errorf("path not found in tree: %s", path)
		}
	}
	return this.sha1, nil
}

func (r Repo) storeTree(n *Node) Object {
	treeObj := Object{otype: TreeObject, sha1: n.sha1}
	for _, c := range n.children {
		if len(c.children) > 0 {
			r.storeTree(c)
			treeObj.contents.WriteString(fmt.Sprintf("tree %s %s\n", c.sha1, c.name))
		} else {
			treeObj.contents.WriteString(fmt.Sprintf("blob %s %s\n", c.sha1, c.name))
		}
	}
	r.Store(treeObj)
	return treeObj
}

func parseTreeEntry(line string) (ObjectType, string, Hash, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return TreeObject, "", Hash(""), fmt.Errorf("malformed tree entry: %q", line)
	}
	return ObjectType(parts[0]), parts[2], Hash(parts[1]), nil
}

// Reconstruct a tree recursively from stored objects.
func (r Repo) loadTree(root *Node, sha1 Hash, path string) error {
	// read object for specified hash; check that it's a tree
	obj, err := r.Fetch(sha1)
	if err != nil {
		return err
	}
	if obj.otype != TreeObject {
		return fmt.Errorf("not a tree: %s", sha1)
	}

	scanner := bufio.NewScanner(&obj.contents)
	for scanner.Scan() {
		this_type, this_name, this_sha1, err := parseTreeEntry(scanner.Text())
		if err != nil {
			return err
		}

		this_path := filepath.Join(path, this_name)
		root.add(this_path, this_sha1)

		// recurse if necessary
		if this_type == TreeObject {
			if err := r.loadTree(root, this_sha1, this_path); err != nil {
				return err
			}
		}
	}
	return nil
}
