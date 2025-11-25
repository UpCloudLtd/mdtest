package utils_test

import (
	"testing"

	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/stretchr/testify/assert"
)

func stringP(s string) *string {
	return &s
}

func TestParseOptions(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name            string
		input           string
		expectedLang    string
		expectedOptions utils.Options
	}{
		{
			name:         "plain value",
			input:        `sh exit_code=313`,
			expectedLang: "sh",
			expectedOptions: utils.Options{
				"exit_code": stringP("313"),
			},
		},
		{
			name:         "quoted value",
			input:        `py filename="example.py"`,
			expectedLang: "py",
			expectedOptions: utils.Options{
				"filename": stringP("example.py"),
			},
		},
		{
			name:         "empty values",
			input:        `txt no_value empty_double_quotes="" empty= empty_single_quotes=''`,
			expectedLang: "txt",
			expectedOptions: utils.Options{
				"empty":               stringP(""),
				"no_value":            nil,
				"empty_single_quotes": stringP(""),
				"empty_double_quotes": stringP(""),
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
