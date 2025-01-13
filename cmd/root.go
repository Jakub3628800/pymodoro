package cmd

import (
	"fmt"
	"os"
	"strings"
	"td/core"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "A simple, efficient Text User Interface (TUI) app for tracking tasks",
	Long: `To-Do ToDay (td) is a simple, efficient Text User Interface (TUI) app for tracking tasks 
with a focus on daily workflow. Seamlessly add and check off tasks while the backend 
stores your progress in easy-to-read markdown files.

Features:
- 📝 Quick task addition and management
- ✅ Simple checkbox-style task completion
- 📁 Markdown file storage for easy version control and portability
- 📆 Daily, weekly, and monthly view options
- 🖥️ Clean and intuitive TUI for distraction-free productivity`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type model struct {
	cursor int
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
	// Implement save functionality if needed
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "e":
			lineNumber, _ := core.ContainsLine(m.date, m.tasks[m.cursor].Line)
			core.OpenEditor(m.date, lineNumber, false) // Add false as the third argument
			a := &m
			a.Refresh()
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
	return m, nil
}

func (m model) View() string {
	s := core.GetHeader(m.date)

	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.tasks[i].Selected {
			checked = "x"
		}

		replaced := strings.ReplaceAll(task.Line, "- [ ]", "")
		replaced = strings.ReplaceAll(replaced, "- [x]", "")
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, replaced)
	}

	// Use the existing helpStyle from pomo.go
	s += "\n" + helpStyle("Press q to quit.")

	return s
}
