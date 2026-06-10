package testcase

import "io"

// safeWriter returns a non-nil io.Writer. If the provided writer is nil, it returns io.Discard.
func safeWriter(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}
