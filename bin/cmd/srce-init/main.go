package main

import (
	"fmt"
	"log"
	"srce"
)

func main() {
	repo := srce.Repo{Dir: srce.DotDir}
	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("srce initialized in %s\n", srce.DotDir)
}
