package schemagen

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meitner-se/publicapis-gen/specification"
)

func TestGenerateSchemas(t *testing.T) {
	// Test basic schema generation
	t.Run("generates all schemas successfully", func(t *testing.T) {
		var buf bytes.Buffer
		service := &specification.Service{
			Name: "TestService",
		}

		err := GenerateSchemas(&buf, service)
		assert.NoError(t, err, "Should generate schemas successfully")
		assert.NotEmpty(t, buf.String(), "Buffer should contain generated schemas")

		// Parse and verify JSON structure
		var schemas map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &schemas)
		assert.NoError(t, err, "Should be valid JSON")

		// Check that all expected schema types are present
		expectedSchemas := []string{"Service", "Enum", "Object", "Resource", "Field", "ResourceField", "Endpoint", "EndpointRequest", "EndpointResponse"}
		for _, schemaName := range expectedSchemas {
			_, ok := schemas[schemaName]
			assert.True(t, ok, "Should have %s schema", schemaName)
		}
	})

	// Test that generated schemas have correct structure
	t.Run("Service schema has correct properties", func(t *testing.T) {
		var buf bytes.Buffer
		service := &specification.Service{
			Name: "TestService",
		}

		err := GenerateSchemas(&buf, service)
		require.NoError(t, err, "Should generate schemas successfully")

		// Parse schemas
		var schemas map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &schemas)
		require.NoError(t, err, "Should be valid JSON")

		// Check Service schema
		serviceSchema, ok := schemas["Service"].(map[string]interface{})
		require.True(t, ok, "Should have Service schema")

		assert.Equal(t, "object", serviceSchema["type"], "Service schema should have object type")

		properties, ok := serviceSchema["properties"].(map[string]interface{})
		assert.True(t, ok, "Service schema should have properties")

		expectedProperties := []string{"name", "enums", "objects", "resources"}
		for _, prop := range expectedProperties {
			_, ok := properties[prop]
			assert.True(t, ok, "Service schema should have '%s' property", prop)
		}
	})

	// Test with nil service
	t.Run("nil service generates schemas without error", func(t *testing.T) {
		var buf bytes.Buffer

		// The function should handle nil service gracefully
		err := GenerateSchemas(&buf, nil)
		assert.NoError(t, err, "Should handle nil service without error")
		assert.NotEmpty(t, buf.String(), "Buffer should contain generated schemas")

		// Parse and verify JSON structure
		var schemas map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &schemas)
		assert.NoError(t, err, "Should be valid JSON")

		// All schema types should still be present
		assert.Greater(t, len(schemas), 0, "Should have generated schemas")
	})

	// Test schema validation capabilities
	t.Run("generated schemas are valid JSON Schema", func(t *testing.T) {
		var buf bytes.Buffer
		service := &specification.Service{
			Name: "TestService",
			Enums: []specification.Enum{
				{
					Name:        "Status",
					Description: "Status enumeration",
					Values: []specification.EnumValue{
						{Name: "Active", Description: "Active status"},
						{Name: "Inactive", Description: "Inactive status"},
					},
				},
			},
			Objects: []specification.Object{
				{
					Name:        "User",
					Description: "User object",
					Fields: []specification.Field{
						{
							Name:        "id",
							Type:        specification.FieldTypeUUID,
							Description: "User ID",
						},
						{
							Name:        "name",
							Type:        specification.FieldTypeString,
							Description: "User name",
						},
					},
				},
			},
			Resources: []specification.Resource{
				{
					Name:        "User",
					Description: "User resource",
					Operations:  []string{specification.OperationCreate, specification.OperationRead},
					Fields: []specification.ResourceField{
						{
							Field: specification.Field{
								Name:        "id",
								Type:        specification.FieldTypeUUID,
								Description: "User ID",
							},
						},
					},
				},
			},
		}

		err := GenerateSchemas(&buf, service)
		assert.NoError(t, err, "Should generate schemas successfully")

		// Parse schemas
		var schemas map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &schemas)
		assert.NoError(t, err, "Should be valid JSON")

		// Verify each schema has required JSON Schema fields
		for schemaName, schemaData := range schemas {
			schema, ok := schemaData.(map[string]interface{})
			assert.True(t, ok, "%s should be a valid schema object", schemaName)

			// Check for type field (required in JSON Schema)
			_, hasType := schema["type"]
			assert.True(t, hasType, "%s schema should have a type field", schemaName)
		}
	})
}

// TestSchemaConsistency verifies that the generated schemas maintain consistency
func TestSchemaConsistency(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	service := &specification.Service{
		Name:        "ConsistencyTest",
		Description: "Test service for consistency",
	}

	// Generate schemas twice
	err1 := GenerateSchemas(&buf1, service)
	err2 := GenerateSchemas(&buf2, service)

	assert.NoError(t, err1, "First generation should succeed")
	assert.NoError(t, err2, "Second generation should succeed")

	// The output should be identical
	assert.Equal(t, buf1.String(), buf2.String(), "Schema generation should be deterministic")
}
