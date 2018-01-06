package main

import (
	"fmt"
	"log"
	"srce"
)

func main() {
	if err := srce.Init(srce.DotDir); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("srce initialized in %s\n", srce.DotDir)
}
