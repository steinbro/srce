package main

import (
	"fmt"
	"log"
	"os"
	"srce"
)

func main() {
	repo := srce.Repo{Dir: srce.DotDir}

	if len(os.Args) == 2 {
		name := os.Args[1]
		ref, err := repo.GetSymbolicRef(name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ref)

	} else if len(os.Args) == 3 {
		name := os.Args[1]
		ref := os.Args[2]
		if err := repo.SetSymbolicRef(name, ref); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal("usage: symbolic-ref <name> [<ref>]")
	}
}
