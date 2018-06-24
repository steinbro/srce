package srce

import (
	"bytes"
	"testing"
)

var badHeaderTests = []string{
	"blob 16 32\u0000something",
	"no size\u0000something",
	"4\u0000something",
	"\u0000else",
	"whatever\u0000",
	"\u0000",
	"",
}

func TestParseObject(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	for _, tt := range badHeaderTests {
		buf := bytes.NewBufferString(tt)
		if _, err := repo.parseObject(buf); err == nil {
			t.Errorf("%q did not raise an error", tt)
		}
	}
}
