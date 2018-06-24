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

const initialCommitHash = Hash("0000000000000000000000000000000000000000")

type ObjectType string

const (
	BlobObject   ObjectType = "blob"
	CommitObject ObjectType = "commit"
	TreeObject   ObjectType = "tree"
)

type Object struct {
	otype    ObjectType
	sha1     Hash
	size     int
	contents bytes.Buffer
}

type Commit struct {
	tree    Hash
	parent  Hash
	author  AuthorStamp
	message string
}

func (o Object) Type() ObjectType {
	return o.otype
}

func (o Object) Size() int {
	return o.size
}

func (o Object) Contents() string {
	return o.contents.String()
}

func timestampedHash(initial string) Hash {
	// SHA1 for tree is hash of dir name concatenated with RFC3339 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)
	sha := sha1.New()
	sha.Write([]byte(fmt.Sprintf("%s %s", initial, timestamp)))
	return Hash(hex.EncodeToString(sha.Sum(nil)))
}

func blobOject(path string) (Object, error) {
	o := Object{otype: BlobObject}

	// Read file contents as byte array
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return o, err
	}
	o.contents.Write(contents)

	// Compute SHA1 hash of file
	sha := sha1.New()
	sha.Write(contents)
	o.sha1 = Hash(hex.EncodeToString(sha.Sum(nil)))

	return o, nil
}

func commitObject(tree Object, parentHash Hash, message string) (Object, error) {
	o := Object{otype: CommitObject}

	// Use current OS user as committer
	committer, err := user.Current()
	if err != nil {
		return o, err
	}

	authorstamp := AuthorStamp{user: committer.Username, timestamp: time.Now()}

	o.contents.WriteString(fmt.Sprintf("tree %s\n", tree.sha1))
	if parentHash != initialCommitHash {
		o.contents.WriteString(fmt.Sprintf("parent %s\n", parentHash))
	}
	o.contents.WriteString(fmt.Sprintf("author %s\n", authorstamp.toString()))
	o.contents.WriteString(fmt.Sprintf("\n%s\n", message))

	// Compute SHA1 hash of commit contents
	sha := sha1.New()
	sha.Write(o.contents.Bytes())
	o.sha1 = Hash(hex.EncodeToString(sha.Sum(nil)))

	return o, nil
}

func parseCommit(contents bytes.Buffer) (commit Commit, err error) {
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
		if len(parts) != 2 {
			return commit, fmt.Errorf("malformed commit data: %q", line)
		}
		key, value := parts[0], parts[1]

		if key == "tree" {
			if commit.tree, err = ValidateHash(value); err != nil {
				return commit, err
			}
		} else if key == "author" {
			if commit.author, err = parseAuthorStamp(value); err != nil {
				return commit, err
			}
		} else if key == "parent" {
			if commit.parent, err = ValidateHash(value); err != nil {
				return commit, err
			}
		} else {
			return commit, fmt.Errorf("unrecognized field in commit header: %q", key)
		}
	}

	return commit, nil
}
