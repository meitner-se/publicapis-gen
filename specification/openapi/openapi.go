package openapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
)

// Error constants
const (
	errorInvalidService  = "invalid service: service cannot be nil"
	errorInvalidDocument = "invalid document: document cannot be nil"
)

// HTTP Status Code constants
const (
	httpStatus400 = "400"
	httpStatus401 = "401"
	httpStatus403 = "403"
	httpStatus404 = "404"
	httpStatus409 = "409"
	httpStatus422 = "422"
	httpStatus429 = "429"
	httpStatus500 = "500"
)

// OpenAPI version constants
const (
	defaultOpenAPIVersion = "3.1.0"
	defaultServiceVersion = "1.0.0"
)

// API generation constants
const (
	apiTitleSuffix        = " API"
	defaultAPIDescription = "Generated API documentation"
)

// Content type constants
const (
	contentTypeJSON = "application/json"
)

// Standard response descriptions
const (
	requestBodyDescription = "Request body"
	successDescription     = "Successful response"
)

// Schema reference format constants
const (
	schemaReferencePrefix       = "#/components/schemas/"
	responseBodyReferencePrefix = "#/components/responses/"
)

// Server description template
const (
	serverDescriptionTemplate = "%s server"
)

// Error response descriptions
const (
	badRequestDescription    = "Bad Request - The request was malformed or contained invalid parameters"
	unauthorizedDescription  = "Unauthorized - The request is missing valid authentication credentials"
	notFoundDescription      = "Not Found - The requested resource does not exist"
	internalErrorDescription = "Internal Server Error - An unexpected server error occurred"
)

// Error code names
const (
	errorCodeBadRequest          = "BadRequest"
	errorCodeUnauthorized        = "Unauthorized"
	errorCodeForbidden           = "Forbidden"
	errorCodeNotFound            = "NotFound"
	errorCodeConflict            = "Conflict"
	errorCodeUnprocessableEntity = "UnprocessableEntity"
	errorCodeRateLimited         = "RateLimited"
	errorCodeInternal            = "Internal"
)

// Schema types
const (
	schemaTypeString  = "string"
	schemaTypeInteger = "integer"
	schemaTypeBoolean = "boolean"
	schemaTypeArray   = "array"
	schemaTypeObject  = "object"
)

// Schema formats
const (
	schemaFormatUUID     = "uuid"
	schemaFormatDate     = "date"
	schemaFormatDateTime = "date-time"
)

// Speakeasy retry configuration constants
const (
	speakeasyRetriesExtension = "x-speakeasy-retries"
)

// Speakeasy timeout configuration constants
const (
	speakeasyTimeoutExtension = "x-speakeasy-timeout"
	defaultTimeoutMs          = 30000 // 30 seconds in milliseconds
)

// Retry configuration field names
const (
	retryFieldStrategy              = "strategy"
	retryFieldBackoff               = "backoff"
	retryFieldInitialInterval       = "initialInterval"
	retryFieldMaxInterval           = "maxInterval"
	retryFieldMaxElapsedTime        = "maxElapsedTime"
	retryFieldExponent              = "exponent"
	retryFieldStatusCodes           = "statusCodes"
	retryFieldRetryConnectionErrors = "retryConnectionErrors"
)

// Speakeasy pagination configuration constants
const (
	speakeasyPaginationExtension = "x-speakeasy-pagination"
	speakeasyPaginationStrategy  = "offsetLimit"
	speakeasyOffsetParamName     = "offset"
	speakeasyLimitParamName      = "limit"
	speakeasyTotalFieldPath      = "pagination.total"
	speakeasyDataFieldName       = "data"
)

// Object and field names
const (
	errorObjectName        = "Error"
	messageFieldName       = "message"
	codeFieldName          = "code"
	errorCodeEnumName      = "ErrorCode"
	requestBodySuffix      = "RequestBody"
	responseBodySuffix     = "ResponseBody"
	errorResponseBodyPrefix = "Error"
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

// newGenerator creates a new OpenAPI generator with default settings.
func newGenerator() *Generator {
	return &Generator{
		Version: defaultOpenAPIVersion,
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

// addSpeakeasyRetryExtension adds Speakeasy retry configuration extension to the OpenAPI document.
func (g *Generator) addSpeakeasyRetryExtension(document *v3.Document, service *specification.Service) {
	// Initialize extensions map if it doesn't exist
	if document.Extensions == nil {
		document.Extensions = orderedmap.New[string, *yaml.Node]()
	}

	// Get retry configuration from service with defaults applied
	retryConfig := service.GetRetryConfigurationWithDefaults()

	// Convert the retry configuration to a YAML node
	retryNode := &yaml.Node{
		Kind: yaml.MappingNode,
	}

	// Create nodes for the retry configuration
	strategyKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldStrategy}
	strategyValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryConfig.Strategy}

	backoffKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldBackoff}
	backoffValueNode := &yaml.Node{Kind: yaml.MappingNode}

	// Backoff sub-nodes
	initialIntervalKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldInitialInterval}
	initialIntervalValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: strconv.Itoa(retryConfig.Backoff.InitialInterval)}
	maxIntervalKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldMaxInterval}
	maxIntervalValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: strconv.Itoa(retryConfig.Backoff.MaxInterval)}
	maxElapsedTimeKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldMaxElapsedTime}
	maxElapsedTimeValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: strconv.Itoa(retryConfig.Backoff.MaxElapsedTime)}
	exponentKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldExponent}
	exponentValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%g", retryConfig.Backoff.Exponent)}

	backoffValueNode.Content = []*yaml.Node{
		initialIntervalKeyNode, initialIntervalValueNode,
		maxIntervalKeyNode, maxIntervalValueNode,
		maxElapsedTimeKeyNode, maxElapsedTimeValueNode,
		exponentKeyNode, exponentValueNode,
	}

	statusCodesKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldStatusCodes}
	statusCodesValueNode := &yaml.Node{Kind: yaml.SequenceNode}
	statusCodesNodes := make([]*yaml.Node, len(retryConfig.StatusCodes))
	for i, statusCode := range retryConfig.StatusCodes {
		statusCodesNodes[i] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: statusCode,
			Tag:   "!!str", // Ensure this is treated as a string
		}
	}
	statusCodesValueNode.Content = statusCodesNodes

	connectionErrorsKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: retryFieldRetryConnectionErrors}
	connectionErrorsValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: strconv.FormatBool(retryConfig.RetryConnectionErrors)}

	// Assemble the main retry configuration node
	retryNode.Content = []*yaml.Node{
		strategyKeyNode, strategyValueNode,
		backoffKeyNode, backoffValueNode,
		statusCodesKeyNode, statusCodesValueNode,
		connectionErrorsKeyNode, connectionErrorsValueNode,
	}

	// Add the extension to the document
	document.Extensions.Set(speakeasyRetriesExtension, retryNode)
}

// addSpeakeasyTimeoutExtension adds Speakeasy timeout configuration extension to the OpenAPI document.
func (g *Generator) addSpeakeasyTimeoutExtension(document *v3.Document, service *specification.Service) {
	// Initialize extensions map if it doesn't exist
	if document.Extensions == nil {
		document.Extensions = orderedmap.New[string, *yaml.Node]()
	}

	// Determine timeout value from service configuration or use default
	timeoutMs := defaultTimeoutMs
	if service.Timeout != nil && service.Timeout.Timeout > 0 {
		timeoutMs = service.Timeout.Timeout
	}

	// Create a YAML node for the timeout value (in milliseconds)
	timeoutNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: strconv.Itoa(timeoutMs),
	}

	// Add the extension to the document
	document.Extensions.Set(speakeasyTimeoutExtension, timeoutNode)
}

// addSpeakeasyPaginationExtension adds Speakeasy pagination configuration extension to an operation.
func (g *Generator) addSpeakeasyPaginationExtension(operation *v3.Operation) {
	// Initialize extensions map if it doesn't exist
	if operation.Extensions == nil {
		operation.Extensions = orderedmap.New[string, *yaml.Node]()
	}

	// Convert the pagination configuration to a YAML node
	paginationNode := &yaml.Node{
		Kind: yaml.MappingNode,
	}

	// Create nodes for the pagination configuration
	strategyKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "strategy"}
	strategyValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: speakeasyPaginationStrategy}

	offsetParamKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "offsetParam"}
	offsetParamValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: speakeasyOffsetParamName}

	limitParamKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "limitParam"}
	limitParamValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: speakeasyLimitParamName}

	totalFieldKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "totalField"}
	totalFieldValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: speakeasyTotalFieldPath}

	dataFieldKeyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "dataField"}
	dataFieldValueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: speakeasyDataFieldName}

	// Assemble the main pagination configuration node
	paginationNode.Content = []*yaml.Node{
		strategyKeyNode, strategyValueNode,
		offsetParamKeyNode, offsetParamValueNode,
		limitParamKeyNode, limitParamValueNode,
		totalFieldKeyNode, totalFieldValueNode,
		dataFieldKeyNode, dataFieldValueNode,
	}

	// Add the extension to the operation
	operation.Extensions.Set(speakeasyPaginationExtension, paginationNode)
}

// isPaginatedOperation determines if an endpoint represents a paginated operation
// by checking for the presence of limit/offset query parameters and data/pagination response fields.
func (g *Generator) isPaginatedOperation(endpoint specification.Endpoint) bool {
	hasLimitParam := false
	hasOffsetParam := false

	// Check query parameters for limit and offset
	for _, param := range endpoint.Request.QueryParams {
		if param.Name == speakeasyLimitParamName {
			hasLimitParam = true
		}
		if param.Name == speakeasyOffsetParamName {
			hasOffsetParam = true
		}
	}

	// Must have both limit and offset parameters
	if !hasLimitParam || !hasOffsetParam {
		return false
	}

	hasDataField := false
	hasPaginationField := false

	// Check response body fields for data and pagination
	for _, field := range endpoint.Response.BodyFields {
		if field.Name == speakeasyDataFieldName {
			hasDataField = true
		}
		if field.Name == "Pagination" || field.Name == "pagination" {
			hasPaginationField = true
		}
	}

	// Must have both data and pagination fields in response
	return hasDataField && hasPaginationField
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
		version = defaultServiceVersion // Default version if not specified
	}

	info := &base.Info{
		Title:       title,
		Description: g.Description,
		Version:     version,
	}

	// Add contact information if available in the service
	if service.Contact != nil {
		contact := &base.Contact{}
		if service.Contact.Name != "" {
			contact.Name = service.Contact.Name
		}
		if service.Contact.URL != "" {
			contact.URL = service.Contact.URL
		}
		if service.Contact.Email != "" {
			contact.Email = service.Contact.Email
		}
		// Only set contact if at least one field is provided
		if service.Contact.Name != "" || service.Contact.URL != "" || service.Contact.Email != "" {
			info.Contact = contact
		}
	}

	// Create Document
	document := &v3.Document{
		Version: g.Version,
		Info:    info,
	}

	// Add Speakeasy retry configuration as extension
	g.addSpeakeasyRetryExtension(document, service)

	// Add Speakeasy timeout configuration as extension
	g.addSpeakeasyTimeoutExtension(document, service)

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
				Description: fmt.Sprintf(serverDescriptionTemplate, title),
			},
		}
		document.Servers = servers
	}

	// Create Components
	components := &v3.Components{
		Schemas:       orderedmap.New[string, *base.SchemaProxy](),
		RequestBodies: orderedmap.New[string, *v3.RequestBody](),
		Responses:     orderedmap.New[string, *v3.Response](),
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

	// Add request bodies to components
	g.addRequestBodiesToComponents(components, service)

	// Add response bodies to components
	g.addResponseBodiesToComponents(components, service)

	document.Components = components

	// Create Paths
	paths := orderedmap.New[string, *v3.PathItem]()
	for _, resource := range service.Resources {
		g.addResourceToPaths(resource, paths, service)
	}
	document.Paths = &v3.Paths{
		PathItems: paths,
	}

	// Create tags from resources
	document.Tags = g.createTagsFromResources(service)

	return document
}

// createTagsFromResources creates a tags array from service resources for top-level document organization.
func (g *Generator) createTagsFromResources(service *specification.Service) []*base.Tag {
	if len(service.Resources) == 0 {
		return nil
	}

	tags := make([]*base.Tag, len(service.Resources))
	for i, resource := range service.Resources {
		tags[i] = &base.Tag{
			Name:        resource.Name,
			Description: resource.Description,
		}
	}

	return tags
}

// createEnumSchema creates a base.Schema for an enum using native types.
func (g *Generator) createEnumSchema(enum specification.Enum) *base.Schema {
	schema := &base.Schema{
		Type:        []string{schemaTypeString},
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
		Type:        []string{schemaTypeObject},
		Description: obj.Description,
		Properties:  orderedmap.New[string, *base.SchemaProxy](),
	}

	requiredFields := []string{}
	for _, field := range obj.Fields {
		fieldSchema := g.createFieldSchema(field, service)
		proxy := base.CreateSchemaProxy(fieldSchema)
		schema.Properties.Set(field.TagJSON(), proxy)

		if field.IsRequired(service) {
			requiredFields = append(requiredFields, field.TagJSON())
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
			Type:        []string{schemaTypeArray},
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
		return &base.Schema{Type: []string{schemaTypeString}}
	case specification.FieldTypeInt:
		return &base.Schema{Type: []string{schemaTypeInteger}}
	case specification.FieldTypeBool:
		return &base.Schema{Type: []string{schemaTypeBoolean}}
	case specification.FieldTypeUUID:
		return &base.Schema{
			Type:   []string{schemaTypeString},
			Format: schemaFormatUUID,
		}
	case specification.FieldTypeDate:
		return &base.Schema{
			Type:   []string{schemaTypeString},
			Format: schemaFormatDate,
		}
	case specification.FieldTypeTimestamp:
		return &base.Schema{
			Type:   []string{schemaTypeString},
			Format: schemaFormatDateTime,
		}
	default:
		// Check if it's a custom object or enum
		if service.HasObject(fieldType) || service.HasEnum(fieldType) {
			// Create a proper $ref schema reference using allOf
			refString := schemaReferencePrefix + fieldType
			refProxy := base.CreateSchemaProxyRef(refString)

			// Return a schema with AllOf that contains the reference
			return &base.Schema{
				AllOf: []*base.SchemaProxy{refProxy},
			}
		}
		// Default to string if unknown type
		return &base.Schema{Type: []string{schemaTypeString}}
	}
}

// addResourceToPaths adds resource endpoints to paths using native v3 types.
func (g *Generator) addResourceToPaths(resource specification.Resource, paths *orderedmap.Map[string, *v3.PathItem], service *specification.Service) {
	// Group endpoints by path
	pathGroups := make(map[string][]*specification.Endpoint)
	for _, endpoint := range resource.Endpoints {
		fullPath := endpoint.GetFullPath(resource.Name)
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
		OperationId: resource.Name + endpoint.Name,
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

	// Request body - for now, use inline until we implement proper references
	if len(endpoint.Request.BodyParams) > 0 {
		operation.RequestBody = g.createRequestBody(endpoint.Request.BodyParams, service)
	}

	// Responses
	responses := orderedmap.New[string, *v3.Response]()

	// Success response
	successResponse := g.createResponseReference(endpoint.Response, resource.Name, endpoint.Name, service)
	responses.Set(strconv.Itoa(endpoint.Response.StatusCode), successResponse)

	// Add error responses
	g.addErrorResponses(responses, endpoint, service)

	operation.Responses = &v3.Responses{
		Codes: responses,
	}

	// Add Speakeasy pagination extension if this is a paginated operation
	if g.isPaginatedOperation(endpoint) {
		g.addSpeakeasyPaginationExtension(operation)
	}

	return operation
}

// createParameter creates a v3.Parameter from a field using native types.
func (g *Generator) createParameter(field specification.Field, location string, service *specification.Service) *v3.Parameter {
	isRequired := field.IsRequired(service)
	param := &v3.Parameter{
		Name:        field.TagJSON(),
		In:          location,
		Description: field.Description,
		Required:    &isRequired,
		Schema:      base.CreateSchemaProxy(g.createFieldSchema(field, service)),
	}

	return param
}

// addRequestBodiesToComponents extracts request bodies from all endpoints and adds them to the components section.
func (g *Generator) addRequestBodiesToComponents(components *v3.Components, service *specification.Service) {
	// Track unique request bodies to avoid duplicates
	requestBodyMap := make(map[string]*v3.RequestBody)

	// Iterate through all resources and endpoints to collect request bodies
	for _, resource := range service.Resources {
		for _, endpoint := range resource.Endpoints {
			if len(endpoint.Request.BodyParams) > 0 {
				requestBodyName := g.createRequestBodyName(resource.Name, endpoint.Name)

				// Only add if we haven't seen this request body before
				if _, exists := requestBodyMap[requestBodyName]; !exists {
					requestBody := g.createComponentRequestBody(endpoint.Request.BodyParams, service)
					requestBodyMap[requestBodyName] = requestBody
					components.RequestBodies.Set(requestBodyName, requestBody)
				}
			}
		}
	}
}

// createRequestBodyName creates a systematic name for request bodies.
func (g *Generator) createRequestBodyName(resourceName, endpointName string) string {
	return resourceName + endpointName + requestBodySuffix
}

// createComponentRequestBody creates a v3.RequestBody for the components section.
func (g *Generator) createComponentRequestBody(bodyParams []specification.Field, service *specification.Service) *v3.RequestBody {
	var schema *base.Schema
	var isRequired bool

	// If there's only one body parameter and it references a component schema,
	// use the component schema directly instead of wrapping it in an object
	if len(bodyParams) == 1 {
		field := bodyParams[0]
		if service.HasObject(field.Type) || service.HasEnum(field.Type) {
			// Create a direct reference to the component schema using allOf
			refString := schemaReferencePrefix + field.Type
			refProxy := base.CreateSchemaProxyRef(refString)
			schema = &base.Schema{
				AllOf:       []*base.SchemaProxy{refProxy},
				Description: field.Description,
			}
			isRequired = field.IsRequired(service)
		}
	}

	// If we didn't create a direct reference schema, fall back to the object wrapper approach
	if schema == nil {
		schema = &base.Schema{
			Type:       []string{schemaTypeObject},
			Properties: orderedmap.New[string, *base.SchemaProxy](),
		}

		requiredFields := []string{}
		for _, field := range bodyParams {
			fieldSchema := g.createFieldSchema(field, service)
			proxy := base.CreateSchemaProxy(fieldSchema)
			schema.Properties.Set(field.TagJSON(), proxy)

			if field.IsRequired(service) {
				requiredFields = append(requiredFields, field.TagJSON())
			}
		}

		if len(requiredFields) > 0 {
			schema.Required = requiredFields
		}

		isRequired = len(requiredFields) > 0
	}

	// Create media type
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(schema),
	}

	content := orderedmap.New[string, *v3.MediaType]()
	content.Set(contentTypeJSON, mediaType)

	return &v3.RequestBody{
		Description: requestBodyDescription,
		Content:     content,
		Required:    &isRequired,
	}
}

// createRequestBody creates a v3.RequestBody from body parameters using native types.
func (g *Generator) createRequestBody(bodyParams []specification.Field, service *specification.Service) *v3.RequestBody {
	var schema *base.Schema
	var isRequired bool

	// If there's only one body parameter and it references a component schema,
	// use the component schema directly instead of wrapping it in an object
	if len(bodyParams) == 1 {
		field := bodyParams[0]
		if service.HasObject(field.Type) || service.HasEnum(field.Type) {
			// Create a direct reference to the component schema using allOf
			refString := schemaReferencePrefix + field.Type
			refProxy := base.CreateSchemaProxyRef(refString)
			schema = &base.Schema{
				AllOf:       []*base.SchemaProxy{refProxy},
				Description: field.Description,
			}
			isRequired = field.IsRequired(service)
		}
	}

	// If we didn't create a direct reference schema, fall back to the object wrapper approach
	if schema == nil {
		schema = &base.Schema{
			Type:       []string{schemaTypeObject},
			Properties: orderedmap.New[string, *base.SchemaProxy](),
		}

		requiredFields := []string{}
		for _, field := range bodyParams {
			fieldSchema := g.createFieldSchema(field, service)
			proxy := base.CreateSchemaProxy(fieldSchema)
			schema.Properties.Set(field.TagJSON(), proxy)

			if field.IsRequired(service) {
				requiredFields = append(requiredFields, field.TagJSON())
			}
		}

		if len(requiredFields) > 0 {
			schema.Required = requiredFields
		}

		isRequired = len(requiredFields) > 0
	}

	// Create media type
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(schema),
	}

	content := orderedmap.New[string, *v3.MediaType]()
	content.Set(contentTypeJSON, mediaType)

	return &v3.RequestBody{
		Description: requestBodyDescription,
		Content:     content,
		Required:    &isRequired,
	}
}

// addResponseBodiesToComponents extracts response bodies from all endpoints and adds them to the components section.
func (g *Generator) addResponseBodiesToComponents(components *v3.Components, service *specification.Service) {
	// Track unique response bodies to avoid duplicates
	responseBodyMap := make(map[string]*v3.Response)

	// Iterate through all resources and endpoints to collect response bodies
	for _, resource := range service.Resources {
		for _, endpoint := range resource.Endpoints {
			// Add success response body if it has content
			if endpoint.Response.BodyObject != nil || len(endpoint.Response.BodyFields) > 0 {
				responseBodyName := g.createResponseBodyName(resource.Name, endpoint.Name, endpoint.Response.StatusCode)

				// Only add if we haven't seen this response body before
				if _, exists := responseBodyMap[responseBodyName]; !exists {
					responseBody := g.createComponentResponse(endpoint.Response, service)
					responseBodyMap[responseBodyName] = responseBody
					components.Responses.Set(responseBodyName, responseBody)
				}
			}
		}
	}

	// Add common error response bodies
	g.addErrorResponseBodiesToComponents(components, service)
}

// createResponseBodyName creates a systematic name for response bodies.
func (g *Generator) createResponseBodyName(resourceName, endpointName string, statusCode int) string {
	return resourceName + endpointName + strconv.Itoa(statusCode) + responseBodySuffix
}

// createComponentResponse creates a v3.Response for the components section.
func (g *Generator) createComponentResponse(response specification.EndpointResponse, service *specification.Service) *v3.Response {
	componentResponse := &v3.Response{
		Description: response.Description,
	}

	// Add response content if present
	if response.BodyObject != nil || len(response.BodyFields) > 0 {
		content := orderedmap.New[string, *v3.MediaType]()

		var schema *base.Schema
		if response.BodyObject != nil {
			// Create a proper $ref schema reference using allOf
			refString := schemaReferencePrefix + *response.BodyObject
			refProxy := base.CreateSchemaProxyRef(refString)
			schema = &base.Schema{
				AllOf: []*base.SchemaProxy{refProxy},
			}
		} else if len(response.BodyFields) > 0 {
			// Inline schema from body fields
			schema = &base.Schema{
				Type:       []string{schemaTypeObject},
				Properties: orderedmap.New[string, *base.SchemaProxy](),
			}

			for _, field := range response.BodyFields {
				fieldSchema := g.createFieldSchema(field, service)
				proxy := base.CreateSchemaProxy(fieldSchema)
				schema.Properties.Set(field.TagJSON(), proxy)
			}
		}

		if schema != nil {
			mediaType := &v3.MediaType{
				Schema: base.CreateSchemaProxy(schema),
			}
			content.Set(contentTypeJSON, mediaType)
			componentResponse.Content = content
		}
	}

	return componentResponse
}

// createResponseReference creates a v3.Response that references a component response body.
// TODO: Implement proper response references once we determine the correct libopenapi approach
func (g *Generator) createResponseReference(response specification.EndpointResponse, resourceName, endpointName string, service *specification.Service) *v3.Response {
	// For now, create inline responses while we populate the components section
	// This ensures the components/responses section is populated for future reference implementation
	return g.createResponse(response, service)
}

// addErrorResponseBodiesToComponents adds common error response bodies to the components section.
func (g *Generator) addErrorResponseBodiesToComponents(components *v3.Components, service *specification.Service) {
	// Create error schema
	var errorSchema *base.Schema
	if service.HasObject(errorObjectName) {
		// Create a proper $ref schema reference using allOf
		refString := schemaReferencePrefix + errorObjectName
		refProxy := base.CreateSchemaProxyRef(refString)
		errorSchema = &base.Schema{
			AllOf: []*base.SchemaProxy{refProxy},
		}
	} else {
		// Fallback generic error schema
		errorSchema = &base.Schema{
			Type:       []string{schemaTypeObject},
			Properties: orderedmap.New[string, *base.SchemaProxy](),
		}
		messageSchema := &base.Schema{Type: []string{schemaTypeString}}
		codeSchema := &base.Schema{Type: []string{schemaTypeString}}
		errorSchema.Properties.Set(messageFieldName, base.CreateSchemaProxy(messageSchema))
		errorSchema.Properties.Set(codeFieldName, base.CreateSchemaProxy(codeSchema))
		errorSchema.Required = []string{messageFieldName, codeFieldName}
	}

	errorContent := orderedmap.New[string, *v3.MediaType]()
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(errorSchema),
	}
	errorContent.Set(contentTypeJSON, mediaType)

	// Define common error responses
	errorResponses := map[string]string{
		badRequestDescription:    httpStatus400,
		unauthorizedDescription:  httpStatus401,
		notFoundDescription:      httpStatus404,
		internalErrorDescription: httpStatus500,
	}

	for description, statusCode := range errorResponses {
		responseBodyName := errorResponseBodyPrefix + statusCode + responseBodySuffix
		errorResponse := &v3.Response{
			Description: description,
			Content:     errorContent,
		}
		components.Responses.Set(responseBodyName, errorResponse)
	}
}

// createResponse creates a v3.Response from an endpoint response using native types.
func (g *Generator) createResponse(response specification.EndpointResponse, service *specification.Service) *v3.Response {
	openAPIResponse := &v3.Response{
		Description: response.Description,
	}

	// Add response content if present
	if response.BodyObject != nil || len(response.BodyFields) > 0 {
		content := orderedmap.New[string, *v3.MediaType]()

		var schema *base.Schema
		if response.BodyObject != nil {
			// Create a proper $ref schema reference using allOf
			refString := schemaReferencePrefix + *response.BodyObject
			refProxy := base.CreateSchemaProxyRef(refString)
			schema = &base.Schema{
				AllOf: []*base.SchemaProxy{refProxy},
			}
		} else if len(response.BodyFields) > 0 {
			// Inline schema from body fields
			schema = &base.Schema{
				Type:       []string{schemaTypeObject},
				Properties: orderedmap.New[string, *base.SchemaProxy](),
			}

			for _, field := range response.BodyFields {
				fieldSchema := g.createFieldSchema(field, service)
				proxy := base.CreateSchemaProxy(fieldSchema)
				schema.Properties.Set(field.TagJSON(), proxy)
			}
		}

		if schema != nil {
			mediaType := &v3.MediaType{
				Schema: base.CreateSchemaProxy(schema),
			}
			content.Set(contentTypeJSON, mediaType)
			openAPIResponse.Content = content
		}
	}

	return openAPIResponse
}

// addErrorResponses adds error responses based on errorCodes from the specification.
func (g *Generator) addErrorResponses(responses *orderedmap.Map[string, *v3.Response], endpoint specification.Endpoint, service *specification.Service) {
	// Check if endpoint has body parameters
	hasBodyParams := len(endpoint.Request.BodyParams) > 0

	// Find ErrorCode enum in the service
	var errorCodeEnum *specification.Enum
	for i, enum := range service.Enums {
		if enum.Name == errorCodeEnumName {
			errorCodeEnum = &service.Enums[i]
			break
		}
	}

	if errorCodeEnum == nil {
		// Fallback to default error responses if ErrorCode enum not found
		g.addDefaultErrorResponseReferences(responses)
		return
	}

	// Create error schema for inline responses (components are still populated for future use)
	var errorSchema *base.Schema
	if service.HasObject(errorObjectName) {
		// Create a proper $ref schema reference using allOf
		refString := schemaReferencePrefix + errorObjectName
		refProxy := base.CreateSchemaProxyRef(refString)
		errorSchema = &base.Schema{
			AllOf: []*base.SchemaProxy{refProxy},
		}
	} else {
		// Fallback generic error schema
		errorSchema = &base.Schema{
			Type:       []string{schemaTypeObject},
			Properties: orderedmap.New[string, *base.SchemaProxy](),
		}
		messageSchema := &base.Schema{Type: []string{schemaTypeString}}
		codeSchema := &base.Schema{Type: []string{schemaTypeString}}
		errorSchema.Properties.Set(messageFieldName, base.CreateSchemaProxy(messageSchema))
		errorSchema.Properties.Set(codeFieldName, base.CreateSchemaProxy(codeSchema))
		errorSchema.Required = []string{messageFieldName, codeFieldName}
	}

	errorContent := orderedmap.New[string, *v3.MediaType]()
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(errorSchema),
	}
	errorContent.Set(contentTypeJSON, mediaType)

	// Generate responses for each error code
	for _, enumValue := range errorCodeEnum.Values {
		statusCode, description := g.mapErrorCodeToStatusAndDescription(enumValue.Name, enumValue.Description)

		// Skip 422 UnprocessableEntity if endpoint has no body parameters
		if statusCode == httpStatus422 && !hasBodyParams {
			continue
		}

		// Create inline error response (components are still populated for future reference use)
		errorResponse := &v3.Response{
			Description: description,
			Content:     errorContent,
		}
		responses.Set(statusCode, errorResponse)
	}
}

// addDefaultErrorResponseReferences adds fallback error response references when ErrorCode enum is not found.
// TODO: Update to use actual references once we determine the correct libopenapi approach
func (g *Generator) addDefaultErrorResponseReferences(responses *orderedmap.Map[string, *v3.Response]) {
	// Create fallback generic error schema for inline responses
	errorSchema := &base.Schema{
		Type:       []string{schemaTypeObject},
		Properties: orderedmap.New[string, *base.SchemaProxy](),
	}
	messageSchema := &base.Schema{Type: []string{schemaTypeString}}
	codeSchema := &base.Schema{Type: []string{schemaTypeString}}
	errorSchema.Properties.Set(messageFieldName, base.CreateSchemaProxy(messageSchema))
	errorSchema.Properties.Set(codeFieldName, base.CreateSchemaProxy(codeSchema))
	errorSchema.Required = []string{messageFieldName, codeFieldName}

	errorContent := orderedmap.New[string, *v3.MediaType]()
	mediaType := &v3.MediaType{
		Schema: base.CreateSchemaProxy(errorSchema),
	}
	errorContent.Set(contentTypeJSON, mediaType)

	// Define default error status codes and descriptions
	defaultErrors := map[string]string{
		httpStatus400: badRequestDescription,
		httpStatus401: unauthorizedDescription,
		httpStatus404: notFoundDescription,
		httpStatus500: internalErrorDescription,
	}

	for statusCode, description := range defaultErrors {
		// Create inline error response (components are still populated for future reference use)
		errorResponse := &v3.Response{
			Description: description,
			Content:     errorContent,
		}
		responses.Set(statusCode, errorResponse)
	}
}

// mapErrorCodeToStatusAndDescription maps error code names to HTTP status codes and descriptions.
func (g *Generator) mapErrorCodeToStatusAndDescription(errorCodeName, errorCodeDescription string) (string, string) {
	switch errorCodeName {
	case errorCodeBadRequest:
		return httpStatus400, errorCodeDescription
	case errorCodeUnauthorized:
		return httpStatus401, errorCodeDescription
	case errorCodeForbidden:
		return httpStatus403, errorCodeDescription
	case errorCodeNotFound:
		return httpStatus404, errorCodeDescription
	case errorCodeConflict:
		return httpStatus409, errorCodeDescription
	case errorCodeUnprocessableEntity:
		return httpStatus422, errorCodeDescription
	case errorCodeRateLimited:
		return httpStatus429, errorCodeDescription
	case errorCodeInternal:
		return httpStatus500, errorCodeDescription
	default:
		// Default to 500 for unknown error codes
		return httpStatus500, errorCodeDescription
	}
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

// GenerateFromSpecificationToJSON is a convenience method that generates an OpenAPI document
// from a specification.Service and returns it as JSON in a single call.
// This method creates a generator with default settings, sets a standard title and description,
// generates the OpenAPI document, and converts it to JSON format.
func GenerateFromSpecificationToJSON(service *specification.Service) ([]byte, error) {
	if service == nil {
		return nil, errors.New(errorInvalidService)
	}

	// Create generator with default configuration
	generator := newGenerator()

	// Set basic configuration based on service
	generator.Title = service.Name + apiTitleSuffix
	generator.Description = defaultAPIDescription

	// Generate OpenAPI document
	document, err := generator.GenerateFromService(service)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	// Convert to JSON
	jsonBytes, err := generator.ToJSON(document)
	if err != nil {
		return nil, fmt.Errorf("failed to convert OpenAPI document to JSON: %w", err)
	}

	return jsonBytes, nil
}
