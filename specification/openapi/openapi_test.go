package openapi

import (
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// NewGenerator Function Tests
// ============================================================================

// TestNewGenerator tests the creation of a new OpenAPI generator.
func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()

	assert.NotNil(t, generator, "Generator should not be nil")
	assert.Equal(t, "3.1.0", generator.Version, "Generator version should be 3.1.0")
	assert.Equal(t, "", generator.Title, "Generator title should be empty by default")
	assert.Equal(t, "", generator.Description, "Generator description should be empty by default")
	assert.Equal(t, "", generator.ServerURL, "Generator server URL should be empty by default")
}

// ============================================================================
// Generator Tests
// ============================================================================

// TestGenerator_GenerateFromService tests OpenAPI document generation from services.
func TestGenerator_GenerateFromService(t *testing.T) {
	// Test with nil service
	t.Run("nil service returns error", func(t *testing.T) {
		generator := NewGenerator()

		document, err := generator.GenerateFromService(nil)

		assert.Nil(t, document, "Document should be nil when service is nil")
		assert.EqualError(t, err, "invalid service: service cannot be nil", "Should return invalid service error")
	})

	// Test with valid service
	t.Run("valid service generates document", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: "TestService",
		}

		document, err := generator.GenerateFromService(service)

		assert.NotNil(t, document, "Document should not be nil with valid service")
		assert.NoError(t, err, "Should not return error with valid service")
		assert.Equal(t, "3.1.0", document.Version, "Document version should be 3.1.0")
		assert.Equal(t, "TestService", document.Info.Title, "Document title should match service name")
	})

	// Test with complex service
	t.Run("complex service with enums and objects", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: "UserAPI",
			Enums: []specification.Enum{
				{
					Name:        "Status",
					Description: "User status enumeration",
					Values: []specification.EnumValue{
						{Name: "Active", Description: "User is active"},
						{Name: "Inactive", Description: "User is inactive"},
					},
				},
			},
			Objects: []specification.Object{
				{
					Name:        "User",
					Description: "User object",
					Fields: []specification.Field{
						{
							Name:        "id",
							Description: "User identifier",
							Type:        specification.FieldTypeUUID,
						},
						{
							Name:        "email",
							Description: "User email address",
							Type:        specification.FieldTypeString,
						},
						{
							Name:        "status",
							Description: "User status",
							Type:        "Status",
						},
					},
				},
			},
			Resources: []specification.Resource{
				{
					Name:        "User",
					Description: "User resource",
					Operations:  []string{specification.OperationCreate, specification.OperationRead},
					Fields: []specification.ResourceField{
						{
							Field: specification.Field{
								Name:        "id",
								Description: "User identifier",
								Type:        specification.FieldTypeUUID,
							},
							Operations: []string{specification.OperationRead},
						},
						{
							Field: specification.Field{
								Name:        "email",
								Description: "User email address",
								Type:        specification.FieldTypeString,
							},
							Operations: []string{specification.OperationCreate, specification.OperationRead},
						},
					},
					Endpoints: []specification.Endpoint{
						{
							Name:        "Create",
							Title:       "Create User",
							Description: "Create a new user",
							Method:      "POST",
							Path:        "",
							Request: specification.EndpointRequest{
								ContentType: "application/json",
								BodyParams: []specification.Field{
									{
										Name:        "email",
										Description: "User email address",
										Type:        specification.FieldTypeString,
									},
								},
							},
							Response: specification.EndpointResponse{
								ContentType: "application/json",
								StatusCode:  201,
								BodyObject:  stringPtr("User"),
							},
						},
					},
				},
			},
		}

		document, err := generator.GenerateFromService(service)

		assert.NoError(t, err, "Should generate document successfully")
		assert.NotNil(t, document, "Document should not be nil")
		assert.Equal(t, "UserAPI", document.Info.Title, "Document title should match service name")

		// Test JSON output contains expected elements
		jsonBytes, err := generator.ToJSON(document)
		assert.NoError(t, err, "Should convert document to JSON successfully")
		jsonString := string(jsonBytes)
		assert.Contains(t, jsonString, "Status", "JSON should contain Status enum")
		assert.Contains(t, jsonString, "User", "JSON should contain User object")
		assert.Contains(t, jsonString, "Active", "JSON should contain enum values")
		assert.Contains(t, jsonString, "/user", "JSON should contain user path")
	})

	// Test with service version and servers
	t.Run("service with version and servers", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name:    "UserAPI",
			Version: "2.0.0",
			Servers: []specification.ServiceServer{
				{
					URL:         "https://api.example.com",
					Description: "Production server",
				},
				{
					URL:         "https://staging-api.example.com",
					Description: "Staging server",
				},
			},
			Objects: []specification.Object{
				{
					Name:        "User",
					Description: "User object",
					Fields: []specification.Field{
						{
							Name:        "email",
							Description: "User email",
							Type:        specification.FieldTypeString,
							Modifiers:   []string{specification.ModifierNullable},
						},
					},
				},
			},
		}

		document, err := generator.GenerateFromService(service)

		assert.NoError(t, err, "Should generate document successfully")
		assert.NotNil(t, document, "Document should not be nil")
		assert.Equal(t, "2.0.0", document.Info.Version, "Document version should come from service")
		assert.Equal(t, 2, len(document.Servers), "Document should have 2 servers from service")
		assert.Equal(t, "https://api.example.com", document.Servers[0].URL, "First server URL should match service")
		assert.Equal(t, "Production server", document.Servers[0].Description, "First server description should match service")
		assert.Equal(t, "https://staging-api.example.com", document.Servers[1].URL, "Second server URL should match service")
		assert.Equal(t, "Staging server", document.Servers[1].Description, "Second server description should match service")

		// Test JSON output
		jsonBytes, err := generator.ToJSON(document)
		assert.NoError(t, err, "Should convert document to JSON successfully")
		jsonString := string(jsonBytes)
		assert.Contains(t, jsonString, "2.0.0", "JSON should contain service version")
		assert.Contains(t, jsonString, "https://api.example.com", "JSON should contain first server URL")
		assert.Contains(t, jsonString, "https://staging-api.example.com", "JSON should contain second server URL")
	})
}

// TestGenerator_ToYAML tests YAML conversion functionality.
func TestGenerator_ToYAML(t *testing.T) {
	// Test with nil document
	t.Run("nil document returns error", func(t *testing.T) {
		generator := NewGenerator()

		yamlBytes, err := generator.ToYAML(nil)

		assert.Nil(t, yamlBytes, "YAML bytes should be nil when document is nil")
		assert.EqualError(t, err, "invalid document: document cannot be nil", "Should return invalid document error")
	})

	// Test with valid document
	t.Run("valid document converts successfully", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: "TestService",
		}

		document, err := generator.GenerateFromService(service)
		assert.NoError(t, err, "Should generate document successfully")

		yamlBytes, err := generator.ToYAML(document)

		assert.NoError(t, err, "Should convert document to YAML successfully")
		assert.NotNil(t, yamlBytes, "YAML bytes should not be nil")
		assert.Contains(t, string(yamlBytes), "TestService", "YAML should contain service name")
		assert.Contains(t, string(yamlBytes), "3.1.0", "YAML should contain OpenAPI version")
	})
}

// TestGenerator_ToJSON tests JSON conversion functionality.
func TestGenerator_ToJSON(t *testing.T) {
	// Test with nil document
	t.Run("nil document returns error", func(t *testing.T) {
		generator := NewGenerator()

		jsonBytes, err := generator.ToJSON(nil)

		assert.Nil(t, jsonBytes, "JSON bytes should be nil when document is nil")
		assert.EqualError(t, err, "invalid document: document cannot be nil", "Should return invalid document error")
	})

	// Test with valid document
	t.Run("valid document converts successfully", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: "TestService",
		}

		document, err := generator.GenerateFromService(service)
		assert.NoError(t, err, "Should generate document successfully")

		jsonBytes, err := generator.ToJSON(document)

		assert.NoError(t, err, "Should convert document to JSON successfully")
		assert.NotNil(t, jsonBytes, "JSON bytes should not be nil")
		assert.Contains(t, string(jsonBytes), "TestService", "JSON should contain service name")
		assert.Contains(t, string(jsonBytes), "3.1.0", "JSON should contain OpenAPI version")
	})
}

// TestGeneratorConfiguration tests generator configuration options.
func TestGeneratorConfiguration(t *testing.T) {
	generator := &Generator{
		Version:     "3.1.0",
		Title:       "Custom API",
		Description: "Custom API Description",
		ServerURL:   "https://custom.example.com",
	}

	assert.Equal(t, "Custom API", generator.Title, "Generator title should match configured value")
	assert.Equal(t, "Custom API Description", generator.Description, "Generator description should match configured value")
	assert.Equal(t, "https://custom.example.com", generator.ServerURL, "Generator server URL should match configured value")
}

// ============================================================================
// Error Response Tests
// ============================================================================

// TestGenerator_addErrorResponses tests error response generation functionality.
func TestGenerator_addErrorResponses(t *testing.T) {
	// Test with ErrorCode enum
	t.Run("with ErrorCode enum generates all responses", func(t *testing.T) {
		generator := NewGenerator()

		// Create service with ErrorCode enum (simulating ApplyOverlay result)
		service := &specification.Service{
			Name: "TestService",
			Enums: []specification.Enum{
				{
					Name:        "ErrorCode",
					Description: "Standard error codes used in API responses",
					Values: []specification.EnumValue{
						{Name: "BadRequest", Description: "The request was malformed or contained invalid parameters. 400 status code"},
						{Name: "Unauthorized", Description: "The request is missing valid authentication credentials. 401 status code"},
						{Name: "Forbidden", Description: "Request is authenticated, but the user is not allowed to perform the operation. 403 status code"},
						{Name: "NotFound", Description: "The requested resource or endpoint does not exist. 404 status code"},
						{Name: "Conflict", Description: "The request could not be completed due to a conflict. 409 status code"},
						{Name: "UnprocessableEntity", Description: "The request was well-formed but failed validation. 422 status code"},
						{Name: "RateLimited", Description: "When the rate limit has been exceeded. 429 status code"},
						{Name: "Internal", Description: "Some serverside issue. 500 status code"},
					},
				},
			},
			Objects: []specification.Object{
				{
					Name:        "Error",
					Description: "Standard error response object containing error code and message",
					Fields: []specification.Field{
						{Name: "Code", Description: "The specific error code", Type: "ErrorCode"},
						{Name: "Message", Description: "Human-readable error message", Type: specification.FieldTypeString},
					},
				},
			},
		}

		// Test endpoint with body parameters (should include 422)
		endpointWithBody := specification.Endpoint{
			Name:   "CreateUser",
			Method: "POST",
			Path:   "/users",
			Request: specification.EndpointRequest{
				BodyParams: []specification.Field{
					{Name: "email", Type: specification.FieldTypeString},
				},
			},
			Response: specification.EndpointResponse{StatusCode: 201},
		}

		responses := orderedmap.New[string, *v3.Response]()
		generator.addErrorResponses(responses, endpointWithBody, service)

		// Should have all error responses including 422
		expectedStatusCodes := []string{"400", "401", "403", "404", "409", "422", "429", "500"}
		assert.Equal(t, len(expectedStatusCodes), responses.Len(), "Should have all error responses for endpoint with body params")

		for _, statusCode := range expectedStatusCodes {
			response := responses.GetOrZero(statusCode)
			assert.NotNil(t, response, "Should have %s error response", statusCode)
			assert.NotEmpty(t, response.Description, "Error response %s should have description", statusCode)
			assert.NotNil(t, response.Content, "Error response %s should have content", statusCode)
			mediaType := response.Content.GetOrZero("application/json")
			assert.NotNil(t, mediaType, "Error response %s should have JSON content", statusCode)
		}
	})

	// Test without body params
	t.Run("without body params excludes 422", func(t *testing.T) {
		generator := NewGenerator()

		// Create service with ErrorCode enum
		service := &specification.Service{
			Name: "TestService",
			Enums: []specification.Enum{
				{
					Name:        "ErrorCode",
					Description: "Standard error codes used in API responses",
					Values: []specification.EnumValue{
						{Name: "BadRequest", Description: "Bad request error"},
						{Name: "UnprocessableEntity", Description: "Validation error"},
						{Name: "NotFound", Description: "Not found error"},
					},
				},
			},
		}

		// Test endpoint without body parameters (should not include 422)
		endpointWithoutBody := specification.Endpoint{
			Name:   "GetUser",
			Method: "GET",
			Path:   "/users/{id}",
			Request: specification.EndpointRequest{
				PathParams: []specification.Field{
					{Name: "id", Type: specification.FieldTypeUUID},
				},
			},
			Response: specification.EndpointResponse{StatusCode: 200},
		}

		responses := orderedmap.New[string, *v3.Response]()
		generator.addErrorResponses(responses, endpointWithoutBody, service)

		// Should have error responses but not 422
		assert.Equal(t, 2, responses.Len(), "Should have 2 error responses (excluding 422)")
		response400 := responses.GetOrZero("400")
		assert.NotNil(t, response400, "Should have 400 error response")
		response404 := responses.GetOrZero("404")
		assert.NotNil(t, response404, "Should have 404 error response")
		response422 := responses.GetOrZero("422")
		assert.Nil(t, response422, "Should not have 422 error response for endpoint without body params")
	})

	// Test without ErrorCode enum
	t.Run("without ErrorCode enum uses fallback responses", func(t *testing.T) {
		generator := NewGenerator()

		// Create service without ErrorCode enum
		service := &specification.Service{
			Name: "TestService",
			Enums: []specification.Enum{
				{Name: "SomeOtherEnum", Description: "Some other enum", Values: []specification.EnumValue{}},
			},
		}

		endpoint := specification.Endpoint{
			Name:     "TestEndpoint",
			Method:   "GET",
			Path:     "/test",
			Response: specification.EndpointResponse{StatusCode: 200},
		}

		responses := orderedmap.New[string, *v3.Response]()
		generator.addErrorResponses(responses, endpoint, service)

		// Should fall back to default error responses
		expectedDefaultStatusCodes := []string{"400", "401", "404", "500"}
		assert.Equal(t, len(expectedDefaultStatusCodes), responses.Len(), "Should have default error responses")

		for _, statusCode := range expectedDefaultStatusCodes {
			response := responses.GetOrZero(statusCode)
			assert.NotNil(t, response, "Should have %s default error response", statusCode)
			assert.NotEmpty(t, response.Description, "Default error response %s should have description", statusCode)
		}
	})
}

// TestMapErrorCodeToStatusAndDescription tests the error code to status code mapping.
func TestGenerator_mapErrorCodeToStatusAndDescription(t *testing.T) {
	generator := NewGenerator()

	testCases := []struct {
		errorCodeName        string
		errorCodeDescription string
		expectedStatus       string
		expectedDescription  string
	}{
		{"BadRequest", "Bad request error", "400", "Bad request error"},
		{"Unauthorized", "Unauthorized error", "401", "Unauthorized error"},
		{"Forbidden", "Forbidden error", "403", "Forbidden error"},
		{"NotFound", "Not found error", "404", "Not found error"},
		{"Conflict", "Conflict error", "409", "Conflict error"},
		{"UnprocessableEntity", "Validation error", "422", "Validation error"},
		{"RateLimited", "Rate limited error", "429", "Rate limited error"},
		{"Internal", "Internal error", "500", "Internal error"},
		{"UnknownError", "Unknown error", "500", "Unknown error"}, // Default to 500
	}

	for _, tc := range testCases {
		t.Run(tc.errorCodeName, func(t *testing.T) {
			status, description := generator.mapErrorCodeToStatusAndDescription(tc.errorCodeName, tc.errorCodeDescription)
			assert.Equal(t, tc.expectedStatus, status, "Status code should match for error code %s", tc.errorCodeName)
			assert.Equal(t, tc.expectedDescription, description, "Description should match for error code %s", tc.errorCodeName)
		})
	}
}

// TestEndToEndErrorResponseGeneration tests complete OpenAPI generation with proper error responses.
func TestEndToEndErrorResponseGeneration(t *testing.T) {
	generator := NewGenerator()

	// Create a service and apply overlay to get ErrorCode enum
	inputService := &specification.Service{
		Name:    "UserAPI",
		Version: "1.0.0",
		Resources: []specification.Resource{
			{
				Name:        "User",
				Description: "User resource",
				Operations:  []string{specification.OperationCreate, specification.OperationRead},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        "id",
							Description: "User identifier",
							Type:        specification.FieldTypeUUID,
						},
						Operations: []string{specification.OperationRead},
					},
					{
						Field: specification.Field{
							Name:        "email",
							Description: "User email address",
							Type:        specification.FieldTypeString,
						},
						Operations: []string{specification.OperationCreate, specification.OperationRead},
					},
				},
			},
		},
	}

	// Apply overlay to generate default endpoints and error handling
	service := specification.ApplyOverlay(inputService)

	// Generate OpenAPI document
	document, err := generator.GenerateFromService(service)
	assert.NoError(t, err, "Should generate OpenAPI document successfully")
	assert.NotNil(t, document, "Generated document should not be nil")

	// Convert to JSON to inspect the structure
	jsonBytes, err := generator.ToJSON(document)
	assert.NoError(t, err, "Should convert document to JSON successfully")

	jsonString := string(jsonBytes)

	// Verify presence of ErrorCode enum
	assert.Contains(t, jsonString, "ErrorCode", "Generated JSON should contain ErrorCode enum")
	assert.Contains(t, jsonString, "BadRequest", "Generated JSON should contain BadRequest error code")
	assert.Contains(t, jsonString, "UnprocessableEntity", "Generated JSON should contain UnprocessableEntity error code")

	// Verify presence of Error object
	assert.Contains(t, jsonString, "Error", "Generated JSON should contain Error object")

	// Verify error responses are present in endpoints
	assert.Contains(t, jsonString, "\"400\"", "Generated JSON should contain 400 error response")
	assert.Contains(t, jsonString, "\"401\"", "Generated JSON should contain 401 error response")
	assert.Contains(t, jsonString, "\"404\"", "Generated JSON should contain 404 error response")

	// Verify 422 is present for POST endpoints (with body params) but check structure
	assert.Contains(t, jsonString, "\"422\"", "Generated JSON should contain 422 error response for endpoints with body params")

	// Verify success responses are also present
	assert.Contains(t, jsonString, "\"200\"", "Generated JSON should contain 200 success response")
	assert.Contains(t, jsonString, "\"201\"", "Generated JSON should contain 201 success response")

	t.Logf("Generated OpenAPI JSON:\n%s", jsonString)
}

// TestCamelCaseParametersInOpenAPI verifies that parameters use camelCase in OpenAPI output
func TestCamelCaseParametersInOpenAPI(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name:    "TestAPI",
		Version: "1.0.0",
		Resources: []specification.Resource{
			{
				Name:        "User",
				Description: "User resource for testing camelCase",
				Operations:  []string{specification.OperationCreate, specification.OperationRead},
				Endpoints: []specification.Endpoint{
					{
						Name:        "GetUser",
						Title:       "Get User",
						Description: "Get user with parameters",
						Method:      "GET",
						Path:        "/{user_id}",
						Request: specification.EndpointRequest{
							PathParams: []specification.Field{
								{
									Name:        "user_id",
									Description: "User identifier",
									Type:        specification.FieldTypeUUID,
								},
							},
							QueryParams: []specification.Field{
								{
									Name:        "include_details",
									Description: "Include user details",
									Type:        specification.FieldTypeBool,
								},
								{
									Name:        "created_at_filter",
									Description: "Filter by creation date",
									Type:        specification.FieldTypeDate,
								},
							},
						},
						Response: specification.EndpointResponse{
							ContentType: "application/json",
							StatusCode:  200,
							BodyFields: []specification.Field{
								{
									Name:        "user_name",
									Description: "Name of the user",
									Type:        specification.FieldTypeString,
								},
								{
									Name:        "created_at",
									Description: "Creation timestamp",
									Type:        specification.FieldTypeTimestamp,
								},
							},
						},
					},
					{
						Name:        "CreateUser",
						Title:       "Create User",
						Description: "Create new user",
						Method:      "POST",
						Path:        "",
						Request: specification.EndpointRequest{
							ContentType: "application/json",
							BodyParams: []specification.Field{
								{
									Name:        "user_email",
									Description: "User email address",
									Type:        specification.FieldTypeString,
								},
								{
									Name:        "first_name",
									Description: "User first name",
									Type:        specification.FieldTypeString,
								},
							},
						},
						Response: specification.EndpointResponse{
							ContentType: "application/json",
							StatusCode:  201,
							BodyFields: []specification.Field{
								{
									Name:        "user_id",
									Description: "Created user ID",
									Type:        specification.FieldTypeUUID,
								},
							},
						},
					},
				},
			},
		},
	}

	document, err := generator.GenerateFromService(service)

	assert.NoError(t, err, "Should generate document successfully")
	assert.NotNil(t, document, "Document should not be nil")

	// Generate JSON to check the actual parameter names
	jsonBytes, err := generator.ToJSON(document)
	assert.NoError(t, err, "Should convert to JSON successfully")
	jsonString := string(jsonBytes)

	// Verify path parameters are in camelCase
	assert.Contains(t, jsonString, "\"userID\"", "Path parameter should be camelCase: userID")
	assert.NotContains(t, jsonString, "\"user_id\"", "Path parameter should not contain underscores: user_id")

	// Verify query parameters are in camelCase
	assert.Contains(t, jsonString, "\"includeDetails\"", "Query parameter should be camelCase: includeDetails")
	assert.Contains(t, jsonString, "\"createdAtFilter\"", "Query parameter should be camelCase: createdAtFilter")
	assert.NotContains(t, jsonString, "\"include_details\"", "Query parameter should not contain underscores: include_details")
	assert.NotContains(t, jsonString, "\"created_at_filter\"", "Query parameter should not contain underscores: created_at_filter")

	// Verify request body properties are in camelCase
	assert.Contains(t, jsonString, "\"userEmail\"", "Request body property should be camelCase: userEmail")
	assert.Contains(t, jsonString, "\"firstName\"", "Request body property should be camelCase: firstName")
	assert.NotContains(t, jsonString, "\"user_email\"", "Request body property should not contain underscores: user_email")
	assert.NotContains(t, jsonString, "\"first_name\"", "Request body property should not contain underscores: first_name")

	// Verify response body properties are in camelCase
	assert.Contains(t, jsonString, "\"userName\"", "Response body property should be camelCase: userName")
	assert.Contains(t, jsonString, "\"createdAt\"", "Response body property should be camelCase: createdAt")
	assert.NotContains(t, jsonString, "\"user_name\"", "Response body property should not contain underscores: user_name")
	assert.NotContains(t, jsonString, "\"created_at\"", "Response body property should not contain underscores: created_at")

	t.Logf("Generated OpenAPI JSON for camelCase verification:\n%s", jsonString)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}

// ============================================================================
// GenerateFromSpecificationToJSON Function Tests
// ============================================================================

// TestGenerateFromSpecificationToJSON tests the convenience method for generating JSON from a specification.
func TestGenerateFromSpecificationToJSON(t *testing.T) {
	// Test with nil service
	t.Run("nil service returns error", func(t *testing.T) {
		jsonData, err := GenerateFromSpecificationToJSON(nil)

		assert.Nil(t, jsonData, "JSON data should be nil when service is nil")
		assert.EqualError(t, err, "invalid service: service cannot be nil", "Should return invalid service error")
	})

	// Test with valid service
	t.Run("valid service generates JSON", func(t *testing.T) {
		service := &specification.Service{
			Name:    "TestService",
			Version: "1.0.0",
		}

		jsonData, err := GenerateFromSpecificationToJSON(service)

		assert.Nil(t, err, "Should not return error for valid service")
		assert.NotNil(t, jsonData, "JSON data should not be nil")
		assert.Greater(t, len(jsonData), 0, "JSON data should not be empty")

		// Verify it's valid JSON by checking basic structure
		jsonString := string(jsonData)
		assert.Contains(t, jsonString, "openapi", "Should contain OpenAPI version field")
		assert.Contains(t, jsonString, "3.1.0", "Should contain OpenAPI 3.1.0 version")
		assert.Contains(t, jsonString, "TestService API", "Should contain service name with API suffix")
		assert.Contains(t, jsonString, "Generated API documentation", "Should contain default description")
	})

	// Test that it produces same result as the multi-step process
	t.Run("produces same result as multi-step process", func(t *testing.T) {
		service := &specification.Service{
			Name:    "ComparisonService",
			Version: "2.0.0",
		}

		// Generate using convenience method
		convenienceJSON, err := GenerateFromSpecificationToJSON(service)
		assert.Nil(t, err, "Convenience method should not return error")

		// Generate using multi-step process
		generator := NewGenerator()
		generator.Title = service.Name + " API"
		generator.Description = "Generated API documentation"

		document, err := generator.GenerateFromService(service)
		assert.Nil(t, err, "Multi-step method should not return error")

		multiStepJSON, err := generator.ToJSON(document)
		assert.Nil(t, err, "Multi-step ToJSON should not return error")

		// Both methods should produce identical results
		assert.Equal(t, multiStepJSON, convenienceJSON, "Both methods should produce identical JSON")
	})
}

// ============================================================================
// End-to-End Integration Tests
// ============================================================================

// TestYAMLToOpenAPIEndToEnd tests the complete pipeline from YAML specification to OpenAPI JSON.
func TestYAMLToOpenAPIEndToEnd(t *testing.T) {
	t.Run("complete YAML to OpenAPI JSON pipeline", func(t *testing.T) {
		// YAML specification input
		yamlContent := `name: "School Management API"
version: "1.0.0"
servers:
  - url: "https://api.school.example.com/v1"
    description: "Production server for School Management API"

enums:
  - name: "StudentStatus"
    description: "Status of a student in the system"
    values:
      - name: "Active"
        description: "Student is actively enrolled"
      - name: "Inactive"  
        description: "Student is not currently enrolled"
      - name: "Graduated"
        description: "Student has graduated"

  - name: "GradeLevel"
    description: "Grade levels in the school"
    values:
      - name: "Elementary"
        description: "Elementary school level"
      - name: "Middle"
        description: "Middle school level"
      - name: "High"
        description: "High school level"

objects:
  - name: "Contact"
    description: "Contact information"
    fields:
      - name: "email"
        type: "String"
        description: "Email address"
      - name: "phone"
        type: "String"
        description: "Phone number"
        modifiers: ["nullable"]

  - name: "Address"
    description: "Physical address information"  
    fields:
      - name: "street"
        type: "String"
        description: "Street address"
      - name: "city"
        type: "String"
        description: "City"
      - name: "state"
        type: "String"
        description: "State or province"
      - name: "zipCode"
        type: "String" 
        description: "ZIP or postal code"

resources:
  - name: "Students"
    description: "Student management resource"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - name: "id"
        type: "UUID"
        description: "Unique student identifier"
        operations: ["Read"]
      - name: "firstName"
        type: "String"
        description: "Student's first name"
        operations: ["Create", "Read", "Update"]
      - name: "lastName"
        type: "String"
        description: "Student's last name"
        operations: ["Create", "Read", "Update"]
      - name: "studentId"
        type: "String"
        description: "School-assigned student ID"
        operations: ["Create", "Read"]
      - name: "status"
        type: "StudentStatus"
        description: "Current status of the student"
        default: "Active"
        operations: ["Create", "Read", "Update"]
      - name: "gradeLevel"
        type: "GradeLevel"
        description: "Current grade level"
        operations: ["Create", "Read", "Update"]
      - name: "contact"
        type: "Contact"
        description: "Student contact information"
        operations: ["Create", "Read", "Update"]
      - name: "address"
        type: "Address"
        description: "Student home address"
        modifiers: ["nullable"]
        operations: ["Create", "Read", "Update"]
      - name: "enrollmentDate"
        type: "Date"
        description: "Date when student was enrolled"
        operations: ["Create", "Read"]
      - name: "graduationDate"
        type: "Date"
        description: "Expected or actual graduation date"
        modifiers: ["nullable"]
        operations: ["Read", "Update"]`

		// Expected OpenAPI JSON output
		expectedJSON := `{
  "openapi": "3.1.0",
  "info": {
    "title": "School Management API API",
    "description": "Generated API documentation",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.school.example.com/v1",
      "description": "Production server for School Management API"
    }
  ],
  "paths": {
    "/students": {
      "get": {
        "tags": [
          "Students"
        ],
        "summary": "List all Students",
        "description": "Returns a paginated list of all ` + "`" + `Students` + "`" + ` in your organization.",
        "operationId": "List",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "The maximum number of items to return (default: 50)",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "The maximum number of items to return (default: 50)",
              "default": 50
            }
          },
          {
            "name": "offset",
            "in": "query",
            "description": "The number of items to skip before starting to return results (default: 0)",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "The number of items to skip before starting to return results (default: 0)",
              "default": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": {
                        "title": "Students"
                      },
                      "description": "Array of Students objects"
                    },
                    "pagination": {
                      "title": "Pagination",
                      "description": "Pagination information"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      },
      "post": {
        "tags": [
          "Students"
        ],
        "summary": "Create Students",
        "description": "Create a new Students",
        "operationId": "Create",
        "parameters": [],
        "requestBody": {
          "description": "Request body",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "firstName": {
                    "type": "string",
                    "description": "Student's first name"
                  },
                  "lastName": {
                    "type": "string",
                    "description": "Student's last name"
                  },
                  "studentId": {
                    "type": "string",
                    "description": "School-assigned student ID"
                  },
                  "status": {
                    "title": "StudentStatus",
                    "description": "Current status of the student",
                    "default": "Active"
                  },
                  "gradeLevel": {
                    "title": "GradeLevel",
                    "description": "Current grade level"
                  },
                  "contact": {
                    "title": "Contact",
                    "description": "Student contact information"
                  },
                  "address": {
                    "title": "Address",
                    "description": "Student home address",
                    "nullable": true
                  },
                  "enrollmentDate": {
                    "type": "string",
                    "format": "date",
                    "description": "Date when student was enrolled"
                  }
                },
                "required": [
                  "firstName",
                  "lastName",
                  "studentId",
                  "gradeLevel",
                  "enrollmentDate"
                ]
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Students"
                }
              }
            }
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "422": {
            "description": "The request was well-formed but failed validation (e.g. invalid field format or constraints), 422 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      }
    },
    "/students/{id}": {
      "get": {
        "tags": [
          "Students"
        ],
        "summary": "Retrieve an existing Students",
        "description": "Retrieves the ` + "`" + `Students` + "`" + ` with the given ID.",
        "operationId": "Get",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "The unique identifier of the Students to retrieve",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid",
              "description": "The unique identifier of the Students to retrieve"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Students"
                }
              }
            }
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Students"
        ],
        "summary": "Delete Students",
        "description": "Delete a Students",
        "operationId": "Delete",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "The unique identifier of the Students to delete",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid",
              "description": "The unique identifier of the Students to delete"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response"
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      },
      "patch": {
        "tags": [
          "Students"
        ],
        "summary": "Update Students",
        "description": "Update a Students",
        "operationId": "Update",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "The unique identifier of the Students to update",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid",
              "description": "The unique identifier of the Students to update"
            }
          }
        ],
        "requestBody": {
          "description": "Request body",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "firstName": {
                    "type": "string",
                    "description": "Student's first name"
                  },
                  "lastName": {
                    "type": "string",
                    "description": "Student's last name"
                  },
                  "status": {
                    "title": "StudentStatus",
                    "description": "Current status of the student",
                    "default": "Active"
                  },
                  "gradeLevel": {
                    "title": "GradeLevel",
                    "description": "Current grade level"
                  },
                  "contact": {
                    "title": "Contact",
                    "description": "Student contact information"
                  },
                  "address": {
                    "title": "Address",
                    "description": "Student home address",
                    "nullable": true
                  },
                  "graduationDate": {
                    "type": "string",
                    "format": "date",
                    "description": "Expected or actual graduation date",
                    "nullable": true
                  }
                },
                "required": [
                  "firstName",
                  "lastName",
                  "gradeLevel"
                ]
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Students"
                }
              }
            }
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "422": {
            "description": "The request was well-formed but failed validation (e.g. invalid field format or constraints), 422 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      }
    },
    "/students/_search": {
      "post": {
        "tags": [
          "Students"
        ],
        "summary": "Search Students",
        "description": "Search for ` + "`" + `Students` + "`" + ` with filtering capabilities.",
        "operationId": "Search",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "The maximum number of items to return (default: 50)",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "The maximum number of items to return (default: 50)",
              "default": 50
            }
          },
          {
            "name": "offset",
            "in": "query",
            "description": "The number of items to skip before starting to return results (default: 0)",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "The number of items to skip before starting to return results (default: 0)",
              "default": 0
            }
          }
        ],
        "requestBody": {
          "description": "Request body",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "filter": {
                    "type": "string",
                    "description": "Filter criteria to search for specific records"
                  }
                },
                "required": [
                  "filter"
                ]
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": {
                        "title": "Students"
                      },
                      "description": "Array of Students objects"
                    },
                    "pagination": {
                      "title": "Pagination",
                      "description": "Pagination information"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request was malformed or contained invalid parameters. 400 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "401": {
            "description": "The request is missing valid authentication credentials. 401 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "403": {
            "description": "Request is authenticated, but the user is not allowed to perform the operation. 403 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "404": {
            "description": "The requested resource or endpoint does not exist. This can happen if a resource ID is invalid or the route is unknown. 404 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "409": {
            "description": "The request could not be completed due to a conflict, such as a resource with dependencies that prevent deletion. 409 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "422": {
            "description": "The request was well-formed but failed validation (e.g. invalid field format or constraints), 422 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "429": {
            "description": "When the rate limit has been exceeded, 429 status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          },
          "500": {
            "description": "Some serverside issue, 5xx status code",
            "content": {
              "application/json": {
                "schema": {
                  "title": "Error"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "StudentStatus": {
        "type": "string",
        "enum": [
          "Active",
          "Inactive",
          "Graduated"
        ],
        "description": "Status of a student in the system"
      },
      "GradeLevel": {
        "type": "string",
        "enum": [
          "Elementary",
          "Middle",
          "High"
        ],
        "description": "Grade levels in the school"
      },
      "ErrorCode": {
        "type": "string",
        "enum": [
          "BadRequest",
          "Unauthorized",
          "Forbidden",
          "NotFound",
          "Conflict",
          "UnprocessableEntity",
          "RateLimited",
          "Internal"
        ],
        "description": "Standard error codes used in API responses"
      },
      "ErrorFieldCode": {
        "type": "string",
        "enum": [
          "AlreadyExists",
          "Required",
          "NotFound",
          "InvalidValue"
        ],
        "description": "Error codes for field-level validation errors"
      },
      "Contact": {
        "type": "object",
        "properties": {
          "email": {
            "type": "string",
            "description": "Email address"
          },
          "phone": {
            "type": "string",
            "description": "Phone number",
            "nullable": true
          }
        },
        "required": [
          "email"
        ],
        "description": "Contact information"
      },
      "Address": {
        "type": "object",
        "properties": {
          "street": {
            "type": "string",
            "description": "Street address"
          },
          "city": {
            "type": "string",
            "description": "City"
          },
          "state": {
            "type": "string",
            "description": "State or province"
          },
          "zipCode": {
            "type": "string",
            "description": "ZIP or postal code"
          }
        },
        "required": [
          "street",
          "city",
          "state",
          "zipCode"
        ],
        "description": "Physical address information"
      },
      "Error": {
        "type": "object",
        "properties": {
          "code": {
            "title": "ErrorCode",
            "description": "The specific error code indicating the type of error"
          },
          "message": {
            "type": "string",
            "description": "Human-readable error message providing additional details"
          }
        },
        "required": [
          "code",
          "message"
        ],
        "description": "Standard error response object containing error code and message"
      },
      "ErrorField": {
        "type": "object",
        "properties": {
          "code": {
            "title": "ErrorFieldCode",
            "description": "The specific error code indicating the type of field validation error"
          },
          "message": {
            "type": "string",
            "description": "Human-readable error message providing details about the field validation error"
          }
        },
        "required": [
          "code",
          "message"
        ],
        "description": "Field-specific error information containing error code and message for validation errors"
      },
      "Pagination": {
        "type": "object",
        "properties": {
          "offset": {
            "type": "integer",
            "description": "Number of items to skip from the beginning of the result set"
          },
          "limit": {
            "type": "integer",
            "description": "Maximum number of items to return in the result set"
          },
          "total": {
            "type": "integer",
            "description": "Total number of items available for pagination"
          }
        },
        "required": [
          "offset",
          "limit",
          "total"
        ],
        "description": "Pagination parameters for controlling result sets in list operations"
      },
      "Meta": {
        "type": "object",
        "properties": {
          "createdAt": {
            "type": "string",
            "format": "date-time",
            "description": "Timestamp when the resource was created"
          },
          "createdBy": {
            "type": "string",
            "format": "uuid",
            "description": "User who created the resource",
            "nullable": true
          },
          "updatedAt": {
            "type": "string",
            "format": "date-time",
            "description": "Timestamp when the resource was last updated"
          },
          "updatedBy": {
            "type": "string",
            "format": "uuid",
            "description": "User who last updated the resource",
            "nullable": true
          }
        },
        "required": [
          "createdAt",
          "updatedAt"
        ],
        "description": "Metadata fields containing creation and update information"
      },
      "Students": {
        "type": "object",
        "properties": {
          "iD": {
            "type": "string",
            "format": "uuid",
            "description": "Unique identifier for the Students"
          },
          "meta": {
            "title": "Meta",
            "description": "Metadata information for the Students"
          },
          "id": {
            "type": "string",
            "format": "uuid",
            "description": "Unique student identifier"
          },
          "firstName": {
            "type": "string",
            "description": "Student's first name"
          },
          "lastName": {
            "type": "string",
            "description": "Student's last name"
          },
          "studentId": {
            "type": "string",
            "description": "School-assigned student ID"
          },
          "status": {
            "title": "StudentStatus",
            "description": "Current status of the student",
            "default": "Active"
          },
          "gradeLevel": {
            "title": "GradeLevel",
            "description": "Current grade level"
          },
          "contact": {
            "title": "Contact",
            "description": "Student contact information"
          },
          "address": {
            "title": "Address",
            "description": "Student home address",
            "nullable": true
          },
          "enrollmentDate": {
            "type": "string",
            "format": "date",
            "description": "Date when student was enrolled"
          },
          "graduationDate": {
            "type": "string",
            "format": "date",
            "description": "Expected or actual graduation date",
            "nullable": true
          }
        },
        "required": [
          "iD",
          "id",
          "firstName",
          "lastName",
          "studentId",
          "gradeLevel",
          "enrollmentDate"
        ],
        "description": "Student management resource"
      },
      "ContactRequestError": {
        "type": "object",
        "properties": {
          "email": {
            "title": "ErrorField",
            "description": "Email address",
            "nullable": true
          },
          "phone": {
            "title": "ErrorField",
            "description": "Phone number",
            "nullable": true
          }
        },
        "description": "Request error object for Contact"
      },
      "AddressRequestError": {
        "type": "object",
        "properties": {
          "street": {
            "title": "ErrorField",
            "description": "Street address",
            "nullable": true
          },
          "city": {
            "title": "ErrorField",
            "description": "City",
            "nullable": true
          },
          "state": {
            "title": "ErrorField",
            "description": "State or province",
            "nullable": true
          },
          "zipCode": {
            "title": "ErrorField",
            "description": "ZIP or postal code",
            "nullable": true
          }
        },
        "description": "Request error object for Address"
      },
      "StudentsCreateRequestError": {
        "type": "object",
        "properties": {
          "firstName": {
            "title": "ErrorField",
            "description": "Student's first name",
            "nullable": true
          },
          "lastName": {
            "title": "ErrorField",
            "description": "Student's last name",
            "nullable": true
          },
          "studentId": {
            "title": "ErrorField",
            "description": "School-assigned student ID",
            "nullable": true
          },
          "status": {
            "title": "ErrorField",
            "description": "Current status of the student",
            "nullable": true
          },
          "gradeLevel": {
            "title": "ErrorField",
            "description": "Current grade level",
            "nullable": true
          },
          "contact": {
            "title": "ContactRequestError",
            "description": "Student contact information",
            "nullable": true
          },
          "address": {
            "title": "AddressRequestError",
            "description": "Student home address",
            "nullable": true
          },
          "enrollmentDate": {
            "title": "ErrorField",
            "description": "Date when student was enrolled",
            "nullable": true
          }
        },
        "description": "Request error object for Students Create endpoint"
      },
      "StudentsUpdateRequestError": {
        "type": "object",
        "properties": {
          "firstName": {
            "title": "ErrorField",
            "description": "Student's first name",
            "nullable": true
          },
          "lastName": {
            "title": "ErrorField",
            "description": "Student's last name",
            "nullable": true
          },
          "status": {
            "title": "ErrorField",
            "description": "Current status of the student",
            "nullable": true
          },
          "gradeLevel": {
            "title": "ErrorField",
            "description": "Current grade level",
            "nullable": true
          },
          "contact": {
            "title": "ContactRequestError",
            "description": "Student contact information",
            "nullable": true
          },
          "address": {
            "title": "AddressRequestError",
            "description": "Student home address",
            "nullable": true
          },
          "graduationDate": {
            "title": "ErrorField",
            "description": "Expected or actual graduation date",
            "nullable": true
          }
        },
        "description": "Request error object for Students Update endpoint"
      },
      "StudentsSearchRequestError": {
        "type": "object",
        "properties": {
          "filter": {
            "type": "string",
            "description": "Filter criteria to search for specific records",
            "nullable": true
          }
        },
        "description": "Request error object for Students Search endpoint"
      },
      "ContactFilter": {
        "type": "object",
        "properties": {
          "equals": {
            "title": "ContactFilterEquals",
            "description": "Equality filters for Contact",
            "nullable": true
          },
          "notEquals": {
            "title": "ContactFilterEquals",
            "description": "Inequality filters for Contact",
            "nullable": true
          },
          "greaterThan": {
            "title": "ContactFilterRange",
            "description": "Greater than filters for Contact",
            "nullable": true
          },
          "smallerThan": {
            "title": "ContactFilterRange",
            "description": "Smaller than filters for Contact",
            "nullable": true
          },
          "greaterOrEqual": {
            "title": "ContactFilterRange",
            "description": "Greater than or equal filters for Contact",
            "nullable": true
          },
          "smallerOrEqual": {
            "title": "ContactFilterRange",
            "description": "Smaller than or equal filters for Contact",
            "nullable": true
          },
          "contains": {
            "title": "ContactFilterContains",
            "description": "Contains filters for Contact",
            "nullable": true
          },
          "notContains": {
            "title": "ContactFilterContains",
            "description": "Not contains filters for Contact",
            "nullable": true
          },
          "like": {
            "title": "ContactFilterLike",
            "description": "LIKE filters for Contact",
            "nullable": true
          },
          "notLike": {
            "title": "ContactFilterLike",
            "description": "NOT LIKE filters for Contact",
            "nullable": true
          },
          "null": {
            "title": "ContactFilterNull",
            "description": "Null filters for Contact",
            "nullable": true
          },
          "notNull": {
            "title": "ContactFilterNull",
            "description": "Not null filters for Contact",
            "nullable": true
          },
          "orCondition": {
            "type": "boolean",
            "description": "OrCondition decides if this filter is within an OR-condition or AND-condition"
          },
          "nestedFilters": {
            "type": "array",
            "items": {
              "title": "ContactFilter"
            },
            "description": "NestedFilters of the Contact, useful for more complex filters"
          }
        },
        "required": [
          "orCondition"
        ],
        "description": "Filter object for Contact"
      },
      "ContactFilterEquals": {
        "type": "object",
        "properties": {
          "email": {
            "type": "string",
            "description": "Email address",
            "nullable": true
          },
          "phone": {
            "type": "string",
            "description": "Phone number",
            "nullable": true
          }
        },
        "description": "Equality/Inequality filter fields for Contact"
      },
      "ContactFilterRange": {
        "type": "object",
        "description": "Range filter fields for Contact"
      },
      "ContactFilterContains": {
        "type": "object",
        "properties": {
          "email": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "Email address"
          },
          "phone": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "Phone number"
          }
        },
        "description": "Contains filter fields for Contact"
      },
      "ContactFilterLike": {
        "type": "object",
        "properties": {
          "email": {
            "type": "string",
            "description": "Email address",
            "nullable": true
          },
          "phone": {
            "type": "string",
            "description": "Phone number",
            "nullable": true
          }
        },
        "description": "LIKE filter fields for Contact"
      },
      "ContactFilterNull": {
        "type": "object",
        "properties": {
          "phone": {
            "type": "boolean",
            "description": "Phone number",
            "nullable": true
          }
        },
        "description": "Null filter fields for Contact"
      },
      "AddressFilter": {
        "type": "object",
        "properties": {
          "equals": {
            "title": "AddressFilterEquals",
            "description": "Equality filters for Address",
            "nullable": true
          },
          "notEquals": {
            "title": "AddressFilterEquals",
            "description": "Inequality filters for Address",
            "nullable": true
          },
          "greaterThan": {
            "title": "AddressFilterRange",
            "description": "Greater than filters for Address",
            "nullable": true
          },
          "smallerThan": {
            "title": "AddressFilterRange",
            "description": "Smaller than filters for Address",
            "nullable": true
          },
          "greaterOrEqual": {
            "title": "AddressFilterRange",
            "description": "Greater than or equal filters for Address",
            "nullable": true
          },
          "smallerOrEqual": {
            "title": "AddressFilterRange",
            "description": "Smaller than or equal filters for Address",
            "nullable": true
          },
          "contains": {
            "title": "AddressFilterContains",
            "description": "Contains filters for Address",
            "nullable": true
          },
          "notContains": {
            "title": "AddressFilterContains",
            "description": "Not contains filters for Address",
            "nullable": true
          },
          "like": {
            "title": "AddressFilterLike",
            "description": "LIKE filters for Address",
            "nullable": true
          },
          "notLike": {
            "title": "AddressFilterLike",
            "description": "NOT LIKE filters for Address",
            "nullable": true
          },
          "null": {
            "title": "AddressFilterNull",
            "description": "Null filters for Address",
            "nullable": true
          },
          "notNull": {
            "title": "AddressFilterNull",
            "description": "Not null filters for Address",
            "nullable": true
          },
          "orCondition": {
            "type": "boolean",
            "description": "OrCondition decides if this filter is within an OR-condition or AND-condition"
          },
          "nestedFilters": {
            "type": "array",
            "items": {
              "title": "AddressFilter"
            },
            "description": "NestedFilters of the Address, useful for more complex filters"
          }
        },
        "required": [
          "orCondition"
        ],
        "description": "Filter object for Address"
      },
      "AddressFilterEquals": {
        "type": "object",
        "properties": {
          "street": {
            "type": "string",
            "description": "Street address",
            "nullable": true
          },
          "city": {
            "type": "string",
            "description": "City",
            "nullable": true
          },
          "state": {
            "type": "string",
            "description": "State or province",
            "nullable": true
          },
          "zipCode": {
            "type": "string",
            "description": "ZIP or postal code",
            "nullable": true
          }
        },
        "description": "Equality/Inequality filter fields for Address"
      },
      "AddressFilterRange": {
        "type": "object",
        "description": "Range filter fields for Address"
      },
      "AddressFilterContains": {
        "type": "object",
        "properties": {
          "street": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "Street address"
          },
          "city": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "City"
          },
          "state": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "State or province"
          },
          "zipCode": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "ZIP or postal code"
          }
        },
        "description": "Contains filter fields for Address"
      },
      "AddressFilterLike": {
        "type": "object",
        "properties": {
          "street": {
            "type": "string",
            "description": "Street address",
            "nullable": true
          },
          "city": {
            "type": "string",
            "description": "City",
            "nullable": true
          },
          "state": {
            "type": "string",
            "description": "State or province",
            "nullable": true
          },
          "zipCode": {
            "type": "string",
            "description": "ZIP or postal code",
            "nullable": true
          }
        },
        "description": "LIKE filter fields for Address"
      },
      "AddressFilterNull": {
        "type": "object",
        "description": "Null filter fields for Address"
      }
    }
  }
}`

		// Step 1: Parse YAML specification
		service, err := specification.ParseServiceFromYAML([]byte(yamlContent))
		assert.NoError(t, err, "Should successfully parse YAML specification")
		assert.NotNil(t, service, "Parsed service should not be nil")

		// Step 2: Generate OpenAPI JSON using the convenience method
		actualJSON, err := GenerateFromSpecificationToJSON(service)
		assert.NoError(t, err, "Should successfully generate OpenAPI JSON from service")
		assert.NotNil(t, actualJSON, "Generated JSON should not be nil")

		// Step 3: Assert exact JSON match
		assert.JSONEq(t, expectedJSON, string(actualJSON), "Generated OpenAPI JSON should exactly match expected output")
	})
}
