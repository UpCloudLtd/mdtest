package cmd

import (
	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/spf13/cobra"
)

var (
	transforms  []string
	outputPath  string
	quoteValues string

	normaliseCmd = &cobra.Command{
		Aliases: []string{"normalize"},
		Use:     "normalise",
		Short:   "Normalise the fenced code block info texts",
		Long:    "Normalise the fenced code block info texts. By default, removes all info texts defined after the language identifier from the starting code-block fence.",
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(normaliseCmd)
	normaliseCmd.Flags().StringArrayVarP(&transforms, "transform", "t", nil, "transform info text key in `old=new` format, e.g., `-t filename=title` would transform `filename` info text to `title` info text")
	normaliseCmd.Flags().StringVar(&quoteValues, "quote-values", "no", "`when` to quote info text values. `always` or `no`")
	normaliseCmd.Flags().StringVarP(&outputPath, "output", "o", "", "`directory` where to save the normalised files")
	_ = normaliseCmd.MarkFlagRequired("output")
	normaliseCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		params := utils.NormalizeParameters{
			OutputPath: outputPath,
			Transforms: transforms,
		}

		return utils.Normalize(args, params, quoteValues)
	}
}
