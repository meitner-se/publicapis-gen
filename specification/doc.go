// Package specification provides Go structs for defining API specifications
// with YAML/JSON tags and JSON schema generation and validation capabilities.
//
// This package contains the core data structures needed to define resource-oriented
// APIs in a structured way. It supports generating OpenAPI specifications from
// these definitions and can produce JSON schemas for validation.
//
// # Core Types
//
// The main types provided by this package are:
//
//   - Service: Represents the overall API service with its enums, objects, and resources
//   - Enum: Defines enumeration types with possible values
//   - Object: Represents shared data structures used across resources
//   - Resource: Defines API resources with their fields and endpoints
//   - Field: Describes individual fields with type information and metadata
//   - ResourceField: Extends Field with operation-specific configuration
//   - Endpoint: Defines individual API endpoints with request/response structure
//
// # JSON Schema Generation
//
// JSON schema generation capabilities are available through the schemagen sub-package,
// which provides a simple function for generating JSON schemas for all specification types.
//
// Basic schema generation example:
//
//	import (
//	    "bytes"
//	    "os"
//	    "github.com/meitner-se/publicapis-gen/specification/schemagen"
//	)
//
//	var buf bytes.Buffer
//	err := schemagen.GenerateSchemas(&buf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Write schemas to file
//	err = os.WriteFile("schemas.json", buf.Bytes(), 0644)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Parsing and Validation
//
// The specification package provides parsing functions that automatically validate
// specifications and apply overlays to ensure complete API definitions.
//
// Parsing example:
//
//	import "github.com/meitner-se/publicapis-gen/specification"
//
//	// Parse from file (supports YAML and JSON)
//	service, err := specification.ParseServiceFromFile("api-spec.yaml")
//	if err != nil {
//	    log.Fatalf("Failed to parse service: %v", err)
//	}
//
//	// Parse from bytes
//	yamlData := []byte(`
//	name: UserAPI
//	enums: []
//	objects: []
//	resources: []`)
//	service, err = specification.ParseServiceFromYAML(yamlData)
//	if err != nil {
//	    log.Fatalf("Failed to parse service: %v", err)
//	}
//
// # Available Parsing Functions
//
// The specification package provides parsing functions for loading and validating specifications:
//
//   - ParseServiceFromFile(filePath) - Parse from YAML or JSON file
//   - ParseServiceFromBytes(data, fileExtension) - Parse from byte data
//   - ParseServiceFromJSON(data) - Parse from JSON byte data
//   - ParseServiceFromYAML(data) - Parse from YAML byte data
//
// All parsing functions automatically:
//   - Validate the specification structure
//   - Apply overlays to generate complete CRUD endpoints
//   - Add default objects (Error, Pagination, Meta)
//   - Generate filter objects for searchable resources
//
// # Error Handling
//
// Parsing functions return detailed error information when validation fails,
// including line and column numbers for YAML files:
//
//	service, err := specification.ParseServiceFromFile("invalid.yaml")
//	if err != nil {
//	    // Error contains specific validation failure details with line numbers
//	    log.Printf("Validation error: %v", err)
//	}
//
// # Serialization Support
//
// All types support both JSON and YAML serialization through appropriate struct tags.
// This makes it easy to store and load API specifications from various file formats.
// The validation functions automatically handle both JSON and YAML input formats.
package specification
