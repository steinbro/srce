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
	if blob, err := repo.Fetch(hash); err != nil {
		t.Errorf("Blob %s unreadable", hash)
	} else if blob.size == 0 {
		t.Errorf("Blob %s empty", hash)
	} else if blob.otype != "blob" {
		t.Errorf("Blob of wrong type (%s)", blob.otype)
	} else if blob.contents.Len() != blob.size {
		t.Errorf(
			"Blob header/size mismatch (%d != %d)", blob.size, blob.contents.Len())
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
