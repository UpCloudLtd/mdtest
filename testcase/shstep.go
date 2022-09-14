package testcase

import (
	"os/exec"
)

type shStep struct {
	script string
}

func (s shStep) Execute() StepResult {
	output, err := exec.Command("sh", "-c", s.script).CombinedOutput()
	if err != nil {
		return StepResult{
			Success: false,
			Error:   err,
			Output:  string(output),
		}
	}

	return StepResult{
		Success: true,
		Output:  string(output),
	}
}

func parseShStep(options []string, content string) (Step, error) {
	return shStep{script: content}, nil
}
