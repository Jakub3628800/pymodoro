/*
Copyright © 2024 Jakub Kriz
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"td/core"
	"time"

	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for today.",
	Long:  "List tasks for today",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

type model struct {
	cursor int // which to-do list item our cursor is pointing at
	tasks  []core.Task
	date   time.Time
}

func initialModel() model {
	tasks, _ := core.LoadLinesWithSelection(time.Now())
	return model{
		tasks: tasks,
		date:  time.Now(),
	}
}

func (m model) Save() {
	//core.UpdateChoices(m.choices, m.selected, m.date)
}

func (m *model) Refresh() {
	tasks, _ := core.LoadLinesWithSelection(time.Now())
	m.tasks = tasks
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "left", "h":
			m.date = core.PreviousDate(m.date)
			a := &m
			a.Refresh()

		case "right", "l":
			m.date = core.NextDate(m.date)
			a := &m
			a.Refresh()

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "e":
			lineNumber, _ := core.ContainsLine(m.date, m.tasks[m.cursor].Line)
			core.OpenEditor(m.date, lineNumber)
			a := &m
			a.Refresh()

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selected := m.tasks[m.cursor].Selected
			if selected {
				m.tasks[m.cursor].Selected = false
				core.UpdateTaskStatus(false, m.tasks[m.cursor].Line, m.date)
			} else {
				m.tasks[m.cursor].Selected = true
				core.UpdateTaskStatus(true, m.tasks[m.cursor].Line, m.date)
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := core.GetHeader(m.date)
	//s := m.date.Format("2006-01-02") + " " + m.date.Weekday().String() + "\n\n"

	// Iterate over our choices
	for i, task := range m.tasks {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.tasks[i].Selected {
			checked = "x" // selected!
		}

		// Render the row
		replaced := strings.ReplaceAll(task.Line, "- [ ]", "")
		replaced = strings.ReplaceAll(replaced, "- [x]", "")
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, replaced)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
