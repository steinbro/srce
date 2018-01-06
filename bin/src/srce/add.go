package srce

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Add(dotDir, path string) error {
	// Check that .srce directory exists
	if _, err := os.Stat(dotDir); os.IsNotExist(err) {
		return fmt.Errorf("not a srce project")
	}

	// Read file contents as byte array
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Compute SHA1 hash of file
	sha := sha1.New()
	sha.Write(contents)
	shaHex := hex.EncodeToString(sha.Sum(nil))

	// Create .srce/objects/00/ directory (where 00 is the first 2 bytes of hash)
	blobFolder := filepath.Join(dotDir, "objects", shaHex[:2])
	if err := os.MkdirAll(blobFolder, 0700); err != nil {
		return err
	}

	// Write compressed file contents to .srce/objects/00/rest_of_hash
	blobPath := filepath.Join(blobFolder, shaHex[2:])
	blobFile, err := os.OpenFile(blobPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	blobData := zlib.NewWriter(blobFile)
	blobData.Write(contents)
	blobData.Close()
	if err := blobFile.Close(); err != nil {
		return err
	}

	// Write "<sha1> blob <path>" to .srce/index
	indexLine := fmt.Sprintf("%s blob %s\n", shaHex, path)
	indexPath := filepath.Join(dotDir, "index")
	indexFile, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := indexFile.Write([]byte(indexLine)); err != nil {
		return err
	}
	if err := indexFile.Close(); err != nil {
		return err
	}

	return nil
}
