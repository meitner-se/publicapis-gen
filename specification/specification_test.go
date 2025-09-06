package specification

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceJSONMarshaling(t *testing.T) {
	service := Service{
		Name: "TestService",
		Enums: []Enum{
			{
				Name:        "Gender",
				Description: "Gender enumeration",
				Values: []EnumValue{
					{Name: "Male", Description: "Male gender"},
					{Name: "Female", Description: "Female gender"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "User",
				Description: "User object",
				Fields: []Field{
					{Name: "id", Type: "UUID", Description: "User ID"},
					{Name: "name", Type: "String", Description: "User name"},
				},
			},
		},
		Resources: []Resource{
			{
				Name:        "Users",
				Description: "User resource",
				Operations:  []string{"Create", "Read", "Update", "Delete"},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "id",
							Type:        "UUID",
							Description: "User ID",
						},
						Operations: []string{"Read"},
					},
				},
				Endpoints: []Endpoint{
					{
						Name:        "GetUser",
						Title:       "Get User",
						Description: "Get a user by ID",
						Method:      "GET",
						Path:        "/{id}",
						Request: EndpointRequest{
							PathParams: []Field{
								{Name: "id", Type: "UUID", Description: "User ID"},
							},
						},
						Response: EndpointResponse{
							StatusCode: 200,
							BodyFields: []Field{
								{Name: "id", Type: "UUID", Description: "User ID"},
								{Name: "name", Type: "String", Description: "User name"},
							},
						},
					},
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(service)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaledService Service
	err = json.Unmarshal(jsonData, &unmarshaledService)
	require.NoError(t, err)
	assert.Equal(t, service.Name, unmarshaledService.Name)
	assert.Equal(t, len(service.Enums), len(unmarshaledService.Enums))
	assert.Equal(t, len(service.Objects), len(unmarshaledService.Objects))
	assert.Equal(t, len(service.Resources), len(unmarshaledService.Resources))
}

func TestServiceYAMLMarshaling(t *testing.T) {
	service := Service{
		Name: "TestService",
		Enums: []Enum{
			{
				Name:        "Status",
				Description: "Status enumeration",
				Values: []EnumValue{
					{Name: "Active", Description: "Active status"},
					{Name: "Inactive", Description: "Inactive status"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "Person",
				Description: "Person object",
				Fields: []Field{
					{Name: "id", Type: "UUID", Description: "Person ID"},
					{Name: "email", Type: "String", Description: "Email address"},
				},
			},
		},
		Resources: []Resource{
			{
				Name:        "People",
				Description: "People resource",
				Operations:  []string{"Create", "Read"},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "id",
							Type:        "UUID",
							Description: "Person ID",
						},
						Operations: []string{"Read"},
					},
				},
				Endpoints: []Endpoint{
					{
						Name:        "ListPeople",
						Title:       "List People",
						Description: "List all people",
						Method:      "GET",
						Path:        "/",
						Request:     EndpointRequest{},
						Response: EndpointResponse{
							StatusCode: 200,
							BodyFields: []Field{
								{Name: "data", Type: "array", Description: "Array of people", Modifiers: []string{"array"}},
							},
						},
					},
				},
			},
		},
	}

	// Test YAML marshaling
	yamlData, err := yaml.Marshal(service)
	require.NoError(t, err)
	assert.NotEmpty(t, yamlData)

	// Test YAML unmarshaling
	var unmarshaledService Service
	err = yaml.Unmarshal(yamlData, &unmarshaledService)
	require.NoError(t, err)
	assert.Equal(t, service.Name, unmarshaledService.Name)
	assert.Equal(t, len(service.Enums), len(unmarshaledService.Enums))
	assert.Equal(t, len(service.Objects), len(unmarshaledService.Objects))
	assert.Equal(t, len(service.Resources), len(unmarshaledService.Resources))
}

func TestEnumStructure(t *testing.T) {
	enum := Enum{
		Name:        "Priority",
		Description: "Task priority levels",
		Values: []EnumValue{
			{Name: "Low", Description: "Low priority"},
			{Name: "Medium", Description: "Medium priority"},
			{Name: "High", Description: "High priority"},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(enum)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "Priority")
	assert.Contains(t, string(jsonData), "Low")
	assert.Contains(t, string(jsonData), "Medium")
	assert.Contains(t, string(jsonData), "High")

	// Test JSON unmarshaling
	var unmarshaledEnum Enum
	err = json.Unmarshal(jsonData, &unmarshaledEnum)
	require.NoError(t, err)
	assert.Equal(t, enum.Name, unmarshaledEnum.Name)
	assert.Equal(t, enum.Description, unmarshaledEnum.Description)
	assert.Equal(t, len(enum.Values), len(unmarshaledEnum.Values))
}

func TestFieldWithModifiers(t *testing.T) {
	field := Field{
		Name:        "tags",
		Description: "List of tags",
		Type:        "String",
		Default:     "[]",
		Example:     `["tag1", "tag2"]`,
		Modifiers:   []string{"array", "nullable"},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(field)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "array")
	assert.Contains(t, string(jsonData), "nullable")

	// Test JSON unmarshaling
	var unmarshaledField Field
	err = json.Unmarshal(jsonData, &unmarshaledField)
	require.NoError(t, err)
	assert.Equal(t, field.Name, unmarshaledField.Name)
	assert.Equal(t, field.Type, unmarshaledField.Type)
	assert.Equal(t, field.Modifiers, unmarshaledField.Modifiers)
}

func TestResourceFieldInheritanceFromField(t *testing.T) {
	resourceField := ResourceField{
		Field: Field{
			Name:        "username",
			Description: "User's username",
			Type:        "String",
			Example:     "johndoe",
		},
		Operations: []string{"Create", "Read", "Update"},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(resourceField)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "username")
	assert.Contains(t, string(jsonData), "Create")
	assert.Contains(t, string(jsonData), "Update")

	// Test JSON unmarshaling
	var unmarshaledResourceField ResourceField
	err = json.Unmarshal(jsonData, &unmarshaledResourceField)
	require.NoError(t, err)
	assert.Equal(t, resourceField.Name, unmarshaledResourceField.Name)
	assert.Equal(t, resourceField.Description, unmarshaledResourceField.Description)
	assert.Equal(t, resourceField.Operations, unmarshaledResourceField.Operations)
}

func TestEndpointCompleteStructure(t *testing.T) {
	endpoint := Endpoint{
		Name:        "CreateUser",
		Title:       "Create User",
		Description: "Create a new user",
		Method:      "POST",
		Path:        "/",
		Request: EndpointRequest{
			BodyParams: []Field{
				{Name: "name", Type: "String", Description: "User name", Example: "John Doe"},
				{Name: "email", Type: "String", Description: "User email", Example: "john@example.com"},
			},
			QueryParams: []Field{
				{Name: "validate", Type: "Bool", Description: "Whether to validate the user", Default: "true"},
			},
		},
		Response: EndpointResponse{
			StatusCode: 201,
			BodyFields: []Field{
				{Name: "id", Type: "UUID", Description: "Created user ID"},
				{Name: "name", Type: "String", Description: "User name"},
				{Name: "email", Type: "String", Description: "User email"},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(endpoint)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "CreateUser")
	assert.Contains(t, string(jsonData), "POST")
	assert.Contains(t, string(jsonData), "\"/\"")

	// Test JSON unmarshaling
	var unmarshaledEndpoint Endpoint
	err = json.Unmarshal(jsonData, &unmarshaledEndpoint)
	require.NoError(t, err)
	assert.Equal(t, endpoint.Name, unmarshaledEndpoint.Name)
	assert.Equal(t, endpoint.Title, unmarshaledEndpoint.Title)
	assert.Equal(t, endpoint.Method, unmarshaledEndpoint.Method)
	assert.Equal(t, endpoint.Path, unmarshaledEndpoint.Path)
	assert.Equal(t, len(endpoint.Request.BodyParams), len(unmarshaledEndpoint.Request.BodyParams))
	assert.Equal(t, len(endpoint.Response.BodyFields), len(unmarshaledEndpoint.Response.BodyFields))
	assert.Equal(t, len(endpoint.Request.QueryParams), len(unmarshaledEndpoint.Request.QueryParams))
}

func TestEndpointRequestStructure(t *testing.T) {
	endpointRequest := EndpointRequest{
		ContentType: "application/json",
		Headers: []Field{
			{Name: "Authorization", Type: "String", Description: "Bearer token"},
		},
		PathParams: []Field{
			{Name: "id", Type: "UUID", Description: "Resource ID"},
		},
		QueryParams: []Field{
			{Name: "limit", Type: "Int", Description: "Number of items to return", Default: "10"},
		},
		BodyParams: []Field{
			{Name: "name", Type: "String", Description: "Resource name"},
			{Name: "active", Type: "Bool", Description: "Is resource active", Default: "true"},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(endpointRequest)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "application/json")
	assert.Contains(t, string(jsonData), "Authorization")
	assert.Contains(t, string(jsonData), "Bearer token")

	// Test JSON unmarshaling
	var unmarshaledRequest EndpointRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	require.NoError(t, err)
	assert.Equal(t, endpointRequest.ContentType, unmarshaledRequest.ContentType)
	assert.Equal(t, len(endpointRequest.Headers), len(unmarshaledRequest.Headers))
	assert.Equal(t, len(endpointRequest.PathParams), len(unmarshaledRequest.PathParams))
	assert.Equal(t, len(endpointRequest.QueryParams), len(unmarshaledRequest.QueryParams))
	assert.Equal(t, len(endpointRequest.BodyParams), len(unmarshaledRequest.BodyParams))
}

func TestEndpointResponseStructure(t *testing.T) {
	endpointResponse := EndpointResponse{
		ContentType: "application/json",
		StatusCode:  201,
		Headers: []Field{
			{Name: "Location", Type: "String", Description: "URL of created resource"},
		},
		BodyFields: []Field{
			{Name: "id", Type: "UUID", Description: "Created resource ID"},
			{Name: "created_at", Type: "Timestamp", Description: "Creation timestamp"},
		},
		BodyObject: nil,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(endpointResponse)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "application/json")
	assert.Contains(t, string(jsonData), "201")
	assert.Contains(t, string(jsonData), "Location")

	// Test JSON unmarshaling
	var unmarshaledResponse EndpointResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	require.NoError(t, err)
	assert.Equal(t, endpointResponse.ContentType, unmarshaledResponse.ContentType)
	assert.Equal(t, endpointResponse.StatusCode, unmarshaledResponse.StatusCode)
	assert.Equal(t, len(endpointResponse.Headers), len(unmarshaledResponse.Headers))
	assert.Equal(t, len(endpointResponse.BodyFields), len(unmarshaledResponse.BodyFields))
}

func TestEndpointResponseWithBodyObject(t *testing.T) {
	objectName := "User"
	endpointResponse := EndpointResponse{
		ContentType: "application/json",
		StatusCode:  200,
		Headers:     []Field{},
		BodyFields:  []Field{},
		BodyObject:  &objectName,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(endpointResponse)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "User")

	// Test JSON unmarshaling
	var unmarshaledResponse EndpointResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	require.NoError(t, err)
	assert.Equal(t, endpointResponse.StatusCode, unmarshaledResponse.StatusCode)
	require.NotNil(t, unmarshaledResponse.BodyObject)
	assert.Equal(t, "User", *unmarshaledResponse.BodyObject)
}

func TestResourceCompleteStructure(t *testing.T) {
	resource := Resource{
		Name:        "Products",
		Description: "Product management resource",
		Operations:  []string{"Create", "Read", "Update", "Delete"},
		Fields: []ResourceField{
			{
				Field: Field{
					Name:        "id",
					Type:        "UUID",
					Description: "Product ID",
				},
				Operations: []string{"Read"},
			},
			{
				Field: Field{
					Name:        "name",
					Type:        "String",
					Description: "Product name",
					Example:     "Widget",
				},
				Operations: []string{"Create", "Read", "Update"},
			},
		},
		Endpoints: []Endpoint{
			{
				Name:        "GetProduct",
				Title:       "Get Product",
				Description: "Get product by ID",
				Method:      "GET",
				Path:        "/{id}",
				Request: EndpointRequest{
					PathParams: []Field{
						{Name: "id", Type: "UUID", Description: "Product ID"},
					},
				},
				Response: EndpointResponse{
					StatusCode: 200,
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(resource)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "Products")
	assert.Contains(t, string(jsonData), "Create")
	assert.Contains(t, string(jsonData), "Delete")

	// Test JSON unmarshaling
	var unmarshaledResource Resource
	err = json.Unmarshal(jsonData, &unmarshaledResource)
	require.NoError(t, err)
	assert.Equal(t, resource.Name, unmarshaledResource.Name)
	assert.Equal(t, resource.Description, unmarshaledResource.Description)
	assert.Equal(t, resource.Operations, unmarshaledResource.Operations)
	assert.Equal(t, len(resource.Fields), len(unmarshaledResource.Fields))
	assert.Equal(t, len(resource.Endpoints), len(unmarshaledResource.Endpoints))
}

func TestEmptyStructures(t *testing.T) {
	// Test empty structures can be marshaled/unmarshaled
	testCases := []struct {
		name string
		data interface{}
	}{
		{"EmptyService", Service{}},
		{"EmptyEnum", Enum{}},
		{"EmptyObject", Object{}},
		{"EmptyResource", Resource{}},
		{"EmptyField", Field{}},
		{"EmptyResourceField", ResourceField{}},
		{"EmptyEndpoint", Endpoint{}},
		{"EmptyEndpointRequest", EndpointRequest{}},
		{"EmptyEndpointResponse", EndpointResponse{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tc.data)
			require.NoError(t, err)
			assert.NotEmpty(t, jsonData)

			// Test YAML marshaling
			yamlData, err := yaml.Marshal(tc.data)
			require.NoError(t, err)
			assert.NotEmpty(t, yamlData)
		})
	}
}

func TestFieldModifiersEdgeCases(t *testing.T) {
	// Test field with multiple modifiers
	field := Field{
		Name:        "complexField",
		Description: "A complex field with multiple modifiers",
		Type:        "String",
		Default:     "",
		Example:     "",
		Modifiers:   []string{"nullable", "array", "optional"},
	}

	jsonData, err := json.Marshal(field)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "nullable")
	assert.Contains(t, string(jsonData), "array")
	assert.Contains(t, string(jsonData), "optional")

	var unmarshaledField Field
	err = json.Unmarshal(jsonData, &unmarshaledField)
	require.NoError(t, err)
	assert.Equal(t, field.Modifiers, unmarshaledField.Modifiers)
}

func TestEndpointRequestWithAllFields(t *testing.T) {
	request := EndpointRequest{
		ContentType: "multipart/form-data",
		Headers: []Field{
			{Name: "X-API-Key", Type: "String", Description: "API key"},
			{Name: "User-Agent", Type: "String", Description: "User agent"},
		},
		PathParams: []Field{
			{Name: "userId", Type: "UUID", Description: "User ID"},
			{Name: "resourceId", Type: "UUID", Description: "Resource ID"},
		},
		QueryParams: []Field{
			{Name: "page", Type: "Int", Description: "Page number", Default: "1"},
			{Name: "limit", Type: "Int", Description: "Items per page", Default: "10"},
		},
		BodyParams: []Field{
			{Name: "data", Type: "Object", Description: "Request data"},
			{Name: "metadata", Type: "Object", Description: "Additional metadata"},
		},
	}

	jsonData, err := json.Marshal(request)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "multipart/form-data")
	assert.Contains(t, string(jsonData), "X-API-Key")
	assert.Contains(t, string(jsonData), "userId")
	assert.Contains(t, string(jsonData), "page")
	assert.Contains(t, string(jsonData), "data")

	var unmarshaledRequest EndpointRequest
	err = json.Unmarshal(jsonData, &unmarshaledRequest)
	require.NoError(t, err)
	assert.Equal(t, request.ContentType, unmarshaledRequest.ContentType)
	assert.Equal(t, len(request.Headers), len(unmarshaledRequest.Headers))
	assert.Equal(t, len(request.PathParams), len(unmarshaledRequest.PathParams))
	assert.Equal(t, len(request.QueryParams), len(unmarshaledRequest.QueryParams))
	assert.Equal(t, len(request.BodyParams), len(unmarshaledRequest.BodyParams))
}

func TestEndpointResponseErrorCodes(t *testing.T) {
	// Test various HTTP status codes
	statusCodes := []int{200, 201, 400, 401, 403, 404, 500}

	for _, code := range statusCodes {
		t.Run(fmt.Sprintf("StatusCode%d", code), func(t *testing.T) {
			response := EndpointResponse{
				ContentType: "application/json",
				StatusCode:  code,
				Headers: []Field{
					{Name: "X-Request-ID", Type: "String", Description: "Request ID"},
				},
				BodyFields: []Field{
					{Name: "message", Type: "String", Description: "Response message"},
				},
			}

			jsonData, err := json.Marshal(response)
			require.NoError(t, err)
			assert.Contains(t, string(jsonData), fmt.Sprintf("%d", code))

			var unmarshaledResponse EndpointResponse
			err = json.Unmarshal(jsonData, &unmarshaledResponse)
			require.NoError(t, err)
			assert.Equal(t, code, unmarshaledResponse.StatusCode)
		})
	}
}

func TestServiceWithComplexHierarchy(t *testing.T) {
	// Create a complex service structure to test deep nesting
	service := Service{
		Name: "ComplexAPI",
		Enums: []Enum{
			{
				Name:        "UserRole",
				Description: "User roles in the system",
				Values: []EnumValue{
					{Name: "Admin", Description: "Administrator role"},
					{Name: "User", Description: "Regular user role"},
					{Name: "Guest", Description: "Guest user role"},
				},
			},
			{
				Name:        "Status",
				Description: "Entity status",
				Values: []EnumValue{
					{Name: "Active", Description: "Entity is active"},
					{Name: "Inactive", Description: "Entity is inactive"},
					{Name: "Pending", Description: "Entity is pending"},
				},
			},
		},
		Objects: []Object{
			{
				Name:        "Address",
				Description: "Address information",
				Fields: []Field{
					{Name: "street", Type: "String", Description: "Street address"},
					{Name: "city", Type: "String", Description: "City"},
					{Name: "zipCode", Type: "String", Description: "ZIP code"},
					{Name: "country", Type: "String", Description: "Country"},
				},
			},
			{
				Name:        "ContactInfo",
				Description: "Contact information",
				Fields: []Field{
					{Name: "email", Type: "String", Description: "Email address"},
					{Name: "phone", Type: "String", Description: "Phone number"},
					{Name: "address", Type: "Address", Description: "Physical address"},
				},
			},
		},
		Resources: []Resource{
			{
				Name:        "Users",
				Description: "User management",
				Operations:  []string{"Create", "Read", "Update", "Delete"},
				Fields: []ResourceField{
					{
						Field: Field{
							Name:        "id",
							Type:        "UUID",
							Description: "User ID",
						},
						Operations: []string{"Read"},
					},
					{
						Field: Field{
							Name:        "role",
							Type:        "UserRole",
							Description: "User role",
							Default:     "User",
						},
						Operations: []string{"Create", "Read", "Update"},
					},
				},
				Endpoints: []Endpoint{
					{
						Name:        "CreateUser",
						Title:       "Create New User",
						Description: "Create a new user account",
						Method:      "POST",
						Path:        "/",
						Request: EndpointRequest{
							ContentType: "application/json",
							BodyParams: []Field{
								{Name: "username", Type: "String", Description: "Username"},
								{Name: "email", Type: "String", Description: "Email"},
								{Name: "role", Type: "UserRole", Description: "User role"},
							},
						},
						Response: EndpointResponse{
							ContentType: "application/json",
							StatusCode:  201,
							BodyFields: []Field{
								{Name: "id", Type: "UUID", Description: "Created user ID"},
								{Name: "username", Type: "String", Description: "Username"},
							},
						},
					},
				},
			},
		},
	}

	// Test JSON serialization of complex structure
	jsonData, err := json.MarshalIndent(service, "", "  ")
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test YAML serialization of complex structure
	yamlData, err := yaml.Marshal(service)
	require.NoError(t, err)
	assert.NotEmpty(t, yamlData)

	// Test deserialization
	var unmarshaledService Service
	err = json.Unmarshal(jsonData, &unmarshaledService)
	require.NoError(t, err)

	// Verify structure integrity
	assert.Equal(t, service.Name, unmarshaledService.Name)
	assert.Equal(t, len(service.Enums), len(unmarshaledService.Enums))
	assert.Equal(t, len(service.Objects), len(unmarshaledService.Objects))
	assert.Equal(t, len(service.Resources), len(unmarshaledService.Resources))

	// Test nested structures
	assert.Equal(t, len(service.Enums[0].Values), len(unmarshaledService.Enums[0].Values))
	assert.Equal(t, len(service.Objects[0].Fields), len(unmarshaledService.Objects[0].Fields))
	assert.Equal(t, len(service.Resources[0].Fields), len(unmarshaledService.Resources[0].Fields))
	assert.Equal(t, len(service.Resources[0].Endpoints), len(unmarshaledService.Resources[0].Endpoints))
}

func TestApplyOverlay(t *testing.T) {
	t.Run("NilInput", func(t *testing.T) {
		result := ApplyOverlay(nil)
		assert.Nil(t, result)
	})

	t.Run("EmptyService", func(t *testing.T) {
		input := &Service{
			Name:      "EmptyService",
			Enums:     []Enum{},
			Objects:   []Object{},
			Resources: []Resource{},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, 0, len(result.Objects))
		assert.Equal(t, 0, len(result.Resources))
	})

	t.Run("ResourceWithReadOperation", func(t *testing.T) {
		input := &Service{
			Name:    "TestService",
			Enums:   []Enum{},
			Objects: []Object{},
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User management resource",
					Operations:  []string{"Create", "Read", "Update", "Delete"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "User ID",
							},
							Operations: []string{"Read"},
						},
						{
							Field: Field{
								Name:        "name",
								Type:        "String",
								Description: "User name",
								Example:     "John Doe",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "password",
								Type:        "String",
								Description: "User password",
							},
							Operations: []string{"Create", "Update"}, // No Read operation
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should have generated one object
		assert.Equal(t, 1, len(result.Objects))

		// Check the generated object
		userObject := result.Objects[0]
		assert.Equal(t, "Users", userObject.Name)
		assert.Equal(t, "User management resource", userObject.Description)
		assert.Equal(t, 2, len(userObject.Fields)) // Only id and name have Read operation

		// Check fields are correct
		assert.Equal(t, "id", userObject.Fields[0].Name)
		assert.Equal(t, "UUID", userObject.Fields[0].Type)
		assert.Equal(t, "name", userObject.Fields[1].Name)
		assert.Equal(t, "String", userObject.Fields[1].Type)
		assert.Equal(t, "John Doe", userObject.Fields[1].Example)
	})

	t.Run("ResourceWithoutReadOperation", func(t *testing.T) {
		input := &Service{
			Name:    "TestService",
			Enums:   []Enum{},
			Objects: []Object{},
			Resources: []Resource{
				{
					Name:        "InternalLogs",
					Description: "Internal logging resource",
					Operations:  []string{"Create", "Delete"}, // No Read operation
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "Log ID",
							},
							Operations: []string{"Create"},
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should not have generated any objects
		assert.Equal(t, 0, len(result.Objects))
	})

	t.Run("MultipleResourcesWithReadOperation", func(t *testing.T) {
		input := &Service{
			Name:    "TestService",
			Enums:   []Enum{},
			Objects: []Object{},
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User resource",
					Operations:  []string{"Create", "Read"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "User ID",
							},
							Operations: []string{"Read"},
						},
					},
				},
				{
					Name:        "Products",
					Description: "Product resource",
					Operations:  []string{"Read", "Update"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "Product ID",
							},
							Operations: []string{"Read"},
						},
						{
							Field: Field{
								Name:        "name",
								Type:        "String",
								Description: "Product name",
							},
							Operations: []string{"Read", "Update"},
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should have generated two objects
		assert.Equal(t, 2, len(result.Objects))

		// Check first object (Users)
		usersObject := result.Objects[0]
		assert.Equal(t, "Users", usersObject.Name)
		assert.Equal(t, 1, len(usersObject.Fields))

		// Check second object (Products)
		productsObject := result.Objects[1]
		assert.Equal(t, "Products", productsObject.Name)
		assert.Equal(t, 2, len(productsObject.Fields)) // Both id and name have Read operation
	})

	t.Run("ExistingObjectWithSameName", func(t *testing.T) {
		input := &Service{
			Name:  "TestService",
			Enums: []Enum{},
			Objects: []Object{
				{
					Name:        "Users",
					Description: "Existing user object",
					Fields: []Field{
						{
							Name:        "existingField",
							Type:        "String",
							Description: "Existing field",
						},
					},
				},
			},
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User resource",
					Operations:  []string{"Create", "Read"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "User ID",
							},
							Operations: []string{"Read"},
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should still have only one object (the existing one)
		assert.Equal(t, 1, len(result.Objects))

		// Should be the original object, not the generated one
		assert.Equal(t, "Users", result.Objects[0].Name)
		assert.Equal(t, "Existing user object", result.Objects[0].Description)
		assert.Equal(t, 1, len(result.Objects[0].Fields))
		assert.Equal(t, "existingField", result.Objects[0].Fields[0].Name)
	})

	t.Run("FieldsWithModifiers", func(t *testing.T) {
		input := &Service{
			Name:    "TestService",
			Enums:   []Enum{},
			Objects: []Object{},
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User resource",
					Operations:  []string{"Create", "Read"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "tags",
								Type:        "String",
								Description: "User tags",
								Default:     "[]",
								Example:     `["admin", "user"]`,
								Modifiers:   []string{"array", "nullable"},
							},
							Operations: []string{"Create", "Read", "Update"},
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should have generated one object
		assert.Equal(t, 1, len(result.Objects))

		// Check field modifiers are preserved
		userObject := result.Objects[0]
		assert.Equal(t, 1, len(userObject.Fields))

		field := userObject.Fields[0]
		assert.Equal(t, "tags", field.Name)
		assert.Equal(t, "String", field.Type)
		assert.Equal(t, "[]", field.Default)
		assert.Equal(t, `["admin", "user"]`, field.Example)
		assert.Equal(t, []string{"array", "nullable"}, field.Modifiers)
	})

	t.Run("PreservesOriginalServiceStructure", func(t *testing.T) {
		input := &Service{
			Name: "TestService",
			Enums: []Enum{
				{
					Name:        "Status",
					Description: "Status enum",
					Values: []EnumValue{
						{Name: "Active", Description: "Active status"},
					},
				},
			},
			Objects: []Object{
				{
					Name:        "ExistingObject",
					Description: "Pre-existing object",
					Fields: []Field{
						{Name: "field1", Type: "String", Description: "Field 1"},
					},
				},
			},
			Resources: []Resource{
				{
					Name:        "Users",
					Description: "User resource",
					Operations:  []string{"Read"},
					Fields: []ResourceField{
						{
							Field: Field{
								Name:        "id",
								Type:        "UUID",
								Description: "User ID",
							},
							Operations: []string{"Read"},
						},
					},
				},
			},
		}

		result := ApplyOverlay(input)
		require.NotNil(t, result)

		// Should preserve all original structure
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, len(input.Enums), len(result.Enums))
		assert.Equal(t, len(input.Resources), len(result.Resources))

		// Should have existing object plus generated object
		assert.Equal(t, 2, len(result.Objects))

		// First object should be the existing one
		assert.Equal(t, "ExistingObject", result.Objects[0].Name)

		// Second object should be the generated one
		assert.Equal(t, "Users", result.Objects[1].Name)
	})
}

func Test_containsOperation(t *testing.T) {
	t.Run("EmptySlice", func(t *testing.T) {
		result := containsOperation([]string{}, "Read")
		assert.False(t, result)
	})

	t.Run("OperationExists", func(t *testing.T) {
		operations := []string{"Create", "Read", "Update", "Delete"}
		result := containsOperation(operations, "Read")
		assert.True(t, result)
	})

	t.Run("OperationDoesNotExist", func(t *testing.T) {
		operations := []string{"Create", "Update", "Delete"}
		result := containsOperation(operations, "Read")
		assert.False(t, result)
	})

	t.Run("CaseSensitive", func(t *testing.T) {
		operations := []string{"read"}
		result := containsOperation(operations, "Read")
		assert.False(t, result) // Should be case-sensitive
	})
}
