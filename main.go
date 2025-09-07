package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/meitner-se/publicapis-gen/constants"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, constants.ErrorFailedToRun, constants.LogKeyError, err)
	}
}

func run(_ context.Context) error {
	return errors.New(constants.ErrorNotImplemented)
}
