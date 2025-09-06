package main

import (
	"context"
	"errors"
	"log/slog"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to run", "error", err)
	}
}

func run(_ context.Context) error {
	return errors.New("not implemented")
}
