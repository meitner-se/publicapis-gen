// Package schema provides JSON schema generation functionality for specification types.
//
// This package contains the SchemaGenerator type which can generate JSON schemas
// for all specification struct types including Service, Enum, Object, Resource,
// Field, ResourceField, Endpoint, EndpointRequest, and EndpointResponse.
//
// Example usage:
//
//	generator := schema.NewSchemaGenerator()
//	jsonSchema, err := generator.GenerateServiceSchema()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	schemaJSON, err := generator.SchemaToJSON(jsonSchema)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(schemaJSON)
//
// The package supports generating individual schemas for each type or all schemas
// at once using the GenerateAllSchemas method.
package schema
