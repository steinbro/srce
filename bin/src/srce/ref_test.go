package srce

import (
	"path/filepath"
	"testing"
)

func TestResolve(t *testing.T) {
	repo := setUp(t)
	repo.Add(filepath.Join(repo.Dir, "HEAD"))
	repo.Commit("test commit")
	hash, _ := repo.Resolve("HEAD")

	// good cases
	for _, input := range []string{"master", hash, hash[:4]} {
		if result, _ := repo.Resolve(input); result != hash {
			t.Errorf("Resolve(%s) = %s (expecting %s)", input, result, hash)
		}
	}

	// bad cases
	for _, input := range []string{"undefined", hash[:3]} {
		if result, err := repo.Resolve(input); err == nil {
			t.Errorf("Resolve(%s) = %s (expecting error)", input, result, hash)
		}
	}

	tearDown(t)
}
