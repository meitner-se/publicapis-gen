# Getting Started Guide

**Task**: Create your first API specification and generate complete CRUD endpoints

## What you'll build

A complete User Management API with:
- ✅ CRUD endpoints (Create, Get, List, Update, Delete, Search)
- ✅ Automatic validation and error handling
- ✅ Advanced filtering capabilities
- ✅ OpenAPI 3.1 specification

**Time**: ~10 minutes

## Step 1: Install the library

```bash
go mod init my-api-project
go get github.com/meitner-se/publicapis-gen
```

## Step 2: Create your API specification

Create `user-api.yaml`:

```yaml
name: "User Management API"
version: "1.0.0"
servers:
  - url: "https://api.example.com/v1"
    description: "Production server"

enums:
  - name: "UserStatus"
    description: "Status of a user"
    values:
      - name: "Active"
        description: "User is active"
      - name: "Inactive"
        description: "User is inactive"

resources:
  - name: "Users"
    description: "User management resource"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "username"
          type: "String"
          description: "Unique username"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "email"
          type: "String"
          description: "User email address"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "status"
          type: "UserStatus"
          description: "Current user status"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "bio"
          type: "String"
          description: "User biography"
          modifiers: ["Nullable"]
        operations: ["Create", "Read", "Update"]
```

## Step 3: Generate complete API specification

Create `main.go`:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
)

func main() {
    // Parse and expand the specification
    service, err := specification.ParseServiceFromFile("user-api.yaml")
    if err != nil {
        log.Fatal("Error parsing specification:", err)
    }
    
    // Print summary of what was generated
    fmt.Printf("🎉 Generated API: %s v%s\n", service.Name, service.Version)
    fmt.Printf("📊 Total enums: %d\n", len(service.Enums))
    fmt.Printf("🏗️  Total objects: %d\n", len(service.Objects))
    fmt.Printf("📦 Total resources: %d\n", len(service.Resources))
    
    // Show endpoints for Users resource
    for _, resource := range service.Resources {
        if resource.Name == "Users" {
            fmt.Printf("\n🔗 Generated %d endpoints for %s:\n", 
                len(resource.Endpoints), resource.Name)
            for _, endpoint := range resource.Endpoints {
                fmt.Printf("  • %s %s - %s\n", 
                    endpoint.Method, endpoint.Path, endpoint.Description)
            }
        }
    }
    
    // Save the complete specification
    output, _ := json.MarshalIndent(service, "", "  ")
    fmt.Printf("\n💾 Complete specification saved to complete-api.json\n")
    // In a real app, you'd save this to a file
    _ = output
}
```

## Step 4: Run and see the results

```bash
go run main.go
```

**Expected Output:**
```
🎉 Generated API: User Management API v1.0.0
📊 Total enums: 3
🏗️  Total objects: 8
📦 Total resources: 1

🔗 Generated 6 endpoints for Users:
  • POST  - Create a new Users
  • PATCH /{id} - Update a Users  
  • DELETE /{id} - Delete a Users
  • GET /{id} - Retrieve an existing Users
  • GET  - List all Users
  • POST /_search - Search Users with filtering capabilities

💾 Complete specification saved to complete-api.json
```

## What was automatically generated

From your simple specification, the library automatically created:

### 🔧 Standard Objects
- **Users** object (for API responses)
- **Error** object (for error responses)
- **ErrorField** object (for field validation errors)  
- **Pagination** object (for paginated responses)
- **Meta** object (with ID, CreatedAt, UpdatedAt, etc.)

### 📝 Error Handling Enums
- **ErrorCode** enum (BadRequest, NotFound, etc.)
- **ErrorFieldCode** enum (Required, InvalidValue, etc.)

### 🔍 Filter System (5 objects per resource)
- **UsersFilter** (main filter with all operations)
- **UsersFilterEquals** (exact matches)
- **UsersFilterRange** (for comparable fields)
- **UsersFilterContains** (array matching)
- **UsersFilterLike** (pattern matching for strings)
- **UsersFilterNull** (null checks)

### ⚠️ Request Validation Objects
- **UsersRequestError** (general validation errors)
- **UsersCreateRequestError** (Create endpoint specific errors)
- **UsersUpdateRequestError** (Update endpoint specific errors)

## Next Steps

- [📋 Generate OpenAPI 3.1 spec](openapi.md) - Create documentation
- [✅ Validate your specifications](schema-validation.md) - Ensure correctness
- [🔍 Build advanced filters](filtering.md) - Create powerful search
- [📖 Full API Reference](api-reference.md) - Complete documentation

## Troubleshooting

**Error: `invalid operation: operation 'create'`**
- **Fix**: Use PascalCase: `"Create"` instead of `"create"`

**Error: `invalid field type: field type 'string'`**  
- **Fix**: Use PascalCase: `"String"` instead of `"string"`

**Error: `file does not exist`**
- **Fix**: Check file path and ensure `.yaml` extension