// Package server provides functionality to generate Go server code from OpenAPI documents.
//
// This package uses oapi-codegen to generate Go server implementations from OpenAPI 3.0+
// specifications. It provides a simple interface with sensible defaults for generating
// production-ready server code using the Gin framework.
//
// Key features:
// - Uses Gin HTTP framework by default
// - Strict mode enabled for better type safety
// - Configurable output path and package name
// - Automatic server type generation with embedded spec
// - Chi middleware support for routing
//
// Example usage:
//
//	config := server.Config{
//	    OutputPath:  "generated/server.go",
//	    PackageName: "api",
//	}
//
//	generator := server.New(config)
//	err := generator.GenerateFromFile("openapi.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
package server
