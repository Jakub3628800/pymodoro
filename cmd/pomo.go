package cmd

import (
	"fmt"
	"os"
	"strings"
	"td/core"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var pomoCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Start a Pomodoro timer",
	Long:  `Start a Pomodoro timer for focused work sessions. Default duration is 25 minutes.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := initialPomoModel()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pomoCmd)
}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type pomoModel struct {
	progress progress.Model
	start    time.Time
	duration time.Duration
}

func initialPomoModel() pomoModel {
	return pomoModel{
		progress: progress.New(
			progress.WithoutPercentage(),
			progress.WithDefaultGradient(),
		),
		duration: 25 * time.Minute,
		start:    time.Now(),
	}
}

func (m pomoModel) Init() tea.Cmd {
	return tickCmd()
}

func (m pomoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		elapsed := time.Since(m.start)
		if elapsed >= m.duration {
			core.SendNotification("pomo session 25m done", false)
			return m, tea.Quit
		}

		percentage := float64(elapsed) / float64(m.duration)
		progressCmd := m.progress.SetPercent(percentage)
		return m, tea.Batch(tickCmd(), progressCmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m pomoModel) View() string {
	elapsed := time.Since(m.start)
	remaining := m.duration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + fmt.Sprintf("%02d:%02d", minutes, seconds) + "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press 'q' to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
