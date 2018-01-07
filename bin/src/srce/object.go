package srce

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Object struct {
	otype    string
	path     string
	sha1     string
	contents bytes.Buffer
}

func blobOject(path string) (Object, error) {
	o := Object{otype: "blob", path: path}

	// Read file contents as byte array
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return o, err
	}

	// Compute SHA1 hash of file
	sha := sha1.New()
	sha.Write(contents)
	o.sha1 = hex.EncodeToString(sha.Sum(nil))

	// Store compressed contents
	w := zlib.NewWriter(&o.contents)
	w.Write(contents)
	w.Close()

	return o, nil
}

func (r Repo) writeObject(o Object) error {
	// Create .srce/objects/00/ directory (where 00 is the first 2 bytes of hash)
	blobFolder := filepath.Join(r.Dir, "objects", o.sha1[:2])
	if err := os.MkdirAll(blobFolder, 0700); err != nil {
		return err
	}

	// Write file contents to .srce/objects/00/rest_of_hash
	blobPath := filepath.Join(blobFolder, o.sha1[2:])
	if err := ioutil.WriteFile(blobPath, o.contents.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
