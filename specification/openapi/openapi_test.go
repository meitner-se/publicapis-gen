package openapi

import (
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	expectedErrorInvalidService  = "invalid service: service cannot be nil"
	expectedErrorInvalidDocument = "invalid document: document cannot be nil"
	expectedVersion              = "3.1.0"
	emptyString                  = ""
	testServiceName              = "TestService"
	userAPIServiceName           = "UserAPI"
	customAPITitle               = "Custom API"
	customAPIDescription         = "Custom API Description"
	customServerURL              = "https://custom.example.com"
	testServiceVersion           = "2.0.0"
	prodServerURL                = "https://api.example.com"
	stagingServerURL             = "https://staging-api.example.com"
	prodServerDescription        = "Production server"
	stagingServerDescription     = "Staging server"
	statusEnumName               = "Status"
	statusEnumDescription        = "User status enumeration"
	activeEnumValue              = "Active"
	inactiveEnumValue            = "Inactive"
	activeEnumDescription        = "User is active"
	inactiveEnumDescription      = "User is inactive"
	userObjectName               = "User"
	userObjectDescription        = "User object"
	idFieldName                  = "id"
	idFieldDescription           = "User identifier"
	emailFieldName               = "email"
	emailFieldDescription        = "User email address"
	statusFieldName              = "status"
	statusFieldDescription       = "User status"
	userResourceDescription      = "User resource"
	createEndpointName           = "Create"
	createEndpointTitle          = "Create User"
	createEndpointDescription    = "Create a new user"
	postMethod                   = "POST"
	pathEmpty                    = ""
	contentTypeJSON              = "application/json"
	statusCode201                = 201

	// Error response test constants
	errorCodeEnumName            = "ErrorCode"
	errorCodeEnumDescription     = "Standard error codes used in API responses"
	badRequestErrorCode          = "BadRequest"
	unauthorizedErrorCode        = "Unauthorized"
	forbiddenErrorCode           = "Forbidden"
	notFoundErrorCode            = "NotFound"
	conflictErrorCode            = "Conflict"
	unprocessableEntityErrorCode = "UnprocessableEntity"
	rateLimitedErrorCode         = "RateLimited"
	internalErrorCode            = "Internal"
	badRequestDescription        = "The request was malformed or contained invalid parameters. 400 status code"
	unauthorizedDescription      = "The request is missing valid authentication credentials. 401 status code"
	forbiddenDescription         = "Request is authenticated, but the user is not allowed to perform the operation. 403 status code"
	notFoundDescription          = "The requested resource or endpoint does not exist. 404 status code"
	conflictDescription          = "The request could not be completed due to a conflict. 409 status code"
	unprocessableDescription     = "The request was well-formed but failed validation. 422 status code"
	rateLimitedDescription       = "When the rate limit has been exceeded. 429 status code"
	internalDescription          = "Some serverside issue. 500 status code"

	// Error object constants
	errorObjectName              = "Error"
	errorObjectDescription       = "Standard error response object containing error code and message"
	errorCodeFieldName           = "Code"
	errorCodeFieldDescription    = "The specific error code"
	errorMessageFieldName        = "Message"
	errorMessageFieldDescription = "Human-readable error message"

	// Test endpoint constants
	createUserEndpointName     = "CreateUser"
	getUserEndpointName        = "GetUser"
	testEndpointName           = "TestEndpoint"
	usersPath                  = "/users"
	userIDPath                 = "/users/{id}"
	testPath                   = "/test"
	getMethod                  = "GET"
	statusCode200              = 200
	expectedStatusCodes        = 8
	expectedDefaultStatusCodes = 4
)

// ============================================================================
// NewGenerator Function Tests
// ============================================================================

// TestNewGenerator tests the creation of a new OpenAPI generator.
func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()

	assert.NotNil(t, generator, "Generator should not be nil")
	assert.Equal(t, expectedVersion, generator.Version, "Generator version should be 3.1.0")
	assert.Equal(t, emptyString, generator.Title, "Generator title should be empty by default")
	assert.Equal(t, emptyString, generator.Description, "Generator description should be empty by default")
	assert.Equal(t, emptyString, generator.ServerURL, "Generator server URL should be empty by default")
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
		assert.EqualError(t, err, expectedErrorInvalidService, "Should return invalid service error")
	})

	// Test with valid service
	t.Run("valid service generates document", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: testServiceName,
		}

		document, err := generator.GenerateFromService(service)

		assert.NotNil(t, document, "Document should not be nil with valid service")
		assert.NoError(t, err, "Should not return error with valid service")
		assert.Equal(t, expectedVersion, document.Version, "Document version should be 3.1.0")
		assert.Equal(t, testServiceName, document.Info.Title, "Document title should match service name")
	})

	// Test with complex service
	t.Run("complex service with enums and objects", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: userAPIServiceName,
			Enums: []specification.Enum{
				{
					Name:        statusEnumName,
					Description: statusEnumDescription,
					Values: []specification.EnumValue{
						{Name: activeEnumValue, Description: activeEnumDescription},
						{Name: inactiveEnumValue, Description: inactiveEnumDescription},
					},
				},
			},
			Objects: []specification.Object{
				{
					Name:        userObjectName,
					Description: userObjectDescription,
					Fields: []specification.Field{
						{
							Name:        idFieldName,
							Description: idFieldDescription,
							Type:        specification.FieldTypeUUID,
						},
						{
							Name:        emailFieldName,
							Description: emailFieldDescription,
							Type:        specification.FieldTypeString,
						},
						{
							Name:        statusFieldName,
							Description: statusFieldDescription,
							Type:        statusEnumName,
						},
					},
				},
			},
			Resources: []specification.Resource{
				{
					Name:        userObjectName,
					Description: userResourceDescription,
					Operations:  []string{specification.OperationCreate, specification.OperationRead},
					Fields: []specification.ResourceField{
						{
							Field: specification.Field{
								Name:        idFieldName,
								Description: idFieldDescription,
								Type:        specification.FieldTypeUUID,
							},
							Operations: []string{specification.OperationRead},
						},
						{
							Field: specification.Field{
								Name:        emailFieldName,
								Description: emailFieldDescription,
								Type:        specification.FieldTypeString,
							},
							Operations: []string{specification.OperationCreate, specification.OperationRead},
						},
					},
					Endpoints: []specification.Endpoint{
						{
							Name:        createEndpointName,
							Title:       createEndpointTitle,
							Description: createEndpointDescription,
							Method:      postMethod,
							Path:        pathEmpty,
							Request: specification.EndpointRequest{
								ContentType: contentTypeJSON,
								BodyParams: []specification.Field{
									{
										Name:        emailFieldName,
										Description: emailFieldDescription,
										Type:        specification.FieldTypeString,
									},
								},
							},
							Response: specification.EndpointResponse{
								ContentType: contentTypeJSON,
								StatusCode:  statusCode201,
								BodyObject:  stringPtr(userObjectName),
							},
						},
					},
				},
			},
		}

		document, err := generator.GenerateFromService(service)

		assert.NoError(t, err, "Should generate document successfully")
		assert.NotNil(t, document, "Document should not be nil")
		assert.Equal(t, userAPIServiceName, document.Info.Title, "Document title should match service name")

		// Test JSON output contains expected elements
		jsonBytes, err := generator.ToJSON(document)
		assert.NoError(t, err, "Should convert document to JSON successfully")
		jsonString := string(jsonBytes)
		assert.Contains(t, jsonString, statusEnumName, "JSON should contain Status enum")
		assert.Contains(t, jsonString, userObjectName, "JSON should contain User object")
		assert.Contains(t, jsonString, activeEnumValue, "JSON should contain enum values")
		assert.Contains(t, jsonString, "/user", "JSON should contain user path")
	})

	// Test with service version and servers
	t.Run("service with version and servers", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name:    userAPIServiceName,
			Version: testServiceVersion,
			Servers: []specification.ServiceServer{
				{
					URL:         prodServerURL,
					Description: prodServerDescription,
				},
				{
					URL:         stagingServerURL,
					Description: stagingServerDescription,
				},
			},
			Objects: []specification.Object{
				{
					Name:        userObjectName,
					Description: userObjectDescription,
					Fields: []specification.Field{
						{
							Name:        emailFieldName,
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
		assert.Equal(t, testServiceVersion, document.Info.Version, "Document version should come from service")
		assert.Equal(t, 2, len(document.Servers), "Document should have 2 servers from service")
		assert.Equal(t, prodServerURL, document.Servers[0].URL, "First server URL should match service")
		assert.Equal(t, prodServerDescription, document.Servers[0].Description, "First server description should match service")
		assert.Equal(t, stagingServerURL, document.Servers[1].URL, "Second server URL should match service")
		assert.Equal(t, stagingServerDescription, document.Servers[1].Description, "Second server description should match service")

		// Test JSON output
		jsonBytes, err := generator.ToJSON(document)
		assert.NoError(t, err, "Should convert document to JSON successfully")
		jsonString := string(jsonBytes)
		assert.Contains(t, jsonString, testServiceVersion, "JSON should contain service version")
		assert.Contains(t, jsonString, prodServerURL, "JSON should contain first server URL")
		assert.Contains(t, jsonString, stagingServerURL, "JSON should contain second server URL")
	})
}

// TestGenerator_ToYAML tests YAML conversion functionality.
func TestGenerator_ToYAML(t *testing.T) {
	// Test with nil document
	t.Run("nil document returns error", func(t *testing.T) {
		generator := NewGenerator()

		yamlBytes, err := generator.ToYAML(nil)

		assert.Nil(t, yamlBytes, "YAML bytes should be nil when document is nil")
		assert.EqualError(t, err, expectedErrorInvalidDocument, "Should return invalid document error")
	})

	// Test with valid document
	t.Run("valid document converts successfully", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: testServiceName,
		}

		document, err := generator.GenerateFromService(service)
		assert.NoError(t, err, "Should generate document successfully")

		yamlBytes, err := generator.ToYAML(document)

		assert.NoError(t, err, "Should convert document to YAML successfully")
		assert.NotNil(t, yamlBytes, "YAML bytes should not be nil")
		assert.Contains(t, string(yamlBytes), testServiceName, "YAML should contain service name")
		assert.Contains(t, string(yamlBytes), expectedVersion, "YAML should contain OpenAPI version")
	})
}

// TestGenerator_ToJSON tests JSON conversion functionality.
func TestGenerator_ToJSON(t *testing.T) {
	// Test with nil document
	t.Run("nil document returns error", func(t *testing.T) {
		generator := NewGenerator()

		jsonBytes, err := generator.ToJSON(nil)

		assert.Nil(t, jsonBytes, "JSON bytes should be nil when document is nil")
		assert.EqualError(t, err, expectedErrorInvalidDocument, "Should return invalid document error")
	})

	// Test with valid document
	t.Run("valid document converts successfully", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name: testServiceName,
		}

		document, err := generator.GenerateFromService(service)
		assert.NoError(t, err, "Should generate document successfully")

		jsonBytes, err := generator.ToJSON(document)

		assert.NoError(t, err, "Should convert document to JSON successfully")
		assert.NotNil(t, jsonBytes, "JSON bytes should not be nil")
		assert.Contains(t, string(jsonBytes), testServiceName, "JSON should contain service name")
		assert.Contains(t, string(jsonBytes), expectedVersion, "JSON should contain OpenAPI version")
	})
}

// TestGeneratorConfiguration tests generator configuration options.
func TestGeneratorConfiguration(t *testing.T) {
	generator := &Generator{
		Version:     expectedVersion,
		Title:       customAPITitle,
		Description: customAPIDescription,
		ServerURL:   customServerURL,
	}

	assert.Equal(t, customAPITitle, generator.Title, "Generator title should match configured value")
	assert.Equal(t, customAPIDescription, generator.Description, "Generator description should match configured value")
	assert.Equal(t, customServerURL, generator.ServerURL, "Generator server URL should match configured value")
}

// TestGenerator_addErrorResponses tests error response generation functionality.
func TestGenerator_addErrorResponses(t *testing.T) {
	generator := NewGenerator()
	service := &specification.Service{
		Name: userAPIServiceName,
		Enums: []specification.Enum{
			{
				Name:        statusEnumName,
				Description: statusEnumDescription,
				Values: []specification.EnumValue{
					{Name: activeEnumValue, Description: activeEnumDescription},
					{Name: inactiveEnumValue, Description: inactiveEnumDescription},
				},
			},
		},
		Objects: []specification.Object{
			{
				Name:        userObjectName,
				Description: userObjectDescription,
				Fields: []specification.Field{
					{
						Name:        idFieldName,
						Description: idFieldDescription,
						Type:        specification.FieldTypeUUID,
					},
					{
						Name:        emailFieldName,
						Description: emailFieldDescription,
						Type:        specification.FieldTypeString,
					},
					{
						Name:        statusFieldName,
						Description: statusFieldDescription,
						Type:        statusEnumName,
					},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name:        userObjectName,
				Description: userResourceDescription,
				Operations:  []string{specification.OperationCreate, specification.OperationRead},
				Fields: []specification.ResourceField{
					{
						Field: specification.Field{
							Name:        idFieldName,
							Description: idFieldDescription,
							Type:        specification.FieldTypeUUID,
						},
						Operations: []string{specification.OperationRead},
					},
					{
						Field: specification.Field{
							Name:        emailFieldName,
							Description: emailFieldDescription,
							Type:        specification.FieldTypeString,
						},
						Operations: []string{specification.OperationCreate, specification.OperationRead},
					},
				},
				Endpoints: []specification.Endpoint{
					{
						Name:        createEndpointName,
						Title:       createEndpointTitle,
						Description: createEndpointDescription,
						Method:      postMethod,
						Path:        pathEmpty,
						Request: specification.EndpointRequest{
							ContentType: contentTypeJSON,
							BodyParams: []specification.Field{
								{
									Name:        emailFieldName,
									Description: emailFieldDescription,
									Type:        specification.FieldTypeString,
								},
							},
						},
						Response: specification.EndpointResponse{
							ContentType: contentTypeJSON,
							StatusCode:  statusCode201,
							BodyObject:  stringPtr(userObjectName),
						},
					},
				},
			},
		},
	}

	document, err := generator.GenerateFromService(service)

	assert.NoError(t, err, "Should generate document successfully")
	assert.NotNil(t, document, "Document should not be nil")
	assert.Equal(t, userAPIServiceName, document.Info.Title, "Document title should match service name")

	// Test JSON output contains expected elements
	jsonBytes, err := generator.ToJSON(document)
	assert.NoError(t, err, "Should convert document to JSON successfully")
	jsonString := string(jsonBytes)
	assert.Contains(t, jsonString, statusEnumName, "JSON should contain Status enum")
	assert.Contains(t, jsonString, userObjectName, "JSON should contain User object")
	assert.Contains(t, jsonString, activeEnumValue, "JSON should contain enum values")
	assert.Contains(t, jsonString, "/user", "JSON should contain user path")
}

// ============================================================================
// Error Response Tests
// ============================================================================

// TestGenerator_addErrorResponses tests error response generation functionality.
func TestGenerator_addErrorResponses(t *testing.T) {
	// Test with ErrorCode enum
	t.Run("with ErrorCode enum generates all responses", func(t *testing.T) {
		generator := NewGenerator()
		service := &specification.Service{
			Name:    userAPIServiceName,
			Version: testServiceVersion,
			Servers: []specification.ServiceServer{
				{
					URL:         prodServerURL,
					Description: prodServerDescription,
				},
				{
					URL:         stagingServerURL,
					Description: stagingServerDescription,
				},
			},
			Objects: []specification.Object{
				{
					Name:        userObjectName,
					Description: userObjectDescription,
					Fields: []specification.Field{
						{
							Name:        emailFieldName,
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
		assert.Equal(t, testServiceVersion, document.Info.Version, "Document version should come from service")
		assert.Equal(t, 2, len(document.Servers), "Document should have 2 servers from service")
		assert.Equal(t, prodServerURL, document.Servers[0].URL, "First server URL should match service")
		assert.Equal(t, prodServerDescription, document.Servers[0].Description, "First server description should match service")
		assert.Equal(t, stagingServerURL, document.Servers[1].URL, "Second server URL should match service")
		assert.Equal(t, stagingServerDescription, document.Servers[1].Description, "Second server description should match service")

		// Test JSON output
		jsonBytes, err := generator.ToJSON(document)
		assert.NoError(t, err, "Should convert document to JSON successfully")
		jsonString := string(jsonBytes)
		assert.Contains(t, jsonString, testServiceVersion, "JSON should contain service version")
		assert.Contains(t, jsonString, prodServerURL, "JSON should contain first server URL")
		assert.Contains(t, jsonString, stagingServerURL, "JSON should contain second server URL")
	})
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
			Name: testServiceName,
			Enums: []specification.Enum{
				{
					Name:        errorCodeEnumName,
					Description: errorCodeEnumDescription,
					Values: []specification.EnumValue{
						{Name: badRequestErrorCode, Description: badRequestDescription},
						{Name: unauthorizedErrorCode, Description: unauthorizedDescription},
						{Name: forbiddenErrorCode, Description: forbiddenDescription},
						{Name: notFoundErrorCode, Description: notFoundDescription},
						{Name: conflictErrorCode, Description: conflictDescription},
						{Name: unprocessableEntityErrorCode, Description: unprocessableDescription},
						{Name: rateLimitedErrorCode, Description: rateLimitedDescription},
						{Name: internalErrorCode, Description: internalDescription},
					},
				},
			},
			Objects: []specification.Object{
				{
					Name:        errorObjectName,
					Description: errorObjectDescription,
					Fields: []specification.Field{
						{Name: errorCodeFieldName, Description: errorCodeFieldDescription, Type: errorCodeEnumName},
						{Name: errorMessageFieldName, Description: errorMessageFieldDescription, Type: specification.FieldTypeString},
					},
				},
			},
		}

		// Test endpoint with body parameters (should include 422)
		endpointWithBody := specification.Endpoint{
			Name:   createUserEndpointName,
			Method: postMethod,
			Path:   usersPath,
			Request: specification.EndpointRequest{
				BodyParams: []specification.Field{
					{Name: emailFieldName, Type: specification.FieldTypeString},
				},
			},
			Response: specification.EndpointResponse{StatusCode: statusCode201},
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
			mediaType := response.Content.GetOrZero(contentTypeJSON)
			assert.NotNil(t, mediaType, "Error response %s should have JSON content", statusCode)
		}
	})

	// Test without body params
	t.Run("without body params excludes 422", func(t *testing.T) {
		generator := NewGenerator()

		// Create service with ErrorCode enum
		service := &specification.Service{
			Name: testServiceName,
			Enums: []specification.Enum{
				{
					Name:        errorCodeEnumName,
					Description: errorCodeEnumDescription,
					Values: []specification.EnumValue{
						{Name: badRequestErrorCode, Description: "Bad request error"},
						{Name: unprocessableEntityErrorCode, Description: "Validation error"},
						{Name: notFoundErrorCode, Description: "Not found error"},
					},
				},
			},
		}

		// Test endpoint without body parameters (should not include 422)
		endpointWithoutBody := specification.Endpoint{
			Name:   getUserEndpointName,
			Method: getMethod,
			Path:   userIDPath,
			Request: specification.EndpointRequest{
				PathParams: []specification.Field{
					{Name: idFieldName, Type: specification.FieldTypeUUID},
				},
			},
			Response: specification.EndpointResponse{StatusCode: statusCode200},
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
			Name: testServiceName,
			Enums: []specification.Enum{
				{Name: "SomeOtherEnum", Description: "Some other enum", Values: []specification.EnumValue{}},
			},
		}

		endpoint := specification.Endpoint{
			Name:     testEndpointName,
			Method:   getMethod,
			Path:     testPath,
			Response: specification.EndpointResponse{StatusCode: statusCode200},
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

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
