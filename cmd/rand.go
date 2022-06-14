package cmd

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var (
	from     int64
	to       int64 = math.MaxInt64
	realRand *bool
)

// randCmd represents the rand command
var randCmd = &cobra.Command{
	Use:   "rand",
	Short: "get a random number",
	Long:  `get a random number`,
	Run: func(cmd *cobra.Command, args []string) {
		if !checkRandParam() {
			cmd.Println(NormalParamErrorHint)
			return
		}

		if !*realRand {
			cmd.Println(from + rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(to-from))
			return
		}

		result, err := crand.Int(crand.Reader, big.NewInt(to-from))
		if err != nil {
			cmd.Println("generate real random number failed")
			return
		}
		cmd.Println(result.Int64() + from)
	},
}

func checkRandParam() bool {
	if from < 0 || to < 0 || from > to {
		return false
	}

	return true
}

func init() {
	rootCmd.AddCommand(randCmd)
	randCmd.PersistentFlags().Int64VarP(&from, "from", "f", 0, "generate random number [from, to)")
	randCmd.PersistentFlags().Int64VarP(&to, "to", "t", math.MaxInt64, "generate random number [from, to)")
	realRand = randCmd.PersistentFlags().BoolP("real", "r", false, "generate real random number")
}
