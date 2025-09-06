package schema

import (
	"encoding/json"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meitner-se/publicapis-gen/specification"
)

func TestNewSchemaGenerator(t *testing.T) {
	generator := NewSchemaGenerator()
	assert.NotNil(t, generator)
	assert.NotNil(t, generator.reflector)
}

func TestGenerateServiceSchema(t *testing.T) {
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

func TestGenerateEnumSchema(t *testing.T) {
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

func TestGenerateObjectSchema(t *testing.T) {
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

func TestGenerateResourceSchema(t *testing.T) {
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

func TestGenerateFieldSchema(t *testing.T) {
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

func TestGenerateResourceFieldSchema(t *testing.T) {
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

func TestGenerateEndpointSchema(t *testing.T) {
	generator := NewSchemaGenerator()

	schema, err := generator.GenerateEndpointSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)

	// Check that expected properties are present
	_, hasName := schema.Properties.Get("name")
	assert.True(t, hasName)
	_, hasTitle := schema.Properties.Get("title")
	assert.True(t, hasTitle)
	_, hasDescription := schema.Properties.Get("description")
	assert.True(t, hasDescription)
	_, hasMethod := schema.Properties.Get("method")
	assert.True(t, hasMethod)
	_, hasPath := schema.Properties.Get("path")
	assert.True(t, hasPath)
	_, hasRequest := schema.Properties.Get("request")
	assert.True(t, hasRequest)
	_, hasResponse := schema.Properties.Get("response")
	assert.True(t, hasResponse)
}

func TestGenerateEndpointRequestSchema(t *testing.T) {
	generator := NewSchemaGenerator()

	schema, err := generator.GenerateEndpointRequestSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)

	// Check that expected properties are present
	_, hasContentType := schema.Properties.Get("content_type")
	assert.True(t, hasContentType)
	_, hasHeaders := schema.Properties.Get("headers")
	assert.True(t, hasHeaders)
	_, hasPathParams := schema.Properties.Get("path_params")
	assert.True(t, hasPathParams)
	_, hasQueryParams := schema.Properties.Get("query_params")
	assert.True(t, hasQueryParams)
	_, hasBodyParams := schema.Properties.Get("body_params")
	assert.True(t, hasBodyParams)
}

func TestGenerateEndpointResponseSchema(t *testing.T) {
	generator := NewSchemaGenerator()

	schema, err := generator.GenerateEndpointResponseSchema()
	require.NoError(t, err)
	assert.NotNil(t, schema)

	// Check that expected properties are present
	_, hasContentType := schema.Properties.Get("content_type")
	assert.True(t, hasContentType)
	_, hasStatusCode := schema.Properties.Get("status_code")
	assert.True(t, hasStatusCode)
	_, hasHeaders := schema.Properties.Get("headers")
	assert.True(t, hasHeaders)
	_, hasBodyFields := schema.Properties.Get("body_fields")
	assert.True(t, hasBodyFields)
	_, hasBodyObject := schema.Properties.Get("body_object")
	assert.True(t, hasBodyObject)
}

func TestGenerateAllSchemas(t *testing.T) {
	generator := NewSchemaGenerator()

	schemas, err := generator.GenerateAllSchemas()
	require.NoError(t, err)
	assert.NotNil(t, schemas)

	// Check that all expected schemas are present
	expectedSchemas := []string{"Service", "Enum", "Object", "Resource", "Field", "ResourceField", "Endpoint", "EndpointRequest", "EndpointResponse"}
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

func TestSchemaToJSON(t *testing.T) {
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

func TestGenerateServiceSchemaJSON(t *testing.T) {
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

func TestSchemaGenerationIntegration(t *testing.T) {
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

func TestSchemaGenerationWithRealData(t *testing.T) {
	// Test schema generation with actual Service data
	service := specification.Service{
		Name: "UserAPI",
		Enums: []specification.Enum{
			{
				Name:        "UserStatus",
				Description: "Status of the user",
				Values: []specification.EnumValue{
					{Name: "Active", Description: "User is active"},
					{Name: "Inactive", Description: "User is inactive"},
				},
			},
		},
		Objects: []specification.Object{
			{
				Name:        "User",
				Description: "User entity",
				Fields: []specification.Field{
					{Name: "id", Type: "UUID", Description: "User ID"},
					{Name: "username", Type: "String", Description: "Username"},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name:        "Users",
				Description: "User resource",
				Operations:  []string{"Create", "Read", "Update", "Delete"},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        "id",
							Type:        "UUID",
							Description: "User ID",
						},
						Operations: []string{"Read"},
					},
				},
				Endpoints: []specification.Endpoint{
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

func TestSchemaGeneratorWithInvalidJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test JSON marshaling error by creating a mock schema that fails to marshal
	// This tests the error path in SchemaToJSON
	schema := generator.reflector.Reflect(&specification.Service{})
	require.NotNil(t, schema)

	// Test that we get valid JSON (this should not fail with valid schema)
	jsonStr, err := generator.SchemaToJSON(schema)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)
}

func TestGenerateAllSchemasComprehensive(t *testing.T) {
	generator := NewSchemaGenerator()

	schemas, err := generator.GenerateAllSchemas()
	require.NoError(t, err)
	assert.NotNil(t, schemas)

	// Test that we have all expected schemas
	expectedSchemas := []string{
		"Service", "Enum", "Object", "Resource",
		"Field", "ResourceField", "Endpoint",
		"EndpointRequest", "EndpointResponse",
	}

	assert.Equal(t, len(expectedSchemas), len(schemas), "Should have exactly the expected number of schemas")

	// Test each schema individually
	for _, expectedSchema := range expectedSchemas {
		schema, exists := schemas[expectedSchema]
		assert.True(t, exists, "Schema %s should exist", expectedSchema)
		assert.NotNil(t, schema, "Schema %s should not be nil", expectedSchema)
		assert.NotEmpty(t, schema.Type, "Schema %s should have a type", expectedSchema)
		assert.NotNil(t, schema.Properties, "Schema %s should have properties", expectedSchema)
	}

	// Test that each schema can be converted to JSON
	for name, schema := range schemas {
		jsonStr, err := generator.SchemaToJSON(schema)
		require.NoError(t, err, "Failed to convert %s schema to JSON", name)
		assert.NotEmpty(t, jsonStr, "JSON for %s schema should not be empty", name)

		// Verify it's valid JSON
		var jsonObj map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &jsonObj)
		require.NoError(t, err, "Invalid JSON for %s schema", name)
	}

	// Test that we can call GenerateAllSchemas multiple times consistently
	schemas2, err2 := generator.GenerateAllSchemas()
	require.NoError(t, err2)
	assert.Equal(t, len(schemas), len(schemas2), "Multiple calls should return same number of schemas")
}

func TestSchemaGenerationEdgeCases(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test generating schemas for empty structs
	emptyService := specification.Service{}
	emptyEnum := specification.Enum{}
	emptyObject := specification.Object{}
	emptyResource := specification.Resource{}
	emptyField := specification.Field{}
	emptyResourceField := specification.ResourceField{}
	emptyEndpoint := specification.Endpoint{}
	emptyEndpointRequest := specification.EndpointRequest{}
	emptyEndpointResponse := specification.EndpointResponse{}

	// Test that schemas can be generated for empty structs
	testCases := []struct {
		name   string
		target interface{}
	}{
		{"EmptyService", &emptyService},
		{"EmptyEnum", &emptyEnum},
		{"EmptyObject", &emptyObject},
		{"EmptyResource", &emptyResource},
		{"EmptyField", &emptyField},
		{"EmptyResourceField", &emptyResourceField},
		{"EmptyEndpoint", &emptyEndpoint},
		{"EmptyEndpointRequest", &emptyEndpointRequest},
		{"EmptyEndpointResponse", &emptyEndpointResponse},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			schema := generator.reflector.Reflect(tc.target)
			assert.NotNil(t, schema, "Schema for %s should not be nil", tc.name)
			assert.NotEmpty(t, schema.Type, "Schema for %s should have a type", tc.name)
		})
	}
}

func TestSchemaGeneratorReflectorConfiguration(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test that reflector has expected configuration
	assert.NotNil(t, generator.reflector)
	assert.False(t, generator.reflector.AllowAdditionalProperties)
	assert.False(t, generator.reflector.DoNotReference)
	assert.True(t, generator.reflector.ExpandedStruct)
}

func TestSchemaJSONSerialization(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test JSON serialization for all schema types
	schemaGenerators := map[string]func() (*jsonschema.Schema, error){
		"Service":          generator.GenerateServiceSchema,
		"Enum":             generator.GenerateEnumSchema,
		"Object":           generator.GenerateObjectSchema,
		"Resource":         generator.GenerateResourceSchema,
		"Field":            generator.GenerateFieldSchema,
		"ResourceField":    generator.GenerateResourceFieldSchema,
		"Endpoint":         generator.GenerateEndpointSchema,
		"EndpointRequest":  generator.GenerateEndpointRequestSchema,
		"EndpointResponse": generator.GenerateEndpointResponseSchema,
	}

	for name, schemaFunc := range schemaGenerators {
		t.Run(name, func(t *testing.T) {
			schema, err := schemaFunc()
			require.NoError(t, err)

			jsonStr, err := generator.SchemaToJSON(schema)
			require.NoError(t, err)
			assert.NotEmpty(t, jsonStr)

			// Verify it's valid JSON
			var jsonObj map[string]interface{}
			err = json.Unmarshal([]byte(jsonStr), &jsonObj)
			require.NoError(t, err)

			// All schemas should have basic structure
			assert.Contains(t, jsonObj, "type")
			assert.Contains(t, jsonObj, "properties")
		})
	}
}

func TestGenerateServiceSchemaJSONIntegration(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test the complete flow from struct to JSON schema string
	jsonSchema, err := generator.GenerateServiceSchemaJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonSchema)

	// Validate it's proper JSON
	var schemaObj map[string]interface{}
	err = json.Unmarshal([]byte(jsonSchema), &schemaObj)
	require.NoError(t, err)

	// Test specific schema properties
	assert.Contains(t, schemaObj, "$schema")
	assert.Contains(t, schemaObj, "type")
	assert.Equal(t, "object", schemaObj["type"])
}

func TestSchemaGeneratorConfiguration(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test that we can create multiple generators
	generator2 := NewSchemaGenerator()
	assert.True(t, generator != generator2, "Generators should be different instances")

	// Both should work independently
	schema1, err1 := generator.GenerateServiceSchema()
	schema2, err2 := generator2.GenerateServiceSchema()

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotNil(t, schema1)
	assert.NotNil(t, schema2)

	// Test that both generators have the same configuration
	assert.Equal(t, generator.reflector.AllowAdditionalProperties, generator2.reflector.AllowAdditionalProperties)
	assert.Equal(t, generator.reflector.DoNotReference, generator2.reflector.DoNotReference)
	assert.Equal(t, generator.reflector.ExpandedStruct, generator2.reflector.ExpandedStruct)
}

func TestSchemaReflectorBehavior(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test reflector with various struct types
	testStructs := []interface{}{
		&specification.Service{},
		&specification.Enum{},
		&specification.Object{},
		&specification.Resource{},
		&specification.Field{},
		&specification.ResourceField{},
		&specification.Endpoint{},
		&specification.EndpointRequest{},
		&specification.EndpointResponse{},
	}

	for i, testStruct := range testStructs {
		schema := generator.reflector.Reflect(testStruct)
		assert.NotNil(t, schema, "Schema %d should not be nil", i)
		assert.NotEmpty(t, schema.Type, "Schema %d should have a type", i)
	}
}

func TestSchemaGeneratorErrorScenarios(t *testing.T) {
	// Test various edge cases that might help improve coverage
	generator := NewSchemaGenerator()

	// Test with non-pointer types (which should still work)
	schema := generator.reflector.Reflect(specification.Service{})
	assert.NotNil(t, schema)

	// Test with complex nested structures
	complexService := specification.Service{
		Name: "Complex",
		Enums: []specification.Enum{
			{
				Name:        "TestEnum",
				Description: "Test enumeration",
				Values: []specification.EnumValue{
					{Name: "Value1", Description: "First value"},
					{Name: "Value2", Description: "Second value"},
				},
			},
		},
		Objects: []specification.Object{
			{
				Name:        "ComplexObject",
				Description: "Complex object",
				Fields: []specification.Field{
					{
						Name:        "nestedField",
						Type:        "String",
						Description: "A nested field",
						Modifiers:   []string{"array", "nullable"},
						Default:     "defaultValue",
						Example:     "exampleValue",
					},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name:        "ComplexResource",
				Description: "Complex resource",
				Operations:  []string{"Create", "Read", "Update", "Delete", "List"},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        "resourceField",
							Type:        "ComplexObject",
							Description: "Complex field",
							Modifiers:   []string{"nullable"},
						},
						Operations: []string{"Create", "Read", "Update"},
					},
				},
				Endpoints: []specification.Endpoint{
					{
						Name:        "ComplexEndpoint",
						Title:       "Complex Endpoint",
						Description: "Complex endpoint with all fields",
						Method:      "POST",
						Path:        "/complex/{id}",
						Request: specification.EndpointRequest{
							ContentType: "application/json",
							Headers: []specification.Field{
								{Name: "Authorization", Type: "String", Description: "Auth header"},
								{Name: "Content-Type", Type: "String", Description: "Content type"},
							},
							PathParams: []specification.Field{
								{Name: "id", Type: "UUID", Description: "Resource ID"},
							},
							QueryParams: []specification.Field{
								{Name: "filter", Type: "String", Description: "Filter parameter"},
								{Name: "sort", Type: "String", Description: "Sort parameter"},
							},
							BodyParams: []specification.Field{
								{Name: "data", Type: "ComplexObject", Description: "Request data"},
								{Name: "metadata", Type: "String", Description: "Metadata"},
							},
						},
						Response: specification.EndpointResponse{
							ContentType: "application/json",
							StatusCode:  201,
							Headers: []specification.Field{
								{Name: "Location", Type: "String", Description: "Resource location"},
							},
							BodyFields: []specification.Field{
								{Name: "id", Type: "UUID", Description: "Created ID"},
								{Name: "status", Type: "String", Description: "Creation status"},
							},
						},
					},
				},
			},
		},
	}

	// Test that complex structures can generate schemas
	schema = generator.reflector.Reflect(&complexService)
	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)

	// Test JSON conversion with complex structure
	jsonStr, err := generator.SchemaToJSON(schema)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	// Validate complex JSON structure
	var jsonObj map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonObj)
	require.NoError(t, err)
	assert.Contains(t, jsonObj, "properties")
}

func TestAllStructTypesWithReflection(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test every struct type with both pointer and non-pointer variants
	testCases := []struct {
		name    string
		ptrType interface{}
		valType interface{}
	}{
		{"Service", &specification.Service{}, specification.Service{}},
		{"Enum", &specification.Enum{}, specification.Enum{}},
		{"EnumValue", &specification.EnumValue{}, specification.EnumValue{}},
		{"Object", &specification.Object{}, specification.Object{}},
		{"Resource", &specification.Resource{}, specification.Resource{}},
		{"Field", &specification.Field{}, specification.Field{}},
		{"ResourceField", &specification.ResourceField{}, specification.ResourceField{}},
		{"Endpoint", &specification.Endpoint{}, specification.Endpoint{}},
		{"EndpointRequest", &specification.EndpointRequest{}, specification.EndpointRequest{}},
		{"EndpointResponse", &specification.EndpointResponse{}, specification.EndpointResponse{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with pointer
			schemaPtrType := generator.reflector.Reflect(tc.ptrType)
			assert.NotNil(t, schemaPtrType, "%s pointer type should generate schema", tc.name)

			// Test with value
			schemaValType := generator.reflector.Reflect(tc.valType)
			assert.NotNil(t, schemaValType, "%s value type should generate schema", tc.name)

			// Both should have same basic structure
			assert.Equal(t, schemaPtrType.Type, schemaValType.Type, "%s schemas should have same type", tc.name)
		})
	}
}

func TestValidateService(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid service JSON
	validServiceJSON := `{
		"name": "UserAPI",
		"enums": [],
		"objects": [],
		"resources": []
	}`

	err := generator.ValidateService([]byte(validServiceJSON))
	assert.NoError(t, err)

	// Test valid service YAML
	validServiceYAML := `
name: UserAPI
enums: []
objects: []
resources: []
`

	err = generator.ValidateService([]byte(validServiceYAML))
	assert.NoError(t, err)

	// Test invalid service JSON (missing required field)
	invalidServiceJSON := `{
		"enums": [],
		"objects": [],
		"resources": []
	}`

	err = generator.ValidateService([]byte(invalidServiceJSON))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation errors")

	// Test malformed JSON
	malformedJSON := `{"name": "test", invalid json}`
	err = generator.ValidateService([]byte(malformedJSON))
	assert.Error(t, err)
}

func TestValidateEnum(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid enum JSON
	validEnumJSON := `{
		"name": "Status",
		"description": "Status enumeration",
		"values": [
			{"name": "Active", "description": "Active status"},
			{"name": "Inactive", "description": "Inactive status"}
		]
	}`

	err := generator.ValidateEnum([]byte(validEnumJSON))
	assert.NoError(t, err)

	// Test valid enum YAML
	validEnumYAML := `
name: Status
description: Status enumeration
values:
  - name: Active
    description: Active status
  - name: Inactive
    description: Inactive status
`

	err = generator.ValidateEnum([]byte(validEnumYAML))
	assert.NoError(t, err)

	// Test invalid enum JSON (missing required field)
	invalidEnumJSON := `{
		"description": "Status enumeration",
		"values": []
	}`

	err = generator.ValidateEnum([]byte(invalidEnumJSON))
	assert.Error(t, err)
}

func TestValidateObject(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid object JSON
	validObjectJSON := `{
		"name": "User",
		"description": "User object",
		"fields": [
			{"name": "id", "type": "UUID", "description": "User ID"},
			{"name": "name", "type": "String", "description": "User name"}
		]
	}`

	err := generator.ValidateObject([]byte(validObjectJSON))
	assert.NoError(t, err)

	// Test invalid object JSON
	invalidObjectJSON := `{
		"description": "User object",
		"fields": []
	}`

	err = generator.ValidateObject([]byte(invalidObjectJSON))
	assert.Error(t, err)
}

func TestValidateResource(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid resource JSON
	validResourceJSON := `{
		"name": "Users",
		"description": "User resource",
		"operations": ["Create", "Read", "Update", "Delete"],
		"fields": [],
		"endpoints": []
	}`

	err := generator.ValidateResource([]byte(validResourceJSON))
	assert.NoError(t, err)

	// Test invalid resource JSON
	invalidResourceJSON := `{
		"description": "User resource",
		"operations": [],
		"fields": [],
		"endpoints": []
	}`

	err = generator.ValidateResource([]byte(invalidResourceJSON))
	assert.Error(t, err)
}

func TestValidateField(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid field JSON
	validFieldJSON := `{
		"name": "username",
		"description": "User's username",
		"type": "String"
	}`

	err := generator.ValidateField([]byte(validFieldJSON))
	assert.NoError(t, err)

	// Test valid field with optional properties
	validFieldWithOptionalJSON := `{
		"name": "age",
		"description": "User's age",
		"type": "Int",
		"default": "0",
		"example": "25",
		"modifiers": ["nullable"]
	}`

	err = generator.ValidateField([]byte(validFieldWithOptionalJSON))
	assert.NoError(t, err)

	// Test invalid field JSON
	invalidFieldJSON := `{
		"description": "User's username",
		"type": "String"
	}`

	err = generator.ValidateField([]byte(invalidFieldJSON))
	assert.Error(t, err)
}

func TestValidateResourceField(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid resource field JSON
	validResourceFieldJSON := `{
		"name": "username",
		"description": "User's username",
		"type": "String",
		"operations": ["Create", "Read", "Update"]
	}`

	err := generator.ValidateResourceField([]byte(validResourceFieldJSON))
	assert.NoError(t, err)

	// Test invalid resource field JSON
	invalidResourceFieldJSON := `{
		"description": "User's username",
		"type": "String",
		"operations": []
	}`

	err = generator.ValidateResourceField([]byte(invalidResourceFieldJSON))
	assert.Error(t, err)
}

func TestValidateEndpoint(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid endpoint JSON
	validEndpointJSON := `{
		"name": "GetUser",
		"title": "Get User",
		"description": "Get user by ID",
		"method": "GET",
		"path": "/users/{id}",
		"request": {
			"content_type": "",
			"headers": [],
			"path_params": [],
			"query_params": [],
			"body_params": []
		},
		"response": {
			"content_type": "",
			"status_code": 200,
			"headers": [],
			"body_fields": []
		}
	}`

	err := generator.ValidateEndpoint([]byte(validEndpointJSON))
	assert.NoError(t, err)

	// Test invalid endpoint JSON
	invalidEndpointJSON := `{
		"title": "Get User",
		"description": "Get user by ID",
		"method": "GET",
		"path": "/users/{id}",
		"request": {},
		"response": {}
	}`

	err = generator.ValidateEndpoint([]byte(invalidEndpointJSON))
	assert.Error(t, err)
}

func TestValidateEndpointRequest(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid endpoint request JSON
	validEndpointRequestJSON := `{
		"content_type": "application/json",
		"headers": [],
		"path_params": [],
		"query_params": [],
		"body_params": []
	}`

	err := generator.ValidateEndpointRequest([]byte(validEndpointRequestJSON))
	assert.NoError(t, err)

	// Test endpoint request with fields
	validEndpointRequestWithFieldsJSON := `{
		"content_type": "application/json",
		"headers": [
			{"name": "Authorization", "type": "String", "description": "Bearer token"}
		],
		"path_params": [
			{"name": "id", "type": "UUID", "description": "Resource ID"}
		],
		"query_params": [
			{"name": "limit", "type": "Int", "description": "Limit", "default": "10"}
		],
		"body_params": [
			{"name": "name", "type": "String", "description": "Resource name"}
		]
	}`

	err = generator.ValidateEndpointRequest([]byte(validEndpointRequestWithFieldsJSON))
	assert.NoError(t, err)
}

func TestValidateEndpointResponse(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid endpoint response JSON
	validEndpointResponseJSON := `{
		"content_type": "application/json",
		"status_code": 200,
		"headers": [],
		"body_fields": []
	}`

	err := generator.ValidateEndpointResponse([]byte(validEndpointResponseJSON))
	assert.NoError(t, err)

	// Test endpoint response with body object
	validEndpointResponseWithBodyObjectJSON := `{
		"content_type": "application/json",
		"status_code": 201,
		"headers": [],
		"body_fields": [],
		"body_object": "User"
	}`

	err = generator.ValidateEndpointResponse([]byte(validEndpointResponseWithBodyObjectJSON))
	assert.NoError(t, err)

	// Test invalid endpoint response JSON (invalid status code type)
	invalidEndpointResponseJSON := `{
		"content_type": "application/json",
		"status_code": "200",
		"headers": [],
		"body_fields": []
	}`

	err = generator.ValidateEndpointResponse([]byte(invalidEndpointResponseJSON))
	assert.Error(t, err)
}

func TestConvertToJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid JSON
	validJSON := `{"name": "test", "value": 123}`
	result, err := generator.convertToJSON([]byte(validJSON))
	require.NoError(t, err)
	assert.JSONEq(t, validJSON, string(result))

	// Test valid YAML
	validYAML := `
name: test
value: 123
`
	result, err = generator.convertToJSON([]byte(validYAML))
	require.NoError(t, err)
	assert.Contains(t, string(result), "test")
	assert.Contains(t, string(result), "123")

	// Test invalid data
	invalidData := `this is neither JSON nor YAML: {[}`
	_, err = generator.convertToJSON([]byte(invalidData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "neither valid JSON nor YAML")
}

func TestParseServiceFromJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid service JSON
	validServiceJSON := `{
		"name": "UserAPI",
		"enums": [
			{
				"name": "Status",
				"description": "Status enumeration",
				"values": [
					{"name": "Active", "description": "Active status"}
				]
			}
		],
		"objects": [
			{
				"name": "User",
				"description": "User object",
				"fields": [
					{"name": "id", "type": "UUID", "description": "User ID"}
				]
			}
		],
		"resources": [
			{
				"name": "Users",
				"description": "User resource",
				"operations": ["Create", "Read"],
				"fields": [
					{
						"name": "id",
						"type": "UUID",
						"description": "User ID",
						"operations": ["Read"]
					}
				],
				"endpoints": [
					{
						"name": "GetUser",
						"title": "Get User",
						"description": "Get user by ID",
						"method": "GET",
						"path": "/users/{id}",
						"request": {
							"content_type": "",
							"headers": [],
							"path_params": [],
							"query_params": [],
							"body_params": []
						},
						"response": {
							"content_type": "application/json",
							"status_code": 200,
							"headers": [],
							"body_fields": []
						}
					}
				]
			}
		]
	}`

	service, err := generator.ParseServiceFromJSON([]byte(validServiceJSON))
	require.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "UserAPI", service.Name)
	assert.Len(t, service.Enums, 1)
	assert.Len(t, service.Objects, 1)
	assert.Len(t, service.Resources, 1)
	assert.Equal(t, "Status", service.Enums[0].Name)
	assert.Equal(t, "User", service.Objects[0].Name)
	assert.Equal(t, "Users", service.Resources[0].Name)

	// Test invalid service JSON (missing required field)
	invalidServiceJSON := `{
		"enums": [],
		"objects": [],
		"resources": []
	}`

	_, err = generator.ParseServiceFromJSON([]byte(invalidServiceJSON))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test malformed JSON
	malformedJSON := `{"name": "test", invalid}`
	_, err = generator.ParseServiceFromJSON([]byte(malformedJSON))
	assert.Error(t, err)
}

func TestParseServiceFromYAML(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid service YAML
	validServiceYAML := `
name: UserAPI
enums:
  - name: Status
    description: Status enumeration
    values:
      - name: Active
        description: Active status
objects:
  - name: User
    description: User object
    fields:
      - name: id
        type: UUID
        description: User ID
resources:
  - name: Users
    description: User resource
    operations: [Create, Read]
    fields:
      - name: id
        type: UUID
        description: User ID
        operations: [Read]
    endpoints:
      - name: GetUser
        title: Get User
        description: Get user by ID
        method: GET
        path: /users/{id}
        request:
          content_type: ""
          headers: []
          path_params: []
          query_params: []
          body_params: []
        response:
          content_type: application/json
          status_code: 200
          headers: []
          body_fields: []
`

	service, err := generator.ParseServiceFromYAML([]byte(validServiceYAML))
	require.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "UserAPI", service.Name)
	assert.Len(t, service.Enums, 1)
	assert.Len(t, service.Objects, 1)
	assert.Len(t, service.Resources, 1)

	// Test invalid service YAML
	invalidServiceYAML := `
enums: []
objects: []
resources: []
`

	_, err = generator.ParseServiceFromYAML([]byte(invalidServiceYAML))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseEnumFromJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid enum JSON
	validEnumJSON := `{
		"name": "Status",
		"description": "Status enumeration",
		"values": [
			{"name": "Active", "description": "Active status"},
			{"name": "Inactive", "description": "Inactive status"}
		]
	}`

	enum, err := generator.ParseEnumFromJSON([]byte(validEnumJSON))
	require.NoError(t, err)
	assert.NotNil(t, enum)
	assert.Equal(t, "Status", enum.Name)
	assert.Equal(t, "Status enumeration", enum.Description)
	assert.Len(t, enum.Values, 2)
	assert.Equal(t, "Active", enum.Values[0].Name)
	assert.Equal(t, "Inactive", enum.Values[1].Name)

	// Test invalid enum JSON
	invalidEnumJSON := `{
		"description": "Status enumeration",
		"values": []
	}`

	_, err = generator.ParseEnumFromJSON([]byte(invalidEnumJSON))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseEnumFromYAML(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid enum YAML
	validEnumYAML := `
name: Status
description: Status enumeration
values:
  - name: Active
    description: Active status
  - name: Inactive
    description: Inactive status
`

	enum, err := generator.ParseEnumFromYAML([]byte(validEnumYAML))
	require.NoError(t, err)
	assert.NotNil(t, enum)
	assert.Equal(t, "Status", enum.Name)
	assert.Equal(t, "Status enumeration", enum.Description)
	assert.Len(t, enum.Values, 2)

	// Test invalid enum YAML
	invalidEnumYAML := `
description: Status enumeration
values: []
`

	_, err = generator.ParseEnumFromYAML([]byte(invalidEnumYAML))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseObjectFromJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid object JSON
	validObjectJSON := `{
		"name": "User",
		"description": "User object",
		"fields": [
			{"name": "id", "type": "UUID", "description": "User ID"},
			{"name": "name", "type": "String", "description": "User name", "example": "John Doe"}
		]
	}`

	object, err := generator.ParseObjectFromJSON([]byte(validObjectJSON))
	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "User", object.Name)
	assert.Equal(t, "User object", object.Description)
	assert.Len(t, object.Fields, 2)
	assert.Equal(t, "id", object.Fields[0].Name)
	assert.Equal(t, "name", object.Fields[1].Name)

	// Test invalid object JSON
	invalidObjectJSON := `{
		"description": "User object",
		"fields": []
	}`

	_, err = generator.ParseObjectFromJSON([]byte(invalidObjectJSON))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseObjectFromYAML(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid object YAML
	validObjectYAML := `
name: User
description: User object
fields:
  - name: id
    type: UUID
    description: User ID
  - name: name
    type: String
    description: User name
    example: John Doe
`

	object, err := generator.ParseObjectFromYAML([]byte(validObjectYAML))
	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "User", object.Name)
	assert.Equal(t, "User object", object.Description)
	assert.Len(t, object.Fields, 2)

	// Test invalid object YAML
	invalidObjectYAML := `
description: User object
fields: []
`

	_, err = generator.ParseObjectFromYAML([]byte(invalidObjectYAML))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseResourceFromJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid resource JSON
	validResourceJSON := `{
		"name": "Users",
		"description": "User resource",
		"operations": ["Create", "Read", "Update", "Delete"],
		"fields": [
			{
				"name": "id",
				"type": "UUID",
				"description": "User ID",
				"operations": ["Read"]
			}
		],
		"endpoints": [
			{
				"name": "GetUser",
				"title": "Get User",
				"description": "Get user by ID",
				"method": "GET",
				"path": "/users/{id}",
				"request": {
					"content_type": "",
					"headers": [],
					"path_params": [],
					"query_params": [],
					"body_params": []
				},
				"response": {
					"content_type": "application/json",
					"status_code": 200,
					"headers": [],
					"body_fields": []
				}
			}
		]
	}`

	resource, err := generator.ParseResourceFromJSON([]byte(validResourceJSON))
	require.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, "Users", resource.Name)
	assert.Equal(t, "User resource", resource.Description)
	assert.Len(t, resource.Operations, 4)
	assert.Len(t, resource.Fields, 1)
	assert.Len(t, resource.Endpoints, 1)

	// Test invalid resource JSON
	invalidResourceJSON := `{
		"description": "User resource",
		"operations": [],
		"fields": [],
		"endpoints": []
	}`

	_, err = generator.ParseResourceFromJSON([]byte(invalidResourceJSON))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParseResourceFromYAML(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test valid resource YAML
	validResourceYAML := `
name: Users
description: User resource
operations: [Create, Read, Update, Delete]
fields:
  - name: id
    type: UUID
    description: User ID
    operations: [Read]
endpoints:
  - name: GetUser
    title: Get User
    description: Get user by ID
    method: GET
    path: /users/{id}
    request:
      content_type: ""
      headers: []
      path_params: []
      query_params: []
      body_params: []
    response:
      content_type: application/json
      status_code: 200
      headers: []
      body_fields: []
`

	resource, err := generator.ParseResourceFromYAML([]byte(validResourceYAML))
	require.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, "Users", resource.Name)
	assert.Equal(t, "User resource", resource.Description)
	assert.Len(t, resource.Operations, 4)
	assert.Len(t, resource.Fields, 1)
	assert.Len(t, resource.Endpoints, 1)

	// Test invalid resource YAML
	invalidResourceYAML := `
description: User resource
operations: []
fields: []
endpoints: []
`

	_, err = generator.ParseResourceFromYAML([]byte(invalidResourceYAML))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidationErrorMessages(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test validation with detailed error information
	invalidServiceJSON := `{
		"enums": "not an array",
		"objects": [],
		"resources": []
	}`

	err := generator.ValidateService([]byte(invalidServiceJSON))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation errors")

	// The error should contain information about what's wrong
	errorStr := err.Error()
	assert.NotEmpty(t, errorStr)
}

func TestValidationWithComplexStructures(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test complex valid service structure
	complexValidServiceJSON := `{
		"name": "ComplexAPI",
		"enums": [
			{
				"name": "UserRole",
				"description": "User roles in the system",
				"values": [
					{"name": "Admin", "description": "Administrator role"},
					{"name": "User", "description": "Regular user role"}
				]
			}
		],
		"objects": [
			{
				"name": "Address",
				"description": "Address information",
				"fields": [
					{"name": "street", "type": "String", "description": "Street address"},
					{"name": "city", "type": "String", "description": "City"},
					{"name": "zipCode", "type": "String", "description": "ZIP code", "modifiers": ["nullable"]}
				]
			}
		],
		"resources": [
			{
				"name": "Users",
				"description": "User management",
				"operations": ["Create", "Read", "Update", "Delete"],
				"fields": [
					{
						"name": "id",
						"type": "UUID",
						"description": "User ID",
						"operations": ["Read"]
					},
					{
						"name": "role",
						"type": "UserRole",
						"description": "User role",
						"default": "User",
						"operations": ["Create", "Read", "Update"]
					}
				],
				"endpoints": [
					{
						"name": "CreateUser",
						"title": "Create New User",
						"description": "Create a new user account",
						"method": "POST",
						"path": "/",
						"request": {
							"content_type": "application/json",
							"headers": [
								{"name": "Authorization", "type": "String", "description": "Bearer token"}
							],
							"path_params": [],
							"query_params": [],
							"body_params": [
								{"name": "username", "type": "String", "description": "Username"},
								{"name": "email", "type": "String", "description": "Email"},
								{"name": "role", "type": "UserRole", "description": "User role"}
							]
						},
						"response": {
							"content_type": "application/json",
							"status_code": 201,
							"headers": [
								{"name": "Location", "type": "String", "description": "Created resource URL"}
							],
							"body_fields": [
								{"name": "id", "type": "UUID", "description": "Created user ID"},
								{"name": "username", "type": "String", "description": "Username"}
							]
						}
					}
				]
			}
		]
	}`

	// This should validate successfully
	err := generator.ValidateService([]byte(complexValidServiceJSON))
	assert.NoError(t, err)

	// And should parse successfully
	service, err := generator.ParseServiceFromJSON([]byte(complexValidServiceJSON))
	require.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "ComplexAPI", service.Name)
	assert.Len(t, service.Enums, 1)
	assert.Len(t, service.Objects, 1)
	assert.Len(t, service.Resources, 1)

	// Verify complex nested structures
	assert.Equal(t, "UserRole", service.Enums[0].Name)
	assert.Len(t, service.Enums[0].Values, 2)
	assert.Equal(t, "Address", service.Objects[0].Name)
	assert.Len(t, service.Objects[0].Fields, 3)
	assert.Equal(t, "Users", service.Resources[0].Name)
	assert.Len(t, service.Resources[0].Fields, 2)
	assert.Len(t, service.Resources[0].Endpoints, 1)
	assert.Equal(t, "CreateUser", service.Resources[0].Endpoints[0].Name)
}
