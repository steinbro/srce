package srce

import "os"

type Repo struct {
  Dir string
}

func (r Repo) IsInitialized() bool {
  // Check that .srce directory exists
  _, err := os.Stat(r.Dir)
  return !os.IsNotExist(err)
}
