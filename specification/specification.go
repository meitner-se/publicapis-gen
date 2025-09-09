package specification

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aarondl/strmangle"
)

// CRUD Operations
const (
	OperationCreate = "Create"
	OperationRead   = "Read"
	OperationUpdate = "Update"
	OperationDelete = "Delete"
)

// Field Types
const (
	FieldTypeUUID      = "UUID"
	FieldTypeDate      = "Date"
	FieldTypeTimestamp = "Timestamp"
	FieldTypeString    = "String"
	FieldTypeInt       = "Int"
	FieldTypeBool      = "Bool"
)

// Field Modifiers
const (
	ModifierNullable = "nullable"
	ModifierArray    = "array"
)

// Filter suffixes
const (
	filterSuffix         = "Filter"
	filterEqualsSuffix   = "FilterEquals"
	filterRangeSuffix    = "FilterRange"
	filterContainsSuffix = "FilterContains"
	filterLikeSuffix     = "FilterLike"
	filterNullSuffix     = "FilterNull"
)

// Filter field names
const (
	filterFieldEquals         = "Equals"
	filterFieldNotEquals      = "NotEquals"
	filterFieldGreaterThan    = "GreaterThan"
	filterFieldSmallerThan    = "SmallerThan"
	filterFieldGreaterOrEqual = "GreaterOrEqual"
	filterFieldSmallerOrEqual = "SmallerOrEqual"
	filterFieldContains       = "Contains"
	filterFieldNotContains    = "NotContains"
	filterFieldLike           = "Like"
	filterFieldNotLike        = "NotLike"
	filterFieldNull           = "Null"
	filterFieldNotNull        = "NotNull"
	filterFieldOrCondition    = "OrCondition"
	filterFieldNestedFilters  = "NestedFilters"
)

// Filter description templates
const (
	descriptionFilterObject                   = "Filter object for "
	descriptionEqualityFilters                = "Equality filters for "
	descriptionInequalityFilters              = "Inequality filters for "
	descriptionGreaterThanFilters             = "Greater than filters for "
	descriptionSmallerThanFilters             = "Smaller than filters for "
	descriptionGreaterOrEqualFilters          = "Greater than or equal filters for "
	descriptionSmallerOrEqualFilters          = "Smaller than or equal filters for "
	descriptionContainsFilters                = "Contains filters for "
	descriptionNotContainsFilters             = "Not contains filters for "
	descriptionLikeFilters                    = "LIKE filters for "
	descriptionNotLikeFilters                 = "NOT LIKE filters for "
	descriptionNullFilters                    = "Null filters for "
	descriptionNotNullFilters                 = "Not null filters for "
	descriptionOrCondition                    = "OrCondition decides if this filter is within an OR-condition or AND-condition"
	descriptionNestedFiltersTemplate          = "NestedFilters of the "
	descriptionNestedFiltersSuffix            = ", useful for more complex filters"
	descriptionEqualityInequalityFilterFields = "Equality/Inequality filter fields for "
	descriptionRangeFilterFields              = "Range filter fields for "
	descriptionContainsFilterFields           = "Contains filter fields for "
	descriptionLikeFilterFields               = "LIKE filter fields for "
	descriptionNullFilterFields               = "Null filter fields for "
)

// Error Code Values
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

// Error Code Descriptions
const (
	descriptionErrorCodeEnum                = "Standard error codes used in API responses"
	descriptionErrorCodeBadRequest          = "The request was malformed or contained invalid parameters. 400 status code"
	descriptionErrorCodeUnauthorized        = "The request is missing valid authentication credentials. 401 status code"
	descriptionErrorCodeForbidden           = "Request is authenticated, but the user is not allowed to perform the operation. 403 status code"
	descriptionErrorCodeNotFound            = "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code"
	descriptionErrorCodeConflict            = "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code"
	descriptionErrorCodeUnprocessableEntity = "The request was well-formed but failed validation (e.g. invalid field format or constraints), 422 status code"
	descriptionErrorCodeRateLimited         = "When the rate limit has been exceeded, 429 status code"
	descriptionErrorCodeInternal            = "Some serverside issue, 5xx status code"
)

// Error object constants
const (
	errorObjectName              = "Error"
	errorObjectDescription       = "Standard error response object containing error code and message"
	errorCodeFieldName           = "Code"
	errorCodeFieldDescription    = "The specific error code indicating the type of error"
	errorMessageFieldName        = "Message"
	errorMessageFieldDescription = "Human-readable error message providing additional details"
	errorCodeEnumName            = "ErrorCode"
)

// ErrorFieldCode Values
const (
	errorFieldCodeAlreadyExists = "AlreadyExists"
	errorFieldCodeRequired      = "Required"
	errorFieldCodeNotFound      = "NotFound"
	errorFieldCodeInvalidValue  = "InvalidValue"
)

// ErrorFieldCode Descriptions
const (
	descriptionErrorFieldCodeEnum          = "Error codes for field-level validation errors"
	descriptionErrorFieldCodeAlreadyExists = "The field value already exists and violates a unique constraint (e.g., duplicate email address or username)"
	descriptionErrorFieldCodeRequired      = "The field is required but is missing or empty in the request"
	descriptionErrorFieldCodeNotFound      = "A referenced resource or relation does not exist (e.g., foreign key constraint violation)"
	descriptionErrorFieldCodeInvalidValue  = "The field contains an invalid value (e.g., invalid enum value, malformed data, or value out of allowed range)"
)

// ErrorField object constants
const (
	errorFieldObjectName              = "ErrorField"
	errorFieldObjectDescription       = "Field-specific error information containing error code and message for validation errors"
	errorFieldCodeFieldName           = "Code"
	errorFieldCodeFieldDescription    = "The specific error code indicating the type of field validation error"
	errorFieldMessageFieldName        = "Message"
	errorFieldMessageFieldDescription = "Human-readable error message providing details about the field validation error"
	errorFieldCodeEnumName            = "ErrorFieldCode"
)

// Pagination object constants
const (
	paginationObjectName        = "Pagination"
	paginationObjectDescription = "Pagination parameters for controlling result sets in list operations"
	offsetFieldName             = "Offset"
	offsetFieldDescription      = "Number of items to skip from the beginning of the result set"
	limitFieldName              = "Limit"
	limitFieldDescription       = "Maximum number of items to return in the result set"
)

// Auto-column constants
const (
	autoColumnIDName            = "ID"
	autoColumnIDDescTemplate    = "Unique identifier for the %s"
	autoColumnCreatedAtName     = "CreatedAt"
	autoColumnCreatedAtTemplate = "Timestamp when the %s was created"
	autoColumnCreatedByName     = "CreatedBy"
	autoColumnCreatedByTemplate = "User who created the %s"
	autoColumnUpdatedAtName     = "UpdatedAt"
	autoColumnUpdatedAtTemplate = "Timestamp when the %s was last updated"
	autoColumnUpdatedByName     = "UpdatedBy"
	autoColumnUpdatedByTemplate = "User who last updated the %s"
)

// HTTP Methods
const (
	httpMethodGet    = "GET"
	httpMethodPost   = "POST"
	httpMethodPatch  = "PATCH"
	httpMethodPut    = "PUT"
	httpMethodDelete = "DELETE"
)

// Content Types
const (
	contentTypeJSON = "application/json"
)

// Create Endpoint Constants
const (
	createEndpointName        = "Create"
	createEndpointPath        = ""
	createEndpointTitlePrefix = "Create "
	createEndpointDescPrefix  = "Create a new "
	createResponseStatusCode  = 201
)

// Update Endpoint Constants
const (
	updateEndpointName        = "Update"
	updateEndpointPath        = "/{id}"
	updateEndpointTitlePrefix = "Update "
	updateEndpointDescPrefix  = "Update a "
	updateResponseStatusCode  = 200
	updateIDParamName         = "id"
	updateIDParamDescTemplate = "The unique identifier of the %s to update"
)

// Delete Endpoint Constants
const (
	deleteEndpointName        = "Delete"
	deleteEndpointPath        = "/{id}"
	deleteEndpointTitlePrefix = "Delete "
	deleteEndpointDescPrefix  = "Delete a "
	deleteResponseStatusCode  = 204
	deleteIDParamName         = "id"
	deleteIDParamDescTemplate = "The unique identifier of the %s to delete"
)

// Get Endpoint Constants
const (
	getEndpointName        = "Get"
	getEndpointPath        = "/{id}"
	getEndpointTitlePrefix = "Retrieve an existing "
	getResponseStatusCode  = 200
	getIDParamName         = "id"
	getIDParamDescTemplate = "The unique identifier of the %s to retrieve"
)

// List Endpoint Constants
const (
	listEndpointName         = "List"
	listEndpointPath         = ""
	listEndpointTitlePrefix  = "List all "
	listEndpointDescTemplate = "Returns a paginated list of all `%s` in your organization."
	listResponseStatusCode   = 200
	listLimitParamName       = "limit"
	listLimitParamDesc       = "The maximum number of items to return (default: 50)"
	listLimitDefaultValue    = "50"
	listOffsetParamName      = "offset"
	listOffsetParamDesc      = "The number of items to skip before starting to return results (default: 0)"
	listOffsetDefaultValue   = "0"
)

// Search Endpoint Constants
const (
	searchEndpointName         = "Search"
	searchEndpointPath         = "/_search"
	searchEndpointTitlePrefix  = "Search "
	searchEndpointDescTemplate = "Search for `%s` with filtering capabilities."
	searchResponseStatusCode   = 200
	searchFilterParamName      = "Filter"
	searchFilterParamDesc      = "Filter criteria to search for specific records"
)

// Request Error Constants
const (
	requestErrorSuffix            = "RequestError"
	requestErrorDescriptionPrefix = "Request error object for "
)

// Comment formatting constants
const (
	commentPrefix     = "// "
	nameDescSeparator = ": "
	newlineChar       = "\n"
	pathSeparator     = "/"
)

// ServiceServer represents a server in the API service.
type ServiceServer struct {
	// URL of the server
	URL string `json:"url"`

	// Description of the server
	Description string `json:"description,omitempty"`
}

// Service is the definition of an API service.
type Service struct {
	// Name of the service
	Name string `json:"name"`

	// Version of the service
	Version string `json:"version,omitempty"`

	// Servers that are part of the service
	Servers []ServiceServer `json:"servers,omitempty"`

	// Enums that are used in the service
	Enums []Enum `json:"enums"`

	// Objects that are used in the service
	Objects []Object `json:"objects"`

	// Resources that are part of the service
	Resources []Resource `json:"resources"`
}

// Enum represents an enumeration with possible values.
type Enum struct {
	// Name of the enum
	Name string `json:"name"`

	// Description of the enum
	Description string `json:"description"`

	// Values that are possible for the enum
	Values []EnumValue `json:"values"`
}

// EnumValue represents a single value in an enumeration.
type EnumValue struct {
	// Name of the enum value, for example Male for the Enum Gender - should be unique in the Enum
	Name string `json:"name"`

	// Description for the enum value
	Description string `json:"description"`
}

// Object is a shared object within the service,
// can be used by multiple resources.
type Object struct {
	// Name of the object, should be unique in the service
	Name string `json:"name"`

	// Description about the object
	Description string `json:"description"`

	// Fields in the object
	Fields []Field `json:"fields"`
}

// Resource represents a resource in the API with its operations and fields.
type Resource struct {
	// Name of the resource, should be unique within the service
	Name string `json:"name"`

	// Description about the resource
	Description string `json:"description"`

	// Operations that are allowed for the resource can be all of Create, Update, Read, Delete
	Operations []string `json:"operations"`

	// Fields of the resource
	Fields []ResourceField `json:"fields"`

	// Endpoints of the resource
	Endpoints []Endpoint `json:"endpoints"`

	// SkipAutoColumns indicates whether to skip generating auto columns (ID, CreatedAt, etc.) for this resource
	SkipAutoColumns bool `json:"skip_auto_columns,omitempty"`
}

// Field contains information about a field within an endpoint or resource or Object.
type Field struct {
	// Name of the field, should be unique in the Resource or Object or Endpoint
	Name string `json:"name"`

	// Description of the field, explain the reason what it is used for and why it's needed
	Description string `json:"description"`

	// Type of the field, can be one of the types (UUID, Date, Timestamp, String, Int, Bool) or one of the custom Objects
	Type string `json:"type"`

	// Default value of the field
	Default string `json:"default,omitempty"`

	// Example value of the field
	Example string `json:"example,omitempty"`

	// Modifiers of the field, can be nullable or array
	Modifiers []string `json:"modifiers,omitempty"`
}

// ResourceField is used within a resource it extends the field with an operations configuration.
type ResourceField struct {
	Field `json:",inline"`

	// Operations that the field is allowed in (Create,Update,Delete,Read)
	Operations []string `json:"operations"`
}

// Endpoint represents an API endpoint within a resource.
type Endpoint struct {
	// Name of the endpoint, Should be unique within the resource.
	// For example: "Get", "Create", "Update", "Delete", "Search"...
	Name string `json:"name"`

	// Title for the endpoint, should be unique within the resource.
	// For example: "Get School", "Create School", "Update School", "Delete School", "Search School"...
	Title string `json:"title"`

	// Description of the endpoint
	Description string `json:"description"`

	// HTTP method of the endpoint
	Method string `json:"method"`

	// Path of the endpoint, "/:id". No need to include the resource name, it will be added automatically.
	Path string `json:"path"`

	// Request that is used in the endpoint
	Request EndpointRequest `json:"request"`

	// Response that is used in the endpoint on success
	Response EndpointResponse `json:"response"`
}

// EndpointRequest represents the request structure for an API endpoint.
type EndpointRequest struct {
	// Content-Type of the request
	ContentType string `json:"content_type"`

	// Headers that are used in the request
	Headers []Field `json:"headers"`

	// Path parameters that are used in the endpoint
	PathParams []Field `json:"path_params"`

	// Query parameters that are used in the endpoint
	QueryParams []Field `json:"query_params"`

	// Body parameters that are used in the endpoint
	BodyParams []Field `json:"body_params"`
}

// EndpointResponse represents the response structure for an API endpoint.
type EndpointResponse struct {
	// Content-Type of the response
	ContentType string `json:"content_type"`

	// HTTP status code this response represents (e.g. 200, 201, 400)
	StatusCode int `json:"status_code"`

	// Headers returned in the response
	Headers []Field `json:"headers"`

	// Body fields returned in the response (flat or object)
	BodyFields []Field `json:"body_fields"`

	// If a full object is returned (instead of individual fields) - can be object or Resource
	BodyObject *string `json:"body_object,omitempty"`
}

// ApplyOverlay applies an overlay to a specification, generating Objects and endpoints from Resources.
// It creates Objects for Resources that have the "Read" operation, including all fields
// that support the "Read" operation in the generated Object.
// It generates standard CRUD endpoints (Create, Read, Update, Delete) and additional endpoints (List, Search)
// based on the operations supported by each Resource.
// It also adds default error handling objects and pagination support to every service.
func ApplyOverlay(input *Service) *Service {
	if input == nil {
		return nil
	}

	// Create a deep copy of the input service
	result := &Service{
		Name:      input.Name,
		Enums:     make([]Enum, 0, len(input.Enums)+2),     // +2 for ErrorCode and ErrorFieldCode enums
		Objects:   make([]Object, 0, len(input.Objects)+3), // +3 for Error, ErrorField, and Pagination objects
		Resources: make([]Resource, len(input.Resources)),
	}

	// Add default enums and objects if they don't already exist
	addDefaultEnumsAndObjects(result, input)

	// Copy resources
	copy(result.Resources, input.Resources)

	// Generate Objects and endpoints from Resources
	generateObjectsFromResources(result, input.Resources)
	generateEndpointsFromResources(result, input.Resources)

	// Generate RequestError objects for types used in body parameters
	// This happens at the end to ensure all objects and endpoints are generated first
	generateRequestErrorObjectsForBodyParams(result)

	return result
}

// addDefaultEnumsAndObjects adds the default error and pagination objects to the service if they don't already exist.
func addDefaultEnumsAndObjects(result *Service, input *Service) {
	// Check if ErrorCode enum, Error object, ErrorFieldCode enum, ErrorField object, and Pagination object already exist
	errorCodeEnumExists := false
	errorObjectExists := false
	errorFieldCodeEnumExists := false
	errorFieldObjectExists := false
	paginationObjectExists := false
	for _, enum := range input.Enums {
		if enum.Name == errorCodeEnumName {
			errorCodeEnumExists = true
		}
		if enum.Name == errorFieldCodeEnumName {
			errorFieldCodeEnumExists = true
		}
	}
	for _, object := range input.Objects {
		if object.Name == errorObjectName {
			errorObjectExists = true
		}
		if object.Name == errorFieldObjectName {
			errorFieldObjectExists = true
		}
		if object.Name == paginationObjectName {
			paginationObjectExists = true
		}
	}

	// Copy existing enums first to preserve order
	result.Enums = append(result.Enums, input.Enums...)

	// Add default ErrorCode enum if it doesn't exist
	if !errorCodeEnumExists {
		errorCodeEnum := Enum{
			Name:        errorCodeEnumName,
			Description: descriptionErrorCodeEnum,
			Values: []EnumValue{
				{Name: errorCodeBadRequest, Description: descriptionErrorCodeBadRequest},
				{Name: errorCodeUnauthorized, Description: descriptionErrorCodeUnauthorized},
				{Name: errorCodeForbidden, Description: descriptionErrorCodeForbidden},
				{Name: errorCodeNotFound, Description: descriptionErrorCodeNotFound},
				{Name: errorCodeConflict, Description: descriptionErrorCodeConflict},
				{Name: errorCodeUnprocessableEntity, Description: descriptionErrorCodeUnprocessableEntity},
				{Name: errorCodeRateLimited, Description: descriptionErrorCodeRateLimited},
				{Name: errorCodeInternal, Description: descriptionErrorCodeInternal},
			},
		}
		result.Enums = append(result.Enums, errorCodeEnum)
	}

	// Add default ErrorFieldCode enum if it doesn't exist
	if !errorFieldCodeEnumExists {
		errorFieldCodeEnum := Enum{
			Name:        errorFieldCodeEnumName,
			Description: descriptionErrorFieldCodeEnum,
			Values: []EnumValue{
				{Name: errorFieldCodeAlreadyExists, Description: descriptionErrorFieldCodeAlreadyExists},
				{Name: errorFieldCodeRequired, Description: descriptionErrorFieldCodeRequired},
				{Name: errorFieldCodeNotFound, Description: descriptionErrorFieldCodeNotFound},
				{Name: errorFieldCodeInvalidValue, Description: descriptionErrorFieldCodeInvalidValue},
			},
		}
		result.Enums = append(result.Enums, errorFieldCodeEnum)
	}

	// Copy existing objects first to preserve order
	result.Objects = append(result.Objects, input.Objects...)

	// Add default Error object if it doesn't exist
	if !errorObjectExists {
		errorObject := Object{
			Name:        errorObjectName,
			Description: errorObjectDescription,
			Fields: []Field{
				{
					Name:        errorCodeFieldName,
					Description: errorCodeFieldDescription,
					Type:        errorCodeEnumName,
				},
				{
					Name:        errorMessageFieldName,
					Description: errorMessageFieldDescription,
					Type:        FieldTypeString,
				},
			},
		}
		result.Objects = append(result.Objects, errorObject)
	}

	// Add default ErrorField object if it doesn't exist
	if !errorFieldObjectExists {
		errorFieldObject := Object{
			Name:        errorFieldObjectName,
			Description: errorFieldObjectDescription,
			Fields: []Field{
				{
					Name:        errorFieldCodeFieldName,
					Description: errorFieldCodeFieldDescription,
					Type:        errorFieldCodeEnumName,
				},
				{
					Name:        errorFieldMessageFieldName,
					Description: errorFieldMessageFieldDescription,
					Type:        FieldTypeString,
				},
			},
		}
		result.Objects = append(result.Objects, errorFieldObject)
	}

	// Add default Pagination object if it doesn't exist
	if !paginationObjectExists {
		paginationObject := Object{
			Name:        paginationObjectName,
			Description: paginationObjectDescription,
			Fields: []Field{
				{
					Name:        offsetFieldName,
					Description: offsetFieldDescription,
					Type:        FieldTypeInt,
				},
				{
					Name:        limitFieldName,
					Description: limitFieldDescription,
					Type:        FieldTypeInt,
				},
			},
		}
		result.Objects = append(result.Objects, paginationObject)
	}
}

// generateObjectsFromResources generates Objects from Resources that have Read operations.
func generateObjectsFromResources(result *Service, resources []Resource) {
	for _, resource := range resources {
		// Check if the resource has Read operation
		if resource.HasReadOperation() {
			// Check if an object with this name already exists
			if !result.HasObject(resource.Name) {
				// Get readable fields from the resource
				fields := resource.GetReadableFields()

				// Add auto-columns to the object if not skipped
				if !resource.ShouldSkipAutoColumns() {
					autoColumns := CreateAutoColumns(resource.Name)
					fields = append(autoColumns, fields...)
				}

				// Create a new Object based on the Resource
				newObject := Object{
					Name:        resource.Name,
					Description: resource.Description,
					Fields:      fields,
				}

				// Add the new object to the result
				result.Objects = append(result.Objects, newObject)
			}
		}
	}
}

// generateEndpointsFromResources generates standard CRUD and additional endpoints for all resources.
func generateEndpointsFromResources(result *Service, resources []Resource) {
	for _, resource := range resources {
		generateCreateEndpoint(result, resource)
		generateUpdateEndpoint(result, resource)
		generateDeleteEndpoint(result, resource)
		generateGetEndpoint(result, resource)
		generateListEndpoint(result, resource)
		generateSearchEndpoint(result, resource)
	}
}

// generateCreateEndpoint generates a Create endpoint for resources that have Create operations.
func generateCreateEndpoint(result *Service, resource Resource) {
	if resource.HasCreateOperation() && !resource.HasEndpoint(createEndpointName) {
		bodyParams := resource.GetCreateBodyParams()
		resourceName := resource.Name
		createEndpoint := Endpoint{
			Name:        createEndpointName,
			Title:       createEndpointTitlePrefix + resource.Name,
			Description: createEndpointDescPrefix + resource.Name,
			Method:      httpMethodPost,
			Path:        createEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{}, bodyParams),
			Response:    createStandardResponse(createResponseStatusCode, &resourceName),
		}

		addEndpointToResource(result, resource.Name, createEndpoint)
	}
}

// generateUpdateEndpoint generates an Update endpoint for resources that have Update operations.
func generateUpdateEndpoint(result *Service, resource Resource) {
	if resource.HasUpdateOperation() && !resource.HasEndpoint(updateEndpointName) {
		bodyParams := resource.GetUpdateBodyParams()
		idParam := Field{
			Name:        updateIDParamName,
			Description: fmt.Sprintf(updateIDParamDescTemplate, resource.Name),
			Type:        FieldTypeUUID,
		}
		resourceName := resource.Name
		updateEndpoint := Endpoint{
			Name:        updateEndpointName,
			Title:       updateEndpointTitlePrefix + resource.Name,
			Description: updateEndpointDescPrefix + resource.Name,
			Method:      httpMethodPatch,
			Path:        updateEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, bodyParams),
			Response:    createStandardResponse(updateResponseStatusCode, &resourceName),
		}

		addEndpointToResource(result, resource.Name, updateEndpoint)
	}
}

// generateDeleteEndpoint generates a Delete endpoint for resources that have Delete operations.
func generateDeleteEndpoint(result *Service, resource Resource) {
	if resource.HasDeleteOperation() && !resource.HasEndpoint(deleteEndpointName) {
		idParam := Field{
			Name:        deleteIDParamName,
			Description: fmt.Sprintf(deleteIDParamDescTemplate, resource.Name),
			Type:        FieldTypeUUID,
		}
		deleteEndpoint := Endpoint{
			Name:        deleteEndpointName,
			Title:       deleteEndpointTitlePrefix + resource.Name,
			Description: deleteEndpointDescPrefix + resource.Name,
			Method:      httpMethodDelete,
			Path:        deleteEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, []Field{}),
			Response:    createStandardResponse(deleteResponseStatusCode, nil), // No body object for delete
		}

		addEndpointToResource(result, resource.Name, deleteEndpoint)
	}
}

// generateGetEndpoint generates a Get endpoint for resources that have Read operations.
func generateGetEndpoint(result *Service, resource Resource) {
	if resource.HasReadOperation() && !resource.HasEndpoint(getEndpointName) {
		idParam := Field{
			Name:        getIDParamName,
			Description: fmt.Sprintf(getIDParamDescTemplate, resource.Name),
			Type:        FieldTypeUUID,
		}
		resourceName := resource.Name
		getEndpoint := Endpoint{
			Name:        getEndpointName,
			Title:       getEndpointTitlePrefix + resource.Name,
			Description: fmt.Sprintf("Retrieves the `%s` with the given ID.", resource.Name),
			Method:      httpMethodGet,
			Path:        getEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, []Field{}),
			Response:    createStandardResponse(getResponseStatusCode, &resourceName),
		}

		addEndpointToResource(result, resource.Name, getEndpoint)
	}
}

// generateListEndpoint generates a List endpoint for resources that have Read operations.
func generateListEndpoint(result *Service, resource Resource) {
	if resource.HasReadOperation() && !resource.HasEndpoint(listEndpointName) {
		limitParam := CreateLimitParam()
		offsetParam := CreateOffsetParam()
		paginationField := CreatePaginationField()
		dataField := CreateDataField(resource.Name)
		pluralResourceName := resource.GetPluralName()

		listEndpoint := Endpoint{
			Name:        listEndpointName,
			Title:       listEndpointTitlePrefix + pluralResourceName,
			Description: fmt.Sprintf(listEndpointDescTemplate, pluralResourceName),
			Method:      httpMethodGet,
			Path:        listEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{limitParam, offsetParam}, []Field{}),
			Response:    createListResponse(listResponseStatusCode, dataField, paginationField),
		}

		addEndpointToResource(result, resource.Name, listEndpoint)
	}
}

// generateSearchEndpoint generates a Search endpoint for resources that have Read operations.
func generateSearchEndpoint(result *Service, resource Resource) {
	if resource.HasReadOperation() && !resource.HasEndpoint(searchEndpointName) {
		limitParam := CreateLimitParam()
		offsetParam := CreateOffsetParam()
		filterParam := Field{
			Name:        searchFilterParamName,
			Description: searchFilterParamDesc,
			Type:        resource.Name + filterSuffix,
		}
		paginationField := CreatePaginationField()
		dataField := CreateDataField(resource.Name)
		pluralResourceName := resource.GetPluralName()

		searchEndpoint := Endpoint{
			Name:        searchEndpointName,
			Title:       searchEndpointTitlePrefix + pluralResourceName,
			Description: fmt.Sprintf(searchEndpointDescTemplate, pluralResourceName),
			Method:      httpMethodPost,
			Path:        searchEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{limitParam, offsetParam}, []Field{filterParam}),
			Response:    createListResponse(searchResponseStatusCode, dataField, paginationField),
		}

		addEndpointToResource(result, resource.Name, searchEndpoint)
	}
}

// addEndpointToResource adds an endpoint to a resource by name.
func addEndpointToResource(result *Service, resourceName string, endpoint Endpoint) {
	for i := range result.Resources {
		if result.Resources[i].Name == resourceName {
			result.Resources[i].Endpoints = append(result.Resources[i].Endpoints, endpoint)
			break
		}
	}
}

// createStandardRequest creates a standard endpoint request with the given path and body parameters.
func createStandardRequest(pathParams []Field, queryParams []Field, bodyParams []Field) EndpointRequest {
	return EndpointRequest{
		ContentType: contentTypeJSON,
		Headers:     []Field{},
		PathParams:  pathParams,
		QueryParams: queryParams,
		BodyParams:  bodyParams,
	}
}

// createStandardResponse creates a standard endpoint response with the given status code and optional body object.
func createStandardResponse(statusCode int, bodyObject *string) EndpointResponse {
	return EndpointResponse{
		ContentType: contentTypeJSON,
		StatusCode:  statusCode,
		Headers:     []Field{},
		BodyFields:  []Field{},
		BodyObject:  bodyObject,
	}
}

// createListResponse creates a standard list endpoint response with pagination and data fields.
func createListResponse(statusCode int, dataField Field, paginationField Field) EndpointResponse {
	return EndpointResponse{
		ContentType: contentTypeJSON,
		StatusCode:  statusCode,
		Headers:     []Field{},
		BodyFields:  []Field{dataField, paginationField},
		BodyObject:  nil,
	}
}

// collectTypesUsedInBodyParams collects all types (including nested) used in request body parameters.
func collectTypesUsedInBodyParams(service *Service) map[string]bool {
	usedTypes := make(map[string]bool)

	// Collect types from all endpoint body parameters
	for _, resource := range service.Resources {
		for _, endpoint := range resource.Endpoints {
			for _, bodyParam := range endpoint.Request.BodyParams {
				collectTypeRecursively(bodyParam.Type, usedTypes, service.Objects)
			}
		}
	}

	return usedTypes
}

// collectTypeRecursively collects a type and all its nested object types recursively.
func collectTypeRecursively(fieldType string, usedTypes map[string]bool, objects []Object) {
	// Skip if already processed
	if usedTypes[fieldType] {
		return
	}

	// Mark this type as used
	usedTypes[fieldType] = true

	// If it's an object type, recursively collect its field types
	for _, obj := range objects {
		if obj.Name == fieldType {
			for _, field := range obj.Fields {
				collectTypeRecursively(field.Type, usedTypes, objects)
			}
			break
		}
	}
}

// generateRequestErrorObjectsForBodyParams generates RequestError objects only for types used in body parameters.
func generateRequestErrorObjectsForBodyParams(service *Service) {
	// Collect all types used in body parameters
	usedTypes := collectTypesUsedInBodyParams(service)

	// Generate RequestError objects for each used type
	for typeName := range usedTypes {
		// Skip primitive types - they don't need their own RequestError objects
		if isPrimitiveType(typeName) {
			continue
		}

		// Find the object definition
		for _, obj := range service.Objects {
			if obj.Name == typeName {
				requestErrorName := obj.Name + requestErrorSuffix
				requestErrorDescription := requestErrorDescriptionPrefix + obj.Name
				requestError := generateRequestErrorObject(requestErrorName, requestErrorDescription, obj.Fields, service.Objects)
				service.Objects = append(service.Objects, requestError)
				break
			}
		}
	}

	// Generate RequestError objects for specific endpoints that have body parameters
	for _, resource := range service.Resources {
		for _, endpoint := range resource.Endpoints {
			if len(endpoint.Request.BodyParams) > 0 {
				requestErrorName := resource.Name + endpoint.Name + requestErrorSuffix
				requestErrorDescription := requestErrorDescriptionPrefix + resource.Name + " " + endpoint.Name + " endpoint"
				requestError := generateRequestErrorObject(requestErrorName, requestErrorDescription, endpoint.Request.BodyParams, service.Objects)
				service.Objects = append(service.Objects, requestError)
			}
		}
	}
}

// isComparableType returns true if the field type supports range operations.
func isComparableType(fieldType string) bool {
	switch fieldType {
	case FieldTypeInt, FieldTypeDate, FieldTypeTimestamp:
		return true
	default:
		return false
	}
}

// isStringType returns true if the field type is a string type that supports LIKE operations.
func isStringType(fieldType string) bool {
	return fieldType == FieldTypeString
}

// canBeNull returns true if the field can be null (has nullable modifier or is an array).
func canBeNull(field Field) bool {
	return field.IsNullable() || field.IsArray()
}

// isPrimitiveType returns true if the field type is a primitive type.
func isPrimitiveType(fieldType string) bool {
	switch fieldType {
	case FieldTypeUUID, FieldTypeDate, FieldTypeTimestamp, FieldTypeString, FieldTypeInt, FieldTypeBool:
		return true
	default:
		return false
	}
}

// isObjectType returns true if the field type is a custom object type.
// This assumes all object types exist in the provided objects slice.
func isObjectType(fieldType string, objects []Object) bool {
	if isPrimitiveType(fieldType) {
		return false
	}

	for _, obj := range objects {
		if obj.Name == fieldType {
			return true
		}
	}
	return false
}

// convertFieldToErrorField converts a field to its error counterpart.
// Primitive types become *ErrorField, object types become their RequestError equivalent.
func convertFieldToErrorField(field Field, objects []Object) Field {
	errorField := Field{
		Name:        field.Name,
		Description: field.Description,
		Type:        errorFieldObjectName,       // Default to ErrorField type
		Modifiers:   []string{ModifierNullable}, // All error fields are nullable
	}

	if isObjectType(field.Type, objects) {
		errorField.Type = field.Type + requestErrorSuffix
	} else if strings.HasSuffix(field.Type, filterSuffix) {
		// Handle filter types (e.g., UsersFilter -> UsersFilterRequestError)
		errorField.Type = field.Type + requestErrorSuffix
	}
	// For primitive types and other types (enums, etc.), use the default ErrorField type

	return errorField
}

// generateRequestErrorObject generates a RequestError object from a list of fields.
func generateRequestErrorObject(objectName string, description string, fields []Field, objects []Object) Object {
	errorFields := make([]Field, 0, len(fields))

	for _, field := range fields {
		errorField := convertFieldToErrorField(field, objects)
		errorFields = append(errorFields, errorField)
	}

	return Object{
		Name:        objectName,
		Description: description,
		Fields:      errorFields,
	}
}

// generateFilterField creates a filter field based on the original field and filter type.
func generateFilterField(originalField Field, isNullable bool, isArray bool) Field {
	modifiers := []string{}

	if isNullable {
		modifiers = append(modifiers, ModifierNullable)
	}

	if isArray {
		modifiers = append(modifiers, ModifierArray)
	}

	return Field{
		Name:        originalField.Name,
		Description: originalField.Description,
		Type:        originalField.Type,
		Modifiers:   modifiers,
	}
}

// generateNestedFilterField creates a filter field for nested objects, using the appropriate filter type.
func generateNestedFilterField(originalField Field, filterSuffix string, isNullable bool, isArray bool, objects []Object) Field {
	modifiers := []string{}

	if isNullable {
		modifiers = append(modifiers, ModifierNullable)
	}

	if isArray {
		modifiers = append(modifiers, ModifierArray)
	}

	// For nested object fields, use the filter version of the object type
	filterType := originalField.Type
	if isObjectType(originalField.Type, objects) {
		filterType = originalField.Type + filterSuffix
	}

	return Field{
		Name:        originalField.Name,
		Description: originalField.Description,
		Type:        filterType,
		Modifiers:   modifiers,
	}
}

// ApplyFilterOverlay applies filter overlay to a specification, generating Filter objects
// from existing Objects. This should be called after ApplyOverlay to ensure all Objects
// are available for filter generation.
func ApplyFilterOverlay(input *Service) *Service {
	if input == nil {
		return nil
	}

	// Create a deep copy of the input service
	result := &Service{
		Name:      input.Name,
		Enums:     make([]Enum, len(input.Enums)),
		Objects:   make([]Object, 0, len(input.Objects)*7), // Estimate for filter objects
		Resources: make([]Resource, len(input.Resources)),
	}

	// Copy enums
	copy(result.Enums, input.Enums)

	// Copy existing objects first
	result.Objects = append(result.Objects, input.Objects...)

	// Copy resources
	copy(result.Resources, input.Resources)

	// Generate Filter objects from existing Objects
	for _, obj := range input.Objects {
		// Generate main filter object
		mainFilter := Object{
			Name:        obj.Name + filterSuffix,
			Description: descriptionFilterObject + obj.Name,
			Fields: []Field{
				{
					Name:        filterFieldEquals,
					Description: descriptionEqualityFilters + obj.Name,
					Type:        obj.Name + filterEqualsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldNotEquals,
					Description: descriptionInequalityFilters + obj.Name,
					Type:        obj.Name + filterEqualsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldGreaterThan,
					Description: descriptionGreaterThanFilters + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldSmallerThan,
					Description: descriptionSmallerThanFilters + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldGreaterOrEqual,
					Description: descriptionGreaterOrEqualFilters + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldSmallerOrEqual,
					Description: descriptionSmallerOrEqualFilters + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldContains,
					Description: descriptionContainsFilters + obj.Name,
					Type:        obj.Name + filterContainsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldNotContains,
					Description: descriptionNotContainsFilters + obj.Name,
					Type:        obj.Name + filterContainsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldLike,
					Description: descriptionLikeFilters + obj.Name,
					Type:        obj.Name + filterLikeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldNotLike,
					Description: descriptionNotLikeFilters + obj.Name,
					Type:        obj.Name + filterLikeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldNull,
					Description: descriptionNullFilters + obj.Name,
					Type:        obj.Name + filterNullSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldNotNull,
					Description: descriptionNotNullFilters + obj.Name,
					Type:        obj.Name + filterNullSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        filterFieldOrCondition,
					Description: descriptionOrCondition,
					Type:        FieldTypeBool,
					Modifiers:   []string{},
				},
				{
					Name:        filterFieldNestedFilters,
					Description: descriptionNestedFiltersTemplate + obj.Name + descriptionNestedFiltersSuffix,
					Type:        obj.Name + filterSuffix,
					Modifiers:   []string{ModifierArray},
				},
			},
		}
		result.Objects = append(result.Objects, mainFilter)

		// Generate FilterEquals object - contains all fields as nullable (used for both Equals and NotEquals)
		equalsFilter := Object{
			Name:        obj.Name + filterEqualsSuffix,
			Description: descriptionEqualityInequalityFilterFields + obj.Name,
			Fields:      make([]Field, 0, len(obj.Fields)),
		}
		for _, field := range obj.Fields {
			if isObjectType(field.Type, input.Objects) {
				// For nested objects, use the filter version
				equalsFilter.Fields = append(equalsFilter.Fields, generateNestedFilterField(field, filterEqualsSuffix, true, false, input.Objects))
			} else {
				// For primitive types, use the original field type
				equalsFilter.Fields = append(equalsFilter.Fields, generateFilterField(field, true, false))
			}
		}
		result.Objects = append(result.Objects, equalsFilter)

		// Generate FilterRange object - only comparable fields and nested objects
		rangeFilter := Object{
			Name:        obj.Name + filterRangeSuffix,
			Description: descriptionRangeFilterFields + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if isComparableType(field.Type) {
				// For comparable primitive types
				rangeFilter.Fields = append(rangeFilter.Fields, generateFilterField(field, true, false))
			} else if isObjectType(field.Type, input.Objects) {
				// For nested objects, include the filter version
				rangeFilter.Fields = append(rangeFilter.Fields, generateNestedFilterField(field, filterRangeSuffix, true, false, input.Objects))
			}
		}
		result.Objects = append(result.Objects, rangeFilter)

		// Generate FilterContains object - all fields except timestamps as arrays
		containsFilter := Object{
			Name:        obj.Name + filterContainsSuffix,
			Description: descriptionContainsFilterFields + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if field.Type != FieldTypeTimestamp {
				if isObjectType(field.Type, input.Objects) {
					// For nested objects, use the filter version (nullable, not array - arrays are for fields inside the nested filter)
					containsFilter.Fields = append(containsFilter.Fields, generateNestedFilterField(field, filterContainsSuffix, true, false, input.Objects))
				} else {
					// For primitive types, use the original field type
					containsFilter.Fields = append(containsFilter.Fields, generateFilterField(field, false, true))
				}
			}
		}
		result.Objects = append(result.Objects, containsFilter)

		// Generate FilterLike object - only string fields and nested objects
		likeFilter := Object{
			Name:        obj.Name + filterLikeSuffix,
			Description: descriptionLikeFilterFields + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if isStringType(field.Type) {
				// For string primitive types
				likeFilter.Fields = append(likeFilter.Fields, generateFilterField(field, true, false))
			} else if isObjectType(field.Type, input.Objects) {
				// For nested objects, include the filter version
				likeFilter.Fields = append(likeFilter.Fields, generateNestedFilterField(field, filterLikeSuffix, true, false, input.Objects))
			}
		}
		result.Objects = append(result.Objects, likeFilter)

		// Generate FilterNull object - only nullable fields or arrays
		nullFilter := Object{
			Name:        obj.Name + filterNullSuffix,
			Description: descriptionNullFilterFields + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if canBeNull(field) {
				if isObjectType(field.Type, input.Objects) {
					// For nested objects, use the filter version
					nestedNullField := generateNestedFilterField(field, filterNullSuffix, true, false, input.Objects)
					// But for null filters, we change the type to Bool to indicate null/not null
					nestedNullField.Type = FieldTypeBool
					nullFilter.Fields = append(nullFilter.Fields, nestedNullField)
				} else {
					// For primitive types, create a boolean field to indicate null/not null
					nullField := generateFilterField(Field{
						Name:        field.Name,
						Description: field.Description,
						Type:        FieldTypeBool,
					}, true, false)
					nullFilter.Fields = append(nullFilter.Fields, nullField)
				}
			}
		}
		result.Objects = append(result.Objects, nullFilter)
	}

	return result
}

// ResourceField methods

// HasCreateOperation checks if the ResourceField supports Create operations.
func (f ResourceField) HasCreateOperation() bool {
	return slices.Contains(f.Operations, OperationCreate)
}

// HasDeleteOperation checks if the ResourceField supports Delete operations.
func (f ResourceField) HasDeleteOperation() bool {
	return slices.Contains(f.Operations, OperationDelete)
}

// HasReadOperation checks if the ResourceField supports Read operations.
func (f ResourceField) HasReadOperation() bool {
	return slices.Contains(f.Operations, OperationRead)
}

// HasUpdateOperation checks if the ResourceField supports Update operations.
func (f ResourceField) HasUpdateOperation() bool {
	return slices.Contains(f.Operations, OperationUpdate)
}

// Field methods

// IsArray checks if the Field has the array modifier.
func (t Field) IsArray() bool {
	return slices.Contains(t.Modifiers, ModifierArray)
}

// IsNullable checks if the Field has the nullable modifier.
func (t Field) IsNullable() bool {
	return slices.Contains(t.Modifiers, ModifierNullable)
}

// TagJSON returns the JSON tag name for the field in camelCase.
func (t Field) TagJSON() string {
	return camelCase(t.Name)
}

// GetComment returns a formatted comment for the field.
func (t Field) GetComment(tabs string) string {
	return getComment(tabs, t.Description, t.Name)
}

// IsRequired determines if the field is required based on its properties and service context.
func (f Field) IsRequired(service *Service) bool {
	if f.IsNullable() {
		return false
	}

	if f.IsArray() {
		return false
	}

	if f.Default != "" {
		return false
	}

	if service.IsObject(f.Type) {
		return false
	}

	return true
}

// EndpointRequest methods

// GetRequiredBodyParams returns the names of required body parameters.
func (e EndpointRequest) GetRequiredBodyParams(service *Service) []string {
	requiredFields := make([]string, 0, len(e.BodyParams))

	for _, field := range e.BodyParams {
		if !field.IsRequired(service) {
			continue
		}

		requiredFields = append(requiredFields, field.Name)
	}

	return requiredFields
}

// Endpoint methods

// GetFullPath returns the full path for the endpoint including the resource name.
func (e Endpoint) GetFullPath(resourceName string) string {
	return pathSeparator + ToKebabCase(resourceName) + e.Path
}

// Resource methods

// HasCreateOperation checks if the Resource supports Create operations.
func (r Resource) HasCreateOperation() bool {
	return slices.Contains(r.Operations, OperationCreate)
}

// HasDeleteOperation checks if the Resource supports Delete operations.
func (r Resource) HasDeleteOperation() bool {
	return slices.Contains(r.Operations, OperationDelete)
}

// HasReadOperation checks if the Resource supports Read operations.
func (r Resource) HasReadOperation() bool {
	return slices.Contains(r.Operations, OperationRead)
}

// HasUpdateOperation checks if the Resource supports Update operations.
func (r Resource) HasUpdateOperation() bool {
	return slices.Contains(r.Operations, OperationUpdate)
}

// ShouldSkipAutoColumns checks if the Resource should skip generating auto columns.
func (r Resource) ShouldSkipAutoColumns() bool {
	return r.SkipAutoColumns
}

// GetPluralName returns the pluralized name of the resource.
func (r Resource) GetPluralName() string {
	return strmangle.Plural(r.Name)
}

// GetCreateBodyParams returns all fields that support Create operations.
func (r Resource) GetCreateBodyParams() []Field {
	bodyParams := make([]Field, 0)
	for _, resourceField := range r.Fields {
		if resourceField.HasCreateOperation() {
			// Convert ResourceField to Field by copying the embedded Field
			field := Field{
				Name:        resourceField.Name,
				Description: resourceField.Description,
				Type:        resourceField.Type,
				Default:     resourceField.Default,
				Example:     resourceField.Example,
				Modifiers:   make([]string, len(resourceField.Modifiers)),
			}
			copy(field.Modifiers, resourceField.Modifiers)
			bodyParams = append(bodyParams, field)
		}
	}
	return bodyParams
}

// GetUpdateBodyParams returns all fields that support Update operations.
func (r Resource) GetUpdateBodyParams() []Field {
	bodyParams := make([]Field, 0)
	for _, resourceField := range r.Fields {
		if resourceField.HasUpdateOperation() {
			// Convert ResourceField to Field by copying the embedded Field
			field := Field{
				Name:        resourceField.Name,
				Description: resourceField.Description,
				Type:        resourceField.Type,
				Default:     resourceField.Default,
				Example:     resourceField.Example,
				Modifiers:   make([]string, len(resourceField.Modifiers)),
			}
			copy(field.Modifiers, resourceField.Modifiers)
			bodyParams = append(bodyParams, field)
		}
	}
	return bodyParams
}

// GetReadableFields returns all fields that support Read operations.
func (r Resource) GetReadableFields() []Field {
	readableFields := make([]Field, 0)
	for _, resourceField := range r.Fields {
		if resourceField.HasReadOperation() {
			// Convert ResourceField to Field by copying the embedded Field
			field := Field{
				Name:        resourceField.Name,
				Description: resourceField.Description,
				Type:        resourceField.Type,
				Default:     resourceField.Default,
				Example:     resourceField.Example,
				Modifiers:   make([]string, len(resourceField.Modifiers)),
			}
			copy(field.Modifiers, resourceField.Modifiers)
			readableFields = append(readableFields, field)
		}
	}
	return readableFields
}

// HasEndpoint checks if the resource has an endpoint with the given name.
func (r Resource) HasEndpoint(name string) bool {
	for _, endpoint := range r.Endpoints {
		if endpoint.Name == name {
			return true
		}
	}
	return false
}

// Service methods

// IsObject checks if the given field type represents a custom object.
func (s *Service) IsObject(fieldType string) bool {
	return isObjectType(fieldType, s.Objects)
}

// HasObject checks if the service contains an object with the given name.
func (s *Service) HasObject(name string) bool {
	for _, obj := range s.Objects {
		if obj.Name == name {
			return true
		}
	}
	return false
}

// HasEnum checks if the service contains an enum with the given name.
func (s *Service) HasEnum(name string) bool {
	for _, enum := range s.Enums {
		if enum.Name == name {
			return true
		}
	}
	return false
}

// GetObject returns the object with the given name, or nil if not found.
func (s *Service) GetObject(name string) *Object {
	for _, obj := range s.Objects {
		if obj.Name == name {
			return &obj
		}
	}
	return nil
}

// Object methods

// HasField checks if the object contains a field with the given name.
func (o Object) HasField(name string) bool {
	for _, field := range o.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

// GetField returns the field with the given name, or nil if not found.
func (o Object) GetField(name string) *Field {
	for _, field := range o.Fields {
		if field.Name == name {
			return &field
		}
	}
	return nil
}

// Utility factory methods

// CreateLimitParam creates a standard limit parameter for pagination.
func CreateLimitParam() Field {
	return Field{
		Name:        listLimitParamName,
		Description: listLimitParamDesc,
		Type:        FieldTypeInt,
		Default:     listLimitDefaultValue,
	}
}

// CreateOffsetParam creates a standard offset parameter for pagination.
func CreateOffsetParam() Field {
	return Field{
		Name:        listOffsetParamName,
		Description: listOffsetParamDesc,
		Type:        FieldTypeInt,
		Default:     listOffsetDefaultValue,
	}
}

// CreatePaginationField creates a standard pagination field for responses.
func CreatePaginationField() Field {
	return Field{
		Name:        paginationObjectName,
		Description: "Pagination information",
		Type:        paginationObjectName,
	}
}

// CreateDataField creates a standard data field for array responses.
func CreateDataField(resourceName string) Field {
	return Field{
		Name:        "data",
		Description: fmt.Sprintf("Array of %s objects", resourceName),
		Type:        resourceName,
		Modifiers:   []string{ModifierArray},
	}
}

// CreateIDParam creates a standard ID parameter for path parameters.
func CreateIDParam(description string) Field {
	return Field{
		Name:        "id",
		Description: description,
		Type:        FieldTypeUUID,
	}
}

// CreateAutoColumnID creates a standard ID field for auto-columns.
func CreateAutoColumnID(resourceName string) Field {
	return Field{
		Name:        autoColumnIDName,
		Description: fmt.Sprintf(autoColumnIDDescTemplate, resourceName),
		Type:        FieldTypeUUID,
	}
}

// CreateAutoColumnCreatedAt creates a standard CreatedAt field for auto-columns.
func CreateAutoColumnCreatedAt(resourceName string) Field {
	return Field{
		Name:        autoColumnCreatedAtName,
		Description: fmt.Sprintf(autoColumnCreatedAtTemplate, resourceName),
		Type:        FieldTypeTimestamp,
	}
}

// CreateAutoColumnCreatedBy creates a standard CreatedBy field for auto-columns.
func CreateAutoColumnCreatedBy(resourceName string) Field {
	return Field{
		Name:        autoColumnCreatedByName,
		Description: fmt.Sprintf(autoColumnCreatedByTemplate, resourceName),
		Type:        FieldTypeUUID,
		Modifiers:   []string{ModifierNullable},
	}
}

// CreateAutoColumnUpdatedAt creates a standard UpdatedAt field for auto-columns.
func CreateAutoColumnUpdatedAt(resourceName string) Field {
	return Field{
		Name:        autoColumnUpdatedAtName,
		Description: fmt.Sprintf(autoColumnUpdatedAtTemplate, resourceName),
		Type:        FieldTypeTimestamp,
	}
}

// CreateAutoColumnUpdatedBy creates a standard UpdatedBy field for auto-columns.
func CreateAutoColumnUpdatedBy(resourceName string) Field {
	return Field{
		Name:        autoColumnUpdatedByName,
		Description: fmt.Sprintf(autoColumnUpdatedByTemplate, resourceName),
		Type:        FieldTypeUUID,
		Modifiers:   []string{ModifierNullable},
	}
}

// CreateAutoColumns creates all standard auto-column fields for resources.
func CreateAutoColumns(resourceName string) []Field {
	return []Field{
		CreateAutoColumnID(resourceName),
		CreateAutoColumnCreatedAt(resourceName),
		CreateAutoColumnCreatedBy(resourceName),
		CreateAutoColumnUpdatedAt(resourceName),
		CreateAutoColumnUpdatedBy(resourceName),
	}
}

// Helper functions

// getComment formats a comment string with proper indentation and prefixes.
func getComment(tabs string, description string, name string) string {
	comment := description

	if !strings.HasPrefix(description, name) {
		comment = name + nameDescSeparator + description
	}

	// Every new line should be prefixed with a //
	lines := strings.Split(comment, newlineChar)
	for i, line := range lines {
		lines[i] = tabs + commentPrefix + line
	}

	comment = strings.Join(lines, newlineChar)
	comment = strings.TrimSuffix(comment, tabs+newlineChar+commentPrefix)

	return comment
}

// camelCase converts a string to camelCase format.
func camelCase(s string) string {
	return strmangle.CamelCase(s)
}

// ToKebabCase converts a string to kebab-case format.
func ToKebabCase(s string) string {
	// Handle empty string
	if s == "" {
		return s
	}

	// First normalize spaces and underscores to hyphens, then handle PascalCase
	normalized := strings.ReplaceAll(s, "_", "-")
	normalized = strings.ReplaceAll(normalized, " ", "-")

	// Convert PascalCase/camelCase to kebab-case
	var result strings.Builder
	runes := []rune(normalized)

	for i, r := range runes {
		// Insert hyphen before uppercase letters in these cases:
		// 1. Before an uppercase letter that follows a lowercase letter or digit
		// 2. Before the last uppercase letter in a sequence of uppercase letters if followed by lowercase
		if i > 0 && r >= 'A' && r <= 'Z' && runes[i-1] != '-' {
			prev := runes[i-1]

			// Case 1: Previous character is lowercase or digit
			if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
				result.WriteByte('-')
			} else if prev >= 'A' && prev <= 'Z' {
				// Case 2: Previous character is uppercase, check if current is followed by lowercase
				if i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
					result.WriteByte('-')
				}
			}
		}
		result.WriteRune(r)
	}

	// Convert to lowercase
	return strings.ToLower(result.String())
}
