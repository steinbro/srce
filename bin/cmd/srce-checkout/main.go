package main

import (
	"flag"
	"fmt"
	"log"
	"srce"
)

func main() {
	branch := flag.String("b", "", "branch")
	flag.Parse()

	repo := srce.Repo{Dir: srce.DotDir}

	if len(*branch) > 0 { // create a new branch, and switch to it
		if err := repo.CreateBranch(*branch); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Switched to a new branch %q\n", *branch)
		}

	} else if len(flag.Args()) == 1 { // checkout whole tree
		if err := repo.CheckoutTree(flag.Args()[0]); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Switched to branch %q\n", flag.Args()[0])
		}

	} else if len(flag.Args()) == 2 { // checkout a single file
		if err := repo.CheckoutFile(flag.Args()[0], flag.Args()[1]); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal("usage: checkout [-b <new_brnach>] <ref> [<path>]")
	}
}
