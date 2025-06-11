package testrun

import (
	"context"

	"github.com/UpCloudLtd/mdtest/id"
	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/progress"
)

func executeTests(ctx context.Context, paths []string, params RunParameters, testLog *progress.Progress, run *RunResult) {
	if len(paths) == 0 {
		return
	}

	jobQueue := make(chan int, params.NumberOfJobs)
	testQueue := paths
	returnChan := make(chan testcase.TestResult)

	for i := 0; i < params.NumberOfJobs; i++ {
		jobQueue <- i
	}

	for {
		select {
		case curJobID := <-jobQueue:
			if len(testQueue) == 0 {
				break
			}
			curTest := testQueue[0]
			testQueue = testQueue[1:]

			go func(jobID int, test string) {
				defer func() {
					jobQueue <- jobID
				}()
				returnChan <- testcase.Execute(ctx, curTest, testcase.TestParameters{
					EnvOverride: params.Env,
					JobID:       jobID,
					RunID:       run.ID,
					TestID:      id.NewTestID(),
					TestLog:     testLog,
				})
			}(curJobID, curTest)
		case res := <-returnChan:
			if res.Success {
				run.SuccessCount++
			} else {
				run.FailureCount++
			}

			run.TestResults = append(run.TestResults, res)

			if len(run.TestResults) >= len(paths) {
				return
			}
		}
	}
}
