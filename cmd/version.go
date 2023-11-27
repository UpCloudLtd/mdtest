package cmd

import (
	"fmt"
	"runtime"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/UpCloudLtd/mdtest/output"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		data := []output.SummaryItem{
			{Key: "Version", Value: globals.GetVersion()},
			{Key: "Build date", Value: globals.BuildDate},
			{Key: "Built with", Value: runtime.Version()},
			{Key: "System", Value: runtime.GOOS},
			{Key: "Architecture", Value: runtime.GOARCH},
		}

		// Omit unknown build date from the output
		if globals.BuildDate == "unknown" {
			data = append(data[:1], data[2:]...)
		}

		fmt.Fprint(versionCmd.OutOrStdout(), output.SummaryTable(data))
	}
}
