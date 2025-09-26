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

	t.Run("missing config file returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err := run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when no config file is found")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
	})

	t.Run("invalid log level returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Create a test config file
		testConfig := Config{
			{
				Specification: "testdata/school-management-api.yaml",
				OpenAPIJSON:   "test-invalid-log-openapi.json",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-invalid-log-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())
		defer os.Remove("test-invalid-log-openapi.json")

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Arrange
		os.Args = []string{"publicapis-gen", "generate", "-config=" + tmpConfigFile.Name(), "-log-level=invalid"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error for invalid log level")
		assert.Contains(t, err.Error(), "invalid log level", "Error should mention invalid log level")
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
		require.NoError(t, err, "OpenAPI JSON file should be created from config")

		_, err = os.Stat("test-config-schema.json")
		require.NoError(t, err, "Schema JSON file should be created from config")
	})

	t.Run("config file with server generation succeeds", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Create a test config file with server generation
		testConfig := Config{
			{
				Specification: "testdata/school-management-api.yaml",
				ServerGo:      "test-config-server.go",
				ServerPackage: "testapi",
			},
		}

		yamlData, err := yaml.Marshal(&testConfig)
		require.NoError(t, err)

		tmpConfigFile, err := os.CreateTemp("", "test-config-server-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpConfigFile.Name())
		defer os.Remove("test-config-server.go")

		_, err = tmpConfigFile.Write(yamlData)
		require.NoError(t, err)
		tmpConfigFile.Close()

		// Arrange command line arguments for config mode with server generation
		os.Args = []string{"publicapis-gen", "generate", "-config=" + tmpConfigFile.Name()}
		ctx := context.Background()

		// Act - run the command
		err = run(ctx)

		// Assert - command should succeed
		require.NoError(t, err, "run() should not return an error for valid config file with server generation")

		// Verify server output file was created
		_, err = os.Stat("test-config-server.go")
		require.NoError(t, err, "Go server file should be created from config")

		// Read and verify the generated server file contains expected content
		serverContent, err := os.ReadFile("test-config-server.go")
		require.NoError(t, err, "Should be able to read generated server file")

		serverContentStr := string(serverContent)
		// Note: servergen always generates "package api", it doesn't use the ServerPackage field
		assert.Contains(t, serverContentStr, "package api", "Generated server should contain package declaration")
		assert.Contains(t, serverContentStr, "func RegisterSchoolManagementAPIAPI[Session any]", "Generated server should contain registration function")
		assert.Contains(t, serverContentStr, "type SchoolManagementAPIAPI[Session any] struct", "Generated server should contain API struct")
		assert.Contains(t, serverContentStr, "type Student struct", "Generated server should contain Student type")
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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange command line arguments with no flags and no default config
		os.Args = []string{"publicapis-gen", "generate"}
		ctx := context.Background()

		// Act
		err = run(ctx)

		// Assert
		require.Error(t, err, "run() should return an error when no config and no flags provided")
		assert.Contains(t, err.Error(), errorInvalidConfig, "Error should mention invalid config")
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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

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
		defer func() { _ = os.Chdir(origDir) }()

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

// stringPtr returns a pointer to a string value - helper for tests
func stringPtr(s string) *string {
	return &s
}

// Test_generateServerFromSpecification_e2e tests the servergen functionality end-to-end
func Test_generateServerFromSpecification_e2e(t *testing.T) {
	// Test constants to avoid hardcoded strings
	const (
		testServerPackage    = "testapi"
		expectedPackageDecl  = "package api"
		expectedImportGin    = `"github.com/gin-gonic/gin"`
		expectedImportTypes  = `"github.com/meitner-se/go-types"`
		expectedRegisterFunc = "func RegisterTestServiceAPI[Session any]"
		expectedErrorType    = "type Error struct {"
		expectedOpenAPIRoute = `routerGroup.StaticFileFS("/openapi.json", "openapi.json", http.FS(api.OpenAPI_JSON))`
	)

	// Create a test service with various features
	service := &specification.Service{
		Name:    "TestService",
		Version: "v1",
		Enums: []specification.Enum{
			{
				Name:        "Status",
				Description: "Status enumeration",
				Values: []specification.EnumValue{
					{Name: "Active", Description: "Active status"},
					{Name: "Inactive", Description: "Inactive status"},
				},
			},
			{
				Name:        "ErrorCode",
				Description: "Error codes",
				Values: []specification.EnumValue{
					{Name: "BadRequest", Description: "Bad request"},
					{Name: "Unauthorized", Description: "Unauthorized"},
					{Name: "Forbidden", Description: "Forbidden"},
					{Name: "NotFound", Description: "Not found"},
					{Name: "Conflict", Description: "Conflict"},
					{Name: "UnprocessableEntity", Description: "Unprocessable entity"},
					{Name: "RateLimited", Description: "Rate limited"},
					{Name: "Internal", Description: "Internal error"},
				},
			},
		},
		Objects: []specification.Object{
			{
				Name:        "User",
				Description: "User object",
				Fields: []specification.Field{
					{
						Name:        "ID",
						Description: "User ID",
						Type:        "UUID",
					},
					{
						Name:        "Name",
						Description: "User name",
						Type:        "String",
					},
					{
						Name:        "Status",
						Description: "User status",
						Type:        "Status",
					},
				},
			},
			{
				Name:        "Error",
				Description: "Error object",
				Fields: []specification.Field{
					{
						Name:        "Code",
						Description: "Error code",
						Type:        "ErrorCode",
					},
					{
						Name:        "Message",
						Description: "Error message",
						Type:        "String",
					},
					{
						Name:        "RequestID",
						Description: "Request ID",
						Type:        "String",
					},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name:        "User",
				Description: "User resource",
				Operations:  []string{"Create", "Read", "Update", "Delete"},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        "ID",
							Description: "User ID",
							Type:        "UUID",
						},
						Operations: []string{"Read"},
					},
					{
						Field: specification.Field{
							Name:        "Name",
							Description: "User name",
							Type:        "String",
						},
						Operations: []string{"Create", "Read", "Update"},
					},
				},
				Endpoints: []specification.Endpoint{
					{
						Name:        "CreateUser",
						Method:      "POST",
						Path:        "/users",
						Title:       "Create User",
						Summary:     "Create a new user",
						Description: "Creates a new user in the system",
						Request: specification.EndpointRequest{
							BodyParams: []specification.Field{
								{
									Name:        "Name",
									Description: "User name",
									Type:        "String",
								},
							},
						},
						Response: specification.EndpointResponse{
							StatusCode: 201,
							BodyObject: stringPtr("User"),
						},
					},
					{
						Name:        "GetUser",
						Method:      "GET",
						Path:        "/users/{id}",
						Title:       "Get User",
						Summary:     "Get a user by ID",
						Description: "Retrieves a user by their ID",
						Request: specification.EndpointRequest{
							PathParams: []specification.Field{
								{
									Name:        "ID",
									Description: "User ID",
									Type:        "UUID",
								},
							},
						},
						Response: specification.EndpointResponse{
							StatusCode: 200,
							BodyObject: stringPtr("User"),
						},
					},
				},
			},
		},
	}

	// Create a temporary directory for test output
	tempDir := os.TempDir()
	outputPath := tempDir + "/test_server.go"
	defer os.Remove(outputPath) // Clean up after test

	ctx := context.Background()

	// Act - Generate server code
	err := generateServerFromSpecification(ctx, service, "test-spec.yaml", outputPath, testServerPackage)

	// Assert
	assert.Nil(t, err, "Expected no error when generating server code")

	// Verify file was created
	_, statErr := os.Stat(outputPath)
	assert.Nil(t, statErr, "Expected generated file to exist")

	// Read and verify generated content
	content, readErr := os.ReadFile(outputPath)
	assert.Nil(t, readErr, "Expected no error reading generated file")

	generatedCode := string(content)
	assert.NotEmpty(t, generatedCode, "Expected generated code to be non-empty")

	// Verify key components of generated code
	assert.Contains(t, generatedCode, expectedPackageDecl, "Generated code should contain package declaration")
	assert.Contains(t, generatedCode, expectedImportGin, "Generated code should import gin")
	assert.Contains(t, generatedCode, expectedImportTypes, "Generated code should import go-types")
	assert.Contains(t, generatedCode, expectedRegisterFunc, "Generated code should contain RegisterAPI function")
	assert.Contains(t, generatedCode, expectedErrorType, "Generated code should contain Error type")
	assert.Contains(t, generatedCode, expectedOpenAPIRoute, "Generated code should contain OpenAPI route")

	// Verify enum generation
	assert.Contains(t, generatedCode, "type Status types.String", "Generated code should contain Status enum")
	assert.Contains(t, generatedCode, "StatusActive", "Generated code should contain StatusActive")
	assert.Contains(t, generatedCode, "Status(types.NewString(\"Active\"))", "Generated code should contain enum values")

	// Verify object generation
	assert.Contains(t, generatedCode, "type User struct {", "Generated code should contain User object")
	assert.Contains(t, generatedCode, "ID types.UUID `json:\"id\"`", "Generated code should contain User fields")

	// Verify endpoint generation (Note: servergen includes resource name in path)
	assert.Contains(t, generatedCode, "routerGroup.POST(\"/user/users\", serveWithResponse(201, api.Server, api.User.CreateUser))", "Generated code should contain POST endpoint")
	assert.Contains(t, generatedCode, "routerGroup.GET(\"/user/users/:id\", serveWithResponse(200, api.Server, api.User.GetUser))", "Generated code should contain GET endpoint")

	// Verify request/response types
	assert.Contains(t, generatedCode, "type UserCreateUserBodyParams struct {", "Generated code should contain request body type")
	assert.Contains(t, generatedCode, "type UserGetUserPathParams struct {", "Generated code should contain path params type")
	// Note: Response types are not generated when using BodyObject

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty service", func(t *testing.T) {
			// Arrange
			emptyService := &specification.Service{
				Name:    "EmptyService",
				Version: "v1",
			}
			outputPath := tempDir + "/empty_server.go"
			defer os.Remove(outputPath)

			// Act
			err := generateServerFromSpecification(ctx, emptyService, "empty-spec.yaml", outputPath, "emptyapi")

			// Assert
			assert.Nil(t, err, "Expected no error with empty service")

			// Verify file was created
			_, statErr := os.Stat(outputPath)
			assert.Nil(t, statErr, "Expected generated file to exist")

			// Read and verify basic structure
			content, _ := os.ReadFile(outputPath)
			generatedCode := string(content)
			assert.Contains(t, generatedCode, expectedPackageDecl, "Should still generate basic structure")
			assert.Contains(t, generatedCode, "func RegisterEmptyServiceAPI[Session any]", "Should generate register function for empty service")
		})

		t.Run("service with no endpoints", func(t *testing.T) {
			// Arrange
			serviceNoEndpoints := &specification.Service{
				Name:    "NoEndpointService",
				Version: "v1",
				Resources: []specification.Resource{
					{
						Name:        "Resource",
						Description: "Resource with no endpoints",
						Endpoints:   []specification.Endpoint{}, // No endpoints
					},
				},
			}
			outputPath := tempDir + "/no_endpoints_server.go"
			defer os.Remove(outputPath)

			// Act
			err := generateServerFromSpecification(ctx, serviceNoEndpoints, "no-endpoints-spec.yaml", outputPath, "noendpointsapi")

			// Assert
			assert.Nil(t, err, "Expected no error with no endpoints")

			// Verify file was created
			_, statErr := os.Stat(outputPath)
			assert.Nil(t, statErr, "Expected generated file to exist")
		})
	})
}

// ============================================================================
// Diff Functionality Tests
// ============================================================================

func TestCompareWithDiskFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "diff_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("identical files return no difference", func(t *testing.T) {
		// Arrange
		content := "line 1\nline 2\nline 3\n"
		testFile := tempDir + "/identical.txt"
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		// Act
		diff, err := compareWithDiskFile(testFile, []byte(content))

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, diff, "Identical files should produce no diff output")
	})

	t.Run("non-existent file returns appropriate message", func(t *testing.T) {
		// Arrange
		nonExistentFile := tempDir + "/does_not_exist.txt"
		content := []byte("some content")

		// Act
		diff, err := compareWithDiskFile(nonExistentFile, content)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, diffFileNotExist, diff)
	})

	t.Run("different files show detailed diff", func(t *testing.T) {
		// Arrange
		diskContent := "line 1\noriginal line 2\nline 3\nline 4\n"
		generatedContent := "line 1\nmodified line 2\nline 3\nline 4\n"
		testFile := tempDir + "/different.txt"
		err := os.WriteFile(testFile, []byte(diskContent), 0644)
		require.NoError(t, err)

		// Act
		diff, err := compareWithDiskFile(testFile, []byte(generatedContent))

		// Assert
		assert.NoError(t, err)
		assert.Contains(t, diff, diffContentDiffers, "Should indicate content differs")
		assert.Contains(t, diff, "First difference found at line 2", "Should show line number of first difference")
		assert.Contains(t, diff, diffGeneratedHeader, "Should show generated content header")
		assert.Contains(t, diff, diffDiskHeader, "Should show disk content header")
		assert.Contains(t, diff, "modified line 2", "Should show generated content")
		assert.Contains(t, diff, "original line 2", "Should show disk content")
		assert.Contains(t, diff, "→ 2:", "Should mark the differing line with arrow")
	})

	t.Run("files with different lengths show diff", func(t *testing.T) {
		// Arrange
		diskContent := "line 1\nline 2\n"
		generatedContent := "line 1\nline 2\nline 3\nextra line\n"
		testFile := tempDir + "/different_length.txt"
		err := os.WriteFile(testFile, []byte(diskContent), 0644)
		require.NoError(t, err)

		// Act
		diff, err := compareWithDiskFile(testFile, []byte(generatedContent))

		// Assert
		assert.NoError(t, err)
		assert.Contains(t, diff, diffContentDiffers, "Should indicate content differs")
		assert.Contains(t, diff, "First difference found at line 3", "Should show first difference at end of shorter file")
		assert.Contains(t, diff, "more line(s) differ", "Should indicate additional differences")
	})
}

func TestGenerateDetailedDiff(t *testing.T) {
	t.Run("single line difference", func(t *testing.T) {
		// Arrange
		generated := []byte("line 1\nmodified line\nline 3")
		disk := []byte("line 1\noriginal line\nline 3")

		// Act
		diff := generateDetailedDiff(generated, disk)

		// Assert
		assert.Contains(t, diff, "First difference found at line 2", "Should identify correct line number")
		assert.Contains(t, diff, "modified line", "Should show generated content")
		assert.Contains(t, diff, "original line", "Should show disk content")
		assert.Contains(t, diff, "→ 2:", "Should mark differing line")
	})

	t.Run("difference at beginning", func(t *testing.T) {
		// Arrange
		generated := []byte("modified first line\nline 2\nline 3")
		disk := []byte("original first line\nline 2\nline 3")

		// Act
		diff := generateDetailedDiff(generated, disk)

		// Assert
		assert.Contains(t, diff, "First difference found at line 1", "Should identify first line difference")
		assert.Contains(t, diff, "→ 1:", "Should mark first line")
	})

	t.Run("multiple differences shows additional count", func(t *testing.T) {
		// Arrange
		generated := []byte("line 1\nmodified 2\nline 3\nmodified 4\nmodified 5")
		disk := []byte("line 1\noriginal 2\nline 3\noriginal 4\noriginal 5")

		// Act
		diff := generateDetailedDiff(generated, disk)

		// Assert
		assert.Contains(t, diff, "First difference found at line 2", "Should show first difference")
		assert.Contains(t, diff, "more line(s) differ", "Should indicate additional differences")
	})

	t.Run("files with different lengths", func(t *testing.T) {
		// Arrange
		generated := []byte("line 1\nline 2\nextra line")
		disk := []byte("line 1\nline 2")

		// Act
		diff := generateDetailedDiff(generated, disk)

		// Assert
		assert.Contains(t, diff, "First difference found at line 3", "Should identify length difference")
		assert.Contains(t, diff, "more line(s) differ", "Should indicate length difference")
	})
}

func TestFindFirstDifference(t *testing.T) {
	t.Run("identical content returns -1", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "line 2", "line 3"}
		lines2 := []string{"line 1", "line 2", "line 3"}

		// Act
		diff := findFirstDifference(lines1, lines2)

		// Assert
		assert.Equal(t, -1, diff, "Identical content should return -1")
	})

	t.Run("different content returns correct line", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "modified", "line 3"}
		lines2 := []string{"line 1", "original", "line 3"}

		// Act
		diff := findFirstDifference(lines1, lines2)

		// Assert
		assert.Equal(t, 1, diff, "Should return index of first different line")
	})

	t.Run("different lengths returns correct position", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "line 2", "extra"}
		lines2 := []string{"line 1", "line 2"}

		// Act
		diff := findFirstDifference(lines1, lines2)

		// Assert
		assert.Equal(t, 2, diff, "Should return position where lengths differ")
	})
}

func TestCountAdditionalDifferences(t *testing.T) {
	t.Run("no additional differences", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "different", "line 3"}
		lines2 := []string{"line 1", "original", "line 3"}

		// Act
		count := countAdditionalDifferences(lines1, lines2, 2) // Start after first diff at index 1

		// Assert
		assert.Equal(t, 0, count, "Should count no additional differences")
	})

	t.Run("multiple differences", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "diff1", "line 3", "diff2", "diff3"}
		lines2 := []string{"line 1", "orig1", "line 3", "orig2", "orig3"}

		// Act
		count := countAdditionalDifferences(lines1, lines2, 2) // Start after first diff at index 1

		// Assert
		assert.Equal(t, 2, count, "Should count remaining differences")
	})

	t.Run("different lengths", func(t *testing.T) {
		// Arrange
		lines1 := []string{"line 1", "diff", "line 3", "extra1", "extra2"}
		lines2 := []string{"line 1", "orig", "line 3"}

		// Act
		count := countAdditionalDifferences(lines1, lines2, 2) // Start after first diff at index 1

		// Assert
		assert.Equal(t, 2, count, "Should count extra lines as differences")
	})
}
