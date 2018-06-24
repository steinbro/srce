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

func (r Repo) UpdateRef(inputRef, inputHash string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	// validate user-specified ref
	ref, err := r.expandRef(inputRef)
	if err != nil {
		return err
	}

	// validate user-specified hash
	hash, err := r.Resolve(inputHash)
	if err != nil {
		return err
	}

	// write hash to e.g. refs/heads/master
	refFile, err := os.OpenFile(
		r.internalPath(ref), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer refFile.Close()

	_, err = refFile.WriteString(fmt.Sprintf("%s\n", hash))
	return err
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

// expandRef resolvews a user-specified ref into a fully-qualified branch name,
// e.g. "master" becomes "refs/heads/master" (if it exists).
func (r Repo) expandRef(input string) (string, error) {
	expandedInput := filepath.Join("refs", "heads", input)
	// prevent relative paths, e.g. "../../HEAD"
	if !strings.HasSuffix(expandedInput, input) {
		return "", fmt.Errorf("bad ref: %s", input)
	}

	if _, err := os.Stat(r.internalPath(expandedInput)); err == nil {
		// input was e.g. "master"
		return expandedInput, nil
	} else if _, err := os.Stat(r.internalPath(input)); err == nil {
		// input was e.g. "refs/heads/master"
		return input, nil
	}

	return "", fmt.Errorf("unrecognized ref: %s", input)
}

func (r Repo) createRef(name string) (err error) {
	_, err = os.OpenFile(
		r.internalPath(name), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	return
}

func (r Repo) Resolve(name string) (Hash, error) {
	if !r.IsInitialized() {
		return Hash(""), fmt.Errorf("not a srce project")
	}

	transformedInput := name
	// is it a special name referring to a branch? (e.g. "HEAD")
	if ref, err := r.GetSymbolicRef(transformedInput); err == nil {
		// HEAD -> refs/heads/master
		transformedInput = ref
	}

	// is it a branch?
	// master -> refs/heads/master
	if ref, err := r.expandRef(transformedInput); err == nil {
		if hash, err := ioutil.ReadFile(r.internalPath(ref)); err == nil {
			// refs/heads/master -> d41d09fa...
			transformedInput = strings.TrimSpace(string(hash))
		} else {
			return Hash(""), err
		}
	}

	// finally, is it an object hash, or unambiguous prefix?
	if hash, err := ValidateHash(transformedInput); err == nil {
		// d41d -> d41d09fa...
		return r.ExpandPartialHash(hash)
	} else {
		return Hash(""), err
	}
}
