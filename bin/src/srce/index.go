package srce

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Index struct {
	path string
}

type IndexEntry struct {
	sha1  Hash
	itype ObjectType
	path  string
}

type IndexEntryOrError struct {
	IndexEntry
	Err error
}

func parseIndexEntry(line string) (ie IndexEntry, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ie, fmt.Errorf("malformed index entry: %q", line)
	}
	return IndexEntry{
		sha1: Hash(parts[0]), itype: ObjectType(parts[1]), path: parts[2]}, nil
}

func (i IndexEntry) toString() string {
	return fmt.Sprintf("%s %s %s\n", i.sha1, i.itype, i.path)
}

func (r Repo) getIndex() Index {
	return Index{path: r.internalPath("index")}
}

func (i Index) read() (<-chan IndexEntryOrError, error) {
	indexFile, err := os.Open(i.path)
	if err != nil {
		return nil, err
	}

	// Return an iterator of IndexEntries
	ch := make(chan IndexEntryOrError)
	go func() {
		scanner := bufio.NewScanner(indexFile)
		for scanner.Scan() {
			ie, err := parseIndexEntry(scanner.Text())
			if err != nil {
				ch <- IndexEntryOrError{Err: err}
				return
			}
			ch <- IndexEntryOrError{IndexEntry: ie}
		}
		close(ch)
		indexFile.Close()
	}()
	return ch, nil
}

func (i Index) add(sha1 Hash, itype ObjectType, path string) (err error) {
	indexFile, err := os.OpenFile(i.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer indexFile.Close()

	entry := IndexEntry{sha1: sha1, itype: itype, path: path}
	_, err = indexFile.WriteString(entry.toString())
	return
}

func (i Index) clear() error {
	indexFile, err := os.OpenFile(i.path, os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	indexFile.Truncate(0)
	return nil
}

func (r Repo) Status() error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	entries, err := r.getIndex().read()
	if os.IsNotExist(err) {
		return nil // missing index, equivalent to empty index
	} else if err != nil {
		return err
	}

	for e := range entries {
		if e.Err != nil {
			return e.Err
		}
		fmt.Printf("M\t%s\n", e.IndexEntry.path)
	}

	return nil
}
