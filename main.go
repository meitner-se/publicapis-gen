package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/meitner-se/publicapis-gen/specification/openapi"
	"github.com/meitner-se/publicapis-gen/specification/schema"
)

// Error messages and log keys
const (
	errorNotImplemented = "not implemented"
	errorFailedToRun    = "failed to run"
	errorInvalidFile    = "invalid input file"
	errorInvalidMode    = "invalid operation mode"
	errorFileWrite      = "failed to write file"
	logKeyError         = "error"
	logKeyFile          = "file"
	logKeyMode          = "mode"
)

// Operation modes
const (
	modeOverlay = "overlay"
	modeOpenAPI = "openapi"
	modeSchema  = "schema"
)

// File extensions
const (
	extYAML = ".yaml"
	extYML  = ".yml"
	extJSON = ".json"
)

// Output file suffixes
const (
	suffixOverlay = "-overlay"
	suffixOpenAPI = "-openapi"
	suffixSchema  = "-schema"
)

// Log levels
const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
	logLevelOff   = "off"
)

// Usage messages
const (
	usageDescription = "publicapis-gen - Generate API specifications and OpenAPI documents"
	usageExample     = "\nExamples:\n  publicapis-gen -file=spec.yaml -mode=overlay\n  publicapis-gen -file=spec.json -mode=openapi\n  publicapis-gen -file=spec.yaml -mode=schema\n  publicapis-gen -file=spec.yaml -mode=openapi -output=api-spec.json\n  publicapis-gen -file=spec.yaml -mode=schema -output=schemas.json\n  publicapis-gen -file=spec.yaml -mode=openapi -log-level=info"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, errorFailedToRun, logKeyError, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Parse command line flags
	var (
		fileFlag     = flag.String("file", "", "Path to input specification file (YAML or JSON)")
		modeFlag     = flag.String("mode", "", "Operation mode: 'overlay', 'openapi', or 'schema'")
		outputFlag   = flag.String("output", "", "Output file path (optional, defaults to input name with suffix)")
		logLevelFlag = flag.String("log-level", logLevelOff, "Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)")
		helpFlag     = flag.Bool("help", false, "Show help message")
	)

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n", usageDescription)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", usageExample)
	}

	flag.Parse()

	// Configure logging
	if err := configureLogging(*logLevelFlag); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	// Show help if requested
	if *helpFlag {
		flag.Usage()
		return nil
	}

	// Validate required flags
	if *fileFlag == "" {
		flag.Usage()
		return fmt.Errorf("%s: file flag is required", errorInvalidFile)
	}

	if *modeFlag == "" {
		flag.Usage()
		return fmt.Errorf("%s: mode flag is required", errorInvalidMode)
	}

	// Validate mode
	if *modeFlag != modeOverlay && *modeFlag != modeOpenAPI && *modeFlag != modeSchema {
		return fmt.Errorf("%s: mode must be '%s', '%s', or '%s'", errorInvalidMode, modeOverlay, modeOpenAPI, modeSchema)
	}

	// Read and parse input file
	service, err := readSpecificationFile(*fileFlag)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully parsed input file", logKeyFile, *fileFlag)

	// Execute the requested operation
	switch *modeFlag {
	case modeOverlay:
		return generateOverlay(ctx, service, *fileFlag, *outputFlag)
	case modeOpenAPI:
		return generateOpenAPI(ctx, service, *fileFlag, *outputFlag)
	case modeSchema:
		return generateSchema(ctx, service, *fileFlag, *outputFlag)
	default:
		return fmt.Errorf("%s: unsupported mode '%s'", errorInvalidMode, *modeFlag)
	}
}

// readSpecificationFile reads and parses a YAML or JSON specification file
// with overlays automatically applied.
func readSpecificationFile(filePath string) (*specification.Service, error) {
	return specification.ParseServiceFromFile(filePath)
}

// generateOverlay generates a specification with overlay applied.
func generateOverlay(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating specification with overlay", logKeyMode, modeOverlay)

	// Service already has overlays applied from parsing
	// This mode outputs the complete specification with overlays

	// Determine output file path
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateOutputPath(inputFile, suffixOverlay)
	}

	// Write output file
	return writeSpecificationFile(ctx, service, outputPath)
}

// generateOpenAPI generates an OpenAPI document from the specification.
func generateOpenAPI(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating OpenAPI document", logKeyMode, modeOpenAPI)

	// Generate OpenAPI document as JSON in a single call
	outputData, err := openapi.GenerateFromSpecificationToJSON(service)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	// Determine output file path - always use JSON for OpenAPI
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateOpenAPIOutputPath(inputFile)
	} else {
		// Ensure output path has .json extension
		outputPath = ensureJSONExtension(outputPath)
	}

	// Write output file
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated OpenAPI document", logKeyFile, outputPath)
	fmt.Printf("OpenAPI document generated: %s\n", outputPath)

	return nil
}

// generateSchema generates JSON schemas from the specification.
func generateSchema(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating JSON schemas", logKeyMode, modeSchema)

	// Create schema generator
	generator := schema.NewSchemaGenerator()

	// Generate all schemas
	schemas, err := generator.GenerateAllSchemas()
	if err != nil {
		return fmt.Errorf("failed to generate schemas: %w", err)
	}

	// Convert all schemas to a combined JSON structure
	schemaMap := make(map[string]interface{})
	for name, schemaObj := range schemas {
		// Convert each schema to JSON and then parse it back to interface{} for clean structure
		jsonStr, err := generator.SchemaToJSON(schemaObj)
		if err != nil {
			return fmt.Errorf("failed to convert %s schema to JSON: %w", name, err)
		}

		var schemaData interface{}
		if err := json.Unmarshal([]byte(jsonStr), &schemaData); err != nil {
			return fmt.Errorf("failed to parse %s schema JSON: %w", name, err)
		}

		schemaMap[name] = schemaData
	}

	// Marshal the combined schema map to JSON with proper indentation
	outputData, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal combined schemas to JSON: %w", err)
	}

	// Determine output file path - always use JSON for schemas
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateSchemaOutputPath(inputFile)
	} else {
		// Ensure output path has .json extension
		outputPath = ensureJSONExtension(outputPath)
	}

	// Write output file
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated JSON schemas", logKeyFile, outputPath)
	fmt.Printf("JSON schemas generated: %s\n", outputPath)

	return nil
}

// writeSpecificationFile writes a specification to a file in the appropriate format.
func writeSpecificationFile(ctx context.Context, service *specification.Service, outputPath string) error {
	// Determine output format based on extension
	ext := strings.ToLower(filepath.Ext(outputPath))
	var outputData []byte
	var err error

	switch ext {
	case extYAML, extYML:
		outputData, err = yaml.Marshal(service)
		if err != nil {
			return fmt.Errorf("failed to marshal specification to YAML: %w", err)
		}
	case extJSON:
		outputData, err = json.MarshalIndent(service, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal specification to JSON: %w", err)
		}
	default:
		// Default to YAML if extension is not recognized
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + extYAML
		outputData, err = yaml.Marshal(service)
		if err != nil {
			return fmt.Errorf("failed to marshal specification to YAML: %w", err)
		}
	}

	// Write output file
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated specification with overlay", logKeyFile, outputPath)
	fmt.Printf("Specification with overlay generated: %s\n", outputPath)

	return nil
}

// generateOutputPath generates an output file path based on input file and suffix.
func generateOutputPath(inputFile, suffix string) string {
	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)
	return base + suffix + ext
}

// generateOpenAPIOutputPath generates an output file path for OpenAPI documents (always JSON).
func generateOpenAPIOutputPath(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + suffixOpenAPI + extJSON
}

// generateSchemaOutputPath generates an output file path for JSON schema documents (always JSON).
func generateSchemaOutputPath(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + suffixSchema + extJSON
}

// ensureJSONExtension ensures the output path has a .json extension.
func ensureJSONExtension(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext != extJSON {
		base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
		return base + extJSON
	}
	return outputPath
}

// configureLogging configures the slog logger based on the specified log level.
func configureLogging(logLevel string) error {
	var level slog.Level

	switch logLevel {
	case logLevelDebug:
		level = slog.LevelDebug
	case logLevelInfo:
		level = slog.LevelInfo
	case logLevelWarn:
		level = slog.LevelWarn
	case logLevelError:
		level = slog.LevelError
	case logLevelOff:
		// Set to a very high level to suppress all logging
		level = slog.Level(1000)
	default:
		return fmt.Errorf("unsupported log level '%s', must be one of: %s, %s, %s, %s, %s",
			logLevel, logLevelDebug, logLevelInfo, logLevelWarn, logLevelError, logLevelOff)
	}

	// Create a handler with the specified level
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
