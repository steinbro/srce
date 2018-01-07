package srce

import (
	"os"
	"path/filepath"
	"testing"
)

func setUp(t *testing.T) Repo {
	repo := Repo{Dir: testFolder}

	// Remove any stale test directories
	if err := os.RemoveAll(repo.Dir); err != nil {
		t.Fatal(err)
	}
	if err := repo.Init(); err != nil {
		t.Fatal(err)
	}
	return repo
}

func tearDown(t *testing.T) {
	// Remove temporary test directory
	if err := os.RemoveAll(testFolder); err != nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	repo := setUp(t)
	// Check no errors are raised
	testFile := filepath.Join(repo.Dir, "HEAD")
	if err := repo.Add(testFile); err != nil {
		t.Fatal(err)
	}

	// Index should be created, with one entry
	entries, err := repo.getIndex().read()
	if err != nil {
		t.Fatal("index file not readable after Add")
	}
	hash := (<-entries).sha1

	// New non-empty blob should exist
	blobPath := filepath.Join(repo.Dir, "objects", hash[:2], hash[2:])
	stat, err := os.Stat(blobPath)
	if err != nil {
		t.Errorf("Blob file %s unreadable", blobPath)
	} else if stat.Size() == 0 {
		t.Errorf("Blob file %s empty", blobPath)
	}
	tearDown(t)
}

func TestAddNonexistent(t *testing.T) {
	repo := setUp(t)
	if err := repo.Add("nonexistent"); err == nil {
		t.Error("Add nonexistent file succeeded")
	} else {
		t.Log(err)
	}
	tearDown(t)
}

func TestAddOutside(t *testing.T) {
	repo := Repo{Dir: testFolder}
	if err := repo.Add("nonexistent"); err == nil {
		t.Fatal("Add to non-project succeeded")
	} else {
		t.Log(err)
	}
}

func TestAddUnwritable(t *testing.T) {
	repo := setUp(t)
	// Make project folder read-only
	if err := os.Chmod(repo.Dir, 0500); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(repo.Dir, "HEAD")
	if err := repo.Add(testFile); err == nil {
		t.Fatal("Add to non-writable project succeeded")
	} else {
		t.Log(err)
	}

	// Restore writability
	if err := os.Chmod(repo.Dir, 0700); err != nil {
		t.Fatal(err)
	}
	tearDown(t)
}
