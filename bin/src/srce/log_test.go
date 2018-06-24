package srce

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// log against empty repo
	if repo.Log() == nil {
		t.Error("log succeeded in empty repo")
	}

	// log against repo with one commit
	repo.commitSomething(t)
	if err := repo.Log(); err != nil {
		t.Error(err)
	}

	// log against repo with missing commit object
	deletedHash, _ := repo.getLastCommit(t)
	repo.commitSomething(t)
	commitObjPath := repo.internalPath(
		"objects", deletedHash.prefix(), deletedHash.remainder())
	if err := os.Remove(commitObjPath); err != nil {
		t.Fatal(err)
	}
	if repo.Log() == nil {
		t.Error("log succeeded with missing commit")
	}

	// log against repo where commit object has been replaced with junk
	o := Object{otype: CommitObject, sha1: deletedHash, size: 0}
	o.contents.WriteString("nonsense")
	if err := repo.Store(o); err != nil {
		t.Fatal(err)
	}
	if repo.Log() == nil {
		t.Error("log succeeded with malformed commit")
	}
}
