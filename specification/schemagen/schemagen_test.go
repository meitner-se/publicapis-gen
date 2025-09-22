package schemagen

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// GenerateSchemas Tests
// ============================================================================

func TestGenerateSchemas(t *testing.T) {
	// Test the main happy path
	var buf bytes.Buffer
	err := GenerateSchemas(&buf)

	require.NoError(t, err, "Expected no error when generating schemas")
	assert.NotEmpty(t, buf.String(), "Buffer should contain generated schemas")

	// Parse the generated JSON
	var schemas map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &schemas)
	require.NoError(t, err, "Generated output should be valid JSON")

	// Check that all expected schemas are present
	expectedSchemas := []string{
		"Service", "Enum", "Object", "Resource",
		"Field", "ResourceField", "Endpoint",
		"EndpointRequest", "EndpointResponse",
	}

	assert.Equal(t, len(expectedSchemas), len(schemas), "Should generate exactly %d schemas", len(expectedSchemas))

	for _, expectedSchemaName := range expectedSchemas {
		schema, exists := schemas[expectedSchemaName]
		assert.True(t, exists, "Schema '%s' should be present in results", expectedSchemaName)
		assert.NotNil(t, schema, "Schema '%s' should not be nil", expectedSchemaName)

		// Each schema should have basic structure
		schemaMap, ok := schema.(map[string]interface{})
		require.True(t, ok, "Schema '%s' should be a map", expectedSchemaName)
		assert.Contains(t, schemaMap, "type", "Schema '%s' should have a type field", expectedSchemaName)
		assert.Contains(t, schemaMap, "properties", "Schema '%s' should have properties field", expectedSchemaName)
	}

	t.Run("edge cases", func(t *testing.T) {
		t.Run("multiple calls generate identical schemas", func(t *testing.T) {
			var buf1, buf2 bytes.Buffer

			err1 := GenerateSchemas(&buf1)
			require.NoError(t, err1)

			err2 := GenerateSchemas(&buf2)
			require.NoError(t, err2)

			assert.JSONEq(t, buf1.String(), buf2.String(), "Multiple calls should generate identical schemas")
		})

		t.Run("buffer reuse", func(t *testing.T) {
			var buf bytes.Buffer

			// First call
			err := GenerateSchemas(&buf)
			require.NoError(t, err)
			firstOutput := buf.String()

			// Reset buffer and call again
			buf.Reset()
			err = GenerateSchemas(&buf)
			require.NoError(t, err)
			secondOutput := buf.String()

			assert.JSONEq(t, firstOutput, secondOutput, "Reusing buffer should work correctly")
		})

		t.Run("buffer with existing content", func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString("existing content")

			err := GenerateSchemas(&buf)
			require.NoError(t, err)

			// The buffer should now contain existing content + schemas
			output := buf.String()
			assert.True(t, len(output) > len("existing content"), "Buffer should have new content appended")

			// The JSON part should be valid
			jsonStart := bytes.Index(buf.Bytes(), []byte("{"))
			assert.True(t, jsonStart >= 0, "Should find JSON start")

			var schemas map[string]interface{}
			err = json.Unmarshal(buf.Bytes()[jsonStart:], &schemas)
			assert.NoError(t, err, "JSON part should be valid")
		})
	})
}

// ============================================================================
// Schema Structure Tests
// ============================================================================

func TestGenerateSchemasStructure(t *testing.T) {
	var buf bytes.Buffer
	err := GenerateSchemas(&buf)
	require.NoError(t, err)

	var schemas map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &schemas)
	require.NoError(t, err)

	// Test Service schema structure
	t.Run("Service schema has correct properties", func(t *testing.T) {
		serviceSchema := schemas["Service"].(map[string]interface{})
		properties := serviceSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "enums", "objects", "resources"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Service schema should have '%s' property", prop)
		}
	})

	// Test Enum schema structure
	t.Run("Enum schema has correct properties", func(t *testing.T) {
		enumSchema := schemas["Enum"].(map[string]interface{})
		properties := enumSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "description", "values"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Enum schema should have '%s' property", prop)
		}
	})

	// Test Object schema structure
	t.Run("Object schema has correct properties", func(t *testing.T) {
		objectSchema := schemas["Object"].(map[string]interface{})
		properties := objectSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "description", "fields"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Object schema should have '%s' property", prop)
		}
	})

	// Test Resource schema structure
	t.Run("Resource schema has correct properties", func(t *testing.T) {
		resourceSchema := schemas["Resource"].(map[string]interface{})
		properties := resourceSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "description", "operations", "fields", "endpoints"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Resource schema should have '%s' property", prop)
		}
	})

	// Test Field schema structure
	t.Run("Field schema has correct properties", func(t *testing.T) {
		fieldSchema := schemas["Field"].(map[string]interface{})
		properties := fieldSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "description", "type", "default", "example", "modifiers"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Field schema should have '%s' property", prop)
		}
	})

	// Test ResourceField has both Field properties and operations
	t.Run("ResourceField schema completeness", func(t *testing.T) {
		resourceFieldSchema := schemas["ResourceField"].(map[string]interface{})
		properties := resourceFieldSchema["properties"].(map[string]interface{})

		// Should have Field properties
		fieldProperties := []string{"name", "description", "type"}
		for _, prop := range fieldProperties {
			assert.Contains(t, properties, prop, "ResourceField should have Field property '%s'", prop)
		}

		// Should also have operations
		assert.Contains(t, properties, "operations", "ResourceField should have operations property")
	})

	// Test Endpoint schema structure
	t.Run("Endpoint schema has correct properties", func(t *testing.T) {
		endpointSchema := schemas["Endpoint"].(map[string]interface{})
		properties := endpointSchema["properties"].(map[string]interface{})

		expectedProperties := []string{"name", "title", "description", "method", "path", "request", "response"}
		for _, prop := range expectedProperties {
			assert.Contains(t, properties, prop, "Endpoint schema should have '%s' property", prop)
		}
	})

	// Test EndpointRequest has all parameter types
	t.Run("EndpointRequest schema completeness", func(t *testing.T) {
		endpointRequestSchema := schemas["EndpointRequest"].(map[string]interface{})
		properties := endpointRequestSchema["properties"].(map[string]interface{})

		paramTypes := []string{"headers", "path_params", "query_params", "body_params"}
		for _, paramType := range paramTypes {
			assert.Contains(t, properties, paramType, "EndpointRequest should have '%s' property", paramType)
		}

		assert.Contains(t, properties, "content_type", "EndpointRequest should have content_type property")
	})

	// Test EndpointResponse has all fields
	t.Run("EndpointResponse schema completeness", func(t *testing.T) {
		endpointResponseSchema := schemas["EndpointResponse"].(map[string]interface{})
		properties := endpointResponseSchema["properties"].(map[string]interface{})

		expectedFields := []string{"content_type", "status_code", "headers", "body_fields", "body_object", "description"}
		for _, field := range expectedFields {
			assert.Contains(t, properties, field, "EndpointResponse should have '%s' property", field)
		}
	})
}

// ============================================================================
// JSON Formatting Tests
// ============================================================================

func TestGenerateSchemasFormattedJSON(t *testing.T) {
	var buf bytes.Buffer
	err := GenerateSchemas(&buf)
	require.NoError(t, err)

	// Check that output is indented (contains newlines and spaces)
	output := buf.String()
	assert.Contains(t, output, "\n", "Output should contain newlines")
	assert.Contains(t, output, "  ", "Output should be indented")

	// Verify it can be unmarshaled and remarshaled
	var schemas map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &schemas)
	require.NoError(t, err)

	// Re-marshal with indentation to compare
	remarshaled, err := json.MarshalIndent(schemas, "", "  ")
	require.NoError(t, err)

	// The output should be the same (normalized comparison)
	assert.JSONEq(t, string(remarshaled), output, "Output should be properly formatted JSON")
}

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestGenerateSchemasValidSchema(t *testing.T) {
	var buf bytes.Buffer
	err := GenerateSchemas(&buf)
	require.NoError(t, err)

	var schemas map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &schemas)
	require.NoError(t, err)

	// Each schema should have basic JSON Schema properties
	for name, schema := range schemas {
		schemaMap := schema.(map[string]interface{})

		// Should have $schema property
		assert.Contains(t, schemaMap, "$schema", "Schema '%s' should have $schema property", name)

		// Should have type property
		assert.Contains(t, schemaMap, "type", "Schema '%s' should have type property", name)
		assert.Equal(t, "object", schemaMap["type"], "Schema '%s' should be of type object", name)

		// Should have properties
		assert.Contains(t, schemaMap, "properties", "Schema '%s' should have properties", name)
		properties, ok := schemaMap["properties"].(map[string]interface{})
		assert.True(t, ok, "Schema '%s' properties should be a map", name)
		assert.NotEmpty(t, properties, "Schema '%s' should have at least one property", name)

		// Should have required fields
		if required, hasRequired := schemaMap["required"]; hasRequired {
			requiredArray, ok := required.([]interface{})
			assert.True(t, ok, "Schema '%s' required should be an array", name)
			assert.NotEmpty(t, requiredArray, "Schema '%s' required array should not be empty", name)
		}

		// Test that additionalProperties is set correctly
		if additionalProps, exists := schemaMap["additionalProperties"]; exists {
			assert.Equal(t, false, additionalProps, "Schema '%s' should not allow additional properties", name)
		}
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestSchemaToJSON(t *testing.T) {
	// This is now an internal function, so we test it indirectly through GenerateSchemas
	var buf bytes.Buffer
	err := GenerateSchemas(&buf)
	require.NoError(t, err)

	// If GenerateSchemas works and produces valid JSON, schemaToJSON must be working
	var schemas map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &schemas)
	require.NoError(t, err, "schemaToJSON must be producing valid JSON")
}
