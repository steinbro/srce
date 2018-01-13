package srce

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/user"
	"regexp"
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

func commitObject(tree Object, message string) (Object, error) {
	o := Object{otype: "commit"}

	// Use current OS user as committer
	committer, err := user.Current()
	if err != nil {
		return o, err
	}
	o.sha1 = timestampedHash(committer.Username)

	o.contents.WriteString(fmt.Sprintf(
		"tree %s\nauthor %s\n\n%s\n", tree.sha1, committer.Username, message))

	return o, nil
}

func (r Repo) parseCommit(contents bytes.Buffer) (Commit, error) {
	pattern, _ := regexp.Compile(
		"^tree ([0-9a-f]{40})\nauthor ([^\n]+)\n\n(.+)\n$")
	m := pattern.FindStringSubmatch(contents.String())
	if len(m) != 4 {
		return Commit{}, fmt.Errorf("malformed commit: %q", contents.String())
	}
	return Commit{tree: m[1], author: m[2], message: m[3]}, nil
}
