package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (r Repo) CreateBranch(branchName string) error {
	if _, err := r.expandRef(branchName); err != nil {
		// point HEAD at new branch
		// refs/head/<branchName> will be created upon commit
		if err := r.SetSymbolicRef(
			"HEAD", filepath.Join("refs", "heads", branchName)); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("%s: already exists", branchName)
	}
}

func (r Repo) loadCommitTree(refname string) (*Node, error) {
	refHash, err := r.Resolve(refname)
	if err != nil {
		return nil, err
	}
	refObj, err := r.Fetch(refHash)
	if err != nil {
		return nil, err
	}
	refCommit, err := parseCommit(refObj.contents)
	if err != nil {
		return nil, err
	}
	tree := makeTree()
	if err := r.loadTree(tree, refCommit.tree, ""); err != nil {
		return nil, err
	}
	return tree, nil
}

func (r Repo) CheckoutTree(refname string) error {
	tree, err := r.loadCommitTree(refname)
	if err != nil {
		return err
	}

	for path := range tree.walk() {
		if err := r.restoreFile(tree, path); err != nil {
			return err
		}
	}

	return nil
}

func (r Repo) CheckoutFile(refname, filepath string) error {
	tree, err := r.loadCommitTree(refname)
	if err != nil {
		return err
	}

	return r.restoreFile(tree, filepath)
}

func (r Repo) restoreFile(tree *Node, filepath string) error {
	fileHash, err := tree.get(filepath)
	if err != nil {
		return err
	}
	fileObj, err := r.Fetch(fileHash)
	if err != nil {
		return err
	}
	if fileObj.otype == TreeObject {
		return os.MkdirAll(filepath, 0644)
	} else if fileObj.otype == BlobObject {
		return ioutil.WriteFile(filepath, []byte(fileObj.Contents()), 0644)
	} else {
		return fmt.Errorf("cannot restore object of type %q", fileObj.otype)
	}
}
