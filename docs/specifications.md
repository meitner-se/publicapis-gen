# Working with Specifications

**Task**: Create, validate, and structure API specifications

## Create a specification from scratch

### Input: Business requirements
- User management system
- Need CRUD operations
- Require status tracking
- Must validate email format

### Output: Structured specification

```yaml
name: "User Management API"
version: "1.0.0"

# Define enums for controlled values
enums:
  - name: "UserStatus" 
    description: "Available user status values"
    values:
      - name: "Active"
        description: "User can access the system"
      - name: "Inactive"
        description: "User account is disabled"
      - name: "Pending"
        description: "User account awaiting verification"

# Define shared objects
objects:
  - name: "Contact"
    description: "Contact information"
    fields:
      - name: "email"
        type: "String"
        description: "Primary email address"
      - name: "phone"
        type: "String"
        description: "Phone number"
        modifiers: ["Nullable"]

# Define resources (main entities)
resources:
  - name: "Users"
    description: "User account management"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "username"
          type: "String"
          description: "Unique username identifier"
        operations: ["Create", "Read"]
      - field:
          name: "contact"
          type: "Contact"
          description: "User contact information"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "status"
          type: "UserStatus"
          description: "Current account status"
        operations: ["Read", "Update"]
      - field:
          name: "tags"
          type: "String"
          description: "User-defined tags"
          modifiers: ["Array", "Nullable"]
        operations: ["Create", "Read", "Update"]
```

## Configure retry behavior

APIs often need retry mechanisms for handling transient failures. You can configure retry behavior directly in your specification:

### Basic retry configuration

```yaml
name: "Resilient API"
version: "1.0.0"

# Configure retry behavior for the entire API
retry:
  strategy: "backoff"
  backoff:
    initial_interval: 500      # Initial delay in milliseconds
    max_interval: 60000        # Maximum delay between retries
    max_elapsed_time: 3600000  # Total timeout in milliseconds
    exponent: 1.5              # Backoff multiplier
  status_codes: ["5XX"]        # HTTP status codes to retry on
  retry_connection_errors: true # Whether to retry connection errors

resources:
  - name: "Users"
    description: "User management with built-in retry support"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "email"
          type: "String"
          description: "User email address"
        operations: ["Create", "Read", "Update"]
```

### Advanced retry configuration

You can customize retry behavior for different scenarios:

```yaml
name: "High Availability API"
version: "1.0.0"

# Custom retry configuration for high-traffic scenarios
retry:
  strategy: "backoff"
  backoff:
    initial_interval: 1000     # Start with 1 second delay
    max_interval: 30000        # Cap at 30 seconds
    max_elapsed_time: 1800000  # Total timeout of 30 minutes
    exponent: 2.0              # Exponential backoff
  status_codes: ["5XX", "429"] # Retry on server errors and rate limits
  retry_connection_errors: false # Don't retry connection errors

resources:
  - name: "Orders"
    description: "Critical order processing"
    operations: ["Create", "Read", "Update"]
    fields:
      - field:
          name: "amount"
          type: "Int"
          description: "Order amount"
        operations: ["Create", "Read"]
```

### Retry configuration reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `strategy` | string | "backoff" | Retry strategy (currently only "backoff" supported) |
| `backoff.initial_interval` | int | 500 | Initial delay between retries in milliseconds |
| `backoff.max_interval` | int | 60000 | Maximum delay between retries in milliseconds |
| `backoff.max_elapsed_time` | int | 3600000 | Maximum total time for all retry attempts |
| `backoff.exponent` | float | 1.5 | Multiplier for exponential backoff |
| `status_codes` | array | ["5XX"] | HTTP status codes that should trigger retries |
| `retry_connection_errors` | boolean | true | Whether to retry on connection errors |

**Note:** If no retry configuration is provided, the system uses sensible defaults shown above.

## Configure timeout behavior

APIs need timeout mechanisms to prevent requests from hanging indefinitely. You can configure timeout behavior directly in your specification:

### Basic timeout configuration

```yaml
name: "Efficient API"
version: "1.0.0"

# Configure timeout behavior for the entire API
timeout:
  timeout: 45000  # Request timeout in milliseconds (45 seconds)

resources:
  - name: "Users"
    description: "User management with custom timeout"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "email"
          type: "String"
          description: "User email address"
        operations: ["Create", "Read", "Update"]
```

### Advanced timeout configuration

You can customize timeout values for different API scenarios:

```yaml
name: "High Performance API"
version: "1.0.0"

# Custom timeout configuration for performance-critical scenarios
timeout:
  timeout: 15000  # Fast response requirement (15 seconds)

resources:
  - name: "RealTimeData"
    description: "Real-time data processing"
    operations: ["Create", "Read"]
    fields:
      - field:
          name: "data"
          type: "String"
          description: "Real-time data payload"
        operations: ["Create", "Read"]
```

### Timeout configuration reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `timeout` | int | 30000 | Request timeout in milliseconds |

**Note:** If no timeout configuration is provided, the system uses a default timeout of 30 seconds (30000 milliseconds).

## Validate specifications

### Task: Check specification correctness

**Input**: YAML/JSON specification file

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
)

func validateSpec(filename string) {
    // Parse with automatic validation
    service, err := specification.ParseServiceFromFile(filename)
    if err != nil {
        log.Printf("❌ Validation failed: %v", err)
        return
    }
    
    fmt.Printf("✅ Specification is valid!\n")
    fmt.Printf("📊 Summary:\n")
    fmt.Printf("  • %d enums\n", len(service.Enums))
    fmt.Printf("  • %d objects\n", len(service.Objects))  
    fmt.Printf("  • %d resources\n", len(service.Resources))
    
    // Check for common issues
    for _, resource := range service.Resources {
        if len(resource.Operations) == 0 {
            fmt.Printf("⚠️  Warning: Resource '%s' has no operations\n", resource.Name)
        }
        if len(resource.Fields) == 0 {
            fmt.Printf("⚠️  Warning: Resource '%s' has no fields\n", resource.Name)
        }
    }
}

func main() {
    validateSpec("user-api.yaml")
}
```

**Output**: Validation results with specific error locations

```
✅ Specification is valid!
📊 Summary:
  • 1 enums
  • 5 objects
  • 1 resources
```

## Structure complex specifications

### Task: Organize large APIs with multiple resources

**Input**: Multi-resource system requirements

```yaml
name: "School Management API"
version: "1.0.0"

# Shared enums across resources
enums:
  - name: "GradeLevel"
    description: "Academic grade levels"
    values:
      - name: "Elementary"
        description: "Grades K-5"
      - name: "Middle" 
        description: "Grades 6-8"
      - name: "High"
        description: "Grades 9-12"

  - name: "StudentStatus"
    description: "Student enrollment status"
    values:
      - name: "Enrolled"
        description: "Currently enrolled"
      - name: "Graduated"
        description: "Successfully graduated"
      - name: "Withdrawn"
        description: "Withdrawn from school"

# Shared objects for reuse
objects:
  - name: "Address"
    description: "Physical address"
    fields:
      - name: "street"
        type: "String" 
        description: "Street address"
      - name: "city"
        type: "String"
        description: "City name"
      - name: "zipCode"
        type: "String"
        description: "Postal code"

  - name: "Person"
    description: "Basic person information"
    fields:
      - name: "firstName"
        type: "String"
        description: "First name"
      - name: "lastName" 
        type: "String"
        description: "Last name"
      - name: "address"
        type: "Address"
        description: "Home address"

# Multiple resources using shared objects
resources:
  - name: "Students"
    description: "Student management"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "person"
          type: "Person"
          description: "Personal information"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "gradeLevel"
          type: "GradeLevel"
          description: "Current grade level"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "status"
          type: "StudentStatus"
          description: "Enrollment status"
        operations: ["Read", "Update"]

  - name: "Teachers"
    description: "Teacher management"
    operations: ["Create", "Read", "Update"]
    fields:
      - field:
          name: "person"
          type: "Person" 
          description: "Personal information"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "subjects"
          type: "String"
          description: "Subjects taught"
          modifiers: ["Array"]
        operations: ["Create", "Read", "Update"]
      - field:
          name: "gradeLevel"
          type: "GradeLevel"
          description: "Grade level taught"
        operations: ["Create", "Read", "Update"]
```

## Parse specifications programmatically

### Task: Load and work with specs in Go code

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
)

func analyzeSpecification(filename string) {
    // Parse specification
    service, err := specification.ParseServiceFromFile(filename)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("📋 API: %s v%s\n", service.Name, service.Version)
    
    // Analyze resources
    for _, resource := range service.Resources {
        fmt.Printf("\n🏗️  Resource: %s\n", resource.Name)
        fmt.Printf("   Operations: %v\n", resource.Operations)
        
        // Show readable fields
        if resource.HasReadOperation() {
            readFields := resource.GetReadableFields()
            fmt.Printf("   Readable fields: %d\n", len(readFields))
            for _, field := range readFields {
                modifiers := ""
                if len(field.Modifiers) > 0 {
                    modifiers = fmt.Sprintf(" [%v]", field.Modifiers)
                }
                fmt.Printf("     • %s: %s%s\n", field.Name, field.Type, modifiers)
            }
        }
        
        // Show generated endpoints
        fmt.Printf("   Endpoints: %d\n", len(resource.Endpoints))
        for _, endpoint := range resource.Endpoints {
            fmt.Printf("     • %s %s\n", endpoint.Method, endpoint.Name)
        }
    }
    
    // Analyze objects
    fmt.Printf("\n📦 Generated Objects:\n")
    for _, obj := range service.Objects {
        fieldTypes := make(map[string]int)
        for _, field := range obj.Fields {
            fieldTypes[field.Type]++
        }
        fmt.Printf("   • %s (%d fields)\n", obj.Name, len(obj.Fields))
    }
}

func main() {
    analyzeSpecification("school-api.yaml")
}
```

**Output**: Detailed analysis of your specification

```
📋 API: School Management API v1.0.0

🏗️  Resource: Students
   Operations: [Create Read Update Delete]
   Readable fields: 4
     • ID: UUID
     • Meta: Meta
     • person: Person
     • gradeLevel: GradeLevel
     • status: StudentStatus
   Endpoints: 6
     • POST Create
     • PATCH Update
     • DELETE Delete
     • GET Get
     • GET List
     • POST Search

📦 Generated Objects:
   • ErrorCode (8 fields)
   • Error (2 fields)
   • ErrorFieldCode (4 fields)
   • ErrorField (2 fields)
   • Pagination (3 fields)
   • Meta (4 fields)
   • Address (3 fields)
   • Person (3 fields)
   • Students (4 fields)
```

## Best Practices

### ✅ Do's
- **Use PascalCase**: `"Create"`, `"String"`, `"Nullable"`  
- **Be descriptive**: Clear field descriptions help generated docs
- **Reuse objects**: Define shared objects to avoid duplication
- **Group operations**: Put related operations on same resource
- **Use enums**: For controlled values like status, type, category

### ❌ Don'ts  
- **Don't use lowercase**: `"create"`, `"string"` will cause validation errors
- **Don't repeat field definitions**: Use shared objects instead
- **Don't skip descriptions**: They become API documentation
- **Don't mix concerns**: Keep resources focused on single entities

## Common Patterns

### Pattern: Status Management
```yaml
enums:
  - name: "Status"
    values:
      - name: "Draft" 
      - name: "Published"
      - name: "Archived"

resources:
  - name: "Posts"
    fields:
      - field:
          name: "status"
          type: "Status"
          description: "Publication status"
        operations: ["Read", "Update"]  # Can't create with status
```

### Pattern: Audit Fields  
```yaml
# Auto-generated by overlay:
# - ID (UUID)
# - Meta object with CreatedAt, UpdatedAt, CreatedBy, UpdatedBy

resources:
  - name: "Users"
    skip_auto_columns: false  # Default: includes audit fields
    fields:
      # Your business fields here
```

### Pattern: Flexible Attributes
```yaml
resources:
  - name: "Products"
    fields:
      - field:
          name: "tags"
          type: "String"
          modifiers: ["Array", "Nullable"]
        operations: ["Create", "Read", "Update"]
      - field:
          name: "metadata"
          type: "String"  # JSON string for flexible data
          modifiers: ["Nullable"]
        operations: ["Create", "Read", "Update"]
```

## Related Tasks

- [📋 Generate OpenAPI specs](openapi.md) - Create documentation from specs
- [🔍 Build advanced filters](filtering.md) - Use generated filter objects
- [✅ Validate with JSON Schema](schema-validation.md) - Ensure spec correctness