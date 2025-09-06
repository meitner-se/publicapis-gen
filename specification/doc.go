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
// # JSON Schema Generation and Validation
//
// JSON schema generation and validation capabilities are available through the schema sub-package,
// which provides the SchemaGenerator type for producing JSON schemas and validating data against them.
//
// Basic schema generation example:
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
// # Validation Support
//
// The schema package provides validation functions that verify JSON/YAML data against
// the generated schemas before parsing. This ensures that only valid data is processed.
//
// Validation example:
//
//	generator := schema.NewSchemaGenerator()
//
//	// Validate JSON data
//	jsonData := []byte(`{"name": "MyAPI", "enums": [], "objects": [], "resources": []}`)
//	err := generator.ValidateService(jsonData)
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
//	// Validate YAML data
//	yamlData := []byte(`
//	name: MyAPI
//	enums: []
//	objects: []
//	resources: []`)
//	err = generator.ValidateService(yamlData)
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
// # Parsing with Validation
//
// The schema package provides parsing functions that combine validation and unmarshaling
// into a single operation, ensuring that only valid, well-formed data is parsed.
//
// Parsing example:
//
//	generator := schema.NewSchemaGenerator()
//
//	// Parse and validate JSON
//	jsonData := []byte(`{
//	    "name": "UserAPI",
//	    "enums": [
//	        {
//	            "name": "Status",
//	            "description": "User status",
//	            "values": [
//	                {"name": "Active", "description": "Active user"},
//	                {"name": "Inactive", "description": "Inactive user"}
//	            ]
//	        }
//	    ],
//	    "objects": [],
//	    "resources": []
//	}`)
//
//	service, err := generator.ParseServiceFromJSON(jsonData)
//	if err != nil {
//	    log.Fatalf("Failed to parse service: %v", err)
//	}
//
//	// Parse and validate YAML
//	yamlData := []byte(`
//	name: UserAPI
//	enums:
//	  - name: Status
//	    description: User status
//	    values:
//	      - name: Active
//	        description: Active user
//	      - name: Inactive
//	        description: Inactive user
//	objects: []
//	resources: []`)
//
//	service, err = generator.ParseServiceFromYAML(yamlData)
//	if err != nil {
//	    log.Fatalf("Failed to parse service: %v", err)
//	}
//
// # Available Validation and Parsing Functions
//
// The schema package provides validation and parsing functions for all main specification types:
//
//   - ValidateService, ParseServiceFromJSON, ParseServiceFromYAML
//   - ValidateEnum, ParseEnumFromJSON, ParseEnumFromYAML
//   - ValidateObject, ParseObjectFromJSON, ParseObjectFromYAML
//   - ValidateResource, ParseResourceFromJSON, ParseResourceFromYAML
//   - ValidateField, ValidateResourceField
//   - ValidateEndpoint, ValidateEndpointRequest, ValidateEndpointResponse
//
// # Error Handling
//
// Validation functions return detailed error information when validation fails,
// including specific details about what constraints were violated:
//
//	err := generator.ValidateService(invalidData)
//	if err != nil {
//	    // Error contains specific validation failure details
//	    log.Printf("Validation error: %v", err)
//	}
//
// # Serialization Support
//
// All types support both JSON and YAML serialization through appropriate struct tags.
// This makes it easy to store and load API specifications from various file formats.
// The validation functions automatically handle both JSON and YAML input formats.
package specification
