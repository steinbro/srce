package srce

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/user"
	"time"
)

type Object struct {
	otype    string
	sha1     string
	size     int
	contents bytes.Buffer
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
	o.sha1 = timestampedHash(committer.Name)

	o.contents.Write([]byte(fmt.Sprintf(
		"tree %s\nauthor %s\n\n%s\n", tree.sha1, committer.Name, message)))

	return o, nil
}
