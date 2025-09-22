// Package openapigen provides functionality to generate OpenAPI 3.1 specifications from specification.Service structs.
//
// This package is designed to convert the internal specification format into valid OpenAPI 3.1
// documents that can be used for API documentation, client generation, and API validation.
// The generated specifications follow the OpenAPI 3.1.0 standard and include comprehensive
// schema definitions, endpoint documentation, and validation rules.
//
// The package leverages the pb33f/libopenapi library for robust OpenAPI 3.1 support, providing
// enterprise-grade functionality for generating, validating, and manipulating OpenAPI specifications.
// By using libopenapi's high-level v3.Document types instead of custom definitions, we ensure
// compatibility with the official OpenAPI 3.1 specification and gain access to powerful
// parsing, validation, and serialization capabilities.
//
// # Future Usage
//
// Once implementation is complete, typical usage will be:
//
//	generator := openapigen.NewGenerator()
//	document, err := generator.GenerateFromService(service)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Convert to YAML
//	yamlBytes, err := generator.ToYAML(document)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Convert to JSON
//	jsonBytes, err := generator.ToJSON(document)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # libopenapi Integration
//
// This package uses types from github.com/pb33f/libopenapi:
//   - v3.Document: High-level OpenAPI 3.1 document representation
//   - base.Info: OpenAPI info section with title, description, and version
//   - base.Schema: OpenAPI schema definitions for objects and types
//   - v3.Components: OpenAPI components section for reusable schemas
//   - v3.PathItem: OpenAPI path definitions with HTTP operations
//
// These types provide enterprise-grade OpenAPI 3.1 support with built-in
// validation, parsing, and serialization capabilities, eliminating the need
// for custom OpenAPI type definitions.
//
// # OpenAPI 3.1 Features
//
// The generated specifications will include:
//   - Complete schema definitions for all objects, enums, and resources
//   - Path definitions with proper HTTP methods and parameters
//   - Request and response body schemas
//   - Parameter validation and constraints
//   - Error response definitions
//   - Security definitions (when applicable)
//   - Server definitions and metadata
//
// # Integration with Specification Package
//
// This package is designed to work seamlessly with the specification.Service struct
// and all its related types (Resource, Object, Enum, Field, etc.), converting them
// into their OpenAPI 3.1 equivalents while preserving all semantic meaning and
// validation rules.
//
// The package will support both the raw specification format and the overlay-processed
// format that includes generated CRUD endpoints, filter objects, and error handling.
package openapigen
