package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.0.2"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of xcli",
	Long:  `All software has versions. This is xcli's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xcli", version)
		cmd.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
