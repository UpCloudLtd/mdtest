//nolint:paralleltest // These can not be run in parallel due to how args are overridden
package cmd

import (
	"testing"
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
			testPath: "../testdata/success_expected_0_got_0.md",
			exitCode: 0,
		},
		{
			testPath: "../testdata",
			exitCode: 2,
		},
	} {
		test := test
		t.Run(test.testPath, func(t *testing.T) {
			rootCmd.SetArgs([]string{test.testPath})
			exitCode := Execute()
			if exitCode != test.exitCode {
				t.Errorf("Expected exit code %d, got %d", test.exitCode, exitCode)
			}
		})
	}
}
