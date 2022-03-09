package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "xcli",
	Short: "A generator for xcli based applications",
	Long:  "A generator for xcli based applications",
}

func Execute() error {
	return rootCmd.Execute()
}
