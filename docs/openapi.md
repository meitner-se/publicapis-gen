# Generate OpenAPI 3.1 Specifications

**Task**: Convert your API specifications into OpenAPI 3.1 documentation

## Generate OpenAPI from specification

### Input: Your API specification
```yaml
name: "Pet Store API"
version: "2.0.0"
servers:
  - url: "https://api.petstore.com/v2"
    description: "Production server"

enums:
  - name: "PetStatus"
    description: "Pet availability status"
    values:
      - name: "Available"
        description: "Pet is available for adoption"
      - name: "Pending" 
        description: "Pet adoption is pending"
      - name: "Sold"
        description: "Pet has been adopted"

resources:
  - name: "Pets"
    description: "Pet management"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "name"
          type: "String"
          description: "Pet name"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "status"
          type: "PetStatus"
          description: "Availability status"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "photoUrls"
          type: "String"
          description: "Pet photos"
          modifiers: ["Array"]
        operations: ["Create", "Read", "Update"]
```

### Output: Complete OpenAPI 3.1 specification

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/openapi"
)

func generateOpenAPI() {
    // Parse your specification
    service, err := specification.ParseServiceFromFile("petstore-api.yaml")
    if err != nil {
        log.Fatal("Failed to parse specification:", err)
    }
    
    // Generate OpenAPI 3.1 document
    generator := openapi.NewGenerator()
    document, err := generator.GenerateFromService(service)
    if err != nil {
        log.Fatal("Failed to generate OpenAPI:", err)
    }
    
    // Convert to YAML (for documentation)
    yamlBytes, err := generator.ToYAML(document)
    if err != nil {
        log.Fatal("Failed to convert to YAML:", err)
    }
    
    // Save OpenAPI specification
    err = os.WriteFile("openapi.yaml", yamlBytes, 0644)
    if err != nil {
        log.Fatal("Failed to save file:", err)
    }
    
    fmt.Println("✅ OpenAPI 3.1 specification generated: openapi.yaml")
    
    // Also generate JSON version
    jsonBytes, err := generator.ToJSON(document)
    if err != nil {
        log.Fatal("Failed to convert to JSON:", err)
    }
    
    err = os.WriteFile("openapi.json", jsonBytes, 0644)
    if err != nil {
        log.Fatal("Failed to save JSON file:", err)
    }
    
    fmt.Println("✅ JSON version generated: openapi.json")
}

func main() {
    generateOpenAPI()
}
```

**Generated OpenAPI structure:**
```yaml
openapi: 3.1.0
info:
  title: Pet Store API
  version: 2.0.0
servers:
  - url: https://api.petstore.com/v2
    description: Production server

components:
  schemas:
    # Enum definitions
    PetStatus:
      type: string
      enum: [Available, Pending, Sold]
      description: Pet availability status
    
    # Object schemas
    Pets:
      type: object
      properties:
        ID:
          type: string
          format: uuid
        Meta:
          $ref: '#/components/schemas/Meta'
        name:
          type: string
          description: Pet name
        status:
          $ref: '#/components/schemas/PetStatus'
        photoUrls:
          type: array
          items:
            type: string
      required: [ID, Meta, name, status, photoUrls]
    
    # Error schemas
    Error:
      type: object
      properties:
        Code:
          $ref: '#/components/schemas/ErrorCode'
        Message:
          type: string
      required: [Code, Message]

paths:
  /pets:
    get:
      summary: List all Pets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 50
        - name: offset
          in: query  
          schema:
            type: integer
            default: 0
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Pets'
                  Pagination:
                    $ref: '#/components/schemas/Pagination'
    
    post:
      summary: Create Pets
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                status:
                  $ref: '#/components/schemas/PetStatus'
                photoUrls:
                  type: array
                  items:
                    type: string
              required: [name, status, photoUrls]
      responses:
        '201':
          description: Created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pets'
        '422':
          description: Validation Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PetsCreateRequestError'

  /pets/{id}:
    get:
      summary: Retrieve an existing Pets
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pets'
    
    patch:
      summary: Update Pets
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                status:
                  $ref: '#/components/schemas/PetStatus'
                photoUrls:
                  type: array
                  items:
                    type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pets'
    
    delete:
      summary: Delete Pets
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: No Content

  /pets/_search:
    post:
      summary: Search Pets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 50
        - name: offset
          in: query
          schema:
            type: integer
            default: 0
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                Filter:
                  $ref: '#/components/schemas/PetsFilter'
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Pets'
                  Pagination:
                    $ref: '#/components/schemas/Pagination'
```

## Customize OpenAPI generation

### Task: Add custom documentation and examples

```go
package main

import (
    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/openapi"
)

func customizeOpenAPI() {
    service, _ := specification.ParseServiceFromFile("api.yaml")
    
    generator := openapi.NewGenerator()
    
    // Add custom configuration
    generator.SetContactInfo("API Team", "api@company.com", "https://company.com")
    generator.SetLicenseInfo("MIT", "https://opensource.org/licenses/MIT")  
    generator.AddTag("pets", "Everything about your Pets")
    generator.AddSecurityScheme("apiKey", openapi.SecurityScheme{
        Type: "apiKey",
        Name: "X-API-Key",  
        In:   "header",
    })
    
    // Generate with customization
    document, _ := generator.GenerateFromService(service)
    
    // Add examples to specific operations
    generator.AddOperationExample("getPet", "application/json", map[string]interface{}{
        "ID":     "123e4567-e89b-12d3-a456-426614174000",
        "name":   "Fluffy",
        "status": "Available",
        "photoUrls": []string{
            "https://example.com/photos/fluffy1.jpg",
            "https://example.com/photos/fluffy2.jpg",
        },
    })
    
    yamlBytes, _ := generator.ToYAML(document)
    // Save customized specification...
}
```

## Validate OpenAPI output

### Task: Ensure generated specification is valid

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
    "github.com/meitner-se/publicapis-gen/specification/openapi"
)

func validateOpenAPI() {
    // Generate OpenAPI
    service, _ := specification.ParseServiceFromFile("api.yaml")
    generator := openapi.NewGenerator()
    document, _ := generator.GenerateFromService(service)
    
    // Validate the generated OpenAPI document
    valid, errors := generator.ValidateDocument(document)
    
    if valid {
        fmt.Println("✅ Generated OpenAPI 3.1 specification is valid")
    } else {
        fmt.Println("❌ OpenAPI validation errors:")
        for _, err := range errors {
            fmt.Printf("  • %s\n", err)
        }
    }
    
    // Check for completeness
    stats := generator.GetDocumentStats(document)
    fmt.Printf("\n📊 Document Statistics:\n")
    fmt.Printf("  • Paths: %d\n", stats.PathCount)
    fmt.Printf("  • Operations: %d\n", stats.OperationCount)
    fmt.Printf("  • Schemas: %d\n", stats.SchemaCount)
    fmt.Printf("  • Parameters: %d\n", stats.ParameterCount)
}
```

## Use with documentation tools

### Task: Generate interactive API documentation

**With Swagger UI:**
```bash
# Serve your OpenAPI spec with Swagger UI
docker run -p 8080:8080 \
  -v $(pwd):/api \
  -e SWAGGER_JSON=/api/openapi.yaml \
  swaggerapi/swagger-ui
```

**With Redoc:**
```bash
# Generate static documentation with Redoc
npx redoc-cli build openapi.yaml --output docs.html
```

**With Postman:**
```bash
# Import into Postman for testing
# File -> Import -> openapi.yaml
```

## Integration with CI/CD

### Task: Automatically update API documentation

**GitHub Actions example:**
```yaml
name: Generate API Documentation
on:
  push:
    paths: ['**/*.yaml', '**/*.yml']

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Generate OpenAPI
        run: |
          go run scripts/generate-openapi.go
          
      - name: Deploy to GitHub Pages  
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

## Common OpenAPI Features

### Generated automatically:
- ✅ **Complete path definitions** for all CRUD endpoints
- ✅ **Request/response schemas** with proper validation  
- ✅ **Parameter definitions** with types and constraints
- ✅ **Error response schemas** for all HTTP status codes
- ✅ **Component schemas** for reusable objects
- ✅ **Enum definitions** with descriptions
- ✅ **Filter schemas** for search endpoints

### Available for customization:
- 🔧 **Server definitions** and environment URLs
- 🔧 **Security schemes** (API keys, OAuth2, etc.)
- 🔧 **Contact and license** information  
- 🔧 **Tags and grouping** for organization
- 🔧 **Examples and descriptions** for better docs
- 🔧 **Custom headers** and metadata

## Best Practices

### ✅ Do's
- **Add server URLs**: Help users understand where to call your API
- **Include examples**: Real data makes documentation more useful
- **Use tags**: Group related operations together
- **Add contact info**: Let users know how to get help
- **Version your specs**: Track changes over time

### ❌ Don'ts
- **Don't skip descriptions**: They become the main documentation
- **Don't ignore validation**: Ensure your OpenAPI is valid
- **Don't forget examples**: They help users understand expected data
- **Don't hardcode URLs**: Use server definitions for flexibility

## Related Tasks

- [🚀 Getting Started](getting-started.md) - Create your first specification
- [⚙️ Working with Specifications](specifications.md) - Structure complex APIs  
- [✅ JSON Schema Validation](schema-validation.md) - Validate your specs
- [🔍 Advanced Filtering](filtering.md) - Document search capabilities