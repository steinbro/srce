package srce

import (
	"os"
	"testing"
)

func TestResolve(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

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
}

func TestUpdateRefBad(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// create initial HEAD and master ref
	repo.commitSomething(t)

	badRefs := []struct{ ref, hash string }{
		{"wont/expand", "notahash"},
		{"refs/heads/master", "notahash"},
	}

	for _, r := range badRefs {
		if err := repo.UpdateRef(r.ref, r.hash); err == nil {
			t.Errorf("UpdateRef(%q, %q) succeeded (expecting error)", r.ref, r.hash)
		}
	}
}

func TestUpdateRefUnwritable(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// create initial HEAD and master ref
	repo.commitSomething(t)
	hash, _ := repo.getLastCommit(t)

	// Make master ref read-only
	if err := os.Chmod(repo.internalPath("refs", "heads", "master"), 0500); err != nil {
		t.Fatal(err)
	}
	// Restore writability when finished
	defer os.Chmod(repo.internalPath(), 0700)

	if err := repo.UpdateRef("refs/heads/master", string(hash)); err == nil {
		t.Error("UpdateRef on unwritable repo succeeded ")
	}
}
