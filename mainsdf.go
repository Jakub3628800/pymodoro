package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultSessionDuration time.Duration = time.Duration(25 * time.Minute)

type session struct {
	Start    time.Time `json:"start"`
	Duration int       `json:"duration"`
	Category string    `json:"category"`
}

func printElapsed(d time.Duration) {
	fmt.Printf("\033[1A\033[K")
	fmt.Println(d.Truncate(1 * time.Second))
}

func runSession(duration time.Duration, category string, timerEnabled bool) session {
	startTime := time.Now()
	elapsed := time.Since(startTime)

	if timerEnabled {
		fmt.Println("=============================")
		fmt.Println("=============================")
	}
	for elapsed < time.Duration(duration) {

		if timerEnabled {
			printElapsed(elapsed)
		}
		time.Sleep(100 * time.Millisecond)
		elapsed = time.Since(startTime)
	}
	return session{Start: startTime, Duration: int(duration.Minutes()), Category: category}
}

func loadSessions(filename string) ([]session, error) {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sessions []session
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessions)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return sessions, nil
}

func saveSessions(filename string, sessions []session) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(sessions)
	if err != nil {
		return err
	}

	return nil
}

func sendNotification(msg string, silent bool) {
	if silent {
		fmt.Println(msg)
		fmt.Println()
	}
	err := exec.Command("notify-send", msg).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func ensureFileExistence(date time.Time, vaultLoc string) string {

	year := strconv.Itoa(date.Year())
	month := date.Month().String()

	_ = os.MkdirAll(vaultLoc+"/"+year, os.ModePerm)
	_ = os.MkdirAll(vaultLoc+"/"+year+"/"+month, os.ModePerm)

	day := strconv.Itoa(date.Day())
	filepath := vaultLoc + "/" + year + "/" + month + "/" + day
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("File already exists.")
		} else {
			fmt.Println("Error creating file:", err)
		}
	}
	defer file.Close()
	return filepath

}

func appendLineToFile(filename, line string) error {
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

func printFileContents(filename string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file contents to stdout (os.Stdout)
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}

	return nil
}

func replaceCheckbox(filename, line string) (bool, error) {
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

func openEditor(filename string) error {
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

//func main() {
//	mainDir := "td"
//	filepath := ensureFileExistence(time.Now(), mainDir)
//	fmt.Println(filepath)
//	//linecontent := "this is an example"
//	//appendLineToFile(filepath, linecontent)
//	printFileContents(filepath)
//	replaceCheckbox(filepath, "- [ ] this is an example")
//	printFileContents(filepath)
//	openEditor(filepath)
//	fmt.Println("Hello Hello")
//}

//type model struct {
//choices  []string         // items on the to-do list
//cursor   int              // which to-do list item our cursor is pointing at
//selected map[int]struct{} // which to-do items are selected
//}

//func initialModel() model {

//ch, sel, _ := loadLinesWithSelection("td/2024/August/13")
////return model{
////	// Our to-do list is a grocery list
////	choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

////	// A map which indicates which choices are selected. We're using
////	// the  map like a mathematical set. The keys refer to the indexes
////	// of the `choices` slice, above.
////	selected: make(map[int]struct{}),
////}
//return model{choices: ch, selected: sel}
//}

//func (m model) Init() tea.Cmd {
//// Just return `nil`, which means "no I/O right now, please."
//return nil
//}

//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//switch msg := msg.(type) {

//// Is it a key press?
//case tea.KeyMsg:

//// Cool, what was the actual key pressed?
//switch msg.String() {

//// These keys should exit the program.
//case "ctrl+c", "q":
//return m, tea.Quit

//// The "up" and "k" keys move the cursor up
//case "up", "k":
//if m.cursor > 0 {
//m.cursor--
//}

//// The "down" and "j" keys move the cursor down
//case "down", "j":
//if m.cursor < len(m.choices)-1 {
//m.cursor++
//}

//// The "enter" key and the spacebar (a literal space) toggle
//// the selected state for the item that the cursor is pointing at.
//case "enter", " ":
//_, ok := m.selected[m.cursor]
//if ok {
//delete(m.selected, m.cursor)
//} else {
//m.selected[m.cursor] = struct{}{}
//}
//}
//}

//// Return the updated model to the Bubble Tea runtime for processing.
//// Note that we're not returning a command.
//return m, nil
//}

//func (m model) View() string {
//// The header
//s := "What should we buy at the market?\n\n"

//// Iterate over our choices
//for i, choice := range m.choices {

//// Is the cursor pointing at this choice?
//cursor := " " // no cursor
//if m.cursor == i {
//cursor = ">" // cursor!
//}

//// Is this choice selected?
//checked := " " // not selected
//if _, ok := m.selected[i]; ok {
//checked = "x" // selected!
//}

//// Render the row
//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
//}

//// The footer
//s += "\nPress q to quit.\n"

//// Send the UI for rendering
//return s
//}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

func main() {
	//p := tea.NewProgram(initialModel())
	//if _, err := p.Run(); err != nil {
	//	fmt.Printf("Alas, there's been an error: %v", err)
	//	os.Exit(1)
	//}

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.1)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
