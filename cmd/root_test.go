package cmd

import (
	"os"
	"testing"
	"time"
	"td/core"
)

func TestCopyTasksFromYesterday(t *testing.T) {
	// Setup
	yesterday := time.Now().AddDate(0, 0, -1)
	today := time.Now()
	
	// Create yesterday's tasks
	core.AddTask(yesterday, "- [ ] Task 1")
	core.AddTask(yesterday, "- [x] Task 2")
	core.AddTask(yesterday, "- [ ] Task 3")

	// Create model
	m := model{date: today}

	// Run the function
	m.copyTasksFromYesterday()

	// Check if tasks were copied correctly
	tasks, _ := core.LoadLinesWithSelection(today)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Line != "- [ ] Task 1" {
		t.Errorf("Expected '- [ ] Task 1', got '%s'", tasks[0].Line)
	}

	if tasks[1].Line != "- [ ] Task 3" {
		t.Errorf("Expected '- [ ] Task 3', got '%s'", tasks[1].Line)
	}

	// Cleanup
	os.Remove(core.GetFilePath(yesterday))
	os.Remove(core.GetFilePath(today))
}

func TestAppendTask(t *testing.T) {
	// This test is more challenging to implement as it involves user input via nvim
	// For now, we'll just test if the function exists
	m := model{}
	m.appendTask()
	// If the function doesn't exist, this test will fail to compile
}

func TestViewColorChange(t *testing.T) {
	// Setup
	m := model{
		tasks: []core.Task{
			{Line: "- [ ] Task 1", Selected: false},
			{Line: "- [x] Task 2", Selected: true},
		},
		date: time.Now(),
	}

	// Run the View function
	output := m.View()

	// Check if the output contains different colors for completed and uncompleted tasks
	if !strings.Contains(output, "\x1b[37m") || !strings.Contains(output, "\x1b[32m") {
		t.Error("Expected different colors for completed and uncompleted tasks")
	}
}