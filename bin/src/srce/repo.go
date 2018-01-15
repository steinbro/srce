package srce

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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
	objFolder := r.internalPath("objects", o.sha1[:2])
	if err := os.MkdirAll(objFolder, 0700); err != nil {
		return err
	}

	// Write file contents to .srce/objects/00/rest_of_hash
	objPath := r.internalPath("objects", o.sha1[:2], o.sha1[2:])
	objFile, err := os.OpenFile(objPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer objFile.Close()

	// Store compressed contents
	w := zlib.NewWriter(objFile)
	// Add header e.g. "blob 16", terminated by null byte
	w.Write([]byte(fmt.Sprintf("%s %d", o.otype, o.contents.Len())))
	w.Write([]byte("\u0000"))
	o.contents.WriteTo(w)
	w.Close()

	return nil
}

func (r Repo) parseObject(buf *bytes.Buffer) (Object, error) {
	var o Object

	// Read raw data until null byte to populate header
	header := new(bytes.Buffer)
	for {
		if c, _, err := buf.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return o, err
			}
		} else if c == '\u0000' {
			break // end of header
		} else {
			header.WriteRune(c) // append rune to header
		}
	}

	// Extract object type and size from header
	pattern, _ := regexp.Compile("^([a-z]+) ([0-9]+)$")
	metadata := pattern.FindStringSubmatch(header.String())
	if len(metadata) != 3 {
		return o, fmt.Errorf("malformed header: %q", header.String())
	}
	o.otype = metadata[1]
	o.size, _ = strconv.Atoi(metadata[2])

	// Remaining raw data is uncompressed object contents
	io.Copy(&o.contents, buf)

	return o, nil
}

func (r Repo) Fetch(sha1 string) (Object, error) {
	f, err := os.Open(r.internalPath("objects", sha1[:2], sha1[2:]))
	if err != nil {
		return Object{}, err
	}
	defer f.Close()

	// Copy decompressed data into buffer
	buf := new(bytes.Buffer)
	z, err := zlib.NewReader(f)
	if err != nil {
		return Object{}, err
	}
	io.Copy(buf, z)

	return r.parseObject(buf)
}
