package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meitner-se/publicapis-gen/specification"
)

// ============================================================================
// Generator Tests
// ============================================================================

func TestNewGenerator(t *testing.T) {
	config := NewDefaultConfig()
	generator := NewGenerator(config)

	assert.NotNil(t, generator)
	assert.Equal(t, config, generator.config)
}

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	assert.Equal(t, defaultPackageName, config.PackageName)
	assert.Equal(t, defaultOutputFile, config.OutputFile)
	assert.True(t, config.GenerateTypes)
	assert.True(t, config.GenerateSpec)
}

func TestGenerator_ValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := NewDefaultConfig()
		generator := NewGenerator(config)

		err := generator.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("empty package name", func(t *testing.T) {
		config := NewDefaultConfig()
		config.PackageName = ""
		generator := NewGenerator(config)

		err := generator.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "package name cannot be empty")
	})

	t.Run("empty output file", func(t *testing.T) {
		config := NewDefaultConfig()
		config.OutputFile = ""
		generator := NewGenerator(config)

		err := generator.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "output file cannot be empty")
	})

	t.Run("types disabled", func(t *testing.T) {
		config := NewDefaultConfig()
		config.GenerateTypes = false
		generator := NewGenerator(config)

		err := generator.ValidateConfig()
		assert.NoError(t, err) // Should still be valid since Gin server is hardcoded
	})
}

func TestGenerator_GenerateFromDocument(t *testing.T) {
	t.Run("nil document returns error", func(t *testing.T) {
		config := NewDefaultConfig()
		generator := NewGenerator(config)

		code, err := generator.GenerateFromDocument(nil)
		assert.Error(t, err)
		assert.Nil(t, code)
		assert.Contains(t, err.Error(), errorInvalidDocument)
	})

	t.Run("valid document generates code", func(t *testing.T) {
		// Generate valid OpenAPI JSON for testing
		validOpenAPIJSON := []byte(`{
			"openapi": "3.0.3",
			"info": {
				"title": "Test API",
				"version": "1.0.0"
			},
			"paths": {
				"/users": {
					"get": {
						"operationId": "getUsers",
						"responses": {
							"200": {
								"description": "Success"
							}
						}
					}
				}
			}
		}`)

		// Generate server code
		config := NewDefaultConfig()
		serverGenerator := NewGenerator(config)

		code, err := serverGenerator.GenerateFromDocument(validOpenAPIJSON)
		assert.NoError(t, err)
		assert.NotNil(t, code)

		// Verify generated code contains expected elements
		codeStr := string(code)
		assert.Contains(t, codeStr, "package "+config.PackageName)
		assert.Contains(t, codeStr, "ServerInterface") // Gin server interface
	})
}

func TestGenerator_GenerateFromService(t *testing.T) {
	t.Run("nil service returns error", func(t *testing.T) {
		config := NewDefaultConfig()
		generator := NewGenerator(config)

		code, err := generator.GenerateFromService(nil)
		assert.Error(t, err)
		assert.Nil(t, code)
		assert.Contains(t, err.Error(), "invalid service")
	})

	t.Run("valid service generates code", func(t *testing.T) {
		// Create a simple test service
		service := createTestService()

		// Generate server code directly from service
		config := NewDefaultConfig()
		serverGenerator := NewGenerator(config)

		code, err := serverGenerator.GenerateFromService(service)
		assert.NoError(t, err)
		assert.NotNil(t, code)

		// Verify generated code contains expected elements
		codeStr := string(code)
		assert.Contains(t, codeStr, "package "+config.PackageName)
		assert.Contains(t, codeStr, "ServerInterface") // Gin server interface
	})
}

func TestGenerator_GenerateToBuffer(t *testing.T) {
	// Generate valid OpenAPI JSON for testing
	validOpenAPIJSON := []byte(`{
		"openapi": "3.0.3",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/users": {
				"get": {
					"operationId": "getUsers",
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		}
	}`)

	// Generate server code to buffer
	config := NewDefaultConfig()
	serverGenerator := NewGenerator(config)

	buffer, err := serverGenerator.GenerateToBuffer(validOpenAPIJSON)
	assert.NoError(t, err)
	assert.NotNil(t, buffer)
	assert.Greater(t, buffer.Len(), 0)

	// Verify buffer contains expected content
	codeStr := buffer.String()
	assert.Contains(t, codeStr, "package "+config.PackageName)
}

func TestGenerator_GenerateToWriter(t *testing.T) {
	// Generate valid OpenAPI JSON for testing
	validOpenAPIJSON := []byte(`{
		"openapi": "3.0.3",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/users": {
				"get": {
					"operationId": "getUsers",
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		}
	}`)

	// Generate server code to string builder
	config := NewDefaultConfig()
	serverGenerator := NewGenerator(config)

	var builder strings.Builder
	err := serverGenerator.GenerateToWriter(validOpenAPIJSON, &builder)
	assert.NoError(t, err)

	// Verify writer contains expected content
	codeStr := builder.String()
	assert.Contains(t, codeStr, "package "+config.PackageName)
	assert.Greater(t, len(codeStr), 0)
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestService creates a simple test service for testing server generation.
func createTestService() *specification.Service {
	return &specification.Service{
		Name:    "Test API",
		Version: "1.0.0",
		Resources: []specification.Resource{
			{
				Name:        "User",
				Description: "User management resource",
				Operations:  []string{specification.OperationCreate, specification.OperationRead},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        "Name",
							Description: "User name",
							Type:        specification.FieldTypeString,
							Example:     "John Doe",
						},
						Operations: []string{specification.OperationCreate, specification.OperationRead},
					},
					{
						Field: specification.Field{
							Name:        "Email",
							Description: "User email",
							Type:        specification.FieldTypeString,
							Example:     "john@example.com",
						},
						Operations: []string{specification.OperationCreate, specification.OperationRead},
					},
				},
			},
		},
	}
}
