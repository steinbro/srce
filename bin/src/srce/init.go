package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DotDir = ".srce"

func Init(dotDir string) error {
	if _, err := os.Stat(dotDir); os.IsNotExist(err) {
		os.Mkdir(dotDir, 0700)
	} else {
		return fmt.Errorf("%s already exists", dotDir)
	}

	subdirs := []string{"objects/info", "objects/pack", "refs/heads", "refs/tags"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(dotDir, subdir), 0700); err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(
		filepath.Join(dotDir, "HEAD"),
		[]byte("ref: refs/heads/master\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}
