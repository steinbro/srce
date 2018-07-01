package srce

import (
	"os"
	"testing"
)

func TestStatus(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	// status before anything has been added/committed
	if err := repo.Status(); err != nil {
		t.Error(err)
	}

	// status with active index
	repo.Add(repo.internalPath("HEAD"))
	if err := repo.Status(); err != nil {
		t.Error(err)
	}

	// status after cleared index
	repo.commitSomething(t)
	if err := repo.Status(); err != nil {
		t.Error(err)
	}

	// append junk to index
	indexFile, err := os.OpenFile(
		repo.internalPath("index"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := indexFile.WriteString("invalid\n"); err != nil {
		t.Fatal(err)
	}
	indexFile.Close()

	if err := repo.Status(); err == nil {
		t.Error("Index with malformed entry succeeded")
	} else {
		t.Log(err)
	}
}
