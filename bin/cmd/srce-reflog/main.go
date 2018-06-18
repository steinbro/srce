package main

import (
	"log"
	"os"
	"srce"
)

func main() {
	ref := "HEAD"
	if len(os.Args) > 1 {
		ref = os.Args[1]
	}

	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.RefLog(ref); err != nil {
		log.Fatal(err)
	}
}
