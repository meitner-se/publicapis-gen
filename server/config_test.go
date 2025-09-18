package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Config Tests
// ============================================================================

func TestConfig_Validate(t *testing.T) {
	// Test valid configuration
	validConfig := Config{
		OutputPath:  "api/server.go",
		PackageName: "api",
	}

	err := validConfig.Validate()
	assert.Nil(t, err, "Expected no error for valid configuration")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty output path returns validation error", func(t *testing.T) {
			invalidConfig := Config{
				OutputPath:  "",
				PackageName: "api",
			}

			err := invalidConfig.Validate()
			assert.Error(t, err, "Expected error for empty output path")
			assert.Contains(t, err.Error(), "output path cannot be empty")
		})

		t.Run("empty package name returns validation error", func(t *testing.T) {
			invalidConfig := Config{
				OutputPath:  "server.go",
				PackageName: "",
			}

			err := invalidConfig.Validate()
			assert.Error(t, err, "Expected error for empty package name")
			assert.Contains(t, err.Error(), "package name cannot be empty")
		})

		t.Run("both empty fields return validation error", func(t *testing.T) {
			invalidConfig := Config{
				OutputPath:  "",
				PackageName: "",
			}

			err := invalidConfig.Validate()
			assert.Error(t, err, "Expected error for both empty fields")
		})
	})
}

func TestConfig_SetDefaults(t *testing.T) {
	// Test setting defaults for empty configuration
	emptyConfig := Config{}
	emptyConfig.SetDefaults()

	assert.Equal(t, defaultOutputPath, emptyConfig.OutputPath, "Expected default output path")
	assert.Equal(t, defaultPackageName, emptyConfig.PackageName, "Expected default package name")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("preserves existing values", func(t *testing.T) {
			customConfig := Config{
				OutputPath:  "custom/path/server.go",
				PackageName: "customapi",
			}

			customConfig.SetDefaults()

			assert.Equal(t, "custom/path/server.go", customConfig.OutputPath, "Should preserve custom output path")
			assert.Equal(t, "customapi", customConfig.PackageName, "Should preserve custom package name")
		})

		t.Run("sets defaults for partial configuration", func(t *testing.T) {
			partialConfig := Config{
				OutputPath: "custom.go",
				// PackageName is empty
			}

			partialConfig.SetDefaults()

			assert.Equal(t, "custom.go", partialConfig.OutputPath, "Should preserve custom output path")
			assert.Equal(t, defaultPackageName, partialConfig.PackageName, "Should set default package name")
		})
	})
}

// ============================================================================
// ConfigError Tests
// ============================================================================

func TestConfigError_Error(t *testing.T) {
	expectedMessage := "test error message"
	configErr := &ConfigError{Message: expectedMessage}

	actualMessage := configErr.Error()
	assert.Equal(t, expectedMessage, actualMessage, "Error message should match expected value")

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty message", func(t *testing.T) {
			emptyErr := &ConfigError{Message: ""}
			assert.Equal(t, "", emptyErr.Error(), "Empty message should return empty string")
		})

		t.Run("multi-line message", func(t *testing.T) {
			multiLineMessage := "line 1\nline 2\nline 3"
			multiLineErr := &ConfigError{Message: multiLineMessage}
			assert.Equal(t, multiLineMessage, multiLineErr.Error(), "Multi-line message should be preserved")
		})
	})
}
