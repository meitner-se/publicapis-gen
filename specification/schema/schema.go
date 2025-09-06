package schema

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"

	"github.com/meitner-se/publicapis-gen/specification"
)

// SchemaGenerator provides functionality to generate JSON schemas from specification structs.
type SchemaGenerator struct {
	reflector *jsonschema.Reflector
}

// NewSchemaGenerator creates a new schema generator with default configuration.
func NewSchemaGenerator() *SchemaGenerator {
	r := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            false,
		ExpandedStruct:            true,
	}

	return &SchemaGenerator{
		reflector: r,
	}
}

// GenerateServiceSchema generates a JSON schema for the Service struct.
func (sg *SchemaGenerator) GenerateServiceSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Service{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Service")
	}

	return schema, nil
}

// GenerateEnumSchema generates a JSON schema for the Enum struct.
func (sg *SchemaGenerator) GenerateEnumSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Enum{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Enum")
	}

	return schema, nil
}

// GenerateObjectSchema generates a JSON schema for the Object struct.
func (sg *SchemaGenerator) GenerateObjectSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Object{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Object")
	}

	return schema, nil
}

// GenerateResourceSchema generates a JSON schema for the Resource struct.
func (sg *SchemaGenerator) GenerateResourceSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Resource{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Resource")
	}

	return schema, nil
}

// GenerateFieldSchema generates a JSON schema for the Field struct.
func (sg *SchemaGenerator) GenerateFieldSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Field{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Field")
	}

	return schema, nil
}

// GenerateResourceFieldSchema generates a JSON schema for the ResourceField struct.
func (sg *SchemaGenerator) GenerateResourceFieldSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.ResourceField{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for ResourceField")
	}

	return schema, nil
}

// GenerateEndpointSchema generates a JSON schema for the Endpoint struct.
func (sg *SchemaGenerator) GenerateEndpointSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Endpoint{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for Endpoint")
	}

	return schema, nil
}

// GenerateEndpointRequestSchema generates a JSON schema for the EndpointRequest struct.
func (sg *SchemaGenerator) GenerateEndpointRequestSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.EndpointRequest{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for EndpointRequest")
	}

	return schema, nil
}

// GenerateEndpointResponseSchema generates a JSON schema for the EndpointResponse struct.
func (sg *SchemaGenerator) GenerateEndpointResponseSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.EndpointResponse{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for EndpointResponse")
	}

	return schema, nil
}

// GenerateAllSchemas generates JSON schemas for all main specification structs.
func (sg *SchemaGenerator) GenerateAllSchemas() (map[string]*jsonschema.Schema, error) {
	schemas := make(map[string]*jsonschema.Schema)

	serviceSchema, err := sg.GenerateServiceSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Service schema: %w", err)
	}
	schemas["Service"] = serviceSchema

	enumSchema, err := sg.GenerateEnumSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Enum schema: %w", err)
	}
	schemas["Enum"] = enumSchema

	objectSchema, err := sg.GenerateObjectSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Object schema: %w", err)
	}
	schemas["Object"] = objectSchema

	resourceSchema, err := sg.GenerateResourceSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Resource schema: %w", err)
	}
	schemas["Resource"] = resourceSchema

	fieldSchema, err := sg.GenerateFieldSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Field schema: %w", err)
	}
	schemas["Field"] = fieldSchema

	resourceFieldSchema, err := sg.GenerateResourceFieldSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ResourceField schema: %w", err)
	}
	schemas["ResourceField"] = resourceFieldSchema

	endpointSchema, err := sg.GenerateEndpointSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Endpoint schema: %w", err)
	}
	schemas["Endpoint"] = endpointSchema

	endpointRequestSchema, err := sg.GenerateEndpointRequestSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate EndpointRequest schema: %w", err)
	}
	schemas["EndpointRequest"] = endpointRequestSchema

	endpointResponseSchema, err := sg.GenerateEndpointResponseSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate EndpointResponse schema: %w", err)
	}
	schemas["EndpointResponse"] = endpointResponseSchema

	return schemas, nil
}

// SchemaToJSON converts a JSON schema to a JSON string.
func (sg *SchemaGenerator) SchemaToJSON(schema *jsonschema.Schema) (string, error) {
	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// GenerateServiceSchemaJSON generates a JSON schema for the Service struct and returns it as a JSON string.
func (sg *SchemaGenerator) GenerateServiceSchemaJSON() (string, error) {
	schema, err := sg.GenerateServiceSchema()
	if err != nil {
		return "", err
	}

	return sg.SchemaToJSON(schema)
}