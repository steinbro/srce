package srce

import (
	"fmt"
	"os"
	"path/filepath"
)

const DotDir = ".srce"

func (r Repo) Init() error {
	if !r.IsInitialized() {
		os.Mkdir(r.Dir, 0700)
	} else {
		return fmt.Errorf("%s already exists", r.Dir)
	}

	subdirs := []string{"objects/info", "objects/pack", "refs/heads", "refs/tags"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(r.Dir, subdir), 0700); err != nil {
			return err
		}
	}

	if err := r.SetSymbolicRef("HEAD", "refs/heads/master"); err != nil {
		return err
	}

	return nil
}
