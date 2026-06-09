package testcase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvStep_parseValue(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		value  string
		output string
		errors []string
	}{
		{
			name:   "Inner single quotes in unquoted value",
			value:  `asd'123'asd`,
			output: `asd'123'asd`,
			errors: []string{"variable values with inner quotes should be quoted"},
		},
		{
			name:   "Inner single quotes in single quoted value",
			value:  `'asd'123'asd'`,
			output: `asd'123'asd`,
			errors: []string{"variable values with inner quotes should be quoted"},
		},
		{
			name:   "Inner double quotes in unquoted value",
			value:  `{"key": "value"}`,
			output: `{"key": "value"}`,
			errors: []string{
				"variable values with inner quotes should be quoted",
				"variable values that contain whitespace should be quoted",
			},
		},
		{
			name:   "Inner double quotes in double quoted value",
			value:  `"{"key": "value"}"`,
			output: `{"key": "value"}`,
			errors: []string{"variable values with inner quotes should be quoted"},
		},
		{
			name:   "Inner double quotes in single quoted value",
			value:  `'{"key": "value"}'`,
			output: `{"key": "value"}`,
		},
		{
			name:   "Hashtag in unquoted value",
			value:  `Enjoying life #blessed`,
			output: `Enjoying life`,
			errors: []string{"variable values that contain whitespace should be quoted"},
		},
		{
			name:   "Hashtag in single quoted value",
			value:  `'Enjoying life #blessed'`,
			output: `Enjoying life #blessed`,
		},
		{
			name:   "Hashtag in double quoted value",
			value:  `"Enjoying life #blessed"`,
			output: `Enjoying life #blessed`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output, errs := parseValue(test.value, nil, nil)
			assert.Equal(t, test.output, output)
			require.Len(t, errs, len(test.errors))

			if len(test.errors) > 0 {
				for i, err := range errs {
					assert.Equal(t, test.errors[i], err.Error())
				}
			}
		})
	}
}
