package testcase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/UpCloudLtd/mdtest/utils"
)

type shStep struct {
	script   string
	cleanup  bool
	when     string
	exitCode int
}

var _ Step = shStep{}

func (s shStep) shParams() string {
	if s.cleanup {
		return "-xc"
	}
	return "-xec"
}

func safeWriter(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}

func unexpectedExitCode(expected, got int) error {
	return fmt.Errorf("expected exit code %d, got %d", expected, got)
}

func (s shStep) shouldSkip(t *testStatus) (bool, error) {
	if s.when == "" {
		return false, nil
	}

	env, err := t.GetCelEnv()
	if err != nil {
		return false, err
	}

	val, err := utils.EvaluateExpression(env, s.when)
	if err != nil {
		return false, err //nolint:wrapcheck // utils package wraps the error from cel-go
	}

	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("when expression did not evaluate to a boolean")
	}

	return !b, nil
}

func (s shStep) Execute(ctx context.Context, t *testStatus) StepResult {
	skip, err := s.shouldSkip(t)
	if err != nil {
		return StepResult{
			Status: StepStatusFailure,
			Error:  err,
		}
	}
	if skip {
		return StepResult{
			Status: StepStatusSkipped,
		}
	}

	cmd := exec.CommandContext(ctx, "sh", s.shParams(), s.script) //nolint:gosec // Here we trust that the user knows what their tests do
	cmd.Cancel = func() error {
		return utils.Terminate(cmd)
	}
	cmd.Dir = getTestDirPath(t.Params)
	cmd.Env = t.GetEnv()
	utils.UseProcessGroup(cmd)

	return RunAndHandleStepOutput(ctx, cmd, t, s.exitCode)
}

func RunAndHandleStepOutput(ctx context.Context, cmd *exec.Cmd, t *testStatus, exitCode int) StepResult {
	var output bytes.Buffer
	cmd.Stdout = io.MultiWriter(&output, safeWriter(t.Params.StdoutWriter))
	cmd.Stderr = io.MultiWriter(&output, safeWriter(t.Params.StderrWriter))

	if err := cmd.Start(); err != nil {
		return StepResult{
			Status: StepStatusFailure,
			Error:  fmt.Errorf("failed to start command: %w", err),
		}
	}

	if err := cmd.Wait(); err != nil {
		var exit *exec.ExitError
		isExit := errors.As(err, &exit)
		if !isExit {
			return StepResult{
				Status: StepStatusFailure,
				Error:  fmt.Errorf("unexpected error (%w)", err),
				Output: output.String(),
			}
		}
		got := exit.ExitCode()

		if ctxErr := ctx.Err(); got == -1 && ctxErr != nil {
			return StepResult{
				Status: StepStatusFailure,
				Error:  utils.GetContextError(ctxErr),
				Output: output.String(),
			}
		}
		if got != exitCode {
			return StepResult{
				Status: StepStatusFailure,
				Error:  unexpectedExitCode(exitCode, got),
				Output: output.String(),
			}
		}
	} else if exitCode != 0 {
		return StepResult{
			Status: StepStatusFailure,
			Error:  unexpectedExitCode(exitCode, 0),
			Output: output.String(),
		}
	}

	return StepResult{
		Status: StepStatusSuccess,
		Output: output.String(),
	}
}

func (s shStep) IsCleanup() bool {
	return s.cleanup
}

func parseShStep(options utils.Options, content string) (Step, error) {
	exitCode, _ := strconv.Atoi(options.GetString("exit_code"))
	cleanup := options.GetBoolean("cleanup")
	when := options.GetString("when")

	return shStep{script: content, cleanup: cleanup, exitCode: exitCode, when: when}, nil
}
