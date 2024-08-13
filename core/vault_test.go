package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFilename(t *testing.T) {
	date := time.Date(2023, time.August, 15, 0, 0, 0, 0, time.UTC)
	expected := ".td/2023/August/15"
	result := GetFilename(date)
	if result != expected {
		t.Errorf("GetFilename() = %v, want %v", result, expected)
	}
}

func TestAppendLineToFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "testfile")

	err = AppendLineToFile(tempFile, "Test line")
	if err != nil {
		t.Fatalf("AppendLineToFile() error = %v", err)
	}

	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "- [ ] Test line\n"
	if string(content) != expected {
		t.Errorf("File content = %v, want %v", string(content), expected)
	}
}

func TestLoadLinesWithSelection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "testfile")

	content := `- [ ] Task 1
- [x] Task 2
- [ ] Task 3`
	err = os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	lines, selected, err := LoadLinesWithSelection(tempFile)
	if err != nil {
		t.Fatalf("LoadLinesWithSelection() error = %v", err)
	}

	expectedLines := []string{"Task 1", "Task 2", "Task 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("LoadLinesWithSelection() returned %d lines, want %d", len(lines), len(expectedLines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d = %v, want %v", i, line, expectedLines[i])
		}
	}

	if _, ok := selected[1]; !ok {
		t.Errorf("Task 2 should be selected")
	}
}

func TestFileContainsLine(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "testfile")

	content := `- [ ] Task 1
- [x] Task 2
- [ ] Task 3`
	err = os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Existing unchecked task", "Task 1", true},
		{"Existing checked task", "Task 2", true},
		{"Non-existing task", "Task 4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FileContainsLine(tempFile, tt.line)
			if err != nil {
				t.Fatalf("FileContainsLine() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("FileContainsLine() = %v, want %v", result, tt.expected)
			}
		})
	}
}
