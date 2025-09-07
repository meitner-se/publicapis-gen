package specification

import "github.com/meitner-se/publicapis-gen/constants"

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
		if containsOperation(resource.Operations, constants.OperationRead) {
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
					if containsOperation(resourceField.Operations, constants.OperationRead) {
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
