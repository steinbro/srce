package srce

import (
	"os"
	"testing"
)

const testFolder = ".srce-test"

func TestInit(t *testing.T) {
	// When finished, remove temporary test directory
	defer os.RemoveAll(testFolder)

	// Check no errors are raised
	repo := Repo{Dir: testFolder}
	if err := repo.Init(); err != nil {
		t.Fatal(err)
	}

	// Check that HEAD file was created
	if _, err := repo.GetSymbolicRef("HEAD"); os.IsNotExist(err) {
		t.Error("HEAD doesn't exist after Init")
	}
}

func TestInitBad(t *testing.T) {
	// Create temporary test "directory"
	os.OpenFile(testFolder, os.O_RDONLY|os.O_CREATE, 0666)
	// When finished, remove temporary test "directory"
	defer os.Remove(testFolder)

	// Check that an error is raised
	repo := Repo{Dir: testFolder}
	if err := repo.Init(); err == nil {
		t.Errorf("Init succeeded when %s already exists", testFolder)
	}
}

// Commands besides srce-init generally shoud raise an error when executed
// outside of a repo
func TestCommandsOutside(t *testing.T) {
	repo := Repo{Dir: testFolder}
	cmds := [](func() error){
		func() error { return repo.Add("HEAD") },
		func() error { return repo.Commit("whatever") },
		repo.Log,
		func() error { return repo.RefLog("HEAD") },
		func() error { _, err := repo.Resolve("HEAD"); return err },
		repo.Status,
		func() error { _, err := repo.GetSymbolicRef("HEAD"); return err },
		func() error { return repo.SetSymbolicRef("HEAD", "whatever") },
		func() error { return repo.UpdateRef("master", "whatever") },
	}
	expected := "not a srce project"
	for _, cmd := range cmds {
		if err := cmd(); err == nil {
			t.Errorf("%s in non-project succeeded", cmd)
		} else if err.Error() != expected {
			t.Errorf("got error %q (expecting %q)", err, expected)
		}
	}
}
