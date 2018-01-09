package srce

import (
  "fmt"
  "os"
  "path/filepath"
)

func (r Repo) UpdateRef(ref, hash string) error {
  // write hash to e.g. .srce/refs/heads/master
  refFile, err := os.OpenFile(
    filepath.Join(r.Dir, ref), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    return err
  }
  if _, err := refFile.Write([]byte(fmt.Sprintf("%s\n", hash))); err != nil {
    return err
  }
  refFile.Close()
  return nil
}
