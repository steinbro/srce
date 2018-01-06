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

	if err := srce.Add(srce.DotDir, path); err != nil {
		log.Fatal(err)
	}
}
