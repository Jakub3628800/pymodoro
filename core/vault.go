package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var vaultLoc string
var intervalMode string //daily, weekly or monthly
var templatePath string

type Task struct {
	Line     string
	Selected bool
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func templateFile() string {
	return vaultLoc + "/" + templatePath
}

func init() {
	vaultLoc = getEnv("TD_VAULT_LOC", ".td")
	intervalMode = getEnv("TD_INTERVAL_MODE", "weekly")
	templatePath = getEnv("TD_TEMPLATE_PATH", ".template") // relative to vault location
}

func getFilename(date time.Time) string {

	year, week := date.ISOWeek()
	month := date.Month().String()

	if intervalMode == "daily" {
		return filepath.Join(vaultLoc, date.Format("2006/January/2"))
	} else if intervalMode == "weekly" {
		return fmt.Sprintf("%s/%d/%s/week%d", vaultLoc, year, month, week)
	}

	return vaultLoc + "/" + strconv.Itoa(year) + "/" + month + "/" + month
}

func GetHeader(date time.Time) string {

	if intervalMode == "daily" {

		return date.Format("2006-01-02") + " " + date.Weekday().String() + "\n\n"
	} else {
		_, week := date.ISOWeek()
		return "Week " + strconv.Itoa(week) + "\n\n"
	}

}

func NextDate(date time.Time) time.Time {
	if intervalMode == "daily" {
		return date.Add(24 * time.Hour)
	} else if intervalMode == "weekly" {

		return date.Add(24 * 7 * time.Hour)
	}
	return date.Add(24 * 7 * 31 * time.Hour) //todo fix this

}

func PreviousDate(date time.Time) time.Time {
	if intervalMode == "daily" {
		return date.Add(-24 * time.Hour)
	} else if intervalMode == "weekly" {

		return date.Add(-24 * 7 * time.Hour)
	}
	return date.Add(24 * 7 * 31 * time.Hour) //todo fix this

}

func openFile(date time.Time) (*os.File, error) {
	filename := getFilename(date)
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if fileExists(templateFile()) {
		cmd := exec.Command("cp", vaultLoc+templatePath, path)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

	}
	return nil
}

// Append line to file corresponding to date
func AddTask(date time.Time, line string) error {
	line = "- [ ] " + line

	filename := getFilename(date)
	if !fileExists(filename) {
		createFile(filename)
	}
	file, err := openFile(date)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the line to the file, followed by a newline character
	_, err = file.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}

func isLineCheckbox(line string) (bool, bool) {
	// matches[1] is the indentation (spaces or tabs)
	// matches[2] is the checkbox status (space for unchecked, 'x' for checked)
	// matches[3] is the content of the line
	selected := false
	pattern := regexp.MustCompile(`^([\t ]*)- \[( |x)\] ?(.*)$`)
	if matches := pattern.FindStringSubmatch(line); matches != nil {
		if matches[2] == "x" {
			selected = true
		}
		return true, selected
	}
	return false, false

}

func linesWithSelection(filename string) ([]Task, error) {
	var tasks []Task

	if !fileExists(filename) {
		if fileExists(templateFile()) {
			return linesWithSelection(templateFile())
		} else {
			return tasks, nil
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Iterate over lines
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || trimmedLine == "- [ ]" || trimmedLine == "- [x]" {
			lineNumber++
			continue // Skip empty lines or lines with only checkbox
		}
		isCheck, selected := isLineCheckbox(line)
		if isCheck {
			tasks = append(tasks, Task{line, selected})
		}
		lineNumber++
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return tasks, fmt.Errorf("error reading file: %v", err)
	}
	return tasks, nil
}

func LoadLinesWithSelection(date time.Time) ([]Task, error) {
	filename := getFilename(date)
	return linesWithSelection(filename)

}

func OpenEditor(date time.Time, lineNumber int) error {
	// Create a command to open the file with vim
	filename := getFilename(date)
	if !fileExists(filename) {
		createFile(filename)
	}
	args := []string{fmt.Sprintf("+%d", lineNumber), filename}
	cmd := exec.Command("vim", args...)

	// Set the command's standard input, output, and error to the current program's ones
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func UpdateTaskStatus(selected bool, taskDescription string, date time.Time) error {
	filename := getFilename(date)
	if !fileExists(filename) {
		createFile(filename)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Split the content into lines
	lines := strings.Split(string(file), "\n")

	// Find and update the line
	lineUpdated := false
	for i, line := range lines {
		if strings.Contains(line, taskDescription) {
			if selected {
				lines[i] = strings.Replace(line, "- [ ]", "- [x]", 1)
			} else {
				lines[i] = strings.Replace(line, "- [x]", "- [ ]", 1)
			}
			lineUpdated = true
			break
		}
	}

	if !lineUpdated {
		return fmt.Errorf("task not found in the file")
	}

	// Join the lines back into a single string
	updatedContent := strings.Join(lines, "\n")

	// Write the updated content back to the file
	err = os.WriteFile(filename, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

//func UpdateChoices(choices []string, selected map[int]struct{}, date time.Time) error {
//filename := getFilename(date)
//if !fileExists(filename) {
//createFile(filename)
//}

//file, err := os.ReadFile(filename)
//if err != nil {
//return fmt.Errorf("error reading file: %w", err)
//}

//lines := strings.Split(string(file), "\n")

//// Process each line
//for i, line := range lines {
//for choiceIndex, choice := range choices {
//if strings.Contains(line, choice) {
//_, isSelected := selected[choiceIndex]
//if isSelected && strings.Contains(line, "[ ]") {
//lines[i] = strings.Replace(line, "[ ]", "[x]", 1)
//} else if !isSelected && strings.Contains(line, "[x]") {
//lines[i] = strings.Replace(line, "[x]", "[ ]", 1)
//}
//}
//}
//}

//// Write the updated content back to the file
//output := strings.Join(lines, "\n")
//err = os.WriteFile(filename, []byte(output), 0644)
//if err != nil {
//return fmt.Errorf("error writing file: %w", err)
//}

//return nil
//}

func containsLine(filename string, searchLine string) (int, error) {

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimPrefix(strings.TrimPrefix(line, "- [ ]"), "- [x]")
		if strings.TrimSpace(trimmedLine) == strings.TrimSpace(searchLine) {
			return lineNumber, nil
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, nil
}

func ContainsLine(date time.Time, searchLine string) (int, error) {
	filename := getFilename(date)

	if !fileExists(filename) {
		return containsLine(templateFile(), searchLine)
	}
	return containsLine(filename, searchLine)
}
