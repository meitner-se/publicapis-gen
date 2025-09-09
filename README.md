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
- [📝 Examples](examples/) - Complete example projects

## Installation

```bash
go get github.com/meitner-se/publicapis-gen
```

## Running Tests

```bash
go test ./... -v
```