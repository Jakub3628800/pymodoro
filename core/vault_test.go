package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFilename(t *testing.T) {
	// Save the original values
	originalVaultLoc := vaultLoc
	originalIntervalMode := intervalMode

	// Restore the original values after the test
	defer func() {
		vaultLoc = originalVaultLoc
		intervalMode = originalIntervalMode
	}()

	tests := []struct {
		name         string
		vaultLoc     string
		intervalMode string
		date         time.Time
		want         string
	}{
		{
			name:         "Daily mode",
			vaultLoc:     "/test/vault",
			intervalMode: "daily",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/30",
		},
		{
			name:         "Weekly mode",
			vaultLoc:     "/test/vault",
			intervalMode: "weekly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/week35",
		},
		{
			name:         "Monthly mode",
			vaultLoc:     "/test/vault",
			intervalMode: "monthly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/August",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vaultLoc = tt.vaultLoc
			intervalMode = tt.intervalMode
			got := getFilename(tt.date)
			if got != tt.want {
				t.Errorf("getFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHeader(t *testing.T) {
	// Save the original value
	originalIntervalMode := intervalMode

	// Restore the original value after the test
	defer func() {
		intervalMode = originalIntervalMode
	}()

	tests := []struct {
		name         string
		intervalMode string
		date         time.Time
		want         string
	}{
		{
			name:         "Daily mode",
			intervalMode: "daily",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "2024-08-30 Friday\n\n",
		},
		{
			name:         "Weekly mode",
			intervalMode: "weekly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "Week 35\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intervalMode = tt.intervalMode
			got := GetHeader(tt.date)
			if got != tt.want {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the test environment
	vaultLoc = tempDir
	intervalMode = "daily"
	templatePath = "test_template"

	// Create a test template file
	templateContent := "Template content\n"
	err = os.WriteFile(filepath.Join(tempDir, templatePath), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Test adding a task
	testDate := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	testTask := "Test task"
	err = AddTask(testDate, testTask)
	if err != nil {
		t.Errorf("AddTask() error = %v", err)
	}

	// Verify the task was added correctly
	filename := getFilename(testDate)
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := templateContent + "- [ ] " + testTask + "\n"
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the test environment
	vaultLoc = tempDir
	intervalMode = "daily"

	// Create a test file with tasks
	testDate := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	filename := getFilename(testDate)
	initialContent := "- [ ] Task 1\n- [ ] Task 2\n- [ ] Task 3\n"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}
	err = os.WriteFile(filename, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test updating task status
	err = UpdateTaskStatus(true, "Task 2", testDate)
	if err != nil {
		t.Errorf("UpdateTaskStatus() error = %v", err)
	}

	// Verify the task status was updated correctly
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := "- [ ] Task 1\n- [x] Task 2\n- [ ] Task 3\n"
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}
}
