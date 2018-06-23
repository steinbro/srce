package srce

import (
	"testing"
	"time"
)

var authorStampTests = map[string]AuthorStamp{
	"steinbro 1": AuthorStamp{user: "steinbro", timestamp: time.Unix(1, 0)},
	"Daniel W. Steinbrook <steinbro@post.harvard.edu> 1529279090": AuthorStamp{
		user:      "Daniel W. Steinbrook <steinbro@post.harvard.edu>",
		timestamp: time.Unix(1529279090, 0)},
}

var badAuthorStamps = []string{
	"no date",
	"really long date 999999999999999999999999",
	"steinbro 1 extra",
}

func TestAuthorStamp(t *testing.T) {
	for input, output := range authorStampTests {
		if result, err := parseAuthorStamp(input); err != nil {
			t.Errorf("parseAuthorStamp(%q) error: %q", input, err)
		} else if result != output {
			t.Errorf("parseAuthorStamp(%q) = %s (expecting %s)", input, result, output)
		}
	}

	for _, input := range badAuthorStamps {
		if result, err := parseAuthorStamp(input); err == nil {
			t.Errorf("parseAuthorStamp(%q) = %q (expected error)", input, result)
		}
	}
}
