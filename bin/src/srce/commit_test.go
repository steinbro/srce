package srce

import (
	"bufio"
	"path/filepath"
	"strings"
	"testing"
)

func (r Repo) checkTree(t *testing.T, sha1 string) {
	obj, err := r.Fetch(sha1)
	if err != nil {
		t.Error(err)
	}

	scanner := bufio.NewScanner(&obj.contents)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		t.Logf("looking for %s %s", parts[0], parts[1])
		if parts[0] == "tree" {
			r.checkTree(t, parts[1])
		} else {
			if _, err := r.Fetch(parts[1]); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestCommit(t *testing.T) {
	repo := setUp(t)

	testFile := filepath.Join(repo.Dir, "HEAD")
	if err := repo.Add(testFile); err != nil {
		t.Fatal(err)
	}

	if err := repo.Commit("test commit"); err != nil {
		t.Fatal(err)
	}

	if entries, _ := repo.getIndex().read(); len(entries) > 0 {
		t.Error("index not empty after commit")
	}

	refHash, err := repo.Resolve("master")
	if err != nil {
		t.Error(err)
	}
	commitObj, err := repo.Fetch(refHash)
	if err != nil {
		t.Errorf("master (%s) not in repo", refHash)
	}

	commitData := commitObj.contents.String()
	treeHash := string(commitData[5:45])
	repo.checkTree(t, treeHash)

	tearDown(t)
}

func TestCommitEmpty(t *testing.T) {
	repo := setUp(t)
	if err := repo.Commit("test commit"); err == nil {
		t.Error("empty commit succeeded")
	}
	tearDown(t)
}
