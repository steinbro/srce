package srce

import (
	"fmt"
	"io/ioutil"
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
	//TODO
	return nil
}

func (r Repo) CheckoutFile(refname, filepath string) error {
	tree, err := r.loadCommitTree(refname)
	if err != nil {
		return err
	}

	fileHash, err := tree.get(filepath)
	if err != nil {
		return err
	}
	fileObj, err := r.Fetch(fileHash)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, []byte(fileObj.Contents()), 0644)
}
