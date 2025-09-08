// Package openapi provides functionality to generate OpenAPI 3.1 specifications from specification.Service structs.
//
// This package is designed to convert the internal specification format into valid OpenAPI 3.1
// documents that can be used for API documentation, client generation, and API validation.
// The generated specifications follow the OpenAPI 3.1.0 standard and include comprehensive
// schema definitions, endpoint documentation, and validation rules.
//
// The package leverages the libopenapi library for robust OpenAPI 3.1 support, providing
// enterprise-grade functionality for generating, validating, and manipulating OpenAPI specifications.
//
// # Future Usage
//
// Once implementation is complete, typical usage will be:
//
//	generator := openapi.NewGenerator()
//	spec, err := generator.GenerateFromService(service)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Convert to YAML
//	yamlBytes, err := generator.ToYAML(spec)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Convert to JSON
//	jsonBytes, err := generator.ToJSON(spec)
//	if err != nil {
//	    log.Fatal(err)
//	}
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
package openapi
