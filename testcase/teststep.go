package testcase

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/mdtest/utils"
)

type StepResult struct {
	Success bool
	Output  string
	Error   error
}

type Step interface {
	Execute(*testStatus) StepResult
}

func parseCodeBlock(lang string, options map[string]string, content string) (Step, error) {
	if options["filename"] != "" {
		return parseFilenameStep(options, content)
	}

	switch lang {
	case "env":
		return parseEnvStep(options, content)
	case "sh":
		return parseShStep(options, content)
	default:
		return nil, nil //nolint:nilnil // "Parsed" non-step code block without errors
	}
}

func parseStep(scanner *bufio.Scanner) (Step, error) {
	line := scanner.Text()
	if !strings.HasPrefix(line, "```") {
		return nil, fmt.Errorf("current scanner position is not at start of a test step")
	}

	lang, options := utils.ParseOptions(line[3:])

	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			return parseCodeBlock(lang, options, content)
		}

		content += line + "\n"
	}
	return nil, fmt.Errorf("could not parse test step (unexpected EOF)")
}
