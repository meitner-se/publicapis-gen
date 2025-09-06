package specification

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewSchemaGenerator(t *testing.T) {
	generator := NewSchemaGenerator()
	assert.NotNil(t, generator)
	assert.NotNil(t, generator.reflector)
}

func Test_GenerateServiceSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateServiceSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that it's a valid schema
	assert.NotEmpty(t, schema.Type)
	assert.NotNil(t, schema.Properties)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasEnums := schema.Properties.Get("enums")
	assert.True(t, hasEnums)
	_, hasObjects := schema.Properties.Get("objects")
	assert.True(t, hasObjects)
	_, hasResources := schema.Properties.Get("resources")
	assert.True(t, hasResources)
}

func Test_GenerateEnumSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateEnumSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasValues := schema.Properties.Get("values")
	assert.True(t, hasValues)
}

func Test_GenerateObjectSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateObjectSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasFields := schema.Properties.Get("fields")
	assert.True(t, hasFields)
}

func Test_GenerateResourceSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateResourceSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasOperations := schema.Properties.Get("operations")
	assert.True(t, hasOperations)
	_, hasFields := schema.Properties.Get("fields")
	assert.True(t, hasFields)
	_, hasEndpoints := schema.Properties.Get("endpoints")
	assert.True(t, hasEndpoints)
}

func Test_GenerateFieldSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateFieldSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasType := schema.Properties.Get("type")
	assert.True(t, hasType)
	_, hasDefault := schema.Properties.Get("default")
	assert.True(t, hasDefault)
	_, hasExample := schema.Properties.Get("example")
	assert.True(t, hasExample)
	_, hasModifiers := schema.Properties.Get("modifiers")
	assert.True(t, hasModifiers)
}

func Test_GenerateResourceFieldSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateResourceFieldSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present (should include Field properties plus operations)
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasType := schema.Properties.Get("type")
	assert.True(t, hasType)
	_, hasOperations := schema.Properties.Get("operations")
	assert.True(t, hasOperations)
}

func Test_GenerateEndpointSchema(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateEndpointSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)
	
	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasMethod := schema.Properties.Get("method")
	assert.True(t, hasMethod)
	_, hasPath := schema.Properties.Get("path")
	assert.True(t, hasPath)
	_, hasRequestFields := schema.Properties.Get("request_fields")
	assert.True(t, hasRequestFields)
	_, hasResponseFields := schema.Properties.Get("response_fields")
	assert.True(t, hasResponseFields)
	_, hasQueryParams := schema.Properties.Get("query_parameters")
	assert.True(t, hasQueryParams)
	_, hasPathParams := schema.Properties.Get("path_parameters")
	assert.True(t, hasPathParams)
}

func Test_GenerateAllSchemas(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schemas, err := generator.GenerateAllSchemas()
	require.NoError(t, err)
	assert.NotNil(t, schemas)
	
	// Check that all expected schemas are present
	expectedSchemas := []string{"Service", "Enum", "Object", "Resource", "Field", "ResourceField", "Endpoint"}
	for _, expectedSchema := range expectedSchemas {
		assert.Contains(t, schemas, expectedSchema, "Schema %s should be present", expectedSchema)
		assert.NotNil(t, schemas[expectedSchema], "Schema %s should not be nil", expectedSchema)
	}
	
	// Check that each schema has the correct structure
	for name, schema := range schemas {
		assert.NotEmpty(t, schema.Type, "Schema %s should have a type", name)
		assert.NotNil(t, schema.Properties, "Schema %s should have properties", name)
	}
}

func Test_SchemaToJSON(t *testing.T) {
	generator := NewSchemaGenerator()
	
	schema, err := generator.GenerateServiceSchema()
	require.NoError(t, err)
	
	jsonStr, err := generator.SchemaToJSON(schema)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)
	
	// Verify it's valid JSON
	var jsonObj map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonObj)
	require.NoError(t, err)
	
	// Check that essential schema elements are present
	assert.Contains(t, jsonObj, "type")
	assert.Contains(t, jsonObj, "properties")
}

func Test_GenerateServiceSchemaJSON(t *testing.T) {
	generator := NewSchemaGenerator()
	
	jsonStr, err := generator.GenerateServiceSchemaJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)
	
	// Verify it's valid JSON
	var jsonObj map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonObj)
	require.NoError(t, err)
	
	// Check that service-specific properties are present
	properties, ok := jsonObj["properties"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, properties, "name")
	assert.Contains(t, properties, "enums")
	assert.Contains(t, properties, "objects")
	assert.Contains(t, properties, "resources")
}

func Test_SchemaGeneration_Integration(t *testing.T) {
	// Test that we can generate schemas for a complete service specification
	generator := NewSchemaGenerator()
	
	// Generate all schemas
	schemas, err := generator.GenerateAllSchemas()
	require.NoError(t, err)
	
	// Convert each to JSON and verify they're valid
	for name, schema := range schemas {
		jsonStr, err := generator.SchemaToJSON(schema)
		require.NoError(t, err, "Failed to convert %s schema to JSON", name)
		
		// Verify it's valid JSON
		var jsonObj map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &jsonObj)
		require.NoError(t, err, "Generated JSON for %s schema is invalid", name)
		
		// Each schema should have basic structure
		assert.Contains(t, jsonObj, "type", "Schema %s should have a type", name)
		assert.Contains(t, jsonObj, "properties", "Schema %s should have properties", name)
	}
}

func Test_SchemaGeneration_WithRealData(t *testing.T) {
	// Test schema generation with actual Service data
	service := Service{
		Name: "UserAPI",
		Enums: []Enum{
			{
				Name:        "UserStatus",
				Description: "Status of the user",
				Values: []EnumValue{
					{Name: "Active", Description: "User is active"},
					{Name: "Inactive", Description: "User is inactive"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "User",
				Description: "User entity",
				Fields: []Field{
					{Name: "id", Type: "UUID", Description: "User ID"},
					{Name: "username", Type: "String", Description: "Username"},
				},
			},
		},
		Resources: []Resource{
			{
				Name:        "Users",
				Description: "User resource",
				Operations:  []string{"Create", "Read", "Update", "Delete"},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "id",
							Type:        "UUID",
							Description: "User ID",
						},
						Operations: []string{"Read"},
					},
				},
				Endpoints: []Endpoint{
					{
						Name:        "GetUser",
						Description: "Get user by ID",
						Method:      "GET",
						Path:        "/users/{id}",
					},
				},
			},
		},
	}

	// First, serialize the service to JSON to make sure our structs work
	serviceJSON, err := json.Marshal(service)
	require.NoError(t, err)
	assert.NotEmpty(t, serviceJSON)

	// Then generate the schema for Service
	generator := NewSchemaGenerator()
	schema, err := generator.GenerateServiceSchema()
	require.NoError(t, err)

	// Convert schema to JSON
	schemaJSON, err := generator.SchemaToJSON(schema)
	require.NoError(t, err)
	assert.NotEmpty(t, schemaJSON)

	// Verify that the schema contains expected field definitions
	var schemaObj map[string]interface{}
	err = json.Unmarshal([]byte(schemaJSON), &schemaObj)
	require.NoError(t, err)

	properties := schemaObj["properties"].(map[string]interface{})
	assert.Contains(t, properties, "name")
	assert.Contains(t, properties, "enums")
	assert.Contains(t, properties, "objects")
	assert.Contains(t, properties, "resources")
}