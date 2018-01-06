package srce

import (
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
	contents []byte
}

func blobOject(path string) (Object, error) {
	o := Object{otype: "blob", path: path}

	// Read file contents as byte array
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return o, err
	}
	o.contents = contents

	// Compute SHA1 hash of file
	sha := sha1.New()
	sha.Write(contents)
	o.sha1 = hex.EncodeToString(sha.Sum(nil))

	return o, nil
}

func (o Object) write(dotDir string) error {
	// Create .srce/objects/00/ directory (where 00 is the first 2 bytes of hash)
	blobFolder := filepath.Join(dotDir, "objects", o.sha1[:2])
	if err := os.MkdirAll(blobFolder, 0700); err != nil {
		return err
	}

	// Write compressed file contents to .srce/objects/00/rest_of_hash
	blobPath := filepath.Join(blobFolder, o.sha1[2:])
	blobFile, err := os.OpenFile(blobPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	blobData := zlib.NewWriter(blobFile)
	blobData.Write(o.contents)
	blobData.Close()
	if err := blobFile.Close(); err != nil {
		return err
	}

	// Write "<sha1> blob <path>" to .srce/index
	if err := getIndex(dotDir).add(o.sha1, o.otype, o.path); err != nil {
		return err
	}

	return nil
}
