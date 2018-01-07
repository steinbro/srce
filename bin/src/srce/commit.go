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

	commitObj, err := commitObject(treeObj, message)
	if err != nil {
		return err
	}
	r.Store(commitObj)

	// Update .srce/refs/heads/master to point to commit object
	if err := r.updateHead(commitObj.sha1); err != nil {
		return err
	}

	// Reset .srce/index
	if err := index.clear(); err != nil {
		return err
	}

	return nil
}