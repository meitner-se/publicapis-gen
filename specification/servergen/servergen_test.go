package servergen

import (
	"bytes"
	"testing"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/stretchr/testify/assert"
)

// Test constants to avoid hardcoded strings
const (
	// Service constants
	testServiceName    = "TestService"
	testServiceVersion = "v1"
	testPathName       = "test-service"

	// Resource constants
	testResourceName       = "User"
	testResourceDesc       = "User resource"
	testResourcePluralName = "users"
	testResourcePath       = "/users"

	// Field constants
	testFieldName          = "Name"
	testFieldDesc          = "User name"
	testFieldType          = "String"
	testFieldTypeUUID      = "UUID"
	testFieldTypeInt       = "Int"
	testFieldTypeBool      = "Bool"
	testFieldTypeDate      = "Date"
	testFieldTypeTimestamp = "Timestamp"
	testFieldTypeCustom    = "CustomObject"

	// Enum constants
	testEnumName       = "UserRole"
	testEnumDesc       = "User role enumeration"
	testEnumValueAdmin = "Admin"
	testEnumValueUser  = "User"
	testEnumDescAdmin  = "Administrator role"
	testEnumDescUser   = "Regular user role"

	// Object constants
	testObjectName      = "Address"
	testObjectDesc      = "Address object"
	testObjectNameError = "Error"

	// Endpoint constants
	testEndpointName         = "CreateUser"
	testEndpointMethod       = "POST"
	testEndpointPath         = "/users"
	testEndpointTitle        = "Create User"
	testEndpointSummary      = "Create a new user"
	testEndpointDesc         = "Creates a new user in the system"
	testEndpointResponseCode = 201

	// Expected generated code fragments
	expectedPackageDecl   = "package api"
	expectedImportStmt    = "import ("
	expectedContextImport = `"context"`
	expectedEmbedImport   = `"embed"`
	expectedJSONImport    = `"encoding/json"`
	expectedHTTPImport    = `"net/http"`
	expectedGinImport     = `"github.com/gin-gonic/gin"`
	expectedTypesImport   = `"github.com/meitner-se/go-types"`

	// Error handling constants
	expectedErrorMethod        = "func (e *Error) Error() string"
	expectedHTTPStatusMethod   = "func (e *Error) HTTPStatusCode() int"
	expectedBadRequestCase     = "case ErrorCodeBadRequest:"
	expectedBadRequestReturn   = "return http.StatusBadRequest"
	expectedUnauthorizedCase   = "case ErrorCodeUnauthorized:"
	expectedUnauthorizedReturn = "return http.StatusUnauthorized"

	// Server registration constants
	expectedRegisterFunc      = "func RegisterTestServiceAPI[Session any]"
	expectedConvertErrorCheck = "if api.Server.ConvertErrorFunc == nil"
	expectedGetSessionCheck   = "if api.Server.GetSessionFunc == nil"
	expectedPanicGetSession   = `panic("GetSessionFunc is nil")`
	expectedRouterGroup       = `routerGroup := router.Group("/test-service/v1")`
	expectedOpenAPIRoute      = `routerGroup.StaticFileFS("/openapi.json", "openapi.json", http.FS(api.OpenAPI_JSON))`

	// Type generation constants
	expectedEnumType   = "type UserRole types.String"
	expectedEnumVar    = "var ("
	expectedEnumValue  = "UserRoleAdmin = UserRole(types.NewString(\"Admin\")) // Administrator role"
	expectedObjectType = "type Address struct {"
	expectedFieldDecl  = "Name types.String `json:\"name\"`"

	// Request/Response type constants
	expectedRequestType     = "type Request[sessionType, pathParamsType, queryParamsType, bodyParamsType any] struct {"
	expectedRequestIDMethod = "func (r Request[sessionType, pathParamsType, queryParamsType, bodyParamsType]) RequestID() string"
	expectedPathParamsType  = "type UserCreateUserPathParams struct {"
	expectedQueryParamsType = "type UserCreateUserQueryParams struct {"
	expectedBodyParamsType  = "type UserCreateUserBodyParams struct {"
	expectedResponseType    = "type UserCreateUserResponse struct {"

	// Utility function constants
	expectedServeWithResponse    = "func serveWithResponse["
	expectedServeWithoutResponse = "func serveWithoutResponse["
	expectedParseRequest         = "func parseRequest["
	expectedDecodeBodyParams     = "func decodeBodyParams[T any](r *http.Request) (T, error)"
	expectedDecodePathParams     = "func decodePathParams[T any](c *gin.Context) (T, error)"
	expectedDecodeQueryParams    = "func decodeQueryParams[T any](c *gin.Context) (T, error)"

	// Comment constants
	expectedFieldComment  = "// Name: User name"
	expectedObjectComment = "// Address object"
)

// ============================================================================
// GenerateServer Tests
// ============================================================================

func TestGenerateServer(t *testing.T) {
	// Arrange
	service := createTestService()
	buf := &bytes.Buffer{}

	// Act
	err := GenerateServer(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating server code")

	generatedCode := buf.String()
	assert.NotEmpty(t, generatedCode, "Expected generated code to be non-empty")

	// Verify package declaration
	assert.Contains(t, generatedCode, expectedPackageDecl, "Generated code should contain package declaration")

	// Verify imports
	assert.Contains(t, generatedCode, expectedImportStmt, "Generated code should contain import statement")
	assert.Contains(t, generatedCode, expectedContextImport, "Generated code should import context")
	assert.Contains(t, generatedCode, expectedEmbedImport, "Generated code should import embed")
	assert.Contains(t, generatedCode, expectedJSONImport, "Generated code should import encoding/json")
	assert.Contains(t, generatedCode, expectedHTTPImport, "Generated code should import net/http")
	assert.Contains(t, generatedCode, expectedGinImport, "Generated code should import gin-gonic/gin")
	assert.Contains(t, generatedCode, expectedTypesImport, "Generated code should import go-types")

	// Verify server registration function
	assert.Contains(t, generatedCode, expectedRegisterFunc, "Generated code should contain RegisterAPI function")

	// Verify enum generation
	assert.Contains(t, generatedCode, expectedEnumType, "Generated code should contain enum type declaration")

	// Verify object generation
	assert.Contains(t, generatedCode, expectedObjectType, "Generated code should contain object type declaration")

	// Verify request/response types
	assert.Contains(t, generatedCode, expectedRequestType, "Generated code should contain Request type")

	// Verify utility functions
	assert.Contains(t, generatedCode, expectedServeWithResponse, "Generated code should contain serveWithResponse function")
	assert.Contains(t, generatedCode, expectedServeWithoutResponse, "Generated code should contain serveWithoutResponse function")
	assert.Contains(t, generatedCode, expectedParseRequest, "Generated code should contain parseRequest function")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty service", func(t *testing.T) {
			// Arrange
			emptyService := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
			}
			buf := &bytes.Buffer{}

			// Act
			err := GenerateServer(buf, emptyService)

			// Assert
			assert.Nil(t, err, "Expected no error with empty service")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, expectedPackageDecl, "Should still generate basic structure")
		})

		t.Run("gofmt formatting", func(t *testing.T) {
			// Arrange
			service := createTestService()
			buf := &bytes.Buffer{}

			// Act
			err := GenerateServer(buf, service)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			// Check that code is properly formatted (proper indentation)
			// Note: Double spaces can appear in comments like "// //"
			assert.Contains(t, generatedCode, "\t", "Generated code should use tabs for indentation")
		})
	})
}

// ============================================================================
// generateEnums Tests
// ============================================================================

func TestGenerateEnums(t *testing.T) {
	// Arrange
	enums := []specification.Enum{
		{
			Name:        testEnumName,
			Description: testEnumDesc,
			Values: []specification.EnumValue{
				{Name: testEnumValueAdmin, Description: testEnumDescAdmin},
				{Name: testEnumValueUser, Description: testEnumDescUser},
			},
		},
	}
	buf := &bytes.Buffer{}

	// Act
	err := generateEnums(buf, enums)

	// Assert
	assert.Nil(t, err, "Expected no error when generating enums")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedEnumType, "Should generate enum type declaration")
	assert.Contains(t, generatedCode, expectedEnumVar, "Should generate var block")
	assert.Contains(t, generatedCode, expectedEnumValue, "Should generate enum value with description")
	assert.Contains(t, generatedCode, "UserRoleUser = UserRole(types.NewString(\"User\")) // Regular user role",
		"Should generate all enum values")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty enums slice", func(t *testing.T) {
			// Arrange
			buf := &bytes.Buffer{}

			// Act
			err := generateEnums(buf, []specification.Enum{})

			// Assert
			assert.Nil(t, err, "Expected no error with empty enums")
			assert.Empty(t, buf.String(), "Should generate nothing for empty enums")
		})

		t.Run("enum with no values", func(t *testing.T) {
			// Arrange
			enumsNoValues := []specification.Enum{
				{
					Name:        testEnumName,
					Description: testEnumDesc,
					Values:      []specification.EnumValue{},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateEnums(buf, enumsNoValues)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, expectedEnumType, "Should still generate enum type")
			assert.Contains(t, generatedCode, expectedEnumVar, "Should generate empty var block")
		})

		t.Run("enum value with special characters", func(t *testing.T) {
			// Arrange
			specialEnums := []specification.Enum{
				{
					Name:        "Status",
					Description: "Status enum",
					Values: []specification.EnumValue{
						{Name: "In-Progress", Description: "Work in progress"},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateEnums(buf, specialEnums)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, `StatusIn-Progress = Status(types.NewString("In-Progress"))`,
				"Should handle special characters in enum names")
		})
	})
}

// ============================================================================
// getTypeForGo Tests
// ============================================================================

func TestGetTypeForGo(t *testing.T) {
	// Arrange
	service := createTestService()

	testCases := []struct {
		name         string
		field        specification.Field
		expectedType string
	}{
		{
			name: "primitive string type",
			field: specification.Field{
				Name: testFieldName,
				Type: testFieldType,
			},
			expectedType: "types.String",
		},
		{
			name: "primitive UUID type",
			field: specification.Field{
				Name: "ID",
				Type: testFieldTypeUUID,
			},
			expectedType: "types.UUID",
		},
		{
			name: "custom object type",
			field: specification.Field{
				Name: "Address",
				Type: testObjectName,
			},
			expectedType: testObjectName,
		},
		{
			name: "nullable object type",
			field: specification.Field{
				Name:      "Address",
				Type:      testObjectName,
				Modifiers: []string{"Nullable"},
			},
			expectedType: "*" + testObjectName,
		},
		{
			name: "array of strings",
			field: specification.Field{
				Name:      "Tags",
				Type:      testFieldType,
				Modifiers: []string{"Array"},
			},
			expectedType: "[]types.String",
		},
		{
			name: "array of objects",
			field: specification.Field{
				Name:      "Addresses",
				Type:      testObjectName,
				Modifiers: []string{"Array"},
			},
			expectedType: "[]" + testObjectName,
		},
		{
			name: "nullable array of objects",
			field: specification.Field{
				Name:      "Addresses",
				Type:      testObjectName,
				Modifiers: []string{"Array", "Nullable"},
			},
			expectedType: "[]*" + testObjectName,
		},
		{
			name: "enum type",
			field: specification.Field{
				Name: "Role",
				Type: testEnumName,
			},
			expectedType: testEnumName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := getTypeForGo(tc.field, service)

			// Assert
			assert.Equal(t, tc.expectedType, result,
				"Expected Go type to be %s for field %s", tc.expectedType, tc.field.Name)
		})
	}

	t.Run("edge cases", func(t *testing.T) {
		t.Run("unknown type", func(t *testing.T) {
			// Arrange
			field := specification.Field{
				Name: "Unknown",
				Type: "UnknownType",
			}

			// Act
			result := getTypeForGo(field, service)

			// Assert
			assert.Equal(t, "types.UnknownType", result,
				"Unknown types should be prefixed with types.")
		})

		t.Run("all modifiers combined", func(t *testing.T) {
			// Arrange
			field := specification.Field{
				Name:      "ComplexField",
				Type:      testObjectName,
				Modifiers: []string{"Array", "Nullable"},
			}

			// Act
			result := getTypeForGo(field, service)

			// Assert
			assert.Equal(t, "[]*"+testObjectName, result,
				"Should handle both array and nullable modifiers correctly")
		})
	})
}

// ============================================================================
// getTypePrefix Tests
// ============================================================================

func TestGetTypePrefix(t *testing.T) {
	// Arrange
	service := createTestService()

	testCases := []struct {
		name           string
		field          specification.Field
		expectedPrefix string
	}{
		{
			name: "primitive type no modifiers",
			field: specification.Field{
				Type: testFieldType,
			},
			expectedPrefix: "types.",
		},
		{
			name: "object type no modifiers",
			field: specification.Field{
				Type: testObjectName,
			},
			expectedPrefix: "",
		},
		{
			name: "array primitive type",
			field: specification.Field{
				Type:      testFieldType,
				Modifiers: []string{"Array"},
			},
			expectedPrefix: "[]types.",
		},
		{
			name: "nullable object",
			field: specification.Field{
				Type:      testObjectName,
				Modifiers: []string{"Nullable"},
			},
			expectedPrefix: "*",
		},
		{
			name: "array nullable object",
			field: specification.Field{
				Type:      testObjectName,
				Modifiers: []string{"Array", "Nullable"},
			},
			expectedPrefix: "[]*",
		},
		{
			name: "enum type",
			field: specification.Field{
				Type: testEnumName,
			},
			expectedPrefix: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := getTypePrefix(tc.field, service)

			// Assert
			assert.Equal(t, tc.expectedPrefix, result,
				"Expected type prefix to be %s", tc.expectedPrefix)
		})
	}

	t.Run("edge cases", func(t *testing.T) {
		t.Run("nullable primitive", func(t *testing.T) {
			// Arrange
			field := specification.Field{
				Type:      testFieldType,
				Modifiers: []string{"Nullable"},
			}

			// Act
			result := getTypePrefix(field, service)

			// Assert
			assert.Equal(t, "types.", result,
				"Nullable primitives should not have pointer prefix")
		})

		t.Run("unknown type treated as primitive", func(t *testing.T) {
			// Arrange
			field := specification.Field{
				Type: "UnknownType",
			}

			// Act
			result := getTypePrefix(field, service)

			// Assert
			assert.Equal(t, "types.", result,
				"Unknown types should be treated as primitives")
		})
	})
}

// ============================================================================
// generateObjects Tests
// ============================================================================

func TestGenerateObjects(t *testing.T) {
	// Arrange
	service := createTestService()
	buf := &bytes.Buffer{}

	// Act
	err := generateObjects(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating objects")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedObjectComment, "Should generate object comment")
	assert.Contains(t, generatedCode, expectedObjectType, "Should generate object struct")
	assert.Contains(t, generatedCode, expectedFieldComment, "Should generate field comment")
	assert.Contains(t, generatedCode, expectedFieldDecl, "Should generate field declaration")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty objects", func(t *testing.T) {
			// Arrange
			serviceNoObjects := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Objects: []specification.Object{},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateObjects(buf, serviceNoObjects)

			// Assert
			assert.Nil(t, err, "Expected no error with empty objects")
			assert.Empty(t, buf.String(), "Should generate nothing for empty objects")
		})

		t.Run("object with no fields", func(t *testing.T) {
			// Arrange
			serviceEmptyObject := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Objects: []specification.Object{
					{
						Name:        "Empty",
						Description: "Empty object",
						Fields:      []specification.Field{},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateObjects(buf, serviceEmptyObject)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "type Empty struct {", "Should generate empty struct")
			assert.Contains(t, generatedCode, "}", "Should close struct properly")
		})

		t.Run("Error object with all status codes", func(t *testing.T) {
			// Arrange
			serviceWithError := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Objects: []specification.Object{
					{
						Name:        testObjectNameError,
						Description: "Error object",
						Fields: []specification.Field{
							{Name: "Code", Type: "ErrorCode"},
							{Name: "Message", Type: testFieldType},
						},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateObjects(buf, serviceWithError)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, expectedUnauthorizedCase, "Should handle Unauthorized case")
			assert.Contains(t, generatedCode, expectedUnauthorizedReturn, "Should return 401 for Unauthorized")
			assert.Contains(t, generatedCode, "default:", "Should have default case")
			assert.Contains(t, generatedCode, "return http.StatusInternalServerError", "Default should return 500")
		})
	})
}

// ============================================================================
// generateServer Tests
// ============================================================================

func TestGenerateServerFunc(t *testing.T) {
	// Arrange
	service := createTestServiceWithEndpoints()
	buf := &bytes.Buffer{}

	// Act
	err := generateServer(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating server function")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedRegisterFunc, "Should generate RegisterAPI function")
	assert.Contains(t, generatedCode, expectedConvertErrorCheck, "Should check ConvertErrorFunc")
	assert.Contains(t, generatedCode, expectedGetSessionCheck, "Should check GetSessionFunc")
	assert.Contains(t, generatedCode, expectedPanicGetSession, "Should panic if GetSessionFunc is nil")
	assert.Contains(t, generatedCode, expectedRouterGroup, "Should create router group with correct path")
	assert.Contains(t, generatedCode, expectedOpenAPIRoute, "Should register OpenAPI route")

	// Check endpoint registration (note: generates singular paths)
	assert.Contains(t, generatedCode, `routerGroup.POST("/user", serveWithResponse(201, api.Server, api.User.CreateUser))`,
		"Should register POST endpoint with response")
	assert.Contains(t, generatedCode, `routerGroup.DELETE("/user/{id}", serveWithoutResponse(204, api.Server, api.User.DeleteUser))`,
		"Should register DELETE endpoint without response")

	// Check type definitions
	assert.Contains(t, generatedCode, "type getSessionFunc[T any] func(ctx context.Context, headers http.Header, requestID string) (T, error)",
		"Should define getSessionFunc type")
	assert.Contains(t, generatedCode, "type TestServiceAPI[Session any] struct {",
		"Should define service API struct")
	assert.Contains(t, generatedCode, "type UserAPI[Session any] interface {",
		"Should define resource API interface")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("service with no resources", func(t *testing.T) {
			// Arrange
			serviceNoResources := &specification.Service{
				Name:      testServiceName,
				Version:   testServiceVersion,
				Resources: []specification.Resource{},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateServer(buf, serviceNoResources)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, expectedRegisterFunc, "Should still generate RegisterAPI function")
			assert.NotContains(t, generatedCode, "routerGroup.POST", "Should not register any endpoints")
		})

		t.Run("endpoint with different HTTP methods", func(t *testing.T) {
			// Arrange
			serviceVariousMethods := createTestServiceWithVariousHTTPMethods()
			buf := &bytes.Buffer{}

			// Act
			err := generateServer(buf, serviceVariousMethods)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "routerGroup.GET(", "Should register GET endpoint")
			assert.Contains(t, generatedCode, "routerGroup.POST(", "Should register POST endpoint")
			assert.Contains(t, generatedCode, "routerGroup.PATCH(", "Should register PATCH endpoint")
			assert.Contains(t, generatedCode, "routerGroup.DELETE(", "Should register DELETE endpoint")
		})
	})
}

// ============================================================================
// generateRequestTypes Tests
// ============================================================================

func TestGenerateRequestTypes(t *testing.T) {
	// Arrange
	service := createTestServiceWithEndpoints()
	buf := &bytes.Buffer{}

	// Act
	err := generateRequestTypes(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating request types")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedRequestType, "Should generate Request generic type")
	assert.Contains(t, generatedCode, expectedRequestIDMethod, "Should generate RequestID() method")
	assert.Contains(t, generatedCode, "return r.requestID", "RequestID method should return requestID field")

	// Check path params type generation
	assert.Contains(t, generatedCode, "type UserDeleteUserPathParams struct {", "Should generate path params type")
	assert.Contains(t, generatedCode, `ID types.UUID `+"`json:\"id\"`", "Should include path param fields")

	// Check query params type generation
	assert.Contains(t, generatedCode, "type UserListUsersQueryParams struct {", "Should generate query params type")
	assert.Contains(t, generatedCode, `Limit types.Int `+"`form:\"limit\"`", "Should use form tag for query params")

	// Check body params type generation
	assert.Contains(t, generatedCode, "type UserCreateUserBodyParams struct {", "Should generate body params type")
	assert.Contains(t, generatedCode, `Name types.String `+"`json:\"name\"`", "Should use json tag for body params")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("endpoint with no params", func(t *testing.T) {
			// Arrange
			serviceNoParams := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Resources: []specification.Resource{
					{
						Name: testResourceName,
						Endpoints: []specification.Endpoint{
							{
								Name:   "Simple",
								Method: "GET",
								Path:   "/simple",
								Request: specification.EndpointRequest{
									PathParams:  []specification.Field{},
									QueryParams: []specification.Field{},
									BodyParams:  []specification.Field{},
								},
							},
						},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateRequestTypes(buf, serviceNoParams)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.NotContains(t, generatedCode, "type UserSimplePathParams", "Should not generate empty param types")
			assert.NotContains(t, generatedCode, "type UserSimpleQueryParams", "Should not generate empty param types")
			assert.NotContains(t, generatedCode, "type UserSimpleBodyParams", "Should not generate empty param types")
		})

		t.Run("fields with custom objects", func(t *testing.T) {
			// Arrange
			serviceCustomFields := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Objects: []specification.Object{
					{Name: testObjectName},
				},
				Resources: []specification.Resource{
					{
						Name: testResourceName,
						Endpoints: []specification.Endpoint{
							{
								Name:   "CreateWithAddress",
								Method: "POST",
								Request: specification.EndpointRequest{
									BodyParams: []specification.Field{
										{
											Name: "HomeAddress",
											Type: testObjectName,
										},
									},
								},
							},
						},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateRequestTypes(buf, serviceCustomFields)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "HomeAddress Address `json:\"homeAddress\"`",
				"Should handle custom object types correctly")
		})
	})
}

// ============================================================================
// generateResponseTypes Tests
// ============================================================================

func TestGenerateResponseTypes(t *testing.T) {
	// Arrange
	service := createTestServiceWithEndpoints()
	buf := &bytes.Buffer{}

	// Act
	err := generateResponseTypes(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating response types")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, "type UserListUsersResponse struct {", "Should generate response type")
	assert.Contains(t, generatedCode, `Data []types.User `+"`json:\"data\"`", "Should include array field")
	assert.Contains(t, generatedCode, `Pagination Pagination `+"`json:\"pagination\"`", "Should include pagination field")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("endpoint with no response body", func(t *testing.T) {
			// Arrange
			serviceNoResponse := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Resources: []specification.Resource{
					{
						Name: testResourceName,
						Endpoints: []specification.Endpoint{
							{
								Name:   "Delete",
								Method: "DELETE",
								Response: specification.EndpointResponse{
									StatusCode: 204,
									BodyFields: []specification.Field{},
								},
							},
						},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateResponseTypes(buf, serviceNoResponse)

			// Assert
			assert.Nil(t, err, "Expected no error")
			assert.NotContains(t, buf.String(), "type UserDeleteResponse",
				"Should not generate response type for endpoints with no body")
		})

		t.Run("response with custom object fields", func(t *testing.T) {
			// Arrange
			serviceCustomResponse := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Objects: []specification.Object{
					{Name: "Stats"},
				},
				Resources: []specification.Resource{
					{
						Name: testResourceName,
						Endpoints: []specification.Endpoint{
							{
								Name: "GetStats",
								Response: specification.EndpointResponse{
									StatusCode: 200,
									BodyFields: []specification.Field{
										{
											Name: "UserStats",
											Type: "Stats",
										},
									},
								},
							},
						},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateResponseTypes(buf, serviceCustomResponse)

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "UserStats Stats `json:\"userStats\"`",
				"Should handle custom object types in responses")
		})
	})
}

// ============================================================================
// generateUtils Tests
// ============================================================================

func TestGenerateUtils(t *testing.T) {
	// Arrange
	buf := &bytes.Buffer{}

	// Act
	err := generateUtils(buf)

	// Assert
	assert.Nil(t, err, "Expected no error when generating utils")

	generatedCode := buf.String()

	// Check all utility functions are generated
	assert.Contains(t, generatedCode, expectedServeWithResponse, "Should generate serveWithResponse")
	assert.Contains(t, generatedCode, expectedServeWithoutResponse, "Should generate serveWithoutResponse")
	assert.Contains(t, generatedCode, expectedParseRequest, "Should generate parseRequest")
	assert.Contains(t, generatedCode, expectedDecodeBodyParams, "Should generate decodeBodyParams")
	assert.Contains(t, generatedCode, expectedDecodePathParams, "Should generate decodePathParams")
	assert.Contains(t, generatedCode, expectedDecodeQueryParams, "Should generate decodeQueryParams")

	// Check function implementations
	assert.Contains(t, generatedCode, `requestID := "test-request-id"`, "Should use test request ID")
	assert.Contains(t, generatedCode, "parseRequest[sessionType, pathParamsType, queryParamsType, bodyParamsType]",
		"Should call parseRequest with generic types")
	assert.Contains(t, generatedCode, "c.JSON(successStatusCode, response)",
		"Should return JSON response with success code")
	assert.Contains(t, generatedCode, "c.JSON(apiError.HTTPStatusCode(), apiError)",
		"Should return error with appropriate status code")

	// Check parseRequest implementation
	assert.Contains(t, generatedCode, "if _, ok := any(request.BodyParams).(struct{}); !ok {",
		"Should check if body params exist")
	assert.Contains(t, generatedCode, "if _, ok := any(request.PathParams).(struct{}); !ok {",
		"Should check if path params exist")
	assert.Contains(t, generatedCode, "if _, ok := any(request.QueryParams).(struct{}); !ok {",
		"Should check if query params exist")

	// Check decode functions
	assert.Contains(t, generatedCode, "json.NewDecoder(r.Body).Decode(&v)",
		"Should decode JSON body")
	assert.Contains(t, generatedCode, "c.ShouldBindQuery(&result)",
		"Should bind query parameters")
	assert.Contains(t, generatedCode, "m := make(map[string]string, len(c.Params))",
		"Should create map for path params")

	t.Run("consistency", func(t *testing.T) {
		// Generate twice and compare
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		err1 := generateUtils(buf1)
		err2 := generateUtils(buf2)

		assert.Nil(t, err1, "First generation should not error")
		assert.Nil(t, err2, "Second generation should not error")
		assert.Equal(t, buf1.String(), buf2.String(), "Utils should be generated consistently")
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

func createTestService() *specification.Service {
	return &specification.Service{
		Name:    testServiceName,
		Version: testServiceVersion,
		Enums: []specification.Enum{
			{
				Name:        testEnumName,
				Description: testEnumDesc,
				Values: []specification.EnumValue{
					{Name: testEnumValueAdmin, Description: testEnumDescAdmin},
					{Name: testEnumValueUser, Description: testEnumDescUser},
				},
			},
		},
		Objects: []specification.Object{
			{
				Name:        testObjectName,
				Description: testObjectDesc,
				Fields: []specification.Field{
					{
						Name:        testFieldName,
						Description: testFieldDesc,
						Type:        testFieldType,
					},
				},
			},
		},
		Resources: []specification.Resource{
			{
				Name:        testResourceName,
				Description: testResourceDesc,
			},
		},
	}
}

func createTestServiceWithEndpoints() *specification.Service {
	service := createTestService()
	service.Resources[0].Endpoints = []specification.Endpoint{
		{
			Name:        "CreateUser",
			Method:      "POST",
			Path:        "",
			Title:       testEndpointTitle,
			Summary:     testEndpointSummary,
			Description: testEndpointDesc,
			Request: specification.EndpointRequest{
				ContentType: "application/json",
				BodyParams: []specification.Field{
					{
						Name:        testFieldName,
						Description: testFieldDesc,
						Type:        testFieldType,
					},
				},
			},
			Response: specification.EndpointResponse{
				ContentType: "application/json",
				StatusCode:  testEndpointResponseCode,
				BodyObject:  &service.Resources[0].Name,
			},
		},
		{
			Name:   "DeleteUser",
			Method: "DELETE",
			Path:   "/{id}",
			Request: specification.EndpointRequest{
				PathParams: []specification.Field{
					{
						Name: "ID",
						Type: testFieldTypeUUID,
					},
				},
			},
			Response: specification.EndpointResponse{
				StatusCode: 204,
			},
		},
		{
			Name:   "ListUsers",
			Method: "GET",
			Path:   "",
			Request: specification.EndpointRequest{
				QueryParams: []specification.Field{
					{
						Name:    "Limit",
						Type:    testFieldTypeInt,
						Default: "50",
					},
				},
			},
			Response: specification.EndpointResponse{
				StatusCode: 200,
				BodyFields: []specification.Field{
					{
						Name:      "Data",
						Type:      testResourceName,
						Modifiers: []string{"Array"},
					},
					{
						Name: "Pagination",
						Type: "Pagination",
					},
				},
			},
		},
	}

	// Add Pagination object for list response
	service.Objects = append(service.Objects, specification.Object{
		Name: "Pagination",
		Fields: []specification.Field{
			{Name: "Offset", Type: testFieldTypeInt},
			{Name: "Limit", Type: testFieldTypeInt},
			{Name: "Total", Type: testFieldTypeInt},
		},
	})

	return service
}

func createTestServiceWithVariousHTTPMethods() *specification.Service {
	return &specification.Service{
		Name:    testServiceName,
		Version: testServiceVersion,
		Resources: []specification.Resource{
			{
				Name: testResourceName,
				Endpoints: []specification.Endpoint{
					{
						Name:   "Get",
						Method: "GET",
						Path:   "/{id}",
						Response: specification.EndpointResponse{
							StatusCode: 200,
						},
					},
					{
						Name:   "Create",
						Method: "POST",
						Path:   "",
						Response: specification.EndpointResponse{
							StatusCode: 201,
						},
					},
					{
						Name:   "Update",
						Method: "PATCH",
						Path:   "/{id}",
						Response: specification.EndpointResponse{
							StatusCode: 200,
						},
					},
					{
						Name:   "Delete",
						Method: "DELETE",
						Path:   "/{id}",
						Response: specification.EndpointResponse{
							StatusCode: 204,
						},
					},
				},
			},
		},
	}
}
