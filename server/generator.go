package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

// Generator handles the generation of Go server code from OpenAPI specifications.
type Generator struct {
	config Config
}

// New creates a new Generator instance with the provided configuration.
// The configuration is validated and defaults are applied if needed.
func New(config Config) (*Generator, error) {
	// Set defaults for empty fields
	config.SetDefaults()

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Generator{
		config: config,
	}, nil
}

// GenerateFromFile generates Go server code from an OpenAPI specification file.
// The file can be in YAML or JSON format.
func (g *Generator) GenerateFromFile(specPath string) error {
	return g.GenerateFromFileWithContext(context.Background(), specPath)
}

// GenerateFromFileWithContext generates Go server code from an OpenAPI specification file with context.
// The file can be in YAML or JSON format.
func (g *Generator) GenerateFromFileWithContext(ctx context.Context, specPath string) error {
	// Read the OpenAPI specification file
	specData, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("%s: %w", errorSpecRead, err)
	}

	return g.GenerateFromDataWithContext(ctx, specData)
}

// GenerateFromReader generates Go server code from an OpenAPI specification reader.
func (g *Generator) GenerateFromReader(reader io.Reader) error {
	return g.GenerateFromReaderWithContext(context.Background(), reader)
}

// GenerateFromReaderWithContext generates Go server code from an OpenAPI specification reader with context.
func (g *Generator) GenerateFromReaderWithContext(ctx context.Context, reader io.Reader) error {
	// Read all data from the reader
	specData, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("%s: %w", errorSpecRead, err)
	}

	return g.GenerateFromDataWithContext(ctx, specData)
}

// GenerateFromData generates Go server code from OpenAPI specification data.
func (g *Generator) GenerateFromData(specData []byte) error {
	return g.GenerateFromDataWithContext(context.Background(), specData)
}

// GenerateFromDataWithContext generates Go server code from OpenAPI specification data with context.
func (g *Generator) GenerateFromDataWithContext(ctx context.Context, specData []byte) error {
	// Parse the OpenAPI specification
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData(specData)
	if err != nil {
		return fmt.Errorf("%s: %w", errorSpecRead, err)
	}

	// Validate the specification
	if err := spec.Validate(ctx); err != nil {
		return fmt.Errorf("invalid OpenAPI specification: %w", err)
	}

	// Create the codegen configuration with hardcoded defaults for Gin and strict mode
	config := codegen.Configuration{
		PackageName: g.config.PackageName,
		Generate: codegen.GenerateOptions{
			EmbeddedSpec:  defaultEmbedSpec,
			EchoServer:    defaultEchoServer,
			GinServer:     defaultGinServer,
			GorillaServer: defaultGorillaServer,
			Strict:        defaultStrictServer,
			Models:        defaultModels,
		},
		OutputOptions: codegen.OutputOptions{
			SkipFmt:   defaultSkipFmt,
			SkipPrune: defaultSkipPrune,
		},
	}

	// Generate the code
	code, err := codegen.Generate(spec, config)
	if err != nil {
		return fmt.Errorf("%s: %w", errorCodegenFailed, err)
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(g.config.OutputPath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write the generated code to the output file
	if err := os.WriteFile(g.config.OutputPath, []byte(code), 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	return nil
}

// GetConfig returns a copy of the generator configuration.
func (g *Generator) GetConfig() Config {
	return g.config
}
