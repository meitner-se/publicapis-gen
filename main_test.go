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
	t.Cleanup(func() {
		slog.SetDefault(slog.Default())
	})

	buf := new(bytes.Buffer)

	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))

	main()

	assert.Contains(t, buf.String(), "level=ERROR msg=\"failed to run\" error=\"not implemented\"\n")
}

func Test_run(t *testing.T) {
	ctx := context.Background()

	err := run(ctx)
	require.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}
