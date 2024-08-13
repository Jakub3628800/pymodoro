/*
Copyright © 2024 Jakub Kriz
*/
package cmd

import (
	"fmt"
	"pomodoro/vault"
	"time"

	"github.com/spf13/cobra"
)

var dateFlag string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add task to today's list.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		date, err := parseDate(dateFlag)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		filename := vault.GetFilename(date, "td")

		contains, _ := vault.FileContainsLine(filename, args[0])
		if contains {
			fmt.Println("This item already exists. Skipping")
		} else {
			vault.AppendLineToFile(filename, args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&dateFlag, "date", "today", "Date (today, tomorrow, yesterday, or YYYY-MM-DD)")
}

func parseDate(input string) (time.Time, error) {
	now := time.Now()
	switch input {
	case "today":
		return now, nil
	case "tomorrow":
		return now.AddDate(0, 0, 1), nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	default:
		return time.Parse("2006-01-02", input)
	}
}
