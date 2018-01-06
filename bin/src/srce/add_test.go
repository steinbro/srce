package srce

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setUp(t *testing.T) {
	// Remove any stale test directories
	if err := os.RemoveAll(testFolder); err != nil {
		t.Fatal(err)
	}

	if err := Init(testFolder); err != nil {
		t.Fatal(err)
	}
}

func tearDown(t *testing.T) {
	// Remove temporary test directory
	if err := os.RemoveAll(testFolder); err != nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	setUp(t)
	// Check no errors are raised
	testFile := filepath.Join(testFolder, "HEAD")
	if err := Add(testFolder, testFile); err != nil {
		t.Fatal(err)
	}

	// Index should be created, with one line
	indexPath := filepath.Join(testFolder, "index")
	indexFIle, err := ioutil.ReadFile(indexPath)
	if err != nil {
		t.Fatal("index file not readable after Add")
	}
	indexLine := strings.Split(string(indexFIle), " ")
	hash := indexLine[0]

	// New non-empty blob should exist
	blobPath := filepath.Join(testFolder, "objects", hash[:2], hash[2:])
	stat, err := os.Stat(blobPath)
	if err != nil {
		t.Fatalf("Blob file %s unreadable", blobPath)
	}
	if stat.Size() == 0 {
		t.Fatalf("Blob file %s empty", blobPath)
	}
	tearDown(t)
}

func TestAddNonexistent(t *testing.T) {
	setUp(t)
	if err := Add(testFolder, "nonexistent"); err == nil {
		t.Fatal("Add nonexistent file succeeded")
	} else {
		t.Log(err)
	}
	tearDown(t)
}

func TestAddOutside(t *testing.T) {
	if err := Add(testFolder, "nonexistent"); err == nil {
		t.Fatal("Add to non-project succeeded")
	} else {
		t.Log(err)
	}
}

func TestAddUnwritable(t *testing.T) {
	setUp(t)
	// Make project folder read-only
	if err := os.Chmod(testFolder, 0500); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(testFolder, "HEAD")
	if err := Add(testFolder, testFile); err == nil {
		t.Fatal("Add to non-writable project succeeded")
	} else {
		t.Log(err)
	}

	// Restore writability
	if err := os.Chmod(testFolder, 0700); err != nil {
		t.Fatal(err)
	}
	tearDown(t)
}
