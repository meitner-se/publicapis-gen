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

	t.Run("invalid log level returns error", func(t *testing.T) {
		// Reset flag package for this test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		// Arrange
		os.Args = []string{"publicapis-gen", "-file=test.yaml", "-mode=overlay", "-log-level=invalid"}
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
		os.Args = []string{"publicapis-gen", "-file=nonexistent.yaml", "-mode=overlay"}
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

		// Create temporary YAML specification file
		yamlContent := `name: "School Management API"
version: "1.0.0"
servers:
  - url: "https://api.school.example.com/v1"
    description: "Production server for School Management API"

enums:
  - name: "StudentStatus"
    description: "Status of a student in the system"
    values:
      - name: "Active"
        description: "Student is actively enrolled"
      - name: "Inactive"  
        description: "Student is not currently enrolled"
      - name: "Graduated"
        description: "Student has graduated"

  - name: "GradeLevel"
    description: "Grade levels in the school"
    values:
      - name: "Elementary"
        description: "Elementary school level"
      - name: "Middle"
        description: "Middle school level"
      - name: "High"
        description: "High school level"

objects:
  - name: "Contact"
    description: "Contact information"
    fields:
      - name: "email"
        type: "String"
        description: "Email address"
      - name: "phone"
        type: "String"
        description: "Phone number"
        modifiers: ["nullable"]

  - name: "Address"
    description: "Physical address information"  
    fields:
      - name: "street"
        type: "String"
        description: "Street address"
      - name: "city"
        type: "String"
        description: "City"
      - name: "state"
        type: "String"
        description: "State or province"
      - name: "zipCode"
        type: "String" 
        description: "ZIP or postal code"

resources:
  - name: "Students"
    description: "Student management resource"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - name: "id"
        type: "UUID"
        description: "Unique student identifier"
        operations: ["Read"]
      - name: "firstName"
        type: "String"
        description: "Student's first name"
        operations: ["Create", "Read", "Update"]
      - name: "lastName"
        type: "String"
        description: "Student's last name"
        operations: ["Create", "Read", "Update"]
      - name: "studentId"
        type: "String"
        description: "School-assigned student ID"
        operations: ["Create", "Read"]
      - name: "status"
        type: "StudentStatus"
        description: "Current status of the student"
        default: "Active"
        operations: ["Create", "Read", "Update"]
      - name: "gradeLevel"
        type: "GradeLevel"
        description: "Current grade level"
        operations: ["Create", "Read", "Update"]
      - name: "contact"
        type: "Contact"
        description: "Student contact information"
        operations: ["Create", "Read", "Update"]
      - name: "address"
        type: "Address"
        description: "Student home address"
        modifiers: ["nullable"]
        operations: ["Create", "Read", "Update"]
      - name: "enrollmentDate"
        type: "Date"
        description: "Date when student was enrolled"
        operations: ["Create", "Read"]
      - name: "graduationDate"
        type: "Date"
        description: "Expected or actual graduation date"
        modifiers: ["nullable"]
        operations: ["Read", "Update"]`

		// Create temporary YAML file
		tmpYAMLFile, err := os.CreateTemp("", "test-spec-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpYAMLFile.Name())

		_, err = tmpYAMLFile.Write([]byte(yamlContent))
		require.NoError(t, err)
		tmpYAMLFile.Close()

		// Determine expected output file path
		expectedOutputPath := generateOpenAPIOutputPath(tmpYAMLFile.Name())
		defer os.Remove(expectedOutputPath) // Clean up output file

		// Arrange command line arguments for OpenAPI generation
		os.Args = []string{"publicapis-gen", "-file=" + tmpYAMLFile.Name(), "-mode=openapi"}
		ctx := context.Background()

		// Act - run the command
		err = run(ctx)

		// Assert - command should succeed
		require.NoError(t, err, "run() should not return an error for valid YAML to OpenAPI conversion")

		// Verify output file was created
		_, err = os.Stat(expectedOutputPath)
		require.NoError(t, err, "Output OpenAPI JSON file should be created")

		// Read and verify the generated JSON content
		outputData, err := os.ReadFile(expectedOutputPath)
		require.NoError(t, err, "Should be able to read generated OpenAPI JSON file")

		// Parse JSON to ensure it's valid
		var openAPIDoc map[string]interface{}
		err = json.Unmarshal(outputData, &openAPIDoc)
		require.NoError(t, err, "Generated file should contain valid JSON")

		// Verify basic OpenAPI structure
		assert.Equal(t, "3.1.0", openAPIDoc["openapi"], "Should have correct OpenAPI version")

		info, ok := openAPIDoc["info"].(map[string]interface{})
		require.True(t, ok, "Should have info section")
		assert.Equal(t, "School Management API API", info["title"], "Should have correct API title")
		assert.Equal(t, "1.0.0", info["version"], "Should have correct version")

		// Verify servers section
		servers, ok := openAPIDoc["servers"].([]interface{})
		require.True(t, ok, "Should have servers section")
		require.Len(t, servers, 1, "Should have one server")
		server := servers[0].(map[string]interface{})
		assert.Equal(t, "https://api.school.example.com/v1", server["url"], "Should have correct server URL")

		// Verify paths section exists and contains expected endpoints
		paths, ok := openAPIDoc["paths"].(map[string]interface{})
		require.True(t, ok, "Should have paths section")

		// Check for generated CRUD endpoints
		assert.Contains(t, paths, "/students", "Should have students collection endpoint")
		assert.Contains(t, paths, "/students/{id}", "Should have students individual resource endpoint")
		assert.Contains(t, paths, "/students/_search", "Should have students search endpoint")

		// Verify components/schemas section
		components, ok := openAPIDoc["components"].(map[string]interface{})
		require.True(t, ok, "Should have components section")
		schemas, ok := components["schemas"].(map[string]interface{})
		require.True(t, ok, "Should have schemas section")

		// Verify enum schemas
		assert.Contains(t, schemas, "StudentStatus", "Should have StudentStatus enum schema")
		assert.Contains(t, schemas, "GradeLevel", "Should have GradeLevel enum schema")

		// Verify object schemas
		assert.Contains(t, schemas, "Contact", "Should have Contact object schema")
		assert.Contains(t, schemas, "Address", "Should have Address object schema")

		// Verify resource schema
		assert.Contains(t, schemas, "Students", "Should have Students resource schema")

		// Verify overlay-generated schemas (error handling, pagination, etc.)
		assert.Contains(t, schemas, "Error", "Should have Error schema")
		assert.Contains(t, schemas, "ErrorCode", "Should have ErrorCode enum schema")
		assert.Contains(t, schemas, "Pagination", "Should have Pagination schema")
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
