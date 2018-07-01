package srce

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type RefLog struct {
	path string
}

type RefLogEntry struct {
	sha1Before Hash
	sha1After  Hash
	author     AuthorStamp
	message    string
}

type RefLogEntryOrError struct {
	RefLogEntry
	Err error
}

func parseRefLogEntry(line string) (rle RefLogEntry, err error) {
	pattern := regexp.MustCompile(
		"^([0-9a-f]{40}) ([0-9a-f]{40}) ([^\t]+)\t(.+)$")
	parts := pattern.FindStringSubmatch(line)
	if len(parts) < 5 {
		return rle, fmt.Errorf("malformed reflog entry: %q", line)
	}
	author, err := parseAuthorStamp(parts[3])
	if err != nil {
		return rle, fmt.Errorf("malformed author/timestamp in reflog: %q", parts[3])
	}
	return RefLogEntry{
		sha1Before: Hash(parts[1]), sha1After: Hash(parts[2]), author: author,
		message: parts[4]}, nil
}

func (rle RefLogEntry) toString() string {
	return fmt.Sprintf(
		"%s %s %s\t%s\n", rle.sha1Before, rle.sha1After, rle.author.toString(),
		rle.message)
}

func (r Repo) getRefLog(ref string) RefLog {
	return RefLog{path: r.internalPath(fmt.Sprintf("logs/%s", ref))}
}

func (rl RefLog) read() (<-chan RefLogEntryOrError, error) {
	refLogFile, err := os.Open(rl.path)
	if err != nil {
		return nil, err
	}

	// Return an iterator of RefLogEntries
	ch := make(chan RefLogEntryOrError)
	go func() {
		scanner := bufio.NewScanner(refLogFile)
		for scanner.Scan() {
			rle, err := parseRefLogEntry(scanner.Text())
			if err != nil {
				ch <- RefLogEntryOrError{Err: err}
				return
			}
			ch <- RefLogEntryOrError{RefLogEntry: rle}
		}
		close(ch)
		refLogFile.Close()
	}()
	return ch, nil
}

func (rl RefLog) add(sha1Before, sha1After Hash, author AuthorStamp, message string) error {
	refLogFile, err := os.OpenFile(
		rl.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer refLogFile.Close()

	entry := RefLogEntry{
		sha1Before: sha1Before, sha1After: sha1After, author: author,
		message: message}
	if _, err := refLogFile.WriteString(entry.toString()); err != nil {
		return err
	}
	return nil
}

func (r Repo) RefLog(input string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	// validate/normalize input (e.g. master -> refs/heads/master)
	ref, err := r.expandRef(input)
	if err != nil {
		return err
	}

	entries, err := r.getRefLog(ref).read()
	if err != nil {
		return err
	}

	i := 0
	for rle := range entries {
		if rle.Err != nil {
			return rle.Err
		}
		fmt.Printf(
			"%s %s@{%d}: %s\n", rle.RefLogEntry.sha1After.abbreviated(), ref, i,
			rle.RefLogEntry.message)
		i++
	}

	return nil
}
