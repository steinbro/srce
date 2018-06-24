package srce

import "testing"

func TestResolve(t *testing.T) {
	repo := setUp(t)
	repo.Add(repo.internalPath("HEAD"))
	repo.Commit("test commit")
	hash, err := repo.Resolve("HEAD")
	if err != nil {
		t.Error(err)
	}

	// good cases
	goodRefs := []string{
		"master",
		"refs/heads/master",
		string(hash),
		hash.abbreviated(),
	}
	for _, input := range goodRefs {
		if result, err := repo.Resolve(input); err != nil {
			t.Errorf("Unexpected error for ref %q: %q", input, err)
		} else if result != hash {
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
