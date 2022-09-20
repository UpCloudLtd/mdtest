package testrun

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/UpCloudLtd/mdtest/id"
	"github.com/UpCloudLtd/mdtest/output"
	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

type RunParameters struct {
	NumberOfJobs int
	OutputTarget io.Writer
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

type PathWarning struct {
	path string
	err  error
}

func (warn PathWarning) Message() messages.Update {
	return messages.Update{
		Message: fmt.Sprintf("Finding %s", warn.path),
		Details: fmt.Sprintf("Error: %s", warn.err.Error()),
		Status:  messages.MessageStatusWarning,
	}
}

func parseFilePaths(rawPaths []string, depth int) ([]string, []PathWarning) {
	paths := []string{}
	warnings := []PathWarning{}
	for _, rawPath := range rawPaths {
		info, err := os.Stat(rawPath)
		if err != nil {
			warnings = append(warnings, PathWarning{rawPath, err})
			if info == nil {
				continue
			}
		}

		if info.Mode().IsDir() && depth != 0 {
			files, err := os.ReadDir(rawPath)
			if err != nil {
				warnings = append(warnings, PathWarning{rawPath, err})
			}

			dirRawPaths := []string{}
			for _, file := range files {
				dirRawPaths = append(dirRawPaths, path.Join(rawPath, file.Name()))
			}

			dirPaths, dirWarnings := parseFilePaths(dirRawPaths, depth-1)
			if dirWarnings != nil {
				warnings = append(warnings, dirWarnings...)
			}

			paths = append(paths, dirPaths...)
		}

		if strings.HasSuffix(rawPath, ".md") {
			paths = append(paths, rawPath)
		}
	}
	return paths, warnings
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
	started := time.Now()
	paths, warnings := parseFilePaths(rawPaths, 1)

	testLog := progress.NewProgress(nil)
	testLog.Start()

	for _, warning := range warnings {
		_ = testLog.Push(warning.Message())
	}

	run := RunResult{
		ID:      id.NewRunID(),
		Started: started,
		Success: true,
	}
	executeTests(paths, params, testLog, &run)

	testLog.Stop()
	run.Success = run.FailureCount == 0
	run.Finished = time.Now()

	PrintSummary(params.OutputTarget, run)
	return run
}
