package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	showType int

	showNormal    = 0
	showUnix      = 1
	showUnixMilli = 2

	showSeconds bool
	showMilli   bool
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "get current time",
	Long:  `get current time`,
	Run: func(cmd *cobra.Command, args []string) {
		if showSeconds {
			showType = showUnix
		}
		if showMilli {
			showType = showUnixMilli
		}

		now := time.Now()
		switch showType {
		case showNormal:
			cmd.Println(now.Format(time.RFC3339))
		case showUnix:
			cmd.Println(now.Unix())
		case showUnixMilli:
			cmd.Println(now.UnixMilli())
		}
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
	timeCmd.PersistentFlags().BoolVarP(&showSeconds, "ts", "", false, "show time stamp in seconds")
	timeCmd.PersistentFlags().BoolVarP(&showMilli, "tms", "", false, "show time stamp in milliseconds")
}
