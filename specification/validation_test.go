package specification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// ValidateService Tests
// ============================================================================

func TestValidateService(t *testing.T) {
	// Test valid service
	validService := &Service{
		Name:    "TestService",
		Version: "1.0.0",
		Enums: []Enum{
			{
				Name:        "Status",
				Description: "Status enum",
				Values: []EnumValue{
					{Name: "Active", Description: "Active status"},
					{Name: "Inactive", Description: "Inactive status"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "User",
				Description: "User object",
				Fields: []Field{
					{
						Name:        "name",
						Description: "User name",
						Type:        FieldTypeString,
					},
				},
			},
		},
		Resources: []Resource{
			{
				Name:        "Users",
				Description: "User resource",
				Operations:  []string{OperationCreate, OperationRead},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "name",
							Description: "User name",
							Type:        FieldTypeString,
						},
						Operations: []string{OperationCreate, OperationRead},
					},
				},
			},
		},
	}

	err := ValidateService(validService)
	assert.NoError(t, err, "Valid service should pass validation")

	t.Run("invalid resource operation", func(t *testing.T) {
		invalidService := *validService
		invalidService.Resources[0].Operations = []string{"invalid"}

		err := ValidateService(&invalidService)
		assert.Error(t, err, "Service with invalid operation should fail validation")
		assert.Contains(t, err.Error(), "invalid operation")
	})

	t.Run("invalid field type", func(t *testing.T) {
		serviceWithInvalidFieldType := &Service{
			Name:    "TestService",
			Version: "1.0.0",
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User resource",
					Operations:  []string{OperationCreate, OperationRead}, // Valid operations
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "name",
								Description: "User name",
								Type:        "InvalidType", // Invalid type
							},
							Operations: []string{OperationCreate, OperationRead},
						},
					},
				},
			},
		}

		err := ValidateService(serviceWithInvalidFieldType)
		assert.Error(t, err, "Service with invalid field type should fail validation")
		assert.Contains(t, err.Error(), "field type")
	})
}

// ============================================================================
// ValidateOperations Tests
// ============================================================================

func TestValidateOperations(t *testing.T) {
	// Test valid operations
	validOperations := []string{OperationCreate, OperationRead, OperationUpdate, OperationDelete}
	err := ValidateOperations(validOperations)
	assert.NoError(t, err, "Valid operations should pass validation")

	// Test empty operations
	err = ValidateOperations([]string{})
	assert.NoError(t, err, "Empty operations should pass validation")

	t.Run("invalid operation", func(t *testing.T) {
		invalidOperations := []string{OperationCreate, "invalid"}
		err := ValidateOperations(invalidOperations)
		assert.Error(t, err, "Invalid operation should fail validation")
		assert.Contains(t, err.Error(), "invalid operation")
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("lowercase operation", func(t *testing.T) {
		invalidOperations := []string{"create", OperationRead}
		err := ValidateOperations(invalidOperations)
		assert.Error(t, err, "Lowercase operation should fail validation")
		assert.Contains(t, err.Error(), "invalid operation")
		assert.Contains(t, err.Error(), "create")
	})

	t.Run("multiple invalid operations", func(t *testing.T) {
		invalidOperations := []string{"create", "read", "invalid"}
		err := ValidateOperations(invalidOperations)
		assert.Error(t, err, "Multiple invalid operations should fail validation")
		assert.Contains(t, err.Error(), "invalid operation")
	})
}

// ============================================================================
// ValidateFieldType Tests
// ============================================================================

func TestValidateFieldType(t *testing.T) {
	service := &Service{
		Enums: []Enum{
			{Name: "Status", Description: "Status enum"},
		},
		Objects: []Object{
			{Name: "User", Description: "User object"},
		},
	}

	// Test valid primitive types
	primitiveTypes := []string{
		FieldTypeUUID, FieldTypeDate, FieldTypeTimestamp,
		FieldTypeString, FieldTypeInt, FieldTypeBool,
	}
	for _, primitiveType := range primitiveTypes {
		err := ValidateFieldType(service, primitiveType)
		assert.NoError(t, err, "Valid primitive type %s should pass validation", primitiveType)
	}

	// Test valid enum type
	err := ValidateFieldType(service, "Status")
	assert.NoError(t, err, "Valid enum type should pass validation")

	// Test valid object type
	err = ValidateFieldType(service, "User")
	assert.NoError(t, err, "Valid object type should pass validation")

	t.Run("invalid type", func(t *testing.T) {
		err := ValidateFieldType(service, "InvalidType")
		assert.Error(t, err, "Invalid type should fail validation")
		assert.Contains(t, err.Error(), "invalid field type")
		assert.Contains(t, err.Error(), "InvalidType")
	})

	t.Run("lowercase primitive type", func(t *testing.T) {
		err := ValidateFieldType(service, "string")
		assert.Error(t, err, "Lowercase primitive type should fail validation")
		assert.Contains(t, err.Error(), "invalid field type")
		assert.Contains(t, err.Error(), "string")
	})
}

// ============================================================================
// ValidateModifiers Tests
// ============================================================================

func TestValidateModifiers(t *testing.T) {
	// Test valid modifiers
	validModifiers := []string{ModifierNullable, ModifierArray}
	err := ValidateModifiers(validModifiers)
	assert.NoError(t, err, "Valid modifiers should pass validation")

	// Test empty modifiers
	err = ValidateModifiers([]string{})
	assert.NoError(t, err, "Empty modifiers should pass validation")

	// Test single valid modifier
	err = ValidateModifiers([]string{ModifierNullable})
	assert.NoError(t, err, "Single valid modifier should pass validation")

	err = ValidateModifiers([]string{ModifierArray})
	assert.NoError(t, err, "Single valid modifier should pass validation")

	t.Run("invalid modifier", func(t *testing.T) {
		invalidModifiers := []string{ModifierNullable, "invalid"}
		err := ValidateModifiers(invalidModifiers)
		assert.Error(t, err, "Invalid modifier should fail validation")
		assert.Contains(t, err.Error(), "invalid modifier")
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("lowercase modifier", func(t *testing.T) {
		invalidModifiers := []string{"nullable", ModifierArray}
		err := ValidateModifiers(invalidModifiers)
		assert.Error(t, err, "Lowercase modifier should fail validation")
		assert.Contains(t, err.Error(), "invalid modifier")
		assert.Contains(t, err.Error(), "nullable")
	})

	t.Run("mixed case modifiers", func(t *testing.T) {
		invalidModifiers := []string{"nullable", "array"}
		err := ValidateModifiers(invalidModifiers)
		assert.Error(t, err, "Mixed case modifiers should fail validation")
		assert.Contains(t, err.Error(), "invalid modifier")
	})
}

// ============================================================================
// ValidateField Tests
// ============================================================================

func TestValidateField(t *testing.T) {
	service := &Service{
		Enums: []Enum{
			{Name: "Status", Description: "Status enum"},
		},
		Objects: []Object{
			{Name: "User", Description: "User object"},
		},
	}

	// Test valid field
	validField := Field{
		Name:        "name",
		Description: "User name",
		Type:        FieldTypeString,
		Modifiers:   []string{ModifierNullable},
	}

	err := ValidateField(service, &validField)
	assert.NoError(t, err, "Valid field should pass validation")

	t.Run("field with invalid type", func(t *testing.T) {
		invalidField := validField
		invalidField.Type = "InvalidType"

		err := ValidateField(service, &invalidField)
		assert.Error(t, err, "Field with invalid type should fail validation")
		assert.Contains(t, err.Error(), "field type")
	})

	t.Run("field with invalid modifier", func(t *testing.T) {
		invalidField := validField
		invalidField.Modifiers = []string{"invalid"}

		err := ValidateField(service, &invalidField)
		assert.Error(t, err, "Field with invalid modifier should fail validation")
		assert.Contains(t, err.Error(), "field modifiers")
	})

	t.Run("field with both invalid type and modifier", func(t *testing.T) {
		invalidField := validField
		invalidField.Type = "InvalidType"
		invalidField.Modifiers = []string{"invalid"}

		err := ValidateField(service, &invalidField)
		assert.Error(t, err, "Field with both invalid type and modifier should fail validation")
	})
}

// ============================================================================
// ValidateResourceField Tests
// ============================================================================

func TestValidateResourceField(t *testing.T) {
	service := &Service{
		Enums:   []Enum{},
		Objects: []Object{},
	}

	// Test valid resource field
	validResourceField := ResourceField{
		Field: Field{
			Name:        "name",
			Description: "User name",
			Type:        FieldTypeString,
			Modifiers:   []string{ModifierNullable},
		},
		Operations: []string{OperationCreate, OperationRead},
	}

	err := ValidateResourceField(service, &validResourceField)
	assert.NoError(t, err, "Valid resource field should pass validation")

	t.Run("resource field with invalid operation", func(t *testing.T) {
		invalidResourceField := validResourceField
		invalidResourceField.Operations = []string{"invalid"}

		err := ValidateResourceField(service, &invalidResourceField)
		assert.Error(t, err, "Resource field with invalid operation should fail validation")
		assert.Contains(t, err.Error(), "field operations")
	})

	t.Run("resource field with invalid field type", func(t *testing.T) {
		invalidResourceField := validResourceField
		invalidResourceField.Type = "InvalidType"

		err := ValidateResourceField(service, &invalidResourceField)
		assert.Error(t, err, "Resource field with invalid field type should fail validation")
		assert.Contains(t, err.Error(), "field type")
	})
}

// ============================================================================
// ValidateResource Tests
// ============================================================================

func TestValidateResource(t *testing.T) {
	service := &Service{
		Enums:   []Enum{},
		Objects: []Object{},
	}

	// Test valid resource
	validResource := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationRead},
		Fields: []ResourceField{
			{
				Field: Field{
					Name:        "name",
					Description: "User name",
					Type:        FieldTypeString,
				},
				Operations: []string{OperationCreate, OperationRead},
			},
		},
		Endpoints: []Endpoint{},
	}

	err := ValidateResource(service, &validResource)
	assert.NoError(t, err, "Valid resource should pass validation")

	t.Run("resource with invalid operation", func(t *testing.T) {
		invalidResource := validResource
		invalidResource.Operations = []string{"invalid"}

		err := ValidateResource(service, &invalidResource)
		assert.Error(t, err, "Resource with invalid operation should fail validation")
		assert.Contains(t, err.Error(), "resource operations")
	})

	t.Run("resource with invalid field", func(t *testing.T) {
		invalidResource := validResource
		invalidResource.Fields[0].Type = "InvalidType"

		err := ValidateResource(service, &invalidResource)
		assert.Error(t, err, "Resource with invalid field should fail validation")
		assert.Contains(t, err.Error(), "field 0 (name)")
	})
}

// ============================================================================
// ValidateObject Tests
// ============================================================================

func TestValidateObject(t *testing.T) {
	service := &Service{
		Enums:   []Enum{},
		Objects: []Object{},
	}

	// Test valid object
	validObject := Object{
		Name:        "User",
		Description: "User object",
		Fields: []Field{
			{
				Name:        "name",
				Description: "User name",
				Type:        FieldTypeString,
			},
		},
	}

	err := ValidateObject(service, &validObject)
	assert.NoError(t, err, "Valid object should pass validation")

	t.Run("object with invalid field", func(t *testing.T) {
		invalidObject := validObject
		invalidObject.Fields[0].Type = "InvalidType"

		err := ValidateObject(service, &invalidObject)
		assert.Error(t, err, "Object with invalid field should fail validation")
		assert.Contains(t, err.Error(), "field 0 (name)")
	})
}

// ============================================================================
// ValidateEndpoint Tests
// ============================================================================

func TestValidateEndpoint(t *testing.T) {
	service := &Service{
		Enums:   []Enum{},
		Objects: []Object{},
	}

	// Test valid endpoint
	validEndpoint := Endpoint{
		Name:        "GetUser",
		Title:       "Get User",
		Description: "Get user by ID",
		Method:      "GET",
		Path:        "/{id}",
		Request: EndpointRequest{
			QueryParams: []Field{
				{
					Name:        "includeDetails",
					Description: "Include user details",
					Type:        FieldTypeBool,
				},
			},
		},
		Response: EndpointResponse{
			BodyFields: []Field{
				{
					Name:        "name",
					Description: "User name",
					Type:        FieldTypeString,
				},
			},
		},
	}

	err := ValidateEndpoint(service, &validEndpoint)
	assert.NoError(t, err, "Valid endpoint should pass validation")

	t.Run("endpoint with invalid request field", func(t *testing.T) {
		invalidEndpoint := validEndpoint
		invalidEndpoint.Request.QueryParams[0].Type = "InvalidType"

		err := ValidateEndpoint(service, &invalidEndpoint)
		assert.Error(t, err, "Endpoint with invalid request field should fail validation")
		assert.Contains(t, err.Error(), "request query param 0")
	})

	t.Run("endpoint with invalid response field", func(t *testing.T) {
		// Create endpoint with only response fields (no request fields to avoid conflicts)
		endpointWithInvalidResponse := Endpoint{
			Name:        "GetUser",
			Title:       "Get User",
			Description: "Get user by ID",
			Method:      "GET",
			Path:        "/{id}",
			Request:     EndpointRequest{}, // Empty request
			Response: EndpointResponse{
				BodyFields: []Field{
					{
						Name:        "name",
						Description: "User name",
						Type:        "InvalidType", // Invalid type
					},
				},
			},
		}

		err := ValidateEndpoint(service, &endpointWithInvalidResponse)
		assert.Error(t, err, "Endpoint with invalid response field should fail validation")
		assert.Contains(t, err.Error(), "response body field 0")
	})
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestValidationIntegration(t *testing.T) {
	t.Run("parsing applies validation", func(t *testing.T) {
		// Test that parsing a service with validation errors fails
		invalidYAML := `
name: "TestService"
version: "1.0.0"
enums: []
objects: []
resources:
  - name: "Users"
    description: "User resource"
    operations: ["invalid"]  # Invalid operation
    fields: []
    endpoints: []
`
		_, err := ParseServiceFromYAML([]byte(invalidYAML))
		assert.Error(t, err, "Parsing invalid YAML should fail due to validation")
		assert.Contains(t, err.Error(), "validation failed")
		assert.Contains(t, err.Error(), "invalid operation")
	})

	t.Run("valid YAML passes validation", func(t *testing.T) {
		validYAML := `
name: "TestService"
version: "1.0.0"
enums: []
objects: []
resources:
  - name: "Users"
    description: "User resource"
    operations: ["Create", "Read"]
    fields:
      - name: "name"
        description: "User name"
        type: "String"
        operations: ["Create", "Read"]
    endpoints: []
`
		service, err := ParseServiceFromYAML([]byte(validYAML))
		assert.NoError(t, err, "Parsing valid YAML should succeed")
		assert.Equal(t, "TestService", service.Name)
	})

	t.Run("validation catches PascalCase requirement", func(t *testing.T) {
		// Test with lowercase modifiers (should fail)
		invalidYAML := `
name: "TestService"
version: "1.0.0"
enums: []
objects: []
resources:
  - name: "Users"
    description: "User resource"
    operations: ["Create", "Read"]
    fields:
      - name: "name"
        description: "User name"
        type: "String"
        modifiers: ["nullable"]  # Should be "Nullable"
        operations: ["Create", "Read"]
    endpoints: []
`
		_, err := ParseServiceFromYAML([]byte(invalidYAML))
		assert.Error(t, err, "Parsing YAML with lowercase modifiers should fail")
		assert.Contains(t, err.Error(), "validation failed")
		assert.Contains(t, err.Error(), "invalid modifier")
	})
}