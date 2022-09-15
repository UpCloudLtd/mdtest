package testcase

import (
	"bufio"
	"fmt"
	"strings"
)

type StepResult struct {
	Success bool
	Output  string
	Error   error
}

type Step interface {
	Execute(*testStatus) StepResult
}

func parseOptions(optionsStr string) (string, map[string]string) {
	optionsList := strings.Split(optionsStr, " ")
	options := make(map[string]string)

	lang := optionsList[0]
	for _, option := range optionsList[1:] {
		items := strings.SplitN(option, "=", 2)

		key := items[0]
		value := ""
		if len(items) > 1 {
			value = items[1]
		}

		options[key] = value
	}

	return lang, options
}

func parseStep(scanner *bufio.Scanner) (Step, error) {
	line := scanner.Text()
	if !strings.HasPrefix(line, "```") {
		return nil, fmt.Errorf("current scanner position is not at start of a test step")
	}

	lang, options := parseOptions(line[3:])

	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			break
		}

		content += line + "\n"
	}

	switch lang {
	case "env":
		return parseEnvStep(options, content)
	case "sh":
		return parseShStep(options, content)
	default:
		return nil, nil
	}
}
