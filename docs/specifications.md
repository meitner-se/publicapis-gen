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