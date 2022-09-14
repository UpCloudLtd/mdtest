package cmd

import (
	"fmt"
	"runtime"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(""+
				"Version:      %s\n"+
				"Go version:   %s\n"+
				"Build date:   %s\n"+
				"System:       %s\n"+
				"Architecture: %s\n",
				globals.Version,
				runtime.Version(),
				globals.BuildDate,
				runtime.GOOS,
				runtime.GOARCH)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
