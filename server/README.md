# Server Package

The `server` package provides functionality to generate Go server code from OpenAPI documents using `oapi-codegen`. It is designed to generate production-ready server implementations using the Gin HTTP framework with sensible defaults.

## Features

- **Gin Framework**: Uses Gin HTTP framework by default for high-performance web services
- **Strict Mode**: Enables strict mode for better type safety and validation
- **Configurable Output**: Allows customization of output path and package name
- **Embedded Spec**: Automatically embeds the OpenAPI specification in generated code
- **Comprehensive Testing**: Full test coverage with integration tests

## Installation

The server package is part of the `publicapis-gen` module:

```go
import "github.com/meitner-se/publicapis-gen/server"
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "github.com/meitner-se/publicapis-gen/server"
)

func main() {
    // Create configuration
    config := server.Config{
        OutputPath:  "generated/api_server.go",
        PackageName: "api",
    }

    // Create generator
    generator, err := server.New(config)
    if err != nil {
        log.Fatal(err)
    }

    // Generate server code from OpenAPI file
    err = generator.GenerateFromFile("openapi.yaml")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Advanced Usage

```go
package main

import (
    "context"
    "os"
    "github.com/meitner-se/publicapis-gen/server"
)

func main() {
    // Read OpenAPI specification
    specData, err := os.ReadFile("api-spec.json")
    if err != nil {
        log.Fatal(err)
    }

    // Create generator with custom config
    config := server.Config{
        OutputPath:  "internal/handlers/server.gen.go",
        PackageName: "handlers",
    }

    generator, err := server.New(config)
    if err != nil {
        log.Fatal(err)
    }

    // Generate with context
    ctx := context.Background()
    err = generator.GenerateFromDataWithContext(ctx, specData)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

The `server.Config` struct has the following configurable fields:

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `OutputPath` | `string` | Path where generated server code will be written | `"server.gen.go"` |
| `PackageName` | `string` | Package name for the generated code | `"api"` |

### Hardcoded Configuration

The following settings are hardcoded for optimal Gin server generation:

- **GinServer**: `true` - Uses Gin HTTP framework
- **Strict**: `true` - Enables strict mode for type safety
- **Models**: `true` - Generates model structs
- **EmbeddedSpec**: `true` - Embeds OpenAPI spec in generated code
- **SkipFmt**: `false` - Formats generated code
- **SkipPrune**: `false` - Prunes unused code

## Methods

### Generator Creation

#### `New(config Config) (*Generator, error)`

Creates a new generator instance with the provided configuration. Applies defaults and validates the configuration.

```go
config := server.Config{
    OutputPath:  "api/server.go",
    PackageName: "myapi",
}

generator, err := server.New(config)
```

### Code Generation

#### `GenerateFromFile(specPath string) error`

Generates Go server code from an OpenAPI specification file (JSON or YAML).

```go
err := generator.GenerateFromFile("openapi.yaml")
```

#### `GenerateFromData(specData []byte) error`

Generates Go server code from OpenAPI specification data in memory.

```go
err := generator.GenerateFromData(openAPIBytes)
```

#### `GenerateFromReader(reader io.Reader) error`

Generates Go server code from an OpenAPI specification reader.

```go
file, _ := os.Open("spec.json")
defer file.Close()
err := generator.GenerateFromReader(file)
```

### Context-Aware Methods

All generation methods have context-aware variants:

- `GenerateFromFileWithContext(ctx context.Context, specPath string) error`
- `GenerateFromDataWithContext(ctx context.Context, specData []byte) error`
- `GenerateFromReaderWithContext(ctx context.Context, reader io.Reader) error`

### Configuration Access

#### `GetConfig() Config`

Returns a copy of the generator's configuration.

```go
config := generator.GetConfig()
fmt.Printf("Output: %s, Package: %s", config.OutputPath, config.PackageName)
```

## Generated Code Structure

The generated server code includes:

- **Server Interface**: Defines handler methods for all API operations
- **Request/Response Types**: Type-safe structs for all API data
- **Handler Registration**: Functions to register handlers with Gin router
- **Embedded Spec**: The original OpenAPI specification embedded as a constant
- **Validation**: Automatic request/response validation based on OpenAPI schema

### Example Generated Interface

```go
type ServerInterface interface {
    // GET /users
    GetUsers(ctx *gin.Context, params GetUsersParams)
    
    // POST /users
    CreateUser(ctx *gin.Context)
    
    // GET /users/{userId}
    GetUserById(ctx *gin.Context, userId string)
}
```

## Error Handling

The package provides structured error handling:

- **ConfigError**: Configuration validation errors
- **Spec Reading Errors**: Issues with OpenAPI specification files
- **Generation Errors**: Problems during code generation
- **File Write Errors**: Issues writing generated code

## Best Practices

1. **Use Absolute Paths**: Prefer absolute paths for `OutputPath` to avoid ambiguity
2. **Validate OpenAPI Spec**: Ensure your OpenAPI specification is valid before generation
3. **Version Control**: Add generated files to `.gitignore` and regenerate during build
4. **Package Naming**: Use descriptive package names that reflect the API purpose

## Examples

### Simple Health Check API

```yaml
# health-api.yaml
openapi: 3.0.0
info:
  title: Health API
  version: 1.0.0
paths:
  /health:
    get:
      operationId: getHealth
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
```

```go
// Generate server code
config := server.Config{
    OutputPath:  "health_server.go",
    PackageName: "health",
}

generator, _ := server.New(config)
generator.GenerateFromFile("health-api.yaml")
```

### Implementation

```go
// Implement the generated interface
type HealthServer struct{}

func (h *HealthServer) GetHealth(ctx *gin.Context) {
    ctx.JSON(200, gin.H{"status": "ok"})
}

// Set up the server
func main() {
    server := &HealthServer{}
    router := gin.Default()
    
    // Register generated handlers
    health.RegisterHandlers(router, server)
    
    router.Run(":8080")
}
```

## Testing

The package includes comprehensive tests covering:

- Configuration validation
- Code generation from various input formats
- Error handling scenarios
- Integration tests with real OpenAPI specifications

Run tests with:

```bash
go test ./server/... -v
```

## Dependencies

- `github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen` - Code generation engine
- `github.com/getkin/kin-openapi/openapi3` - OpenAPI specification parsing
- `github.com/gin-gonic/gin` - HTTP framework (generated code dependency)

## License

This package is part of the `publicapis-gen` project and follows the same licensing terms.