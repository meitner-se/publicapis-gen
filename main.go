package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/meitner-se/publicapis-gen/specification/openapi"
)

// Error messages and log keys
const (
	errorNotImplemented    = "not implemented"
	errorFailedToRun       = "failed to run"
	errorInvalidFile       = "invalid input file"
	errorUnsupportedFormat = "unsupported file format"
	errorInvalidMode       = "invalid operation mode"
	errorFileRead          = "failed to read file"
	errorFileParse         = "failed to parse file"
	errorFileWrite         = "failed to write file"
	logKeyError            = "error"
	logKeyFile             = "file"
	logKeyMode             = "mode"
)

// Operation modes
const (
	modeOverlay = "overlay"
	modeOpenAPI = "openapi"
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
)

// Usage messages
const (
	usageDescription = "publicapis-gen - Generate API specifications and OpenAPI documents"
	usageExample     = "\nExamples:\n  publicapis-gen -file=spec.yaml -mode=overlay\n  publicapis-gen -file=spec.json -mode=openapi\n  publicapis-gen -file=spec.yaml -mode=openapi -output=api-spec.yaml"
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
		fileFlag   = flag.String("file", "", "Path to input specification file (YAML or JSON)")
		modeFlag   = flag.String("mode", "", "Operation mode: 'overlay' or 'openapi'")
		outputFlag = flag.String("output", "", "Output file path (optional, defaults to input name with suffix)")
		helpFlag   = flag.Bool("help", false, "Show help message")
	)

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n", usageDescription)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", usageExample)
	}

	flag.Parse()

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
	if *modeFlag != modeOverlay && *modeFlag != modeOpenAPI {
		return fmt.Errorf("%s: mode must be '%s' or '%s'", errorInvalidMode, modeOverlay, modeOpenAPI)
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
	default:
		return fmt.Errorf("%s: unsupported mode '%s'", errorInvalidMode, *modeFlag)
	}
}

// readSpecificationFile reads and parses a YAML or JSON specification file.
func readSpecificationFile(filePath string) (*specification.Service, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: file does not exist: %s", errorInvalidFile, filePath)
	}

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorFileRead, err)
	}

	// Determine file format by extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var service specification.Service

	switch ext {
	case extYAML, extYML:
		if err := yaml.Unmarshal(data, &service); err != nil {
			return nil, fmt.Errorf("%s: YAML parsing error: %w", errorFileParse, err)
		}
	case extJSON:
		if err := json.Unmarshal(data, &service); err != nil {
			return nil, fmt.Errorf("%s: JSON parsing error: %w", errorFileParse, err)
		}
	default:
		return nil, fmt.Errorf("%s: file must have .yaml, .yml, or .json extension", errorUnsupportedFormat)
	}

	return &service, nil
}

// generateOverlay generates a specification with overlay applied.
func generateOverlay(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating specification with overlay", logKeyMode, modeOverlay)

	// Apply overlay to the service
	overlayedService := specification.ApplyOverlay(service)
	if overlayedService == nil {
		return errors.New("failed to apply overlay to specification")
	}

	// Apply filter overlay as well
	finalService := specification.ApplyFilterOverlay(overlayedService)
	if finalService == nil {
		return errors.New("failed to apply filter overlay to specification")
	}

	// Determine output file path
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateOutputPath(inputFile, suffixOverlay)
	}

	// Write output file
	return writeSpecificationFile(ctx, finalService, outputPath)
}

// generateOpenAPI generates an OpenAPI document from the specification.
func generateOpenAPI(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating OpenAPI document", logKeyMode, modeOpenAPI)

	// Apply overlay to the service first
	overlayedService := specification.ApplyOverlay(service)
	if overlayedService == nil {
		return errors.New("failed to apply overlay to specification")
	}

	// Apply filter overlay as well
	finalService := specification.ApplyFilterOverlay(overlayedService)
	if finalService == nil {
		return errors.New("failed to apply filter overlay to specification")
	}

	// Create OpenAPI generator
	generator := openapi.NewGenerator()

	// Set basic configuration
	generator.Title = finalService.Name + " API"
	generator.Description = "Generated API documentation"

	// Add server information if available
	if len(finalService.Servers) > 0 {
		// Server information is handled by the generator from the service
	}

	// Generate OpenAPI document
	document, err := generator.GenerateFromService(finalService)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	// Determine output file path
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateOutputPath(inputFile, suffixOpenAPI)
	}

	// Determine output format based on extension
	ext := strings.ToLower(filepath.Ext(outputPath))
	var outputData []byte

	switch ext {
	case extYAML, extYML:
		outputData, err = generator.ToYAML(document)
		if err != nil {
			return fmt.Errorf("failed to convert OpenAPI document to YAML: %w", err)
		}
	case extJSON:
		outputData, err = generator.ToJSON(document)
		if err != nil {
			return fmt.Errorf("failed to convert OpenAPI document to JSON: %w", err)
		}
	default:
		// Default to YAML if extension is not recognized
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + extYAML
		outputData, err = generator.ToYAML(document)
		if err != nil {
			return fmt.Errorf("failed to convert OpenAPI document to YAML: %w", err)
		}
	}

	// Write output file
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated OpenAPI document", logKeyFile, outputPath)
	fmt.Printf("OpenAPI document generated: %s\n", outputPath)

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
