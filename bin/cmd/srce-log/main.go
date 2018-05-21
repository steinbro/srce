package main

import (
	"log"
	"srce"
)

func main() {
	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.Log(); err != nil {
		log.Fatal(err)
	}
}
