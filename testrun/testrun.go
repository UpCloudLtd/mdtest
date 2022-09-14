package testrun

import (
	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/progress"
)

type RunResult struct {
	Success      bool
	SuccessCount int
	FailureCount int
	TestResults  []testcase.TestResult
}

func Execute(paths []string) RunResult {
	testLog := progress.NewProgress(nil)
	testLog.Start()

	run := RunResult{Success: true}
	for _, path := range paths {
		res := testcase.Execute(path, testLog)
		if res.Success {
			run.SuccessCount++
		} else {
			run.FailureCount++
		}

		run.TestResults = append(run.TestResults, res)
	}

	testLog.Stop()
	run.Success = run.FailureCount == 0
	return run
}
