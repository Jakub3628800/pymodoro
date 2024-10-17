package cmd

import (
	"fmt"
	"os"
	"td/core"
	"time"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit today's task file",
	Long:  `Open today's task file in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		date := time.Now()
		err := core.OpenEditor(date, 1) // Start at line 1
		if err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
