package common //nolint:revive,nolintlint

import (
	"io"
	"log/slog"
)

// Ptr returns a pointer to the original value.
func Ptr[T any](value T) *T {
	return &value
}

// SafeClose closes an io.Closer and logs any error.
// Use this in a defer statement to satisfy revive linter.
func SafeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("failed to close", "err", err)
	}
}
