package srce

import "fmt"

func (r Repo) Add(path string) error {
	if !r.IsInitialized() {
		return fmt.Errorf("not a srce project")
	}

	o, err := blobOject(path)
	if err != nil {
		return err
	}

	if err := r.writeObject(o); err != nil {
		return err
	}

	// Write "<sha1> blob <path>" to .srce/index
	if err := r.getIndex().add(o.sha1, o.otype, o.path); err != nil {
		return err
	}

	return nil
}
