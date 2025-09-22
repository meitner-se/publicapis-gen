// Package servergen provides server code generation capabilities from specification types.
//
// This package generates server implementations directly from specification.Service
// definitions, complementing the existing root-level server package that works with
// OpenAPI documents. It offers a specification-native approach to server generation
// with deep integration into the specification type system.
//
// Key features:
// - Direct generation from specification.Service structs
// - Native support for specification types (Resource, Field, Enum, Object)
// - Automatic CRUD endpoint generation based on resource operations
// - Custom endpoint support for complex business logic
// - Type-safe parameter and response handling
// - Built-in validation using specification field definitions
//
// # Generation Workflow
//
// The typical workflow involves loading a specification and generating server code:
//
//	import (
//	    "github.com/meitner-se/publicapis-gen/specification"
//	    "github.com/meitner-se/publicapis-gen/specification/servergen"
//	)
//
//	// Load specification from YAML/JSON
//	service, err := specification.LoadFromFile("api-spec.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Generate server code
//	generator := servergen.NewGenerator()
//	err = generator.Generate(service, "generated/server.go")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Resource-Based Generation
//
// Unlike OpenAPI-based generation, this package understands specification resources
// and can automatically generate appropriate handlers for each supported operation:
//
//	// For a resource with operations: ["Create", "Read", "Update", "Delete"]
//	// Automatically generates:
//	// - POST /users (Create)
//	// - GET /users/{id} (Read)
//	// - PUT /users/{id} (Update)
//	// - DELETE /users/{id} (Delete)
//	// - GET /users (List with pagination)
//
// # Custom Endpoints
//
// Resources can define custom endpoints beyond standard CRUD operations:
//
//	// Custom endpoint: POST /users/{id}/activate
//	// Generates handler with proper parameter extraction and validation
//
// # Type Integration
//
// Generated code leverages specification types for:
// - Request/response struct generation from resource fields
// - Enum validation using specification enum definitions
// - Object composition using specification object types
// - Field validation based on type and modifier specifications
//
// This package is designed to work alongside the existing server package,
// providing an alternative generation approach that's more tightly integrated
// with the specification type system and optimized for resource-oriented APIs.
package servergen
