package openapi

import (
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
)

const (
	expectedErrorNotImplemented  = "not implemented"
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

// TestGenerateFromServiceWithValidService tests stub behavior with valid service.
func TestGenerateFromServiceWithValidService(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "TestService",
	}

	document, err := generator.GenerateFromService(service)

	assert.Nil(t, document, "Document should be nil in stub implementation")
	assert.EqualError(t, err, expectedErrorNotImplemented, "Should return not implemented error")
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
