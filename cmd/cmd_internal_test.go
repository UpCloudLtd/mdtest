//nolint:paralleltest // These can not be run in parallel due to how args are overridden
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testdataExpectedJUnitXML() string {
	timeoutExitCodeFailure := ""
	xEchoCommandWithQuotes := "+ echo &#39;Environment variable VAR is set to VALUE&#39;"
	if runtime.GOOS == "linux" {
		xEchoCommandWithQuotes = "+ echo Environment variable VAR is set to VALUE"
	}
	if runtime.GOOS == "windows" {
		timeoutExitCodeFailure = "\n    <failure>expected exit code 0, got 3221225786</failure>"
	}
	return fmt.Sprintf(`<testsuite name="Test JUnit XML output" tests="9" failures="5" errors="1" skipped="0" time="ELAPSED" timestamp="STARTED">
  <testcase classname="Test JUnit XML output" name="Fail: expected 0, got 3" time="ELAPSED">
    <failure>expected exit code 0, got 3</failure>
    <system-out># Step 1:&#xA;+ exit 3&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Fail: expected 1, got 0" time="ELAPSED">
    <failure>expected exit code 1, got 0</failure>
    <system-out># Step 1:&#xA;+ exit 0&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Fail: invalid test step">
    <error>could not parse test step (unexpected EOF)</error>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Fail: test environment variable values" time="ELAPSED">
    <failure>expected exit code 0, got 1</failure>
    <system-out># Step 1:&#xA;berry=banana&#xA;fruit=apple&#xA;# Step 2:&#xA;+ test banana = strawberry&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Fail: failing step, skipped step, and (failing) cleanup step" time="ELAPSED">
    <failure>expected exit code 0, got 4</failure>
    <failure>expected exit code 0, got 2</failure>
    <failure>expected exit code 0, got 1</failure>
    <system-out># Step 1:&#xA;+ exit 4&#xA;# Step 2:&#xA;# Skipped&#xA;# Step 3:&#xA;+ exit 2&#xA;# Step 4:&#xA;+ exit 1&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Success: expected 0, got 0" time="ELAPSED">
    <system-out># Step 1:&#xA;+ exit 0&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Success: normalise info texts" time="ELAPSED">
    <system-out># Step 1:&#xA;# No output&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Success: skip sh step with when expression" time="ELAPSED">
    <system-out># Step 1:&#xA;+ exit 0&#xA;# Step 2:&#xA;# Skipped&#xA;# Step 3:&#xA;VAR=VALUE&#xA;# Step 4:&#xA;%s&#xA;Environment variable VAR is set to VALUE&#xA;# Step 5:&#xA;# Skipped&#xA;</system-out>
  </testcase>
  <testcase classname="Test JUnit XML output" name="Sleep" time="ELAPSED">
    <failure>test run timeout exceeded</failure>%s
    <system-out># Step 1:&#xA;+ sleep 600&#xA;# Step 2:&#xA;# Skipped&#xA;</system-out>
  </testcase>
</testsuite>`, xEchoCommandWithQuotes, timeoutExitCodeFailure)
}

func readJUnitXML(t *testing.T, path string) string {
	t.Helper()

	reportBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read JUnit XML file: %v", err)
	}
	report := string(reportBytes)

	timeRe := regexp.MustCompile(`time="[^"]*"`)
	timestampRe := regexp.MustCompile(`timestamp="[^"]*"`)

	report = timeRe.ReplaceAllString(report, `time="ELAPSED"`)
	report = timestampRe.ReplaceAllString(report, `timestamp="STARTED"`)

	return report
}

func TestRoot_testdata(t *testing.T) {
	for _, test := range []struct {
		testPath  string
		exitCode  int
		extraArgs []string
		junitXML  string
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
			testPath: "../testdata/fail_with_cleanup.md",
			exitCode: 1,
		},
		{
			testPath: "../testdata/success_expected_0_got_0.md",
			exitCode: 0,
		},
		{
			testPath:  "../testdata",
			exitCode:  6,
			extraArgs: []string{"--name", "Test JUnit XML output"},
			junitXML:  testdataExpectedJUnitXML(),
		},
		{
			testPath: "../testdata/fail_test_environment_variable_values.md",
			exitCode: 1,
		},
		{
			testPath:  "../testdata/fail_test_environment_variable_values.md",
			extraArgs: []string{"-e", "berry=strawberry"},
			exitCode:  1,
		},
		{
			testPath:  "../testdata/fail_test_environment_variable_values.md",
			extraArgs: []string{"-e", "berry=strawberry", "--env", "fruit=orange"},
			exitCode:  0,
		},
	} {
		t.Run(test.testPath, func(t *testing.T) {
			args := []string{"--timeout", "1s"}

			if len(test.extraArgs) > 0 {
				args = append(args, test.extraArgs...)
			}

			junitPath := ""
			if test.junitXML != "" {
				f, err := os.CreateTemp("", "mdtest-junit.*.xml")
				if err != nil {
					t.Fatal("Cannot create temporary file", err)
				}

				junitPath = f.Name()
				defer os.Remove(junitPath)

				args = append(args, "--jobs", "1", "--junit-xml", junitPath)
			}
			rootCmd.SetArgs(append(args, test.testPath))
			exitCode := Execute()
			assert.Equal(t, test.exitCode, exitCode)

			if test.junitXML != "" {
				actual := readJUnitXML(t, junitPath)
				assert.Equal(t, test.junitXML, actual)
			}
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
