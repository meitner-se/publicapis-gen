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

### Example Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/meitner-se/publicapis-gen/specification"
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
    generator := specification.NewSchemaGenerator()
    schema, err := generator.GenerateServiceSchema()
    if err != nil {
        log.Fatal(err)
    }
    
    schemaJSON, err := generator.SchemaToJSON(schema)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Schema JSON:", schemaJSON)
}
```

### Running Tests

```bash
go test ./specification/... -v
```