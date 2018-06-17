package main

import (
	"log"
	"os"
	"srce"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: update-ref <ref> <hash>")
	}
	ref := os.Args[1]
	hash, err := srce.ValidateHash(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	repo := srce.Repo{Dir: srce.DotDir}

	fullHash, err := repo.ExpandPartialHash(hash)
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.UpdateRef(ref, fullHash); err != nil {
		log.Fatal(err)
	}
}
