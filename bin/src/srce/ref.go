package srce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (r Repo) internalPath(parts ...string) string {
	things := append([]string{r.Dir}, parts...)
	return filepath.Join(things...)
}

func (r Repo) UpdateRef(ref string, hash Hash) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	// write hash to e.g. .srce/refs/heads/master
	refFile, err := os.OpenFile(
		r.internalPath(ref), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := refFile.WriteString(fmt.Sprintf("%s\n", hash)); err != nil {
		return err
	}
	refFile.Close()
	return nil
}

func (r Repo) GetSymbolicRef(name string) (string, error) {
	if !r.IsInitialized() {
		return "", fmt.Errorf("not a srce project")
	}

	// get value of "ref: " from .srce/HEAD
	data, err := ioutil.ReadFile(r.internalPath(name))
	if err != nil {
		return "", err
	}
	pattern := regexp.MustCompile("ref: (.+)\n")
	ref := pattern.FindSubmatch(data)
	if len(ref) != 2 {
		return "", fmt.Errorf("malformed ref: %q", data)
	}
	return string(ref[1]), nil
}

func (r Repo) SetSymbolicRef(name, ref string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	return ioutil.WriteFile(
		r.internalPath(name),
		[]byte(fmt.Sprintf("ref: %s\n", ref)), 0644)
}

func (r Repo) Resolve(name string) (Hash, error) {
	if !r.IsInitialized() {
		return Hash(""), fmt.Errorf("not a srce project")
	}

	// is it a branch, or a special name referring to a branch?
	possibleRef := r.internalPath("refs", "heads", name)
	// prevent relative paths, e.g. "../../HEAD"
	if !strings.HasSuffix(possibleRef, name) {
		return Hash(""), fmt.Errorf("bad ref: %s", name)
	}
	if ref, err := r.GetSymbolicRef(name); err == nil {
		// already includes the refs/heads part
		possibleRef = r.internalPath(ref)
	}

	// is it a branch?
	if refValue, err := ioutil.ReadFile(possibleRef); err == nil {
		// case "master"
		return ValidateHash(strings.TrimSpace(string(refValue)))
	} else if refValue, err := ioutil.ReadFile(r.internalPath(name)); err == nil {
		// case "refs/heads/master"
		return ValidateHash(strings.TrimSpace(string(refValue)))
	}

	// finally, is it an object hash, or unambiguous prefix?
	if hash, err := ValidateHash(name); err == nil {
		return r.ExpandPartialHash(hash)
	} else {
		return Hash(""), err
	}
}
