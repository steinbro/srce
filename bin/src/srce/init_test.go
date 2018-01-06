package srce

import (
	"os"
	"path/filepath"
	"testing"
)

const testFolder = ".srce-test"

func TestInit(t *testing.T) {
	// Remove any stale test directories
	os.RemoveAll(testFolder)

	// Check no errors are raised
	if err := Init(testFolder); err != nil {
		t.Fatal(err)
	}

	// Check that HEAD file was created
	headFile := filepath.Join(testFolder, "HEAD")
	if _, err := os.Stat(headFile); os.IsNotExist(err) {
		t.Fatalf("%s doesn't exist after Init", headFile)
	}

	// Remove temporary test directory
	os.RemoveAll(testFolder)
}

func TestInitBad(t *testing.T) {
	// Create temporary test "directory"
	os.OpenFile(testFolder, os.O_RDONLY|os.O_CREATE, 0666)

	// Check that an error is raised
	if err := Init(testFolder); err == nil {
		t.Fatalf("Init succeeded when %s already exists", testFolder)
	}

	// Remove temporary test "directory"
	os.Remove(testFolder)
}
