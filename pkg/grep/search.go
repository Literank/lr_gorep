package grep

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func Grep(pattern string, filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matchingLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match, err := regexp.MatchString(pattern, line)
		if err != nil {
			return nil, err
		}
		if match {
			matchingLines = append(matchingLines, strings.TrimLeft(line, " \t"))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}
