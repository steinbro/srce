package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func (r Repo) UpdateRef(ref, hash string) error {
	// write hash to e.g. .srce/refs/heads/master
	refFile, err := os.OpenFile(
		filepath.Join(r.Dir, ref), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := refFile.Write([]byte(fmt.Sprintf("%s\n", hash))); err != nil {
		return err
	}
	refFile.Close()
	return nil
}

func (r Repo) GetSymbolicRef(name string) (string, error) {
	// get value of "ref: " from .srce/HEAD
	data, err := ioutil.ReadFile(filepath.Join(r.Dir, name))
	if err != nil {
		return "", err
	}
	pattern, _ := regexp.Compile("ref: (.+)\n")
	ref := pattern.FindSubmatch(data)
	return string(ref[1]), nil
}

func (r Repo) SetSymbolicRef(name, ref string) error {
	return ioutil.WriteFile(
		filepath.Join(r.Dir, name),
		[]byte(fmt.Sprintf("ref: %s\n", ref)), 0644)
}
