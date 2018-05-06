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

func (r Repo) UpdateRef(ref, hash string) error {
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
	pattern, _ := regexp.Compile("ref: (.+)\n")
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

func (r Repo) Resolve(name string) (string, error) {
	if !r.IsInitialized() {
		return "", fmt.Errorf("not a srce project")
	}

	// is it a branch, or a special name referring to a branch?
	possibleRef := r.internalPath("refs", "heads", name)
	// prevent relative paths, e.g. "../../HEAD"
	if !strings.HasSuffix(possibleRef, name) {
		return "", fmt.Errorf("bad ref: %s", name)
	}
	if ref, err := r.GetSymbolicRef(name); err == nil {
		// already includes the refs/heads part
		possibleRef = r.internalPath(ref)
	}

	// is it a branch?
	if refValue, err := ioutil.ReadFile(possibleRef); err == nil {
		return strings.TrimSpace(string(refValue)), nil
	}

	// is it an object hash, or unambiguous prefix?
	if len(name) > 3 {
		pattern := filepath.Join(r.Dir, "objects", name[:2], name[2:]+"*")
		if matches, _ := filepath.Glob(pattern); len(matches) == 1 {
			return name[:2] + filepath.Base(matches[0]), nil
		} else if len(matches) > 1 {
			return "", fmt.Errorf("ambiguous name: %s", name)
		}
	}

	// nothing matched
	return "", fmt.Errorf("cannot resolve %s", name)
}
