package main

import (
	"flag"
	"fmt"
	"log"
	"srce"
)

func main() {
	prettyPtr := flag.Bool("p", false, "pretty-print")
	sizePtr := flag.Bool("s", false, "size")
	typePtr := flag.Bool("t", false, "type")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("usage: cat-file [-t | -s | -p] <name>")
	}

	repo := srce.Repo{Dir: srce.DotDir}
	hash, err := repo.Resolve(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	o, err := repo.Fetch(hash)
	if err != nil {
		log.Fatal(err)
	}

	if *prettyPtr {
		fmt.Print(o.Contents())
	} else if *sizePtr {
		fmt.Println(o.Size())
	} else if *typePtr {
		fmt.Println(o.Type())
	}
}
