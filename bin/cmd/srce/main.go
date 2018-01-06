package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Usage: %s <command> [<args>]\n", os.Args[0])
		os.Exit(1)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	subcommand := filepath.Join(filepath.Dir(ex),
		fmt.Sprintf("srce-%s", args[0]))
	if _, err := os.Stat(subcommand); os.IsNotExist(err) {
		log.Fatalf("%s: no such command", args[0])
	}
	if err := syscall.Exec(subcommand, args[1:], os.Environ()); err != nil {
		log.Fatal(err)
	}
}
