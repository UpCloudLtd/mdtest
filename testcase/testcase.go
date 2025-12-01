package testcase

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

type EnvSource string

const (
	EnvSourceBuiltIn  EnvSource = "builtin"
	EnvSourceCommand  EnvSource = "command"
	EnvSourceEnviron  EnvSource = "environ"
	EnvSourceTestcase EnvSource = "testcase"
)

type testStatus struct {
	Params TestParameters
	Env    map[EnvSource][]string
}

func (t *testStatus) GetEnv() []string {
	env := t.Env[EnvSourceEnviron]
	env = append(env, t.Env[EnvSourceTestcase]...)
	env = append(env, t.Env[EnvSourceCommand]...)
	env = append(env, t.Env[EnvSourceBuiltIn]...)
	return env
}

func (t *testStatus) GetCelEnv() (*cel.Env, error) {
	m := utils.EnvEntriesAsMap(t.GetEnv())
	opts := make([]cel.EnvOption, 0, len(m))
	for k, v := range m {
		opts = append(opts, cel.Constant(k, cel.StringType, types.String(v)))
	}

	env, err := cel.NewEnv(opts...)
	if err != nil {
		err = fmt.Errorf("failed to initialize CEL environment: %w", err)
	}
	return env, err
}

func NewTestStatus(params TestParameters) testStatus {
	status := testStatus{
		Env:    make(map[EnvSource][]string),
		Params: params,
	}

	status.Env[EnvSourceEnviron] = os.Environ()
	status.Env[EnvSourceBuiltIn] = []string{
		fmt.Sprintf("MDTEST_JOBID=%d", params.JobID),
		"MDTEST_RUNID=" + params.RunID,
		"MDTEST_TESTID=" + params.TestID,
		"MDTEST_VERSION=" + globals.Version,
		"MDTEST_WORKSPACE=" + getTestDirPath(params),
	}
	status.Env[EnvSourceCommand] = params.EnvOverride

	return status
}

type TestParameters struct {
	EnvOverride []string
	JobID       int
	RunID       string
	TestID      string
	TestLog     *progress.Progress
}

type TestResult struct {
	Name         string
	Started      time.Time
	Finished     time.Time
	Success      bool
	SuccessCount int
	FailureCount int
	StepsCount   int
	Results      []StepResult
	Error        error
}

func (t TestResult) SkippedCount() int {
	return t.StepsCount - t.SuccessCount - t.FailureCount
}

func parse(path string) (string, []Step, error) {
	name := path
	file, err := os.Open(path)
	if err != nil {
		return name, nil, fmt.Errorf(`failed to open test file at "%s" (%w)`, path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	steps := []Step{}
	for scanner.Scan() {
		line := scanner.Text()
		if name == path && strings.HasPrefix(line, "# ") {
			name = strings.TrimPrefix(line, "# ")
			name = strings.TrimSpace(name)
		}

		if strings.HasPrefix(line, "```") {
			step, err := parseStep(scanner)
			if err != nil {
				return name, nil, err
			}
			if step != nil {
				steps = append(steps, step)
			}
		}
	}

	return name, steps, nil
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

	skippedCount := test.SkippedCount()
	if skippedCount > 0 {
		details += fmt.Sprintf(" (%d skipped)", skippedCount)
	}

	return details
}

func errorMessage(key string, err error) messages.Update {
	return messages.Update{
		Key:     key,
		Status:  messages.MessageStatusError,
		Details: fmt.Sprintf("Error: %s", err.Error()),
	}
}

func Execute(ctx context.Context, path string, params TestParameters) TestResult {
	testLog := params.TestLog
	started := time.Now()

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Parsing %s", path), Status: messages.MessageStatusStarted})

	name, steps, err := parse(path)
	if err != nil {
		_ = testLog.Push(errorMessage(path, err))
		return TestResult{Name: name, Error: err}
	}

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Creating temporary directory for %s", path), Status: messages.MessageStatusStarted})

	err = createTestDir(params)
	if err != nil {
		_ = testLog.Push(errorMessage(path, err))
		return TestResult{Name: name, Error: err}
	}

	_ = testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Running %s", path)})

	test := TestResult{
		Name:       name,
		Started:    started,
		StepsCount: len(steps),
	}
	status := NewTestStatus(params)
	for i, step := range steps {
		_ = testLog.Push(messages.Update{
			Key:             path,
			ProgressMessage: fmt.Sprintf("(Step %d of %d)", i+1, len(steps)),
		})

		if err := ctx.Err(); err != nil {
			test.Error = utils.GetContextError(err)
		}

		if test.FailureCount > 0 && !step.IsCleanup() {
			test.Results = append(test.Results, StepResult{})
			continue
		}

		res := step.Execute(ctx, &status)
		test.Results = append(test.Results, res)
		switch res.Status {
		case StepStatusSuccess:
			test.SuccessCount++
		case StepStatusFailure:
			test.FailureCount++
		}
	}

	test.Finished = time.Now()
	test.Success = test.FailureCount == 0 && test.Error == nil
	if test.Success {
		_ = testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusSuccess})
		_ = removeTestDir(params)
	} else {
		_ = testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusError, Details: getFailureDetails(test)})
	}
	return test
}
