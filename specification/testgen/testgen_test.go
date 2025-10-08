package testgen

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
	testResourceName = "Student"
	testResourceDesc = "Student resource"

	// Field constants
	testFieldName     = "Name"
	testFieldDesc     = "Student name"
	testFieldType     = "String"
	testFieldTypeUUID = "UUID"
	testFieldTypeInt  = "Int"
	testFieldTypeBool = "Bool"

	// Endpoint constants
	testEndpointName         = "CreateStudent"
	testEndpointMethod       = "POST"
	testEndpointPath         = ""
	testEndpointTitle        = "Create Student"
	testEndpointSummary      = "Create a new student"
	testEndpointDesc         = "Creates a new student in the system"
	testEndpointResponseCode = 201

	// Expected generated code fragments
	expectedPackageDecl    = "package main"
	expectedImportStmt     = "import ("
	expectedContextImport  = `"context"`
	expectedIOImport       = "\"io\""
	expectedHTTPImport     = "\"net/http\""
	expectedHTTPTestImport = "\"net/http/httptest\""
	expectedURLImport      = "\"net/url\""
	expectedTestingImport  = "\"testing\""
	expectedGinImport      = "\"github.com/gin-gonic/gin\""
	expectedUUIDImport     = "\"github.com/google/uuid\""
	expectedTypesImport    = "\"github.com/meitner-se/go-types\""
	expectedTestifyImport  = "\"github.com/stretchr/testify/assert\""

	// Test function constants
	expectedTestFunction               = "func TestStudentCreateStudent(t *testing.T) {"
	expectedTestConstants              = "const ("
	expectedTestRequestID              = "testRequestID = \"test-request-id-123\""
	expectedTestSessionUserID          = "testSessionUserID = \"test-session-user-id\""
	expectedTestRequestIDFormatted     = "testRequestID     = \"test-request-id-123\""
	expectedTestSessionUserIDFormatted = "testSessionUserID = \"test-session-user-id\""
	expectedGinTestMode                = "gin.SetMode(gin.TestMode)"

	// Mock constants
	expectedMockInterface  = "type MockStudentAPI struct {"
	expectedMockMethod     = "func (m *MockStudentAPI) CreateStudent"
	expectedFuncField      = "CreateStudentFunc func(ctx context.Context"
	expectedRequestCapture = "capturedRequest = request"

	// HTTP request constants
	expectedHTTPRequest       = "http.NewRequestWithContext(ctx, \"POST\""
	expectedHTTPStatusAssert  = "if resp.StatusCode != 201"
	expectedServerSetup       = "server := httptest.NewServer(router)"
	expectedURLConstruction   = "requestURL := server.URL + \"/"
	expectedResponseBodyCheck = "var responseBody map[string]interface{}"
)

// ============================================================================
// GenerateInternalTests Tests
// ============================================================================

func TestGenerateInternalTests(t *testing.T) {
	// Arrange
	service := createTestService()
	buf := &bytes.Buffer{}
	packageName := "myapi"

	// Act
	err := GenerateInternalTests(buf, service, packageName)

	// Assert
	assert.Nil(t, err, "Expected no error when generating internal tests")

	generatedCode := buf.String()

	// Check package declaration
	assert.Contains(t, generatedCode, "package myapi", "Should use specified package name")

	// Check imports (no external API import)
	assert.Contains(t, generatedCode, "\"github.com/gin-gonic/gin\"", "Should import gin")
	assert.Contains(t, generatedCode, "\"github.com/stretchr/testify/assert\"", "Should import testify")
	assert.NotContains(t, generatedCode, "\"./api\"", "Should not import external API package")

	// Check endpoint tests (no package prefixes)
	assert.Contains(t, generatedCode, "func TestStudentCreateStudent(t *testing.T) {", "Should generate endpoint test")
	assert.Contains(t, generatedCode, "var capturedRequest Request[any,", "Should use Request type without prefix")
	assert.Contains(t, generatedCode, "RegisterTestServiceAPI(router, &TestServiceAPI[any]{", "Should register API without prefix")

	// Check mock structures (no package prefixes)
	assert.Contains(t, generatedCode, "type MockStudentAPI struct {", "Should generate mock struct")
	assert.Contains(t, generatedCode, "func (m *MockStudentAPI) CreateStudent(", "Should generate mock method")

	// Check utility function tests (no package prefixes)
	assert.Contains(t, generatedCode, "func Test_serveWithResponse(t *testing.T) {", "Should generate serveWithResponse test")
	assert.Contains(t, generatedCode, "func Test_serveWithoutResponse(t *testing.T) {", "Should generate serveWithoutResponse test")
	assert.Contains(t, generatedCode, "func Test_handleRequest(t *testing.T) {", "Should generate handleRequest test")
	assert.Contains(t, generatedCode, "func Test_decodeBodyParams(t *testing.T) {", "Should generate decodeBodyParams test")
	assert.Contains(t, generatedCode, "func Test_decodePathParams(t *testing.T) {", "Should generate decodePathParams test")
	assert.Contains(t, generatedCode, "func Test_decodeQueryParams(t *testing.T) {", "Should generate decodeQueryParams test")

	// Verify no package prefixes are used
	assert.Contains(t, generatedCode, "handler := serveWithResponse(", "Should call serveWithResponse without prefix")
	assert.Contains(t, generatedCode, "handleRequest[any, struct{}, struct{}, struct{}](", "Should call handleRequest")
	assert.Contains(t, generatedCode, "decodeBodyParams[TestBody](", "Should call decodeBodyParams with struct type")
	assert.Contains(t, generatedCode, "decodePathParams[TestPathParams](", "Should call decodePathParams with struct type")
	assert.Contains(t, generatedCode, "decodeQueryParams[TestQueryParams](", "Should call decodeQueryParams with struct type")
}

// ============================================================================
// GenerateTests Tests
// ============================================================================

func TestGenerateTests(t *testing.T) {
	// Arrange
	service := createTestService()
	packageName := "main"
	buf := &bytes.Buffer{}

	// Act
	err := GenerateTests(buf, service, packageName, "api", "./api")

	// Assert
	assert.Nil(t, err, "Expected no error when generating tests")

	generatedCode := buf.String()
	assert.NotEmpty(t, generatedCode, "Expected generated code to be non-empty")

	// Verify package declaration
	assert.Contains(t, generatedCode, expectedPackageDecl, "Generated code should contain package declaration")

	// Verify imports
	assert.Contains(t, generatedCode, expectedImportStmt, "Generated code should contain import statement")
	assert.Contains(t, generatedCode, expectedContextImport, "Generated code should import context")
	assert.Contains(t, generatedCode, expectedIOImport, "Generated code should import io")
	assert.Contains(t, generatedCode, expectedHTTPImport, "Generated code should import net/http")
	assert.Contains(t, generatedCode, expectedHTTPTestImport, "Generated code should import net/http/httptest")
	assert.Contains(t, generatedCode, expectedURLImport, "Generated code should import net/url")
	assert.Contains(t, generatedCode, expectedTestingImport, "Generated code should import testing")
	assert.Contains(t, generatedCode, expectedGinImport, "Generated code should import gin")
	assert.Contains(t, generatedCode, expectedUUIDImport, "Generated code should import uuid")
	assert.Contains(t, generatedCode, expectedTypesImport, "Generated code should import go-types")
	assert.Contains(t, generatedCode, expectedTestifyImport, "Generated code should import testify assert")

	// Verify test constants
	assert.Contains(t, generatedCode, expectedTestConstants, "Generated code should contain test constants")
	assert.Contains(t, generatedCode, expectedTestRequestIDFormatted, "Generated code should contain test request ID")
	assert.Contains(t, generatedCode, expectedTestSessionUserIDFormatted, "Generated code should contain test session user ID")

	// Verify test function generation
	assert.Contains(t, generatedCode, expectedTestFunction, "Generated code should contain test function")
	assert.Contains(t, generatedCode, expectedGinTestMode, "Generated code should set gin test mode")

	// Verify mock generation
	assert.Contains(t, generatedCode, expectedMockInterface, "Generated code should contain mock interface")
	assert.Contains(t, generatedCode, expectedMockMethod, "Generated code should contain mock method")
	assert.Contains(t, generatedCode, expectedFuncField, "Generated code should contain function field")
	assert.Contains(t, generatedCode, expectedRequestCapture, "Generated code should capture request")

	// Verify HTTP request generation
	assert.Contains(t, generatedCode, expectedHTTPRequest, "Generated code should contain HTTP request creation")
	assert.Contains(t, generatedCode, expectedHTTPStatusAssert, "Generated code should contain HTTP status assertion")
	assert.Contains(t, generatedCode, expectedServerSetup, "Generated code should contain server setup")
	assert.Contains(t, generatedCode, expectedURLConstruction, "Generated code should contain URL construction")
	assert.Contains(t, generatedCode, expectedResponseBodyCheck, "Generated code should contain response body check")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty service", func(t *testing.T) {
			// Arrange
			emptyService := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
			}
			buf := &bytes.Buffer{}

			// Act
			err := GenerateTests(buf, emptyService, packageName, "api", "./api")

			// Assert
			assert.Nil(t, err, "Expected no error with empty service")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, expectedPackageDecl, "Should still generate basic structure")
		})

		t.Run("custom package name", func(t *testing.T) {
			// Arrange
			service := createTestService()
			customPackage := "customapi"
			buf := &bytes.Buffer{}

			// Act
			err := GenerateTests(buf, service, customPackage, "api", "./api")

			// Assert
			assert.Nil(t, err, "Expected no error with custom package")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "package "+customPackage, "Should use custom package name")
		})

		t.Run("gofmt formatting", func(t *testing.T) {
			// Arrange
			service := createTestService()
			buf := &bytes.Buffer{}

			// Act
			err := GenerateTests(buf, service, packageName, "api", "./api")

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			// Check that code is properly formatted (proper indentation)
			assert.Contains(t, generatedCode, "\t", "Generated code should use tabs for indentation")
		})
	})
}

// ============================================================================
// generateImports Tests
// ============================================================================

func TestGenerateImports(t *testing.T) {
	// Arrange
	buf := &bytes.Buffer{}

	// Act
	err := generateImports(buf, "api", "./api")

	// Assert
	assert.Nil(t, err, "Expected no error when generating imports")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedImportStmt, "Should generate import statement")
	assert.Contains(t, generatedCode, expectedContextImport, "Should import context")
	assert.Contains(t, generatedCode, expectedIOImport, "Should import io")
	assert.Contains(t, generatedCode, expectedHTTPImport, "Should import net/http")
	assert.Contains(t, generatedCode, expectedHTTPTestImport, "Should import net/http/httptest")
	assert.Contains(t, generatedCode, expectedURLImport, "Should import net/url")
	assert.Contains(t, generatedCode, expectedTestingImport, "Should import testing")
	assert.Contains(t, generatedCode, expectedGinImport, "Should import gin")
	assert.Contains(t, generatedCode, expectedUUIDImport, "Should import uuid")
	assert.Contains(t, generatedCode, expectedTypesImport, "Should import go-types")
	assert.Contains(t, generatedCode, expectedTestifyImport, "Should import testify assert")
}

// ============================================================================
// generateTestConstants Tests
// ============================================================================

func TestGenerateTestConstants(t *testing.T) {
	// Arrange
	service := createTestService()
	buf := &bytes.Buffer{}

	// Act
	err := generateTestConstants(buf, service)

	// Assert
	assert.Nil(t, err, "Expected no error when generating test constants")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedTestConstants, "Should generate test constants")
	assert.Contains(t, generatedCode, expectedTestRequestID, "Should generate test request ID")
	assert.Contains(t, generatedCode, expectedTestSessionUserID, "Should generate test session user ID")
	assert.Contains(t, generatedCode, "testTimeout", "Should generate test timeout")
}

// ============================================================================
// generateEndpointTest Tests
// ============================================================================

func TestGenerateEndpointTest(t *testing.T) {
	// Arrange
	service := createTestService()
	resource := service.Resources[0]
	endpoint := resource.Endpoints[0]
	buf := &bytes.Buffer{}

	// Act
	err := generateEndpointTest(buf, service, resource, endpoint, "api")

	// Assert
	assert.Nil(t, err, "Expected no error when generating endpoint test")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedTestFunction, "Should generate test function")
	assert.Contains(t, generatedCode, expectedGinTestMode, "Should set gin test mode")
	assert.Contains(t, generatedCode, "// Arrange", "Should contain test sections")
	assert.Contains(t, generatedCode, "// Act", "Should contain act section")
	assert.Contains(t, generatedCode, "// Assert", "Should contain assert section")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("endpoint with path parameters", func(t *testing.T) {
			// Arrange
			service := createTestServiceWithPathParams()
			resource := service.Resources[0]
			endpoint := resource.Endpoints[0]
			buf := &bytes.Buffer{}

			// Act
			err := generateEndpointTest(buf, service, resource, endpoint, "api")

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "// Path parameters", "Should generate path parameter section")
			assert.Contains(t, generatedCode, "testPathID", "Should generate path parameter variable")
		})

		t.Run("endpoint with query parameters", func(t *testing.T) {
			// Arrange
			service := createTestServiceWithQueryParams()
			resource := service.Resources[0]
			endpoint := resource.Endpoints[0]
			buf := &bytes.Buffer{}

			// Act
			err := generateEndpointTest(buf, service, resource, endpoint, "api")

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "// Query parameters", "Should generate query parameter section")
			assert.Contains(t, generatedCode, "testQueryLimit", "Should generate query parameter variable")
		})

		t.Run("endpoint with body parameters", func(t *testing.T) {
			// Arrange
			service := createTestService()
			resource := service.Resources[0]
			endpoint := resource.Endpoints[0]
			buf := &bytes.Buffer{}

			// Act
			err := generateEndpointTest(buf, service, resource, endpoint, "api")

			// Assert
			assert.Nil(t, err, "Expected no error")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "// Body parameters", "Should generate body parameter section")
			assert.Contains(t, generatedCode, "testBody := map[string]interface{}", "Should generate body map")
		})
	})
}

// ============================================================================
// generateHelperFunctions Tests
// ============================================================================

func TestGenerateHelperFunctions(t *testing.T) {
	// Arrange
	service := createTestService()
	buf := &bytes.Buffer{}

	// Act
	err := generateHelperFunctions(buf, service, "api")

	// Assert
	assert.Nil(t, err, "Expected no error when generating helper functions")

	generatedCode := buf.String()
	assert.Contains(t, generatedCode, expectedMockInterface, "Should generate mock interface")
	assert.Contains(t, generatedCode, expectedMockMethod, "Should generate mock method")
	assert.Contains(t, generatedCode, expectedFuncField, "Should generate function field")

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
			err := generateHelperFunctions(buf, serviceNoResources, "api")

			// Assert
			assert.Nil(t, err, "Expected no error with no resources")
			generatedCode := buf.String()
			assert.Contains(t, generatedCode, "Mock interfaces", "Should still generate header comment")
		})

		t.Run("resource with no endpoints", func(t *testing.T) {
			// Arrange
			serviceNoEndpoints := &specification.Service{
				Name:    testServiceName,
				Version: testServiceVersion,
				Resources: []specification.Resource{
					{
						Name:        testResourceName,
						Description: testResourceDesc,
						Endpoints:   []specification.Endpoint{},
					},
				},
			}
			buf := &bytes.Buffer{}

			// Act
			err := generateHelperFunctions(buf, serviceNoEndpoints, "api")

			// Assert
			assert.Nil(t, err, "Expected no error with no endpoints")
			generatedCode := buf.String()
			assert.NotContains(t, generatedCode, expectedMockInterface, "Should not generate mock for resource with no endpoints")
		})
	})
}

// ============================================================================
// getJSONKey Tests
// ============================================================================

func TestGetJSONKey(t *testing.T) {
	testCases := []struct {
		name        string
		fieldName   string
		expectedKey string
	}{
		{
			name:        "simple field name",
			fieldName:   "Name",
			expectedKey: "name",
		},
		{
			name:        "camelCase field name",
			fieldName:   "FirstName",
			expectedKey: "firstName",
		},
		{
			name:        "ID field",
			fieldName:   "ID",
			expectedKey: "id",
		},
		{
			name:        "complex field name",
			fieldName:   "CreatedAt",
			expectedKey: "createdAt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := getJSONKey(tc.fieldName)

			// Assert
			assert.Equal(t, tc.expectedKey, result,
				"Expected JSON key to be %s for field %s", tc.expectedKey, tc.fieldName)
		})
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func createTestService() *specification.Service {
	return &specification.Service{
		Name:    testServiceName,
		Version: testServiceVersion,
		Resources: []specification.Resource{
			{
				Name:        testResourceName,
				Description: testResourceDesc,
				Endpoints: []specification.Endpoint{
					{
						Name:        testEndpointName,
						Method:      testEndpointMethod,
						Path:        testEndpointPath,
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
							BodyObject:  getResourceNamePtr(),
						},
					},
				},
			},
		},
	}
}

func createTestServiceWithPathParams() *specification.Service {
	service := createTestService()
	service.Resources[0].Endpoints[0].Request.PathParams = []specification.Field{
		{
			Name:        "ID",
			Description: "Student ID",
			Type:        testFieldTypeUUID,
		},
	}
	service.Resources[0].Endpoints[0].Path = "/{id}"
	return service
}

func createTestServiceWithQueryParams() *specification.Service {
	service := createTestService()
	service.Resources[0].Endpoints[0].Request.QueryParams = []specification.Field{
		{
			Name:        "Limit",
			Description: "Limit results",
			Type:        testFieldTypeInt,
		},
	}
	return service
}

func getResourceNamePtr() *string {
	name := testResourceName
	return &name
}

// ============================================================================
// Header Generation Tests
// ============================================================================

func TestGenerateHeaderPopulation(t *testing.T) {
	testCases := []struct {
		name             string
		responseHeaders  []specification.Field
		packagePrefix    string
		expectedContains []string
	}{
		{
			name: "String type header",
			responseHeaders: []specification.Field{
				{Name: "X-Custom-Header", Type: "String"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XCustomHeader: types.NewString(",
				"test-x-custom-header-value",
			},
		},
		{
			name: "Int type header",
			responseHeaders: []specification.Field{
				{Name: "X-Rate-Limit", Type: "Int"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XRateLimit: types.NewInt(12345)",
			},
		},
		{
			name: "Int64 type header",
			responseHeaders: []specification.Field{
				{Name: "RateLimit-Reset", Type: "Int64"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"RateLimitReset: types.NewInt64(67890)",
			},
		},
		{
			name: "Bool type header",
			responseHeaders: []specification.Field{
				{Name: "X-Is-Active", Type: "Bool"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XIsActive: types.NewBool(true)",
			},
		},
		{
			name: "UUID type header",
			responseHeaders: []specification.Field{
				{Name: "X-Request-ID", Type: "UUID"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XRequestID: types.NewUUID(uuid.New())",
			},
		},
		{
			name: "Date type header",
			responseHeaders: []specification.Field{
				{Name: "X-Created-Date", Type: "Date"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XCreatedDate: types.NewDate(\"2024-01-15\")",
			},
		},
		{
			name: "Timestamp type header",
			responseHeaders: []specification.Field{
				{Name: "X-Created-At", Type: "Timestamp"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XCreatedAt: types.NewTimestamp(\"2024-01-15T10:30:00Z\")",
			},
		},
		{
			name: "Multiple headers with different types",
			responseHeaders: []specification.Field{
				{Name: "X-Rate-Limit", Type: "Int"},
				{Name: "RateLimit-Reset", Type: "Int64"},
				{Name: "X-Request-ID", Type: "UUID"},
			},
			packagePrefix: "",
			expectedContains: []string{
				"XRateLimit: types.NewInt(12345)",
				"RateLimitReset: types.NewInt64(67890)",
				"XRequestID: types.NewUUID(uuid.New())",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			service := &specification.Service{
				ResponseHeaders: tc.responseHeaders,
			}
			buf := &bytes.Buffer{}

			// Act
			generateHeaderPopulation(buf, service, tc.packagePrefix)

			// Assert
			generatedCode := buf.String()
			for _, expected := range tc.expectedContains {
				assert.Contains(t, generatedCode, expected, "Generated code should contain expected string")
			}
		})
	}

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty response headers", func(t *testing.T) {
			service := &specification.Service{
				ResponseHeaders: []specification.Field{},
			}
			buf := &bytes.Buffer{}

			generateHeaderPopulation(buf, service, "")

			assert.Empty(t, buf.String(), "Should generate nothing for empty response headers")
		})

		t.Run("nil response headers", func(t *testing.T) {
			service := &specification.Service{
				ResponseHeaders: nil,
			}
			buf := &bytes.Buffer{}

			generateHeaderPopulation(buf, service, "")

			assert.Empty(t, buf.String(), "Should generate nothing for nil response headers")
		})
	})
}

func TestGenerateHeaderAssertions(t *testing.T) {
	testCases := []struct {
		name             string
		responseHeaders  []specification.Field
		expectedContains []string
	}{
		{
			name: "String type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Custom-Header", Type: "String"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"test-x-custom-header-value\", w.Header().Get(\"X-Custom-Header\")",
			},
		},
		{
			name: "Int type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Rate-Limit", Type: "Int"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"12345\", w.Header().Get(\"X-Rate-Limit\")",
			},
		},
		{
			name: "Int64 type header assertion",
			responseHeaders: []specification.Field{
				{Name: "RateLimit-Reset", Type: "Int64"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"67890\", w.Header().Get(\"RateLimit-Reset\")",
			},
		},
		{
			name: "Bool type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Is-Active", Type: "Bool"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"true\", w.Header().Get(\"X-Is-Active\")",
			},
		},
		{
			name: "UUID type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Request-ID", Type: "UUID"},
			},
			expectedContains: []string{
				"assert.NotEmpty(t, w.Header().Get(\"X-Request-ID\")",
			},
		},
		{
			name: "Date type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Created-Date", Type: "Date"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"2024-01-15\", w.Header().Get(\"X-Created-Date\")",
			},
		},
		{
			name: "Timestamp type header assertion",
			responseHeaders: []specification.Field{
				{Name: "X-Created-At", Type: "Timestamp"},
			},
			expectedContains: []string{
				"assert.Equal(t, \"2024-01-15T10:30:00Z\", w.Header().Get(\"X-Created-At\")",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			service := &specification.Service{
				ResponseHeaders: tc.responseHeaders,
			}
			buf := &bytes.Buffer{}

			// Act
			generateHeaderAssertions(buf, service)

			// Assert
			generatedCode := buf.String()
			for _, expected := range tc.expectedContains {
				assert.Contains(t, generatedCode, expected, "Generated code should contain expected assertion")
			}
		})
	}

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty response headers", func(t *testing.T) {
			service := &specification.Service{
				ResponseHeaders: []specification.Field{},
			}
			buf := &bytes.Buffer{}

			generateHeaderAssertions(buf, service)

			assert.Empty(t, buf.String(), "Should generate nothing for empty response headers")
		})
	})
}
