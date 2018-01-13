package srce

import (
	"compress/zlib"
	"os"
	"path/filepath"
)

type Repo struct {
	Dir string
}

func (r Repo) IsInitialized() bool {
	// Check that .srce directory exists
	_, err := os.Stat(r.Dir)
	return !os.IsNotExist(err)
}

func (r Repo) Store(o Object) error {
	// Create .srce/objects/00/ directory (where 00 is the first 2 bytes of hash)
	objFolder := filepath.Join(r.Dir, "objects", o.sha1[:2])
	if err := os.MkdirAll(objFolder, 0700); err != nil {
		return err
	}

	// Write file contents to .srce/objects/00/rest_of_hash
	objPath := filepath.Join(objFolder, o.sha1[2:])
	objFile, err := os.OpenFile(objPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer objFile.Close()

	// Store compressed contents
	w := zlib.NewWriter(objFile)
	o.contents.WriteTo(w)
	w.Close()

	return nil
}
