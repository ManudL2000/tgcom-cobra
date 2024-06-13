package modfile

import (
	"bufio"
	"os"
	"testing"
)

var dryRun = false

func TestChangeFile(t *testing.T) {

	t.Run("SingleLine", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		conf := Config{
			Filename:   tmpFile.Name(),
			LineNum:    "2",
			StartLabel: "",
			EndLabel:   "",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     dryRun,
		}
		ChangeFile(conf)

		expected := "Line 1\n// Line 2\nLine 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})
}

func TestSelectCommentChars(t *testing.T) {
	tests := []struct {
		filename      string
		expectedChars string
		shouldErr     bool
	}{
		{"testfile.go", "//", false},
		{"testfile.false", "", true},
	}

	for _, tt := range tests {
		commentChars, _ := selectCommentChars(tt.filename, "")
		if commentChars == "" && !tt.shouldErr {
			t.Errorf("selectCommentChars(%s) error: unexpected empty result", tt.filename)
		} else if commentChars != tt.expectedChars {
			t.Errorf("selectCommentChars(%s) = %v, want %v", tt.filename, commentChars, tt.expectedChars)
		}
	}
}

// Utility functions

func createTempFile(t testing.TB, content string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	return tmpFile, func() { os.Remove(tmpFile.Name()) }
}

func assertFileContent(t testing.TB, filename string, expected string) {
	t.Helper()
	modified, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}
	if string(modified) != expected {
		t.Errorf("Unexpected file content:\nGot:\n%s\nExpected:\n%s", string(modified), expected)
	}
}

func setupTestFile(t testing.TB, filename string, lines []string) {
	t.Helper()
	if err := createTestFile(filename); err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	if err := writeTestContent(filename, lines); err != nil {
		t.Fatalf("Error writing test content: %v", err)
	}
}

func createTestFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return f.Close()
}

func writeTestContent(filename string, lines []string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
