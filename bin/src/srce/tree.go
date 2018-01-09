package srce

import (
	"fmt"
	"strings"
)

type Node struct {
	name     string
	sha1     string
	children map[string]*Node
}

func newNode(name string) *Node {
	return &Node{
		name: name, sha1: timestampedHash(name), children: make(map[string]*Node)}
}

func makeTree() *Node {
	return newNode("root")
}

func (n *Node) add(path, sha1 string) {
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

func (r Repo) storeTree(n *Node) Object {
	treeObj := Object{sha1: n.sha1}
	for _, c := range n.children {
		if len(c.children) > 0 {
			r.storeTree(c)
			treeObj.contents.Write([]byte(
				fmt.Sprintf("tree %s %s\n", c.sha1, c.name)))
		} else {
			treeObj.contents.Write([]byte(
				fmt.Sprintf("blob %s %s\n", c.sha1, c.name)))
		}
	}
	r.Store(treeObj)
	return treeObj
}
