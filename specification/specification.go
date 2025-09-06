package specification

// Service is the definition of an API service.
type Service struct {
	// Name of the service
	Name string `json:"name" yaml:"name"`

	// Enums that are used in the service
	Enums []Enum `json:"enums" yaml:"enums"`

	// Objects that are used in the service
	Objects []Object `json:"objects" yaml:"objects"`

	// Resources that are part of the service
	Resources []Resource `json:"resources" yaml:"resources"`
}

// Enum represents an enumeration with possible values.
type Enum struct {
	// Name of the enum
	Name string `json:"name" yaml:"name"`

	// Description of the enum
	Description string `json:"description" yaml:"description"`

	// Values that are possible for the enum
	Values []EnumValue `json:"values" yaml:"values"`
}

// EnumValue represents a single value in an enumeration.
type EnumValue struct {
	// Name of the enum value, for example Male for the Enum Gender - should be unique in the Enum
	Name string `json:"name" yaml:"name"`

	// Description for the enum value
	Description string `json:"description" yaml:"description"`
}

// Object is a shared object within the service,
// can be used by multiple resources.
type Object struct {
	// Name of the object, should be unique in the service
	Name string `json:"name" yaml:"name"`

	// Description about the object
	Description string `json:"description" yaml:"description"`

	// Fields in the object
	Fields []Field `json:"fields" yaml:"fields"`
}

// Resource represents a resource in the API with its operations and fields.
type Resource struct {
	// Name of the resource, should be unique within the service
	Name string `json:"name" yaml:"name"`

	// Description about the resource
	Description string `json:"description" yaml:"description"`

	// Operations that are allowed for the resource can be all of Create, Update, Read, Delete
	Operations []string `json:"operations" yaml:"operations"`

	// Fields of the resource
	Fields []ResourceField `json:"fields" yaml:"fields"`

	// Endpoints of the resource
	Endpoints []Endpoint `json:"endpoints" yaml:"endpoints"`
}

// Field contains information about a field within an endpoint or resource or Object.
type Field struct {
	// Name of the field, should be unique in the Resource or Object or Endpoint
	Name string `json:"name" yaml:"name"`

	// Description of the field, explain the reason what it is used for and why it's needed
	Description string `json:"description" yaml:"description"`

	// Type of the field, can be one of the types (UUID, Date, Timestamp, String, Int, Bool) or one of the custom Objects
	Type string `json:"type" yaml:"type"`

	// Default value of the field
	Default string `json:"default,omitempty" yaml:"default,omitempty"`

	// Example value of the field
	Example string `json:"example,omitempty" yaml:"example,omitempty"`

	// Modifiers of the field, can be nullable or array
	Modifiers []string `json:"modifiers,omitempty" yaml:"modifiers,omitempty"`
}

// ResourceField is used within a resource it extends the field with an operations configuration.
type ResourceField struct {
	Field `json:",inline" yaml:",inline"`

	// Operations that the field is allowed in (Create,Update,Delete,Read)
	Operations []string `json:"operations" yaml:"operations"`
}

// Endpoint represents an API endpoint within a resource.
type Endpoint struct {
	// Name of the endpoint, should be unique within the resource
	Name string `json:"name" yaml:"name"`

	// Description of the endpoint
	Description string `json:"description" yaml:"description"`

	// HTTP method for the endpoint (GET, POST, PUT, DELETE, PATCH)
	Method string `json:"method" yaml:"method"`

	// Path pattern for the endpoint
	Path string `json:"path" yaml:"path"`

	// Request body fields for the endpoint
	RequestFields []Field `json:"request_fields,omitempty" yaml:"request_fields,omitempty"`

	// Response fields for the endpoint
	ResponseFields []Field `json:"response_fields,omitempty" yaml:"response_fields,omitempty"`

	// Query parameters for the endpoint
	QueryParameters []Field `json:"query_parameters,omitempty" yaml:"query_parameters,omitempty"`

	// Path parameters for the endpoint
	PathParameters []Field `json:"path_parameters,omitempty" yaml:"path_parameters,omitempty"`
}