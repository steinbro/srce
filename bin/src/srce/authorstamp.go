package srce

import (
	"fmt"
	"os/user"
	"regexp"
	"strconv"
	"time"
)

type AuthorStamp struct {
	user      string
	timestamp time.Time
}

func parseAuthorStamp(input string) (a AuthorStamp, err error) {
	pattern := regexp.MustCompile("^(.+) ([0-9]+)$")
	userAndDate := pattern.FindStringSubmatch(input)
	if len(userAndDate) < 3 {
		err = fmt.Errorf("invalid author stamp: %q", input)
		return
	}
	a.user = userAndDate[1]
	timestamp, err := strconv.ParseInt(userAndDate[2], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid commit timestamp: %q", userAndDate[2])
		return
	}
	a.timestamp = time.Unix(timestamp, 0)
	return
}

func (a AuthorStamp) toString() string {
	return fmt.Sprintf("%s %d", a.user, a.timestamp.Unix())
}

func currentUserAndTime() AuthorStamp {
	user, _ := user.Current()
	return AuthorStamp{user: user.Username, timestamp: time.Now()}
}
