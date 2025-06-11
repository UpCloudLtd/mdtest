package testcase

import (
	"context"
	"strings"
)

type envStep struct {
	envUpdates []string
}

func (s envStep) Execute(_ context.Context, t *testStatus) StepResult {
	t.Env[EnvSourceTestcase] = append(t.Env[EnvSourceTestcase], s.envUpdates...)

	return StepResult{
		Success: true,
		Output:  strings.Join(s.envUpdates, "\n"),
	}
}

func parseEnvStep(_ map[string]string, content string) (Step, error) {
	return envStep{envUpdates: strings.Split(content, "\n")}, nil
}
