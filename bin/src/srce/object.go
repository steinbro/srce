package srce

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"
	"time"
)

type Object struct {
	otype    string
	sha1     string
	size     int
	contents bytes.Buffer
}

type Commit struct {
	tree    string
	author  string
	parent  string
	message string
}

func (o Object) Type() string {
	return o.otype
}

func (o Object) Size() int {
	return o.size
}

func (o Object) Contents() string {
	return o.contents.String()
}

func timestampedHash(initial string) string {
	// SHA1 for tree is hash of dir name concatenated with RFC3339 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)
	sha := sha1.New()
	sha.Write([]byte(fmt.Sprintf("%s %s", initial, timestamp)))
	return hex.EncodeToString(sha.Sum(nil))
}

func blobOject(path string) (Object, error) {
	o := Object{otype: "blob"}

	// Read file contents as byte array
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return o, err
	}
	o.contents.Write(contents)

	// Compute SHA1 hash of file
	sha := sha1.New()
	sha.Write(contents)
	o.sha1 = hex.EncodeToString(sha.Sum(nil))

	return o, nil
}

func commitObject(tree Object, parentHash string, message string) (Object, error) {
	o := Object{otype: "commit"}

	// Use current OS user as committer
	committer, err := user.Current()
	if err != nil {
		return o, err
	}
	o.sha1 = timestampedHash(committer.Username)

	o.contents.WriteString(fmt.Sprintf("tree %s\n", tree.sha1))
	if parentHash != "" {
		o.contents.WriteString(fmt.Sprintf("parent %s\n", parentHash))
	}
	o.contents.WriteString(fmt.Sprintf("author %s\n", committer.Username))
	o.contents.WriteString(fmt.Sprintf("\n%s\n", message))

	return o, nil
}

func parseCommit(contents bytes.Buffer) (Commit, error) {
	commit := Commit{}
	scanner := bufio.NewScanner(bytes.NewReader(contents.Bytes()))

	for scanner.Scan() {
		line := scanner.Text()

		// blank line separates header from commit message
		if line == "" {
			for scanner.Scan() {
				commit.message += scanner.Text()
			}
			break // nothing left in commit body
		}

		parts := strings.SplitN(line, " ", 2)
		key, value := parts[0], parts[1]

		if key == "tree" {
			commit.tree = value
		} else if key == "author" {
			commit.author = value
		} else if key == "parent" {
			commit.parent = value
		} else {
			return Commit{}, fmt.Errorf("unrecognized field in commit header: %q", key)
		}
	}

	return commit, nil
}
