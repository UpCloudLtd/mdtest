package testrun

import (
	"github.com/UpCloudLtd/mdtest/id"
	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/progress"
)

func executeTests(paths []string, params RunParameters, testLog *progress.Progress, run *RunResult) {
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
		case curJobId := <-jobQueue:
			if len(testQueue) == 0 {
				break
			}
			curTest := testQueue[0]
			testQueue = testQueue[1:]

			go func(jobId int, test string) {
				defer func() {
					jobQueue <- jobId
				}()
				returnChan <- testcase.Execute(curTest, testcase.TestParameters{
					JobId:   jobId,
					RunId:   run.Id,
					TestId:  id.NewTestId(),
					TestLog: testLog,
				})
			}(curJobId, curTest)
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
