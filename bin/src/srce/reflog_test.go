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
	emptyRefLog := RefLogEntry{}

	for input, output := range goodRefLogs {
		if result := parseRefLogEntry(input); result == emptyRefLog {
			t.Errorf("parseRefLogEntry(%q) failed to parse", input)
		} else if result != output {
			t.Errorf("parseRefLogEntry(%q) = %s (expecting %s)", input, result, output)
		}
	}

	for _, input := range badRefLogs {
		if result := parseRefLogEntry(input); result != emptyRefLog {
			t.Errorf("parseRefLogEntry(%q) = %q (expected error)", input, result)
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

func TestMissingRefLog(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	repo.commitSomething(t)
	os.Remove(repo.internalPath("logs", "refs", "heads", "master"))

	if err := repo.RefLog("master"); err == nil {
		t.Error("RefLog with missing reflog succeeded")
	} else {
		t.Log(err)
	}
}
