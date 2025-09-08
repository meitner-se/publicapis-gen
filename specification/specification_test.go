package specification

import (
	"encoding/json"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Service Tests
// ============================================================================

func TestService(t *testing.T) {
	t.Run("JSON marshaling and unmarshaling", func(t *testing.T) {
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
	})

	t.Run("YAML marshaling and unmarshaling", func(t *testing.T) {
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
					Operations:  []string{"Create", "Read"},
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
	})

	t.Run("with complex hierarchy", func(t *testing.T) {
		service := Service{
			Name: "ComplexService",
			Enums: []Enum{
				{
					Name:        "UserRole",
					Description: "User role enumeration",
					Values: []EnumValue{
						{Name: "Admin", Description: "Administrator role"},
						{Name: "User", Description: "Regular user role"},
						{Name: "Guest", Description: "Guest user role"},
					},
				},
				{
					Name:        "Status",
					Description: "General status enumeration",
					Values: []EnumValue{
						{Name: "Active", Description: "Active status"},
						{Name: "Inactive", Description: "Inactive status"},
						{Name: "Pending", Description: "Pending status"},
						{Name: "Deleted", Description: "Deleted status"},
					},
				},
			},
			Objects: []Object{
				{
					Name:        "Address",
					Description: "Address object",
					Fields: []Field{
						{Name: "street", Type: "String", Description: "Street address"},
						{Name: "city", Type: "String", Description: "City name"},
						{Name: "zipCode", Type: "String", Description: "ZIP code"},
						{Name: "country", Type: "String", Description: "Country name"},
					},
				},
				{
					Name:        "Contact",
					Description: "Contact information object",
					Fields: []Field{
						{Name: "email", Type: "String", Description: "Email address"},
						{Name: "phone", Type: "String", Description: "Phone number"},
						{Name: "address", Type: "Address", Description: "Physical address"},
					},
				},
				{
					Name:        "User",
					Description: "User object",
					Fields: []Field{
						{Name: "id", Type: "UUID", Description: "User ID"},
						{Name: "username", Type: "String", Description: "Username"},
						{Name: "role", Type: "UserRole", Description: "User role"},
						{Name: "status", Type: "Status", Description: "User status"},
						{Name: "contact", Type: "Contact", Description: "Contact information"},
						{Name: "tags", Type: "String", Modifiers: []string{"array"}, Description: "User tags"},
						{Name: "metadata", Type: "String", Modifiers: []string{"nullable"}, Description: "Additional metadata"},
					},
				},
			},
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
								Name:        "username",
								Type:        "String",
								Description: "Username",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "role",
								Type:        "UserRole",
								Description: "User role",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "status",
								Type:        "Status",
								Description: "User status",
								Default:     "Active",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "contact",
								Type:        "Contact",
								Description: "Contact information",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "tags",
								Type:        "String",
								Modifiers:   []string{"array"},
								Description: "User tags",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
						{
							Field: Field{
								Name:        "metadata",
								Type:        "String",
								Modifiers:   []string{"nullable"},
								Description: "Additional metadata",
							},
							Operations: []string{"Create", "Read", "Update"},
						},
					},
				},
			},
		}

		// Verify structure integrity
		assert.NotEmpty(t, service.Name)
		assert.Len(t, service.Enums, 2)
		assert.Len(t, service.Objects, 3)
		assert.Len(t, service.Resources, 1)

		// Verify enum structure
		userRoleEnum := service.Enums[0]
		assert.Equal(t, "UserRole", userRoleEnum.Name)
		assert.Len(t, userRoleEnum.Values, 3)

		// Verify object nesting
		userObject := service.Objects[2]
		assert.Equal(t, "User", userObject.Name)
		assert.Len(t, userObject.Fields, 7)

		// Verify resource field operations
		userResource := service.Resources[0]
		assert.Equal(t, "Users", userResource.Name)
		assert.Len(t, userResource.Fields, 7)

		// Test JSON marshaling of complex structure
		jsonData, err := json.Marshal(service)
		require.NoError(t, err)
		assert.NotEmpty(t, jsonData)

		// Test JSON unmarshaling of complex structure
		var unmarshaledService Service
		err = json.Unmarshal(jsonData, &unmarshaledService)
		require.NoError(t, err)
		assert.Equal(t, service.Name, unmarshaledService.Name)
		assert.Equal(t, len(service.Enums), len(unmarshaledService.Enums))
		assert.Equal(t, len(service.Objects), len(unmarshaledService.Objects))
		assert.Equal(t, len(service.Resources), len(unmarshaledService.Resources))
	})
}

func TestService_IsObject(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty collections", func(t *testing.T) {
			emptyService := Service{
				Name:      "EmptyService",
				Enums:     []Enum{},
				Objects:   []Object{},
				Resources: []Resource{},
			}

			assert.False(t, emptyService.IsObject("AnyType"), "Should return false for any type when no objects")
		})

		t.Run("nil collections", func(t *testing.T) {
			nilService := Service{
				Name:      "NilService",
				Enums:     nil,
				Objects:   nil,
				Resources: nil,
			}

			assert.False(t, nilService.IsObject("AnyType"), "Should return false for any type when objects is nil")
		})

		t.Run("case sensitivity", func(t *testing.T) {
			caseService := Service{
				Name: "CaseService",
				Objects: []Object{
					{Name: "User", Description: "User object"},
				},
			}

			// Should be case sensitive
			assert.True(t, caseService.IsObject("User"), "Should find exact case match")
			assert.False(t, caseService.IsObject("user"), "Should not find different case")
			assert.False(t, caseService.IsObject("USER"), "Should not find different case")
		})
	})
}

func TestService_HasObject(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty objects", func(t *testing.T) {
			emptyService := Service{
				Name:    "EmptyService",
				Objects: []Object{},
			}

			assert.False(t, emptyService.HasObject("AnyObject"), "Should return false for any object when no objects")
		})

		t.Run("duplicate object names", func(t *testing.T) {
			duplicateService := Service{
				Name: "DuplicateService",
				Objects: []Object{
					{Name: "User", Description: "First user"},
					{Name: "User", Description: "Second user"},
				},
			}

			// Should find the first occurrence
			assert.True(t, duplicateService.HasObject("User"), "Should find object with duplicate name")
		})
	})
}

func TestService_HasEnum(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty enums", func(t *testing.T) {
			emptyService := Service{
				Name:  "EmptyService",
				Enums: []Enum{},
			}

			assert.False(t, emptyService.HasEnum("AnyEnum"), "Should return false for any enum when no enums")
		})

		t.Run("duplicate enum names", func(t *testing.T) {
			duplicateService := Service{
				Name: "DuplicateService",
				Enums: []Enum{
					{Name: "Status", Description: "First status"},
					{Name: "Status", Description: "Second status"},
				},
			}

			// Should find the first occurrence
			assert.True(t, duplicateService.HasEnum("Status"), "Should find enum with duplicate name")
		})
	})
}

func TestService_GetObject(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty objects", func(t *testing.T) {
			emptyService := Service{
				Name:    "EmptyService",
				Objects: []Object{},
			}

			assert.Nil(t, emptyService.GetObject("AnyObject"), "Should return nil for any object when no objects")
		})

		t.Run("duplicate object names", func(t *testing.T) {
			duplicateService := Service{
				Name: "DuplicateService",
				Objects: []Object{
					{Name: "User", Description: "First user"},
					{Name: "User", Description: "Second user"},
				},
			}

			foundObject := duplicateService.GetObject("User")
			assert.NotNil(t, foundObject, "Should return first object with duplicate name")
			assert.Equal(t, "First user", foundObject.Description, "Should return first object with duplicate name")
		})
	})
}

// ============================================================================
// Enum Tests
// ============================================================================

func TestEnum(t *testing.T) {
	t.Run("structure", func(t *testing.T) {
		enum := Enum{
			Name:        "UserRole",
			Description: "User role enumeration",
			Values: []EnumValue{
				{Name: "Admin", Description: "Administrator role"},
				{Name: "User", Description: "Regular user role"},
				{Name: "Guest", Description: "Guest user role"},
			},
		}

		assert.Equal(t, "UserRole", enum.Name)
		assert.Equal(t, "User role enumeration", enum.Description)
		assert.Len(t, enum.Values, 3)

		// Check enum values
		assert.Equal(t, "Admin", enum.Values[0].Name)
		assert.Equal(t, "Administrator role", enum.Values[0].Description)
		assert.Equal(t, "User", enum.Values[1].Name)
		assert.Equal(t, "Regular user role", enum.Values[1].Description)
		assert.Equal(t, "Guest", enum.Values[2].Name)
		assert.Equal(t, "Guest user role", enum.Values[2].Description)
	})
}

// ============================================================================
// Object Tests
// ============================================================================

func TestObject_HasField(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty fields", func(t *testing.T) {
			emptyObject := Object{
				Name:        "EmptyObject",
				Description: "Object with no fields",
				Fields:      []Field{},
			}

			assert.False(t, emptyObject.HasField("anyField"), "Should return false for any field when no fields")
		})

		t.Run("nil fields", func(t *testing.T) {
			nilFieldsObject := Object{
				Name:        "NilFieldsObject",
				Description: "Object with nil fields",
				Fields:      nil,
			}

			assert.False(t, nilFieldsObject.HasField("anyField"), "Should return false for any field when fields is nil")
		})

		t.Run("case sensitivity", func(t *testing.T) {
			caseObject := Object{
				Name: "CaseObject",
				Fields: []Field{
					{Name: "userName", Type: FieldTypeString, Description: "User name"},
				},
			}

			// Should be case sensitive
			assert.True(t, caseObject.HasField("userName"), "Should find exact case match")
			assert.False(t, caseObject.HasField("username"), "Should not find different case")
			assert.False(t, caseObject.HasField("USERNAME"), "Should not find different case")
			assert.False(t, caseObject.HasField("UserName"), "Should not find different case")
		})
	})
}

func TestObject_GetField(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty fields", func(t *testing.T) {
			emptyObject := Object{
				Name:        "EmptyObject",
				Description: "Object with no fields",
				Fields:      []Field{},
			}

			assert.Nil(t, emptyObject.GetField("anyField"), "Should return nil for any field when no fields")
		})

		t.Run("duplicate field names", func(t *testing.T) {
			duplicateObject := Object{
				Name:        "DuplicateFieldsObject",
				Description: "Object with duplicate field names",
				Fields: []Field{
					{Name: "id", Type: FieldTypeUUID, Description: "First ID"},
					{Name: "id", Type: FieldTypeString, Description: "Second ID"},
				},
			}

			// Should find the first occurrence
			assert.True(t, duplicateObject.HasField("id"), "Should find field with duplicate name")

			foundField := duplicateObject.GetField("id")
			assert.NotNil(t, foundField, "Should return first field with duplicate name")
			assert.Equal(t, FieldTypeUUID, foundField.Type, "Should return first field with duplicate name")
			assert.Equal(t, "First ID", foundField.Description, "Should return first field with duplicate name")
		})
	})
}

// ============================================================================
// Field Tests
// ============================================================================

func TestField_IsArray(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("multiple modifiers including array", func(t *testing.T) {
			field := Field{
				Name:      "multiModifierField",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierArray, ModifierNullable},
			}

			assert.True(t, field.IsArray(), "Should be array with both modifiers")
		})

		t.Run("duplicate modifiers", func(t *testing.T) {
			field := Field{
				Name:      "duplicateModifierField",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierArray, ModifierArray, ModifierNullable, ModifierNullable},
			}

			assert.True(t, field.IsArray(), "Should handle duplicate array modifiers")
		})

		t.Run("empty modifiers", func(t *testing.T) {
			field := Field{
				Name:      "noModifierField",
				Type:      FieldTypeString,
				Modifiers: []string{},
			}

			assert.False(t, field.IsArray(), "Should not be array with empty modifiers")
		})

		t.Run("nil modifiers", func(t *testing.T) {
			field := Field{
				Name:      "nilModifierField",
				Type:      FieldTypeString,
				Modifiers: nil,
			}

			assert.False(t, field.IsArray(), "Should not be array with nil modifiers")
		})
	})
}

func TestField_IsNullable(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("multiple modifiers including nullable", func(t *testing.T) {
			field := Field{
				Name:      "multiModifierField",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierArray, ModifierNullable},
			}

			assert.True(t, field.IsNullable(), "Should be nullable with both modifiers")
		})

		t.Run("duplicate modifiers", func(t *testing.T) {
			field := Field{
				Name:      "duplicateModifierField",
				Type:      FieldTypeString,
				Modifiers: []string{ModifierArray, ModifierArray, ModifierNullable, ModifierNullable},
			}

			assert.True(t, field.IsNullable(), "Should handle duplicate nullable modifiers")
		})
	})
}

func TestField_TagJSON(t *testing.T) {
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

func TestField_IsRequired(t *testing.T) {
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("field with empty vs no default", func(t *testing.T) {
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
		})

		t.Run("field with whitespace default", func(t *testing.T) {
			fieldWhitespaceDefault := Field{
				Name:    "whitespaceDefault",
				Type:    FieldTypeString,
				Default: "   ",
			}
			assert.False(t, fieldWhitespaceDefault.IsRequired(service), "Field with whitespace default should not be required")
		})
	})
}

func TestField_GetComment(t *testing.T) {
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
