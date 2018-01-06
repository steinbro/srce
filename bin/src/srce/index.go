package srce

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	path string
}

type IndexEntry struct {
	sha1  string
	itype string
	path  string
}

func parseIndexEntry(line string) IndexEntry {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return IndexEntry{}
	}
	return IndexEntry{sha1: parts[0], itype: parts[1], path: parts[2]}
}

func (i IndexEntry) toBytes() []byte {
	return []byte(fmt.Sprintf("%s %s %s\n", i.sha1, i.itype, i.path))
}

func getIndex(dotDir string) Index {
	return Index{path: filepath.Join(dotDir, "index")}
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

func (i Index) add(hash, itype, path string) error {
  entry := IndexEntry{sha1: hash, itype: itype, path: path}
	indexFile, err := os.OpenFile(i.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := indexFile.Write(entry.toBytes()); err != nil {
		return err
	}
	indexFile.Close()
	return nil
}

func (i Index) clear() error {
	indexFile, err := os.OpenFile(i.path, os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	indexFile.Truncate(0)
	indexFile.Close()
	return nil
}
