package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "get current time",
	Long:  `get current time`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(time.Now().Format(time.RFC3339))
		// cmd.Println(time.Now().Format("1"))
		// cmd.Println(time.Now().Format("2006"))
		// cmd.Println(time.Now().Format("2"))
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
