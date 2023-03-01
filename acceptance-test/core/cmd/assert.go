package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/spec_validate"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_runner"
	"time"
)

var assertCmd = &cobra.Command{
	Use:  "assert",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Setenv(common.PhaseEnv, string(common.Asserting))
		testrunner := test_runner.Testrunner{}

		status, err := test_runner.RunOption.Status()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if common.CurrentPhase() == common.Asserting {
			test_report.BuildReport()
		}
		start := time.Now()

		if run(testrunner, status) {
			return
		}

		if common.CurrentPhase() == common.Asserting {
			defer func() {
				if !test_report.GetReport().Statistic.Pass {
					os.Exit(2)
				}
			}()
			duration := time.Since(start)
			test_report.GetReport().SetTotalRunTime(duration)
			if _, err := test_report.GetReport().Gen("."); err != nil {
				fmt.Println(err.Error())
				return
			}
			test_report.GetReport().Print()
		}

		contract()
	},
}

func run(testrunner test_runner.Testrunner, status test_runner.RunnerType) bool {
	switch status {
	case test_runner.ROOT:
		for _, root := range test_runner.RunOption.Roots {
			if err := testrunner.RunByRoot(root); err != nil {
				fmt.Println(err.Error())
				return true
			}
		}
		break
	case test_runner.SUITE:
		for _, suite := range test_runner.RunOption.Suites {
			if err := testrunner.RunBySuite(suite); err != nil {
				fmt.Println(err.Error())
				return true
			}
		}
		break
	case test_runner.CASE:
		for _, caseName := range test_runner.RunOption.Cases {
			if err := testrunner.RunByCase(caseName); err != nil {
				fmt.Println(err.Error())
				return true
			}
		}
		break
	}
	return false
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&test_runner.RunOption.Roots, "roots", "r", []string{}, "specify root names")
	rootCmd.PersistentFlags().StringSliceVarP(&test_runner.RunOption.Suites, "suites", "s", []string{}, "specify suite names")
	rootCmd.PersistentFlags().StringSliceVarP(&test_runner.RunOption.Cases, "cases", "c", []string{}, "specify case names")
	rootCmd.AddCommand(assertCmd)
}

func contract() {
	enable := os.Getenv("CONTRACT_ENABLE")
	owner := os.Getenv("CONTRACT_REPO_OWNER")
	name := os.Getenv("CONTRACT_REPO_NAME")
	branch := os.Getenv("CONTRACT_REPO_BRANCH")
	file := os.Getenv("CONTRACT_REPO_PATH")
	if enable == "true" {
		spec, err := spec_validate.QuerySpec(spec_validate.GithubQuery(), owner, name, fmt.Sprintf("%s:%s", branch, file))
		validator, err := spec_validate.Load("lephora", []byte(spec))
		if err != nil {
			panic(err)
		}

		for _, flow := range spec_validate.FlowPool {
			_ = validator.Validate(flow.Req, flow.Resp)
		}
		validator.Report()
	}
}
