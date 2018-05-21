package srce

import (
	"fmt"
)

func (r Repo) Log() error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	hash, _ := r.Resolve("HEAD")
	for hash != "" {
		o, _ := r.Fetch(hash)
		c, _ := parseCommit(o.contents)
		fmt.Printf("commit %s\n", hash)
		fmt.Printf("author %s\n", c.author)
		fmt.Printf("\n\t%s\n\n", c.message)
		hash = c.parent
	}

	return nil
}
