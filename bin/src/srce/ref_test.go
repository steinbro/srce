package srce

import "testing"

func TestResolve(t *testing.T) {
	repo := setUp(t)
	repo.Add(repo.internalPath("HEAD"))
	repo.Commit("test commit")
	hash, _ := repo.Resolve("HEAD")

	// good cases
	for _, input := range []string{"master", string(hash), string(hash[:4])} {
		if result, _ := repo.Resolve(input); result != hash {
			t.Errorf("Resolve(%q) = %q (expecting %q)", input, result, hash)
		}
	}

	// bad cases
	badRefs := []string{
		"undefined",
		string(hash[:3]),
		"index",
		"../../HEAD",
		"../../../../../../../../../../etc/passwd",
	}
	for _, input := range badRefs {
		if result, err := repo.Resolve(input); err == nil {
			t.Errorf("Resolve(%q) = %q (expecting error)", input, result)
		}
	}

	tearDown(t)
}
