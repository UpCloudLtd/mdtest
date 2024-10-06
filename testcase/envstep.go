package testcase

import (
	"strings"
)

type envStep struct {
	envUpdates []string
}

func (s envStep) Execute(t *testStatus) StepResult {
	t.Env = append(t.Env, s.envUpdates...)

	return StepResult{
		Success: true,
	}
}

func parseEnvStep(_ map[string]string, content string) (Step, error) {
	return envStep{envUpdates: strings.Split(content, "\n")}, nil
}
