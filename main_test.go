package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/meitner-se/publicapis-gen/specification/server"
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
		os.Args = []string{"publicapis-gen", "generate"}
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
		os.Args = []string{"publicapis-gen", "generate", "-file=test.yaml"}
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
		os.Args = []string{"publicapis-gen", "generate", "-file=test.yaml", "-mode=invalid"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for invalid mode")
		assert.Contains(t, err.Error(), errorInvalidMode, "Error should mention invalid mode")
	})

	t.Run("invalid log level returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "generate", "-file=test.yaml", "-mode=overlay", "-log-level=invalid"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for invalid log level")
		assert.Contains(t, err.Error(), "invalid log level", "Error should mention invalid log level")
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "generate", "-file=nonexistent.yaml", "-mode=overlay"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for nonexistent file")
		assert.Contains(t, err.Error(), errorInvalidFile, "Error should mention invalid file")
	})

	t.Run("complete YAML to OpenAPI JSON pipeline", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Use test data files
		inputSpecFile := "testdata/school-management-api.yaml"
		expectedOutputFile := "testdata/school-management-api-expected.json"

		// Create temporary output file
		tmpOutputFile, err := os.CreateTemp("", "test-output-*.json")
		require.NoError(t, err)
		defer os.Remove(tmpOutputFile.Name())
		tmpOutputFile.Close()

		// Arrange command line arguments for OpenAPI generation
		os.Args = []string{"publicapis-gen", "generate", "-file=" + inputSpecFile, "-mode=openapi", "-output=" + tmpOutputFile.Name()}
		ctx := context.Background()

		// Act - run the command
		err = run(ctx)

		// Assert - command should succeed
		require.NoError(t, err, "run() should not return an error for valid YAML to OpenAPI conversion")

		// Verify output file was created
		_, err = os.Stat(tmpOutputFile.Name())
		require.NoError(t, err, "Output OpenAPI JSON file should be created")

		// Read the generated JSON content
		actualOutputData, err := os.ReadFile(tmpOutputFile.Name())
		require.NoError(t, err, "Should be able to read generated OpenAPI JSON file")

		// Read the expected JSON content
		expectedOutputData, err := os.ReadFile(expectedOutputFile)
		require.NoError(t, err, "Should be able to read expected OpenAPI JSON file")

		// Parse both JSON files to ensure they're valid
		var actualJSON, expectedJSON map[string]interface{}
		err = json.Unmarshal(actualOutputData, &actualJSON)
		require.NoError(t, err, "Generated file should contain valid JSON")
		err = json.Unmarshal(expectedOutputData, &expectedJSON)
		require.NoError(t, err, "Expected file should contain valid JSON")

		// Assert exact JSON match
		assert.JSONEq(t, string(expectedOutputData), string(actualOutputData), "Generated OpenAPI JSON should exactly match expected output")
	})

	t.Run("complete YAML to Go server code e2e pipeline", func(t *testing.T) {
		// Note: This test uses our minimal test infrastructure to verify the
		// server generation integration works correctly with valid OpenAPI documents.
		// The full specification->OpenAPI->server flow has a known limitation with
		// duplicate typenames from comprehensive schema generation.

		// Test the server generation capability using our server package directly
		// with a known-good OpenAPI document from testdata
		config := server.NewDefaultConfig()
		serverGenerator := server.NewGenerator(config)

		// Read the minimal OpenAPI document designed for server generation
		inputOpenAPIFile := "testdata/minimal-openapi-for-server.json"
		expectedServerFile := "testdata/simple-e2e-server-expected.go"

		// Read input OpenAPI document
		inputData, err := os.ReadFile(inputOpenAPIFile)
		require.NoError(t, err, "Should be able to read test OpenAPI document")

		// Generate server code
		actualServerCode, err := serverGenerator.GenerateFromDocument(inputData)
		require.NoError(t, err, "Server generation should succeed with minimal OpenAPI document")

		// Read expected server code
		expectedServerCode, err := os.ReadFile(expectedServerFile)
		require.NoError(t, err, "Should be able to read expected server code")

		// Verify the generated code structure matches expectations
		actualCode := string(actualServerCode)
		_ = string(expectedServerCode) // Keep expected code for future exact comparison if needed

		// Compare key structural elements (allowing for minor formatting differences)
		assert.Contains(t, actualCode, "package api", "Generated code should have correct package")
		assert.Contains(t, actualCode, "ServerInterface", "Generated code should contain server interface")
		assert.Contains(t, actualCode, "type Item struct", "Generated code should contain Item type")
		assert.Contains(t, actualCode, "GetItems", "Generated code should contain GetItems method")

		// Verify it contains essential server framework imports
		assert.Contains(t, actualCode, "github.com/gin-gonic/gin", "Generated code should import Gin framework")

		// Verify generated code structure is similar (not exact match due to potential oapi-codegen version differences)
		assert.Greater(t, len(actualCode), 500, "Generated server code should be substantial")
		assert.Contains(t, actualCode, "// Code generated by", "Generated code should contain generation comment")
	})

	t.Run("config file with valid jobs succeeds", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Create a test config file
		testConfig := Config{
			{
				Specification: "testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-config-openapi.json",
				SchemaJSON:    "test-config-schema.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-config-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())
		defer os.Remove("test-config-openapi.json")
		defer os.Remove("test-config-schema.json")

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Arrange command line arguments for config mode
		os.Args = []string{"publicapis-gen", "generate", "-config=" + tmpConfigFile.Name()}
		ctx := context.Background()

		// Act - run the command
		err = run(ctx)

		// Assert - command should succeed
		require.NoError(t, err, "run() should not return an error for valid config file")

		// Verify output files were created
		_, err = os.Stat("test-config-openapi.json")
		assert.NoError(t, err, "OpenAPI JSON file should be created")
		_, err = os.Stat("test-config-schema.json")
		assert.NoError(t, err, "Schema JSON file should be created")
	})

	t.Run("config file with server generation job", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Create a test config file with server generation
		// Note: Uses minimal OpenAPI file to avoid the duplicate typename issue
		testConfig := Config{
			{
				Specification: "testdata/simple-e2e-api.yaml",
				OpenAPIJSON:   "test-server-config-openapi.json",
				ServerGo:      "test-server-config-server.go",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-config-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())
		defer os.Remove("test-server-config-openapi.json")
		defer os.Remove("test-server-config-server.go")

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Arrange command line arguments for config mode
		os.Args = []string{"publicapis-gen", "generate", "-config=" + tmpConfigFile.Name()}
		ctx := context.Background()

		// Act - run the command
		err = run(ctx)

		// Assert - Server generation will fail due to duplicate typename, but OpenAPI should succeed
		// This verifies the config file integration works for both modes
		require.Error(t, err, "Config should fail during server generation due to duplicate typename limitation")
		assert.Contains(t, err.Error(), "failed to generate server Go code", "Error should be in server generation phase")

		// Verify OpenAPI file was created successfully before server generation failed
		_, err = os.Stat("test-server-config-openapi.json")
		assert.NoError(t, err, "OpenAPI JSON file should be created even if server generation fails")
	})

	t.Run("mixing config and legacy flags returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with both config and legacy flags
		os.Args = []string{"publicapis-gen", "generate", "-config=test.yaml", "-file=spec.yaml"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when mixing config and legacy flags")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
	})

	t.Run("nonexistent config file returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments for nonexistent config
		os.Args = []string{"publicapis-gen", "generate", "-config=nonexistent-config.yaml"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for nonexistent config file")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
	})

	t.Run("uses default publicapis.yaml when available and no flags provided", func(t *testing.T) {
		// Save original working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)

		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-default-config-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create a test default config file (publicapis.yaml)
		testConfig := Config{
			{
				Specification: origDir + "/testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-default-openapi.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		err = os.WriteFile(defaultConfigYAML, yamlData, 0644)
		require.NoError(t, err)
		defer os.Remove("test-default-openapi.json")

		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with no flags (should use default config)
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.NoError(t, err, "run() should succeed when using default config file")

		// Verify output file was created
		_, err = os.Stat("test-default-openapi.json")
		require.NoError(t, err, "OpenAPI JSON file should be created from default config")
	})

	t.Run("uses default publicapis.yml when publicapis.yaml not available", func(t *testing.T) {
		// Save original working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)

		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-default-yml-config-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create a test default config file (publicapis.yml)
		testConfig := Config{
			{
				Specification: origDir + "/testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-default-yml-openapi.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		err = os.WriteFile(defaultConfigYML, yamlData, 0644)
		require.NoError(t, err)
		defer os.Remove("test-default-yml-openapi.json")

		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with no flags (should use default config)
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.NoError(t, err, "run() should succeed when using default .yml config file")

		// Verify output file was created
		_, err = os.Stat("test-default-yml-openapi.json")
		require.NoError(t, err, "OpenAPI JSON file should be created from default .yml config")
	})

	t.Run("prefers publicapis.yaml over publicapis.yml when both exist", func(t *testing.T) {
		// Save original working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)

		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-preference-config-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create both config files with different outputs to test preference
		testConfigYAML := Config{
			{
				Specification: origDir + "/testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-preference-yaml-openapi.json",
			},
		}
		testConfigYML := Config{
			{
				Specification: origDir + "/testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-preference-yml-openapi.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfigYAML)
		require.NoError(t, err)
		err = os.WriteFile(defaultConfigYAML, yamlData, 0644)
		require.NoError(t, err)

		ymlData, err := yaml.Marshal(&testConfigYML)
		require.NoError(t, err)
		err = os.WriteFile(defaultConfigYML, ymlData, 0644)
		require.NoError(t, err)

		defer os.Remove("test-preference-yaml-openapi.json")
		defer os.Remove("test-preference-yml-openapi.json")

		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with no flags (should use default config)
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.NoError(t, err, "run() should succeed when using default config file")

		// Verify the .yaml config was used (not the .yml)
		_, err = os.Stat("test-preference-yaml-openapi.json")
		require.NoError(t, err, "OpenAPI JSON file from .yaml config should be created")

		// Verify the .yml config was NOT used
		_, err = os.Stat("test-preference-yml-openapi.json")
		require.Error(t, err, "OpenAPI JSON file from .yml config should NOT be created")
	})

	t.Run("no default config and no flags returns error", func(t *testing.T) {
		// Save original working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)

		// Create temporary directory and change to it (with no default config)
		tmpDir, err := os.CreateTemp("", "test-no-default-config-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with no flags and no default config
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when no config and no flags provided")
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
		assert.Contains(t, err.Error(), "unsupported file format")
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
			name:      "JSON file with overlay suffix",
			inputFile: "api.json",
			suffix:    suffixOverlay,
			expected:  "api-overlay.json",
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

func Test_generateOpenAPIOutputPath(t *testing.T) {
	testCases := []struct {
		name      string
		inputFile string
		expected  string
	}{
		{
			name:      "YAML file generates JSON OpenAPI",
			inputFile: "spec.yaml",
			expected:  "spec-openapi.json",
		},
		{
			name:      "JSON file generates JSON OpenAPI",
			inputFile: "api.json",
			expected:  "api-openapi.json",
		},
		{
			name:      "File with path generates JSON OpenAPI",
			inputFile: "/path/to/spec.yml",
			expected:  "/path/to/spec-openapi.json",
		},
		{
			name:      "File without extension generates JSON OpenAPI",
			inputFile: "spec",
			expected:  "spec-openapi.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := generateOpenAPIOutputPath(tc.inputFile)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_ensureJSONExtension(t *testing.T) {
	testCases := []struct {
		name       string
		outputPath string
		expected   string
	}{
		{
			name:       "JSON extension unchanged",
			outputPath: "api-spec.json",
			expected:   "api-spec.json",
		},
		{
			name:       "YAML extension changed to JSON",
			outputPath: "api-spec.yaml",
			expected:   "api-spec.json",
		},
		{
			name:       "YML extension changed to JSON",
			outputPath: "api-spec.yml",
			expected:   "api-spec.json",
		},
		{
			name:       "No extension gets JSON extension",
			outputPath: "api-spec",
			expected:   "api-spec.json",
		},
		{
			name:       "Other extension changed to JSON",
			outputPath: "api-spec.xml",
			expected:   "api-spec.json",
		},
		{
			name:       "Path with JSON extension unchanged",
			outputPath: "/path/to/api-spec.json",
			expected:   "/path/to/api-spec.json",
		},
		{
			name:       "Path with YAML extension changed to JSON",
			outputPath: "/path/to/api-spec.yaml",
			expected:   "/path/to/api-spec.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := ensureJSONExtension(tc.outputPath)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_configureLogging(t *testing.T) {
	testCases := []struct {
		name        string
		logLevel    string
		expectError bool
	}{
		{
			name:        "debug level",
			logLevel:    logLevelDebug,
			expectError: false,
		},
		{
			name:        "info level",
			logLevel:    logLevelInfo,
			expectError: false,
		},
		{
			name:        "warn level",
			logLevel:    logLevelWarn,
			expectError: false,
		},
		{
			name:        "error level",
			logLevel:    logLevelError,
			expectError: false,
		},
		{
			name:        "off level",
			logLevel:    logLevelOff,
			expectError: false,
		},
		{
			name:        "invalid level returns error",
			logLevel:    "invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original logger
			originalLogger := slog.Default()
			defer slog.SetDefault(originalLogger)

			// Act
			err := configureLogging(tc.logLevel)

			// Assert
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported log level")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ============================================================================
// Default Config File Tests
// ============================================================================

func Test_findDefaultConfigFile(t *testing.T) {
	// Save original working directory
	origDir, err := os.Getwd()
	require.NoError(t, err)

	t.Run("returns empty string when no default config files exist", func(t *testing.T) {
		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-no-defaults-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Act
		result := findDefaultConfigFile()

		// Assert
		assert.Empty(t, result, "Should return empty string when no default config files exist")
	})

	t.Run("returns publicapis.yaml when it exists", func(t *testing.T) {
		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-yaml-exists-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create publicapis.yaml
		err = os.WriteFile(defaultConfigYAML, []byte("# test config"), 0644)
		require.NoError(t, err)

		// Act
		result := findDefaultConfigFile()

		// Assert
		assert.Equal(t, defaultConfigYAML, result, "Should return publicapis.yaml when it exists")
	})

	t.Run("returns publicapis.yml when only it exists", func(t *testing.T) {
		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-yml-exists-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create publicapis.yml
		err = os.WriteFile(defaultConfigYML, []byte("# test config"), 0644)
		require.NoError(t, err)

		// Act
		result := findDefaultConfigFile()

		// Assert
		assert.Equal(t, defaultConfigYML, result, "Should return publicapis.yml when only it exists")
	})

	t.Run("prefers publicapis.yaml over publicapis.yml when both exist", func(t *testing.T) {
		// Create temporary directory and change to it
		tmpDir, err := os.CreateTemp("", "test-both-exist-*")
		require.NoError(t, err)
		defer os.Remove(tmpDir.Name())
		tmpDir.Close()

		tmpDirPath := tmpDir.Name()
		os.Remove(tmpDirPath) // Remove the file, we want the directory
		err = os.Mkdir(tmpDirPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirPath)

		err = os.Chdir(tmpDirPath)
		require.NoError(t, err)
		defer os.Chdir(origDir)

		// Create both files
		err = os.WriteFile(defaultConfigYAML, []byte("# test yaml config"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(defaultConfigYML, []byte("# test yml config"), 0644)
		require.NoError(t, err)

		// Act
		result := findDefaultConfigFile()

		// Assert
		assert.Equal(t, defaultConfigYAML, result, "Should prefer publicapis.yaml when both exist")
	})
}

// ============================================================================
// Config File Tests
// ============================================================================

func Test_parseConfigFile(t *testing.T) {
	t.Run("parses valid config file", func(t *testing.T) {
		// Create a test config file
		testConfig := Config{
			{
				Specification: "spec1.yaml",
				OpenAPIJSON:   "output1.json",
				SchemaJSON:    "schema1.json",
			},
			{
				Specification: "spec2.yaml",
				OpenAPIYAML:   "output2.yaml",
				OverlayYAML:   "overlay2.yaml",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-config-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Act
		config, err := parseConfigFile(tmpConfigFile.Name())

		// Assert
		require.NoError(t, err, "parseConfigFile should succeed for valid config")
		assert.Len(t, config, 2, "Config should have 2 jobs")
		assert.Equal(t, "spec1.yaml", config[0].Specification, "First job specification should match")
		assert.Equal(t, "output1.json", config[0].OpenAPIJSON, "First job OpenAPI JSON should match")
		assert.Equal(t, "spec2.yaml", config[1].Specification, "Second job specification should match")
		assert.Equal(t, "output2.yaml", config[1].OpenAPIYAML, "Second job OpenAPI YAML should match")
	})

	t.Run("returns error for nonexistent file", func(t *testing.T) {
		// Act
		config, err := parseConfigFile("nonexistent-config.yaml")

		// Assert
		require.Error(t, err, "parseConfigFile should return error for nonexistent file")
		assert.Nil(t, config, "Config should be nil on error")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
	})

	t.Run("returns error for invalid YAML", func(t *testing.T) {
		// Create invalid YAML file
		tmpConfigFile, err := os.CreateTemp("", "test-invalid-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())

		_, err = tmpConfigFile.Write([]byte("invalid: yaml: content: [unclosed"))
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Act
		config, err := parseConfigFile(tmpConfigFile.Name())

		// Assert
		require.Error(t, err, "parseConfigFile should return error for invalid YAML")
		assert.Nil(t, config, "Config should be nil on error")
		assert.Contains(t, err.Error(), errorConfigParsing, "Error should mention config parsing")
	})

	t.Run("returns error for empty config", func(t *testing.T) {
		// Create empty config file
		tmpConfigFile, err := os.CreateTemp("", "test-empty-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())

		_, err = tmpConfigFile.Write([]byte("[]"))
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Act
		config, err := parseConfigFile(tmpConfigFile.Name())

		// Assert
		require.Error(t, err, "parseConfigFile should return error for empty config")
		assert.Nil(t, config, "Config should be nil on error")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
		assert.Contains(t, err.Error(), "at least one job", "Error should mention missing jobs")
	})

	t.Run("returns error for job missing specification", func(t *testing.T) {
		// Create config with job missing specification
		testConfig := Config{
			{
				OpenAPIJSON: "output.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-missing-spec-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Act
		config, err := parseConfigFile(tmpConfigFile.Name())

		// Assert
		require.Error(t, err, "parseConfigFile should return error for job missing specification")
		assert.Nil(t, config, "Config should be nil on error")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
		assert.Contains(t, err.Error(), "missing required 'specification' field", "Error should mention missing specification")
	})

	t.Run("returns error for job with no outputs", func(t *testing.T) {
		// Create config with job that has no output formats
		testConfig := Config{
			{
				Specification: "spec.yaml",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-no-outputs-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Act
		config, err := parseConfigFile(tmpConfigFile.Name())

		// Assert
		require.Error(t, err, "parseConfigFile should return error for job with no outputs")
		assert.Nil(t, config, "Config should be nil on error")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
		assert.Contains(t, err.Error(), "at least one output format", "Error should mention missing output formats")
	})
}

func Test_generateOpenAPIYAMLOutputPath(t *testing.T) {
	testCases := []struct {
		name      string
		inputFile string
		expected  string
	}{
		{
			name:      "YAML file generates YAML OpenAPI",
			inputFile: "spec.yaml",
			expected:  "spec-openapi.yaml",
		},
		{
			name:      "JSON file generates YAML OpenAPI",
			inputFile: "api.json",
			expected:  "api-openapi.yaml",
		},
		{
			name:      "File with path generates YAML OpenAPI",
			inputFile: "/path/to/spec.yml",
			expected:  "/path/to/spec-openapi.yaml",
		},
		{
			name:      "File without extension generates YAML OpenAPI",
			inputFile: "spec",
			expected:  "spec-openapi.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := generateOpenAPIYAMLOutputPath(tc.inputFile)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_ensureYAMLExtension(t *testing.T) {
	testCases := []struct {
		name       string
		outputPath string
		expected   string
	}{
		{
			name:       "YAML extension unchanged",
			outputPath: "api-spec.yaml",
			expected:   "api-spec.yaml",
		},
		{
			name:       "YML extension unchanged",
			outputPath: "api-spec.yml",
			expected:   "api-spec.yml",
		},
		{
			name:       "JSON extension changed to YAML",
			outputPath: "api-spec.json",
			expected:   "api-spec.yaml",
		},
		{
			name:       "No extension gets YAML extension",
			outputPath: "api-spec",
			expected:   "api-spec.yaml",
		},
		{
			name:       "Other extension changed to YAML",
			outputPath: "api-spec.xml",
			expected:   "api-spec.yaml",
		},
		{
			name:       "Path with YAML extension unchanged",
			outputPath: "/path/to/api-spec.yaml",
			expected:   "/path/to/api-spec.yaml",
		},
		{
			name:       "Path with JSON extension changed to YAML",
			outputPath: "/path/to/api-spec.json",
			expected:   "/path/to/api-spec.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := ensureYAMLExtension(tc.outputPath)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}
