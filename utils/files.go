package utils

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/progress/messages"
)

var quotedValueRegex = regexp.MustCompile(`^["'](.*)["']$`)

type PathWarning struct {
	path string
	err  error
}

func (warn PathWarning) Message() messages.Update {
	return messages.Update{
		Message: fmt.Sprintf("Finding %s", warn.path),
		Details: fmt.Sprintf("Error: %s", warn.err.Error()),
		Status:  messages.MessageStatusWarning,
	}
}

func ParseFilePaths(rawPaths []string, depth int) ([]string, []PathWarning) {
	paths := []string{}
	warnings := []PathWarning{}
	for _, rawPath := range rawPaths {
		info, err := os.Stat(rawPath)
		if err != nil {
			warnings = append(warnings, PathWarning{rawPath, err})
			if info == nil {
				continue
			}
		}

		if info.Mode().IsDir() && depth != 0 {
			files, err := os.ReadDir(rawPath)
			if err != nil {
				warnings = append(warnings, PathWarning{rawPath, err})
			}

			dirRawPaths := []string{}
			for _, file := range files {
				dirRawPaths = append(dirRawPaths, path.Join(rawPath, file.Name()))
			}

			dirPaths, dirWarnings := ParseFilePaths(dirRawPaths, depth-1)
			if dirWarnings != nil {
				warnings = append(warnings, dirWarnings...)
			}

			paths = append(paths, dirPaths...)
		}

		if strings.HasSuffix(rawPath, ".md") {
			paths = append(paths, rawPath)
		}
	}
	return paths, warnings
}

func OptionToBoolean(value string) bool {
	return strings.ToLower(value) == "true"
}

func ParseOptions(optionsStr string) (string, map[string]string) {
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

		if quotedValueRegex.MatchString(value) {
			value = value[1 : len(value)-1]
		}

		options[key] = value
	}

	return lang, options
}
