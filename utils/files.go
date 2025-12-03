package utils

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/progress/messages"
)

var (
	optionRegex      = regexp.MustCompile(`([^=\s]+(=(('[^']*')|("[^"]*"))|\S+){0,1})`)
	quotedValueRegex = regexp.MustCompile(`^["'](.*)["']$`)
)

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

type Options map[string]*string

func (o Options) GetString(key string) string {
	val := o[key]
	if val == nil {
		return ""
	}
	return *val
}

func (o Options) GetBoolean(key string) bool {
	val, ok := o[key]

	// Option not set
	if !ok {
		return false
	}

	// Option set but no value, e.g. "cleanup"
	if val == nil {
		return true
	}

	// Option set with a value, e.g. "cleanup=true"
	return strings.EqualFold(*val, "true")
}

func splitOptions(optionsStr string) []string {
	return optionRegex.FindAllString(optionsStr, -1)
}

func ParseOptions(optionsStr string) (string, Options) {
	optionsList := splitOptions(optionsStr)
	options := make(Options)

	lang := optionsList[0]
	for _, option := range optionsList[1:] {
		items := strings.SplitN(option, "=", 2)

		key := items[0]
		var value *string = nil
		if len(items) > 1 {
			value = &items[1]
		}

		if value != nil && quotedValueRegex.MatchString(*value) {
			str := (*value)[1 : len(*value)-1]
			value = &str
		}

		options[key] = value
	}

	return lang, options
}
