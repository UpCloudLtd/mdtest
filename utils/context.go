package utils

import (
	"context"
	"errors"
	"fmt"
)

func GetContextError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("test run timeout exceeded")
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("test run was canceled with interrupt signal")
	default:
		return fmt.Errorf("unexpected context error (%w)", err)
	}
}
