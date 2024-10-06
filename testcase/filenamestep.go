package testcase

import (
	"fmt"
	"os"
	"path/filepath"
)

type filenameStep struct {
	content  string
	filename string
}

func (s filenameStep) Execute(t *testStatus) StepResult {
	target := filepath.Join(getTestDirPath(t.Params), s.filename)
	f, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return StepResult{
			Success: false,
			Error:   fmt.Errorf("failed to open file: %w", err),
		}
	}

	defer f.Close()

	if _, err = f.WriteString(s.content); err != nil {
		return StepResult{
			Success: false,
			Error:   fmt.Errorf("failed to write code block content to file: %w", err),
		}
	}

	return StepResult{
		Success: true,
	}
}

func parseFilenameStep(options map[string]string, content string) (Step, error) {
	return filenameStep{content: content, filename: options["filename"]}, nil
}