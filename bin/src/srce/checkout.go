package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (r Repo) changeHead(newName string) error {
	oldName, err := r.GetSymbolicRef("HEAD")
	if err != nil {
		return err
	}
	// if current HEAD has commits, update reflog
	if oldHash, err := r.Resolve(oldName); err == nil {
		// hash sometimes doesn't change, e.g. when switching to new branch
		newHash, _ := r.Resolve(newName)

		// update reflog
		refMessage := fmt.Sprintf(
			"checkout: moving from %s to %s", filepath.Base(oldName), newName)
		headRefLog := r.getRefLog("HEAD")
		if err := headRefLog.add(oldHash, newHash, currentUserAndTime(), refMessage); err != nil {
			return err
		}
	}

	// point HEAD at new branch
	return r.SetSymbolicRef("HEAD", newName)
}

func (r Repo) CreateBranch(branchName string) error {
	if _, err := r.expandRef(branchName); err == nil {
		return fmt.Errorf("%s: already exists", branchName)
	}

	// refs/head/<branchName> will actually be created upon commit
	return r.changeHead(filepath.Join("refs", "heads", branchName))
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

	return r.changeHead(refname)
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
