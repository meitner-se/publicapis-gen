package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	ctx := context.Background()

	err := run(ctx)
	require.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}
