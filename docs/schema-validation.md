# JSON Schema Validation

**Task**: Validate API specifications and requests using generated JSON schemas

## Generate JSON schemas

### Input: API specification

```yaml
name: "Blog API"
version: "1.0.0"

enums:
  - name: "PostStatus"
    description: "Blog post status"
    values:
      - name: "Draft"
        description: "Post is in draft"
      - name: "Published"
        description: "Post is published"
      - name: "Archived"  
        description: "Post is archived"

resources:
  - name: "Posts"
    description: "Blog post management"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "title"
          type: "String"
          description: "Post title"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "content"
          type: "String"
          description: "Post content"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "status"
          type: "PostStatus"
          description: "Publication status"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "tags"
          type: "String"
          description: "Post tags"
          modifiers: ["Array", "Nullable"]
        operations: ["Create", "Read", "Update"]
```

### Output: Complete JSON schemas

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification/schema"
)

func generateSchemas() {
    generator := schema.NewSchemaGenerator()
    
    // Generate schemas for all types
    schemas, err := generator.GenerateAllSchemas()
    if err != nil {
        log.Fatal("Failed to generate schemas:", err)
    }
    
    fmt.Printf("📋 Generated %d JSON schemas:\n", len(schemas))
    for name, schema := range schemas {
        fmt.Printf("  • %s\n", name)
    }
    
    // Convert Service schema to JSON for inspection
    serviceSchemaJSON, err := generator.SchemaToJSON(schemas["Service"])
    if err != nil {
        log.Fatal("Failed to convert schema:", err)
    }
    
    fmt.Printf("\n🔍 Service schema (first 500 chars):\n%s...\n", 
        serviceSchemaJSON[:500])
}

func main() {
    generateSchemas()
}
```

**Generated schemas:**
```
📋 Generated 9 JSON schemas:
  • Service
  • Enum
  • EnumValue
  • Object
  • Resource
  • Field
  • ResourceField
  • Endpoint
  • EndpointRequest
  • EndpointResponse

🔍 Service schema (first 500 chars):
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Name of the service"
    },
    "version": {
      "type": "string",
      "description": "Version of the service"
    },
    "servers": {
      "type": "array",
      "items": {
        "$ref": "#/$defs/ServiceServer"
      }
    },
    "enums": {
      "type": "array",
      "items": {
        "$ref": "#/$defs/Enum"
      }
    }
  }
}...
```

## Validate specification files

### Task: Check YAML/JSON specification correctness

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/meitner-se/publicapis-gen/specification/schema"
)

func validateSpecFile(filename string) {
    // Read specification file
    data, err := os.ReadFile(filename)
    if err != nil {
        log.Fatal("Failed to read file:", err)
    }
    
    // Create schema generator
    generator := schema.NewSchemaGenerator()
    
    // Validate the specification
    err = generator.ValidateService(data)
    if err != nil {
        fmt.Printf("❌ Validation failed for %s:\n", filename)
        fmt.Printf("   %v\n", err)
        return
    }
    
    fmt.Printf("✅ %s is valid!\n", filename)
    
    // Parse and show summary
    service, err := generator.ParseServiceFromYAML(data)
    if err != nil {
        fmt.Printf("⚠️  Parse warning: %v\n", err)
        return
    }
    
    fmt.Printf("📊 Summary:\n")
    fmt.Printf("   • Name: %s v%s\n", service.Name, service.Version)
    fmt.Printf("   • Enums: %d\n", len(service.Enums))
    fmt.Printf("   • Objects: %d\n", len(service.Objects))
    fmt.Printf("   • Resources: %d\n", len(service.Resources))
}

func main() {
    // Validate multiple specification files
    files := []string{
        "blog-api.yaml",
        "user-api.yaml", 
        "invalid-api.yaml",
    }
    
    for _, file := range files {
        validateSpecFile(file)
        fmt.Println()
    }
}
```

**Output with validation results:**
```
✅ blog-api.yaml is valid!
📊 Summary:
   • Name: Blog API v1.0.0
   • Enums: 1
   • Objects: 8
   • Resources: 1

✅ user-api.yaml is valid!
📊 Summary:
   • Name: User API v1.0.0  
   • Enums: 2
   • Objects: 6
   • Resources: 1

❌ Validation failed for invalid-api.yaml:
   validation error at line 15, column 18 (operations): invalid operation: operation 'create' must be one of: [Create Read Update Delete]
```

## Validate API requests

### Task: Validate incoming HTTP requests against generated schemas

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/schema"
)

// Middleware to validate request bodies
func ValidateRequestBody(schemaGenerator *schema.SchemaGenerator, schemaName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method != "POST" && r.Method != "PATCH" && r.Method != "PUT" {
                next.ServeHTTP(w, r)
                return
            }
            
            // Read request body
            var body json.RawMessage
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
            }
            
            // Validate against schema
            if err := validateJSONAgainstSchema(schemaGenerator, schemaName, body); err != nil {
                validationError := map[string]interface{}{
                    "Code": "BadRequest",
                    "Message": "Request validation failed",
                    "Details": err.Error(),
                }
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(validationError)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func validateJSONAgainstSchema(generator *schema.SchemaGenerator, schemaName string, data json.RawMessage) error {
    // This would validate the JSON data against the named schema
    // Implementation depends on your schema validation library
    
    switch schemaName {
    case "CreatePost":
        return generator.ValidateCreatePost(data)
    case "UpdatePost":
        return generator.ValidateUpdatePost(data)
    case "SearchFilter":
        return generator.ValidateSearchFilter(data)
    default:
        return fmt.Errorf("unknown schema: %s", schemaName)
    }
}

func createPost(w http.ResponseWriter, r *http.Request) {
    // Request has already been validated by middleware
    var request CreatePostRequest
    json.NewDecoder(r.Body).Decode(&request)
    
    // Process the validated request
    post := &Post{
        ID:      generateID(),
        Title:   request.Title,
        Content: request.Content,
        Status:  request.Status,
        Tags:    request.Tags,
        Meta: Meta{
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }
    
    // Save post...
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func main() {
    // Load specification and generate schemas
    service, _ := specification.ParseServiceFromFile("blog-api.yaml")
    generator := schema.NewSchemaGenerator()
    
    // Setup routes with validation middleware
    http.Handle("/posts", ValidateRequestBody(generator, "CreatePost")(http.HandlerFunc(createPost)))
    
    fmt.Println("🔒 API server with validation running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Test request validation

### Task: Verify validation catches invalid requests

**Valid request:**
```json
POST /posts
Content-Type: application/json

{
  "title": "My Blog Post",
  "content": "This is the content of my blog post.",
  "status": "Draft",
  "tags": ["programming", "go"]
}
```

**Response**: `201 Created`

**Invalid request (missing required field):**
```json  
POST /posts
Content-Type: application/json

{
  "content": "This is the content of my blog post.",
  "status": "Draft",
  "tags": ["programming", "go"]
}
```

**Response**: `400 Bad Request`
```json
{
  "Code": "BadRequest",
  "Message": "Request validation failed", 
  "Details": "missing required field 'title'"
}
```

**Invalid request (wrong enum value):**
```json
POST /posts  
Content-Type: application/json

{
  "title": "My Blog Post",
  "content": "This is the content of my blog post.",
  "status": "published",  // Should be "Published" (PascalCase)
  "tags": ["programming", "go"]
}
```

**Response**: `400 Bad Request`
```json
{
  "Code": "BadRequest", 
  "Message": "Request validation failed",
  "Details": "field 'status': 'published' is not a valid value, must be one of: [Draft, Published, Archived]"
}
```

## Generate request validation schemas

### Task: Create specific schemas for each endpoint

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/schema"
)

func generateRequestSchemas() {
    // Load expanded specification (includes generated endpoints)
    service, err := specification.ParseServiceFromFile("blog-api.yaml") 
    if err != nil {
        log.Fatal(err)
    }
    
    generator := schema.NewSchemaGenerator()
    
    // Generate schemas for each endpoint's request body
    for _, resource := range service.Resources {
        for _, endpoint := range resource.Endpoints {
            if len(endpoint.Request.BodyParams) > 0 {
                schemaName := resource.Name + endpoint.Name + "Request"
                schema := generateEndpointRequestSchema(generator, endpoint.Request.BodyParams)
                
                // Save schema to file
                schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
                filename := fmt.Sprintf("schemas/%s.json", schemaName)
                os.WriteFile(filename, schemaJSON, 0644)
                
                fmt.Printf("📄 Generated %s\n", filename)
            }
        }
    }
}

func generateEndpointRequestSchema(generator *schema.SchemaGenerator, bodyParams []specification.Field) map[string]interface{} {
    schema := map[string]interface{}{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "type":    "object",
        "properties": make(map[string]interface{}),
        "required": []string{},
    }
    
    properties := schema["properties"].(map[string]interface{})
    required := []string{}
    
    for _, field := range bodyParams {
        // Generate field schema based on type and modifiers
        fieldSchema := generateFieldSchema(field)
        properties[field.Name] = fieldSchema
        
        // Check if field is required
        if !field.IsNullable() && field.Default == "" && !field.IsArray() {
            required = append(required, field.Name)
        }
    }
    
    schema["required"] = required
    return schema
}

func generateFieldSchema(field specification.Field) map[string]interface{} {
    fieldSchema := map[string]interface{}{}
    
    switch field.Type {
    case "String":
        fieldSchema["type"] = "string"
    case "Int":
        fieldSchema["type"] = "integer"  
    case "Bool":
        fieldSchema["type"] = "boolean"
    case "UUID":
        fieldSchema["type"] = "string"
        fieldSchema["format"] = "uuid"
    case "Date":
        fieldSchema["type"] = "string"
        fieldSchema["format"] = "date"
    case "Timestamp":
        fieldSchema["type"] = "string"
        fieldSchema["format"] = "date-time"
    default:
        // Custom type - reference to another schema
        fieldSchema["$ref"] = fmt.Sprintf("#/$defs/%s", field.Type)
    }
    
    if field.Description != "" {
        fieldSchema["description"] = field.Description
    }
    
    if field.IsArray() {
        arraySchema := map[string]interface{}{
            "type":  "array", 
            "items": fieldSchema,
        }
        if field.IsNullable() {
            arraySchema["nullable"] = true
        }
        return arraySchema
    }
    
    if field.IsNullable() {
        fieldSchema["nullable"] = true
    }
    
    return fieldSchema
}

func main() {
    os.MkdirAll("schemas", 0755)
    generateRequestSchemas()
}
```

**Generated schema files:**
```
📄 Generated schemas/PostsCreateRequest.json
📄 Generated schemas/PostsUpdateRequest.json
📄 Generated schemas/PostsSearchRequest.json
```

**Example PostsCreateRequest.json:**
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "Post title"
    },
    "content": {
      "type": "string", 
      "description": "Post content"
    },
    "status": {
      "$ref": "#/$defs/PostStatus",
      "description": "Publication status"
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "nullable": true,
      "description": "Post tags"
    }
  },
  "required": ["title", "content", "status"]
}
```

## Validate with external tools

### Task: Use JSON schema validation tools

**Command line validation with ajv-cli:**
```bash
# Install ajv-cli
npm install -g ajv-cli

# Validate a specification file
ajv validate -s schemas/ServiceSchema.json -d blog-api.json
```

**Python validation script:**
```python
import json
import jsonschema
from jsonschema import validate

def validate_api_spec(spec_file, schema_file):
    # Load schema
    with open(schema_file, 'r') as f:
        schema = json.load(f)
    
    # Load specification
    with open(spec_file, 'r') as f:
        spec = json.load(f)
    
    try:
        validate(instance=spec, schema=schema)
        print(f"✅ {spec_file} is valid")
    except jsonschema.ValidationError as e:
        print(f"❌ Validation error in {spec_file}:")
        print(f"   {e.message}")
        print(f"   Path: {' -> '.join(str(x) for x in e.path)}")

# Validate specification
validate_api_spec("blog-api.json", "schemas/ServiceSchema.json")
```

## Integration with CI/CD

### Task: Automatically validate specifications in CI

**GitHub Actions example:**
```yaml
name: Validate API Specifications
on:
  pull_request:
    paths: ['**/*.yaml', '**/*.yml', '**/*.json']

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Validate specifications
        run: |
          go run scripts/validate-specs.go
          if [ $? -ne 0 ]; then
            echo "❌ Specification validation failed"
            exit 1
          fi
          echo "✅ All specifications are valid"
      
      - name: Generate schemas
        run: |
          go run scripts/generate-schemas.go
          git diff --exit-code schemas/
          if [ $? -ne 0 ]; then
            echo "❌ Generated schemas are out of date"
            echo "Run 'go run scripts/generate-schemas.go' and commit changes"
            exit 1
          fi
```

## Common validation patterns

### Pattern: Required field validation
```json
{
  "type": "object",
  "properties": {
    "title": {"type": "string"},
    "content": {"type": "string"}
  },
  "required": ["title", "content"]
}
```

### Pattern: Enum validation
```json
{
  "type": "object", 
  "properties": {
    "status": {
      "type": "string",
      "enum": ["Draft", "Published", "Archived"]
    }
  }
}
```

### Pattern: Array with constraints
```json
{
  "type": "object",
  "properties": {
    "tags": {
      "type": "array",
      "items": {"type": "string"},
      "minItems": 1,
      "maxItems": 10,
      "uniqueItems": true
    }
  }
}
```

### Pattern: Nested object validation
```json
{
  "type": "object",
  "properties": {
    "author": {
      "$ref": "#/$defs/User"
    }
  },
  "$defs": {
    "User": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "email": {"type": "string", "format": "email"}
      },
      "required": ["name", "email"]
    }
  }
}
```

## Best Practices

### ✅ Do's
- **Validate early**: Check specifications at build time, not runtime
- **Use specific schemas**: Generate schemas for each endpoint's requests/responses
- **Provide clear errors**: Include field path and expected values in error messages
- **Cache schemas**: Load and compile schemas once, reuse for multiple validations
- **Test validation**: Include invalid data in your test suites

### ❌ Don'ts  
- **Don't skip validation**: Always validate user input against schemas
- **Don't ignore errors**: Handle validation failures gracefully
- **Don't hardcode schemas**: Generate them from your specifications
- **Don't forget updates**: Regenerate schemas when specifications change

## Related Tasks

- [🚀 Getting Started](getting-started.md) - Create your first validated specification
- [⚙️ Working with Specifications](specifications.md) - Structure specs for better validation
- [📋 Generate OpenAPI](openapi.md) - Create validated API documentation  
- [🔍 Advanced Filtering](filtering.md) - Validate complex filter requests