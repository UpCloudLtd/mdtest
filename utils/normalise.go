package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
)

type NormalizeParameters struct {
	OutputPath string
	Transforms []string
}

func Normalize(rawPaths []string, params NormalizeParameters) error {
	paths, warnings := ParseFilePaths(rawPaths, 1)

	info, err := os.Stat(params.OutputPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.MkdirAll(params.OutputPath, 0o755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}
		return fmt.Errorf(`failed to stat output path "%s" (%w)`, params.OutputPath, err)
	}
	if info != nil && !info.IsDir() {
		return fmt.Errorf(`output directory "%s" must be a directory`, params.OutputPath)
	}

	normLog := progress.NewProgress(nil)
	normLog.Start()

	for _, warning := range warnings {
		_ = normLog.Push(warning.Message())
	}

	transformMap := make(map[string]string)
	for _, s := range params.Transforms {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) > 1 {
			transformMap[parts[0]] = parts[1]
		}
	}

	for _, path := range paths {
		_ = normLog.Push(messages.Update{
			Key:     path,
			Message: fmt.Sprintf("Normalising %s", path),
			Status:  messages.MessageStatusStarted,
		})
		err := normalize(path, params.OutputPath, transformMap)
		if err != nil {
			_ = normLog.Push(messages.Update{
				Key:     path,
				Details: fmt.Sprintf("Error: %s", err.Error()),
				Status:  messages.MessageStatusError,
			})
		} else {
			_ = normLog.Push(messages.Update{
				Key:    path,
				Status: messages.MessageStatusSuccess,
			})
		}
	}

	normLog.Stop()
	return nil
}

func transformOptions(options, transforms map[string]string) string {
	output := ""
	for key, value := range options {
		if newKey := transforms[key]; newKey != "" {
			if value != "" {
				output += fmt.Sprintf(" %s=%s", newKey, value)
			} else {
				output += fmt.Sprintf(" %s", newKey)
			}
		}
	}

	return output
}

func normalize(path, outputDir string, transforms map[string]string) error {
	input, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`failed to open input file at "%s" (%w)`, path, err)
	}
	defer input.Close()

	output, err := os.Create(filepath.Join(outputDir, filepath.Base(path)))
	if err != nil {
		return fmt.Errorf(`failed to open output file at "%s" (%w)`, path, err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			info := line[3:]
			if len(info) > 0 {
				lang, options := ParseOptions(info)
				line = fmt.Sprintf("```%s%s", lang, transformOptions(options, transforms))
			}
		}
		_, err = output.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf(`failed to write to output file "%s" (%w)`, output.Name(), err)
		}
	}

	return nil
}
