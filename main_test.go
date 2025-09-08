package main

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	// Arrange
	expectedLogLevel := "level=ERROR"
	expectedLogMessage := `msg="failed to run"`
	expectedErrorMessage := `error="not implemented"`

	t.Cleanup(func() {
		slog.SetDefault(slog.Default())
	})

	buf := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))

	// Act
	main()

	// Assert
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput, "Log output should not be empty")
	assert.Contains(t, logOutput, expectedLogLevel, "Log should contain ERROR level")
	assert.Contains(t, logOutput, expectedLogMessage, "Log should contain 'failed to run' message")
	assert.Contains(t, logOutput, expectedErrorMessage, "Log should contain 'not implemented' error")
}

func Test_run(t *testing.T) {
	// Arrange
	expectedErrorMessage := errorNotImplemented
	ctx := context.Background()

	// Act
	err := run(ctx)

	// Assert
	require.Error(t, err, "run() should return an error")
	assert.Equal(t, expectedErrorMessage, err.Error(), "Error message should match expected 'not implemented' message")
}
