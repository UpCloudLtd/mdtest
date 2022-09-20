package cmd

import (
	"errors"
	"runtime"

	"github.com/UpCloudLtd/mdtest/testrun"
	"github.com/spf13/cobra"
)

var (
	numberOfJobs int

	rootCmd = &cobra.Command{
		Use:   "mdtest [flags] path ...",
		Short: "A testing tool with markdown testcases",
		Long:  "Tool for combining examples and test cases. Parses markdown files for test steps and uses these to test command line applications.",
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.Flags().IntVarP(&numberOfJobs, "jobs", "j", runtime.NumCPU()*2, "number of jobs to use for executing tests in parallel")
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		params := testrun.RunParameters{
			NumberOfJobs: numberOfJobs,
			OutputTarget: rootCmd.OutOrStdout(),
		}

		res := testrun.Execute(args, params)
		err := testrun.NewRunError(res)
		if err != nil {
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
