# publicapis-gen

**publicapis-gen** is a lightweight Go library for defining **resource-oriented APIs** with automatic code generation. Define your API once using simple Go structs and YAML, then automatically generate complete OpenAPI specifications, JSON schemas, CRUD endpoints, and advanced filter systems.

## What it does

- **Define APIs** using simple YAML/JSON specifications or Go structs
- **Auto-generate** complete CRUD endpoints (Create, Read, Update, Delete, List, Search)
- **Create** OpenAPI 3.1 specifications for documentation and client generation
- **Generate** JSON schemas for validation
- **Build** comprehensive filter systems for advanced search capabilities
- **Validate** specifications against generated schemas

## When to use it

Perfect for teams who want:
- ✅ Structured API design with minimal boilerplate
- ✅ Automatic generation of standard CRUD operations
- ✅ Consistent API patterns across services
- ✅ OpenAPI documentation without manual maintenance
- ✅ Advanced search capabilities with minimal effort

## Quick Example

Create a simple API specification:

```yaml
# user-api.yaml
name: "User API"
version: "1.0.0"

resources:
  - name: "Users"
    description: "User management"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "username"
          type: "String"
          description: "Username"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "email"
          type: "String"
          description: "Email address"
        operations: ["Create", "Read", "Update"]
```

Generate complete API specification with CRUD endpoints:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
)

func main() {
    // Parse specification from file
    service, err := specification.ParseServiceFromFile("user-api.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Automatically includes:
    // - Complete CRUD endpoints (Create, Get, List, Update, Delete, Search)
    // - Response objects and request validation
    // - Filter systems for advanced search
    // - Pagination and error handling
    
    fmt.Printf("Generated %d endpoints for Users resource\n", 
        len(service.Resources[0].Endpoints))
}
```

## Error Cases

Common errors and solutions:

| Error | Cause | Solution |
|-------|--------|----------|
| `invalid operation: operation 'create' must be one of: [Create Read Update Delete]` | Operation names must use PascalCase | Use `Create` instead of `create` |
| `invalid field type: field type 'string' must be one of...` | Field types must use PascalCase | Use `String` instead of `string` |
| `validation failed: YAML parsing error` | Malformed YAML syntax | Check YAML indentation and syntax |
| `file does not exist` | File path is incorrect | Verify file path and extension |

## Limits

- **Field Types**: UUID, String, Int, Bool, Date, Timestamp + custom Objects/Enums
- **Operations**: Create, Read, Update, Delete only
- **Modifiers**: Nullable, Array only
- **File Formats**: YAML (.yaml, .yml) and JSON (.json)
- **Validation**: Requires valid Go struct tags for JSON schema generation

## Documentation

- [📖 Complete API Reference](docs/api-reference.md) - Full API documentation
- [🚀 Getting Started Guide](docs/getting-started.md) - Step-by-step tutorial
- [⚙️ Working with Specifications](docs/specifications.md) - Creating and validating specs  
- [🔍 Advanced Filtering](docs/filtering.md) - Building complex search functionality
- [📋 OpenAPI Generation](docs/openapi.md) - Generating OpenAPI 3.1 specs
- [✅ JSON Schema Validation](docs/schema-validation.md) - Validating specifications

## Installation

### Using the Go library
```bash
go get github.com/meitner-se/publicapis-gen
```

### Using the CLI tool
```bash
go install github.com/meitner-se/publicapis-gen@latest
```

## CLI Usage

The CLI tool provides `generate` and `diff` commands to process specification files and create OpenAPI documents, JSON schemas, overlays, and server code.

### Basic Usage
```bash
# Show available commands
publicapis-gen help

# Show help for specific commands
publicapis-gen help generate
publicapis-gen help diff
publicapis-gen generate -help
publicapis-gen diff -help

# Generate from specification file (legacy mode)
publicapis-gen generate -file=api-spec.yaml -mode=openapi -output=openapi.json
publicapis-gen generate -file=api-spec.yaml -mode=schema -output=schemas.json
publicapis-gen generate -file=api-spec.yaml -mode=overlay -output=complete-spec.yaml
publicapis-gen generate -file=api-spec.yaml -mode=server -output=server.go

# Using config file (recommended)
publicapis-gen generate -config=build-config.yaml

# Auto-detect default config file
publicapis-gen generate  # Looks for publicapis.yaml or publicapis.yml

# Check for differences between generated files and disk files
publicapis-gen diff -config=build-config.yaml
publicapis-gen diff  # Uses default config file
```

### Configuration File Example
Create a `publicapis.yaml` config file to process multiple specifications:

```yaml
- specification: "users-api.yaml"
  openapi_json: "dist/users-openapi.json"
  schema_json: "dist/users-schema.json"
  overlay_yaml: "dist/users-complete.yaml"
  server_go: "dist/users-server.go"
  server_package: "api"

- specification: "products-api.yaml"  
  openapi_yaml: "dist/products-openapi.yaml"
  schema_json: "dist/products-schema.json"
  server_go: "dist/products-server.go"
```

### Available Modes
- **`openapi`** - Generate OpenAPI 3.1 specification (JSON)
- **`schema`** - Generate JSON schemas for validation  
- **`overlay`** - Generate complete specification with overlays applied
- **`server`** - Generate Go server code with Gin framework

### Options
- **`-file`** - Path to input specification file (YAML or JSON)
- **`-mode`** - Operation mode (openapi, schema, overlay, server)
- **`-output`** - Output file path (optional, auto-generated if not specified)
- **`-config`** - Path to YAML config file for batch processing
- **`-log-level`** - Logging verbosity (debug, info, warn, error, off)

### Commands
- **`generate`** - Generate API specifications and output files
- **`diff`** - Check for differences between generated content and files on disk
- **`help`** - Show help information for commands

## Running Tests

```bash
go test ./... -v
```