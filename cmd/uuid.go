package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/spf13/cobra"
)

var isTrim *bool

// uuidCmd represents the uuid command
var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "generate uuid string",
	Long:  `generate uuid string`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := uuid.GenerateUUID()
		if err != nil {
			fmt.Printf("uuid generate error: %v", err)
			return
		}

		if !*isTrim {
			fmt.Println(id)
			return
		}

		fmt.Println(strings.ReplaceAll(id, "-", ""))
	},
}

func init() {
	rootCmd.AddCommand(uuidCmd)
	isTrim = uuidCmd.PersistentFlags().BoolP("trim", "t", false, "trim uuid '-' char")
}
