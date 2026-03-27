package cmd

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/UpCloudLtd/mdtest/testrun"
	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/spf13/cobra"
)

var (
	env              []string
	name             string
	numberOfJobs     int
	timeout          time.Duration
	junitXML         string
	outputToTerminal bool

	rootCmd = &cobra.Command{
		Use:   "mdtest [flags] path ...",
		Short: "A testing tool with markdown testcases",
		Long:  "Tool for combining examples and test cases. Parses markdown files for test steps and uses these to test command line applications.",
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.Flags().StringArrayVarP(&env, "env", "e", nil, "set environment variables for the test run in `key=value` format, e.g., `-e TARGET=test`. The variables set in the command will override the variables defined in the test files")
	rootCmd.Flags().IntVarP(&numberOfJobs, "jobs", "j", runtime.NumCPU()*2, "number of jobs to use for executing tests in parallel")
	rootCmd.Flags().StringVar(&name, "name", "", "name for the testsuite to be printed into the console and to be used as the testsuite name in JUnit XML report")
	rootCmd.Flags().StringVarP(&junitXML, "junit-xml", "x", "", "generate JUnit XML report to the specified `path`")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 0, "timeout for the test run as a `duration` string, e.g., 1s, 1m, 1h")
	rootCmd.Flags().BoolVar(&outputToTerminal, "output-to-terminal", false, "print output from `sh` blocks to the terminal in real-time. Only available when running tests non-concurrently: use either `--jobs=1` or target a single test file when using this flag.")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if outputToTerminal && numberOfJobs > 1 {
			files, _ := utils.ParseFilePaths(args, 1)
			if len(files) > 1 {
				cmd.SilenceUsage = true
				return fmt.Errorf("--output-to-terminal cannot be used when running tests in parallel. Please set --jobs=1 or target a single test file when using --output-to-terminal flag. Currently --jobs is set to %d and targeting %d files.", numberOfJobs, len(files))
			}
		}

		return nil
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		params := testrun.RunParameters{
			Env:              env,
			Name:             name,
			NumberOfJobs:     numberOfJobs,
			OutputTarget:     rootCmd.OutOrStdout(),
			Timeout:          timeout,
			JUnitXML:         junitXML,
			OutputToTerminal: outputToTerminal,
		}

		res := testrun.Execute(args, params)
		if err := testrun.NewRunError(res); err != nil {
			return err
		}

		return nil
	}
}

func Execute() int {
	err := rootCmd.Execute()
	if err != nil {
		var runerr *testrun.RunError
		isRunerr := errors.As(err, &runerr)
		if isRunerr {
			return runerr.ExitCode()
		}
		return 100
	}

	return 0
}
