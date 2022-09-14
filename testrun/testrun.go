package testrun

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

type RunResult struct {
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

func Execute(rawPaths []string) RunResult {
	paths, warnings := parseFilePaths(rawPaths, 1)

	testLog := progress.NewProgress(nil)
	testLog.Start()

	for _, warning := range warnings {
		testLog.Push(warning.Message())
	}

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
