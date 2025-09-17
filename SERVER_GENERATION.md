# Server Generation from OpenAPI Specifications

This document describes the new server generation functionality that allows generating Go server code directly from specification files.

## Overview

The server generation feature extends the existing OpenAPI document generation to include Go server code generation using [oapi-codegen/v2](https://github.com/oapi-codegen/oapi-codegen/v2). The complete flow is:

1. **Parse specification file** (YAML/JSON)
2. **Generate OpenAPI document** (OpenAPI 3.0.3)
3. **Generate Go server code** (Echo framework by default)

## Usage

### Command Line (Legacy Mode)

Generate server code directly from a specification file:

```bash
./publicapis-gen generate -file=api-spec.yaml -mode=server -output=server.go
```

### Configuration File (Recommended)

Create a configuration file with server generation:

```yaml
# example-config.yaml
- specification: api-spec.yaml
  openapi_json: api-openapi.json
  server_go: api-server.go
```

Then run:

```bash
./publicapis-gen generate -config=example-config.yaml
```

### Complete Generation Example

Generate all formats including server code:

```yaml
# complete-config.yaml
- specification: school-api.yaml
  openapi_json: school-api-openapi.json
  openapi_yaml: school-api-openapi.yaml
  schema_json: school-api-schemas.json
  overlay_yaml: school-api-overlay.yaml
  server_go: school-api-server.go
```

## Generated Server Code

The generated server code includes:

- **Type definitions** for all objects and enums
- **Echo framework server interface** (`ServerInterface`)
- **Request/response types** for all endpoints
- **Embedded OpenAPI specification** for runtime documentation
- **Validation helpers** and utilities

### Example Generated Interface

```go
package api

// ServerInterface represents all server handlers.
type ServerInterface interface {
    // Create a new user
    // (POST /user)
    UserCreate(ctx echo.Context) error
    
    // Get a user by ID
    // (GET /user/{id})
    UserGet(ctx echo.Context, id string) error
    
    // List users with pagination
    // (GET /user)
    UserList(ctx echo.Context, params UserListParams) error
}
```

## Implementation Notes

### Framework Support

The current implementation generates server code for:
- ✅ **Echo framework** (default)
- ✅ **Type definitions**
- ✅ **Embedded OpenAPI spec**

Additional frameworks supported by oapi-codegen can be enabled by modifying the `server.GeneratorConfig`:
- Chi framework
- Gin framework

### OpenAPI Version Compatibility

The server generation uses **OpenAPI 3.0.3** for maximum compatibility with code generation tools. While the base system supports OpenAPI 3.1.0, server generation requires 3.0.3 due to oapi-codegen compatibility requirements.

### Error Handling

The generated server includes comprehensive error handling:
- Standard HTTP error responses (400, 401, 403, 404, 409, 422, 429, 500)
- Request validation errors with detailed field-level feedback
- Type-safe error objects matching the specification

## Known Limitations

1. **Schema Name Conflicts**: Complex specifications may generate duplicate type names. This can be resolved by simplifying the specification or adjusting the OpenAPI generation logic.

2. **OpenAPI 3.1 → 3.0.3 Compatibility**: Some OpenAPI 3.1 features may not be fully compatible with 3.0.3. The system handles this gracefully by skipping validation.

## Integration Example

Here's how to use the generated server:

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/labstack/echo/v4"
    "./api" // Your generated API package
)

// Implement the ServerInterface
type Server struct{}

func (s *Server) UserCreate(ctx echo.Context) error {
    // Your implementation here
    return ctx.JSON(201, map[string]string{"status": "created"})
}

func (s *Server) UserGet(ctx echo.Context, id string) error {
    // Your implementation here
    return ctx.JSON(200, map[string]string{"id": id})
}

func (s *Server) UserList(ctx echo.Context, params api.UserListParams) error {
    // Your implementation here
    return ctx.JSON(200, map[string]interface{}{
        "data": []map[string]string{},
        "pagination": map[string]int{"offset": 0, "limit": 50, "total": 0},
    })
}

func main() {
    e := echo.New()
    
    server := &Server{}
    api.RegisterHandlers(e, server)
    
    log.Fatal(e.Start(":8080"))
}
```

## Dependencies

The server generation adds these dependencies:
- `github.com/oapi-codegen/oapi-codegen/v2` - Server code generation
- `github.com/getkin/kin-openapi` - OpenAPI document parsing (required by oapi-codegen)

These dependencies are automatically included when using the server generation feature.