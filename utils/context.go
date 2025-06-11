package utils

import (
	"context"
	"errors"
	"fmt"
)

type CanceledError string

var _ error = CanceledError("")

func (e CanceledError) Error() string {
	switch e {
	case CanceledError("timeout"):
		return "test run timeout exceeded"
	case CanceledError("interrupt"):
		return "test run was canceled with interrupt signal"
	default:
		return "test run was canceled"
	}
}

func GetContextError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return CanceledError("timeout")
	case errors.Is(err, context.Canceled):
		return CanceledError("interrupt")
	default:
		return fmt.Errorf("unexpected context error (%w)", err)
	}
}

func IsCanceled(err error) bool {
	switch {
	case errors.Is(err, CanceledError("timeout")):
		return true
	case errors.Is(err, CanceledError("interrupt")):
		return true
	default:
		return false
	}
}
