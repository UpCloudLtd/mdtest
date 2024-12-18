//nolint:paralleltest // These can not be run in parallel due to how args are overridden
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot_testdata(t *testing.T) {
	for _, test := range []struct {
		testPath string
		exitCode int
	}{
		{
			testPath: "../testdata/fail_expected_0_got_3.md",
			exitCode: 1,
		},
		{
			testPath: "../testdata/fail_expected_1_got_0.md",
			exitCode: 1,
		},
		{
			testPath: "../testdata/fail_invalid_test_step.md",
			exitCode: 1,
		},
		{
			testPath: "../testdata/success_expected_0_got_0.md",
			exitCode: 0,
		},
		{
			testPath: "../testdata",
			exitCode: 3,
		},
	} {
		test := test
		t.Run(test.testPath, func(t *testing.T) {
			rootCmd.SetArgs([]string{test.testPath})
			exitCode := Execute()
			assert.Equal(t, test.exitCode, exitCode)
		})
	}
}

func TestNormalise_testdata(t *testing.T) {
	for _, test := range []struct {
		testPath      string
		transformArgs []string
		exitCode      int
		output        string
	}{
		{
			testPath:      "../testdata/success_normalise_infotexts.md",
			transformArgs: []string{"-t", "filename=title"},
			exitCode:      0,
			output: `# Success: normalise info texts

The normalise command with ` + "`" + `-t filename=title` + "`" + ` transform argument should remove and ` + "`" + `no_value` + "`" + ` and ` + "`" + `key=value` + "`" + ` args and replace ` + "`" + `filename` + "`" + ` key with ` + "`" + `title` + "`" + `,

` + "```" + `sh title=true.sh
exit 0
` + "```" + `
`,
		},
		{
			testPath:      "../testdata/success_normalise_infotexts.md",
			transformArgs: []string{"-t", "filename=title", "--quote-values=always"},
			exitCode:      0,
			output: `# Success: normalise info texts

The normalise command with ` + "`" + `-t filename=title` + "`" + ` transform argument should remove and ` + "`" + `no_value` + "`" + ` and ` + "`" + `key=value` + "`" + ` args and replace ` + "`" + `filename` + "`" + ` key with ` + "`" + `title` + "`" + `,

` + "```" + `sh title="true.sh"
exit 0
` + "```" + `
`,
		},
	} {
		test := test
		t.Run(test.testPath, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "example")
			if err != nil {
				t.Errorf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(dir)

			var args []string
			args = append(args, "normalise", "-o", dir)
			args = append(args, test.transformArgs...)
			args = append(args, test.testPath)
			rootCmd.SetArgs(args)

			exitCode := Execute()
			assert.Equal(t, test.exitCode, exitCode)

			outputFile := filepath.Join(dir, filepath.Base(test.testPath))
			assert.FileExists(t, outputFile)

			outputBytes, err := os.ReadFile(filepath.Join(dir, filepath.Base(test.testPath)))
			if err != nil {
				t.Errorf("Failed to read output file: %v", err)
			}
			assert.Equal(t, test.output, string(outputBytes))
		})
	}
}
