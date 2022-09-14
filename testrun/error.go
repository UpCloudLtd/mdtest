package testrun

import "fmt"

type RunError struct {
	FailureCount int
}

func NewRunError(run RunResult) *RunError {
	if run.Success {
		return nil
	}
	return &RunError{FailureCount: run.FailureCount}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (err RunError) ExitCode() int {
	return min(err.FailureCount, 99)
}

func (err RunError) Error() string {
	return fmt.Sprintf("Test run failed with %d failing test case(s)", err.FailureCount)
}
