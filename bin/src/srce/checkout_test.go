package srce

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBranch(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	if err := repo.CreateBranch("foo"); err != nil {
		t.Error(err)
	}
	if ref, err := repo.GetSymbolicRef("HEAD"); err != nil {
		t.Error(err)
	} else if ref != "refs/heads/foo" {
		t.Errorf("HEAD = %q (expecting refs/heads/foo)", ref)
	}

	repo.commitSomething(t)
	if err := repo.CreateBranch("foo"); err == nil {
		t.Error("CreateBranch succeeded when branch exists")
	}
	if _, err := repo.Resolve("master"); err == nil {
		t.Error("commit on non-master branch appeared on master")
	}

	lastHash, _ := repo.getLastCommit(t)
	if branchHash, err := repo.Resolve("foo"); err != nil {
		t.Error(err)
	} else if lastHash != branchHash {
		t.Error("non-master branch hash mismatch")
	}
}

func TestCheckoutFile(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// write "hello world" to .srce-test/foo
	originalContents := []byte("hello world")
	testFile := repo.internalPath("foo")
	ioutil.WriteFile(testFile, originalContents, 0644)

	// commit .srce-test/foo
	if err := repo.Add(testFile); err != nil {
		t.Fatal(err)
	}
	if err := repo.Commit("test commit"); err != nil {
		t.Fatal(err)
	}

	// replace contents of .srce-test/foo with "goodbye world"
	if err := ioutil.WriteFile(testFile, []byte("goodbye world"), 0644); err != nil {
		t.Fatal(err)
	}

	// checkout committed version of .srce-test
	if err := repo.CheckoutFile("nonexistent", testFile); err == nil {
		t.Fatal("checked out file from nonexistent branch")
	}
	if err := repo.CheckoutFile("master", "nonexistent"); err == nil {
		t.Fatal("checked out nonexistent file")
	}
	if err := repo.CheckoutFile("master", testFile); err != nil {
		t.Fatal(err)
	}

	if contents, err := ioutil.ReadFile(testFile); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(contents, originalContents) != 0 {
		t.Errorf(
			"post-CheckoutFile %q contents = %q (expecting %q)",
			testFile, contents, originalContents)
	}
}

func TestCheckoutTree(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// write "hello world" to foo/bar (file within a tree)
	originalContents := []byte("hello world")
	testDir := "foo"
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	testFile := filepath.Join(testDir, "bar")
	if err := ioutil.WriteFile(testFile, originalContents, 0644); err != nil {
		t.Fatal(err)
	}

	// commit .srce-test/foo
	if err := repo.Add(testFile); err != nil {
		t.Fatal(err)
	}
	if err := repo.Commit("test commit"); err != nil {
		t.Fatal(err)
	}

	// replace contents of .srce-test/foo with "goodbye world"
	if err := ioutil.WriteFile(testFile, []byte("goodbye world"), 0644); err != nil {
		t.Fatal(err)
	}

	// checkout committed version of .srce-test
	if err := repo.CheckoutTree("nonexistent"); err == nil {
		t.Fatal("checked out file from nonexistent branch")
	}
	if err := repo.CheckoutTree("master"); err != nil {
		t.Fatal(err)
	}

	if contents, err := ioutil.ReadFile(testFile); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(contents, originalContents) != 0 {
		t.Errorf(
			"post-CheckoutFile %q contents = %q (expecting %q)",
			testFile, contents, originalContents)
	}
}
