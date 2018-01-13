package srce

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func (r Repo) getObject(sha1 string) (io.Reader, error) {
	objData, _ := os.Open(filepath.Join(r.Dir, "objects", sha1[:2], sha1[2:]))
	return zlib.NewReader(objData)
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

	refHash, err := repo.Resolve("master")
	if err != nil {
		t.Error(err)
	}
	commitFile, err := repo.getObject(refHash)
	if err != nil {
		t.Errorf("master (%s) not in repo", refHash)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(commitFile)
	commitData := buf.String()
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
