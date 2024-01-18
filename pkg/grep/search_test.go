package grep

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGrep(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write some content to the temporary file
	content := "line1\nline2\nline3\npattern\nline4"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	// Close the file before testing
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test Grep function
	pattern := "pattern"
	options := &Options{}
	result, err := Grep(pattern, tmpfile.Name(), options)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := MatchResult{
		tmpfile.Name(): {
			{LineNumber: 4, Line: "pattern"},
		},
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, result)
	}
}

func TestGrepCount(t *testing.T) {
	// Test cases for GrepCount function
	t.Run("EmptyResult", func(t *testing.T) {
		result := make(MatchResult)
		count := GrepCount(result)
		expectedCount := 0
		if count != expectedCount {
			t.Errorf("Expected count %v, but got %v", expectedCount, count)
		}
	})

	t.Run("NonEmptyResult", func(t *testing.T) {
		result := MatchResult{
			"file1.txt": {
				{LineNumber: 1, Line: "pattern"},
				{LineNumber: 5, Line: "pattern"},
			},
			"file2.txt": {
				{LineNumber: 3, Line: "pattern"},
			},
		}
		count := GrepCount(result)
		expectedCount := 3
		if count != expectedCount {
			t.Errorf("Expected count %v, but got %v", expectedCount, count)
		}
	})
}

func TestGrepRecursive(t *testing.T) {
	// Create a temporary directory for testing
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create nested files and directories
	files := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
		"subdir/file4.txt",
	}

	for _, file := range files {
		filePath := filepath.Join(tmpdir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			t.Fatal(err)
		}

		// Write some content to the files
		content := "pattern"
		if err := os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}

	// Test GrepRecursive function
	pattern := "pattern"
	options := &Options{}
	result, err := GrepRecursive(pattern, tmpdir, options)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := MatchResult{
		filepath.Join(tmpdir, "file1.txt"): {
			{LineNumber: 1, Line: "pattern"},
		},
		filepath.Join(tmpdir, "file2.txt"): {
			{LineNumber: 1, Line: "pattern"},
		},
		filepath.Join(tmpdir, "subdir/file3.txt"): {
			{LineNumber: 1, Line: "pattern"},
		},
		filepath.Join(tmpdir, "subdir/file4.txt"): {
			{LineNumber: 1, Line: "pattern"},
		},
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, result)
	}
}
