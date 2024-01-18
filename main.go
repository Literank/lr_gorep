package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Literank/gorep/pkg/grep"
)

func main() {
	// Set custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] pattern [file_path]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	// Optional flags
	countFlag := flag.Bool("c", false, "count\nOnly a count of selected lines is written to standard output.")
	ignoreCaseFlag := flag.Bool("i", false, "ignore-case\nPerform case insensitive matching. By default, it is case sensitive.")
	lineNumberFlag := flag.Bool("n", false, "line-number\nEach output line is preceded by its relative line number in the file, starting at line 1. This option is ignored if -count is specified.")
	recursiveFlag := flag.Bool("r", false, "recursive\nRecursively search subdirectories listed.")
	invertMatchFlag := flag.Bool("v", false, "invert-match\nSelected lines are those not matching any of the specified patterns.")

	flag.Parse()

	// Retrieve positional arguments
	// pattern - The pattern to search for
	// file_path - The path to the file to search in
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Argument `pattern` is required.")
		flag.Usage()
		os.Exit(0)
	}
	pattern, filePath := args[0], ""
	if len(args) > 1 {
		filePath = args[1]
	}

	options := &grep.Options{}
	if *ignoreCaseFlag {
		options.IgnoreCase = true
	}
	if *invertMatchFlag {
		options.InvertMatch = true
	}

	var result grep.MatchResult
	var err error

	if *recursiveFlag && filePath != "" {
		result, err = grep.GrepRecursive(pattern, filePath, options)
		if err != nil {
			log.Fatal("Failed to do recursive grep, error:", err)
		}
	} else {
		result, err = grep.Grep(pattern, filePath, options)
		if err != nil {
			log.Fatal("Failed to grep, error:", err)
		}
	}

	if *countFlag {
		fmt.Println(grep.GrepCount(result))
	} else {
		printResult(result, *lineNumberFlag)
	}
}

func printResult(result grep.MatchResult, lineNumberOption bool) {
	currentFile := ""
	fileCount := len(result)

	for filePath, items := range result {
		for _, item := range items {
			if fileCount > 1 && filePath != currentFile {
				currentFile = filePath
				fmt.Printf("\n%s:\n", filePath)
			}
			if lineNumberOption {
				fmt.Printf("%d: %s\n", item.LineNumber, item.Line)
			} else {
				fmt.Println(item.Line)
			}
		}
	}
}
