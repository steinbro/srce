package main

import (
	"log"
	"os"
	"srce"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no message specified")
	}
	message := os.Args[1]

	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.Commit(message); err != nil {
		log.Fatal(err)
	}
}
