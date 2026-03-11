package testcase

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/mdtest/utils"
)

type envStep struct {
	envUpdates []string
}

var _ Step = envStep{}

var (
	whitespaceRe  = regexp.MustCompile(`\s`)
	singlequoteRe = regexp.MustCompile(`^'[^']*'$`)
	doublequoteRe = regexp.MustCompile(`^".*"$`)
)

func parseValue(value string, env EnvBySource, envUpdates []string) (string, error) {
	var err error

	// Do not expand variables in single quotes, but remove the quotes
	if singlequoteRe.MatchString(value) {
		return value[1 : len(value)-1], nil
	}

	// Remove double quotes
	if doublequoteRe.MatchString(value) {
		value = value[1 : len(value)-1]
	} else if whitespaceRe.MatchString(value) {
		// Return an error (with the parsed value) if unquoted value contains whitespace, as it would fail when executed in shell. The error is shown as a warning in the test output.
		err = errors.New("variable values that contain whitespace should be quoted")
	}

	// Expand variables in unquoted and double-quoted values
	return os.Expand(value, func(key string) string {
		env = maps.Clone(env)
		env[EnvSourceTestcase] = append(env[EnvSourceTestcase], envUpdates...)

		m := utils.EnvEntriesAsMap(env.Merge())
		return m[key]
	}), err
}

func (s envStep) Execute(_ context.Context, t *testStatus) StepResult {
	var envUpdates []string
	var output strings.Builder
	var warnings []string

	for _, line := range s.envUpdates {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines without adding them to output
		if trimmed == "" {
			continue
		}

		// Skip comments after adding them to output
		if strings.HasPrefix(trimmed, "#") {
			output.WriteString(fmt.Sprintf("%s\n", trimmed))
			continue
		}

		output.WriteString(fmt.Sprintf("+ %s\n", trimmed))

		parts := strings.SplitN(trimmed, "=", 2)
		if whitespaceRe.MatchString(parts[0]) {
			warnings = append(warnings, fmt.Sprintf(`variable key should not contain whitespace: %s`, trimmed))
		}

		if len(parts) < 2 {
			envUpdates = append(envUpdates, trimmed)
			output.WriteString(fmt.Sprintf("%s=", parts[0]))
			continue
		}

		value, err := parseValue(parts[1], t.Env, envUpdates)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf(`%s: %s`, err.Error(), trimmed))
		}
		envUpdate := fmt.Sprintf("%s=%s", parts[0], value)

		envUpdates = append(envUpdates, envUpdate)
		output.WriteString(fmt.Sprintln(envUpdate))
	}

	t.Env[EnvSourceTestcase] = append(t.Env[EnvSourceTestcase], envUpdates...)

	return StepResult{
		Status:   StepStatusSuccess,
		Output:   output.String(),
		Warnings: warnings,
	}
}

func (s envStep) IsCleanup() bool {
	return false
}

func parseEnvStep(_ utils.Options, content string) (Step, error) {
	return envStep{envUpdates: strings.Split(content, "\n")}, nil
}
