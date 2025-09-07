package specification

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

// Service is the definition of an API service.
type Service struct {
	// Name of the service
	Name string `json:"name"`

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

// containsOperation checks if a slice of operations contains a specific operation.
func containsOperation(operations []string, operation string) bool {
	for _, op := range operations {
		if op == operation {
			return true
		}
	}
	return false
}

// ApplyOverlay applies an overlay to a specification, generating Objects from Resources.
// It creates Objects for Resources that have the "Read" operation, including all fields
// that support the "Read" operation in the generated Object.
// It also adds default ErrorCode enum, Error object, ErrorFieldCode enum, and ErrorField object to every service.
func ApplyOverlay(input *Service) *Service {
	if input == nil {
		return nil
	}

	// Create a deep copy of the input service
	result := &Service{
		Name:      input.Name,
		Enums:     make([]Enum, 0, len(input.Enums)+2),     // +2 for ErrorCode and ErrorFieldCode enums
		Objects:   make([]Object, 0, len(input.Objects)+2), // +2 for Error and ErrorField objects
		Resources: make([]Resource, len(input.Resources)),
	}

	// Check if ErrorCode enum, Error object, ErrorFieldCode enum, and ErrorField object already exist
	errorCodeEnumExists := false
	errorObjectExists := false
	errorFieldCodeEnumExists := false
	errorFieldObjectExists := false
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

	// Copy resources
	copy(result.Resources, input.Resources)

	// Generate Objects from Resources that have Read operations
	for _, resource := range input.Resources {
		// Check if the resource has Read operation
		if containsOperation(resource.Operations, OperationRead) {
			// Check if an object with this name already exists
			objectExists := false
			for _, existingObj := range result.Objects {
				if existingObj.Name == resource.Name {
					objectExists = true
					break
				}
			}

			// Only create the object if it doesn't already exist
			if !objectExists {
				// Create a new Object based on the Resource
				newObject := Object{
					Name:        resource.Name,
					Description: resource.Description,
					Fields:      make([]Field, 0),
				}

				// Add all fields that support Read operation
				for _, resourceField := range resource.Fields {
					if containsOperation(resourceField.Operations, OperationRead) {
						// Convert ResourceField to Field by copying the embedded Field
						field := Field{
							Name:        resourceField.Field.Name,
							Description: resourceField.Field.Description,
							Type:        resourceField.Field.Type,
							Default:     resourceField.Field.Default,
							Example:     resourceField.Field.Example,
							Modifiers:   make([]string, len(resourceField.Field.Modifiers)),
						}
						copy(field.Modifiers, resourceField.Field.Modifiers)
						newObject.Fields = append(newObject.Fields, field)
					}
				}

				// Add the new object to the result
				result.Objects = append(result.Objects, newObject)
			}
		}
	}

	return result
}

// containsModifier checks if a slice of modifiers contains a specific modifier.
func containsModifier(modifiers []string, modifier string) bool {
	for _, mod := range modifiers {
		if mod == modifier {
			return true
		}
	}
	return false
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
	return containsModifier(field.Modifiers, ModifierNullable) || containsModifier(field.Modifiers, ModifierArray)
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
