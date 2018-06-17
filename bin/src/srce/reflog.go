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
	sha1Before string
	sha1After  string
	author     string
	message    string
}

func parseRefLogEntry(line string) RefLogEntry {
	pattern := regexp.MustCompile("([0-9a-f]{40}) ([0-9a-f]{40}) ([^\t]+)\t(.+)")
	parts := pattern.FindStringSubmatch(line)
	if len(parts) < 5 {
		return RefLogEntry{}
	}
	return RefLogEntry{
		sha1Before: parts[1], sha1After: parts[2], author: parts[3],
		message: parts[4]}
}

func (rle RefLogEntry) toString() string {
	return fmt.Sprintf(
		"%s %s %s\t%s\n", rle.sha1Before, rle.sha1After, rle.author, rle.message)
}

func (r Repo) getRefLog(ref string) RefLog {
	return RefLog{path: r.internalPath(fmt.Sprintf("logs/%s", ref))}
}

func (rl RefLog) read() (<-chan RefLogEntry, error) {
	refLogFile, err := os.Open(rl.path)
	if err != nil {
		return nil, err
	}

	// Return an iterator of RefLogEntries
	ch := make(chan RefLogEntry)
	go func() {
		scanner := bufio.NewScanner(refLogFile)
		for scanner.Scan() {
			ch <- parseRefLogEntry(scanner.Text())
		}
		close(ch)
		refLogFile.Close()
	}()
	return ch, nil
}

func (rl RefLog) add(sha1Before, sha1After, author, message string) error {
	entry := RefLogEntry{
		sha1Before: sha1Before, sha1After: sha1After, author: author,
		message: message}
	refLogFile, err := os.OpenFile(
		rl.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := refLogFile.WriteString(entry.toString()); err != nil {
		return err
	}
	refLogFile.Close()
	return nil
}

func (r Repo) RefLog(ref string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	rl := r.getRefLog(ref)
	entries, err := rl.read()
	if err != nil {
		return err
	}

	i := 0
	for rle := range entries {
		fmt.Printf("%s %s@{%d}: %s\n", rle.sha1After[:8], ref, i, rle.message)
		i++
	}

	return nil
}
