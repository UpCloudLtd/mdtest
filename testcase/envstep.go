package testcase

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	singlequoteRe = regexp.MustCompile(`^'.*'$`)
	doublequoteRe = regexp.MustCompile(`^".*"$`)
	innerQuotesRe = regexp.MustCompile(`(^"(.*".*)+"$)|(^'(.*'.*)'$)|(^[^'"](.*['"].*)+[^'"]$)`)
	commentRe     = regexp.MustCompile(`\s+#.*`)
)

func parseValue(value string, env EnvBySource, envUpdates []string) (string, []error) {
	var errs []error

	// Check if value contains inner quotes as these are maintained by default (unlike in shell).
	if innerQuotesRe.MatchString(value) {
		errs = append(errs, errors.New("variable values with inner quotes should be quoted"))
	}

	// Do not expand variables in single quotes, but remove the quotes
	if singlequoteRe.MatchString(value) {
		return value[1 : len(value)-1], errs
	}

	// Remove double quotes
	if doublequoteRe.MatchString(value) {
		value = value[1 : len(value)-1]
	} else {
		// Remove comments from unquoted values
		value = commentRe.ReplaceAllString(value, "")

		if whitespaceRe.MatchString(value) {
			// Return an error (with the parsed value) if unquoted value contains whitespace, as it would fail when executed in shell. The error is shown as a warning in the test output.
			errs = append(errs, errors.New("variable values that contain whitespace should be quoted"))
		}
	}

	// Expand variables in unquoted and double-quoted values
	return os.Expand(value, func(key string) string {
		env = maps.Clone(env)
		env[EnvSourceTestcase] = append(env[EnvSourceTestcase], envUpdates...)

		m := utils.EnvEntriesAsMap(env.Merge())
		return m[key]
	}), errs
}

func (s envStep) Execute(_ context.Context, t *testStatus) StepResult {
	var envUpdates []string
	var output strings.Builder
	var warnings []string

	var writer io.Writer = &output
	if t.Params.OutputToTerminal {
		writer = io.MultiWriter(&output, safeWriter(t.Params.StderrWriter))
	}

	for _, line := range s.envUpdates {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines without adding them to output
		if trimmed == "" {
			continue
		}

		// Skip comments after adding them to output
		if strings.HasPrefix(trimmed, "#") {
			fmt.Fprintf(writer, "%s\n", trimmed)
			continue
		}

		fmt.Fprintf(writer, "+ %s\n", trimmed)

		parts := strings.SplitN(trimmed, "=", 2)
		if whitespaceRe.MatchString(parts[0]) {
			warnings = append(warnings, fmt.Sprintf(`variable key should not contain whitespace: %s`, trimmed))
		}

		if len(parts) < 2 {
			envUpdates = append(envUpdates, trimmed)
			fmt.Fprintf(writer, "%s=", parts[0])
			continue
		}

		value, errs := parseValue(parts[1], t.Env, envUpdates)
		if len(errs) > 0 {
			for _, err := range errs {
				warnings = append(warnings, fmt.Sprintf(`%s: %s`, err.Error(), trimmed))
			}
		}
		envUpdate := fmt.Sprintf("%s=%s", parts[0], value)

		envUpdates = append(envUpdates, envUpdate)
		fmt.Fprintf(writer, "%s\n", envUpdate)
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
