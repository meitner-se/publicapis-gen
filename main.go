package main

import (
	"context"
	"errors"
	"log/slog"
)

// Error messages and log keys
const (
	errorNotImplemented = "not implemented"
	errorFailedToRun    = "failed to run"
	logKeyError         = "error"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, errorFailedToRun, logKeyError, err)
	}
}

func run(_ context.Context) error {
	return errors.New(errorNotImplemented)
}
