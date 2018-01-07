package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Repo struct {
	Dir string
}

func (r Repo) IsInitialized() bool {
	// Check that .srce directory exists
	_, err := os.Stat(r.Dir)
	return !os.IsNotExist(err)
}

func (r Repo) Store(o Object) error {
	// Create .srce/objects/00/ directory (where 00 is the first 2 bytes of hash)
	objFolder := filepath.Join(r.Dir, "objects", o.sha1[:2])
	if err := os.MkdirAll(objFolder, 0700); err != nil {
		return err
	}

	// Write file contents to .srce/objects/00/rest_of_hash
	objPath := filepath.Join(objFolder, o.sha1[2:])
	if err := ioutil.WriteFile(objPath, o.contents.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (r Repo) updateHead(sha1 string) error {
	// get value of "ref: " from .srce/HEAD
	head, err := ioutil.ReadFile(filepath.Join(r.Dir, "HEAD"))
	if err != nil {
		return err
	}
	currentBranch := strings.TrimSpace(string(head[5:]))

	// write hash to e.g. .srce/refs/heads/master
	refFile, err := os.OpenFile(
		filepath.Join(r.Dir, currentBranch), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := refFile.Write([]byte(fmt.Sprintf("%s\n", sha1))); err != nil {
		return err
	}
	refFile.Close()
	return nil
}
