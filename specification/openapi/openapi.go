package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
	"gopkg.in/yaml.v3"
)

// Error constants
const (
	errorNotImplemented    = "not implemented"
	errorInvalidService    = "invalid service: service cannot be nil"
	errorInvalidDocument   = "invalid document: document cannot be nil"
	errorFailedToMarshal   = "failed to marshal specification"
	errorFailedToUnmarshal = "failed to unmarshal specification"
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

	// lastLibDocument stores the last generated libopenapi Document for serialization
	lastLibDocument libopenapi.Document
}

// NewGenerator creates a new OpenAPI generator with default settings.
func NewGenerator() *Generator {
	return &Generator{
		Version: "3.1.0",
	}
}

// OpenAPIDocument represents a simplified OpenAPI document structure for JSON marshaling
type OpenAPIDocument struct {
	OpenAPI    string                     `json:"openapi"`
	Info       OpenAPIInfo                `json:"info"`
	Servers    []OpenAPIServer            `json:"servers,omitempty"`
	Paths      map[string]OpenAPIPathItem `json:"paths"`
	Components OpenAPIComponents          `json:"components,omitempty"`
}

type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type OpenAPIComponents struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}

type OpenAPIPathItem struct {
	Get    *OpenAPIOperation `json:"get,omitempty"`
	Post   *OpenAPIOperation `json:"post,omitempty"`
	Put    *OpenAPIOperation `json:"put,omitempty"`
	Patch  *OpenAPIOperation `json:"patch,omitempty"`
	Delete *OpenAPIOperation `json:"delete,omitempty"`
}

type OpenAPIOperation struct {
	OperationID string                     `json:"operationId,omitempty"`
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Parameters  []OpenAPIParameter         `json:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse `json:"responses"`
}

type OpenAPIParameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      interface{} `json:"schema"`
}

type OpenAPIRequestBody struct {
	Description string                      `json:"description,omitempty"`
	Required    bool                        `json:"required,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content"`
}

type OpenAPIResponse struct {
	Description string                      `json:"description"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty"`
}

type OpenAPIMediaType struct {
	Schema interface{} `json:"schema"`
}

// GenerateFromService generates an OpenAPI 3.1 document from a specification.Service.
func (g *Generator) GenerateFromService(service *specification.Service) (*v3.Document, error) {
	if service == nil {
		return nil, errors.New(errorInvalidService)
	}

	// Build the OpenAPI specification as a plain Go map
	openAPISpec := g.buildOpenAPIDocument(service)

	// Marshal to JSON bytes
	specBytes, err := json.Marshal(openAPISpec)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	// Create a libopenapi Document from the bytes
	libDoc, err := libopenapi.NewDocument(specBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create libopenapi document: %w", err)
	}

	// Store the libopenapi Document for later serialization
	g.lastLibDocument = libDoc

	// Build the high-level v3 model
	model, errs := libDoc.BuildV3Model()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to build v3 model: %v", errs)
	}

	return &model.Model, nil
}

// buildOpenAPIDocument creates the actual OpenAPI document structure
func (g *Generator) buildOpenAPIDocument(service *specification.Service) *OpenAPIDocument {
	// Create document title
	title := g.Title
	if title == "" {
		title = service.Name
	}

	doc := &OpenAPIDocument{
		OpenAPI: g.Version,
		Info: OpenAPIInfo{
			Title:       title,
			Description: g.Description,
			Version:     "1.0.0",
		},
		Paths:      make(map[string]OpenAPIPathItem),
		Components: OpenAPIComponents{Schemas: make(map[string]interface{})},
	}

	// Add server if specified
	if g.ServerURL != "" {
		doc.Servers = []OpenAPIServer{
			{
				URL:         g.ServerURL,
				Description: fmt.Sprintf("%s server", title),
			},
		}
	}

	// Convert enums to schemas
	for _, enum := range service.Enums {
		doc.Components.Schemas[enum.Name] = g.buildEnumSchema(enum)
	}

	// Convert objects to schemas
	for _, obj := range service.Objects {
		doc.Components.Schemas[obj.Name] = g.buildObjectSchema(obj, service)
	}

	// Convert resources to paths
	for _, resource := range service.Resources {
		g.addResourceToPaths(doc, resource, service)
	}

	return doc
}

// buildEnumSchema creates an OpenAPI schema for an enum
func (g *Generator) buildEnumSchema(enum specification.Enum) map[string]interface{} {
	enumValues := make([]string, len(enum.Values))
	for i, value := range enum.Values {
		enumValues[i] = value.Name
	}

	return map[string]interface{}{
		"type":        "string",
		"description": enum.Description,
		"enum":        enumValues,
	}
}

// buildObjectSchema creates an OpenAPI schema for an object
func (g *Generator) buildObjectSchema(obj specification.Object, service *specification.Service) map[string]interface{} {
	schema := map[string]interface{}{
		"type":        "object",
		"description": obj.Description,
		"properties":  make(map[string]interface{}),
	}

	properties := make(map[string]interface{})
	required := []string{}

	for _, field := range obj.Fields {
		properties[field.Name] = g.buildFieldSchema(field, service)
		if field.IsRequired(service) {
			required = append(required, field.Name)
		}
	}

	schema["properties"] = properties
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// buildFieldSchema creates an OpenAPI schema for a field
func (g *Generator) buildFieldSchema(field specification.Field, service *specification.Service) map[string]interface{} {
	var schema map[string]interface{}

	// Handle array modifier
	if field.IsArray() {
		schema = map[string]interface{}{
			"type":  "array",
			"items": g.buildTypeSchema(field.Type, service),
		}
	} else {
		schema = g.buildTypeSchema(field.Type, service)
	}

	// Add description if not already set
	if _, exists := schema["description"]; !exists && field.Description != "" {
		schema["description"] = field.Description
	}

	// Handle nullable modifier
	if field.IsNullable() {
		if fieldType, exists := schema["type"]; exists {
			if typeStr, ok := fieldType.(string); ok {
				schema["type"] = []string{typeStr, "null"}
			}
		}
	}

	// Add default and example if present
	if field.Default != "" {
		schema["default"] = field.Default
	}
	if field.Example != "" {
		schema["example"] = field.Example
	}

	return schema
}

// buildTypeSchema creates an OpenAPI schema for a specific type
func (g *Generator) buildTypeSchema(fieldType string, service *specification.Service) map[string]interface{} {
	switch fieldType {
	case specification.FieldTypeString:
		return map[string]interface{}{"type": "string"}
	case specification.FieldTypeInt:
		return map[string]interface{}{"type": "integer"}
	case specification.FieldTypeBool:
		return map[string]interface{}{"type": "boolean"}
	case specification.FieldTypeUUID:
		return map[string]interface{}{
			"type":   "string",
			"format": "uuid",
		}
	case specification.FieldTypeDate:
		return map[string]interface{}{
			"type":   "string",
			"format": "date",
		}
	case specification.FieldTypeTimestamp:
		return map[string]interface{}{
			"type":   "string",
			"format": "date-time",
		}
	default:
		// Check if it's a custom object or enum
		if service.HasObject(fieldType) || service.HasEnum(fieldType) {
			return map[string]interface{}{
				"$ref": fmt.Sprintf("#/components/schemas/%s", fieldType),
			}
		}
		// Default to string if unknown type
		return map[string]interface{}{"type": "string"}
	}
}

// addResourceToPaths adds resource endpoints to the OpenAPI paths
func (g *Generator) addResourceToPaths(doc *OpenAPIDocument, resource specification.Resource, service *specification.Service) {
	basePath := "/" + strings.ToLower(resource.Name)

	// Group endpoints by path
	pathGroups := make(map[string][]*specification.Endpoint)
	for _, endpoint := range resource.Endpoints {
		fullPath := basePath + endpoint.Path
		pathGroups[fullPath] = append(pathGroups[fullPath], &endpoint)
	}

	// Create PathItem for each unique path
	for path, endpoints := range pathGroups {
		pathItem := OpenAPIPathItem{}

		for _, endpoint := range endpoints {
			operation := g.buildOperation(*endpoint, resource, service)

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

		doc.Paths[path] = pathItem
	}
}

// buildOperation creates an OpenAPI operation from an endpoint
func (g *Generator) buildOperation(endpoint specification.Endpoint, resource specification.Resource, service *specification.Service) *OpenAPIOperation {
	operation := &OpenAPIOperation{
		OperationID: endpoint.Name,
		Summary:     endpoint.Title,
		Description: endpoint.Description,
		Tags:        []string{resource.Name},
		Responses:   make(map[string]OpenAPIResponse),
	}

	// Add parameters
	var parameters []OpenAPIParameter

	// Path parameters
	for _, param := range endpoint.Request.PathParams {
		parameters = append(parameters, OpenAPIParameter{
			Name:        param.Name,
			In:          "path",
			Description: param.Description,
			Required:    param.IsRequired(service),
			Schema:      g.buildFieldSchema(param, service),
		})
	}

	// Query parameters
	for _, param := range endpoint.Request.QueryParams {
		parameters = append(parameters, OpenAPIParameter{
			Name:        param.Name,
			In:          "query",
			Description: param.Description,
			Required:    param.IsRequired(service),
			Schema:      g.buildFieldSchema(param, service),
		})
	}

	operation.Parameters = parameters

	// Request body
	if len(endpoint.Request.BodyParams) > 0 {
		operation.RequestBody = g.buildRequestBody(endpoint.Request.BodyParams, service)
	}

	// Response
	operation.Responses[strconv.Itoa(endpoint.Response.StatusCode)] = g.buildResponse(endpoint.Response, service)

	// Add common error responses
	g.addErrorResponses(operation, service)

	return operation
}

// buildRequestBody creates an OpenAPI request body
func (g *Generator) buildRequestBody(bodyParams []specification.Field, service *specification.Service) *OpenAPIRequestBody {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	properties := make(map[string]interface{})
	required := []string{}

	for _, field := range bodyParams {
		properties[field.Name] = g.buildFieldSchema(field, service)
		if field.IsRequired(service) {
			required = append(required, field.Name)
		}
	}

	schema["properties"] = properties
	if len(required) > 0 {
		schema["required"] = required
	}

	return &OpenAPIRequestBody{
		Description: "Request body",
		Required:    len(required) > 0,
		Content: map[string]OpenAPIMediaType{
			"application/json": {Schema: schema},
		},
	}
}

// buildResponse creates an OpenAPI response
func (g *Generator) buildResponse(response specification.EndpointResponse, service *specification.Service) OpenAPIResponse {
	openAPIResponse := OpenAPIResponse{
		Description: "Successful response",
	}

	if response.BodyObject != nil || len(response.BodyFields) > 0 {
		var schema interface{}

		if response.BodyObject != nil {
			schema = map[string]interface{}{
				"$ref": fmt.Sprintf("#/components/schemas/%s", *response.BodyObject),
			}
		} else if len(response.BodyFields) > 0 {
			schemaMap := map[string]interface{}{
				"type":       "object",
				"properties": make(map[string]interface{}),
			}

			properties := make(map[string]interface{})
			for _, field := range response.BodyFields {
				properties[field.Name] = g.buildFieldSchema(field, service)
			}
			schemaMap["properties"] = properties
			schema = schemaMap
		}

		openAPIResponse.Content = map[string]OpenAPIMediaType{
			"application/json": {Schema: schema},
		}
	}

	return openAPIResponse
}

// addErrorResponses adds common error responses to the operation
func (g *Generator) addErrorResponses(operation *OpenAPIOperation, service *specification.Service) {
	var errorSchema interface{}
	if service.HasObject("Error") {
		errorSchema = map[string]interface{}{
			"$ref": "#/components/schemas/Error",
		}
	} else {
		errorSchema = map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{"type": "string"},
				"code":    map[string]interface{}{"type": "string"},
			},
			"required": []string{"message", "code"},
		}
	}

	errorContent := map[string]OpenAPIMediaType{
		"application/json": {Schema: errorSchema},
	}

	operation.Responses["400"] = OpenAPIResponse{
		Description: "Bad Request - The request was malformed or contained invalid parameters",
		Content:     errorContent,
	}
	operation.Responses["401"] = OpenAPIResponse{
		Description: "Unauthorized - The request is missing valid authentication credentials",
		Content:     errorContent,
	}
	operation.Responses["404"] = OpenAPIResponse{
		Description: "Not Found - The requested resource does not exist",
		Content:     errorContent,
	}
	operation.Responses["500"] = OpenAPIResponse{
		Description: "Internal Server Error - An unexpected server error occurred",
		Content:     errorContent,
	}
}

// ToYAML converts an OpenAPI document to YAML format.
func (g *Generator) ToYAML(document *v3.Document) ([]byte, error) {
	if document == nil {
		return nil, errors.New(errorInvalidDocument)
	}

	if g.lastLibDocument == nil {
		return nil, errors.New("no document has been generated")
	}

	// Use libopenapi's Render method to serialize the document
	yamlBytes, err := g.lastLibDocument.Render()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	return yamlBytes, nil
}

// ToJSON converts an OpenAPI document to JSON format.
func (g *Generator) ToJSON(document *v3.Document) ([]byte, error) {
	if document == nil {
		return nil, errors.New(errorInvalidDocument)
	}

	if g.lastLibDocument == nil {
		return nil, errors.New("no document has been generated")
	}

	// Use libopenapi's Render method to get YAML, then convert to JSON
	yamlBytes, err := g.lastLibDocument.Render()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	// Parse YAML and convert to JSON with proper formatting
	var yamlData interface{}
	err = yaml.Unmarshal(yamlBytes, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFailedToUnmarshal, err)
	}

	jsonBytes, err := json.MarshalIndent(yamlData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	return jsonBytes, nil
}
