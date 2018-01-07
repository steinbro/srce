package srce

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func (r Repo) getObject(sha1 string) (io.Reader, error) {
	return os.Open(filepath.Join(r.Dir, "objects", sha1[:2], sha1[2:]))
}

func (r Repo) checkTree(t *testing.T, sha1 string) {
	objData, err := r.getObject(sha1)
	if err != nil {
		t.Error(err)
	}

	scanner := bufio.NewScanner(objData)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		t.Logf("looking for %s %s", parts[0], parts[1])
		if parts[0] == "tree" {
			r.checkTree(t, parts[1])
		} else {
			if _, err := r.getObject(parts[1]); err != nil {
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

	ref := filepath.Join(repo.Dir, "refs", "heads", "master")
	refValue, err := ioutil.ReadFile(ref)
	if err != nil {
		t.Error(err)
	} else if len(refValue) == 0 {
		t.Error("refs/heads/master not updated")
	}

	refHash := strings.TrimSpace(string(refValue))
	commitFile, err := repo.getObject(refHash)
	if err != nil {
		t.Errorf("master (%s) not in repo", refValue)
	}

	commitData, _ := ioutil.ReadAll(commitFile)
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
