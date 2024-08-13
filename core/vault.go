package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const vaultLoc = ".td"
const intervalMode = "daily"

// Get filename for a given time
func GetFilename(date time.Time) string {

	year, week := date.ISOWeek()
	month := date.Month().String()

	if intervalMode == "daily" {
		return filepath.Join(vaultLoc, date.Format("2006/January/2"))
	} else if intervalMode == "weekly" {
		return fmt.Sprintf("%s/%d/%s/week%d", vaultLoc, year, month, week)
	}

	return vaultLoc + "/" + strconv.Itoa(year) + "/" + month + "/" + month
}

// Append line to file
func AppendLineToFile(filename string, line string) error {
	line = "- [ ] " + line
	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func LoadLinesWithSelection(filename string) ([]string, map[int]struct{}, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var lines []string
	selectedMap := make(map[int]struct{})

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		// Check for "- [ ]" or "- [x]" prefix
		if strings.HasPrefix(line, "- [ ]") {
			line = strings.TrimPrefix(line, "- [ ] ")
			lines = append(lines, line)
		} else if strings.HasPrefix(line, "- [x]") {
			line = strings.TrimPrefix(line, "- [x] ")
			lines = append(lines, line)
			// Mark the line as selected
			selectedMap[i] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return lines, selectedMap, nil
}

func OpenEditor(filename string) error {
	// Create a command to open the file with vim
	cmd := exec.Command("vim", filename)

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

func ReplaceCheckbox(filename, line string) (bool, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the file line by line and store lines in a slice
	var lines []string
	found := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine == line {
			// Check if the line starts with "- [ ] "
			if strings.HasPrefix(currentLine, "- [ ] ") {
				// Replace with "- [x]"
				currentLine = strings.Replace(currentLine, "- [ ] ", "- [x] ", 1)
				found = true
			}
		}
		lines = append(lines, currentLine)
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}
	fmt.Println(found)
	// If the line was not found or did not need replacing, return false
	if !found {
		return false, nil
	}

	// Write the modified lines back to the file
	if err := os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return false, err
	}

	return true, nil
}

func UpdateChoices(choices []string, selected map[int]struct{}, filename string) error {
	// Read the file
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	lines := strings.Split(string(file), "\n")

	// Process each line
	for i, line := range lines {
		for choiceIndex, choice := range choices {
			if strings.Contains(line, choice) {
				_, isSelected := selected[choiceIndex]
				if isSelected && strings.Contains(line, "[ ]") {
					lines[i] = strings.Replace(line, "[ ]", "[x]", 1)
				} else if !isSelected && strings.Contains(line, "[x]") {
					lines[i] = strings.Replace(line, "[x]", "[ ]", 1)
				}
			}
		}
	}

	// Write the updated content back to the file
	output := strings.Join(lines, "\n")
	err = os.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// Does file contain line?
func FileContainsLine(filename, searchLine string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimPrefix(strings.TrimPrefix(line, "- [ ]"), "- [x]")
		if strings.TrimSpace(trimmedLine) == strings.TrimSpace(searchLine) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
