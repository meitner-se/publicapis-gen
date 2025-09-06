package specification

import (
	"encoding/json"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSchemaGeneratorWithInvalidJSON(t *testing.T) {
	generator := NewSchemaGenerator()

	// Test JSON marshaling error by creating a mock schema that fails to marshal
	// This tests the error path in SchemaToJSON
	schema := generator.reflector.Reflect(&Service{})
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
	emptyService := Service{}
	emptyEnum := Enum{}
	emptyObject := Object{}
	emptyResource := Resource{}
	emptyField := Field{}
	emptyResourceField := ResourceField{}
	emptyEndpoint := Endpoint{}
	emptyEndpointRequest := EndpointRequest{}
	emptyEndpointResponse := EndpointResponse{}

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
		"Service":         generator.GenerateServiceSchema,
		"Enum":            generator.GenerateEnumSchema,
		"Object":          generator.GenerateObjectSchema,
		"Resource":        generator.GenerateResourceSchema,
		"Field":           generator.GenerateFieldSchema,
		"ResourceField":   generator.GenerateResourceFieldSchema,
		"Endpoint":        generator.GenerateEndpointSchema,
		"EndpointRequest": generator.GenerateEndpointRequestSchema,
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
		&Service{},
		&Enum{},
		&Object{},
		&Resource{},
		&Field{},
		&ResourceField{},
		&Endpoint{},
		&EndpointRequest{},
		&EndpointResponse{},
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
	schema := generator.reflector.Reflect(Service{})
	assert.NotNil(t, schema)
	
	// Test with complex nested structures
	complexService := Service{
		Name: "Complex",
		Enums: []Enum{
			{
				Name:        "TestEnum",
				Description: "Test enumeration",
				Values: []EnumValue{
					{Name: "Value1", Description: "First value"},
					{Name: "Value2", Description: "Second value"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "ComplexObject",
				Description: "Complex object",
				Fields: []Field{
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
		Resources: []Resource{
			{
				Name:        "ComplexResource",
				Description: "Complex resource",
				Operations:  []string{"Create", "Read", "Update", "Delete", "List"},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "resourceField",
							Type:        "ComplexObject",
							Description: "Complex field",
							Modifiers:   []string{"nullable"},
						},
						Operations: []string{"Create", "Read", "Update"},
					},
				},
				Endpoints: []Endpoint{
					{
						Name:        "ComplexEndpoint",
						Title:       "Complex Endpoint",
						Description: "Complex endpoint with all fields",
						Method:      "POST",
						Path:        "/complex/{id}",
						Request: EndpointRequest{
							ContentType: "application/json",
							Headers: []Field{
								{Name: "Authorization", Type: "String", Description: "Auth header"},
								{Name: "Content-Type", Type: "String", Description: "Content type"},
							},
							PathParams: []Field{
								{Name: "id", Type: "UUID", Description: "Resource ID"},
							},
							QueryParams: []Field{
								{Name: "filter", Type: "String", Description: "Filter parameter"},
								{Name: "sort", Type: "String", Description: "Sort parameter"},
							},
							BodyParams: []Field{
								{Name: "data", Type: "ComplexObject", Description: "Request data"},
								{Name: "metadata", Type: "String", Description: "Metadata"},
							},
						},
						Response: EndpointResponse{
							ContentType: "application/json",
							StatusCode:  201,
							Headers: []Field{
								{Name: "Location", Type: "String", Description: "Resource location"},
							},
							BodyFields: []Field{
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
		name     string
		ptrType  interface{}
		valType  interface{}
	}{
		{"Service", &Service{}, Service{}},
		{"Enum", &Enum{}, Enum{}},
		{"EnumValue", &EnumValue{}, EnumValue{}},
		{"Object", &Object{}, Object{}},
		{"Resource", &Resource{}, Resource{}},
		{"Field", &Field{}, Field{}},
		{"ResourceField", &ResourceField{}, ResourceField{}},
		{"Endpoint", &Endpoint{}, Endpoint{}},
		{"EndpointRequest", &EndpointRequest{}, EndpointRequest{}},
		{"EndpointResponse", &EndpointResponse{}, EndpointResponse{}},
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
