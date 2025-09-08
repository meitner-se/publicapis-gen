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

// ============================================================================
// Factory Function Tests
// ============================================================================

func TestCreateLimitParam(t *testing.T) {
	// Act
	limitParam := CreateLimitParam()

	// Assert
	assert.Equal(t, listLimitParamName, limitParam.Name, "Limit parameter should have correct name")
	assert.Equal(t, listLimitParamDesc, limitParam.Description, "Limit parameter should have correct description")
	assert.Equal(t, FieldTypeInt, limitParam.Type, "Limit parameter should have Int type")
	assert.Equal(t, listLimitDefaultValue, limitParam.Default, "Limit parameter should have correct default value")

	t.Run("consistency", func(t *testing.T) {
		// Test that factory methods always return consistent results
		limit1 := CreateLimitParam()
		limit2 := CreateLimitParam()

		assert.Equal(t, limit1, limit2, "CreateLimitParam should return consistent results")
	})
}

func TestCreateOffsetParam(t *testing.T) {
	// Act
	offsetParam := CreateOffsetParam()

	// Assert
	assert.Equal(t, listOffsetParamName, offsetParam.Name, "Offset parameter should have correct name")
	assert.Equal(t, listOffsetParamDesc, offsetParam.Description, "Offset parameter should have correct description")
	assert.Equal(t, FieldTypeInt, offsetParam.Type, "Offset parameter should have Int type")
	assert.Equal(t, listOffsetDefaultValue, offsetParam.Default, "Offset parameter should have correct default value")

	t.Run("consistency", func(t *testing.T) {
		offset1 := CreateOffsetParam()
		offset2 := CreateOffsetParam()

		assert.Equal(t, offset1, offset2, "CreateOffsetParam should return consistent results")
	})
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

	t.Run("consistency", func(t *testing.T) {
		pagination1 := CreatePaginationField()
		pagination2 := CreatePaginationField()

		assert.Equal(t, pagination1, pagination2, "CreatePaginationField should return consistent results")
	})
}

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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty string", func(t *testing.T) {
			dataField := CreateDataField("")

			assert.Equal(t, "data", dataField.Name, "Data field should have 'data' name even with empty resource name")
			assert.Equal(t, "Array of  objects", dataField.Description, "Data field should handle empty resource name in description")
			assert.Equal(t, "", dataField.Type, "Data field type should match empty resource name")
			assert.True(t, dataField.IsArray(), "Data field should always be array")
		})
	})
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

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty description", func(t *testing.T) {
			idParam := CreateIDParam("")

			assert.Equal(t, "id", idParam.Name, "ID param should always have 'id' name")
			assert.Equal(t, "", idParam.Description, "ID param should accept empty description")
			assert.Equal(t, FieldTypeUUID, idParam.Type, "ID param should always have UUID type")
		})
	})
}

// ============================================================================
// Utility Function Tests
// ============================================================================

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

	t.Run("edge cases", func(t *testing.T) {
		edgeCases := []struct {
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

		for _, tc := range edgeCases {
			t.Run("camelCase_"+tc.input, func(t *testing.T) {
				result := camelCase(tc.input)
				assert.Equal(t, tc.expected, result, "CamelCase of '%s' should be '%s'", tc.input, tc.expected)
			})
		}
	})
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

	t.Run("edge cases", func(t *testing.T) {
		edgeCases := []struct {
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

		for _, tc := range edgeCases {
			t.Run("toKebabCase_"+tc.input, func(t *testing.T) {
				result := toKebabCase(tc.input)
				assert.Equal(t, tc.expected, result, "ToKebabCase of '%s' should be '%s'", tc.input, tc.expected)
			})
		}
	})
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

	t.Run("edge cases", func(t *testing.T) {
		edgeCases := []struct {
			name        string
			tabs        string
			description string
			fieldName   string
			expected    string
		}{
			{
				name:        "empty description",
				tabs:        "\t",
				description: "",
				fieldName:   "field",
				expected:    "\t// field: ",
			},
			{
				name:        "empty field name",
				tabs:        "\t",
				description: "Some description",
				fieldName:   "",
				expected:    "\t// Some description",
			},
			{
				name:        "description with only newlines",
				tabs:        "\t",
				description: "\n\n",
				fieldName:   "field",
				expected:    "\t// field: \n\t// \n\t// ",
			},
			{
				name:        "trailing newline in description",
				tabs:        "\t",
				description: "Description with trailing newline\n",
				fieldName:   "field",
				expected:    "\t// field: Description with trailing newline\n\t// ",
			},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				result := getComment(tc.tabs, tc.description, tc.fieldName)
				assert.Equal(t, tc.expected, result, "Comment formatting should match expected for case '%s'", tc.name)
			})
		}
	})
}
