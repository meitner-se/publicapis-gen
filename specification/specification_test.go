package specification

import (
	"encoding/json"
	"testing"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Service_JSONMarshaling(t *testing.T) {
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
						Description: "Get a user by ID",
						Method:      "GET",
						Path:        "/users/{id}",
						PathParameters: []Field{
							{Name: "id", Type: "UUID", Description: "User ID"},
						},
						ResponseFields: []Field{
							{Name: "id", Type: "UUID", Description: "User ID"},
							{Name: "name", Type: "String", Description: "User name"},
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

func Test_Service_YAMLMarshaling(t *testing.T) {
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
						Description: "List all people",
						Method:      "GET",
						Path:        "/people",
						ResponseFields: []Field{
							{Name: "data", Type: "array", Description: "Array of people", Modifiers: []string{"array"}},
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

func Test_Enum_Structure(t *testing.T) {
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

func Test_Field_WithModifiers(t *testing.T) {
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

func Test_ResourceField_InheritanceFromField(t *testing.T) {
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

func Test_Endpoint_CompleteStructure(t *testing.T) {
	endpoint := Endpoint{
		Name:        "CreateUser",
		Description: "Create a new user",
		Method:      "POST",
		Path:        "/users",
		RequestFields: []Field{
			{Name: "name", Type: "String", Description: "User name", Example: "John Doe"},
			{Name: "email", Type: "String", Description: "User email", Example: "john@example.com"},
		},
		ResponseFields: []Field{
			{Name: "id", Type: "UUID", Description: "Created user ID"},
			{Name: "name", Type: "String", Description: "User name"},
			{Name: "email", Type: "String", Description: "User email"},
		},
		QueryParameters: []Field{
			{Name: "validate", Type: "Bool", Description: "Whether to validate the user", Default: "true"},
		},
		PathParameters: []Field{},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(endpoint)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "CreateUser")
	assert.Contains(t, string(jsonData), "POST")
	assert.Contains(t, string(jsonData), "/users")

	// Test JSON unmarshaling
	var unmarshaledEndpoint Endpoint
	err = json.Unmarshal(jsonData, &unmarshaledEndpoint)
	require.NoError(t, err)
	assert.Equal(t, endpoint.Name, unmarshaledEndpoint.Name)
	assert.Equal(t, endpoint.Method, unmarshaledEndpoint.Method)
	assert.Equal(t, endpoint.Path, unmarshaledEndpoint.Path)
	assert.Equal(t, len(endpoint.RequestFields), len(unmarshaledEndpoint.RequestFields))
	assert.Equal(t, len(endpoint.ResponseFields), len(unmarshaledEndpoint.ResponseFields))
	assert.Equal(t, len(endpoint.QueryParameters), len(unmarshaledEndpoint.QueryParameters))
}

func Test_Resource_CompleteStructure(t *testing.T) {
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
				Description: "Get product by ID",
				Method:      "GET",
				Path:        "/products/{id}",
				PathParameters: []Field{
					{Name: "id", Type: "UUID", Description: "Product ID"},
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