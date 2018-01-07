package main

import (
	"log"
	"os"
	"srce"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no path specified")
	}
	path := os.Args[1]

	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.Add(path); err != nil {
		log.Fatal(err)
	}
}
