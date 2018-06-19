package srce

import "testing"

func TestValidateHash(t *testing.T) {
	goodHashes := []string{
		"012b",
		"0cbc6f95d498948b48bfbb0aa294c37349bdff89",
		"feedface",
	}
	for _, input := range goodHashes {
		if _, err := ValidateHash(input); err != nil {
			t.Errorf("Unexpected error for hash %q: %q", input, err)
		}
	}

	badHashes := []string{
		"123",
		"0cbc6f95d498948b48bfbb0aa294c37349bdff89abc",
		"exactlyfortycharactersbutnottherightones",
	}
	for _, input := range badHashes {
		if _, err := ValidateHash(input); err == nil {
			t.Errorf("Bad hash validated: %q", input)
		}
	}
}

func TestExpandPartialHash(t *testing.T) {
	repo := setUp(t)

	input := Hash("deadbeef")
	if result, err := repo.ExpandPartialHash(input); err == nil {
		t.Errorf("Bad partial hash expanded: %q -> %q", input, result)
	}

	// one unambiguous entry
	h := Hash("deadbeeffeedface")
	repo.Store(Object{otype: BlobObject, sha1: h})
	if result, err := repo.ExpandPartialHash(input); err != nil {
		t.Errorf("ExpandPartialHash(%q) failed: %q", input, err)
	} else if result != h {
		t.Errorf("ExpandPartialHash(%q) = %q (expected %q)", input, result, h)
	}

	// now make hash prefix ambiguous
	repo.Store(Object{otype: BlobObject, sha1: Hash("deadbeefbadc0ded")})
	if result, err := repo.ExpandPartialHash(input); err == nil {
		t.Errorf("Ambiguous hash expansion succeeded: %q -> %q", input, result)
	}

	tearDown(t)
}
