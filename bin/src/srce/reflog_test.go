package srce

import (
	"os"
	"testing"
	"time"
)

var goodRefLogs = map[string]RefLogEntry{
	"0000000000000000000000000000000000000000 444d27d3ae655154a91ce7522bcd24a26baf4599 steinbro 1529293210	commit: first": RefLogEntry{
		sha1Before: Hash("0000000000000000000000000000000000000000"),
		sha1After:  Hash("444d27d3ae655154a91ce7522bcd24a26baf4599"),
		author:     AuthorStamp{user: "steinbro", timestamp: time.Unix(1529293210, 0)},
		message:    "commit: first"},
}

var badRefLogs = []string{
	"1234 abcd author 123	commit: message",
	"0000000000000000000000000000000000000000 444d27d3ae655154a91ce7522bcd24a26baf4599 bad authorstamp	commit: first",
}

func TestParseRefLogEntry(t *testing.T) {
	for input, output := range goodRefLogs {
		if result, err := parseRefLogEntry(input); err != nil {
			t.Errorf("parseRefLogEntry(%q) failed: %q", input, err)
		} else if result != output {
			t.Errorf("parseRefLogEntry(%q) = %v (expecting %v)", input, result, output)
		}
	}

	for _, input := range badRefLogs {
		if result, err := parseRefLogEntry(input); err == nil {
			t.Errorf("parseRefLogEntry(%q) = %v (expected error)", input, result)
		}
	}
}

func TestRefLog(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)
	hash1, _ := repo.getLastCommit(t)

	// check HEAD reflog
	headRefLog := repo.getRefLog("HEAD")
	entries, err := headRefLog.read()
	if err != nil {
		t.Error(err)
	} else {
		for rle := range entries {
			if rle.sha1After != hash1 {
				t.Errorf(
					"HEAD reflog/commit hash mismatch (%q != %q)", rle.sha1After, hash1)
			}
		}
	}

	// check master reflog
	headRefLog = repo.getRefLog("refs/heads/master")
	entries, err = headRefLog.read()
	if err != nil {
		t.Error(err)
	} else {
		for rle := range entries {
			if rle.sha1After != hash1 {
				t.Errorf(
					"master reflog/commit hash mismatch (%q != %q)", rle.sha1After, hash1)
			}
		}
	}
}

func TestRefLogCommand(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)

	// check HEAD reflog
	for _, ref := range []string{"HEAD", "master", "refs/heads/master"} {
		if err := repo.RefLog(ref); err != nil {
			t.Error(err)
		}
	}

	// bad refs
	for _, ref := range []string{"../../../foo", "mazter", "refs/heads/foo"} {
		if err := repo.RefLog(ref); err == nil {
			t.Errorf("RefLog(%q) raised no error", ref)
		}
	}
}

func TestMalformedMissingRefLog(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)

	// append junk to reflog
	refLogFile, err := os.OpenFile(
		repo.internalPath("logs", "refs", "heads", "master"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := refLogFile.WriteString("invalid\n"); err != nil {
		t.Fatal(err)
	}
	refLogFile.Close()
	if err := repo.RefLog("master"); err == nil {
		t.Error("RefLog with malformed reflog succeeded")
	} else {
		t.Log(err)
	}

	// destroy reflog entirely
	os.Remove(repo.internalPath("logs", "refs", "heads", "master"))
	if err := repo.RefLog("master"); err == nil {
		t.Error("RefLog with missing reflog succeeded")
	} else {
		t.Log(err)
	}
}
