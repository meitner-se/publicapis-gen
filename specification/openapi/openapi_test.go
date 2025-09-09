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
