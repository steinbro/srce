package srce

import "testing"

func TestRefLog(t *testing.T) {
	repo := setUp(t)

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

	tearDown(t)
}

func TestRefLogCommand(t *testing.T) {
	repo := setUp(t)

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

	tearDown(t)
}
