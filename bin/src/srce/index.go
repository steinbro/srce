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

func parseIndexEntry(line string) IndexEntry {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return IndexEntry{}
	}
	return IndexEntry{
		sha1: Hash(parts[0]), itype: ObjectType(parts[1]), path: parts[2]}
}

func (i IndexEntry) toString() string {
	return fmt.Sprintf("%s %s %s\n", i.sha1, i.itype, i.path)
}

func (r Repo) getIndex() Index {
	return Index{path: r.internalPath("index")}
}

func (i Index) read() (<-chan IndexEntry, error) {
	indexFile, err := os.Open(i.path)
	if err != nil {
		return nil, err
	}

	// Return an iterator of IndexEntries
	ch := make(chan IndexEntry)
	go func() {
		scanner := bufio.NewScanner(indexFile)
		for scanner.Scan() {
			ch <- parseIndexEntry(scanner.Text())
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
	if err != nil {
		return err
	}

	for e := range entries {
		fmt.Printf("M\t%s\n", e.path)
	}

	return nil
}
