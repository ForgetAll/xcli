package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/spf13/cobra"
)

const (
	UUIDLength = 36
	UUIDTrim   = 32
)

var (
	isTrim *bool
	count  int8 = 36
)

// uuidCmd represents the uuid command
var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "generate uuid string",
	Long:  `generate uuid string`,
	Run: func(cmd *cobra.Command, args []string) {
		if !checkParam() {
			cmd.Println(NormalParamErrorHint)
			return
		}

		id, err := uuid.GenerateUUID()
		if err != nil {
			fmt.Printf("uuid generate error: %v", err)
			return
		}

		fmt.Println(handleUID(id))
	},
}

func checkParam() bool {
	if count < 0 || count > UUIDLength || (*isTrim && count > UUIDTrim) {
		return false
	}

	return true
}

func handleUID(uid string) string {
	if !*isTrim {
		return uid[:count]
	}

	id := strings.ReplaceAll(uid, "-", "")
	return id[:count]
}

func init() {
	rootCmd.AddCommand(uuidCmd)
	isTrim = uuidCmd.PersistentFlags().BoolP("trim", "t", false, "trim uuid '-' char")
	uuidCmd.PersistentFlags().Int8VarP(&count, "count", "c", UUIDLength, "return count length uuid")
}
