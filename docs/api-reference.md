# Complete API Reference

**Reference**: Full documentation for all packages and functions

## Package: specification

Core package for defining API specifications.

### Types

#### Service
Main service definition containing all API components.

```go
type Service struct {
    Name      string          `json:"name"`              // Service name
    Version   string          `json:"version,omitempty"` // Service version  
    Servers   []ServiceServer `json:"servers,omitempty"` // Server definitions
    Enums     []Enum          `json:"enums"`             // Enum definitions
    Objects   []Object        `json:"objects"`           // Shared objects
    Resources []Resource      `json:"resources"`         // API resources
}
```

**Methods:**
- `IsObject(fieldType string) bool` - Check if field type is custom object
- `HasObject(name string) bool` - Check if service contains object
- `HasEnum(name string) bool` - Check if service contains enum
- `GetObject(name string) *Object` - Get object by name

#### Resource
Defines an API resource with operations and fields.

```go
type Resource struct {
    Name            string          `json:"name"`                      // Resource name
    Description     string          `json:"description"`               // Resource description
    Operations      []string        `json:"operations"`                // Allowed operations
    Fields          []ResourceField `json:"fields"`                    // Resource fields
    Endpoints       []Endpoint      `json:"endpoints"`                 // Custom endpoints
    SkipAutoColumns bool            `json:"skip_auto_columns,omitempty"` // Skip auto fields
}
```

**Methods:**
- `HasCreateOperation() bool` - Check for Create operation
- `HasReadOperation() bool` - Check for Read operation  
- `HasUpdateOperation() bool` - Check for Update operation
- `HasDeleteOperation() bool` - Check for Delete operation
- `GetPluralName() string` - Get pluralized resource name
- `GetCreateBodyParams() []Field` - Get fields for Create endpoint
- `GetUpdateBodyParams() []Field` - Get fields for Update endpoint
- `GetReadableFields() []Field` - Get fields for Read operations
- `HasEndpoint(name string) bool` - Check if endpoint exists
- `ShouldSkipAutoColumns() bool` - Check if auto-columns should be skipped

#### Field
Basic field definition with type and metadata.

```go
type Field struct {
    Name        string   `json:"name"`                // Field name
    Description string   `json:"description"`         // Field description
    Type        string   `json:"type"`                // Field type
    Default     string   `json:"default,omitempty"`   // Default value
    Example     string   `json:"example,omitempty"`   // Example value
    Modifiers   []string `json:"modifiers,omitempty"` // Field modifiers
}
```

**Methods:**
- `IsArray() bool` - Check if field has Array modifier
- `IsNullable() bool` - Check if field has Nullable modifier  
- `TagJSON() string` - Get JSON tag name (camelCase)
- `IsRequired(service *Service) bool` - Check if field is required

#### ResourceField  
Field with operation-specific configuration.

```go
type ResourceField struct {
    Field                          // Embedded field
    Operations []string `json:"operations"` // Operations field supports
}
```

**Methods:**
- `HasCreateOperation() bool` - Check for Create operation
- `HasReadOperation() bool` - Check for Read operation
- `HasUpdateOperation() bool` - Check for Update operation  
- `HasDeleteOperation() bool` - Check for Delete operation

#### Object
Shared object definition for reuse across resources.

```go
type Object struct {
    Name        string  `json:"name"`        // Object name
    Description string  `json:"description"` // Object description
    Fields      []Field `json:"fields"`      // Object fields
}
```

**Methods:**
- `HasField(name string) bool` - Check if object has field
- `GetField(name string) *Field` - Get field by name

#### Enum
Enumeration with possible values.

```go
type Enum struct {
    Name        string      `json:"name"`        // Enum name
    Description string      `json:"description"` // Enum description
    Values      []EnumValue `json:"values"`      // Possible values
}
```

#### EnumValue
Single value in an enumeration.

```go
type EnumValue struct {
    Name        string `json:"name"`        // Value name
    Description string `json:"description"` // Value description
}
```

#### Endpoint
Custom endpoint definition.

```go
type Endpoint struct {
    Name        string            `json:"name"`        // Endpoint name
    Title       string            `json:"title"`       // Endpoint title
    Summary     string            `json:"summary"`     // Endpoint summary (short plain text)
    Description string            `json:"description"` // Endpoint description (longer, supports markdown)
    Method      string            `json:"method"`      // HTTP method
    Path        string            `json:"path"`        // URL path
    Request     EndpointRequest   `json:"request"`     // Request definition
    Response    EndpointResponse  `json:"response"`    // Response definition
}
```

**Methods:**
- `GetFullPath(resourceName string) string` - Get full path including resource

#### EndpointRequest
Request structure for an endpoint.

```go
type EndpointRequest struct {
    ContentType string  `json:"content_type"`        // Request content type
    Headers     []Field `json:"headers"`             // Request headers
    PathParams  []Field `json:"path_params"`         // Path parameters
    QueryParams []Field `json:"query_params"`        // Query parameters
    BodyParams  []Field `json:"body_params"`         // Body parameters
}
```

**Methods:**
- `GetRequiredBodyParams(service *Service) []string` - Get required parameter names

#### EndpointResponse
Response structure for an endpoint.

```go
type EndpointResponse struct {
    ContentType string  `json:"content_type"`          // Response content type
    StatusCode  int     `json:"status_code"`           // HTTP status code
    Headers     []Field `json:"headers"`               // Response headers
    BodyFields  []Field `json:"body_fields"`           // Response body fields
    BodyObject  *string `json:"body_object,omitempty"` // Response object name
}
```

### Functions

#### ApplyOverlay
Applies overlay to generate objects and CRUD endpoints from resources.

```go
func ApplyOverlay(input *Service) *Service
```

**What it generates:**
- Objects for resources with Read operations
- Create endpoints for resources with Create operations  
- Update endpoints for resources with Update operations
- Delete endpoints for resources with Delete operations
- Get endpoints for resources with Read operations
- List endpoints for resources with Read operations  
- Search endpoints for resources with Read operations
- Standard error handling objects (Error, ErrorCode, etc.)
- Pagination objects for list operations
- Request validation objects for all endpoint body parameters

#### ApplyFilterOverlay
Applies filter overlay to generate comprehensive filter objects.

```go
func ApplyFilterOverlay(input *Service) *Service
```

**What it generates:**
- Main filter objects (e.g., `UsersFilter`)
- Equals filter objects (e.g., `UsersFilterEquals`)
- Range filter objects (e.g., `UsersFilterRange`) 
- Contains filter objects (e.g., `UsersFilterContains`)
- LIKE filter objects (e.g., `UsersFilterLike`)
- Null filter objects (e.g., `UsersFilterNull`)

#### Parsing Functions
Parse specifications from files and data.

```go
func ParseServiceFromFile(filePath string) (*Service, error)
func ParseServiceFromBytes(data []byte, fileExtension string) (*Service, error)
func ParseServiceFromJSON(data []byte) (*Service, error)
func ParseServiceFromYAML(data []byte) (*Service, error)
```

**Note**: All parsing functions automatically apply overlays for complete specifications.

#### Validation Functions
Validate specifications with detailed error reporting.

```go
func ValidateServiceWithPosition(data []byte, fileExtension string) error
```

**Features:**
- YAML/JSON format support
- Line and column error reporting
- Operation validation (PascalCase required)
- Field type validation
- Modifier validation

### Constants

#### Operations
```go
const (
    OperationCreate = "Create"
    OperationRead   = "Read"
    OperationUpdate = "Update"
    OperationDelete = "Delete"
)
```

#### Field Types
```go
const (
    FieldTypeUUID      = "UUID"
    FieldTypeDate      = "Date"
    FieldTypeTimestamp = "Timestamp"
    FieldTypeString    = "String"
    FieldTypeInt       = "Int"
    FieldTypeBool      = "Bool"
)
```

#### Modifiers
```go
const (
    ModifierNullable = "Nullable"
    ModifierArray    = "Array"
)
```

---

## Package: specification/schemagen

JSON schema generation and validation package.

### Types

#### SchemaGenerator
Main type for generating and validating JSON schemas.

```go
type SchemaGenerator struct {
    // Private fields
}
```

### Functions

#### NewSchemaGenerator
Creates a new schema generator with default configuration.

```go
func NewSchemaGenerator() *SchemaGenerator
```

#### Schema Generation Methods
Generate JSON schemas for specification types.

```go
func (sg *SchemaGenerator) GenerateServiceSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateEnumSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateObjectSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateResourceSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateFieldSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateResourceFieldSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateEndpointSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateEndpointRequestSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateEndpointResponseSchema() (*jsonschema.Schema, error)
func (sg *SchemaGenerator) GenerateAllSchemas() (map[string]*jsonschema.Schema, error)
```

#### Schema Utility Methods
Convert schemas to JSON and manage validation.

```go
func (sg *SchemaGenerator) SchemaToJSON(schema *jsonschema.Schema) (string, error)
```

#### Validation Methods  
Validate JSON/YAML data against generated schemas.

```go
func (sg *SchemaGenerator) ValidateService(data []byte) error
func (sg *SchemaGenerator) ValidateEnum(data []byte) error
func (sg *SchemaGenerator) ValidateObject(data []byte) error
func (sg *SchemaGenerator) ValidateResource(data []byte) error
func (sg *SchemaGenerator) ValidateField(data []byte) error
func (sg *SchemaGenerator) ValidateResourceField(data []byte) error
func (sg *SchemaGenerator) ValidateEndpoint(data []byte) error
func (sg *SchemaGenerator) ValidateEndpointRequest(data []byte) error
func (sg *SchemaGenerator) ValidateEndpointResponse(data []byte) error
```

#### Parsing Methods
Parse and validate data into Go structs.

```go
func (sg *SchemaGenerator) ParseServiceFromJSON(data []byte) (*specification.Service, error)
func (sg *SchemaGenerator) ParseServiceFromYAML(data []byte) (*specification.Service, error)
func (sg *SchemaGenerator) ParseEnumFromJSON(data []byte) (*specification.Enum, error)
func (sg *SchemaGenerator) ParseEnumFromYAML(data []byte) (*specification.Enum, error)
func (sg *SchemaGenerator) ParseObjectFromJSON(data []byte) (*specification.Object, error)
func (sg *SchemaGenerator) ParseObjectFromYAML(data []byte) (*specification.Object, error)
func (sg *SchemaGenerator) ParseResourceFromJSON(data []byte) (*specification.Resource, error)
func (sg *SchemaGenerator) ParseResourceFromYAML(data []byte) (*specification.Resource, error)
```

---

## Package: specification/openapigen

OpenAPI 3.1 generation package.

### Types

#### Generator
Main type for generating OpenAPI 3.1 specifications.

```go
type Generator struct {
    // Private fields
}
```

### Functions

#### NewGenerator
Creates a new OpenAPI generator.

```go
func NewGenerator() *Generator
```

#### Generation Methods
Generate OpenAPI documents from specifications.

```go
func (g *Generator) GenerateFromService(service *specification.Service) (*v3.Document, error)
```

#### Output Methods
Convert OpenAPI documents to different formats.

```go
func (g *Generator) ToYAML(document *v3.Document) ([]byte, error)
func (g *Generator) ToJSON(document *v3.Document) ([]byte, error)
```

#### Customization Methods
Add custom information to OpenAPI documents.

```go
func (g *Generator) SetContactInfo(name, email, url string)
func (g *Generator) SetLicenseInfo(name, url string)
func (g *Generator) AddTag(name, description string)
func (g *Generator) AddSecurityScheme(name string, scheme SecurityScheme)
func (g *Generator) AddOperationExample(operationId, mediaType string, example interface{})
```

#### Validation Methods
Validate generated OpenAPI documents.

```go
func (g *Generator) ValidateDocument(document *v3.Document) (bool, []string)
func (g *Generator) GetDocumentStats(document *v3.Document) DocumentStats
```

### Types

#### SecurityScheme
Security scheme definition for OpenAPI.

```go
type SecurityScheme struct {
    Type string `json:"type"`
    Name string `json:"name,omitempty"`
    In   string `json:"in,omitempty"`
}
```

#### DocumentStats
Statistics about generated OpenAPI document.

```go
type DocumentStats struct {
    PathCount      int
    OperationCount int
    SchemaCount    int
    ParameterCount int
}
```

---

## Usage Patterns

### Basic Usage Pattern
```go
// 1. Parse specification
service, err := specification.ParseServiceFromFile("api.yaml")
if err != nil {
    log.Fatal(err)
}

// 2. Generate OpenAPI
generator := openapi.NewGenerator()
document, err := generator.GenerateFromService(service)
if err != nil {
    log.Fatal(err)
}

// 3. Save documentation
yamlBytes, _ := generator.ToYAML(document)
os.WriteFile("openapi.yaml", yamlBytes, 0644)
```

### Validation Pattern
```go
// 1. Create schema generator
generator := schema.NewSchemaGenerator()

// 2. Validate specification file
err := generator.ValidateService(data)
if err != nil {
    log.Printf("Validation failed: %v", err)
    return
}

// 3. Parse if valid
service, err := generator.ParseServiceFromYAML(data)
if err != nil {
    log.Fatal(err)
}
```

### Filter Generation Pattern
```go
// 1. Parse base specification
service, _ := specification.ParseServiceFromFile("api.yaml")

// 2. Apply overlays to generate complete specification
overlayedService := specification.ApplyOverlay(service)
completeService := specification.ApplyFilterOverlay(overlayedService)

// 3. Use generated filter objects
for _, obj := range completeService.Objects {
    if strings.HasSuffix(obj.Name, "Filter") {
        fmt.Printf("Filter: %s with %d fields\n", obj.Name, len(obj.Fields))
    }
}
```

## Error Handling

### Validation Errors
All validation functions return detailed errors with context:

```go
type ValidationError struct {
    Message string
    Line    int    
    Column  int
    Path    string
}
```

### Common Validation Errors
- `invalid operation: operation 'create' must be one of: [Create Read Update Delete]`
- `invalid field type: field type 'string' must be one of the primitive types`
- `invalid modifier: modifier 'nullable' must be one of: [Nullable Array]`

### Parsing Errors
- `file does not exist: <filepath>`
- `unsupported file format: file must have .yaml, .yml, or .json extension`
- `failed to read file: <error details>`
- `failed to parse file: YAML/JSON parsing error`

## File Format Support

### Supported Extensions
- `.yaml` - YAML format
- `.yml` - YAML format  
- `.json` - JSON format

### Example File Formats

**YAML format:**
```yaml
name: "My API"
version: "1.0.0"
resources:
  - name: "Users"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "username"
          type: "String"
        operations: ["Create", "Read"]
```

**JSON format:**
```json
{
  "name": "My API",
  "version": "1.0.0",
  "resources": [
    {
      "name": "Users",
      "operations": ["Create", "Read", "Update", "Delete"],
      "fields": [
        {
          "field": {
            "name": "username",
            "type": "String"
          },
          "operations": ["Create", "Read"]
        }
      ]
    }
  ]
}
```

## Performance Considerations

### Schema Generation
- Schemas are generated using reflection - cache results for repeated use
- Schema generation is fast but can be optimized by pre-generating schemas

### Validation
- JSON schema validation has overhead - validate once, cache results
- Use streaming validation for large files
- Consider async validation for multiple files

### Overlay Application  
- Overlay functions create deep copies - avoid repeated applications
- Apply overlays once and cache the result
- Filter overlay depends on base overlay - apply in correct order

## Related Documentation

- [🚀 Getting Started Guide](getting-started.md) - Step-by-step tutorial
- [⚙️ Working with Specifications](specifications.md) - Creating and structuring specs
- [🔍 Advanced Filtering](filtering.md) - Using generated filter objects
- [📋 OpenAPI Generation](openapi.md) - Creating documentation
- [✅ JSON Schema Validation](schema-validation.md) - Validating specifications