package main

import (
	"fmt"
	"log"
	"os"
	"srce"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: rev-parse <name>")
	}
	name := os.Args[1]

	repo := srce.Repo{Dir: srce.DotDir}
	if hash, err := repo.Resolve(name); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(hash)
	}
}
