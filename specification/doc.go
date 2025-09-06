// Package specification provides Go structs for defining API specifications
// with YAML/JSON tags and JSON schema generation capabilities.
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
// JSON schema generation capabilities are available through the schema sub-package,
// which provides the SchemaGenerator type for producing JSON schemas for all specification types.
//
// Example usage:
//
//	import "github.com/meitner-se/publicapis-gen/specification/schema"
//
//	generator := schema.NewSchemaGenerator()
//	jsonSchema, err := generator.GenerateServiceSchema()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	schemaJSON, err := generator.SchemaToJSON(jsonSchema)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Serialization Support
//
// All types support both JSON and YAML serialization through appropriate struct tags.
// This makes it easy to store and load API specifications from various file formats.
package specification
