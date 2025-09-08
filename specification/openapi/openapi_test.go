package openapi

import (
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
)

const (
	expectedErrorInvalidService  = "invalid service: service cannot be nil"
	expectedErrorInvalidDocument = "invalid document: document cannot be nil"
)

// TestNewGenerator tests the creation of a new OpenAPI generator.
func TestNewGenerator(t *testing.T) {
	expectedVersion := "3.1.0"

	generator := NewGenerator()

	assert.NotNil(t, generator, "Generator should not be nil")
	assert.Equal(t, expectedVersion, generator.Version, "Generator version should be 3.1.0")
	assert.Equal(t, "", generator.Title, "Generator title should be empty by default")
	assert.Equal(t, "", generator.Description, "Generator description should be empty by default")
	assert.Equal(t, "", generator.ServerURL, "Generator server URL should be empty by default")
}

// TestGenerateFromServiceWithNilService tests error handling when service is nil.
func TestGenerateFromServiceWithNilService(t *testing.T) {
	generator := NewGenerator()

	document, err := generator.GenerateFromService(nil)

	assert.Nil(t, document, "Document should be nil when service is nil")
	assert.EqualError(t, err, expectedErrorInvalidService, "Should return invalid service error")
}

// TestGenerateFromServiceWithValidService tests generating OpenAPI document from valid service.
func TestGenerateFromServiceWithValidService(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "TestService",
	}

	document, err := generator.GenerateFromService(service)

	assert.NotNil(t, document, "Document should not be nil with valid service")
	assert.NoError(t, err, "Should not return error with valid service")
	assert.Equal(t, "3.1.0", document.Version, "Document version should be 3.1.0")
	assert.Equal(t, "TestService", document.Info.Title, "Document title should match service name")
}

// TestToYAMLWithNilDocument tests error handling when document is nil.
func TestToYAMLWithNilDocument(t *testing.T) {
	generator := NewGenerator()

	yamlBytes, err := generator.ToYAML(nil)

	assert.Nil(t, yamlBytes, "YAML bytes should be nil when document is nil")
	assert.EqualError(t, err, expectedErrorInvalidDocument, "Should return invalid document error")
}

// TestToJSONWithNilDocument tests error handling when document is nil.
func TestToJSONWithNilDocument(t *testing.T) {
	generator := NewGenerator()

	jsonBytes, err := generator.ToJSON(nil)

	assert.Nil(t, jsonBytes, "JSON bytes should be nil when document is nil")
	assert.EqualError(t, err, expectedErrorInvalidDocument, "Should return invalid document error")
}

// TestGeneratorConfiguration tests generator configuration options.
func TestGeneratorConfiguration(t *testing.T) {
	expectedTitle := "Custom API"
	expectedDescription := "Custom API Description"
	expectedServerURL := "https://custom.example.com"

	generator := &Generator{
		Version:     "3.1.0",
		Title:       expectedTitle,
		Description: expectedDescription,
		ServerURL:   expectedServerURL,
	}

	assert.Equal(t, expectedTitle, generator.Title, "Generator title should match configured value")
	assert.Equal(t, expectedDescription, generator.Description, "Generator description should match configured value")
	assert.Equal(t, expectedServerURL, generator.ServerURL, "Generator server URL should match configured value")
}

// TestToJSONWithValidDocument tests JSON conversion with valid document.
func TestToJSONWithValidDocument(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "TestService",
	}

	document, err := generator.GenerateFromService(service)
	assert.NoError(t, err, "Should generate document successfully")

	jsonBytes, err := generator.ToJSON(document)

	assert.NoError(t, err, "Should convert document to JSON successfully")
	assert.NotNil(t, jsonBytes, "JSON bytes should not be nil")
	assert.Contains(t, string(jsonBytes), "TestService", "JSON should contain service name")
	assert.Contains(t, string(jsonBytes), "3.1.0", "JSON should contain OpenAPI version")
}

// TestToYAMLWithValidDocument tests YAML conversion with valid document.
func TestToYAMLWithValidDocument(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "TestService",
	}

	document, err := generator.GenerateFromService(service)
	assert.NoError(t, err, "Should generate document successfully")

	yamlBytes, err := generator.ToYAML(document)

	assert.NoError(t, err, "Should convert document to YAML successfully")
	assert.NotNil(t, yamlBytes, "YAML bytes should not be nil")
	assert.Contains(t, string(yamlBytes), "TestService", "YAML should contain service name")
	assert.Contains(t, string(yamlBytes), "3.1.0", "YAML should contain OpenAPI version")
}

// TestGenerateFromServiceWithComplexService tests object generation with enums, objects, and resources.
func TestGenerateFromServiceWithComplexService(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "UserAPI",
		Enums: []specification.Enum{
			{
				Name:        "Status",
				Description: "User status enumeration",
				Values: []specification.EnumValue{
					{Name: "Active", Description: "User is active"},
					{Name: "Inactive", Description: "User is inactive"},
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
						Description: "User identifier",
						Type:        specification.FieldTypeUUID,
					},
					{
						Name:        "email",
						Description: "User email address",
						Type:        specification.FieldTypeString,
					},
					{
						Name:        "status",
						Description: "User status",
						Type:        "Status",
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
							Description: "User identifier",
							Type:        specification.FieldTypeUUID,
						},
						Operations: []string{specification.OperationRead},
					},
					{
						Field: specification.Field{
							Name:        "email",
							Description: "User email address",
							Type:        specification.FieldTypeString,
						},
						Operations: []string{specification.OperationCreate, specification.OperationRead},
					},
				},
				Endpoints: []specification.Endpoint{
					{
						Name:        "Create",
						Title:       "Create User",
						Description: "Create a new user",
						Method:      "POST",
						Path:        "",
						Request: specification.EndpointRequest{
							ContentType: "application/json",
							BodyParams: []specification.Field{
								{
									Name:        "email",
									Description: "User email address",
									Type:        specification.FieldTypeString,
								},
							},
						},
						Response: specification.EndpointResponse{
							ContentType: "application/json",
							StatusCode:  201,
							BodyObject:  stringPtr("User"),
						},
					},
				},
			},
		},
	}

	document, err := generator.GenerateFromService(service)

	assert.NoError(t, err, "Should generate document successfully")
	assert.NotNil(t, document, "Document should not be nil")
	assert.Equal(t, "UserAPI", document.Info.Title, "Document title should match service name")

	// Test JSON output contains expected elements
	jsonBytes, err := generator.ToJSON(document)
	assert.NoError(t, err, "Should convert document to JSON successfully")
	jsonString := string(jsonBytes)
	assert.Contains(t, jsonString, "Status", "JSON should contain Status enum")
	assert.Contains(t, jsonString, "User", "JSON should contain User object")
	assert.Contains(t, jsonString, "Active", "JSON should contain enum values")
	assert.Contains(t, jsonString, "/user", "JSON should contain user path")
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
