// Package schemagen provides JSON schema generation, validation, and parsing functionality for specification types.
//
// This package contains the SchemaGenerator type which can generate JSON schemas
// for all specification struct types including Service, Enum, Object, Resource,
// Field, ResourceField, Endpoint, EndpointRequest, and EndpointResponse.
//
// In addition to schema generation, this package provides comprehensive validation
// and parsing capabilities that ensure data conforms to the generated schemas
// before being unmarshaled into Go structs.
//
// # Schema Generation
//
// The package exports a single function GenerateSchemas that writes all generated
// schemas to a bytes.Buffer:
//
//	import (
//	    "bytes"
//	    "github.com/meitner-se/publicapis-gen/specification"
//	    "github.com/meitner-se/publicapis-gen/specification/schemagen"
//	)
//
//	// Load specification
//	service, err := specification.LoadFromFile("api-spec.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Generate schemas
//	var buf bytes.Buffer
//	err = schemagen.GenerateSchemas(&buf, service)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Write to file
//	err = os.WriteFile("schemas.json", buf.Bytes(), 0644)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Validation
//
// Validation functions verify that JSON or YAML data conforms to the generated schemas:
//
//	generator := schemagen.NewSchemaGenerator()
//
//	// Validate JSON data against Service schema
//	jsonData := []byte(`{"name": "MyAPI", "enums": [], "objects": [], "resources": []}`)
//	err := generator.ValidateService(jsonData)
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
//	// Validate YAML data against Enum schema
//	yamlData := []byte(`
//	name: Status
//	description: Status enumeration
//	values:
//	  - name: Active
//	    description: Active status`)
//	err = generator.ValidateEnum(yamlData)
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
// # Parsing with Validation
//
// Parsing functions combine validation and unmarshaling in a single operation:
//
//	generator := schemagen.NewSchemaGenerator()
//
//	// Parse and validate Service from JSON
//	jsonData := []byte(`{"name": "UserAPI", "enums": [], "objects": [], "resources": []}`)
//	service, err := generator.ParseServiceFromJSON(jsonData)
//	if err != nil {
//	    log.Fatalf("Failed to parse service: %v", err)
//	}
//
//	// Parse and validate Enum from YAML
//	yamlData := []byte(`
//	name: Priority
//	description: Task priority levels
//	values:
//	  - name: Low
//	    description: Low priority
//	  - name: High
//	    description: High priority`)
//	enum, err := generator.ParseEnumFromYAML(yamlData)
//	if err != nil {
//	    log.Fatalf("Failed to parse enum: %v", err)
//	}
//
// # Available Functions
//
// The SchemaGenerator provides the following categories of functions:
//
// Schema Generation:
//   - GenerateServiceSchema, GenerateEnumSchema, GenerateObjectSchema
//   - GenerateResourceSchema, GenerateFieldSchema, GenerateResourceFieldSchema
//   - GenerateEndpointSchema, GenerateEndpointRequestSchema, GenerateEndpointResponseSchema
//   - GenerateAllSchemas
//
// Validation:
//   - ValidateService, ValidateEnum, ValidateObject, ValidateResource
//   - ValidateField, ValidateResourceField, ValidateEndpoint
//   - ValidateEndpointRequest, ValidateEndpointResponse
//
// Parsing with Validation:
//   - ParseServiceFromJSON, ParseServiceFromYAML
//   - ParseEnumFromJSON, ParseEnumFromYAML
//   - ParseObjectFromJSON, ParseObjectFromYAML
//   - ParseResourceFromJSON, ParseResourceFromYAML
//
// # Error Handling
//
// Validation and parsing functions provide detailed error messages when failures occur:
//   - Schema generation errors indicate issues with reflection or schema creation
//   - Validation errors include specific details about constraint violations
//   - Parsing errors cover both validation failures and JSON/YAML unmarshaling issues
//
// The package supports generating individual schemas for each type or all schemas
// at once using the GenerateAllSchemas method. All validation and parsing functions
// automatically handle both JSON and YAML input formats.
package schemagen
