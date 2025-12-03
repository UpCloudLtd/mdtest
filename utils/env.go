package utils

import "strings"

func EnvEntriesAsMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, entry := range env {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else {
			m[parts[0]] = ""
		}
	}
	return m
}
