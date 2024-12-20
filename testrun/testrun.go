package testrun

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/UpCloudLtd/mdtest/id"
	"github.com/UpCloudLtd/mdtest/output"
	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/UpCloudLtd/progress"
)

type RunParameters struct {
	NumberOfJobs int
	OutputTarget io.Writer
	Timeout      time.Duration
}

type RunResult struct {
	ID           string
	Started      time.Time
	Finished     time.Time
	Success      bool
	SuccessCount int
	FailureCount int
	TestResults  []testcase.TestResult
}

func PrintSummary(target io.Writer, run RunResult) {
	tests := output.Total(len(run.TestResults))
	if run.SuccessCount > 0 {
		tests = output.Passed(run.SuccessCount) + ", " + tests
	}
	if run.FailureCount > 0 {
		tests = output.Failed(run.FailureCount) + ", " + tests
	}

	elapsed := fmt.Sprintf("%.3f s", run.Finished.Sub(run.Started).Seconds())

	data := []output.SummaryItem{
		{Key: "Tests", Value: tests},
		{Key: "Elapsed", Value: elapsed},
	}

	fmt.Fprintf(target, "\n%s", output.SummaryTable((data)))
}

func Execute(rawPaths []string, params RunParameters) RunResult {
	ctx, cancel := context.WithCancel(context.Background())
	if params.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), params.Timeout)
		defer cancel()
	}

	started := time.Now()
	paths, warnings := utils.ParseFilePaths(rawPaths, 1)

	testLog := progress.NewProgress(nil)
	testLog.Start()

	// Handle possible interrupts during execution
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		cancel()
	}()

	for _, warning := range warnings {
		_ = testLog.Push(warning.Message())
	}

	run := RunResult{
		ID:      id.NewRunID(),
		Started: started,
		Success: true,
	}
	executeTests(ctx, paths, params, testLog, &run)

	testLog.Stop()
	run.Success = run.FailureCount == 0
	run.Finished = time.Now()

	PrintSummary(params.OutputTarget, run)
	return run
}
