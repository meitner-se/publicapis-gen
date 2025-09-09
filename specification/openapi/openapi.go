package openapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
)

// Error constants
const (
	errorInvalidService  = "invalid service: service cannot be nil"
	errorInvalidDocument = "invalid document: document cannot be nil"
)

// Generator handles OpenAPI 3.1 specification generation from specification.Service.
type Generator struct {
	// Version specifies the OpenAPI version to generate (default: "3.1.0")
	Version string

	// Title specifies the API title (defaults to service name if not set)
	Title string

	// Description specifies the API description
	Description string

	// ServerURL specifies the base server URL for the API
	ServerURL string
}

// NewGenerator creates a new OpenAPI generator with default settings.
func NewGenerator() *Generator {
	return &Generator{
		Version: "3.1.0",
	}
}

// GenerateFromService generates an OpenAPI 3.1 document from a specification.Service.
func (g *Generator) GenerateFromService(service *specification.Service) (*v3.Document, error) {
	if service == nil {
		return nil, errors.New(errorInvalidService)
	}

	// Build document using native libopenapi v3 types
	return g.buildV3Document(service), nil
}

// buildV3Document creates a v3.Document using native libopenapi types.
func (g *Generator) buildV3Document(service *specification.Service) *v3.Document {
	// Create document title
	title := g.Title
	if title == "" {
		title = service.Name
	}

	// Create Info section
	version := service.Version
	if version == "" {
		version = "1.0.0" // Default version if not specified
	}

	info := &base.Info{
		Title:       title,
		Description: g.Description,
		Version:     version,
	}

	// Create Document
	document := &v3.Document{
		Version: g.Version,
		Info:    info,
	}

	// Add servers from service specification
	if len(service.Servers) > 0 {
		servers := make([]*v3.Server, len(service.Servers))
		for i, server := range service.Servers {
			servers[i] = &v3.Server{
				URL:         server.URL,
				Description: server.Description,
			}
		}
		document.Servers = servers
	} else if g.ServerURL != "" {
		// Fallback to generator's ServerURL for backwards compatibility
		servers := []*v3.Server{
			{
				URL:         g.ServerURL,
				Description: fmt.Sprintf("%s server", title),
			},
		}
		document.Servers = servers
	}

	// Create Components
	components := &v3.Components{
		Schemas: orderedmap.New[string, *base.SchemaProxy](),
	}

	// Add enums to components
	for _, enum := range service.Enums {
		schema := g.createEnumSchema(enum)
		proxy := base.CreateSchemaProxy(schema)
		components.Schemas.Set(enum.Name, proxy)
	}

	// Add objects to components
	for _, obj := range service.Objects {
		schema := g.createObjectSchema(obj, service)
		proxy := base.CreateSchemaProxy(schema)
		components.Schemas.Set(obj.Name, proxy)
	}

	document.Components = components

	// Create Paths
	paths := orderedmap.New[string, *v3.PathItem]()
	for _, resource := range service.Resources {
		g.addResourceToPaths(resource, paths, service)
	}
	document.Paths = &v3.Paths{
		PathItems: paths,
	}

	return document
}

// createEnumSchema creates a base.Schema for an enum using native types.
func (g *Generator) createEnumSchema(enum specification.Enum) *base.Schema {
	schema := &base.Schema{
		Type:        []string{"string"},
		Description: enum.Description,
	}

	// Add enum values
	enumValues := make([]*yaml.Node, len(enum.Values))
	for i, value := range enum.Values {
		node := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value.Name,
		}
		enumValues[i] = node
	}
	schema.Enum = enumValues

	return schema
}

// createObjectSchema creates a base.Schema for an object using native types.
func (g *Generator) createObjectSchema(obj specification.Object, service *specification.Service) *base.Schema {
	schema := &base.Schema{
		Type:        []string{"object"},
		Description: obj.Description,
		Properties:  orderedmap.New[string, *base.SchemaProxy](),
	}

	requiredFields := []string{}
	for _, field := range obj.Fields {
		fieldSchema := g.createFieldSchema(field, service)
		proxy := base.CreateSchemaProxy(fieldSchema)
		schema.Properties.Set(field.Name, proxy)

		if field.IsRequired(service) {
			requiredFields = append(requiredFields, field.Name)
		}
	}

	if len(requiredFields) > 0 {
		schema.Required = requiredFields
	}

	return schema
}

// createFieldSchema creates a base.Schema for a field using native types.
func (g *Generator) createFieldSchema(field specification.Field, service *specification.Service) *base.Schema {
	var schema *base.Schema

	// Handle array modifier
	if field.IsArray() {
		schema = &base.Schema{
			Type:        []string{"array"},
			Description: field.Description,
		}

		itemSchema := g.getTypeSchema(field.Type, service)
		schema.Items = &base.DynamicValue[*base.SchemaProxy, bool]{
			N: 0, // Single schema (not boolean)
			A: base.CreateSchemaProxy(itemSchema),
		}
	} else {
		schema = g.getTypeSchema(field.Type, service)
		schema.Description = field.Description
	}

	// Handle nullable modifier
	if field.IsNullable() {
		// Use the Nullable field instead of appending "null" to type array
		nullable := true
		schema.Nullable = &nullable
	}

	// Add default value if present
	if field.Default != "" {
		defaultNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: field.Default,
		}
		schema.Default = defaultNode
	}

	// Add example if present
	if field.Example != "" {
		exampleNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: field.Example,
		}
		schema.Examples = []*yaml.Node{exampleNode}
	}

	return schema
}

// getTypeSchema returns a base.Schema for the given field type.
func (g *Generator) getTypeSchema(fieldType string, service *specification.Service) *base.Schema {
	switch fieldType {
	case specification.FieldTypeString:
		return &base.Schema{Type: []string{"string"}}
	case specification.FieldTypeInt:
		return &base.Schema{Type: []string{"integer"}}
	case specification.FieldTypeBool:
		return &base.Schema{Type: []string{"boolean"}}
	case specification.FieldTypeUUID:
		return &base.Schema{
			Type:   []string{"string"},
			Format: "uuid",
		}
	case specification.FieldTypeDate:
		return &base.Schema{
			Type:   []string{"string"},
			Format: "date",
		}
	case specification.FieldTypeTimestamp:
		return &base.Schema{
			Type:   []string{"string"},
			Format: "date-time",
		}
	default:
		// Check if it's a custom object or enum
		if service.HasObject(fieldType) || service.HasEnum(fieldType) {
			// Create a reference schema
			return &base.Schema{
				Title: fieldType, // Temporary - proper $ref handling would need low-level API
			}
		}
		// Default to string if unknown type
		return &base.Schema{Type: []string{"string"}}
	}
}

// addResourceToPaths adds resource endpoints to paths using native v3 types.
func (g *Generator) addResourceToPaths(resource specification.Resource, paths *orderedmap.Map[string, *v3.PathItem], service *specification.Service) {
	basePath := "/" + strings.ToLower(resource.Name)

	// Group endpoints by path
	pathGroups := make(map[string][]*specification.Endpoint)
	for _, endpoint := range resource.Endpoints {
		fullPath := basePath + endpoint.Path
		pathGroups[fullPath] = append(pathGroups[fullPath], &endpoint)
	}

	// Create PathItem for each unique path
	for path, endpoints := range pathGroups {
		pathItem := &v3.PathItem{}

		for _, endpoint := range endpoints {
			operation := g.createOperation(*endpoint, resource, service)

			// Set operation based on HTTP method
			switch strings.ToUpper(endpoint.Method) {
			case http.MethodGet:
				pathItem.Get = operation
			case http.MethodPost:
				pathItem.Post = operation
			case http.MethodPatch:
				pathItem.Patch = operation
			case http.MethodPut:
				pathItem.Put = operation
			case http.MethodDelete:
				pathItem.Delete = operation
			}
		}

		paths.Set(path, pathItem)
	}
}

// createOperation creates a v3.Operation from an endpoint using native types.
func (g *Generator) createOperation(endpoint specification.Endpoint, resource specification.Resource, service *specification.Service) *v3.Operation {
	operation := &v3.Operation{
		OperationId: endpoint.Name,
		Summary:     endpoint.Title,
		Description: endpoint.Description,
		Tags:        []string{resource.Name},
	}

	// Add parameters
	parameters := []*v3.Parameter{}

	// Path parameters
	for _, param := range endpoint.Request.PathParams {
		parameters = append(parameters, g.createParameter(param, "path", service))
	}

	// Query parameters
	for _, param := range endpoint.Request.QueryParams {
		parameters = append(parameters, g.createParameter(param, "query", service))
	}

	operation.Parameters = parameters

	// Request body
	if len(endpoint.Request.BodyParams) > 0 {
		operation.RequestBody = g.createRequestBody(endpoint.Request.BodyParams, service)
	}

	// Responses
	responses := orderedmap.New[string, *v3.Response]()

	// Success response
	successResponse := g.createResponse(endpoint.Response, service)
	responses.Set(strconv.Itoa(endpoint.Response.StatusCode), successResponse)

	// Add error responses
	g.addErrorResponses(responses, service)

	operation.Responses = &v3.Responses{
		Codes: responses,
	}

	return operation
}

// createParameter creates a v3.Parameter from a field using native types.
func (g *Generator) createParameter(field specification.Field, location string, service *specification.Service) *v3.Parameter {
	isRequired := field.IsRequired(service)
	param := &v3.Parameter{
		Name:        field.Name,
		In:          location,
		Description: field.Description,
		Required:    &isRequired,
		Schema:      base.CreateSchemaProxy(g.createFieldSchema(field, service)),
	}

	return param
}

// createRequestBody creates a v3.RequestBody from body parameters using native types.
func (g *Generator) createRequestBody(bodyParams []specification.Field, service *specification.Service) *v3.RequestBody {
	// Create schema from body parameters
	schema := &base.Schema{
		Type:       []string{"object"},
		Properties: orderedmap.New[string, *base.SchemaProxy](),
	}

	requiredFields := []string{}
	for _, field := range bodyParams {
		fieldSchema := g.createFieldSchema(field, service)
		proxy := base.CreateSchemaProxy(fieldSchema)
		schema.Properties.Set(field.Name, proxy)

		if field.IsRequired(service) {
			requiredFields = append(requiredFields, field.Name)
		}
	}

	if len(requiredFields) > 0 {
		schema.Required = requiredFields
	}

	// Create media type
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(schema),
	}

	content := orderedmap.New[string, *v3.MediaType]()
	content.Set("application/json", mediaType)

	isRequired := len(requiredFields) > 0
	return &v3.RequestBody{
		Description: "Request body",
		Content:     content,
		Required:    &isRequired,
	}
}

// createResponse creates a v3.Response from an endpoint response using native types.
func (g *Generator) createResponse(response specification.EndpointResponse, service *specification.Service) *v3.Response {
	openAPIResponse := &v3.Response{
		Description: "Successful response",
	}

	// Add response content if present
	if response.BodyObject != nil || len(response.BodyFields) > 0 {
		content := orderedmap.New[string, *v3.MediaType]()

		var schema *base.Schema
		if response.BodyObject != nil {
			// Reference to existing schema
			schema = &base.Schema{
				Title: *response.BodyObject, // Temporary - proper $ref handling would need low-level API
			}
		} else if len(response.BodyFields) > 0 {
			// Inline schema from body fields
			schema = &base.Schema{
				Type:       []string{"object"},
				Properties: orderedmap.New[string, *base.SchemaProxy](),
			}

			for _, field := range response.BodyFields {
				fieldSchema := g.createFieldSchema(field, service)
				proxy := base.CreateSchemaProxy(fieldSchema)
				schema.Properties.Set(field.Name, proxy)
			}
		}

		if schema != nil {
			mediaType := &v3.MediaType{
				Schema: base.CreateSchemaProxy(schema),
			}
			content.Set("application/json", mediaType)
			openAPIResponse.Content = content
		}
	}

	return openAPIResponse
}

// addErrorResponses adds common error responses using native types.
func (g *Generator) addErrorResponses(responses *orderedmap.Map[string, *v3.Response], service *specification.Service) {
	// Create error schema
	var errorSchema *base.Schema
	if service.HasObject("Error") {
		errorSchema = &base.Schema{
			Title: "Error", // Temporary - proper $ref handling would need low-level API
		}
	} else {
		// Fallback generic error schema
		errorSchema = &base.Schema{
			Type:       []string{"object"},
			Properties: orderedmap.New[string, *base.SchemaProxy](),
		}
		messageSchema := &base.Schema{Type: []string{"string"}}
		codeSchema := &base.Schema{Type: []string{"string"}}
		errorSchema.Properties.Set("message", base.CreateSchemaProxy(messageSchema))
		errorSchema.Properties.Set("code", base.CreateSchemaProxy(codeSchema))
		errorSchema.Required = []string{"message", "code"}
	}

	errorContent := orderedmap.New[string, *v3.MediaType]()
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(errorSchema),
	}
	errorContent.Set("application/json", mediaType)

	// 400 Bad Request
	badRequestResponse := &v3.Response{
		Description: "Bad Request - The request was malformed or contained invalid parameters",
		Content:     errorContent,
	}
	responses.Set("400", badRequestResponse)

	// 401 Unauthorized
	unauthorizedResponse := &v3.Response{
		Description: "Unauthorized - The request is missing valid authentication credentials",
		Content:     errorContent,
	}
	responses.Set("401", unauthorizedResponse)

	// 404 Not Found
	notFoundResponse := &v3.Response{
		Description: "Not Found - The requested resource does not exist",
		Content:     errorContent,
	}
	responses.Set("404", notFoundResponse)

	// 500 Internal Server Error
	internalErrorResponse := &v3.Response{
		Description: "Internal Server Error - An unexpected server error occurred",
		Content:     errorContent,
	}
	responses.Set("500", internalErrorResponse)
}

// ToYAML converts an OpenAPI document to YAML format.
func (g *Generator) ToYAML(document *v3.Document) ([]byte, error) {
	if document == nil {
		return nil, errors.New(errorInvalidDocument)
	}

	// Use libopenapi's native Render method
	return document.Render()
}

// ToJSON converts an OpenAPI document to JSON format.
func (g *Generator) ToJSON(document *v3.Document) ([]byte, error) {
	if document == nil {
		return nil, errors.New(errorInvalidDocument)
	}

	// Use libopenapi's native RenderJSON method
	return document.RenderJSON("  ")
}
