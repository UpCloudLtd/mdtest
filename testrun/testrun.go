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
	"github.com/UpCloudLtd/progress/messages"
)

type RunParameters struct {
	Env              []string
	JUnitXML         string
	Name             string
	NumberOfJobs     int
	OutputTarget     io.Writer
	Timeout          time.Duration
	WarningsAsErrors bool
}

type RunResult struct {
	ID           string
	Name         string
	Started      time.Time
	Finished     time.Time
	Success      bool
	SuccessCount int
	FailureCount int
	TestResults  []testcase.TestResult
}

func PrintName(target io.Writer, name string) {
	data := []output.SummaryItem{
		{Key: "Name", Value: name},
	}

	fmt.Fprintf(target, "%s\n", output.SummaryTable((data)))
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

func newRunResult(params RunParameters) RunResult {
	started := time.Now()
	runID := id.NewRunID()
	name := params.Name
	if name == "" {
		name = runID
	}

	return RunResult{
		ID:      runID,
		Name:    name,
		Started: started,
		Success: true,
	}
}

func Execute(rawPaths []string, params RunParameters) RunResult {
	ctx, cancel := context.WithCancel(context.Background())
	if params.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), params.Timeout)
		defer cancel()
	}

	if params.Name != "" {
		PrintName(params.OutputTarget, params.Name)
	}

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

	run := newRunResult(params)

	executeTests(ctx, paths, params, testLog, &run)

	run.Success = run.FailureCount == 0
	run.Finished = time.Now()

	if params.JUnitXML != "" {
		warning := func(err error) messages.Update {
			return messages.Update{
				Message: "Generating JUnit XML report",
				Details: fmt.Sprintf("Error: %s", err.Error()),
				Status:  messages.MessageStatusWarning,
			}
		}
		err := junitReport(run, params)
		if err != nil {
			_ = testLog.Push(warning(err))
		}
	}
	testLog.Stop()

	PrintSummary(params.OutputTarget, run)
	return run
}
