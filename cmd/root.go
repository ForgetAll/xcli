package cmd

import "github.com/spf13/cobra"

const (
	NormalParamErrorHint = "param error"
)

var rootCmd = &cobra.Command{
	Use:     "xcli",
	Short:   "A generator for xcli based applications",
	Long:    "A generator for xcli based applications",
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}
