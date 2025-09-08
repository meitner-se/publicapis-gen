package specification

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Edge case and comprehensive tests for better coverage

func TestResourceEdgeCases(t *testing.T) {
	t.Run("empty resource operations", func(t *testing.T) {
		resource := Resource{
			Name:        "EmptyResource",
			Description: "Resource with no operations",
			Operations:  []string{},
		}

		assert.False(t, resource.HasCreateOperation(), "Empty operations should return false for Create")
		assert.False(t, resource.HasReadOperation(), "Empty operations should return false for Read")
		assert.False(t, resource.HasUpdateOperation(), "Empty operations should return false for Update")
		assert.False(t, resource.HasDeleteOperation(), "Empty operations should return false for Delete")
	})

	t.Run("resource with all operations", func(t *testing.T) {
		resource := Resource{
			Name:        "FullResource",
			Description: "Resource with all operations",
			Operations:  []string{OperationCreate, OperationRead, OperationUpdate, OperationDelete},
		}

		assert.True(t, resource.HasCreateOperation(), "Should have Create operation")
		assert.True(t, resource.HasReadOperation(), "Should have Read operation")
		assert.True(t, resource.HasUpdateOperation(), "Should have Update operation")
		assert.True(t, resource.HasDeleteOperation(), "Should have Delete operation")
	})

	t.Run("resource with no fields", func(t *testing.T) {
		resource := Resource{
			Name:        "NoFieldsResource",
			Description: "Resource with no fields",
			Operations:  []string{OperationCreate, OperationRead, OperationUpdate},
			Fields:      []ResourceField{},
		}

		createParams := resource.GetCreateBodyParams()
		updateParams := resource.GetUpdateBodyParams()
		readableFields := resource.GetReadableFields()

		assert.Empty(t, createParams, "Should return empty slice for create params when no fields")
		assert.Empty(t, updateParams, "Should return empty slice for update params when no fields")
		assert.Empty(t, readableFields, "Should return empty slice for readable fields when no fields")
	})

	t.Run("resource with no endpoints", func(t *testing.T) {
		resource := Resource{
			Name:        "NoEndpointsResource",
			Description: "Resource with no endpoints",
			Endpoints:   []Endpoint{},
		}

		assert.False(t, resource.HasEndpoint("AnyEndpoint"), "Should return false for any endpoint name when no endpoints")
	})

	t.Run("resource with duplicate operations", func(t *testing.T) {
		resource := Resource{
			Name:        "DuplicateOpsResource",
			Description: "Resource with duplicate operations",
			Operations:  []string{OperationCreate, OperationCreate, OperationRead, OperationRead},
		}

		// Should still work correctly with duplicates
		assert.True(t, resource.HasCreateOperation(), "Should handle duplicate Create operations")
		assert.True(t, resource.HasReadOperation(), "Should handle duplicate Read operations")
		assert.False(t, resource.HasUpdateOperation(), "Should not have Update operation")
		assert.False(t, resource.HasDeleteOperation(), "Should not have Delete operation")
	})
}

func TestFieldEdgeCases(t *testing.T) {
	t.Run("field with multiple modifiers", func(t *testing.T) {
		field := Field{
			Name:      "multiModifierField",
			Type:      FieldTypeString,
			Modifiers: []string{ModifierArray, ModifierNullable},
		}

		assert.True(t, field.IsArray(), "Should be array with both modifiers")
		assert.True(t, field.IsNullable(), "Should be nullable with both modifiers")
	})

	t.Run("field with duplicate modifiers", func(t *testing.T) {
		field := Field{
			Name:      "duplicateModifierField",
			Type:      FieldTypeString,
			Modifiers: []string{ModifierArray, ModifierArray, ModifierNullable, ModifierNullable},
		}

		assert.True(t, field.IsArray(), "Should handle duplicate array modifiers")
		assert.True(t, field.IsNullable(), "Should handle duplicate nullable modifiers")
	})

	t.Run("field with empty modifiers", func(t *testing.T) {
		field := Field{
			Name:      "noModifierField",
			Type:      FieldTypeString,
			Modifiers: []string{},
		}

		assert.False(t, field.IsArray(), "Should not be array with empty modifiers")
		assert.False(t, field.IsNullable(), "Should not be nullable with empty modifiers")
	})

	t.Run("field with nil modifiers", func(t *testing.T) {
		field := Field{
			Name:      "nilModifierField",
			Type:      FieldTypeString,
			Modifiers: nil,
		}

		assert.False(t, field.IsArray(), "Should not be array with nil modifiers")
		assert.False(t, field.IsNullable(), "Should not be nullable with nil modifiers")
	})

	t.Run("field required edge cases", func(t *testing.T) {
		service := &Service{
			Objects: []Object{
				{Name: "CustomObject", Fields: []Field{{Name: "field1", Type: FieldTypeString}}},
			},
		}

		// Test field with empty default vs no default
		fieldEmptyDefault := Field{
			Name:    "emptyDefault",
			Type:    FieldTypeString,
			Default: "",
		}
		fieldNoDefault := Field{
			Name: "noDefault",
			Type: FieldTypeString,
		}

		assert.True(t, fieldNoDefault.IsRequired(service), "Field with no default should be required")
		assert.True(t, fieldEmptyDefault.IsRequired(service), "Field with empty default should be required")

		// Test field with whitespace-only default
		fieldWhitespaceDefault := Field{
			Name:    "whitespaceDefault",
			Type:    FieldTypeString,
			Default: "   ",
		}
		assert.False(t, fieldWhitespaceDefault.IsRequired(service), "Field with whitespace default should not be required")
	})
}

func TestResourceFieldEdgeCases(t *testing.T) {
	t.Run("resource field with empty operations", func(t *testing.T) {
		resourceField := ResourceField{
			Field: Field{
				Name: "emptyOpsField",
				Type: FieldTypeString,
			},
			Operations: []string{},
		}

		assert.False(t, resourceField.HasCreateOperation(), "Empty operations should return false for Create")
		assert.False(t, resourceField.HasReadOperation(), "Empty operations should return false for Read")
		assert.False(t, resourceField.HasUpdateOperation(), "Empty operations should return false for Update")
		assert.False(t, resourceField.HasDeleteOperation(), "Empty operations should return false for Delete")
	})

	t.Run("resource field with all operations", func(t *testing.T) {
		resourceField := ResourceField{
			Field: Field{
				Name: "allOpsField",
				Type: FieldTypeString,
			},
			Operations: []string{OperationCreate, OperationRead, OperationUpdate, OperationDelete},
		}

		assert.True(t, resourceField.HasCreateOperation(), "Should have Create operation")
		assert.True(t, resourceField.HasReadOperation(), "Should have Read operation")
		assert.True(t, resourceField.HasUpdateOperation(), "Should have Update operation")
		assert.True(t, resourceField.HasDeleteOperation(), "Should have Delete operation")
	})

	t.Run("resource field with duplicate operations", func(t *testing.T) {
		resourceField := ResourceField{
			Field: Field{
				Name: "duplicateOpsField",
				Type: FieldTypeString,
			},
			Operations: []string{OperationCreate, OperationCreate, OperationRead, OperationRead},
		}

		assert.True(t, resourceField.HasCreateOperation(), "Should handle duplicate Create operations")
		assert.True(t, resourceField.HasReadOperation(), "Should handle duplicate Read operations")
		assert.False(t, resourceField.HasUpdateOperation(), "Should not have Update operation")
		assert.False(t, resourceField.HasDeleteOperation(), "Should not have Delete operation")
	})
}

func TestServiceEdgeCases(t *testing.T) {
	t.Run("service with empty collections", func(t *testing.T) {
		service := Service{
			Name:      "EmptyService",
			Enums:     []Enum{},
			Objects:   []Object{},
			Resources: []Resource{},
		}

		assert.False(t, service.IsObject("AnyType"), "Should return false for any type when no objects")
		assert.False(t, service.HasObject("AnyObject"), "Should return false for any object when no objects")
		assert.False(t, service.HasEnum("AnyEnum"), "Should return false for any enum when no enums")
		assert.Nil(t, service.GetObject("AnyObject"), "Should return nil for any object when no objects")
	})

	t.Run("service with nil collections", func(t *testing.T) {
		service := Service{
			Name:      "NilService",
			Enums:     nil,
			Objects:   nil,
			Resources: nil,
		}

		assert.False(t, service.IsObject("AnyType"), "Should return false for any type when objects is nil")
		assert.False(t, service.HasObject("AnyObject"), "Should return false for any object when objects is nil")
		assert.False(t, service.HasEnum("AnyEnum"), "Should return false for any enum when enums is nil")
		assert.Nil(t, service.GetObject("AnyObject"), "Should return nil for any object when objects is nil")
	})

	t.Run("service with duplicate object/enum names", func(t *testing.T) {
		service := Service{
			Name: "DuplicateService",
			Enums: []Enum{
				{Name: "Status", Description: "First status"},
				{Name: "Status", Description: "Second status"},
			},
			Objects: []Object{
				{Name: "User", Description: "First user"},
				{Name: "User", Description: "Second user"},
			},
		}

		// Should find the first occurrence
		assert.True(t, service.HasEnum("Status"), "Should find enum with duplicate name")
		assert.True(t, service.HasObject("User"), "Should find object with duplicate name")

		foundObject := service.GetObject("User")
		assert.NotNil(t, foundObject, "Should return first object with duplicate name")
		assert.Equal(t, "First user", foundObject.Description, "Should return first object with duplicate name")
	})

	t.Run("service case sensitivity", func(t *testing.T) {
		service := Service{
			Name: "CaseService",
			Enums: []Enum{
				{Name: "Status", Description: "Status enum"},
			},
			Objects: []Object{
				{Name: "User", Description: "User object"},
			},
		}

		// Should be case sensitive
		assert.True(t, service.HasEnum("Status"), "Should find exact case match")
		assert.False(t, service.HasEnum("status"), "Should not find different case")
		assert.False(t, service.HasEnum("STATUS"), "Should not find different case")

		assert.True(t, service.HasObject("User"), "Should find exact case match")
		assert.False(t, service.HasObject("user"), "Should not find different case")
		assert.False(t, service.HasObject("USER"), "Should not find different case")
	})
}

func TestObjectEdgeCases(t *testing.T) {
	t.Run("object with empty fields", func(t *testing.T) {
		object := Object{
			Name:        "EmptyObject",
			Description: "Object with no fields",
			Fields:      []Field{},
		}

		assert.False(t, object.HasField("anyField"), "Should return false for any field when no fields")
		assert.Nil(t, object.GetField("anyField"), "Should return nil for any field when no fields")
	})

	t.Run("object with nil fields", func(t *testing.T) {
		object := Object{
			Name:        "NilFieldsObject",
			Description: "Object with nil fields",
			Fields:      nil,
		}

		assert.False(t, object.HasField("anyField"), "Should return false for any field when fields is nil")
		assert.Nil(t, object.GetField("anyField"), "Should return nil for any field when fields is nil")
	})

	t.Run("object with duplicate field names", func(t *testing.T) {
		object := Object{
			Name:        "DuplicateFieldsObject",
			Description: "Object with duplicate field names",
			Fields: []Field{
				{Name: "id", Type: FieldTypeUUID, Description: "First ID"},
				{Name: "id", Type: FieldTypeString, Description: "Second ID"},
			},
		}

		// Should find the first occurrence
		assert.True(t, object.HasField("id"), "Should find field with duplicate name")

		foundField := object.GetField("id")
		assert.NotNil(t, foundField, "Should return first field with duplicate name")
		assert.Equal(t, FieldTypeUUID, foundField.Type, "Should return first field with duplicate name")
		assert.Equal(t, "First ID", foundField.Description, "Should return first field with duplicate name")
	})

	t.Run("object field case sensitivity", func(t *testing.T) {
		object := Object{
			Name: "CaseObject",
			Fields: []Field{
				{Name: "userName", Type: FieldTypeString, Description: "User name"},
			},
		}

		// Should be case sensitive
		assert.True(t, object.HasField("userName"), "Should find exact case match")
		assert.False(t, object.HasField("username"), "Should not find different case")
		assert.False(t, object.HasField("USERNAME"), "Should not find different case")
		assert.False(t, object.HasField("UserName"), "Should not find different case")
	})
}

func TestUtilityFunctionEdgeCases(t *testing.T) {
	t.Run("camelCase edge cases", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"", ""},
			{"a", "a"},
			{"_", ""},
			{"__", ""},
			{"a_", "a"},
			{"_a", "a"},
			{"a__b", "aB"},
			{"multiple___underscores", "multipleUnderscores"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("camelCase_%s", tc.input), func(t *testing.T) {
				result := camelCase(tc.input)
				assert.Equal(t, tc.expected, result, "CamelCase of '%s' should be '%s'", tc.input, tc.expected)
			})
		}
	})

	t.Run("toKebabCase edge cases", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"", ""},
			{"a", "a"},
			{"A", "a"},
			{"_", "-"},
			{" ", "-"},
			{"__", "--"},
			{"  ", "--"},
			{"a_b c", "a-b-c"},
			{"Multiple   Spaces", "multiple---spaces"},
			{"Mixed_Case String", "mixed-case-string"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("toKebabCase_%s", tc.input), func(t *testing.T) {
				result := toKebabCase(tc.input)
				assert.Equal(t, tc.expected, result, "ToKebabCase of '%s' should be '%s'", tc.input, tc.expected)
			})
		}
	})

	t.Run("getComment edge cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			tabs        string
			description string
			fieldName   string
			expected    string
		}{
			{
				name:        "empty description",
				tabs:        "\\t",
				description: "",
				fieldName:   "field",
				expected:    "\\t// field: ",
			},
			{
				name:        "empty field name",
				tabs:        "\\t",
				description: "Some description",
				fieldName:   "",
				expected:    "\\t// Some description",
			},
			{
				name:        "description with only newlines",
				tabs:        "\\t",
				description: "\\n\\n",
				fieldName:   "field",
				expected:    "\\t// field: \\n\\n",
			},
			{
				name:        "trailing newline in description",
				tabs:        "\\t",
				description: "Description with trailing newline\\n",
				fieldName:   "field",
				expected:    "\\t// field: Description with trailing newline\\n",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := getComment(tc.tabs, tc.description, tc.fieldName)
				assert.Equal(t, tc.expected, result, "Comment formatting should match expected for case '%s'", tc.name)
			})
		}
	})
}

func TestEndpointEdgeCases(t *testing.T) {
	t.Run("endpoint getFullPath edge cases", func(t *testing.T) {
		testCases := []struct {
			name         string
			resourceName string
			endpointPath string
			expected     string
		}{
			{
				name:         "empty resource name",
				resourceName: "",
				endpointPath: "/test",
				expected:     "//test",
			},
			{
				name:         "empty endpoint path",
				resourceName: "Users",
				endpointPath: "",
				expected:     "/users",
			},
			{
				name:         "both empty",
				resourceName: "",
				endpointPath: "",
				expected:     "/",
			},
			{
				name:         "resource with special characters",
				resourceName: "User_Profiles",
				endpointPath: "/{id}",
				expected:     "/user-profiles/{id}",
			},
			{
				name:         "path starting with slash",
				resourceName: "Users",
				endpointPath: "/additional/path",
				expected:     "/users/additional/path",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := Endpoint{
					Name: "TestEndpoint",
					Path: tc.endpointPath,
				}
				result := endpoint.GetFullPath(tc.resourceName)
				assert.Equal(t, tc.expected, result, "Full path should match expected for case '%s'", tc.name)
			})
		}
	})
}

func TestEndpointRequestEdgeCases(t *testing.T) {
	t.Run("endpoint request with no body params", func(t *testing.T) {
		service := &Service{}
		endpointRequest := EndpointRequest{
			BodyParams: []Field{},
		}

		requiredParams := endpointRequest.GetRequiredBodyParams(service)
		assert.Empty(t, requiredParams, "Should return empty slice when no body params")
	})

	t.Run("endpoint request with nil body params", func(t *testing.T) {
		service := &Service{}
		endpointRequest := EndpointRequest{
			BodyParams: nil,
		}

		requiredParams := endpointRequest.GetRequiredBodyParams(service)
		assert.Empty(t, requiredParams, "Should return empty slice when body params is nil")
	})

	t.Run("endpoint request all params not required", func(t *testing.T) {
		service := &Service{
			Objects: []Object{
				{Name: "Address", Fields: []Field{{Name: "street", Type: FieldTypeString}}},
			},
		}

		endpointRequest := EndpointRequest{
			BodyParams: []Field{
				{Name: "nullable", Type: FieldTypeString, Modifiers: []string{ModifierNullable}},
				{Name: "array", Type: FieldTypeString, Modifiers: []string{ModifierArray}},
				{Name: "withDefault", Type: FieldTypeString, Default: "default"},
				{Name: "object", Type: "Address"},
			},
		}

		requiredParams := endpointRequest.GetRequiredBodyParams(service)
		assert.Empty(t, requiredParams, "Should return empty slice when no params are required")
	})
}

func TestFactoryMethodEdgeCases(t *testing.T) {
	t.Run("factory methods consistency", func(t *testing.T) {
		// Test that factory methods always return consistent results
		limit1 := CreateLimitParam()
		limit2 := CreateLimitParam()

		assert.Equal(t, limit1, limit2, "CreateLimitParam should return consistent results")

		offset1 := CreateOffsetParam()
		offset2 := CreateOffsetParam()

		assert.Equal(t, offset1, offset2, "CreateOffsetParam should return consistent results")

		pagination1 := CreatePaginationField()
		pagination2 := CreatePaginationField()

		assert.Equal(t, pagination1, pagination2, "CreatePaginationField should return consistent results")
	})

	t.Run("createDataField with empty string", func(t *testing.T) {
		dataField := CreateDataField("")

		assert.Equal(t, "data", dataField.Name, "Data field should have 'data' name even with empty resource name")
		assert.Equal(t, "Array of  objects", dataField.Description, "Data field should handle empty resource name in description")
		assert.Equal(t, "", dataField.Type, "Data field type should match empty resource name")
		assert.True(t, dataField.IsArray(), "Data field should always be array")
	})

	t.Run("createIDParam with empty description", func(t *testing.T) {
		idParam := CreateIDParam("")

		assert.Equal(t, "id", idParam.Name, "ID param should always have 'id' name")
		assert.Equal(t, "", idParam.Description, "ID param should accept empty description")
		assert.Equal(t, FieldTypeUUID, idParam.Type, "ID param should always have UUID type")
	})
}
