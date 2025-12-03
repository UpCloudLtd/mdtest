package utils_test

import (
	"testing"

	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/stretchr/testify/assert"
)

func TestEnvEntriesAsMap(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		input    []string
		expected map[string]string
	}{
		{
			name:  "basic entries",
			input: []string{"KEY1=VALUE1", "KEY2=VALUE2"},
			expected: map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
			},
		},
		{
			name:  "empty value",
			input: []string{"KEY1"},
			expected: map[string]string{
				"KEY1": "",
			},
		},
		{
			name:  "whitespace and comments",
			input: []string{"KEY1=VALUE1", "# This is a comment", "    ", "KEY2=VALUE2"},
			expected: map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
			},
		},
		{
			name:  "duplicate keys",
			input: []string{"KEY1=VALUE1", "KEY1=VALUE2"},
			expected: map[string]string{
				"KEY1": "VALUE2",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := utils.EnvEntriesAsMap(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
