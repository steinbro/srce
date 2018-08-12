package srce

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func (r Repo) checkTree(t *testing.T, sha1 Hash) {
	obj, err := r.Fetch(sha1)
	if err != nil {
		t.Error(err)
	}

	scanner := bufio.NewScanner(&obj.contents)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		otype := ObjectType(parts[0])
		ohash := Hash(parts[1])
		//t.Logf("looking for %s %s", otype, ohash)
		if otype == TreeObject {
			r.checkTree(t, ohash)
		} else {
			if o, err := r.Fetch(ohash); err != nil {
				t.Error(err)
			} else if otype != o.Type() {
				t.Errorf("object type mismatch (expected %s, got %s)", parts[0], o.Type())
			}
		}
	}
}

func (r Repo) commitSomething(t *testing.T) {
	testFile := r.internalPath("HEAD")
	if err := r.Add(testFile); err != nil {
		t.Fatal(err)
	}

	if err := r.Commit("test commit"); err != nil {
		t.Fatal(err)
	}
}

func (r Repo) getLastCommit(t *testing.T) (Hash, Commit) {
	refHash, err := r.Resolve("master")
	if err != nil {
		t.Error(err)
	}
	commitObj, err := r.Fetch(refHash)
	if err != nil {
		t.Errorf("master (%s) not in repo", refHash)
	}

	commit, err := parseCommit(commitObj.contents)
	if err != nil {
		t.Error(err)
	}
	return refHash, commit
}

func TestCommit(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)

	if entries, _ := repo.getIndex().read(); len(entries) > 0 {
		t.Error("index not empty after commit")
	}

	_, commit := repo.getLastCommit(t)
	repo.checkTree(t, commit.tree)
}

func TestCommitEmpty(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	if err := repo.Commit("test commit"); err == nil {
		t.Error("empty commit succeeded")
	}
}

func TestCommitParent(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)
	hash1, commit1 := repo.getLastCommit(t)
	if commit1.parent != "" {
		t.Error("parent of first commit not empty")
	}

	repo.commitSomething(t)
	_, commit2 := repo.getLastCommit(t)
	if commit2.parent != hash1 {
		t.Error("first commit is not parent of second commit")
	}
	if commit2.message != "test commit" {
		t.Errorf("unexpected commit message (%q)", commit2.message)
	}
}

func TestCommitMalformed(t *testing.T) {
	testCases := []string{
		"tree notahash\n",
		"author notanauthorstamp\n",
		"parent notahash\n",
		"non sense\n",
	}

	for _, b := range testCases {
		var buf bytes.Buffer
		buf.WriteString(b)
		if _, err := parseCommit(buf); err == nil {
			t.Errorf("parseCommit(%q) succeeded (expected error)", b)
		}
	}
}

func TestCommitUnwritableIndex(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// Add a file to the index
	testFile := repo.internalPath("HEAD")
	if err := repo.Add(testFile); err != nil {
		t.Error(err)
	}

	// Make index read-only
	if err := os.Chmod(repo.getIndex().path, 0500); err != nil {
		t.Fatal(err)
	}
	// Restore writability when finished
	defer os.Chmod(repo.getIndex().path, 0700)

	// Should fail when attempting to clear the index
	if err := repo.Commit("test"); err == nil {
		t.Errorf("Commit with non-writable index succeeded")
	}
}
