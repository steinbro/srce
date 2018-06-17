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
		parentHash = INITIAL_COMMIT_HASH
	}

	commitObj, err := commitObject(treeObj, parentHash, message)
	if err != nil {
		return err
	}
	r.Store(commitObj)

	// Update .srce/refs/heads/master to point to commit object
	branch, err := r.GetSymbolicRef("HEAD")
	if err != nil {
		return err
	}
	r.UpdateRef(branch, commitObj.sha1)

	// update reflog (for both HEAD and branch)
	refMessage := fmt.Sprintf("commit: %s", message)
	c, _ := parseCommit(commitObj.contents)

	headRefLog := r.getRefLog("HEAD")
	headRefLog.add(parentHash, commitObj.sha1, c.author, refMessage)

	oldBranchHash, err := r.Resolve(branch)
	if err != nil {
		oldBranchHash = INITIAL_COMMIT_HASH
	}
	branchRefLog := r.getRefLog(branch)
	branchRefLog.add(oldBranchHash, commitObj.sha1, c.author, refMessage)

	// Reset .srce/index
	if err := index.clear(); err != nil {
		return err
	}

	return nil
}
