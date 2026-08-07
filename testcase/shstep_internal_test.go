package testcase

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/mdtest/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringP(s string) *string {
	return &s
}

func TestShStep_parseTimeout(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		value   *string
		timeout time.Duration
		errors  bool
	}{
		{
			value:   stringP("1s"),
			timeout: time.Second,
			errors:  false,
		},
		{
			value:  stringP("invalid"),
			errors: true,
		},
		{
			value:   nil,
			timeout: 0,
			errors:  false,
		},
	} {
		name := "nil"
		if test.value != nil {
			name = *test.value
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var options utils.Options
			if test.value == nil {
				_, options = utils.ParseOptions("sh")
			} else {
				_, options = utils.ParseOptions("sh timeout=" + *test.value)
			}

			output, err := parseShStep(options, "echo test")
			if test.errors {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				step, ok := output.(shStep)
				require.True(t, ok)
				assert.Equal(t, test.timeout, step.timeout)
			}
		})
	}
}
