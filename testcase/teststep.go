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
	Execute() StepResult
}

func parseStep(scanner *bufio.Scanner) (Step, error) {
	line := scanner.Text()
	if !strings.HasPrefix(line, "```") {
		return nil, fmt.Errorf("current scanner position is not at start of a test step")
	}

	options := strings.Split(line[3:], " ")

	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			break
		}

		content += line + "\n"
	}

	switch options[0] {
	case "sh":
		return parseShStep(options, content)
	default:
		return nil, nil
	}

}
