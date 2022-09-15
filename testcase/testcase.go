package testcase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/mdtest/globals"
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

type testStatus struct {
	TestId string
	Env    []string
}

func NewTestStatus() testStatus {
	id := testId()
	return testStatus{
		Env: append(os.Environ(),
			"MDTEST_TESTID="+id,
			"MDTEST_VERSION="+globals.Version,
		),
		TestId: id,
	}
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
		return nil, err
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

func getFailureDetails(test TestResult) string {
	details := "Failures:"
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

	return details
}

func Execute(path string, testLog *progress.Progress) TestResult {
	testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Parsing %s", path), Status: messages.MessageStatusStarted})

	steps, err := parse(path)
	if err != nil {
		testLog.Push(messages.Update{
			Key:     path,
			Status:  messages.MessageStatusError,
			Details: fmt.Sprintf("Error: %s", err.Error()),
		})
		return TestResult{Error: err}
	}

	testLog.Push(messages.Update{Key: path, Message: fmt.Sprintf("Running %s", path)})

	test := TestResult{StepsCount: len(steps)}
	status := testStatus{Env: os.Environ(), TestId: testId()}
	for i, step := range steps {
		testLog.Push(messages.Update{
			Key:             path,
			ProgressMessage: fmt.Sprintf("(Step %d of %d)", i+1, len(steps)),
		})

		res := step.Execute(&status)
		if res.Success {
			test.SuccessCount++
		} else {
			test.FailureCount++
		}

		test.Results = append(test.Results, res)
	}

	test.Success = test.FailureCount == 0
	if test.Success {
		testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusSuccess})
	} else {
		testLog.Push(messages.Update{Key: path, Status: messages.MessageStatusError, Details: getFailureDetails(test)})
	}
	return test
}
