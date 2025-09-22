package schemagen

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"

	"github.com/meitner-se/publicapis-gen/specification"
)

// Error messages
const (
	errorFailedToGenerate = "failed to generate schema for"
	errorFailedToMarshal  = "failed to marshal schema to JSON"
	errorFailedToConvert  = "failed to convert schema to JSON"
)

// GenerateSchemas generates JSON schemas for all specification types and writes them to the buffer.
// The output is a JSON object with schema names as keys and their JSON schema definitions as values.
func GenerateSchemas(buf *bytes.Buffer) error {
	reflector := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            false,
		ExpandedStruct:            true,
	}

	// Generate all schemas
	schemas := make(map[string]interface{})

	// Service schema
	if schema := reflector.Reflect(&specification.Service{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Service: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Service: %w", errorFailedToMarshal, err)
		}
		schemas["Service"] = schemaData
	} else {
		return fmt.Errorf("%s Service", errorFailedToGenerate)
	}

	// Enum schema
	if schema := reflector.Reflect(&specification.Enum{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Enum: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Enum: %w", errorFailedToMarshal, err)
		}
		schemas["Enum"] = schemaData
	} else {
		return fmt.Errorf("%s Enum", errorFailedToGenerate)
	}

	// Object schema
	if schema := reflector.Reflect(&specification.Object{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Object: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Object: %w", errorFailedToMarshal, err)
		}
		schemas["Object"] = schemaData
	} else {
		return fmt.Errorf("%s Object", errorFailedToGenerate)
	}

	// Resource schema
	if schema := reflector.Reflect(&specification.Resource{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Resource: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Resource: %w", errorFailedToMarshal, err)
		}
		schemas["Resource"] = schemaData
	} else {
		return fmt.Errorf("%s Resource", errorFailedToGenerate)
	}

	// Field schema
	if schema := reflector.Reflect(&specification.Field{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Field: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Field: %w", errorFailedToMarshal, err)
		}
		schemas["Field"] = schemaData
	} else {
		return fmt.Errorf("%s Field", errorFailedToGenerate)
	}

	// ResourceField schema
	if schema := reflector.Reflect(&specification.ResourceField{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s ResourceField: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s ResourceField: %w", errorFailedToMarshal, err)
		}
		schemas["ResourceField"] = schemaData
	} else {
		return fmt.Errorf("%s ResourceField", errorFailedToGenerate)
	}

	// Endpoint schema
	if schema := reflector.Reflect(&specification.Endpoint{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s Endpoint: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s Endpoint: %w", errorFailedToMarshal, err)
		}
		schemas["Endpoint"] = schemaData
	} else {
		return fmt.Errorf("%s Endpoint", errorFailedToGenerate)
	}

	// EndpointRequest schema
	if schema := reflector.Reflect(&specification.EndpointRequest{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s EndpointRequest: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s EndpointRequest: %w", errorFailedToMarshal, err)
		}
		schemas["EndpointRequest"] = schemaData
	} else {
		return fmt.Errorf("%s EndpointRequest", errorFailedToGenerate)
	}

	// EndpointResponse schema
	if schema := reflector.Reflect(&specification.EndpointResponse{}); schema != nil {
		schemaJSON, err := schemaToJSON(schema)
		if err != nil {
			return fmt.Errorf("%s EndpointResponse: %w", errorFailedToConvert, err)
		}
		var schemaData interface{}
		if err := json.Unmarshal([]byte(schemaJSON), &schemaData); err != nil {
			return fmt.Errorf("%s EndpointResponse: %w", errorFailedToMarshal, err)
		}
		schemas["EndpointResponse"] = schemaData
	} else {
		return fmt.Errorf("%s EndpointResponse", errorFailedToGenerate)
	}

	// Marshal the combined schema map to JSON with proper indentation
	outputData, err := json.MarshalIndent(schemas, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal combined schemas to JSON: %w", err)
	}

	// Write to buffer
	buf.Write(outputData)

	return nil
}

// schemaToJSON converts a JSON schema to a JSON string.
func schemaToJSON(schema *jsonschema.Schema) (string, error) {
	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%s: %w", errorFailedToMarshal, err)
	}

	return string(jsonBytes), nil
}
