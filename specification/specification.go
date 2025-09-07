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
func ApplyOverlay(input *Service) *Service {
	if input == nil {
		return nil
	}

	// Create a deep copy of the input service
	result := &Service{
		Name:      input.Name,
		Enums:     make([]Enum, len(input.Enums)),
		Objects:   make([]Object, len(input.Objects)),
		Resources: make([]Resource, len(input.Resources)),
	}

	// Copy enums
	copy(result.Enums, input.Enums)

	// Copy existing objects
	copy(result.Objects, input.Objects)

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
			Description: "Filter object for " + obj.Name,
			Fields: []Field{
				{
					Name:        "Equals",
					Description: "Equality filters for " + obj.Name,
					Type:        obj.Name + filterEqualsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "NotEquals",
					Description: "Inequality filters for " + obj.Name,
					Type:        obj.Name + filterEqualsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "GreaterThan",
					Description: "Greater than filters for " + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "SmallerThan",
					Description: "Smaller than filters for " + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "GreaterOrEqual",
					Description: "Greater than or equal filters for " + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "SmallerOrEqual",
					Description: "Smaller than or equal filters for " + obj.Name,
					Type:        obj.Name + filterRangeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "Contains",
					Description: "Contains filters for " + obj.Name,
					Type:        obj.Name + filterContainsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "NotContains",
					Description: "Not contains filters for " + obj.Name,
					Type:        obj.Name + filterContainsSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "Like",
					Description: "LIKE filters for " + obj.Name,
					Type:        obj.Name + filterLikeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "NotLike",
					Description: "NOT LIKE filters for " + obj.Name,
					Type:        obj.Name + filterLikeSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "Null",
					Description: "Null filters for " + obj.Name,
					Type:        obj.Name + filterNullSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "NotNull",
					Description: "Not null filters for " + obj.Name,
					Type:        obj.Name + filterNullSuffix,
					Modifiers:   []string{ModifierNullable},
				},
				{
					Name:        "OrCondition",
					Description: "OrCondition decides if this filter is within an OR-condition or AND-condition",
					Type:        FieldTypeBool,
					Modifiers:   []string{},
				},
				{
					Name:        "NestedFilters",
					Description: "NestedFilters of the " + obj.Name + ", useful for more complex filters",
					Type:        obj.Name + filterSuffix,
					Modifiers:   []string{ModifierArray},
				},
			},
		}
		result.Objects = append(result.Objects, mainFilter)

		// Generate FilterEquals object - contains all fields as nullable (used for both Equals and NotEquals)
		equalsFilter := Object{
			Name:        obj.Name + filterEqualsSuffix,
			Description: "Equality/Inequality filter fields for " + obj.Name,
			Fields:      make([]Field, 0, len(obj.Fields)),
		}
		for _, field := range obj.Fields {
			equalsFilter.Fields = append(equalsFilter.Fields, generateFilterField(field, true, false))
		}
		result.Objects = append(result.Objects, equalsFilter)

		// Generate FilterRange object - only comparable fields
		rangeFilter := Object{
			Name:        obj.Name + filterRangeSuffix,
			Description: "Range filter fields for " + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if isComparableType(field.Type) {
				rangeFilter.Fields = append(rangeFilter.Fields, generateFilterField(field, true, false))
			}
		}
		result.Objects = append(result.Objects, rangeFilter)

		// Generate FilterContains object - all fields except timestamps as arrays
		containsFilter := Object{
			Name:        obj.Name + filterContainsSuffix,
			Description: "Contains filter fields for " + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if field.Type != FieldTypeTimestamp {
				containsFilter.Fields = append(containsFilter.Fields, generateFilterField(field, false, true))
			}
		}
		result.Objects = append(result.Objects, containsFilter)

		// Generate FilterLike object - only string fields
		likeFilter := Object{
			Name:        obj.Name + filterLikeSuffix,
			Description: "LIKE filter fields for " + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if isStringType(field.Type) {
				likeFilter.Fields = append(likeFilter.Fields, generateFilterField(field, true, false))
			}
		}
		result.Objects = append(result.Objects, likeFilter)

		// Generate FilterNull object - only nullable fields or arrays
		nullFilter := Object{
			Name:        obj.Name + filterNullSuffix,
			Description: "Null filter fields for " + obj.Name,
			Fields:      make([]Field, 0),
		}
		for _, field := range obj.Fields {
			if canBeNull(field) {
				nullField := generateFilterField(Field{
					Name:        field.Name,
					Description: field.Description,
					Type:        FieldTypeBool,
				}, true, false)
				nullFilter.Fields = append(nullFilter.Fields, nullField)
			}
		}
		result.Objects = append(result.Objects, nullFilter)
	}

	return result
}
