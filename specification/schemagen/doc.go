// Package schemagen provides JSON schema generation functionality for specification types.
//
// This package follows the same minimalistic API pattern as the servergen package,
// providing a single exported function that generates JSON schemas for all
// specification struct types including Service, Enum, Object, Resource, Field,
// ResourceField, Endpoint, EndpointRequest, and EndpointResponse.
//
// # Schema Generation
//
// The package provides a single function for generating all schemas:
//
//	var buf bytes.Buffer
//	err := schemagen.GenerateSchemas(&buf)
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
// The generated output is a JSON object with schema names as keys and their
// JSON schema definitions as values:
//
//	{
//	  "Service": {
//	    "$schema": "https://json-schema.org/draft/2020-12/schema",
//	    "type": "object",
//	    "properties": {
//	      "name": { "type": "string" },
//	      "enums": { "type": "array", "items": { "$ref": "#/$defs/Enum" } },
//	      ...
//	    },
//	    "required": ["name"]
//	  },
//	  "Enum": { ... },
//	  ...
//	}
//
// # Usage Pattern
//
// This package follows the servergen pattern where the main file is responsible
// for deciding file names and handling the generated content:
//
//	func generateSchemaFiles() error {
//	    var buf bytes.Buffer
//
//	    // Generate schemas
//	    if err := schemagen.GenerateSchemas(&buf); err != nil {
//	        return fmt.Errorf("failed to generate schemas: %w", err)
//	    }
//
//	    // Main file decides the output path
//	    outputPath := "api-schemas.json"
//
//	    // Write the buffer contents to file
//	    if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
//	        return fmt.Errorf("failed to write schemas: %w", err)
//	    }
//
//	    return nil
//	}
//
// # Generated Schemas
//
// The function generates JSON schemas for the following specification types:
//   - Service: The top-level API service definition
//   - Enum: Enumeration type definitions
//   - Object: Object/model definitions
//   - Resource: Resource definitions with CRUD operations
//   - Field: Field definitions used in objects and resources
//   - ResourceField: Field definitions specific to resources
//   - Endpoint: API endpoint definitions
//   - EndpointRequest: Request structure for endpoints
//   - EndpointResponse: Response structure for endpoints
//
// All schemas are generated with the following configuration:
//   - Additional properties are not allowed (strict validation)
//   - References are expanded inline for clarity
//   - Proper JSON Schema draft 2020-12 compatibility
package schemagen
