package srce

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AuthorStamp struct {
	user      string
	timestamp time.Time
}

func parseAuthorStamp(input string) (a AuthorStamp, err error) {
	userAndDate := strings.SplitN(input, " ", 2)
	a.user = userAndDate[0]
	timestamp, err := strconv.ParseInt(userAndDate[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid commit timestamp: %q", userAndDate[1])
		return
	}
	a.timestamp = time.Unix(timestamp, 0)
	return
}

func (a AuthorStamp) toString() string {
	return fmt.Sprintf("%s %d", a.user, a.timestamp.Unix())
}
