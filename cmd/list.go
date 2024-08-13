/*
Copyright © 2024 Jakub Kriz
*/
package cmd

import (
	"fmt"
	"os"
	"pomodoro/vault"
	"time"

	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
	filename string
	date     time.Time
}

func initialModel() model {
	fname := vault.GetFilename(time.Now(), "td")
	ch, sel, _ := vault.LoadLinesWithSelection(fname)
	return model{
		// Our to-do list is a grocery list
		choices: ch,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: sel,
		filename: fname,
		date:     time.Now(),
	}
}

func (m model) Save() {
	vault.UpdateChoices(m.choices, m.selected, m.filename)
}

func (m *model) Refresh() {
	m.filename = vault.GetFilename(m.date, "td")
	ch, sel, _ := vault.LoadLinesWithSelection(m.filename)
	m.choices = ch
	m.selected = sel
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
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
			m.date = m.date.Add(-24 * time.Hour)
			a := &m
			a.Refresh()

		case "right", "l":
			m.date = m.date.Add(24 * time.Hour)
			a := &m
			a.Refresh()

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "e":
			vault.OpenEditor(m.filename)
			a := &m
			a.Refresh()

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			m.Save()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := m.date.Format("2006-01-02") + " " + m.date.Weekday().String() + "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"
	s += "Press e to edit file.\n"

	// Send the UI for rendering
	return s
}
