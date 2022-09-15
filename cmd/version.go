package cmd

import (
	"fmt"
	"runtime"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/UpCloudLtd/mdtest/output"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			data := []output.SummaryItem{
				{Key: "Version", Value: globals.Version},
				{Key: "Build date", Value: globals.BuildDate},
				{Key: "Built with", Value: runtime.Version()},
				{Key: "System", Value: runtime.GOOS},
				{Key: "Architecture", Value: runtime.GOARCH},
			}

			fmt.Print(output.SummaryTable(data))
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
