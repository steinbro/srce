package main

import (
	"log"
	"os"
	"srce"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: update-ref <ref> <hash> ")
	}
	ref := os.Args[1]
  hash := os.Args[2]

	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.UpdateRef(ref, hash); err != nil {
		log.Fatal(err)
	}
}
