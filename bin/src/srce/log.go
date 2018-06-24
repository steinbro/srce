package srce

import (
	"fmt"
)

func (r Repo) Log() error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	hash, err := r.Resolve("HEAD")
	if err != nil {
		return err
	}

	for hash != "" {
		o, err := r.Fetch(hash)
		if err != nil {
			return err
		}

		c, err := parseCommit(o.contents)
		if err != nil {
			return err
		}

		fmt.Printf("commit %s\n", hash)
		fmt.Printf("Author: %s\n", c.author.user)
		fmt.Printf(
			"Date:   %s\n", c.author.timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Printf("\n\t%s\n\n", c.message)

		hash = c.parent
	}

	return nil
}
