package specification

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aarondl/strmangle"
	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
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

// Default field examples for primitive types
const (
	defaultExampleUUID      = "123e4567-e89b-12d3-a456-426614174000"
	defaultExampleDate      = "2024-01-15"
	defaultExampleTimestamp = "2024-01-15T10:30:00Z"
	defaultExampleString    = "example"
	defaultExampleInt       = "42"
	defaultExampleBool      = "true"
)

// Field Modifiers
const (
	ModifierNullable = "Nullable"
	ModifierArray    = "Array"
)

// Retry Strategies
const (
	RetryStrategyBackoff = "backoff"
)

// Default retry configuration values
const (
	defaultRetryInitialInterval  = 500
	defaultRetryMaxInterval      = 60000
	defaultRetryMaxElapsedTime   = 3600000
	defaultRetryExponent         = 1.5
	defaultRetryStatusCodes      = "5XX"
	defaultRetryConnectionErrors = true
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
	totalFieldName              = "Total"
	totalFieldDescription       = "Total number of items available for pagination"
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

// Meta object constants
const (
	metaObjectName        = "Meta"
	metaObjectDescription = "Meta contains information about the creation and modification of a resource for auditing purposes"
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
	createEndpointName          = "Create"
	createEndpointPath          = ""
	createEndpointTitlePrefix   = "Create "
	createEndpointSummaryPrefix = "Create a new "
	createEndpointDescPrefix    = "Create a new "
	createResponseStatusCode    = 201
)

// Update Endpoint Constants
const (
	updateEndpointName          = "Update"
	updateEndpointPath          = "/{id}"
	updateEndpointTitlePrefix   = "Update "
	updateEndpointSummaryPrefix = "Update a "
	updateEndpointDescPrefix    = "Update a "
	updateResponseStatusCode    = 200
	updateIDParamName           = "id"
	updateIDParamDescTemplate   = "The unique identifier of the %s to update"
)

// Delete Endpoint Constants
const (
	deleteEndpointName          = "Delete"
	deleteEndpointPath          = "/{id}"
	deleteEndpointTitlePrefix   = "Delete "
	deleteEndpointSummaryPrefix = "Delete a "
	deleteEndpointDescPrefix    = "Delete a "
	deleteResponseStatusCode    = 204
	deleteIDParamName           = "id"
	deleteIDParamDescTemplate   = "The unique identifier of the %s to delete"
)

// Get Endpoint Constants
const (
	getEndpointName          = "Get"
	getEndpointPath          = "/{id}"
	getEndpointTitlePrefix   = "Retrieve an existing "
	getEndpointSummaryPrefix = "Get a "
	getResponseStatusCode    = 200
	getIDParamName           = "id"
	getIDParamDescTemplate   = "The unique identifier of the %s to retrieve"
)

// List Endpoint Constants
const (
	listEndpointName            = "List"
	listEndpointPath            = ""
	listEndpointTitlePrefix     = "List all "
	listEndpointSummaryPrefix   = "List "
	listEndpointDescTemplate    = "Returns a paginated list of all `%s` in your organization."
	listResponseStatusCode      = 200
	listLimitParamName          = "limit"
	listLimitParamDesc          = "The maximum number of items to return (default: 50)"
	listLimitParamDescTemplate  = "The maximum number of %s to return (default: 50) when listing %s"
	listLimitDefaultValue       = "50"
	listLimitExampleValue       = "1"
	listOffsetParamName         = "offset"
	listOffsetParamDesc         = "The number of items to skip before starting to return results (default: 0)"
	listOffsetParamDescTemplate = "The number of %s to skip before starting to return results (default: 0) when listing %s"
	listOffsetDefaultValue      = "0"
	listOffsetExampleValue      = "0"
)

// Search Endpoint Constants
const (
	searchEndpointName            = "Search"
	searchEndpointPath            = "/_search"
	searchEndpointTitlePrefix     = "Search "
	searchEndpointSummaryPrefix   = "Search "
	searchEndpointDescTemplate    = "Search for `%s` with filtering capabilities."
	searchResponseStatusCode      = 200
	searchFilterParamName         = "Filter"
	searchFilterParamDesc         = "Filter criteria to search for specific records"
	searchLimitParamDescTemplate  = "The maximum number of %s to return (default: 50) when searching %s"
	searchOffsetParamDescTemplate = "The number of %s to skip before starting to return results (default: 0) when searching %s"
)

// Response Description Constants
const (
	createResponseDescTemplate = "Successfully created the %s"
	updateResponseDescTemplate = "Successfully updated the %s"
	deleteResponseDescTemplate = "Successfully deleted the %s"
	getResponseDescTemplate    = "Successfully retrieved the %s"
	listResponseDescTemplate   = "Successfully retrieved the list of %s"
	searchResponseDescTemplate = "Successfully searched for %s"
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

// File parsing constants
const (
	errorInvalidFile       = "invalid input file"
	errorUnsupportedFormat = "unsupported file format"
	errorFileRead          = "failed to read file"
	errorFileParse         = "failed to parse file"

	// Validation error constants
	errorInvalidOperation = "invalid operation"
	errorInvalidFieldType = "invalid field type"
	errorInvalidModifier  = "invalid modifier"
	errorValidationFailed = "validation failed"
	errorYAMLParsing      = "YAML parsing failed"
)

// File extension constants
const (
	extYAML = ".yaml"
	extYML  = ".yml"
	extJSON = ".json"
)

// ServiceServer represents a server in the API service.
type ServiceServer struct {
	// URL of the server
	URL string `json:"url"`

	// Description of the server
	Description string `json:"description,omitempty"`

	// ID is a unique identifier for the server for SDK generation
	ID string `json:"id,omitempty"`
}

// ServiceContact represents the contact information for the API service.
type ServiceContact struct {
	// Name of the contact person/organization
	Name string `json:"name,omitempty"`

	// URL pointing to the contact information
	URL string `json:"url,omitempty"`

	// Email address of the contact person/organization
	Email string `json:"email,omitempty"`
}

// ServiceLicense represents the license information for the API service.
type ServiceLicense struct {
	// Name of the license used for the API (required)
	Name string `json:"name"`

	// URL pointing to the license used for the API
	URL string `json:"url,omitempty"`

	// SPDX license identifier for the license used for the API
	Identifier string `json:"identifier,omitempty"`
}

// SecurityScheme represents a security scheme definition.
type SecurityScheme struct {
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
}

// SecurityRequirement represents scheme names that must be satisfied together.
type SecurityRequirement []string

// RetryBackoffConfiguration defines the backoff behavior for retry attempts.
type RetryBackoffConfiguration struct {
	// InitialInterval is the initial interval between retries in milliseconds
	InitialInterval int `json:"initial_interval"`

	// MaxInterval is the maximum interval between retries in milliseconds
	MaxInterval int `json:"max_interval"`

	// MaxElapsedTime is the maximum total time for retry attempts in milliseconds
	MaxElapsedTime int `json:"max_elapsed_time"`

	// Exponent is the multiplier for exponential backoff
	Exponent float64 `json:"exponent"`
}

// RetryConfiguration defines the retry behavior for API calls.
type RetryConfiguration struct {
	// Strategy defines the retry strategy (e.g., "backoff")
	Strategy string `json:"strategy"`

	// Backoff configuration for exponential backoff strategy
	Backoff RetryBackoffConfiguration `json:"backoff"`

	// StatusCodes defines which HTTP status codes should trigger retries (e.g., "5XX")
	StatusCodes []string `json:"status_codes"`

	// RetryConnectionErrors indicates whether to retry on connection errors
	RetryConnectionErrors bool `json:"retry_connection_errors"`
}

// TimeoutConfiguration defines the timeout behavior for API calls.
type TimeoutConfiguration struct {
	// Timeout is the request timeout in milliseconds
	Timeout int `json:"timeout"`
}

// Service is the definition of an API service.
type Service struct {
	// Name of the service
	Name string `json:"name"`

	// Version of the service
	Version string `json:"version,omitempty"`

	// Contact information for the service
	Contact *ServiceContact `json:"contact,omitempty"`

	// License information for the service
	License *ServiceLicense `json:"license,omitempty"`

	// Servers that are part of the service
	Servers []ServiceServer `json:"servers,omitempty"`

	// SecuritySchemes defines available security schemes
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`

	// Security defines global security requirements (OR logic between requirements)
	Security []SecurityRequirement `json:"security,omitempty"`

	// Retry configuration for the service
	Retry *RetryConfiguration `json:"retry,omitempty"`

	// Timeout configuration for the service
	Timeout *TimeoutConfiguration `json:"timeout,omitempty"`

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

	// Summary is a short plain text description of the endpoint
	// For example: "Create a new school", "Update an existing school", "Delete a school"...
	Summary string `json:"summary"`

	// Description of the endpoint (can be in markdown format for longer explanations)
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

	// Description of the response
	Description string `json:"description"`

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
		Name:            input.Name,
		Version:         input.Version,
		Contact:         input.Contact,                               // Copy contact information
		License:         input.License,                               // Copy license information
		Servers:         append([]ServiceServer{}, input.Servers...), // Copy servers slice
		SecuritySchemes: input.SecuritySchemes,                       // Copy security schemes
		Security:        input.Security,                              // Copy security requirements
		Retry:           input.Retry,                                 // Copy retry configuration
		Timeout:         input.Timeout,                               // Copy timeout configuration
		Enums:           make([]Enum, 0, len(input.Enums)+2),         // +2 for ErrorCode and ErrorFieldCode enums
		Objects:         make([]Object, 0, len(input.Objects)+3),     // +3 for Error, ErrorField, and Pagination objects
		Resources:       make([]Resource, len(input.Resources)),
	}

	// Add default enums and objects if they don't already exist
	addDefaultEnumsAndObjects(result, input)

	// Copy resources
	copy(result.Resources, input.Resources)

	// Generate Objects and endpoints from Resources
	generateObjectsFromResources(result, input.Resources)
	// Generate filter objects for resources that have Read operations (needed for search endpoints)
	generateFilterObjectsForSearchableResources(result, input.Resources)
	generateEndpointsFromResources(result, input.Resources)

	// Generate RequestError objects for types used in body parameters
	// This happens at the end to ensure all objects and endpoints are generated first
	generateRequestErrorObjectsForBodyParams(result)

	return result
}

// addDefaultEnumsAndObjects adds the default error, pagination, and meta objects to the service if they don't already exist.
func addDefaultEnumsAndObjects(result *Service, input *Service) {
	// Check if ErrorCode enum, Error object, ErrorFieldCode enum, ErrorField object, Pagination object, and Meta object already exist
	errorCodeEnumExists := false
	errorObjectExists := false
	errorFieldCodeEnumExists := false
	errorFieldObjectExists := false
	paginationObjectExists := false
	metaObjectExists := false
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
		if object.Name == metaObjectName {
			metaObjectExists = true
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
					Example:     "0",
				},
				{
					Name:        limitFieldName,
					Description: limitFieldDescription,
					Type:        FieldTypeInt,
					Example:     "1",
				},
				{
					Name:        totalFieldName,
					Description: totalFieldDescription,
					Type:        FieldTypeInt,
					Example:     "100",
				},
			},
		}
		result.Objects = append(result.Objects, paginationObject)
	}

	// Add default Meta object if it doesn't exist
	if !metaObjectExists {
		metaObject := createDefaultMeta()
		result.Objects = append(result.Objects, metaObject)
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
					autoColumns := createAutoColumnsWithMeta(resource.Name)
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
			Summary:     createEndpointSummaryPrefix + resource.Name,
			Description: createEndpointDescPrefix + resource.Name,
			Method:      httpMethodPost,
			Path:        createEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{}, bodyParams),
			Response:    createStandardResponse(createResponseStatusCode, fmt.Sprintf(createResponseDescTemplate, resourceName), &resourceName),
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
			Summary:     updateEndpointSummaryPrefix + resource.Name,
			Description: updateEndpointDescPrefix + resource.Name,
			Method:      httpMethodPatch,
			Path:        updateEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, bodyParams),
			Response:    createStandardResponse(updateResponseStatusCode, fmt.Sprintf(updateResponseDescTemplate, resourceName), &resourceName),
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
			Summary:     deleteEndpointSummaryPrefix + resource.Name,
			Description: deleteEndpointDescPrefix + resource.Name,
			Method:      httpMethodDelete,
			Path:        deleteEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, []Field{}),
			Response:    createStandardResponse(deleteResponseStatusCode, fmt.Sprintf(deleteResponseDescTemplate, resource.Name), nil), // No body object for delete
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
			Summary:     getEndpointSummaryPrefix + resource.Name,
			Description: fmt.Sprintf("Retrieves the `%s` with the given ID.", resource.Name),
			Method:      httpMethodGet,
			Path:        getEndpointPath,
			Request:     createStandardRequest([]Field{idParam}, []Field{}, []Field{}),
			Response:    createStandardResponse(getResponseStatusCode, fmt.Sprintf(getResponseDescTemplate, resourceName), &resourceName),
		}

		addEndpointToResource(result, resource.Name, getEndpoint)
	}
}

// generateListEndpoint generates a List endpoint for resources that have Read operations.
func generateListEndpoint(result *Service, resource Resource) {
	if resource.HasReadOperation() && !resource.HasEndpoint(listEndpointName) {
		limitParam := createListLimitParamForResource(resource)
		offsetParam := createListOffsetParamForResource(resource)
		paginationField := createPaginationField()
		dataField := createDataField(resource.Name)
		pluralResourceName := resource.GetPluralName()

		listEndpoint := Endpoint{
			Name:        listEndpointName,
			Title:       listEndpointTitlePrefix + pluralResourceName,
			Summary:     listEndpointSummaryPrefix + pluralResourceName,
			Description: fmt.Sprintf(listEndpointDescTemplate, pluralResourceName),
			Method:      httpMethodGet,
			Path:        listEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{limitParam, offsetParam}, []Field{}),
			Response:    createListResponse(listResponseStatusCode, fmt.Sprintf(listResponseDescTemplate, pluralResourceName), dataField, paginationField),
		}

		addEndpointToResource(result, resource.Name, listEndpoint)
	}
}

// generateSearchEndpoint generates a Search endpoint for resources that have Read operations.
func generateSearchEndpoint(result *Service, resource Resource) {
	if resource.HasReadOperation() && !resource.HasEndpoint(searchEndpointName) {
		limitParam := createSearchLimitParamForResource(resource)
		offsetParam := createSearchOffsetParamForResource(resource)
		filterParam := Field{
			Name:        searchFilterParamName,
			Description: searchFilterParamDesc,
			Type:        resource.Name + filterSuffix,
		}
		paginationField := createPaginationField()
		dataField := createDataField(resource.Name)
		pluralResourceName := resource.GetPluralName()

		searchEndpoint := Endpoint{
			Name:        searchEndpointName,
			Title:       searchEndpointTitlePrefix + pluralResourceName,
			Summary:     searchEndpointSummaryPrefix + pluralResourceName,
			Description: fmt.Sprintf(searchEndpointDescTemplate, pluralResourceName),
			Method:      httpMethodPost,
			Path:        searchEndpointPath,
			Request:     createStandardRequest([]Field{}, []Field{limitParam, offsetParam}, []Field{filterParam}),
			Response:    createListResponse(searchResponseStatusCode, fmt.Sprintf(searchResponseDescTemplate, pluralResourceName), dataField, paginationField),
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

// createStandardResponse creates a standard endpoint response with the given status code, description, and optional body object.
func createStandardResponse(statusCode int, description string, bodyObject *string) EndpointResponse {
	return EndpointResponse{
		ContentType: contentTypeJSON,
		StatusCode:  statusCode,
		Description: description,
		Headers:     []Field{},
		BodyFields:  []Field{},
		BodyObject:  bodyObject,
	}
}

// createListResponse creates a standard list endpoint response with pagination and data fields.
func createListResponse(statusCode int, description string, dataField Field, paginationField Field) EndpointResponse {
	return EndpointResponse{
		ContentType: contentTypeJSON,
		StatusCode:  statusCode,
		Description: description,
		Headers:     []Field{},
		BodyFields:  []Field{dataField, paginationField},
		BodyObject:  nil,
	}
}

// collectTypesUsedInBodyParams collects all types (including nested) used in request body parameters.
func collectTypesUsedInBodyParams(service *Service) map[string]bool {
	usedTypes := make(map[string]bool)

	// Collect types from all endpoint body parameters, excluding search endpoints
	for _, resource := range service.Resources {
		for _, endpoint := range resource.Endpoints {
			// Skip search endpoints since they don't use errorFields
			isSearchEndpoint := endpoint.Name == searchEndpointName || endpoint.Name == "AdvancedSearch"
			if !isSearchEndpoint {
				for _, bodyParam := range endpoint.Request.BodyParams {
					collectTypeRecursively(bodyParam.Type, usedTypes, service.Objects)
				}
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
			if obj.IsFilter() {
				continue // Do not generate RequestError objects for filter objects
			}

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
				// Skip generating RequestError objects for search endpoints since they don't use errorFields
				isSearchEndpoint := endpoint.Name == searchEndpointName || endpoint.Name == "AdvancedSearch"
				if !isSearchEndpoint {
					requestErrorName := resource.Name + endpoint.Name + requestErrorSuffix
					requestErrorDescription := requestErrorDescriptionPrefix + resource.Name + " " + endpoint.Name + " endpoint"
					requestError := generateRequestErrorObject(requestErrorName, requestErrorDescription, endpoint.Request.BodyParams, service.Objects)
					service.Objects = append(service.Objects, requestError)
				}
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

// generateFilterObjectsForSearchableResources generates filter objects for resources that have Read operations.
// This ensures filter objects exist before search endpoints are generated.
func generateFilterObjectsForSearchableResources(service *Service, resources []Resource) {
	for _, resource := range resources {
		// Only generate filter objects for resources that have Read operations (which will have search endpoints)
		if resource.HasReadOperation() {
			// Find the corresponding object for this resource
			for _, obj := range service.Objects {
				if obj.Name == resource.Name {
					// Generate all filter objects for this resource object
					filterObjects := generateFilterObjectsForObject(obj, service.Objects)
					service.Objects = append(service.Objects, filterObjects...)
					break
				}
			}
		}
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
		Name:            input.Name,
		Version:         input.Version,
		Contact:         input.Contact,                               // Copy contact information
		License:         input.License,                               // Copy license information
		Servers:         append([]ServiceServer{}, input.Servers...), // Copy servers slice
		SecuritySchemes: input.SecuritySchemes,                       // Copy security schemes
		Security:        input.Security,                              // Copy security requirements
		Retry:           input.Retry,                                 // Copy retry configuration
		Timeout:         input.Timeout,                               // Copy timeout configuration
		Enums:           make([]Enum, len(input.Enums)),
		Objects:         make([]Object, 0, len(input.Objects)*7), // Estimate for filter objects
		Resources:       make([]Resource, len(input.Resources)),
	}

	// Copy enums
	copy(result.Enums, input.Enums)

	// Copy existing objects first
	result.Objects = append(result.Objects, input.Objects...)

	// Copy resources
	copy(result.Resources, input.Resources)

	// Collect types used in request body parameters
	usedTypes := collectTypesUsedInBodyParams(result)

	// Generate Filter objects only for Objects used in request body parameters
	for _, obj := range input.Objects {
		// Skip objects that are not used in request body parameters
		if !usedTypes[obj.Name] {
			continue
		}
		// Generate all filter objects for this object
		filterObjects := generateFilterObjectsForObject(obj, input.Objects)
		result.Objects = append(result.Objects, filterObjects...)
	}

	return result
}

// generateFilterObjectsForObject generates all filter-related objects for a given object.
func generateFilterObjectsForObject(obj Object, allObjects []Object) []Object {
	var filterObjects []Object

	// Generate main filter object
	mainFilter := generateMainFilterObject(obj)
	filterObjects = append(filterObjects, mainFilter)

	// Generate specialized filter objects
	equalsFilter := generateEqualsFilterObject(obj, allObjects)
	rangeFilter := generateRangeFilterObject(obj, allObjects)
	containsFilter := generateContainsFilterObject(obj, allObjects)
	likeFilter := generateLikeFilterObject(obj, allObjects)
	nullFilter := generateNullFilterObject(obj, allObjects)

	filterObjects = append(filterObjects, equalsFilter, rangeFilter, containsFilter, likeFilter, nullFilter)

	return filterObjects
}

// generateMainFilterObject creates the main filter object with all filter type references.
func generateMainFilterObject(obj Object) Object {
	return Object{
		Name:        obj.Name + filterSuffix,
		Description: descriptionFilterObject + obj.Name,
		Fields: []Field{
			createFilterField(filterFieldEquals, descriptionEqualityFilters+obj.Name, obj.Name+filterEqualsSuffix, true),
			createFilterField(filterFieldNotEquals, descriptionInequalityFilters+obj.Name, obj.Name+filterEqualsSuffix, true),
			createFilterField(filterFieldGreaterThan, descriptionGreaterThanFilters+obj.Name, obj.Name+filterRangeSuffix, true),
			createFilterField(filterFieldSmallerThan, descriptionSmallerThanFilters+obj.Name, obj.Name+filterRangeSuffix, true),
			createFilterField(filterFieldGreaterOrEqual, descriptionGreaterOrEqualFilters+obj.Name, obj.Name+filterRangeSuffix, true),
			createFilterField(filterFieldSmallerOrEqual, descriptionSmallerOrEqualFilters+obj.Name, obj.Name+filterRangeSuffix, true),
			createFilterField(filterFieldContains, descriptionContainsFilters+obj.Name, obj.Name+filterContainsSuffix, true),
			createFilterField(filterFieldNotContains, descriptionNotContainsFilters+obj.Name, obj.Name+filterContainsSuffix, true),
			createFilterField(filterFieldLike, descriptionLikeFilters+obj.Name, obj.Name+filterLikeSuffix, true),
			createFilterField(filterFieldNotLike, descriptionNotLikeFilters+obj.Name, obj.Name+filterLikeSuffix, true),
			createFilterField(filterFieldNull, descriptionNullFilters+obj.Name, obj.Name+filterNullSuffix, true),
			createFilterField(filterFieldNotNull, descriptionNotNullFilters+obj.Name, obj.Name+filterNullSuffix, true),
			createFilterField(filterFieldOrCondition, descriptionOrCondition, FieldTypeBool, false),
			createFilterField(filterFieldNestedFilters, descriptionNestedFiltersTemplate+obj.Name+descriptionNestedFiltersSuffix, obj.Name+filterSuffix, false, ModifierArray),
		},
	}
}

// createFilterField is a helper to create filter fields with consistent structure.
func createFilterField(name, description, fieldType string, nullable bool, extraModifiers ...string) Field {
	modifiers := make([]string, 0, 2)
	if nullable {
		modifiers = append(modifiers, ModifierNullable)
	}
	modifiers = append(modifiers, extraModifiers...)

	return Field{
		Name:        name,
		Description: description,
		Type:        fieldType,
		Modifiers:   modifiers,
	}
}

// generateEqualsFilterObject generates the FilterEquals object for equality comparisons.
func generateEqualsFilterObject(obj Object, allObjects []Object) Object {
	fields := processFieldsForFilter(obj.Fields, allObjects, func(field Field, objects []Object) (Field, bool) {
		// All fields are included in equals filter
		if isObjectType(field.Type, objects) {
			return generateNestedFilterField(field, filterEqualsSuffix, true, false, objects), true
		}
		return generateFilterField(field, true, false), true
	})

	return Object{
		Name:        obj.Name + filterEqualsSuffix,
		Description: descriptionEqualityInequalityFilterFields + obj.Name,
		Fields:      fields,
	}
}

// generateRangeFilterObject generates the FilterRange object for range comparisons.
func generateRangeFilterObject(obj Object, allObjects []Object) Object {
	fields := processFieldsForFilter(obj.Fields, allObjects, func(field Field, objects []Object) (Field, bool) {
		// Only comparable types and nested objects
		if isComparableType(field.Type) {
			return generateFilterField(field, true, false), true
		}
		if isObjectType(field.Type, objects) {
			return generateNestedFilterField(field, filterRangeSuffix, true, false, objects), true
		}
		return Field{}, false
	})

	return Object{
		Name:        obj.Name + filterRangeSuffix,
		Description: descriptionRangeFilterFields + obj.Name,
		Fields:      fields,
	}
}

// generateContainsFilterObject generates the FilterContains object for contains operations.
func generateContainsFilterObject(obj Object, allObjects []Object) Object {
	fields := processFieldsForFilter(obj.Fields, allObjects, func(field Field, objects []Object) (Field, bool) {
		// All fields except timestamps
		if field.Type == FieldTypeTimestamp {
			return Field{}, false
		}
		if isObjectType(field.Type, objects) {
			return generateNestedFilterField(field, filterContainsSuffix, true, false, objects), true
		}
		return generateFilterField(field, false, true), true
	})

	return Object{
		Name:        obj.Name + filterContainsSuffix,
		Description: descriptionContainsFilterFields + obj.Name,
		Fields:      fields,
	}
}

// generateLikeFilterObject generates the FilterLike object for LIKE operations.
func generateLikeFilterObject(obj Object, allObjects []Object) Object {
	fields := processFieldsForFilter(obj.Fields, allObjects, func(field Field, objects []Object) (Field, bool) {
		// Only string types and nested objects
		if isStringType(field.Type) {
			return generateFilterField(field, true, false), true
		}
		if isObjectType(field.Type, objects) {
			return generateNestedFilterField(field, filterLikeSuffix, true, false, objects), true
		}
		return Field{}, false
	})

	return Object{
		Name:        obj.Name + filterLikeSuffix,
		Description: descriptionLikeFilterFields + obj.Name,
		Fields:      fields,
	}
}

// generateNullFilterObject generates the FilterNull object for null checks.
func generateNullFilterObject(obj Object, allObjects []Object) Object {
	fields := processFieldsForFilter(obj.Fields, allObjects, func(field Field, objects []Object) (Field, bool) {
		// Only nullable fields or arrays
		if !canBeNull(field) {
			return Field{}, false
		}

		if isObjectType(field.Type, objects) {
			// For nested objects, create boolean field for null check
			nestedField := generateNestedFilterField(field, filterNullSuffix, true, false, objects)
			nestedField.Type = FieldTypeBool
			return nestedField, true
		}
		// For primitive types, create boolean field
		return generateFilterField(Field{
			Name:        field.Name,
			Description: field.Description,
			Type:        FieldTypeBool,
		}, true, false), true
	})

	return Object{
		Name:        obj.Name + filterNullSuffix,
		Description: descriptionNullFilterFields + obj.Name,
		Fields:      fields,
	}
}

// processFieldsForFilter is a generic helper that processes fields with a custom filter function.
func processFieldsForFilter(fields []Field, allObjects []Object, filterFunc func(Field, []Object) (Field, bool)) []Field {
	var result []Field
	for _, field := range fields {
		if filterField, include := filterFunc(field, allObjects); include {
			result = append(result, filterField)
		}
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
	return CamelCase(t.Name)
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

// getDefaultExample returns the default example for a primitive field type.
func getDefaultExample(fieldType string) string {
	switch fieldType {
	case FieldTypeUUID:
		return defaultExampleUUID
	case FieldTypeDate:
		return defaultExampleDate
	case FieldTypeTimestamp:
		return defaultExampleTimestamp
	case FieldTypeString:
		return defaultExampleString
	case FieldTypeInt:
		return defaultExampleInt
	case FieldTypeBool:
		return defaultExampleBool
	default:
		slog.Warn("no default example available for field type, consider adding support", "fieldType", fieldType)
		return ""
	}
}

// ensureExample ensures that the field has an example, setting a default one for primitive types if none exists.
func (f *Field) ensureExample() {
	// Only set default examples for primitive types and only if no example already exists
	if f.Example == "" && isPrimitiveType(f.Type) {
		f.Example = getDefaultExample(f.Type)
	}
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
	return pathSeparator + toKebabCase(resourceName) + e.Path
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
	return r.getFieldsByOperation(ResourceField.HasCreateOperation)
}

// GetUpdateBodyParams returns all fields that support Update operations.
func (r Resource) GetUpdateBodyParams() []Field {
	return r.getFieldsByOperation(ResourceField.HasUpdateOperation)
}

// GetReadableFields returns all fields that support Read operations.
func (r Resource) GetReadableFields() []Field {
	return r.getFieldsByOperation(ResourceField.HasReadOperation)
}

// getFieldsByOperation is a helper method that filters ResourceFields by operation and converts them to Fields.
func (r Resource) getFieldsByOperation(operationCheck func(ResourceField) bool) []Field {
	var result []Field
	for _, resourceField := range r.Fields {
		if operationCheck(resourceField) {
			field := r.convertResourceFieldToField(resourceField)
			result = append(result, field)
		}
	}
	return result
}

// convertResourceFieldToField converts a ResourceField to a Field by copying the embedded Field data.
func (r Resource) convertResourceFieldToField(resourceField ResourceField) Field {
	field := Field{
		Name:        resourceField.Name,
		Description: resourceField.Description,
		Type:        resourceField.Type,
		Default:     resourceField.Default,
		Example:     resourceField.Example,
		Modifiers:   make([]string, len(resourceField.Modifiers)),
	}
	copy(field.Modifiers, resourceField.Modifiers)
	field.ensureExample()
	return field
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

// IsFilter checks if the object is a filter object.
func (o Object) IsFilter() bool {
	return strings.HasSuffix(o.Name, filterSuffix) ||
		strings.HasSuffix(o.Name, filterEqualsSuffix) ||
		strings.HasSuffix(o.Name, filterRangeSuffix) ||
		strings.HasSuffix(o.Name, filterContainsSuffix) ||
		strings.HasSuffix(o.Name, filterLikeSuffix) ||
		strings.HasSuffix(o.Name, filterNullSuffix)
}

// Utility factory methods

// createLimitParam creates a standard limit parameter for pagination.
func createLimitParam() Field {
	return Field{
		Name:        listLimitParamName,
		Description: listLimitParamDesc,
		Type:        FieldTypeInt,
		Default:     listLimitDefaultValue,
		Example:     listLimitExampleValue,
	}
}

// createListLimitParamForResource creates a limit parameter with resource-specific description for List operations.
func createListLimitParamForResource(resource Resource) Field {
	pluralResourceName := resource.GetPluralName()
	resourceSpecificDescription := fmt.Sprintf(listLimitParamDescTemplate, pluralResourceName, pluralResourceName)
	return Field{
		Name:        listLimitParamName,
		Description: resourceSpecificDescription,
		Type:        FieldTypeInt,
		Default:     listLimitDefaultValue,
		Example:     listLimitExampleValue,
	}
}

// createSearchLimitParamForResource creates a limit parameter with resource-specific description for Search operations.
func createSearchLimitParamForResource(resource Resource) Field {
	pluralResourceName := resource.GetPluralName()
	resourceSpecificDescription := fmt.Sprintf(searchLimitParamDescTemplate, pluralResourceName, pluralResourceName)
	return Field{
		Name:        listLimitParamName,
		Description: resourceSpecificDescription,
		Type:        FieldTypeInt,
		Default:     listLimitDefaultValue,
		Example:     listLimitExampleValue,
	}
}

// createOffsetParam creates a standard offset parameter for pagination.
func createOffsetParam() Field {
	return Field{
		Name:        listOffsetParamName,
		Description: listOffsetParamDesc,
		Type:        FieldTypeInt,
		Default:     listOffsetDefaultValue,
		Example:     listOffsetExampleValue,
	}
}

// createListOffsetParamForResource creates an offset parameter with resource-specific description for List operations.
func createListOffsetParamForResource(resource Resource) Field {
	pluralResourceName := resource.GetPluralName()
	resourceSpecificDescription := fmt.Sprintf(listOffsetParamDescTemplate, pluralResourceName, pluralResourceName)
	return Field{
		Name:        listOffsetParamName,
		Description: resourceSpecificDescription,
		Type:        FieldTypeInt,
		Default:     listOffsetDefaultValue,
		Example:     listOffsetExampleValue,
	}
}

// createSearchOffsetParamForResource creates an offset parameter with resource-specific description for Search operations.
func createSearchOffsetParamForResource(resource Resource) Field {
	pluralResourceName := resource.GetPluralName()
	resourceSpecificDescription := fmt.Sprintf(searchOffsetParamDescTemplate, pluralResourceName, pluralResourceName)
	return Field{
		Name:        listOffsetParamName,
		Description: resourceSpecificDescription,
		Type:        FieldTypeInt,
		Default:     listOffsetDefaultValue,
		Example:     listOffsetExampleValue,
	}
}

// createPaginationField creates a standard pagination field for responses.
func createPaginationField() Field {
	return Field{
		Name:        paginationObjectName,
		Description: "Pagination information",
		Type:        paginationObjectName,
	}
}

// createDataField creates a standard data field for array responses.
func createDataField(resourceName string) Field {
	return Field{
		Name:        "data",
		Description: fmt.Sprintf("Array of %s objects", resourceName),
		Type:        resourceName,
		Modifiers:   []string{ModifierArray},
	}
}

// createIDParam creates a standard ID parameter for path parameters.
func createIDParam(description string) Field {
	return Field{
		Name:        "id",
		Description: description,
		Type:        FieldTypeUUID,
	}
}

// createAutoColumnID creates a standard ID field for auto-columns.
func createAutoColumnID(resourceName string) Field {
	return Field{
		Name:        autoColumnIDName,
		Description: fmt.Sprintf(autoColumnIDDescTemplate, resourceName),
		Type:        FieldTypeUUID,
		Example:     "123e4567-e89b-12d3-a456-426614174000",
	}
}

// createAutoColumnCreatedAt creates a standard CreatedAt field for auto-columns.
func createAutoColumnCreatedAt(resourceName string) Field {
	return Field{
		Name:        autoColumnCreatedAtName,
		Description: fmt.Sprintf(autoColumnCreatedAtTemplate, resourceName),
		Type:        FieldTypeTimestamp,
		Example:     "2024-01-15T10:30:00Z",
	}
}

// createAutoColumnCreatedBy creates a standard CreatedBy field for auto-columns.
func createAutoColumnCreatedBy(resourceName string) Field {
	return Field{
		Name:        autoColumnCreatedByName,
		Description: fmt.Sprintf(autoColumnCreatedByTemplate, resourceName),
		Type:        FieldTypeUUID,
		Modifiers:   []string{ModifierNullable},
		Example:     "987fcdeb-51a2-43d1-b567-123456789abc",
	}
}

// createAutoColumnUpdatedAt creates a standard UpdatedAt field for auto-columns.
func createAutoColumnUpdatedAt(resourceName string) Field {
	return Field{
		Name:        autoColumnUpdatedAtName,
		Description: fmt.Sprintf(autoColumnUpdatedAtTemplate, resourceName),
		Type:        FieldTypeTimestamp,
		Example:     "2024-01-15T14:45:00Z",
	}
}

// createAutoColumnUpdatedBy creates a standard UpdatedBy field for auto-columns.
func createAutoColumnUpdatedBy(resourceName string) Field {
	return Field{
		Name:        autoColumnUpdatedByName,
		Description: fmt.Sprintf(autoColumnUpdatedByTemplate, resourceName),
		Type:        FieldTypeUUID,
		Modifiers:   []string{ModifierNullable},
		Example:     "987fcdeb-51a2-43d1-b567-123456789abc",
	}
}

// createAutoColumns creates all standard auto-column fields for resources.
func createAutoColumns(resourceName string) []Field {
	return []Field{
		createAutoColumnID(resourceName),
		createAutoColumnCreatedAt(resourceName),
		createAutoColumnCreatedBy(resourceName),
		createAutoColumnUpdatedAt(resourceName),
		createAutoColumnUpdatedBy(resourceName),
	}
}

// createDefaultMeta creates a standard Meta object containing creation and update metadata fields.
func createDefaultMeta() Object {
	return Object{
		Name:        metaObjectName,
		Description: metaObjectDescription,
		Fields: []Field{
			{
				Name:        autoColumnCreatedAtName,
				Description: "Timestamp when the resource was created",
				Type:        FieldTypeTimestamp,
				Example:     "2024-01-15T10:30:00Z",
			},
			{
				Name:        autoColumnCreatedByName,
				Description: "User who created the resource",
				Type:        FieldTypeUUID,
				Modifiers:   []string{ModifierNullable},
				Example:     "987fcdeb-51a2-43d1-b567-123456789abc",
			},
			{
				Name:        autoColumnUpdatedAtName,
				Description: "Timestamp when the resource was last updated",
				Type:        FieldTypeTimestamp,
				Example:     "2024-01-15T14:45:00Z",
			},
			{
				Name:        autoColumnUpdatedByName,
				Description: "User who last updated the resource",
				Type:        FieldTypeUUID,
				Modifiers:   []string{ModifierNullable},
				Example:     "987fcdeb-51a2-43d1-b567-123456789abc",
			},
		},
	}
}

// createAutoColumnsWithMeta creates auto-column fields using Meta object for metadata fields.
func createAutoColumnsWithMeta(resourceName string) []Field {
	return []Field{
		createAutoColumnID(resourceName),
		{
			Name:        metaObjectName,
			Description: fmt.Sprintf("Metadata information for the %s", resourceName),
			Type:        metaObjectName,
		},
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

// CamelCase converts a string to camelCase format.
// Special cases:
// - "ID" becomes "id" instead of "iD"
// - Consecutive capital letters at the start are lowercased (e.g., "CSNSchoolCode" -> "csnSchoolCode")
func CamelCase(s string) string {
	if s == "ID" {
		return "id"
	}

	result := strmangle.CamelCase(s)

	// Handle consecutive capital letters at the beginning
	// Convert sequences like "cSNSchool" to "csnSchool"
	if len(result) > 1 {
		runes := []rune(result)

		// Find the end of consecutive uppercase letters at the beginning (after first char)
		consecutiveEnd := 1
		for consecutiveEnd < len(runes) && runes[consecutiveEnd] >= 'A' && runes[consecutiveEnd] <= 'Z' {
			consecutiveEnd++
		}

		// If we have consecutive uppercase letters, convert them to lowercase
		// except possibly the last one if it's followed by lowercase letters
		if consecutiveEnd > 2 {
			// If the sequence is followed by lowercase letters, keep the last uppercase letter
			// Example: "cSNSchool" -> "csnSchool" (convert "SN" to "sn", keep "S" before "chool")
			if consecutiveEnd < len(runes) && runes[consecutiveEnd] >= 'a' && runes[consecutiveEnd] <= 'z' {
				consecutiveEnd--
			}

			// Convert consecutive uppercase letters to lowercase
			for i := 1; i < consecutiveEnd; i++ {
				if runes[i] >= 'A' && runes[i] <= 'Z' {
					runes[i] = runes[i] - 'A' + 'a'
				}
			}

			result = string(runes)
		}
	}

	return result
}

// toKebabCase converts a string to kebab-case format.
func toKebabCase(s string) string {
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

// Validation functions

// ValidationError represents a validation error with position information.
type ValidationError struct {
	Message string
	Line    int
	Column  int
	Path    string
}

func (e *ValidationError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("validation error at line %d, column %d (%s): %s", e.Line, e.Column, e.Path, e.Message)
	}
	return fmt.Sprintf("validation error (%s): %s", e.Path, e.Message)
}

// ValidateServiceWithPosition validates a service using YAML node position information.
func ValidateServiceWithPosition(data []byte, fileExtension string) error {
	if fileExtension != extYAML && fileExtension != extYML {
		// For non-YAML files, fall back to regular validation
		var service Service
		if fileExtension == extJSON {
			if err := json.Unmarshal(data, &service); err != nil {
				return fmt.Errorf("%s: JSON parsing error: %w", errorFileParse, err)
			}
		}
		return validateService(&service)
	}

	// Parse YAML into struct first
	var service Service
	if err := yaml.Unmarshal(data, &service); err != nil {
		return fmt.Errorf("%s: YAML parsing error: %w", errorFileParse, err)
	}

	// Parse YAML into AST to get position information for errors
	file, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("%s: %w", errorYAMLParsing, err)
	}

	// Create a position tracker to map validation errors to line numbers
	positionTracker := newPositionTracker(file)

	// Run validation and enhance errors with position information
	if validationErr := validateService(&service); validationErr != nil {
		enhancedErr := positionTracker.EnhanceError(validationErr, &service)
		return enhancedErr
	}

	return nil
}

// PositionTracker helps map validation errors to line numbers in YAML.
type PositionTracker struct {
	file *ast.File
}

// newPositionTracker creates a new position tracker from a parsed YAML file.
func newPositionTracker(file *ast.File) *PositionTracker {
	return &PositionTracker{file: file}
}

// EnhanceError enhances a validation error with position information.
func (pt *PositionTracker) EnhanceError(err error, service *Service) error {
	if err == nil || pt.file == nil {
		return err
	}

	// Try to extract meaningful error information and find corresponding YAML position
	errorMsg := err.Error()

	// Handle different error patterns
	if strings.Contains(errorMsg, "invalid operation") {
		return pt.findOperationError(errorMsg, service)
	}
	if strings.Contains(errorMsg, "invalid field type") {
		return pt.findFieldTypeError(errorMsg, service)
	}
	if strings.Contains(errorMsg, "invalid modifier") {
		return pt.findModifierError(errorMsg, service)
	}

	return err
}

// findOperationError tries to find line number for operation validation errors.
func (pt *PositionTracker) findOperationError(errorMsg string, service *Service) error {
	// Extract the invalid operation from the error message
	// Error format: "resource 0 (Students): resource operations: invalid operation: operation 'create' must be one of: [Create Read Update Delete]"
	if pos := pt.findInvalidValue("operations", errorMsg); pos != nil {
		return &ValidationError{
			Message: errorMsg,
			Line:    pos.Line,
			Column:  pos.Column,
			Path:    "operations",
		}
	}
	// If we couldn't find a position, still return the original error message
	return fmt.Errorf("%s", errorMsg)
}

// findFieldTypeError tries to find line number for field type validation errors.
func (pt *PositionTracker) findFieldTypeError(errorMsg string, service *Service) error {
	// Extract the invalid type from error message
	if pos := pt.findInvalidValue("type", errorMsg); pos != nil {
		return &ValidationError{
			Message: errorMsg,
			Line:    pos.Line,
			Column:  pos.Column,
			Path:    "type",
		}
	}
	return fmt.Errorf("%s", errorMsg)
}

// findModifierError tries to find line number for modifier validation errors.
func (pt *PositionTracker) findModifierError(errorMsg string, service *Service) error {
	// Extract the invalid modifier from error message
	if pos := pt.findInvalidValue("modifiers", errorMsg); pos != nil {
		return &ValidationError{
			Message: errorMsg,
			Line:    pos.Line,
			Column:  pos.Column,
			Path:    "modifiers",
		}
	}
	return fmt.Errorf("%s", errorMsg)
}

// findInvalidValue searches for the first occurrence of a value in the YAML that might be invalid.
func (pt *PositionTracker) findInvalidValue(fieldName string, errorMsg string) *token.Position {
	// This is a simplified implementation that finds the first occurrence
	// of specific field names in the YAML and returns their position
	if len(pt.file.Docs) > 0 {
		return pt.findFieldInAST(fieldName, pt.file.Docs[0].Body)
	}
	return nil
}

// findFieldInAST recursively searches for a field name in the AST.
func (pt *PositionTracker) findFieldInAST(fieldName string, node ast.Node) *token.Position {
	if node == nil {
		return nil
	}

	nodeType := node.Type()
	if nodeType.String() == "Mapping" {
		if mappingNode, ok := node.(*ast.MappingNode); ok {
			for _, mappingValue := range mappingNode.Values {
				// Check if this is the field we're looking for
				if keyNode := mappingValue.Key; keyNode != nil && keyNode.Type().String() == "String" {
					if stringNode, ok := keyNode.(*ast.StringNode); ok {
						keyValue := stringNode.GetToken().Value
						if keyValue == fieldName {
							// Return position of the key (where the field name is)
							if token := stringNode.GetToken(); token != nil && token.Position != nil {
								return token.Position
							}
						}
					}
				}
				// Recursively search in the value
				if pos := pt.findFieldInAST(fieldName, mappingValue.Value); pos != nil {
					return pos
				}
			}
		}
	} else if nodeType.String() == "Sequence" {
		if sequenceNode, ok := node.(*ast.SequenceNode); ok {
			for _, entry := range sequenceNode.Values {
				if pos := pt.findFieldInAST(fieldName, entry); pos != nil {
					return pos
				}
			}
		}
	}

	return nil
}

// validateService validates the entire service specification against the defined rules.
func validateService(service *Service) error {
	// Validate retry configuration
	if service.Retry != nil {
		if err := validateRetryConfiguration(service.Retry); err != nil {
			return fmt.Errorf("retry configuration: %w", err)
		}
	}

	// Validate resources
	for i, resource := range service.Resources {
		if err := validateResource(service, &resource); err != nil {
			return fmt.Errorf("resource %d (%s): %w", i, resource.Name, err)
		}
	}

	// Validate objects
	for i, object := range service.Objects {
		if err := validateObject(service, &object); err != nil {
			return fmt.Errorf("object %d (%s): %w", i, object.Name, err)
		}
	}

	return nil
}

// validateResource validates a resource and its fields against the defined rules.
func validateResource(service *Service, resource *Resource) error {
	// Validate operations
	if err := validateOperations(resource.Operations); err != nil {
		return fmt.Errorf("resource operations: %w", err)
	}

	// Validate resource fields
	for i, field := range resource.Fields {
		if err := validateResourceField(service, &field); err != nil {
			return fmt.Errorf("field %d (%s): %w", i, field.Name, err)
		}
	}

	// Validate endpoints
	for i, endpoint := range resource.Endpoints {
		if err := validateEndpoint(service, &endpoint); err != nil {
			return fmt.Errorf("endpoint %d (%s): %w", i, endpoint.Name, err)
		}
	}

	return nil
}

// validateObject validates an object and its fields against the defined rules.
func validateObject(service *Service, object *Object) error {
	// Validate object fields
	for i, field := range object.Fields {
		if err := validateField(service, &field); err != nil {
			return fmt.Errorf("field %d (%s): %w", i, field.Name, err)
		}
	}

	return nil
}

// validateResourceField validates a resource field against the defined rules.
func validateResourceField(service *Service, field *ResourceField) error {
	// Validate operations
	if err := validateOperations(field.Operations); err != nil {
		return fmt.Errorf("field operations: %w", err)
	}

	// Validate the embedded field
	return validateField(service, &field.Field)
}

// validateField validates a field against the defined rules.
func validateField(service *Service, field *Field) error {
	// Validate field type
	if err := validateFieldType(service, field.Type); err != nil {
		return fmt.Errorf("field type: %w", err)
	}

	// Validate modifiers
	if err := validateModifiers(field.Modifiers); err != nil {
		return fmt.Errorf("field modifiers: %w", err)
	}

	return nil
}

// validateEndpoint validates an endpoint against the defined rules.
func validateEndpoint(service *Service, endpoint *Endpoint) error {
	// Validate request body params
	for i, field := range endpoint.Request.BodyParams {
		if err := validateField(service, &field); err != nil {
			return fmt.Errorf("request body param %d (%s): %w", i, field.Name, err)
		}
	}

	// Validate request query params
	for i, field := range endpoint.Request.QueryParams {
		if err := validateField(service, &field); err != nil {
			return fmt.Errorf("request query param %d (%s): %w", i, field.Name, err)
		}
	}

	// Validate request headers
	for i, field := range endpoint.Request.Headers {
		if err := validateField(service, &field); err != nil {
			return fmt.Errorf("request header %d (%s): %w", i, field.Name, err)
		}
	}

	// Validate response body fields
	for i, field := range endpoint.Response.BodyFields {
		if err := validateField(service, &field); err != nil {
			return fmt.Errorf("response body field %d (%s): %w", i, field.Name, err)
		}
	}

	return nil
}

// validateRetryConfiguration validates a retry configuration against the defined rules.
func validateRetryConfiguration(retry *RetryConfiguration) error {
	if retry == nil {
		return nil
	}

	// Validate strategy
	if retry.Strategy != "" && retry.Strategy != RetryStrategyBackoff {
		return fmt.Errorf("retry strategy '%s' must be 'backoff'", retry.Strategy)
	}

	// Validate backoff configuration
	if retry.Backoff.InitialInterval < 0 {
		return fmt.Errorf("retry backoff initial interval must be non-negative, got: %d", retry.Backoff.InitialInterval)
	}

	if retry.Backoff.MaxInterval < 0 {
		return fmt.Errorf("retry backoff max interval must be non-negative, got: %d", retry.Backoff.MaxInterval)
	}

	if retry.Backoff.MaxElapsedTime < 0 {
		return fmt.Errorf("retry backoff max elapsed time must be non-negative, got: %d", retry.Backoff.MaxElapsedTime)
	}

	if retry.Backoff.Exponent < 0 {
		return fmt.Errorf("retry backoff exponent must be non-negative, got: %f", retry.Backoff.Exponent)
	}

	// Validate that max interval is not less than initial interval (if both are specified)
	if retry.Backoff.InitialInterval > 0 && retry.Backoff.MaxInterval > 0 &&
		retry.Backoff.InitialInterval > retry.Backoff.MaxInterval {
		return fmt.Errorf("retry backoff initial interval (%d) cannot be greater than max interval (%d)",
			retry.Backoff.InitialInterval, retry.Backoff.MaxInterval)
	}

	// Validate status codes
	for _, statusCode := range retry.StatusCodes {
		if statusCode == "" {
			return fmt.Errorf("retry status codes cannot contain empty strings")
		}
		// Allow common patterns like "5XX", "4XX", or specific codes like "429", "503"
		if !isValidStatusCode(statusCode) {
			return fmt.Errorf("retry status code '%s' is not valid", statusCode)
		}
	}

	return nil
}

// isValidStatusCode checks if a status code is valid (either a pattern like "5XX" or specific code).
func isValidStatusCode(code string) bool {
	// Allow patterns like "5XX", "4XX", "3XX", "2XX", "1XX"
	if len(code) == 3 && (code[1] == 'X' && code[2] == 'X') {
		first := code[0]
		return first >= '1' && first <= '5'
	}

	// Allow specific HTTP status codes (100-599)
	if len(code) == 3 {
		for _, r := range code {
			if r < '0' || r > '9' {
				return false
			}
		}
		// Convert to int and check range
		var statusCode int
		for i, r := range code {
			if i == 0 {
				statusCode = int(r-'0') * 100
			} else if i == 1 {
				statusCode += int(r-'0') * 10
			} else {
				statusCode += int(r - '0')
			}
		}
		return statusCode >= 100 && statusCode <= 599
	}

	return false
}

// validateOperations validates that all operations are in PascalCase and are valid CRUD operations.
func validateOperations(operations []string) error {
	validOperations := []string{OperationCreate, OperationRead, OperationUpdate, OperationDelete}

	for _, operation := range operations {
		if !slices.Contains(validOperations, operation) {
			return fmt.Errorf("%s: operation '%s' must be one of: %v", errorInvalidOperation, operation, validOperations)
		}
	}

	return nil
}

// validateFieldType validates that the field type is either a valid primitive type, enum, or object.
func validateFieldType(service *Service, fieldType string) error {
	// Check if it's a primitive type
	validPrimitiveTypes := []string{
		FieldTypeUUID, FieldTypeDate, FieldTypeTimestamp,
		FieldTypeString, FieldTypeInt, FieldTypeBool,
	}

	if slices.Contains(validPrimitiveTypes, fieldType) {
		return nil
	}

	// Check if it's a valid enum
	if service.HasEnum(fieldType) {
		return nil
	}

	// Check if it's a valid object
	if service.HasObject(fieldType) {
		return nil
	}

	return fmt.Errorf("%s: field type '%s' must be one of the primitive types %v, or a valid enum/object", errorInvalidFieldType, fieldType, validPrimitiveTypes)
}

// validateModifiers validates that all modifiers are in PascalCase and are valid modifiers.
func validateModifiers(modifiers []string) error {
	validModifiers := []string{ModifierNullable, ModifierArray}

	for _, modifier := range modifiers {
		if !slices.Contains(validModifiers, modifier) {
			return fmt.Errorf("%s: modifier '%s' must be one of: %v", errorInvalidModifier, modifier, validModifiers)
		}
	}

	return nil
}

// Parsing functions

// ParseServiceFromFile reads and parses a YAML or JSON specification file,
// automatically applying overlays to ensure complete specification.
func ParseServiceFromFile(filePath string) (*Service, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: file does not exist: %s", errorInvalidFile, filePath)
	}

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFileRead, err)
	}

	// Parse based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	service, err := parseServiceFromBytes(data, ext)
	if err != nil {
		return nil, err
	}

	// Apply overlays to ensure complete specification
	return applyDefaultOverlays(service), nil
}

// ParseServiceFromBytes parses a service from byte data and file extension,
// automatically applying overlays to ensure complete specification.
func ParseServiceFromBytes(data []byte, fileExtension string) (*Service, error) {
	service, err := parseServiceFromBytes(data, fileExtension)
	if err != nil {
		return nil, err
	}

	// Apply overlays to ensure complete specification
	return applyDefaultOverlays(service), nil
}

// ParseServiceFromJSON parses a service from JSON data,
// automatically applying overlays to ensure complete specification.
func ParseServiceFromJSON(data []byte) (*Service, error) {
	service, err := parseServiceFromBytes(data, extJSON)
	if err != nil {
		return nil, err
	}

	// Apply overlays to ensure complete specification
	return applyDefaultOverlays(service), nil
}

// ParseServiceFromYAML parses a service from YAML data,
// automatically applying overlays to ensure complete specification.
func ParseServiceFromYAML(data []byte) (*Service, error) {
	service, err := parseServiceFromBytes(data, extYAML)
	if err != nil {
		return nil, err
	}

	// Apply overlays to ensure complete specification
	return applyDefaultOverlays(service), nil
}

// parseServiceFromBytes is the internal parsing function without overlay application.
func parseServiceFromBytes(data []byte, fileExtension string) (*Service, error) {
	// Validate with position information first
	if err := ValidateServiceWithPosition(data, fileExtension); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	// If validation passed, parse normally
	var service Service

	switch fileExtension {
	case extYAML, extYML:
		if err := yaml.Unmarshal(data, &service); err != nil {
			return nil, fmt.Errorf("%s: YAML parsing error: %w", errorFileParse, err)
		}
	case extJSON:
		if err := json.Unmarshal(data, &service); err != nil {
			return nil, fmt.Errorf("%s: JSON parsing error: %w", errorFileParse, err)
		}
	default:
		return nil, fmt.Errorf("%s: file must have .yaml, .yml, or .json extension", errorUnsupportedFormat)
	}

	return &service, nil
}

// ensureAllFieldsHaveExamples ensures that all fields in the service have examples set.
// This applies default examples to primitive field types that don't already have examples.
func ensureAllFieldsHaveExamples(service *Service) {
	if service == nil {
		return
	}

	// Apply to object fields
	for i := range service.Objects {
		for j := range service.Objects[i].Fields {
			service.Objects[i].Fields[j].ensureExample()
		}
	}

	// Apply to resource fields
	for i := range service.Resources {
		for j := range service.Resources[i].Fields {
			service.Resources[i].Fields[j].Field.ensureExample()
		}
		// Also apply to endpoint fields
		for j := range service.Resources[i].Endpoints {
			endpoint := &service.Resources[i].Endpoints[j]
			for k := range endpoint.Request.PathParams {
				endpoint.Request.PathParams[k].ensureExample()
			}
			for k := range endpoint.Request.QueryParams {
				endpoint.Request.QueryParams[k].ensureExample()
			}
			for k := range endpoint.Request.BodyParams {
				endpoint.Request.BodyParams[k].ensureExample()
			}
			for k := range endpoint.Request.Headers {
				endpoint.Request.Headers[k].ensureExample()
			}
			for k := range endpoint.Response.BodyFields {
				endpoint.Response.BodyFields[k].ensureExample()
			}
			for k := range endpoint.Response.Headers {
				endpoint.Response.Headers[k].ensureExample()
			}
		}
	}
}

// applyDefaultOverlays applies the standard overlays to a service to ensure
// resource objects and CRUD endpoints are generated.
func applyDefaultOverlays(service *Service) *Service {
	// Apply overlay to generate resource objects and standard endpoints
	overlayedService := ApplyOverlay(service)
	if overlayedService == nil {
		// If overlay application fails, return original service
		return service
	}

	// Apply filter overlay to generate filter objects
	finalService := ApplyFilterOverlay(overlayedService)
	if finalService == nil {
		// If filter overlay application fails, return overlayed service
		return overlayedService
	}

	// Ensure all fields have examples
	ensureAllFieldsHaveExamples(finalService)

	return finalService
}

// HasRetryConfiguration checks if the service has retry configuration defined.
func (s Service) HasRetryConfiguration() bool {
	return s.Retry != nil
}

// GetRetryConfigurationWithDefaults returns the retry configuration with default values applied.
func (s Service) GetRetryConfigurationWithDefaults() RetryConfiguration {
	if s.Retry == nil {
		return createDefaultRetryConfiguration()
	}

	config := *s.Retry

	// Apply defaults for missing values
	if config.Strategy == "" {
		config.Strategy = RetryStrategyBackoff
	}

	if config.Backoff.InitialInterval == 0 {
		config.Backoff.InitialInterval = defaultRetryInitialInterval
	}

	if config.Backoff.MaxInterval == 0 {
		config.Backoff.MaxInterval = defaultRetryMaxInterval
	}

	if config.Backoff.MaxElapsedTime == 0 {
		config.Backoff.MaxElapsedTime = defaultRetryMaxElapsedTime
	}

	if config.Backoff.Exponent == 0 {
		config.Backoff.Exponent = defaultRetryExponent
	}

	if len(config.StatusCodes) == 0 {
		config.StatusCodes = []string{defaultRetryStatusCodes}
	}

	// Note: RetryConnectionErrors is a bool, so we keep the zero value (false) if not set
	// Users must explicitly set it to true if they want connection error retries

	return config
}

// createDefaultRetryConfiguration creates a retry configuration with all default values.
func createDefaultRetryConfiguration() RetryConfiguration {
	return RetryConfiguration{
		Strategy: RetryStrategyBackoff,
		Backoff: RetryBackoffConfiguration{
			InitialInterval: defaultRetryInitialInterval,
			MaxInterval:     defaultRetryMaxInterval,
			MaxElapsedTime:  defaultRetryMaxElapsedTime,
			Exponent:        defaultRetryExponent,
		},
		StatusCodes:           []string{defaultRetryStatusCodes},
		RetryConnectionErrors: defaultRetryConnectionErrors,
	}
}
