package srce

import (
	"fmt"
	"os"
)

func Add(dotDir, path string) error {
	// Check that .srce directory exists
	if _, err := os.Stat(dotDir); os.IsNotExist(err) {
		return fmt.Errorf("not a srce project")
	}

	o, err := blobOject(path)
	if err != nil {
		return err
	}

	if err := o.write(dotDir); err != nil {
		return err
	}

	return nil
}
