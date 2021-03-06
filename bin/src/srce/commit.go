package srce

import (
	"fmt"
)

func (r Repo) Commit(message string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	index := r.getIndex()
	entries, err := index.read()
	if err != nil {
		return fmt.Errorf("nothing to commit")
	}

	root := makeTree()
	items := 0
	for e := range entries {
		items++
		root.add(e.path, e.sha1)
	}
	if items == 0 {
		return fmt.Errorf("nothing to commit")
	}
	treeObj := r.storeTree(root)

	// Use current HEAD as parent commit
	parentHash, err := r.Resolve("HEAD")
	if err != nil {
		// no valid HEAD (is this the first commit?)
		parentHash = initialCommitHash
	}

	commitObj, err := commitObject(treeObj, parentHash, message)
	if err != nil {
		return err
	}
	if err := r.Store(commitObj); err != nil {
		return err
	}

	// Update .srce/refs/heads/master to point to commit object
	branch, err := r.GetSymbolicRef("HEAD")
	if err != nil {
		return err
	}
	// create refs/heads/<branch> if it doesn't yet exist
	if _, err := r.expandRef(branch); err != nil {
		if err := r.createRef(branch); err != nil {
			return err
		}
	}
	if err := r.UpdateRef(branch, string(commitObj.sha1)); err != nil {
		return err
	}

	// update reflog (for both HEAD and branch)
	refMessage := fmt.Sprintf("commit: %s", message)
	c, _ := parseCommit(commitObj.contents)

	headRefLog := r.getRefLog("HEAD")
	headRefLog.add(parentHash, commitObj.sha1, c.author, refMessage)

	oldBranchHash, err := r.Resolve(branch)
	if err != nil {
		oldBranchHash = initialCommitHash
	}
	branchRefLog := r.getRefLog(branch)
	branchRefLog.add(oldBranchHash, commitObj.sha1, c.author, refMessage)

	// Reset .srce/index
	if err := index.clear(); err != nil {
		return err
	}

	return nil
}
