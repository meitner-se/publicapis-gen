package schemagen

import (
	"bytes"
	"encoding/json"
	"fmt"

	yaml "github.com/goccy/go-yaml"
	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"

	"github.com/meitner-se/publicapis-gen/specification"
)

// Error messages
const (
	errorFailedToGenerate  = "failed to generate schema for"
	errorValidationFailed  = "validation failed"
	errorValidationErrors  = "validation errors"
	errorFailedToMarshal   = "failed to marshal schema to JSON"
	errorFailedToConvert   = "failed to convert schema to JSON"
	errorFailedToUnmarshal = "failed to unmarshal"
	errorDataNotValid      = "data is neither valid JSON nor YAML"
	errorConversionFailed  = "failed to convert YAML to JSON"
)

// schemaGenerator provides functionality to generate JSON schemas from specification structs.
type schemaGenerator struct {
	reflector *jsonschema.Reflector
}

// newSchemaGenerator creates a new schema generator with default configuration.
func newSchemaGenerator() *schemaGenerator {
	r := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            false,
		ExpandedStruct:            true,
	}

	return &schemaGenerator{
		reflector: r,
	}
}

// generateServiceSchema generates a JSON schema for the Service struct.
func (sg *schemaGenerator) generateServiceSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Service{})
	if schema == nil {
		return nil, fmt.Errorf("%s Service", errorFailedToGenerate)
	}

	return schema, nil
}

// generateEnumSchema generates a JSON schema for the Enum struct.
func (sg *schemaGenerator) generateEnumSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Enum{})
	if schema == nil {
		return nil, fmt.Errorf("%s Enum", errorFailedToGenerate)
	}

	return schema, nil
}

// generateObjectSchema generates a JSON schema for the Object struct.
func (sg *schemaGenerator) generateObjectSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Object{})
	if schema == nil {
		return nil, fmt.Errorf("%s Object", errorFailedToGenerate)
	}

	return schema, nil
}

// generateResourceSchema generates a JSON schema for the Resource struct.
func (sg *schemaGenerator) generateResourceSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Resource{})
	if schema == nil {
		return nil, fmt.Errorf("%s Resource", errorFailedToGenerate)
	}

	return schema, nil
}

// generateFieldSchema generates a JSON schema for the Field struct.
func (sg *schemaGenerator) generateFieldSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Field{})
	if schema == nil {
		return nil, fmt.Errorf("%s Field", errorFailedToGenerate)
	}

	return schema, nil
}

// generateResourceFieldSchema generates a JSON schema for the ResourceField struct.
func (sg *schemaGenerator) generateResourceFieldSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.ResourceField{})
	if schema == nil {
		return nil, fmt.Errorf("%s ResourceField", errorFailedToGenerate)
	}

	return schema, nil
}

// generateEndpointSchema generates a JSON schema for the Endpoint struct.
func (sg *schemaGenerator) generateEndpointSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.Endpoint{})
	if schema == nil {
		return nil, fmt.Errorf("%s Endpoint", errorFailedToGenerate)
	}

	return schema, nil
}

// generateEndpointRequestSchema generates a JSON schema for the EndpointRequest struct.
func (sg *schemaGenerator) generateEndpointRequestSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.EndpointRequest{})
	if schema == nil {
		return nil, fmt.Errorf("%s EndpointRequest", errorFailedToGenerate)
	}

	return schema, nil
}

// generateEndpointResponseSchema generates a JSON schema for the EndpointResponse struct.
func (sg *schemaGenerator) generateEndpointResponseSchema() (*jsonschema.Schema, error) {
	schema := sg.reflector.Reflect(&specification.EndpointResponse{})
	if schema == nil {
		return nil, fmt.Errorf("%s EndpointResponse", errorFailedToGenerate)
	}

	return schema, nil
}

// generateAllSchemas generates JSON schemas for all main specification structs.
func (sg *schemaGenerator) generateAllSchemas() (map[string]*jsonschema.Schema, error) {
	schemas := make(map[string]*jsonschema.Schema)

	serviceSchema, err := sg.generateServiceSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Service schema: %w", errorFailedToGenerate, err)
	}
	schemas["Service"] = serviceSchema

	enumSchema, err := sg.generateEnumSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Enum schema: %w", errorFailedToGenerate, err)
	}
	schemas["Enum"] = enumSchema

	objectSchema, err := sg.generateObjectSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Object schema: %w", errorFailedToGenerate, err)
	}
	schemas["Object"] = objectSchema

	resourceSchema, err := sg.generateResourceSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Resource schema: %w", errorFailedToGenerate, err)
	}
	schemas["Resource"] = resourceSchema

	fieldSchema, err := sg.generateFieldSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Field schema: %w", errorFailedToGenerate, err)
	}
	schemas["Field"] = fieldSchema

	resourceFieldSchema, err := sg.generateResourceFieldSchema()
	if err != nil {
		return nil, fmt.Errorf("%s ResourceField schema: %w", errorFailedToGenerate, err)
	}
	schemas["ResourceField"] = resourceFieldSchema

	endpointSchema, err := sg.generateEndpointSchema()
	if err != nil {
		return nil, fmt.Errorf("%s Endpoint schema: %w", errorFailedToGenerate, err)
	}
	schemas["Endpoint"] = endpointSchema

	endpointRequestSchema, err := sg.generateEndpointRequestSchema()
	if err != nil {
		return nil, fmt.Errorf("%s EndpointRequest schema: %w", errorFailedToGenerate, err)
	}
	schemas["EndpointRequest"] = endpointRequestSchema

	endpointResponseSchema, err := sg.generateEndpointResponseSchema()
	if err != nil {
		return nil, fmt.Errorf("%s EndpointResponse schema: %w", errorFailedToGenerate, err)
	}
	schemas["EndpointResponse"] = endpointResponseSchema

	return schemas, nil
}

// schemaToJSON converts a JSON schema to a JSON string.
func (sg *schemaGenerator) schemaToJSON(schema *jsonschema.Schema) (string, error) {
	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	return string(jsonBytes), nil
}

// GenerateServiceSchemaJSON generates a JSON schema for the Service struct and returns it as a JSON string.
func (sg *schemaGenerator) GenerateServiceSchemaJSON() (string, error) {
	schema, err := sg.generateServiceSchema()
	if err != nil {
		return "", err
	}

	return sg.schemaToJSON(schema)
}

// ValidateService validates a JSON/YAML representation of a Service against its schema.
func (sg *schemaGenerator) ValidateService(data []byte) error {
	schema, err := sg.generateServiceSchema()
	if err != nil {
		return fmt.Errorf("%s Service schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateEnum validates a JSON/YAML representation of an Enum against its schema.
func (sg *schemaGenerator) ValidateEnum(data []byte) error {
	schema, err := sg.generateEnumSchema()
	if err != nil {
		return fmt.Errorf("%s Enum schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateObject validates a JSON/YAML representation of an Object against its schema.
func (sg *schemaGenerator) ValidateObject(data []byte) error {
	schema, err := sg.generateObjectSchema()
	if err != nil {
		return fmt.Errorf("%s Object schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateResource validates a JSON/YAML representation of a Resource against its schema.
func (sg *schemaGenerator) ValidateResource(data []byte) error {
	schema, err := sg.generateResourceSchema()
	if err != nil {
		return fmt.Errorf("%s Resource schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateField validates a JSON/YAML representation of a Field against its schema.
func (sg *schemaGenerator) ValidateField(data []byte) error {
	schema, err := sg.generateFieldSchema()
	if err != nil {
		return fmt.Errorf("%s Field schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateResourceField validates a JSON/YAML representation of a ResourceField against its schema.
func (sg *schemaGenerator) ValidateResourceField(data []byte) error {
	schema, err := sg.generateResourceFieldSchema()
	if err != nil {
		return fmt.Errorf("%s ResourceField schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateEndpoint validates a JSON/YAML representation of an Endpoint against its schema.
func (sg *schemaGenerator) ValidateEndpoint(data []byte) error {
	schema, err := sg.generateEndpointSchema()
	if err != nil {
		return fmt.Errorf("%s Endpoint schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateEndpointRequest validates a JSON/YAML representation of an EndpointRequest against its schema.
func (sg *schemaGenerator) ValidateEndpointRequest(data []byte) error {
	schema, err := sg.generateEndpointRequestSchema()
	if err != nil {
		return fmt.Errorf("%s EndpointRequest schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// ValidateEndpointResponse validates a JSON/YAML representation of an EndpointResponse against its schema.
func (sg *schemaGenerator) ValidateEndpointResponse(data []byte) error {
	schema, err := sg.generateEndpointResponseSchema()
	if err != nil {
		return fmt.Errorf("%s EndpointResponse schema: %w", errorFailedToGenerate, err)
	}

	return sg.validateWithSchema(schema, data)
}

// validateWithSchema is a helper function that validates data against a JSON schema.
func (sg *schemaGenerator) validateWithSchema(schema *jsonschema.Schema, data []byte) error {
	// Convert schema to JSON string
	schemaJSON, err := sg.schemaToJSON(schema)
	if err != nil {
		return fmt.Errorf("%s: %w", errorFailedToConvert, err)
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)

	// Convert data to JSON if it might be YAML
	jsonData, err := sg.convertToJSON(data)
	if err != nil {
		return fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	// Create document loader
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	if !result.Valid() {
		return fmt.Errorf("%s: %v", errorValidationErrors, result.Errors())
	}

	return nil
}

// convertToJSON converts YAML or JSON data to JSON format.
func (sg *schemaGenerator) convertToJSON(data []byte) ([]byte, error) {
	// First, try to parse as JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err == nil {
		// It's valid JSON, return as-is
		return data, nil
	}

	// Try to parse as YAML
	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, fmt.Errorf("%s: %w", errorDataNotValid, err)
	}

	// Convert YAML data to JSON
	jsonBytes, err := json.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	return jsonBytes, nil
}

// ParseServiceFromJSON parses and validates a Service from JSON data.
func (sg *schemaGenerator) ParseServiceFromJSON(data []byte) (*specification.Service, error) {
	// Validate against schema first
	if err := sg.ValidateService(data); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	// Parse the JSON
	var service specification.Service
	if err := json.Unmarshal(data, &service); err != nil {
		return nil, fmt.Errorf("%s JSON: %w", errorFailedToUnmarshal, err)
	}

	return &service, nil
}

// ParseServiceFromYAML parses and validates a Service from YAML data.
func (sg *schemaGenerator) ParseServiceFromYAML(data []byte) (*specification.Service, error) {
	// Convert YAML to JSON for validation
	jsonData, err := sg.convertToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	// Validate against schema
	if err := sg.ValidateService(jsonData); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	// Parse the YAML directly
	var service specification.Service
	if err := yaml.Unmarshal(data, &service); err != nil {
		return nil, fmt.Errorf("%s YAML: %w", errorFailedToUnmarshal, err)
	}

	return &service, nil
}

// ParseEnumFromJSON parses and validates an Enum from JSON data.
func (sg *schemaGenerator) ParseEnumFromJSON(data []byte) (*specification.Enum, error) {
	if err := sg.ValidateEnum(data); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var enum specification.Enum
	if err := json.Unmarshal(data, &enum); err != nil {
		return nil, fmt.Errorf("%s JSON: %w", errorFailedToUnmarshal, err)
	}

	return &enum, nil
}

// ParseEnumFromYAML parses and validates an Enum from YAML data.
func (sg *schemaGenerator) ParseEnumFromYAML(data []byte) (*specification.Enum, error) {
	jsonData, err := sg.convertToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	if err := sg.ValidateEnum(jsonData); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var enum specification.Enum
	if err := yaml.Unmarshal(data, &enum); err != nil {
		return nil, fmt.Errorf("%s YAML: %w", errorFailedToUnmarshal, err)
	}

	return &enum, nil
}

// ParseObjectFromJSON parses and validates an Object from JSON data.
func (sg *schemaGenerator) ParseObjectFromJSON(data []byte) (*specification.Object, error) {
	if err := sg.ValidateObject(data); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var object specification.Object
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("%s JSON: %w", errorFailedToUnmarshal, err)
	}

	return &object, nil
}

// ParseObjectFromYAML parses and validates an Object from YAML data.
func (sg *schemaGenerator) ParseObjectFromYAML(data []byte) (*specification.Object, error) {
	jsonData, err := sg.convertToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	if err := sg.ValidateObject(jsonData); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var object specification.Object
	if err := yaml.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("%s YAML: %w", errorFailedToUnmarshal, err)
	}

	return &object, nil
}

// ParseResourceFromJSON parses and validates a Resource from JSON data.
func (sg *schemaGenerator) ParseResourceFromJSON(data []byte) (*specification.Resource, error) {
	if err := sg.ValidateResource(data); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var resource specification.Resource
	if err := json.Unmarshal(data, &resource); err != nil {
		return nil, fmt.Errorf("%s JSON: %w", errorFailedToUnmarshal, err)
	}

	return &resource, nil
}

// ParseResourceFromYAML parses and validates a Resource from YAML data.
func (sg *schemaGenerator) ParseResourceFromYAML(data []byte) (*specification.Resource, error) {
	jsonData, err := sg.convertToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorConversionFailed, err)
	}

	if err := sg.ValidateResource(jsonData); err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidationFailed, err)
	}

	var resource specification.Resource
	if err := yaml.Unmarshal(data, &resource); err != nil {
		return nil, fmt.Errorf("%s YAML: %w", errorFailedToUnmarshal, err)
	}

	return &resource, nil
}

// GenerateSchemas generates JSON schemas for all specification types
// and writes them to the provided buffer as a JSON object.
func GenerateSchemas(buf *bytes.Buffer, service *specification.Service) error {
	// Create schema generator
	generator := newSchemaGenerator()

	// Generate all schemas
	schemas, err := generator.generateAllSchemas()
	if err != nil {
		return fmt.Errorf("failed to generate schemas: %w", err)
	}

	// Convert all schemas to a combined JSON structure
	schemaMap := make(map[string]interface{})
	for name, schemaObj := range schemas {
		// Convert each schema to JSON and then parse it back to interface{} for clean structure
		jsonStr, err := generator.schemaToJSON(schemaObj)
		if err != nil {
			return fmt.Errorf("failed to convert %s schema to JSON: %w", name, err)
		}

		var schemaData interface{}
		if err := json.Unmarshal([]byte(jsonStr), &schemaData); err != nil {
			return fmt.Errorf("failed to parse %s schema JSON: %w", name, err)
		}

		schemaMap[name] = schemaData
	}

	// Marshal the combined schema map to JSON with proper indentation
	outputData, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal combined schemas to JSON: %w", err)
	}

	// Write to buffer
	buf.Write(outputData)

	return nil
}
