package utils_test

import (
	"testing"

	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name            string
		input           string
		expectedLang    string
		expectedOptions map[string]string
	}{
		{
			name:         "plain value",
			input:        `sh exit_code=313`,
			expectedLang: "sh",
			expectedOptions: map[string]string{
				"exit_code": "313",
			},
		},
		{
			name:         "quoted value",
			input:        `py filename="example.py"`,
			expectedLang: "py",
			expectedOptions: map[string]string{
				"filename": "example.py",
			},
		},
		{
			name:         "empty values",
			input:        `txt no_value empty_double_quotes="" empty= empty_single_quotes=''`,
			expectedLang: "txt",
			expectedOptions: map[string]string{
				"empty":               "",
				"no_value":            "",
				"empty_single_quotes": "",
				"empty_double_quotes": "",
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			lang, options := utils.ParseOptions(test.input)
			assert.Equal(t, test.expectedLang, lang)
			assert.Equal(t, test.expectedOptions, options)
		})
	}
}
