package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	// Save original os.Args and flag state
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Reset flag package for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Arrange
	os.Args = []string{"publicapis-gen", "-help"}

	t.Cleanup(func() {
		slog.SetDefault(slog.Default())
	})

	buf := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))

	// Act
	main()

	// Assert - When help is shown, the program should exit cleanly without errors
	logOutput := buf.String()
	assert.Empty(t, logOutput, "No error logs should be generated when showing help")
}

func Test_run(t *testing.T) {
	// Save original os.Args and flag state
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	t.Run("help flag shows usage", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "-help"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		assert.Nil(t, err, "run() should not return an error when showing help")
	})

	t.Run("missing file flag returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when file flag is missing")
		assert.Contains(t, err.Error(), errorInvalidFile, "Error should mention invalid file")
	})

	t.Run("missing mode flag returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "-file=test.yaml"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when mode flag is missing")
		assert.Contains(t, err.Error(), errorInvalidMode, "Error should mention invalid mode")
	})

	t.Run("invalid mode returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "-file=test.yaml", "-mode=invalid"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for invalid mode")
		assert.Contains(t, err.Error(), errorInvalidMode, "Error should mention invalid mode")
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "-file=nonexistent.yaml", "-mode=overlay"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for nonexistent file")
		assert.Contains(t, err.Error(), errorInvalidFile, "Error should mention invalid file")
	})
}

func Test_readSpecificationFile(t *testing.T) {
	t.Run("reads valid YAML file", func(t *testing.T) {
		// Create a temporary YAML file
		testService := specification.Service{
			Name:      "TestService",
			Version:   "1.0.0",
			Resources: []specification.Resource{},
			Objects:   []specification.Object{},
			Enums:     []specification.Enum{},
		}

		yamlData, err := yaml.Marshal(&testService)
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp("", "test-spec-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(yamlData)
		require.NoError(t, err)
		tmpFile.Close()

		// Test reading the file
		service, err := readSpecificationFile(tmpFile.Name())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "TestService", service.Name)
		assert.Equal(t, "1.0.0", service.Version)
	})

	t.Run("reads valid JSON file", func(t *testing.T) {
		// Create a temporary JSON file
		testService := specification.Service{
			Name:      "TestService",
			Version:   "1.0.0",
			Resources: []specification.Resource{},
			Objects:   []specification.Object{},
			Enums:     []specification.Enum{},
		}

		jsonData, err := json.MarshalIndent(&testService, "", "  ")
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp("", "test-spec-*.json")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(jsonData)
		require.NoError(t, err)
		tmpFile.Close()

		// Test reading the file
		service, err := readSpecificationFile(tmpFile.Name())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "TestService", service.Name)
		assert.Equal(t, "1.0.0", service.Version)
	})

	t.Run("returns error for nonexistent file", func(t *testing.T) {
		// Act
		service, err := readSpecificationFile("nonexistent.yaml")

		// Assert
		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), errorInvalidFile)
	})

	t.Run("returns error for unsupported file extension", func(t *testing.T) {
		// Create a temporary file with unsupported extension
		tmpFile, err := os.CreateTemp("", "test-spec-*.txt")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		// Act
		service, err := readSpecificationFile(tmpFile.Name())

		// Assert
		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), errorUnsupportedFormat)
	})
}

func Test_generateOutputPath(t *testing.T) {
	testCases := []struct {
		name      string
		inputFile string
		suffix    string
		expected  string
	}{
		{
			name:      "YAML file with overlay suffix",
			inputFile: "spec.yaml",
			suffix:    suffixOverlay,
			expected:  "spec-overlay.yaml",
		},
		{
			name:      "JSON file with OpenAPI suffix",
			inputFile: "api.json",
			suffix:    suffixOpenAPI,
			expected:  "api-openapi.json",
		},
		{
			name:      "File with path and overlay suffix",
			inputFile: "/path/to/spec.yml",
			suffix:    suffixOverlay,
			expected:  "/path/to/spec-overlay.yml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := generateOutputPath(tc.inputFile, tc.suffix)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}
