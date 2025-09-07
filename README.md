# publicapis-gen

**publicapis-gen** is a lightweight tool for generating **opinionated OpenAPI specifications** for resource-oriented APIs with minimal effort. Its goal is to let developers define APIs with as little overhead as possible, while still producing a complete and standardized OpenAPI document.

## What you can do with publicapis-gen

With `publicapis-gen`, you can:

* Define APIs in a resource-oriented style with minimal boilerplate.
* Automatically generate consistent and maintainable OpenAPI specifications.
* Focus on your API's design and functionality instead of wrestling with spec details.

## Who is this for?

This project is ideal for teams who want the benefits of OpenAPI without the complexity of handcrafting full specifications.

## API Specification Package

The `specification` package provides Go structs for defining API specifications with YAML/JSON serialization support and JSON schema generation capabilities.

### Core Types

* **Service**: Represents the overall API service with its enums, objects, and resources
* **Enum**: Defines enumeration types with possible values
* **Object**: Represents shared data structures used across resources
* **Resource**: Defines API resources with their fields and endpoints
* **Field**: Describes individual fields with type information and metadata
* **ResourceField**: Extends Field with operation-specific configuration
* **Endpoint**: Defines individual API endpoints with request/response structure

### Features

* ✅ Full JSON and YAML serialization support
* ✅ JSON schema generation for all types
* ✅ Comprehensive test coverage
* ✅ Resource-oriented API design
* ✅ Type-safe field definitions with modifiers (nullable, array)
* ✅ Specification overlay functionality for auto-generating Objects from Resources

### Example Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/schema"
)

func main() {
    // Define a service specification
    service := specification.Service{
        Name: "UserAPI",
        Enums: []specification.Enum{
            {
                Name:        "UserStatus",
                Description: "Status of the user",
                Values: []specification.EnumValue{
                    {Name: "Active", Description: "User is active"},
                    {Name: "Inactive", Description: "User is inactive"},
                },
            },
        },
        Objects: []specification.Object{
            {
                Name:        "User",
                Description: "User entity",
                Fields: []specification.Field{
                    {Name: "id", Type: "UUID", Description: "User ID"},
                    {Name: "username", Type: "String", Description: "Username"},
                },
            },
        },
        Resources: []specification.Resource{
            {
                Name:        "Users",
                Description: "User resource",
                Operations:  []string{"Create", "Read", "Update", "Delete"},
                Fields: []specification.ResourceField{
                    {
                        Field: specification.Field{
                            Name:        "id",
                            Type:        "UUID",
                            Description: "User ID",
                        },
                        Operations: []string{"Read"},
                    },
                },
                Endpoints: []specification.Endpoint{
                    {
                        Name:        "GetUser",
                        Description: "Get user by ID",
                        Method:      "GET",
                        Path:        "/users/{id}",
                    },
                },
            },
        },
    }

    // Serialize to JSON
    jsonData, err := json.MarshalIndent(service, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Service JSON:", string(jsonData))

    // Generate JSON schema
    generator := schema.NewSchemaGenerator()
    jsonSchema, err := generator.GenerateServiceSchema()
    if err != nil {
        log.Fatal(err)
    }
    
    schemaJSON, err := generator.SchemaToJSON(jsonSchema)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Schema JSON:", schemaJSON)
}
```

### Specification Overlay

The `ApplyOverlay` function provides a powerful way to automatically generate `Object` definitions from your `Resource` definitions. This is particularly useful for creating consistent data models for read operations.

#### How It Works

When you call `specification.ApplyOverlay()` on a service specification:

1. **Analyzes Resources**: Examines each resource in the specification
2. **Checks for Read Operations**: Identifies resources that have the "Read" operation
3. **Generates Objects**: For each resource with Read operations, creates a new Object with the same name
4. **Includes Read Fields**: Only includes fields that support the "Read" operation in the generated Object
5. **Preserves Existing**: Doesn't overwrite existing Objects with the same name

#### Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/meitner-se/publicapis-gen/specification"
)

func main() {
    // Define a service with resources
    input := &specification.Service{
        Name: "UserAPI",
        Resources: []specification.Resource{
            {
                Name:        "Users",
                Description: "User management resource",
                Operations:  []string{"Create", "Read", "Update", "Delete"},
                Fields: []specification.ResourceField{
                    {
                        Field: specification.Field{
                            Name:        "id",
                            Type:        "UUID",
                            Description: "User ID",
                        },
                        Operations: []string{"Read"},
                    },
                    {
                        Field: specification.Field{
                            Name:        "name",
                            Type:        "String",
                            Description: "User name",
                        },
                        Operations: []string{"Create", "Read", "Update"},
                    },
                    {
                        Field: specification.Field{
                            Name:        "password",
                            Type:        "String",
                            Description: "User password",
                        },
                        Operations: []string{"Create", "Update"}, // No Read - won't be included
                    },
                },
            },
        },
    }

    // Apply overlay to generate Objects from Resources
    result := specification.ApplyOverlay(input)

    // The result now contains a "Users" Object with only the fields that support Read:
    // - id (UUID)
    // - name (String)
    // (password field is excluded since it doesn't have Read operation)

    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated specification:", string(jsonData))
}
```

#### Use Cases

* **API Response Models**: Automatically generate response object schemas from your resource definitions
* **Consistent Data Models**: Ensure your Objects match your Resource field definitions
* **Security**: Exclude sensitive fields (like passwords) that don't support Read operations
* **DRY Principle**: Define your data structure once in Resources, generate Objects automatically

## Constants Usage

This project follows a **zero hardcoded strings** policy for maintainability and consistency. All string literals used in the codebase are defined as package-local constants within the same package where they are used.

### Package-Local Constants Approach

Constants are defined within each package where they are used, including:

- **Error Messages**: All error strings and log messages
- **CRUD Operations**: "Create", "Read", "Update", "Delete" operations  
- **Field Types**: Data types like "UUID", "String", "Int", "Bool", "Timestamp"
- **HTTP Methods**: "GET", "POST", "PUT", "DELETE", etc.
- **Content Types**: "application/json", "multipart/form-data", etc.
- **Field Modifiers**: "array", "nullable", "optional"
- **Schema Properties**: JSON schema field names

### Usage Examples

```go
// ❌ BAD - Hardcoded strings
return errors.New("not implemented")
if containsOperation(operations, "Read") {
slog.ErrorContext(ctx, "failed to run", "error", err)

// ✅ GOOD - Use package-local constants
const (
    errorNotImplemented = "not implemented"
    errorFailedToRun    = "failed to run"  
    logKeyError         = "error"
    OperationRead       = "Read"  // Exported for cross-package usage
)

return errors.New(errorNotImplemented)
if containsOperation(operations, OperationRead) {
slog.ErrorContext(ctx, errorFailedToRun, logKeyError, err)
```

### Adding New Constants

When contributing new functionality:

1. Define constants within the same package where they are used
2. Use clear, descriptive names following Go conventions  
3. Group related constants together in const blocks
4. Use camelCase for unexported constants
5. Export constants only when needed by other packages

### Benefits

- **Consistency**: Uniform error messages and strings across the codebase
- **Maintainability**: Package-local constants reduce dependencies and improve modularity  
- **Encapsulation**: Constants are scoped to where they're actually used
- **Internationalization**: Easy to implement i18n in the future
- **Refactoring**: IDE support for renaming and finding usage

### Running Tests

```bash
go test ./specification/... -v
```