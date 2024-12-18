package testcase

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

type testStatus struct {
	Params TestParameters
	Env    []string
}

func NewTestStatus(params TestParameters) testStatus {
	return testStatus{
		Env: append(os.Environ(),
			fmt.Sprintf("MDTEST_JOBID=%d", params.JobID),
			"MDTEST_RUNID="+params.RunID,
			"MDTEST_TESTID="+params.TestID,
			"MDTEST_VERSION="+globals.Version,
			"MDTEST_WORKSPACE="+getTestDirPath(params),
		),
		Params: params,
	}
}

type TestParameters struct {
	JobID   int
	RunID   string
	TestID  string
	TestLog *progress.Progress
}

type TestResult struct {
	Success      bool
	SuccessCount int
	FailureCount int
	StepsCount   int
	Results      []StepResult
	Error        error
}

func parse(path string) ([]Step, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(`failed to open test file at "%s" (%w)`, path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	steps := []Step{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			step, err := parseStep(scanner)
			if err != nil {
				return nil, err
			}
			if step != nil {
				steps = append(steps, step)
			}
		}
	}

	return steps, nil
}

func getTestDirPath(params TestParameters) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("mdtest_%s_%s", params.RunID, params.TestID))
}

func createTestDir(params TestParameters) error {
	dirPath := getTestDirPath(params)
	err := os.MkdirAll(dirPath, 0o750)
	if err != nil {
		return fmt.Errorf(`failed to create temporary directory "%s": "%w`, dirPath, err)
	}

	return nil
}

func removeTestDir(params TestParameters) error {
	dirPath := getTestDirPath(params)
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf(`failed to remove directory "%s": "%w`, dirPath, err)
	}

	return nil
}

func getFailureDetails(test TestResult) string {
	details := ""
	if test.FailureCount > 0 {
		details += "Failures:"
	} else if test.Error != nil {
		details += "Canceled: " + test.Error.Error()
	}
	for i, res := range test.Results {
		if err := res.Error; err != nil {
			details += fmt.Sprintf("\n\nStep %d: %s", i+1, err.Error())

			if res.Output != "" {
				details += "\n\nOutput:\n\n"
				details += res.Output
			}
		}
	}

	details += fmt.Sprintf("\n\n%d of %d test steps failed", test.FailureCount, test.StepsCount)

	skippedCount := test.StepsCount - test.SuccessCount - test.FailureCount
	if skippedCount > 0 {
		details += fmt.Sprintf(" (%d skipped)", skippedCount)
	}

	return details
}

func Execute(ctx context.Context, path string, params TestParameters) TestResult {
	testLog := params.TestLog

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Parsing %s", path), Status: messages.MessageStatusStarted})

	steps, err := parse(path)
	if err != nil {
		_ = testLog.Push(messages.Update{
			Key:     path,
			Status:  messages.MessageStatusError,
			Details: fmt.Sprintf("Error: %s", err.Error()),
		})
		return TestResult{Error: err}
	}

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Creating temporary directory for %s", path), Status: messages.MessageStatusStarted})

	err = createTestDir(params)
	if err != nil {
		_ = testLog.Push(messages.Update{
			Key:     path,
			Status:  messages.MessageStatusError,
			Details: fmt.Sprintf("Error: %s", err.Error()),
		})
		return TestResult{Error: err}
	}

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Running %s", path)})

	test := TestResult{StepsCount: len(steps)}
	status := NewTestStatus(params)
	for i, step := range steps {
		_ = testLog.Push(messages.Update{
			Key:             path,
			ProgressMessage: fmt.Sprintf("(Step %d of %d)", i+1, len(steps)),
		})

		if err := ctx.Err(); err == nil {
			res := step.Execute(ctx, &status)
			test.Results = append(test.Results, res)
			if res.Success {
				test.SuccessCount++
			} else {
				test.FailureCount++
			}
		} else {
			test.Error = utils.GetContextError(err)
		}
	}

	test.Success = test.SuccessCount == test.StepsCount
	if test.Success {
		_ = testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusSuccess})
		_ = removeTestDir(params)
	} else {
		_ = testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusError, Details: getFailureDetails(test)})
	}
	return test
}
