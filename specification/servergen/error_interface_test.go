package servergen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	testErrorStructName        = "Error"
	testErrorCodeField         = "Code"
	testErrorMessageField      = "Message"
	testErrorRequestIDField    = "RequestID"
	testErrorCodeEnumName      = "ErrorCode"
	testErrorMethodSignature   = "func (e *Error) Error() string"
	testErrorMethodBody        = "return e.Message.String()"
	testHTTPStatusCodeMethod   = "func (e *Error) HTTPStatusCode() int"
)

// ============================================================================
// Error Interface Implementation Tests
// ============================================================================

func TestErrorInterfaceImplementation(t *testing.T) {
	// Arrange
	service := &specification.Service{
		Name:    "TestService",
		Version: "v1",
		Objects: []specification.Object{
			{
				Name:        testErrorStructName,
				Description: "Error response object",
				Fields: []specification.Field{
					{Name: testErrorCodeField, Type: testErrorCodeEnumName},
					{Name: testErrorMessageField, Type: "String"},
					{Name: testErrorRequestIDField, Type: "String"},
				},
			},
		},
		Enums: []specification.Enum{
			{
				Name: testErrorCodeEnumName,
				Values: []specification.EnumValue{
					{Name: "BadRequest", Description: "Bad request"},
					{Name: "Unauthorized", Description: "Unauthorized"},
					{Name: "Forbidden", Description: "Forbidden"},
					{Name: "NotFound", Description: "Not found"},
					{Name: "Conflict", Description: "Conflict"},
					{Name: "UnprocessableEntity", Description: "Unprocessable entity"},
					{Name: "RateLimited", Description: "Rate limited"},
					{Name: "Internal", Description: "Internal server error"},
				},
			},
		},
	}

	// Act
	var buf bytes.Buffer
	err := GenerateServer(&buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating server with Error object")
	
	generatedCode := buf.String()
	
	// Verify the Error struct is generated
	assert.Contains(t, generatedCode, "type Error struct {",
		"Generated code should contain Error struct definition")
	
	// Verify the Error() method is generated with correct signature
	assert.Contains(t, generatedCode, testErrorMethodSignature,
		"Generated code should contain Error() method with correct signature")
	
	// Verify the Error() method returns the message
	assert.Contains(t, generatedCode, testErrorMethodBody,
		"Error() method should return e.Message.String()")
	
	// Verify HTTPStatusCode method is also generated
	assert.Contains(t, generatedCode, testHTTPStatusCodeMethod,
		"Generated code should contain HTTPStatusCode() method")
	
	// Verify all error code cases are handled
	errorCodes := []string{
		"BadRequest", "Unauthorized", "Forbidden", "NotFound",
		"Conflict", "UnprocessableEntity", "RateLimited", "Internal",
	}
	
	for _, code := range errorCodes {
		assert.Contains(t, generatedCode, "case ErrorCode"+code+":",
			"HTTPStatusCode method should handle ErrorCode"+code)
	}
	
	// Verify default case
	assert.Contains(t, generatedCode, "default:",
		"HTTPStatusCode method should have a default case")
	assert.Contains(t, generatedCode, "return http.StatusInternalServerError",
		"Default case should return 500 Internal Server Error")

	t.Run("method placement", func(t *testing.T) {
		// Verify the Error() method appears immediately after the Error struct
		errorStructIndex := strings.Index(generatedCode, "type Error struct {")
		errorMethodIndex := strings.Index(generatedCode, testErrorMethodSignature)
		httpStatusMethodIndex := strings.Index(generatedCode, testHTTPStatusCodeMethod)
		
		assert.Greater(t, errorMethodIndex, errorStructIndex,
			"Error() method should appear after Error struct definition")
		assert.Greater(t, httpStatusMethodIndex, errorMethodIndex,
			"HTTPStatusCode() method should appear after Error() method")
	})

	t.Run("no Error object", func(t *testing.T) {
		// Arrange
		serviceNoError := &specification.Service{
			Name:    "TestService",
			Version: "v1",
			Objects: []specification.Object{
				{
					Name:        "User",
					Description: "User object",
					Fields: []specification.Field{
						{Name: "ID", Type: "UUID"},
						{Name: "Name", Type: "String"},
					},
				},
			},
		}
		
		// Act
		var buf bytes.Buffer
		err := GenerateServer(&buf, serviceNoError)
		
		// Assert
		assert.Nil(t, err, "Expected no error when generating server without Error object")
		
		generatedCode := buf.String()
		assert.NotContains(t, generatedCode, testErrorMethodSignature,
			"Should not generate Error() method when there's no Error object")
		assert.NotContains(t, generatedCode, testHTTPStatusCodeMethod,
			"Should not generate HTTPStatusCode() method when there's no Error object")
	})

	t.Run("Error object with minimal fields", func(t *testing.T) {
		// Arrange
		serviceMinimalError := &specification.Service{
			Name:    "TestService",
			Version: "v1",
			Objects: []specification.Object{
				{
					Name:        testErrorStructName,
					Description: "Minimal error object",
					Fields: []specification.Field{
						{Name: testErrorMessageField, Type: "String"},
					},
				},
			},
		}
		
		// Act
		var buf bytes.Buffer
		err := GenerateServer(&buf, serviceMinimalError)
		
		// Assert
		assert.Nil(t, err, "Expected no error with minimal Error object")
		
		generatedCode := buf.String()
		assert.Contains(t, generatedCode, testErrorMethodSignature,
			"Should still generate Error() method for minimal Error object")
		assert.Contains(t, generatedCode, testErrorMethodBody,
			"Error() method should still return e.Message.String()")
	})
}

// ============================================================================
// Error Method Correctness Tests
// ============================================================================

func TestErrorMethodCorrectness(t *testing.T) {
	// Arrange
	service := &specification.Service{
		Name:    "TestService", 
		Version: "v1",
		Objects: []specification.Object{
			{
				Name:        testErrorStructName,
				Description: "Error object for testing",
				Fields: []specification.Field{
					{Name: testErrorCodeField, Type: testErrorCodeEnumName},
					{Name: testErrorMessageField, Type: "String"},
					{Name: "Details", Type: "String", Modifiers: []string{"Nullable"}},
				},
			},
		},
		Enums: []specification.Enum{
			{
				Name: testErrorCodeEnumName,
				Values: []specification.EnumValue{
					{Name: "BadRequest", Description: "Bad request"},
					{Name: "Internal", Description: "Internal error"},
				},
			},
		},
	}

	// Act
	var buf bytes.Buffer
	err := generateObjects(&buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating objects")
	
	generatedCode := buf.String()
	
	// Verify the complete Error() method implementation
	errorMethodStart := strings.Index(generatedCode, "func (e *Error) Error() string {")
	assert.NotEqual(t, -1, errorMethodStart, "Should find Error() method")
	
	errorMethodEnd := strings.Index(generatedCode[errorMethodStart:], "}")
	assert.NotEqual(t, -1, errorMethodEnd, "Should find end of Error() method")
	
	errorMethod := generatedCode[errorMethodStart : errorMethodStart+errorMethodEnd+1]
	
	// Verify method structure
	assert.Contains(t, errorMethod, "func (e *Error) Error() string",
		"Method should have correct signature")
	assert.Contains(t, errorMethod, "return e.Message.String()",
		"Method should return the message as a string")
	assert.True(t, strings.Count(errorMethod, "return") == 1,
		"Method should have exactly one return statement")

	t.Run("pointer receiver", func(t *testing.T) {
		// Verify that the method uses a pointer receiver
		assert.Contains(t, generatedCode, "func (e *Error)",
			"Error() method should use pointer receiver")
		assert.NotContains(t, generatedCode, "func (e Error)",
			"Error() method should not use value receiver")
	})
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestErrorInterfaceIntegration(t *testing.T) {
	// This test verifies that the Error type works correctly in the context
	// of the full generated server code
	
	// Arrange
	service := &specification.Service{
		Name:    "TestAPI",
		Version: "v1",
		Objects: []specification.Object{
			{
				Name:        testErrorStructName,
				Description: "API error response",
				Fields: []specification.Field{
					{Name: testErrorCodeField, Type: testErrorCodeEnumName},
					{Name: testErrorMessageField, Type: "String"},
					{Name: testErrorRequestIDField, Type: "String"},
				},
			},
		},
		Enums: []specification.Enum{
			{
				Name: testErrorCodeEnumName,
				Values: []specification.EnumValue{
					{Name: "BadRequest", Description: "Bad request"},
					{Name: "Unauthorized", Description: "Unauthorized"},
					{Name: "Internal", Description: "Internal error"},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name: "User",
				Endpoints: []specification.Endpoint{
					{
						Name:   "GetUser",
						Method: "GET",
						Path:   "/{id}",
						Request: specification.EndpointRequest{
							PathParams: []specification.Field{
								{Name: "ID", Type: "UUID"},
							},
						},
						Response: specification.EndpointResponse{
							StatusCode: 200,
						},
					},
				},
			},
		},
	}

	// Act
	var buf bytes.Buffer
	err := GenerateServer(&buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating full server")
	
	generatedCode := buf.String()
	
	// Verify Error is used in ConvertErrorFunc
	assert.Contains(t, generatedCode, "ConvertErrorFunc func(err error, requestID string) *Error",
		"ConvertErrorFunc should return *Error")
	
	// Verify Error is used in error handling
	assert.Contains(t, generatedCode, "return &Error{",
		"Error handling should create Error instances")
	
	// Verify the Error type can be used as an error
	assert.Contains(t, generatedCode, "apiError := server.ConvertErrorFunc(err, requestID)",
		"Should be able to assign Error to error variables")
	
	// Verify HTTPStatusCode is used
	assert.Contains(t, generatedCode, "apiError.HTTPStatusCode()",
		"Should use HTTPStatusCode() method for error responses")
}