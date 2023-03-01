package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_runner"
)

var recordCmd = &cobra.Command{
	Use:  "record",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Setenv(common.PhaseEnv, string(common.Recording))
		testrunner := test_runner.Testrunner{}

		status, err := test_runner.RunOption.Status()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if run(testrunner, status) {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
