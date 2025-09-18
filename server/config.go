package server

// Constants for server generation configuration
const (
	// Default configuration values
	defaultGenerate          = true
	defaultStrictServer      = true
	defaultEmbedSpec         = true
	defaultEchoServer        = false
	defaultGinServer         = true
	defaultGorillaServer     = false
	defaultModels            = true
	defaultServerInterface   = true
	defaultSkipFmt           = false
	defaultSkipPrune         = false
	defaultAliasTypes        = false
	defaultGenerateTypeAlias = false
	defaultWithGoose         = false
	defaultResponseTypeAlias = false

	// Package names and defaults
	defaultPackageName = "api"
	defaultOutputPath  = "server.gen.go"

	// Error messages
	errorInvalidOutputPath  = "output path cannot be empty"
	errorInvalidPackageName = "package name cannot be empty"
	errorCodegenFailed      = "code generation failed"
	errorSpecRead           = "failed to read OpenAPI specification"
	errorConfigCreate       = "failed to create generator configuration"
	errorFileWrite          = "failed to write generated code"
)

// Config represents the configuration for server code generation.
// Most settings are hardcoded to sensible defaults for Gin server generation
// with strict mode enabled. Only OutputPath and PackageName are configurable.
type Config struct {
	// OutputPath specifies where the generated server code will be written.
	// This is one of the two configurable parameters.
	OutputPath string

	// PackageName specifies the package name for the generated code.
	// This is one of the two configurable parameters.
	PackageName string
}

// Validate checks if the configuration is valid and sets defaults if needed.
func (c *Config) Validate() error {
	if c.OutputPath == "" {
		return &ConfigError{Message: errorInvalidOutputPath}
	}

	if c.PackageName == "" {
		return &ConfigError{Message: errorInvalidPackageName}
	}

	return nil
}

// SetDefaults sets default values for empty fields.
func (c *Config) SetDefaults() {
	if c.OutputPath == "" {
		c.OutputPath = defaultOutputPath
	}

	if c.PackageName == "" {
		c.PackageName = defaultPackageName
	}
}

// ConfigError represents a configuration validation error.
type ConfigError struct {
	Message string
}

// Error returns the error message.
func (e *ConfigError) Error() string {
	return e.Message
}
