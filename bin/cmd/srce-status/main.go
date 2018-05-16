package main

import (
	"log"
	"srce"
)

func main() {
	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.Status(); err != nil {
		log.Fatal(err)
	}
}
