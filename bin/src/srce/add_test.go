package srce

import (
	"os"
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
	defer tearDown(t)

	// Check no errors are raised
	testFile := repo.internalPath("HEAD")
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
	} else if blob.Size() == 0 {
		t.Errorf("Blob %s empty", hash)
	} else if blob.Type() != BlobObject {
		t.Errorf("Blob of wrong type (%s)", blob.Type())
	} else if blob.Size() != len(blob.Contents()) {
		t.Errorf(
			"Blob header/size mismatch (%d != %d)", blob.Size(), len(blob.Contents()))
	}
}

func TestAddNonexistent(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	if err := repo.Add("nonexistent"); err == nil {
		t.Error("Add nonexistent file succeeded")
	}
}

func TestAddUnwritable(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// Test both failure to write index, and write object
	for _, unwritableDir := range []string{repo.Dir, repo.internalPath("objects")} {
		// Make folder read-only
		if err := os.Chmod(unwritableDir, 0500); err != nil {
			t.Fatal(err)
		}

		testFile := repo.internalPath("HEAD")
		if err := repo.Add(testFile); err == nil {
			t.Errorf("Add to non-writable %q succeeded", unwritableDir)
		}

		// Restore writability, and recreate repo
		os.Chmod(unwritableDir, 0700)
		tearDown(t)
		repo = setUp(t)
	}
}
