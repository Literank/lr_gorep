package grep

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MatchItem represents a match in grep searching
type MatchItem struct {
	LineNumber int
	Line       string
}

// MatchResult represents all matches of all files of a grep search
type MatchResult = map[string][]*MatchItem

// Options struct represents the control options of a grep search
type Options struct {
	CountOnly   bool
	IgnoreCase  bool
	InvertMatch bool
}

func GrepMulti(pattern string, filePaths []string, options *Options) (MatchResult, error) {
	if len(filePaths) == 0 {
		return Grep(pattern, "", options)
	}
	result := make(MatchResult)
	for _, filePath := range filePaths {
		grepResult, err := Grep(pattern, filePath, options)
		if err != nil {
			return nil, err
		}
		for k, v := range grepResult {
			result[k] = v
		}
	}
	return result, nil
}

func Grep(pattern string, filePath string, options *Options) (MatchResult, error) {
	lines, err := readFileLines(filePath)
	if err != nil {
		return nil, err
	}

	var matchingLines []*MatchItem
	patternRegex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// normal grep
	if options == nil {
		matchingLines = filterLines(patternRegex, lines, true)
	} else {
		if options.IgnoreCase {
			patternRegex, err = regexp.Compile("(?i)" + pattern)
			if err != nil {
				return nil, err
			}
		}
		if options.InvertMatch {
			matchingLines = filterLines(patternRegex, lines, false)
		} else {
			matchingLines = filterLines(patternRegex, lines, true)
		}
	}
	return MatchResult{filePath: matchingLines}, nil
}

func GrepCount(result MatchResult) int {
	count := 0
	for _, v := range result {
		count += len(v)
	}
	return count
}

func GrepRecursiveMulti(pattern string, dirPaths []string, options *Options) (MatchResult, error) {
	result := make(MatchResult)
	for _, dirPath := range dirPaths {
		grepResult, err := GrepRecursive(pattern, dirPath, options)
		if err != nil {
			return nil, err
		}
		for k, v := range grepResult {
			result[k] = v
		}
	}
	return result, nil
}

func GrepRecursive(pattern string, directoryPath string, options *Options) (MatchResult, error) {
	results := make(MatchResult)
	err := filepath.Walk(directoryPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			result, grepErr := Grep(pattern, filePath, options)
			if grepErr != nil {
				return grepErr
			}
			results[filePath] = result[filePath]
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func readFileLines(filePath string) ([]string, error) {
	var scanner *bufio.Scanner
	if filePath == "" { // Read from standard input
		scanner = bufio.NewScanner(os.Stdin)
	} else { // Read from the file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Create a scanner to read the file line by line
		scanner = bufio.NewScanner(file)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func filterLines(pattern *regexp.Regexp, lines []string, flag bool) []*MatchItem {
	var filteredLines []*MatchItem
	for lineNumber, line := range lines {
		if flag == pattern.MatchString(line) {
			filteredLines = append(filteredLines, &MatchItem{lineNumber + 1, strings.TrimLeft(line, " \t")})
		}
	}
	return filteredLines
}
