package testcase

import (
	"strings"
)

type envStep struct {
	envUpdates []string
}

func (e envStep) Execute(t *testStatus) StepResult {
	t.Env = append(t.Env, e.envUpdates...)

	return StepResult{
		Success: true,
	}
}

func parseEnvStep(options map[string]string, content string) (Step, error) {
	return envStep{envUpdates: strings.Split(content, "\n")}, nil
}
