package srce

import (
	"fmt"
	"path/filepath"
	"regexp"
)

type Hash string

func ValidateHash(input string) (Hash, error) {
	f := regexp.MustCompile("^[0-9a-f]{4,40}$")
	if !f.MatchString(input) {
		return Hash(""), fmt.Errorf("invalid hash: %q", input)
	}
	return Hash(input), nil
}

func (r Repo) ExpandPartialHash(h Hash) (Hash, error) {
	// check if specified prefix is unambiguous
	pattern := filepath.Join(r.Dir, "objects", h.prefix(), h.remainder()+"*")
	if matches, _ := filepath.Glob(pattern); len(matches) == 1 {
		return Hash(h.prefix() + filepath.Base(matches[0])), nil
	} else if len(matches) > 1 {
		return Hash(""), fmt.Errorf("ambiguous name: %s", h)
	} else {
		return Hash(""), fmt.Errorf("no matching objects in repo: %s", h)
	}
}

func (h Hash) abbreviated() string {
	return string(h[:8])
}

func (h Hash) prefix() string {
	return string(h[:2])
}

func (h Hash) remainder() string {
	return string(h[2:])
}
