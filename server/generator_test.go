package server

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Generator Tests
// ============================================================================

func TestNew(t *testing.T) {
	// Test creating generator with valid configuration
	validConfig := Config{
		OutputPath:  "api/server.go",
		PackageName: "api",
	}

	generator, err := New(validConfig)
	assert.Nil(t, err, "Expected no error when creating generator with valid config")
	assert.NotNil(t, generator, "Expected generator to be created")
	assert.Equal(t, validConfig.OutputPath, generator.config.OutputPath, "Output path should match")
	assert.Equal(t, validConfig.PackageName, generator.config.PackageName, "Package name should match")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("configuration with defaults still validates successfully", func(t *testing.T) {
			configWithEmptyPath := Config{
				OutputPath:  "",
				PackageName: "api",
			}

			generator, err := New(configWithEmptyPath)
			assert.Nil(t, err, "Expected no error as defaults should be applied")
			assert.NotNil(t, generator, "Expected generator to be created")
			assert.Equal(t, defaultOutputPath, generator.config.OutputPath, "Should use default output path")
		})

		t.Run("applies defaults to empty configuration", func(t *testing.T) {
			emptyConfig := Config{}

			generator, err := New(emptyConfig)
			assert.Nil(t, err, "Expected no error when defaults are applied")
			assert.NotNil(t, generator, "Expected generator to be created")
			assert.Equal(t, defaultOutputPath, generator.config.OutputPath, "Should use default output path")
			assert.Equal(t, defaultPackageName, generator.config.PackageName, "Should use default package name")
		})
	})
}

func TestGenerator_GetConfig(t *testing.T) {
	expectedConfig := Config{
		OutputPath:  "custom/server.go",
		PackageName: "customapi",
	}

	generator, err := New(expectedConfig)
	assert.Nil(t, err, "Expected no error creating generator")

	actualConfig := generator.GetConfig()
	assert.Equal(t, expectedConfig.OutputPath, actualConfig.OutputPath, "Output path should match")
	assert.Equal(t, expectedConfig.PackageName, actualConfig.PackageName, "Package name should match")

	t.Run("config independence", func(t *testing.T) {
		// Modify the returned config to ensure it's a copy
		actualConfig.OutputPath = "modified.go"
		actualConfig.PackageName = "modified"

		// Original generator config should be unchanged
		originalConfig := generator.GetConfig()
		assert.Equal(t, expectedConfig.OutputPath, originalConfig.OutputPath, "Original config should be unchanged")
		assert.Equal(t, expectedConfig.PackageName, originalConfig.PackageName, "Original config should be unchanged")
	})
}

func TestGenerator_GenerateFromData(t *testing.T) {
	// Create a simple OpenAPI specification for testing
	simpleOpenAPISpec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/health": {
				"get": {
					"operationId": "getHealth",
					"responses": {
						"200": {
							"description": "OK",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"status": {
												"type": "string"
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	// Create temporary directory for test output
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "server.go")

	config := Config{
		OutputPath:  outputPath,
		PackageName: "testapi",
	}

	generator, err := New(config)
	assert.Nil(t, err, "Expected no error creating generator")

	// Generate code from the OpenAPI specification
	err = generator.GenerateFromData([]byte(simpleOpenAPISpec))
	assert.Nil(t, err, "Expected no error generating code")

	// Verify the output file was created
	assert.FileExists(t, outputPath, "Generated file should exist")

	// Read and verify the generated content
	content, err := os.ReadFile(outputPath)
	assert.Nil(t, err, "Expected no error reading generated file")

	contentStr := string(content)
	assert.Contains(t, contentStr, "package testapi", "Generated code should have correct package name")
	assert.Contains(t, contentStr, "ServerInterface", "Generated code should contain server interface")
	assert.Contains(t, contentStr, "GetHealth", "Generated code should contain GetHealth method")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("invalid OpenAPI specification returns error", func(t *testing.T) {
			invalidSpec := `{"invalid": "json", "missing": "openapi"}`

			tempDir := t.TempDir()
			outputPath := filepath.Join(tempDir, "invalid.go")

			config := Config{
				OutputPath:  outputPath,
				PackageName: "invalid",
			}

			generator, err := New(config)
			assert.Nil(t, err, "Expected no error creating generator")

			err = generator.GenerateFromData([]byte(invalidSpec))
			assert.Error(t, err, "Expected error for invalid OpenAPI specification")
			assert.True(t, strings.Contains(err.Error(), "invalid OpenAPI specification") || strings.Contains(err.Error(), "failed to read OpenAPI specification"), "Error should mention OpenAPI specification issue")
		})

		t.Run("creates output directory if it doesn't exist", func(t *testing.T) {
			tempDir := t.TempDir()
			nestedPath := filepath.Join(tempDir, "nested", "deep", "server.go")

			config := Config{
				OutputPath:  nestedPath,
				PackageName: "nested",
			}

			generator, err := New(config)
			assert.Nil(t, err, "Expected no error creating generator")

			err = generator.GenerateFromData([]byte(simpleOpenAPISpec))
			assert.Nil(t, err, "Expected no error generating code with nested path")

			assert.FileExists(t, nestedPath, "Generated file should exist in nested directory")
		})
	})
}

func TestGenerator_GenerateFromReader(t *testing.T) {
	simpleOpenAPISpec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Reader Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/ping": {
				"get": {
					"operationId": "getPing",
					"responses": {
						"200": {
							"description": "Pong response",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"message": {
												"type": "string"
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	reader := strings.NewReader(simpleOpenAPISpec)
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "reader_server.go")

	config := Config{
		OutputPath:  outputPath,
		PackageName: "readerapi",
	}

	generator, err := New(config)
	assert.Nil(t, err, "Expected no error creating generator")

	err = generator.GenerateFromReader(reader)
	assert.Nil(t, err, "Expected no error generating from reader")

	// Verify the output file was created
	assert.FileExists(t, outputPath, "Generated file should exist")

	// Read and verify the generated content
	content, err := os.ReadFile(outputPath)
	assert.Nil(t, err, "Expected no error reading generated file")

	contentStr := string(content)
	assert.Contains(t, contentStr, "package readerapi", "Generated code should have correct package name")
	assert.Contains(t, contentStr, "GetPing", "Generated code should contain GetPing method")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty reader returns error", func(t *testing.T) {
			emptyReader := bytes.NewReader([]byte{})

			config := Config{
				OutputPath:  filepath.Join(t.TempDir(), "empty.go"),
				PackageName: "empty",
			}

			generator, err := New(config)
			assert.Nil(t, err, "Expected no error creating generator")

			err = generator.GenerateFromReader(emptyReader)
			assert.Error(t, err, "Expected error for empty reader")
		})
	})
}

func TestGenerator_GenerateFromFile(t *testing.T) {
	// Create a temporary OpenAPI spec file
	simpleOpenAPISpec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "File Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/status": {
				"get": {
					"operationId": "getStatus",
					"responses": {
						"200": {
							"description": "Status response",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"healthy": {
												"type": "boolean"
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "api.json")
	outputPath := filepath.Join(tempDir, "file_server.go")

	// Write the spec to file
	err := os.WriteFile(specPath, []byte(simpleOpenAPISpec), 0644)
	assert.Nil(t, err, "Expected no error writing spec file")

	config := Config{
		OutputPath:  outputPath,
		PackageName: "fileapi",
	}

	generator, err := New(config)
	assert.Nil(t, err, "Expected no error creating generator")

	err = generator.GenerateFromFile(specPath)
	assert.Nil(t, err, "Expected no error generating from file")

	// Verify the output file was created
	assert.FileExists(t, outputPath, "Generated file should exist")

	// Read and verify the generated content
	content, err := os.ReadFile(outputPath)
	assert.Nil(t, err, "Expected no error reading generated file")

	contentStr := string(content)
	assert.Contains(t, contentStr, "package fileapi", "Generated code should have correct package name")
	assert.Contains(t, contentStr, "GetStatus", "Generated code should contain GetStatus method")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("non-existent file returns error", func(t *testing.T) {
			nonExistentPath := filepath.Join(t.TempDir(), "does-not-exist.json")

			config := Config{
				OutputPath:  filepath.Join(t.TempDir(), "error.go"),
				PackageName: "error",
			}

			generator, err := New(config)
			assert.Nil(t, err, "Expected no error creating generator")

			err = generator.GenerateFromFile(nonExistentPath)
			assert.Error(t, err, "Expected error for non-existent file")
			assert.Contains(t, err.Error(), "failed to read OpenAPI specification", "Error should mention spec read failure")
		})

		t.Run("handles YAML specification files", func(t *testing.T) {
			yamlSpec := `openapi: 3.0.0
info:
  title: YAML Test API
  version: 1.0.0
paths:
  /yaml:
    get:
      operationId: getYaml
      responses:
        '200':
          description: YAML response
          content:
            application/json:
              schema:
                type: object
                properties:
                  format:
                    type: string`

			tempDir := t.TempDir()
			specPath := filepath.Join(tempDir, "api.yaml")
			outputPath := filepath.Join(tempDir, "yaml_server.go")

			err := os.WriteFile(specPath, []byte(yamlSpec), 0644)
			assert.Nil(t, err, "Expected no error writing YAML spec file")

			config := Config{
				OutputPath:  outputPath,
				PackageName: "yamlapi",
			}

			generator, err := New(config)
			assert.Nil(t, err, "Expected no error creating generator")

			err = generator.GenerateFromFile(specPath)
			assert.Nil(t, err, "Expected no error generating from YAML file")

			assert.FileExists(t, outputPath, "Generated file should exist")

			content, err := os.ReadFile(outputPath)
			assert.Nil(t, err, "Expected no error reading generated file")

			contentStr := string(content)
			assert.Contains(t, contentStr, "package yamlapi", "Generated code should have correct package name")
			assert.Contains(t, contentStr, "GetYaml", "Generated code should contain GetYaml method")
		})
	})
}
