package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Literank/gorep/pkg/grep"
)

func main() {
	flag.Parse()

	// Retrieve positional arguments
	// pattern - The pattern to search for
	// file_path - The path to the file to search in
	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("Both pattern and file_path are required")
	}

	result, err := grep.Grep(args[0], args[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range result {
		fmt.Println(line)
	}
}
