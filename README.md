# publicapis-gen

**publicapis-gen** is a lightweight Go library and specification format for defining **opinionated resource-oriented APIs** with comprehensive code generation capabilities. Its goal is to let developers define APIs with as little overhead as possible, while automatically generating complete API specifications, validation schemas, and filter systems.

## What you can do with publicapis-gen

With `publicapis-gen`, you can:

* Define APIs in a resource-oriented style with minimal boilerplate using Go structs.
* Automatically generate consistent Objects from Resources for response schemas.
* Auto-generate complete CRUD endpoints (Create, Read, Update, Delete, Get, List, Search) from resource definitions.
* Automatically create comprehensive filter systems for search functionality.
* Generate standard error handling objects and request validation schemas.
* Validate JSON/YAML specifications against generated JSON schemas.
* Focus on your API's design and functionality instead of wrestling with implementation details.

## Project Status

**Note**: This project is currently in development. The main executable is not yet implemented - the library is focused on the specification format and code generation functionality. The core `specification` package and `schema` package provide the full functionality described in this documentation.

## Who is this for?

This project is ideal for teams who want the benefits of structured API design with automatic code generation, comprehensive validation, and minimal manual specification work.

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
* ✅ JSON schema generation and validation for all types
* ✅ Comprehensive test coverage
* ✅ Resource-oriented API design
* ✅ Type-safe field definitions with modifiers (nullable, array)
* ✅ Specification overlay functionality for auto-generating Objects from Resources
* ✅ Auto-generation of complete CRUD endpoints (Create, Update, Delete, Get, List, Search)
* ✅ Comprehensive filter system generation for advanced search capabilities
* ✅ Standard error handling objects (Error, ErrorField, ErrorCode enums)
* ✅ Request validation error objects for all body parameters
* ✅ Automatic pagination object generation for list operations

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

    // Apply overlays to generate complete specification
    completeSpec := specification.ApplyOverlay(&service)
    completeSpecWithFilters := specification.ApplyFilterOverlay(completeSpec)

    // Serialize the complete specification to JSON
    completeJSON, err := json.MarshalIndent(completeSpecWithFilters, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Complete Specification:", string(completeJSON))

    // Generate and validate JSON schema
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

The `ApplyOverlay` function provides a powerful way to automatically generate `Object` definitions, `Create` endpoints, and `Update` endpoints from your `Resource` definitions. This is particularly useful for creating consistent data models for read operations and standardized CRUD endpoints.

#### How It Works

When you call `specification.ApplyOverlay()` on a service specification:

1. **Analyzes Resources**: Examines each resource in the specification
2. **Generates Objects**: For each resource with "Read" operations, creates a new Object with the same name
3. **Includes Read Fields**: Only includes fields that support the "Read" operation in the generated Object
4. **Generates Create Endpoints**: For each resource with "Create" operations, creates a standardized Create endpoint
5. **Generates Update Endpoints**: For each resource with "Update" operations, creates a standardized Update endpoint
6. **Includes Operation Fields**: Uses all fields that support the respective operation as body parameters
7. **Returns Resource Object**: Both Create and Update endpoints return the resource object (201 Created / 200 OK)
8. **Preserves Existing**: Doesn't overwrite existing Objects or endpoints with the same name

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

#### Create Endpoint Generation

The overlay automatically generates Create endpoints with the following characteristics:

* **HTTP Method**: POST
* **Path**: Empty (to be combined with resource base path)
* **Request**: JSON body with all fields that have Create operation
* **Response**: 201 Created status code, returns the created resource object
* **Content Type**: application/json for both request and response

Example of generated Create endpoint for a Users resource:

```json
{
  "name": "Create",
  "title": "Create Users",
  "description": "Create a new Users",
  "method": "POST",
  "path": "",
  "request": {
    "content_type": "application/json",
    "body_params": [
      {
        "name": "name",
        "type": "String",
        "description": "User name"
      },
      {
        "name": "email", 
        "type": "String",
        "description": "User email"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 201,
    "body_object": "Users"
  }
}
```

#### Update Endpoint Generation

The overlay also automatically generates Update endpoints with the following characteristics:

* **HTTP Method**: PATCH
* **Path**: /{id} (includes ID as path parameter)
* **Request**: JSON body with all fields that have Update operation, plus ID as path parameter
* **Response**: 200 OK status code, returns the updated resource object
* **Content Type**: application/json for both request and response

Example of generated Update endpoint for a Users resource:

```json
{
  "name": "Update",
  "title": "Update Users",
  "description": "Update a Users",
  "method": "PATCH",
  "path": "/{id}",
  "request": {
    "content_type": "application/json",
    "path_params": [
      {
        "name": "id",
        "type": "UUID",
        "description": "The unique identifier of the resource to update"
      }
    ],
    "body_params": [
      {
        "name": "name",
        "type": "String",
        "description": "User name"
      },
      {
        "name": "email",
        "type": "String",
        "description": "User email"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 200,
    "body_object": "Users"
  }
}
```

#### Delete Endpoint Generation

The overlay automatically generates Delete endpoints for resources with "Delete" operations:

* **HTTP Method**: DELETE
* **Path**: /{id} (includes ID as path parameter)
* **Request**: Only ID as path parameter, no body
* **Response**: 204 No Content status code, no response body
* **Content Type**: application/json for request

Example of generated Delete endpoint for a Users resource:

```json
{
  "name": "Delete",
  "title": "Delete Users",
  "description": "Delete a Users",
  "method": "DELETE",
  "path": "/{id}",
  "request": {
    "content_type": "application/json",
    "path_params": [
      {
        "name": "id",
        "type": "UUID",
        "description": "The unique identifier of the resource to delete"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 204
  }
}
```

#### Get Endpoint Generation

The overlay automatically generates Get endpoints for resources with "Read" operations:

* **HTTP Method**: GET
* **Path**: /{id} (includes ID as path parameter)
* **Request**: Only ID as path parameter, no body
* **Response**: 200 OK status code, returns the resource object
* **Content Type**: application/json for response

Example of generated Get endpoint for a Users resource:

```json
{
  "name": "Get",
  "title": "Retrieve an existing Users",
  "description": "Retrieves the `Users` with the given ID.",
  "method": "GET",
  "path": "/{id}",
  "request": {
    "content_type": "application/json",
    "path_params": [
      {
        "name": "id",
        "type": "UUID",
        "description": "The unique identifier of the Users to retrieve"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 200,
    "body_object": "Users"
  }
}
```

#### List Endpoint Generation

The overlay automatically generates List endpoints for resources with "Read" operations:

* **HTTP Method**: GET
* **Path**: (empty - base resource path)
* **Request**: Query parameters for pagination (limit, offset)
* **Response**: 200 OK status code, returns paginated array of resource objects
* **Content Type**: application/json for response
* **Pagination**: Automatic limit (default: 50) and offset (default: 0) parameters

Example of generated List endpoint for a Users resource:

```json
{
  "name": "List",
  "title": "List all Users",
  "description": "Returns a paginated list of all `Users` in your organization.",
  "method": "GET",
  "path": "",
  "request": {
    "content_type": "application/json",
    "query_params": [
      {
        "name": "limit",
        "type": "Int",
        "description": "The maximum number of items to return (default: 50)",
        "default": "50"
      },
      {
        "name": "offset",
        "type": "Int", 
        "description": "The number of items to skip before starting to return results (default: 0)",
        "default": "0"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 200,
    "body_fields": [
      {
        "name": "data",
        "description": "Array of Users objects",
        "type": "Users",
        "modifiers": ["array"]
      },
      {
        "name": "Pagination",
        "description": "Pagination information",
        "type": "Pagination"
      }
    ]
  }
}
```

#### Search Endpoint Generation

The overlay automatically generates Search endpoints for resources with "Read" operations:

* **HTTP Method**: POST
* **Path**: /_search
* **Request**: Filter body parameter of type `<Resource>Filter`, plus pagination query parameters
* **Response**: 200 OK status code, returns paginated array of resource objects (same as List)
* **Content Type**: application/json for both request and response
* **Filter System**: Uses comprehensive filter objects (requires `ApplyFilterOverlay`)

Example of generated Search endpoint for a Users resource:

```json
{
  "name": "Search",
  "title": "Search Users",
  "description": "Search for `Users` with filtering capabilities.",
  "method": "POST",
  "path": "/_search",
  "request": {
    "content_type": "application/json",
    "query_params": [
      {
        "name": "limit",
        "type": "Int",
        "description": "The maximum number of items to return (default: 50)",
        "default": "50"
      },
      {
        "name": "offset",
        "type": "Int",
        "description": "The number of items to skip before starting to return results (default: 0)", 
        "default": "0"
      }
    ],
    "body_params": [
      {
        "name": "Filter",
        "description": "Filter criteria to search for specific records",
        "type": "UsersFilter"
      }
    ]
  },
  "response": {
    "content_type": "application/json",
    "status_code": 200,
    "body_fields": [
      {
        "name": "data",
        "description": "Array of Users objects",
        "type": "Users",
        "modifiers": ["array"]
      },
      {
        "name": "Pagination", 
        "description": "Pagination information",
        "type": "Pagination"
      }
    ]
  }
}
```

#### Use Cases

* **API Response Models**: Automatically generate response object schemas from your resource definitions
* **Consistent Data Models**: Ensure your Objects match your Resource field definitions
* **Complete CRUD Operations**: Generate all standard endpoints (Create, Read, Update, Delete, Get, List, Search)
* **RESTful API Design**: Generate proper REST endpoints with correct HTTP methods and status codes
* **Paginated Operations**: Automatic pagination support for List and Search endpoints
* **Advanced Search**: Comprehensive filter system for complex search operations
* **Security**: Exclude sensitive fields (like passwords) that don't support Read operations from response objects
* **DRY Principle**: Define your data structure once in Resources, generate Objects and endpoints automatically

## Filter System and ApplyFilterOverlay

The `ApplyFilterOverlay` function provides a comprehensive filter system for advanced search capabilities. This should be called after `ApplyOverlay` to ensure all Objects are available for filter generation.

### How ApplyFilterOverlay Works

When you call `specification.ApplyFilterOverlay()` on a service specification:

1. **Analyzes Objects**: Examines each Object in the specification (including those generated by `ApplyOverlay`)
2. **Generates Filter Objects**: Creates comprehensive filter structures for each Object
3. **Creates Filter Types**: Generates multiple specialized filter objects per Object type
4. **Supports All Field Types**: Handles primitive types, custom objects, and nested relationships
5. **Provides Advanced Operations**: Supports equality, range, contains, LIKE, and null checks
6. **Enables Complex Logic**: Supports OR conditions and nested filtering

### Generated Filter Types

For each Object (e.g., `Users`), the filter overlay generates:

* **`UsersFilter`**: Main filter object with all filter operations
* **`UsersFilterEquals`**: Equality and inequality filters
* **`UsersFilterRange`**: Range filters for comparable types (Int, Date, Timestamp)
* **`UsersFilterContains`**: Contains/not-contains filters for arrays and nested objects
* **`UsersFilterLike`**: LIKE/NOT LIKE pattern matching for strings
* **`UsersFilterNull`**: Null/not-null checks for nullable fields

### Filter Operations

#### Main Filter Structure

Each main filter (e.g., `UsersFilter`) includes:

```go
type UsersFilter struct {
    Equals          *UsersFilterEquals   // Equality filters
    NotEquals       *UsersFilterEquals   // Inequality filters  
    GreaterThan     *UsersFilterRange    // Greater than filters
    SmallerThan     *UsersFilterRange    // Less than filters
    GreaterOrEqual  *UsersFilterRange    // Greater or equal filters
    SmallerOrEqual  *UsersFilterRange    // Less or equal filters
    Contains        *UsersFilterContains // Contains filters
    NotContains     *UsersFilterContains // Not contains filters
    Like            *UsersFilterLike     // LIKE pattern filters
    NotLike         *UsersFilterLike     // NOT LIKE pattern filters
    Null            *UsersFilterNull     // Null checks
    NotNull         *UsersFilterNull     // Not null checks
    OrCondition     bool                 // OR vs AND logic
    NestedFilters   []UsersFilter        // Nested filter conditions
}
```

#### Equality Filters (UsersFilterEquals)

Contains all Object fields as nullable types for exact matching:

```go
type UsersFilterEquals struct {
    Id       *string  // UUID field
    Name     *string  // String field  
    Email    *string  // String field
    Age      *int     // Int field
    // Nested objects use their filter equivalents
    Profile  *ProfileFilterEquals
}
```

#### Range Filters (UsersFilterRange) 

Contains only comparable fields (Int, Date, Timestamp) plus nested objects:

```go
type UsersFilterRange struct {
    Age        *int     // Int field for range operations
    CreatedAt  *string  // Timestamp field for range operations
    // Nested objects use their filter equivalents  
    Profile    *ProfileFilterRange
}
```

#### Contains Filters (UsersFilterContains)

Contains all fields except timestamps as arrays for "contains any of" operations:

```go
type UsersFilterContains struct {
    Id    []string  // Array of UUIDs to match against
    Name  []string  // Array of names to match against  
    Age   []int     // Array of ages to match against
    // Nested objects use their filter equivalents
    Profile  *ProfileFilterContains
}
```

#### LIKE Filters (UsersFilterLike)

Contains only string fields plus nested objects for pattern matching:

```go
type UsersFilterLike struct {
    Name     *string  // String field for LIKE operations
    Email    *string  // String field for LIKE operations
    // Nested objects use their filter equivalents
    Profile  *ProfileFilterLike
}
```

#### Null Filters (UsersFilterNull)

Contains only nullable fields and arrays for null/not-null checks:

```go
type UsersFilterNull struct {
    Description *bool   // Check if nullable description is null/not null
    Tags        *bool   // Check if array field is null/not null
    // Nested objects use boolean checks
    Profile     *bool   // Check if nested profile is null/not null
}
```

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
                            Name:        "email",
                            Type:        "String",
                            Description: "User email",
                        },
                        Operations: []string{"Create", "Read", "Update"},
                    },
                    {
                        Field: specification.Field{
                            Name:        "age",
                            Type:        "Int", 
                            Description: "User age",
                        },
                        Operations: []string{"Create", "Read", "Update"},
                    },
                },
            },
        },
    }

    // First, apply the main overlay to generate Objects and endpoints
    result := specification.ApplyOverlay(input)

    // Then, apply the filter overlay to generate filter objects
    filterResult := specification.ApplyFilterOverlay(result)

    // The result now contains comprehensive filter objects:
    // - UsersFilter (main filter with all operations)
    // - UsersFilterEquals (for exact matches)
    // - UsersFilterRange (for age range queries)  
    // - UsersFilterContains (for array matching)
    // - UsersFilterLike (for name/email pattern matching)
    // - UsersFilterNull (for null checks)

    jsonData, err := json.MarshalIndent(filterResult, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated specification with filters:", string(jsonData))
}
```

### Filter System Benefits

* **Comprehensive Coverage**: Every field gets appropriate filter types based on its data type
* **Type Safety**: Filters respect field types and only offer valid operations
* **Nested Support**: Automatically handles nested object filtering
* **Performance Optimized**: Different filter types for different query patterns
* **Complex Logic**: Supports OR conditions and nested filter combinations
* **RESTful Integration**: Designed to work seamlessly with generated Search endpoints

### Integration with Search Endpoints

The filter system is designed to work with the automatically generated Search endpoints. The Search endpoint accepts a filter object of type `<Resource>Filter` in the request body, enabling complex queries like:

```json
{
  "Filter": {
    "Equals": {
      "name": "John"
    },
    "GreaterThan": {
      "age": 25
    },
    "Like": {
      "email": "%@company.com"
    },
    "OrCondition": false
  }
}
```

## Standard Objects and Error Handling

The `ApplyOverlay` function automatically generates standard objects that are essential for API operations. These objects are created automatically if they don't already exist in your specification.

### Error Handling Objects

#### ErrorCode Enum

Automatically generated enum with standard HTTP error codes:

```go
type ErrorCode string

const (
    BadRequest          ErrorCode = "BadRequest"          // 400 status code
    Unauthorized        ErrorCode = "Unauthorized"        // 401 status code  
    Forbidden           ErrorCode = "Forbidden"           // 403 status code
    NotFound            ErrorCode = "NotFound"            // 404 status code
    Conflict            ErrorCode = "Conflict"            // 409 status code
    UnprocessableEntity ErrorCode = "UnprocessableEntity" // 422 status code
    RateLimited         ErrorCode = "RateLimited"         // 429 status code
    Internal            ErrorCode = "Internal"            // 5xx status code
)
```

#### Error Object

Standard error response object containing error code and message:

```go
type Error struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
}
```

#### ErrorFieldCode Enum

Enum for field-level validation errors:

```go
type ErrorFieldCode string

const (
    AlreadyExists ErrorFieldCode = "AlreadyExists" // Unique constraint violation
    Required      ErrorFieldCode = "Required"      // Missing required field
    NotFound      ErrorFieldCode = "NotFound"      // Referenced resource not found
    InvalidValue  ErrorFieldCode = "InvalidValue"  // Invalid field value
)
```

#### ErrorField Object

Field-specific error information for validation errors:

```go
type ErrorField struct {
    Code    ErrorFieldCode `json:"code"`
    Message string         `json:"message"`
}
```

### Pagination Object

Automatically generated for List and Search endpoints:

```go
type Pagination struct {
    Offset int `json:"offset"` // Number of items skipped
    Limit  int `json:"limit"`  // Maximum items returned
}
```

### Request Error Objects

The overlay also automatically generates `RequestError` objects for types used in request body parameters. These provide validation error information specific to each request type.

#### Object-Based Request Errors

For each Object used in request body parameters, a corresponding `<Object>RequestError` is generated:

```go
// If Users object is used in request body, generates:
type UsersRequestError struct {
    Id    *ErrorField `json:"id,omitempty"`
    Name  *ErrorField `json:"name,omitempty"`  
    Email *ErrorField `json:"email,omitempty"`
    // All fields become nullable ErrorField types
}
```

#### Endpoint-Specific Request Errors

For each endpoint with body parameters, a specific `<Resource><Endpoint>RequestError` is generated:

```go
// For Users Create endpoint, generates:
type UsersCreateRequestError struct {
    Name     *ErrorField `json:"name,omitempty"`
    Email    *ErrorField `json:"email,omitempty"`
    Password *ErrorField `json:"password,omitempty"`
    // Only fields used in this endpoint's body parameters
}
```

### Benefits of Standard Objects

* **Consistency**: Uniform error responses across all endpoints
* **Validation**: Detailed field-level error information
* **Standards Compliance**: HTTP status codes mapped to semantic error types
* **Developer Experience**: Clear error messages for API consumers
* **Automatic Generation**: No manual definition required
* **Type Safety**: Strong typing for all error conditions

## Constants Usage

This project follows a **zero hardcoded strings** policy for maintainability and consistency. All string literals used in the codebase are defined as package-local constants within the same package where they are used.

### Package-Local Constants Approach

Constants are defined within each package where they are used, including:

- **Error Messages**: All error strings and log messages
- **CRUD Operations**: "Create", "Read", "Update", "Delete" operations  
- **Field Types**: Data types like "UUID", "String", "Int", "Bool", "Timestamp", "Date"
- **HTTP Methods**: "GET", "POST", "PATCH", "PUT", "DELETE"
- **Content Types**: "application/json" and other MIME types
- **Field Modifiers**: "array", "nullable" modifiers
- **HTTP Status Codes**: Response status codes (200, 201, 204, etc.)
- **Endpoint Constants**: Names, paths, titles, and descriptions for generated endpoints
- **Filter Constants**: Suffixes and field names for filter system
- **Error Code Values**: Standard error codes and descriptions for API responses
- **Pagination Constants**: Field names and descriptions for pagination objects

### Usage Examples

```go
// ❌ BAD - Hardcoded strings
return errors.New("not implemented")
if containsOperation(operations, "Read") {
slog.ErrorContext(ctx, "failed to run", "error", err)
endpoint.Method = "POST"
endpoint.Path = "/_search"
response.StatusCode = 200
field.Type = "UUID"

// ✅ GOOD - Use package-local constants
const (
    // Error messages and log keys
    errorNotImplemented = "not implemented"
    errorFailedToRun    = "failed to run"  
    logKeyError         = "error"
    
    // CRUD Operations (exported for cross-package usage)
    OperationRead       = "Read"
    OperationCreate     = "Create"
    
    // HTTP Methods and Status Codes
    httpMethodPost      = "POST"
    searchResponseStatusCode = 200
    
    // Endpoint paths and field types
    searchEndpointPath  = "/_search"
    FieldTypeUUID       = "UUID"
)

return errors.New(errorNotImplemented)
if containsOperation(operations, OperationRead) {
slog.ErrorContext(ctx, errorFailedToRun, logKeyError, err)
endpoint.Method = httpMethodPost
endpoint.Path = searchEndpointPath
response.StatusCode = searchResponseStatusCode
field.Type = FieldTypeUUID
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

## Schema Package

The `schema` package provides JSON schema generation and validation capabilities for all specification types. This enables validation of YAML/JSON configuration files against the specification format.

### Schema Generation and Validation Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/meitner-se/publicapis-gen/specification/schema"
)

func main() {
    generator := schema.NewSchemaGenerator()

    // Generate schemas for all types
    schemas, err := generator.GenerateAllSchemas()
    if err != nil {
        log.Fatal(err)
    }

    // Convert Service schema to JSON
    serviceSchemaJSON, err := generator.SchemaToJSON(schemas["Service"])
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Service Schema:", serviceSchemaJSON)

    // Validate YAML/JSON data against Service schema
    yamlData := []byte(`
name: "UserAPI"
resources:
  - name: "Users"
    description: "User management resource"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "id"
          type: "UUID"
          description: "User ID"
        operations: ["Read"]
`)

    // Validate the YAML data
    err = generator.ValidateService(yamlData)
    if err != nil {
        log.Fatal("Validation failed:", err)
    }
    fmt.Println("YAML data is valid!")

    // Parse and validate YAML into a Service struct
    service, err := generator.ParseServiceFromYAML(yamlData)
    if err != nil {
        log.Fatal("Parse failed:", err)
    }
    fmt.Printf("Parsed service: %+v\n", service)
}
```

### Available Schema Operations

* **Generate schemas** for Service, Enum, Object, Resource, Field, ResourceField, Endpoint, EndpointRequest, EndpointResponse
* **Validate JSON/YAML** data against generated schemas
* **Parse and validate** JSON/YAML data into Go structs
* **Convert between formats** (YAML to JSON automatically)

### Running Tests

```bash
go test ./... -v
```