package openapi

import (
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
)

const (
	expectedErrorNotImplemented = "not implemented"
	expectedErrorInvalidService = "invalid service: service cannot be nil"
	expectedErrorInvalidSpec    = "invalid specification: spec cannot be nil"
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

	spec, err := generator.GenerateFromService(nil)

	assert.Nil(t, spec, "Specification should be nil when service is nil")
	assert.EqualError(t, err, expectedErrorInvalidService, "Should return invalid service error")
}

// TestGenerateFromServiceWithValidService tests stub behavior with valid service.
func TestGenerateFromServiceWithValidService(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: "TestService",
	}

	spec, err := generator.GenerateFromService(service)

	assert.Nil(t, spec, "Specification should be nil in stub implementation")
	assert.EqualError(t, err, expectedErrorNotImplemented, "Should return not implemented error")
}

// TestToYAMLWithNilSpec tests error handling when specification is nil.
func TestToYAMLWithNilSpec(t *testing.T) {
	generator := NewGenerator()

	yamlBytes, err := generator.ToYAML(nil)

	assert.Nil(t, yamlBytes, "YAML bytes should be nil when spec is nil")
	assert.EqualError(t, err, expectedErrorInvalidSpec, "Should return invalid spec error")
}

// TestToYAMLWithValidSpec tests stub behavior with valid specification.
func TestToYAMLWithValidSpec(t *testing.T) {
	generator := NewGenerator()
	spec := &Specification{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	yamlBytes, err := generator.ToYAML(spec)

	assert.Nil(t, yamlBytes, "YAML bytes should be nil in stub implementation")
	assert.EqualError(t, err, expectedErrorNotImplemented, "Should return not implemented error")
}

// TestToJSONWithNilSpec tests error handling when specification is nil.
func TestToJSONWithNilSpec(t *testing.T) {
	generator := NewGenerator()

	jsonBytes, err := generator.ToJSON(nil)

	assert.Nil(t, jsonBytes, "JSON bytes should be nil when spec is nil")
	assert.EqualError(t, err, expectedErrorInvalidSpec, "Should return invalid spec error")
}

// TestToJSONWithValidSpec tests stub behavior with valid specification.
func TestToJSONWithValidSpec(t *testing.T) {
	generator := NewGenerator()
	spec := &Specification{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}

	jsonBytes, err := generator.ToJSON(spec)

	assert.Nil(t, jsonBytes, "JSON bytes should be nil in stub implementation")
	assert.EqualError(t, err, expectedErrorNotImplemented, "Should return not implemented error")
}

// TestSpecificationStructure tests the OpenAPI specification structure.
func TestSpecificationStructure(t *testing.T) {
	expectedOpenAPI := "3.1.0"
	expectedTitle := "Test API"
	expectedVersion := "1.0.0"
	expectedDescription := "Test API Description"
	expectedServerURL := "https://api.example.com"
	expectedServerDescription := "Production server"

	spec := &Specification{
		OpenAPI: expectedOpenAPI,
		Info: Info{
			Title:       expectedTitle,
			Version:     expectedVersion,
			Description: expectedDescription,
		},
		Servers: []Server{
			{
				URL:         expectedServerURL,
				Description: expectedServerDescription,
			},
		},
		Paths:      make(map[string]PathItem),
		Components: &Components{Schemas: make(map[string]Schema)},
	}

	assert.Equal(t, expectedOpenAPI, spec.OpenAPI, "OpenAPI version should match")
	assert.Equal(t, expectedTitle, spec.Info.Title, "Info title should match")
	assert.Equal(t, expectedVersion, spec.Info.Version, "Info version should match")
	assert.Equal(t, expectedDescription, spec.Info.Description, "Info description should match")
	assert.Len(t, spec.Servers, 1, "Should have one server")
	assert.Equal(t, expectedServerURL, spec.Servers[0].URL, "Server URL should match")
	assert.Equal(t, expectedServerDescription, spec.Servers[0].Description, "Server description should match")
	assert.NotNil(t, spec.Paths, "Paths should not be nil")
	assert.NotNil(t, spec.Components, "Components should not be nil")
	assert.NotNil(t, spec.Components.Schemas, "Schemas should not be nil")
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
