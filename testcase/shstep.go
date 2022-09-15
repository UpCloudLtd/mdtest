package testcase

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type shStep struct {
	script   string
	exitCode int
}

func unexpectedExitCode(expected, got int) error {
	return fmt.Errorf("expected exit code %d, got %d", expected, got)
}

func (s shStep) Execute(t *testStatus) StepResult {
	cmd := exec.Command("sh", "-xec", s.script)
	cmd.Env = t.Env
	output, err := cmd.CombinedOutput()
	if err != nil {
		var exit *exec.ExitError
		isExit := errors.As(err, &exit)
		if !isExit {
			return StepResult{
				Success: false,
				Error:   fmt.Errorf("unexpected error (%s)", err.Error()),
				Output:  string(output),
			}
		} else if got := exit.ExitCode(); got != s.exitCode {
			return StepResult{
				Success: false,
				Error:   unexpectedExitCode(s.exitCode, got),
				Output:  string(output),
			}
		}
	} else {
		if s.exitCode != 0 {
			return StepResult{
				Success: false,
				Error:   unexpectedExitCode(s.exitCode, 0),
				Output:  string(output),
			}
		}
	}

	return StepResult{
		Success: true,
		Output:  string(output),
	}
}

func parseShStep(options map[string]string, content string) (Step, error) {
	exitCode, _ := strconv.Atoi(options["exit_code"])

	return shStep{script: content, exitCode: exitCode}, nil
}
