package specification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Resource method tests

func TestResourceHasCreateOperation(t *testing.T) {
	// Test resource with Create operation
	resourceWithCreate := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationRead},
	}

	result := resourceWithCreate.HasCreateOperation()
	assert.True(t, result, "Resource with Create operation should return true")

	// Test resource without Create operation
	resourceWithoutCreate := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationRead, OperationUpdate},
	}

	result = resourceWithoutCreate.HasCreateOperation()
	assert.False(t, result, "Resource without Create operation should return false")

	// Test resource with empty operations
	resourceEmpty := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{},
	}

	result = resourceEmpty.HasCreateOperation()
	assert.False(t, result, "Resource with no operations should return false")
}

func TestResourceHasReadOperation(t *testing.T) {
	// Test resource with Read operation
	resourceWithRead := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationRead},
	}

	result := resourceWithRead.HasReadOperation()
	assert.True(t, result, "Resource with Read operation should return true")

	// Test resource without Read operation
	resourceWithoutRead := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationUpdate},
	}

	result = resourceWithoutRead.HasReadOperation()
	assert.False(t, result, "Resource without Read operation should return false")
}

func TestResourceHasUpdateOperation(t *testing.T) {
	// Test resource with Update operation
	resourceWithUpdate := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationUpdate, OperationRead},
	}

	result := resourceWithUpdate.HasUpdateOperation()
	assert.True(t, result, "Resource with Update operation should return true")

	// Test resource without Update operation
	resourceWithoutUpdate := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationRead},
	}

	result = resourceWithoutUpdate.HasUpdateOperation()
	assert.False(t, result, "Resource without Update operation should return false")
}

func TestResourceHasDeleteOperation(t *testing.T) {
	// Test resource with Delete operation
	resourceWithDelete := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationDelete, OperationRead},
	}

	result := resourceWithDelete.HasDeleteOperation()
	assert.True(t, result, "Resource with Delete operation should return true")

	// Test resource without Delete operation
	resourceWithoutDelete := Resource{
		Name:        "Users",
		Description: "User resource",
		Operations:  []string{OperationCreate, OperationRead},
	}

	result = resourceWithoutDelete.HasDeleteOperation()
	assert.False(t, result, "Resource without Delete operation should return false")
}

func TestResourceGetPluralName(t *testing.T) {
	testCases := []struct {
		resourceName   string
		expectedPlural string
	}{
		{"User", "Users"},
		{"Category", "Categories"},
		{"Person", "Persons"},
		{"Child", "Childs"},
		{"Company", "Companies"},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceName, func(t *testing.T) {
			resource := Resource{Name: tc.resourceName}
			result := resource.GetPluralName()
			assert.Equal(t, tc.expectedPlural, result, "Plural name for '%s' should be '%s'", tc.resourceName, tc.expectedPlural)
		})
	}
}

func TestResourceGetCreateBodyParams(t *testing.T) {
	// Arrange
	expectedFieldName1 := "username"
	expectedFieldName2 := "email"
	expectedFieldDescription1 := "User's username"
	expectedFieldDescription2 := "User's email"
	expectedFieldType := FieldTypeString

	resource := Resource{
		Name: "Users",
		Fields: []ResourceField{
			{
				Field: Field{
					Name:        expectedFieldName1,
					Description: expectedFieldDescription1,
					Type:        expectedFieldType,
				},
				Operations: []string{OperationCreate, OperationRead},
			},
			{
				Field: Field{
					Name:        expectedFieldName2,
					Description: expectedFieldDescription2,
					Type:        expectedFieldType,
				},
				Operations: []string{OperationCreate, OperationUpdate},
			},
			{
				Field: Field{
					Name:        "id",
					Description: "User ID",
					Type:        FieldTypeUUID,
				},
				Operations: []string{OperationRead}, // No Create operation
			},
		},
	}

	// Act
	createParams := resource.GetCreateBodyParams()

	// Assert
	assert.Len(t, createParams, 2, "Should return exactly 2 fields with Create operations")

	// Check first field
	assert.Equal(t, expectedFieldName1, createParams[0].Name, "First field name should match")
	assert.Equal(t, expectedFieldDescription1, createParams[0].Description, "First field description should match")
	assert.Equal(t, expectedFieldType, createParams[0].Type, "First field type should match")

	// Check second field
	assert.Equal(t, expectedFieldName2, createParams[1].Name, "Second field name should match")
	assert.Equal(t, expectedFieldDescription2, createParams[1].Description, "Second field description should match")
	assert.Equal(t, expectedFieldType, createParams[1].Type, "Second field type should match")
}

func TestResourceGetUpdateBodyParams(t *testing.T) {
	// Arrange
	expectedFieldName := "email"
	expectedFieldDescription := "User's email"
	expectedFieldType := FieldTypeString

	resource := Resource{
		Name: "Users",
		Fields: []ResourceField{
			{
				Field: Field{
					Name:        "username",
					Description: "User's username",
					Type:        FieldTypeString,
				},
				Operations: []string{OperationCreate, OperationRead}, // No Update
			},
			{
				Field: Field{
					Name:        expectedFieldName,
					Description: expectedFieldDescription,
					Type:        expectedFieldType,
				},
				Operations: []string{OperationCreate, OperationUpdate, OperationRead},
			},
			{
				Field: Field{
					Name:        "id",
					Description: "User ID",
					Type:        FieldTypeUUID,
				},
				Operations: []string{OperationRead}, // No Update
			},
		},
	}

	// Act
	updateParams := resource.GetUpdateBodyParams()

	// Assert
	assert.Len(t, updateParams, 1, "Should return exactly 1 field with Update operations")
	assert.Equal(t, expectedFieldName, updateParams[0].Name, "Field name should match")
	assert.Equal(t, expectedFieldDescription, updateParams[0].Description, "Field description should match")
	assert.Equal(t, expectedFieldType, updateParams[0].Type, "Field type should match")
}

func TestResourceGetReadableFields(t *testing.T) {
	// Arrange
	expectedReadableCount := 2

	resource := Resource{
		Name: "Users",
		Fields: []ResourceField{
			{
				Field: Field{
					Name:        "id",
					Description: "User ID",
					Type:        FieldTypeUUID,
				},
				Operations: []string{OperationRead},
			},
			{
				Field: Field{
					Name:        "username",
					Description: "User's username",
					Type:        FieldTypeString,
				},
				Operations: []string{OperationCreate, OperationRead, OperationUpdate},
			},
			{
				Field: Field{
					Name:        "password",
					Description: "User's password",
					Type:        FieldTypeString,
				},
				Operations: []string{OperationCreate, OperationUpdate}, // No Read
			},
		},
	}

	// Act
	readableFields := resource.GetReadableFields()

	// Assert
	assert.Len(t, readableFields, expectedReadableCount, "Should return exactly 2 readable fields")
	assert.Equal(t, "id", readableFields[0].Name, "First readable field should be 'id'")
	assert.Equal(t, "username", readableFields[1].Name, "Second readable field should be 'username'")
}

func TestResourceHasEndpoint(t *testing.T) {
	// Arrange
	existingEndpointName := "GetUser"
	nonExistentEndpointName := "DeleteUser"

	resource := Resource{
		Name: "Users",
		Endpoints: []Endpoint{
			{
				Name:        existingEndpointName,
				Description: "Get user by ID",
				Method:      httpMethodGet,
				Path:        "/{id}",
			},
			{
				Name:        "CreateUser",
				Description: "Create new user",
				Method:      httpMethodPost,
				Path:        "",
			},
		},
	}

	// Act & Assert - existing endpoint
	result := resource.HasEndpoint(existingEndpointName)
	assert.True(t, result, "Should return true for existing endpoint '%s'", existingEndpointName)

	// Act & Assert - non-existent endpoint
	result = resource.HasEndpoint(nonExistentEndpointName)
	assert.False(t, result, "Should return false for non-existent endpoint '%s'", nonExistentEndpointName)
}

// Field method tests

func TestFieldIsArray(t *testing.T) {
	// Test field with array modifier
	fieldWithArray := Field{
		Name:      "tags",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierArray},
	}

	result := fieldWithArray.IsArray()
	assert.True(t, result, "Field with array modifier should return true")

	// Test field without array modifier
	fieldWithoutArray := Field{
		Name:      "name",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierNullable},
	}

	result = fieldWithoutArray.IsArray()
	assert.False(t, result, "Field without array modifier should return false")

	// Test field with no modifiers
	fieldNoModifiers := Field{
		Name: "name",
		Type: FieldTypeString,
	}

	result = fieldNoModifiers.IsArray()
	assert.False(t, result, "Field with no modifiers should return false")
}

func TestFieldIsNullable(t *testing.T) {
	// Test field with nullable modifier
	fieldWithNullable := Field{
		Name:      "description",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierNullable},
	}

	result := fieldWithNullable.IsNullable()
	assert.True(t, result, "Field with nullable modifier should return true")

	// Test field without nullable modifier
	fieldWithoutNullable := Field{
		Name:      "name",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierArray},
	}

	result = fieldWithoutNullable.IsNullable()
	assert.False(t, result, "Field without nullable modifier should return false")

	// Test field with both modifiers
	fieldWithBoth := Field{
		Name:      "tags",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierArray, ModifierNullable},
	}

	result = fieldWithBoth.IsNullable()
	assert.True(t, result, "Field with nullable modifier (among others) should return true")
}

func TestFieldTagJSON(t *testing.T) {
	testCases := []struct {
		fieldName   string
		expectedTag string
	}{
		{"user_name", "userName"},
		{"first_name", "firstName"},
		{"id", "id"},
		{"created_at", "createdAt"},
		{"user_id", "userID"},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			field := Field{Name: tc.fieldName}
			result := field.TagJSON()
			assert.Equal(t, tc.expectedTag, result, "JSON tag for field '%s' should be '%s'", tc.fieldName, tc.expectedTag)
		})
	}
}

func TestFieldIsRequired(t *testing.T) {
	service := &Service{
		Objects: []Object{
			{Name: "User", Fields: []Field{{Name: "id", Type: FieldTypeUUID}}},
		},
	}

	// Test required field (no nullable, no array, no default, not object type)
	requiredField := Field{
		Name: "username",
		Type: FieldTypeString,
	}

	result := requiredField.IsRequired(service)
	assert.True(t, result, "Field without nullable, array, default, or object type should be required")

	// Test nullable field (not required)
	nullableField := Field{
		Name:      "description",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierNullable},
	}

	result = nullableField.IsRequired(service)
	assert.False(t, result, "Nullable field should not be required")

	// Test array field (not required)
	arrayField := Field{
		Name:      "tags",
		Type:      FieldTypeString,
		Modifiers: []string{ModifierArray},
	}

	result = arrayField.IsRequired(service)
	assert.False(t, result, "Array field should not be required")

	// Test field with default (not required)
	fieldWithDefault := Field{
		Name:    "status",
		Type:    FieldTypeString,
		Default: "active",
	}

	result = fieldWithDefault.IsRequired(service)
	assert.False(t, result, "Field with default value should not be required")

	// Test object type field (not required)
	objectField := Field{
		Name: "user",
		Type: "User", // Matches object in service
	}

	result = objectField.IsRequired(service)
	assert.False(t, result, "Object type field should not be required")
}

// ResourceField method tests

func TestResourceFieldHasCreateOperation(t *testing.T) {
	// Test ResourceField with Create operation
	fieldWithCreate := ResourceField{
		Field: Field{
			Name: "username",
			Type: FieldTypeString,
		},
		Operations: []string{OperationCreate, OperationRead},
	}

	result := fieldWithCreate.HasCreateOperation()
	assert.True(t, result, "ResourceField with Create operation should return true")

	// Test ResourceField without Create operation
	fieldWithoutCreate := ResourceField{
		Field: Field{
			Name: "id",
			Type: FieldTypeUUID,
		},
		Operations: []string{OperationRead},
	}

	result = fieldWithoutCreate.HasCreateOperation()
	assert.False(t, result, "ResourceField without Create operation should return false")
}

func TestResourceFieldHasReadOperation(t *testing.T) {
	// Test ResourceField with Read operation
	fieldWithRead := ResourceField{
		Field: Field{
			Name: "username",
			Type: FieldTypeString,
		},
		Operations: []string{OperationCreate, OperationRead},
	}

	result := fieldWithRead.HasReadOperation()
	assert.True(t, result, "ResourceField with Read operation should return true")

	// Test ResourceField without Read operation
	fieldWithoutRead := ResourceField{
		Field: Field{
			Name: "password",
			Type: FieldTypeString,
		},
		Operations: []string{OperationCreate, OperationUpdate},
	}

	result = fieldWithoutRead.HasReadOperation()
	assert.False(t, result, "ResourceField without Read operation should return false")
}

func TestResourceFieldHasUpdateOperation(t *testing.T) {
	// Test ResourceField with Update operation
	fieldWithUpdate := ResourceField{
		Field: Field{
			Name: "email",
			Type: FieldTypeString,
		},
		Operations: []string{OperationCreate, OperationUpdate, OperationRead},
	}

	result := fieldWithUpdate.HasUpdateOperation()
	assert.True(t, result, "ResourceField with Update operation should return true")

	// Test ResourceField without Update operation
	fieldWithoutUpdate := ResourceField{
		Field: Field{
			Name: "id",
			Type: FieldTypeUUID,
		},
		Operations: []string{OperationRead},
	}

	result = fieldWithoutUpdate.HasUpdateOperation()
	assert.False(t, result, "ResourceField without Update operation should return false")
}

func TestResourceFieldHasDeleteOperation(t *testing.T) {
	// Test ResourceField with Delete operation
	fieldWithDelete := ResourceField{
		Field: Field{
			Name: "adminField",
			Type: FieldTypeString,
		},
		Operations: []string{OperationDelete, OperationRead},
	}

	result := fieldWithDelete.HasDeleteOperation()
	assert.True(t, result, "ResourceField with Delete operation should return true")

	// Test ResourceField without Delete operation
	fieldWithoutDelete := ResourceField{
		Field: Field{
			Name: "username",
			Type: FieldTypeString,
		},
		Operations: []string{OperationCreate, OperationRead},
	}

	result = fieldWithoutDelete.HasDeleteOperation()
	assert.False(t, result, "ResourceField without Delete operation should return false")
}

// Service method tests

func TestServiceIsObject(t *testing.T) {
	// Arrange
	objectName := "User"
	nonObjectType := FieldTypeString

	service := Service{
		Objects: []Object{
			{Name: objectName, Description: "User object"},
			{Name: "Company", Description: "Company object"},
		},
	}

	// Act & Assert - existing object
	result := service.IsObject(objectName)
	assert.True(t, result, "Should return true for existing object type '%s'", objectName)

	// Act & Assert - primitive type
	result = service.IsObject(nonObjectType)
	assert.False(t, result, "Should return false for primitive type '%s'", nonObjectType)

	// Act & Assert - non-existent object
	result = service.IsObject("NonExistentObject")
	assert.False(t, result, "Should return false for non-existent object type")
}

func TestServiceHasObject(t *testing.T) {
	// Arrange
	existingObjectName := "User"
	nonExistentObjectName := "NonExistent"

	service := Service{
		Objects: []Object{
			{Name: existingObjectName, Description: "User object"},
			{Name: "Company", Description: "Company object"},
		},
	}

	// Act & Assert - existing object
	result := service.HasObject(existingObjectName)
	assert.True(t, result, "Should return true for existing object '%s'", existingObjectName)

	// Act & Assert - non-existent object
	result = service.HasObject(nonExistentObjectName)
	assert.False(t, result, "Should return false for non-existent object '%s'", nonExistentObjectName)
}

func TestServiceHasEnum(t *testing.T) {
	// Arrange
	existingEnumName := "UserRole"
	nonExistentEnumName := "NonExistent"

	service := Service{
		Enums: []Enum{
			{Name: existingEnumName, Description: "User role enumeration"},
			{Name: "Status", Description: "Status enumeration"},
		},
	}

	// Act & Assert - existing enum
	result := service.HasEnum(existingEnumName)
	assert.True(t, result, "Should return true for existing enum '%s'", existingEnumName)

	// Act & Assert - non-existent enum
	result = service.HasEnum(nonExistentEnumName)
	assert.False(t, result, "Should return false for non-existent enum '%s'", nonExistentEnumName)
}

func TestServiceGetObject(t *testing.T) {
	// Arrange
	existingObjectName := "User"
	expectedDescription := "User object"
	nonExistentObjectName := "NonExistent"

	service := Service{
		Objects: []Object{
			{Name: existingObjectName, Description: expectedDescription},
			{Name: "Company", Description: "Company object"},
		},
	}

	// Act & Assert - existing object
	result := service.GetObject(existingObjectName)
	assert.NotNil(t, result, "Should return object for existing object name '%s'", existingObjectName)
	assert.Equal(t, existingObjectName, result.Name, "Returned object should have correct name")
	assert.Equal(t, expectedDescription, result.Description, "Returned object should have correct description")

	// Act & Assert - non-existent object
	result = service.GetObject(nonExistentObjectName)
	assert.Nil(t, result, "Should return nil for non-existent object '%s'", nonExistentObjectName)
}

// Object method tests

func TestObjectHasField(t *testing.T) {
	// Arrange
	existingFieldName := "id"
	nonExistentFieldName := "nonExistent"

	object := Object{
		Name: "User",
		Fields: []Field{
			{Name: existingFieldName, Type: FieldTypeUUID},
			{Name: "username", Type: FieldTypeString},
		},
	}

	// Act & Assert - existing field
	result := object.HasField(existingFieldName)
	assert.True(t, result, "Should return true for existing field '%s'", existingFieldName)

	// Act & Assert - non-existent field
	result = object.HasField(nonExistentFieldName)
	assert.False(t, result, "Should return false for non-existent field '%s'", nonExistentFieldName)
}

func TestObjectGetField(t *testing.T) {
	// Arrange
	existingFieldName := "id"
	expectedFieldType := FieldTypeUUID
	nonExistentFieldName := "nonExistent"

	object := Object{
		Name: "User",
		Fields: []Field{
			{Name: existingFieldName, Type: expectedFieldType},
			{Name: "username", Type: FieldTypeString},
		},
	}

	// Act & Assert - existing field
	result := object.GetField(existingFieldName)
	assert.NotNil(t, result, "Should return field for existing field name '%s'", existingFieldName)
	assert.Equal(t, existingFieldName, result.Name, "Returned field should have correct name")
	assert.Equal(t, expectedFieldType, result.Type, "Returned field should have correct type")

	// Act & Assert - non-existent field
	result = object.GetField(nonExistentFieldName)
	assert.Nil(t, result, "Should return nil for non-existent field '%s'", nonExistentFieldName)
}

// Factory method tests

func TestCreateLimitParam(t *testing.T) {
	// Act
	limitParam := CreateLimitParam()

	// Assert
	assert.Equal(t, listLimitParamName, limitParam.Name, "Limit parameter should have correct name")
	assert.Equal(t, listLimitParamDesc, limitParam.Description, "Limit parameter should have correct description")
	assert.Equal(t, FieldTypeInt, limitParam.Type, "Limit parameter should have Int type")
	assert.Equal(t, listLimitDefaultValue, limitParam.Default, "Limit parameter should have correct default value")
}

func TestCreateOffsetParam(t *testing.T) {
	// Act
	offsetParam := CreateOffsetParam()

	// Assert
	assert.Equal(t, listOffsetParamName, offsetParam.Name, "Offset parameter should have correct name")
	assert.Equal(t, listOffsetParamDesc, offsetParam.Description, "Offset parameter should have correct description")
	assert.Equal(t, FieldTypeInt, offsetParam.Type, "Offset parameter should have Int type")
	assert.Equal(t, listOffsetDefaultValue, offsetParam.Default, "Offset parameter should have correct default value")
}

func TestCreatePaginationField(t *testing.T) {
	// Arrange
	expectedName := paginationObjectName
	expectedDescription := "Pagination information"
	expectedType := paginationObjectName

	// Act
	paginationField := CreatePaginationField()

	// Assert
	assert.Equal(t, expectedName, paginationField.Name, "Pagination field should have correct name")
	assert.Equal(t, expectedDescription, paginationField.Description, "Pagination field should have correct description")
	assert.Equal(t, expectedType, paginationField.Type, "Pagination field should have correct type")
}

// Endpoint method tests

func TestEndpointGetFullPath(t *testing.T) {
	testCases := []struct {
		name         string
		resourceName string
		endpointPath string
		expectedPath string
	}{
		{
			name:         "endpoint with id path",
			resourceName: "Users",
			endpointPath: "/{id}",
			expectedPath: "/users/{id}",
		},
		{
			name:         "endpoint with empty path",
			resourceName: "Companies",
			endpointPath: "",
			expectedPath: "/companies",
		},
		{
			name:         "endpoint with search path",
			resourceName: "Products",
			endpointPath: "/_search",
			expectedPath: "/products/_search",
		},
		{
			name:         "multi-word resource name",
			resourceName: "UserProfiles",
			endpointPath: "/{id}/avatar",
			expectedPath: "/userprofiles/{id}/avatar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := Endpoint{
				Name: "TestEndpoint",
				Path: tc.endpointPath,
			}

			result := endpoint.GetFullPath(tc.resourceName)
			assert.Equal(t, tc.expectedPath, result, "Full path should match expected for resource '%s' and endpoint path '%s'", tc.resourceName, tc.endpointPath)
		})
	}
}

// EndpointRequest method tests

func TestEndpointRequestGetRequiredBodyParams(t *testing.T) {
	// Arrange
	service := &Service{
		Objects: []Object{
			{Name: "Address", Fields: []Field{{Name: "street", Type: FieldTypeString}}},
		},
	}

	endpointRequest := EndpointRequest{
		BodyParams: []Field{
			{
				Name: "username",
				Type: FieldTypeString,
				// Required: no nullable, no array, no default, not object type
			},
			{
				Name:      "description",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierNullable}, // Not required
			},
			{
				Name:    "status",
				Type:    FieldTypeString,
				Default: "active", // Not required
			},
			{
				Name: "address",
				Type: "Address", // Object type - not required
			},
			{
				Name:      "tags",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierArray}, // Not required
			},
		},
	}

	// Act
	requiredParams := endpointRequest.GetRequiredBodyParams(service)

	// Assert
	expectedRequiredParams := []string{"username"}
	assert.Equal(t, expectedRequiredParams, requiredParams, "Should return only required body parameter names")
	assert.Len(t, requiredParams, 1, "Should return exactly 1 required parameter")
	assert.Contains(t, requiredParams, "username", "Should contain 'username' as required parameter")
}

// Utility function tests

func TestCamelCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"user_name", "userName"},
		{"first_name", "firstName"},
		{"id", "id"},
		{"created_at", "createdAt"},
		{"user_id", "userID"},
		{"api_key", "apiKey"},
		{"username", "username"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := camelCase(tc.input)
			assert.Equal(t, tc.expected, result, "CamelCase conversion for '%s' should be '%s'", tc.input, tc.expected)
		})
	}
}

func TestToKebabCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"UserProfile", "userprofile"},
		{"user_profile", "user-profile"},
		{"User Profile", "user-profile"},
		{"UserAPI", "userapi"},
		{"API_KEY", "api-key"},
		{"simple", "simple"},
		{"Multi Word String", "multi-word-string"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := toKebabCase(tc.input)
			assert.Equal(t, tc.expected, result, "KebabCase conversion for '%s' should be '%s'", tc.input, tc.expected)
		})
	}
}

func TestGetComment(t *testing.T) {
	testCases := []struct {
		name        string
		tabs        string
		description string
		fieldName   string
		expected    string
	}{
		{
			name:        "simple comment with tabs",
			tabs:        "\t",
			description: "User's username",
			fieldName:   "username",
			expected:    "\t// username: User's username",
		},
		{
			name:        "comment already prefixed with field name",
			tabs:        "\t\t",
			description: "id is the unique identifier",
			fieldName:   "id",
			expected:    "\t\t// id is the unique identifier",
		},
		{
			name:        "multiline description",
			tabs:        "\t",
			description: "First line\nSecond line",
			fieldName:   "field",
			expected:    "\t// field: First line\n\t// Second line",
		},
		{
			name:        "no tabs",
			tabs:        "",
			description: "Simple description",
			fieldName:   "name",
			expected:    "// name: Simple description",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getComment(tc.tabs, tc.description, tc.fieldName)
			assert.Equal(t, tc.expected, result, "Comment formatting should match expected for case '%s'", tc.name)
		})
	}
}

func TestFieldGetComment(t *testing.T) {
	// Test Field.GetComment method which uses the getComment helper
	testCases := []struct {
		name     string
		field    Field
		tabs     string
		expected string
	}{
		{
			name: "basic field comment",
			field: Field{
				Name:        "username",
				Description: "User's username",
			},
			tabs:     "\t",
			expected: "\t// username: User's username",
		},
		{
			name: "field with description starting with field name",
			field: Field{
				Name:        "id",
				Description: "id is the unique identifier",
			},
			tabs:     "\t\t",
			expected: "\t\t// id is the unique identifier",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.field.GetComment(tc.tabs)
			assert.Equal(t, tc.expected, result, "Field comment should match expected for case '%s'", tc.name)
		})
	}
}

// Additional factory method tests

func TestCreateDataField(t *testing.T) {
	testCases := []struct {
		resourceName        string
		expectedName        string
		expectedDescription string
		expectedType        string
		expectedModifiers   []string
	}{
		{
			resourceName:        "User",
			expectedName:        "data",
			expectedDescription: "Array of User objects",
			expectedType:        "User",
			expectedModifiers:   []string{ModifierArray},
		},
		{
			resourceName:        "Product",
			expectedName:        "data",
			expectedDescription: "Array of Product objects",
			expectedType:        "Product",
			expectedModifiers:   []string{ModifierArray},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceName, func(t *testing.T) {
			dataField := CreateDataField(tc.resourceName)

			assert.Equal(t, tc.expectedName, dataField.Name, "Data field should have correct name")
			assert.Equal(t, tc.expectedDescription, dataField.Description, "Data field should have correct description")
			assert.Equal(t, tc.expectedType, dataField.Type, "Data field should have correct type")
			assert.Equal(t, tc.expectedModifiers, dataField.Modifiers, "Data field should have correct modifiers")
			assert.True(t, dataField.IsArray(), "Data field should be an array")
		})
	}
}

func TestCreateIDParam(t *testing.T) {
	testCases := []struct {
		name         string
		description  string
		expectedName string
		expectedType string
	}{
		{
			name:         "user ID parameter",
			description:  "The unique identifier of the user",
			expectedName: "id",
			expectedType: FieldTypeUUID,
		},
		{
			name:         "product ID parameter",
			description:  "Product identifier",
			expectedName: "id",
			expectedType: FieldTypeUUID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			idParam := CreateIDParam(tc.description)

			assert.Equal(t, tc.expectedName, idParam.Name, "ID parameter should have correct name")
			assert.Equal(t, tc.description, idParam.Description, "ID parameter should have correct description")
			assert.Equal(t, tc.expectedType, idParam.Type, "ID parameter should have UUID type")
			assert.Empty(t, idParam.Modifiers, "ID parameter should have no modifiers")
			assert.Empty(t, idParam.Default, "ID parameter should have no default value")
			assert.Empty(t, idParam.Example, "ID parameter should have no example")
		})
	}
}
