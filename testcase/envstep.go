package testcase

import (
	"context"
	"strings"

	"github.com/UpCloudLtd/mdtest/utils"
)

type envStep struct {
	envUpdates []string
}

var _ Step = envStep{}

func (s envStep) Execute(_ context.Context, t *testStatus) StepResult {
	t.Env[EnvSourceTestcase] = append(t.Env[EnvSourceTestcase], s.envUpdates...)

	return StepResult{
		Status: StepStatusSuccess,
		Output: strings.Join(s.envUpdates, "\n"),
	}
}

func (s envStep) IsCleanup() bool {
	return false
}

func parseEnvStep(_ utils.Options, content string) (Step, error) {
	return envStep{envUpdates: strings.Split(content, "\n")}, nil
}
