package testcase

import (
	"errors"
	"os/exec"
	"strconv"
)

type shStep struct {
	script   string
	exitCode int
}

func (s shStep) Execute() StepResult {
	output, err := exec.Command("sh", "-c", s.script).CombinedOutput()
	if err != nil {
		var exit *exec.ExitError
		if isExit := errors.As(err, &exit); isExit && exit.ExitCode() != s.exitCode || !isExit {
			return StepResult{
				Success: false,
				Error:   err,
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
