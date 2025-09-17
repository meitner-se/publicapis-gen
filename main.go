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

	yaml "github.com/goccy/go-yaml"
	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/meitner-se/publicapis-gen/specification/openapi"
	"github.com/meitner-se/publicapis-gen/specification/schema"
	"github.com/meitner-se/publicapis-gen/specification/server"
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
	modeServer  = "server"
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
	suffixServer  = "-server"
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
	usageExample     = "\nExamples:\n  # Using command line flags (legacy)\n  publicapis-gen generate -file=spec.yaml -mode=overlay\n  publicapis-gen generate -file=spec.json -mode=openapi\n  publicapis-gen generate -file=spec.yaml -mode=schema\n  publicapis-gen generate -file=spec.yaml -mode=server\n  publicapis-gen generate -file=spec.yaml -mode=openapi -output=api-spec.json\n  publicapis-gen generate -file=spec.yaml -mode=schema -output=schemas.json\n  publicapis-gen generate -file=spec.yaml -mode=server -output=server.go\n  publicapis-gen generate -file=spec.yaml -mode=openapi -log-level=info\n\n  # Using config file (recommended)\n  publicapis-gen generate -config=build-config.yaml\n  publicapis-gen generate -config=build-config.yaml -log-level=info\n\n  # Using default config file (automatically detects publicapis.yaml or publicapis.yml)\n  publicapis-gen generate\n  publicapis-gen generate -log-level=info"
)

// Config file constants
const (
	configFileFlag     = "config"
	errorInvalidConfig = "invalid config file"
	errorConfigParsing = "failed to parse config file"
	defaultConfigYAML  = "publicapis.yaml"
	defaultConfigYML   = "publicapis.yml"
)

// Command constants
const (
	commandGenerate     = "generate"
	commandHelp         = "help"
	errorInvalidCommand = "invalid command"
	errorMissingCommand = "missing command"
)

// Command usage messages
const (
	mainUsageDescription     = "publicapis-gen - Generate API specifications and OpenAPI documents"
	mainUsageCommands        = "\nAvailable Commands:\n  generate    Generate API specifications and OpenAPI documents\n  help        Show help for commands\n\nUse \"publicapis-gen [command] --help\" for more information about a command."
	generateUsageDescription = "Generate API specifications and OpenAPI documents from specification files"
)

// Job represents a single generation job in the config file
type Job struct {
	Specification string `yaml:"specification" json:"specification"`
	OpenAPIJSON   string `yaml:"openapi_json,omitempty" json:"openapi_json,omitempty"`
	OpenAPIYAML   string `yaml:"openapi_yaml,omitempty" json:"openapi_yaml,omitempty"`
	SchemaJSON    string `yaml:"schema_json,omitempty" json:"schema_json,omitempty"`
	OverlayYAML   string `yaml:"overlay_yaml,omitempty" json:"overlay_yaml,omitempty"`
	OverlayJSON   string `yaml:"overlay_json,omitempty" json:"overlay_json,omitempty"`
	ServerGo      string `yaml:"server_go,omitempty" json:"server_go,omitempty"`
}

// Config represents the configuration file structure
type Config []Job

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, errorFailedToRun, logKeyError, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Parse command line to get the subcommand
	if len(os.Args) < 2 {
		showMainUsage()
		return fmt.Errorf("%s: must specify a command", errorMissingCommand)
	}

	command := os.Args[1]

	switch command {
	case commandGenerate:
		return runGenerateCommand(ctx, os.Args[2:])
	case commandHelp:
		if len(os.Args) >= 3 {
			return showCommandHelp(os.Args[2])
		}
		showMainUsage()
		return nil
	case "-help", "--help", "-h":
		showMainUsage()
		return nil
	default:
		showMainUsage()
		return fmt.Errorf("%s: unknown command '%s'", errorInvalidCommand, command)
	}
}

func showMainUsage() {
	fmt.Fprintf(os.Stderr, "%s\n\n", mainUsageDescription)
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s\n", mainUsageCommands)
}

func showCommandHelp(command string) error {
	switch command {
	case commandGenerate:
		showGenerateUsage()
		return nil
	default:
		showMainUsage()
		return fmt.Errorf("%s: unknown command '%s'", errorInvalidCommand, command)
	}
}

func showGenerateUsage() {
	fmt.Fprintf(os.Stderr, "%s\n\n", generateUsageDescription)
	fmt.Fprintf(os.Stderr, "Usage: %s generate [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -file string\n        Path to input specification file (YAML or JSON) - legacy mode\n")
	fmt.Fprintf(os.Stderr, "  -mode string\n        Operation mode: 'overlay', 'openapi', 'schema', or 'server' - legacy mode\n")
	fmt.Fprintf(os.Stderr, "  -output string\n        Output file path (optional, defaults to input name with suffix) - legacy mode\n")
	fmt.Fprintf(os.Stderr, "  -config string\n        Path to YAML config file containing multiple jobs\n")
	fmt.Fprintf(os.Stderr, "  -log-level string\n        Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)\n")
	fmt.Fprintf(os.Stderr, "  -help\n        Show this help message\n")
	fmt.Fprintf(os.Stderr, "%s\n", usageExample)
}

func runGenerateCommand(ctx context.Context, args []string) error {
	// Create a new FlagSet for the generate command
	generateFlags := flag.NewFlagSet(commandGenerate, flag.ContinueOnError)
	generateFlags.Usage = showGenerateUsage

	// Parse command line flags for generate command
	var (
		fileFlag     = generateFlags.String("file", "", "Path to input specification file (YAML or JSON) - legacy mode")
		modeFlag     = generateFlags.String("mode", "", "Operation mode: 'overlay', 'openapi', or 'schema' - legacy mode")
		outputFlag   = generateFlags.String("output", "", "Output file path (optional, defaults to input name with suffix) - legacy mode")
		configFlag   = generateFlags.String(configFileFlag, "", "Path to YAML config file containing multiple jobs")
		logLevelFlag = generateFlags.String("log-level", logLevelOff, "Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)")
		helpFlag     = generateFlags.Bool("help", false, "Show help message")
	)

	if err := generateFlags.Parse(args); err != nil {
		return err
	}

	// Configure logging
	if err := configureLogging(*logLevelFlag); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	// Show help if requested
	if *helpFlag {
		showGenerateUsage()
		return nil
	}

	// Determine if using config file or legacy mode
	usingConfig := *configFlag != ""
	usingLegacy := *fileFlag != "" || *modeFlag != ""

	if usingConfig && usingLegacy {
		showGenerateUsage()
		return fmt.Errorf("%s: cannot use both config file and legacy flags (file/mode) at the same time", errorInvalidConfig)
	}

	// If no explicit config or legacy flags provided, try to find default config file
	if !usingConfig && !usingLegacy {
		defaultConfigPath := findDefaultConfigFile()
		if defaultConfigPath != "" {
			slog.InfoContext(ctx, "Using default config file", logKeyFile, defaultConfigPath)
			return runConfigMode(ctx, defaultConfigPath)
		}

		// No default config file found, require explicit configuration
		showGenerateUsage()
		return fmt.Errorf("%s: either config file or legacy flags (file and mode) are required", errorInvalidFile)
	}

	if usingConfig {
		return runConfigMode(ctx, *configFlag)
	}

	// Legacy mode - validate required flags
	if *fileFlag == "" {
		showGenerateUsage()
		return fmt.Errorf("%s: file flag is required", errorInvalidFile)
	}

	if *modeFlag == "" {
		showGenerateUsage()
		return fmt.Errorf("%s: mode flag is required", errorInvalidMode)
	}

	// Validate mode
	if *modeFlag != modeOverlay && *modeFlag != modeOpenAPI && *modeFlag != modeSchema && *modeFlag != modeServer {
		return fmt.Errorf("%s: mode must be '%s', '%s', '%s', or '%s'", errorInvalidMode, modeOverlay, modeOpenAPI, modeSchema, modeServer)
	}

	return runLegacyMode(ctx, *fileFlag, *modeFlag, *outputFlag)
}

// runConfigMode processes jobs from a config file
func runConfigMode(ctx context.Context, configPath string) error {
	// Parse config file
	config, err := parseConfigFile(configPath)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully parsed config file", logKeyFile, configPath)

	// Process each job in the config
	for i, job := range config {
		slog.InfoContext(ctx, "Processing job", "job_index", i+1, "specification", job.Specification)

		if err := processJob(ctx, job); err != nil {
			return fmt.Errorf("failed to process job %d (spec: %s): %w", i+1, job.Specification, err)
		}
	}

	slog.InfoContext(ctx, "Successfully processed all jobs", "total_jobs", len(config))
	fmt.Printf("Successfully processed %d jobs from config file: %s\n", len(config), configPath)

	return nil
}

// runLegacyMode processes a single job using the legacy command line flags
func runLegacyMode(ctx context.Context, filePath, mode, outputPath string) error {
	// Read and parse input file
	service, err := readSpecificationFile(filePath)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully parsed input file", logKeyFile, filePath)

	// Execute the requested operation
	switch mode {
	case modeOverlay:
		return generateOverlay(ctx, service, filePath, outputPath)
	case modeOpenAPI:
		return generateOpenAPI(ctx, service, filePath, outputPath)
	case modeSchema:
		return generateSchema(ctx, service, filePath, outputPath)
	case modeServer:
		return generateServer(ctx, service, filePath, outputPath)
	default:
		return fmt.Errorf("%s: unsupported mode '%s'", errorInvalidMode, mode)
	}
}

// findDefaultConfigFile searches for default config files in the current directory
// Returns the path to the first found config file, or empty string if none found
func findDefaultConfigFile() string {
	// Check for publicapis.yaml first, then publicapis.yml
	candidates := []string{defaultConfigYAML, defaultConfigYML}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

// parseConfigFile reads and parses a YAML config file
func parseConfigFile(configPath string) (Config, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: config file does not exist: %s", errorInvalidConfig, configPath)
		}
		return nil, fmt.Errorf("%s: cannot access config file: %w", errorInvalidConfig, err)
	}

	// Read file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %w", errorConfigParsing, err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%s: failed to parse YAML: %w", errorConfigParsing, err)
	}

	// Validate config
	if len(config) == 0 {
		return nil, fmt.Errorf("%s: config file must contain at least one job", errorInvalidConfig)
	}

	// Validate each job
	for i, job := range config {
		if job.Specification == "" {
			return nil, fmt.Errorf("%s: job %d is missing required 'specification' field", errorInvalidConfig, i+1)
		}

		// Check if at least one output format is specified
		if job.OpenAPIJSON == "" && job.OpenAPIYAML == "" && job.SchemaJSON == "" && job.OverlayYAML == "" && job.OverlayJSON == "" && job.ServerGo == "" {
			return nil, fmt.Errorf("%s: job %d must specify at least one output format (openapi_json, openapi_yaml, schema_json, overlay_yaml, overlay_json, or server_go)", errorInvalidConfig, i+1)
		}
	}

	return config, nil
}

// processJob processes a single job from the config file
func processJob(ctx context.Context, job Job) error {
	// Read and parse the specification file
	service, err := readSpecificationFile(job.Specification)
	if err != nil {
		return fmt.Errorf("failed to read specification file '%s': %w", job.Specification, err)
	}

	slog.InfoContext(ctx, "Successfully parsed specification file", logKeyFile, job.Specification)

	// Generate each requested output format
	if job.OpenAPIJSON != "" {
		if err := generateOpenAPI(ctx, service, job.Specification, job.OpenAPIJSON); err != nil {
			return fmt.Errorf("failed to generate OpenAPI JSON to '%s': %w", job.OpenAPIJSON, err)
		}
	}

	if job.OpenAPIYAML != "" {
		if err := generateOpenAPIYAML(ctx, service, job.Specification, job.OpenAPIYAML); err != nil {
			return fmt.Errorf("failed to generate OpenAPI YAML to '%s': %w", job.OpenAPIYAML, err)
		}
	}

	if job.SchemaJSON != "" {
		if err := generateSchema(ctx, service, job.Specification, job.SchemaJSON); err != nil {
			return fmt.Errorf("failed to generate schema JSON to '%s': %w", job.SchemaJSON, err)
		}
	}

	if job.OverlayYAML != "" {
		if err := generateOverlay(ctx, service, job.Specification, job.OverlayYAML); err != nil {
			return fmt.Errorf("failed to generate overlay YAML to '%s': %w", job.OverlayYAML, err)
		}
	}

	if job.OverlayJSON != "" {
		if err := generateOverlay(ctx, service, job.Specification, job.OverlayJSON); err != nil {
			return fmt.Errorf("failed to generate overlay JSON to '%s': %w", job.OverlayJSON, err)
		}
	}

	if job.ServerGo != "" {
		if err := generateServer(ctx, service, job.Specification, job.ServerGo); err != nil {
			return fmt.Errorf("failed to generate server Go code to '%s': %w", job.ServerGo, err)
		}
	}

	return nil
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

// generateOpenAPIYAML generates an OpenAPI document in YAML format from the specification.
func generateOpenAPIYAML(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating OpenAPI YAML document", logKeyMode, "openapi-yaml")

	// Generate OpenAPI document as JSON first
	outputData, err := openapi.GenerateFromSpecificationToJSON(service)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	// Parse JSON to interface{} so we can convert to YAML
	var openAPIDoc interface{}
	if err := json.Unmarshal(outputData, &openAPIDoc); err != nil {
		return fmt.Errorf("failed to parse generated OpenAPI JSON: %w", err)
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(openAPIDoc)
	if err != nil {
		return fmt.Errorf("failed to convert OpenAPI document to YAML: %w", err)
	}

	// Determine output file path - ensure YAML extension
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateOpenAPIYAMLOutputPath(inputFile)
	} else {
		// Ensure output path has .yaml extension
		outputPath = ensureYAMLExtension(outputPath)
	}

	// Write output file
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated OpenAPI YAML document", logKeyFile, outputPath)
	fmt.Printf("OpenAPI YAML document generated: %s\n", outputPath)

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

// generateServer generates Go server code from the specification.
func generateServer(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating Go server code", logKeyMode, modeServer)

	// Generate OpenAPI document using the same method as standard OpenAPI generation
	openapiGenerator := &openapi.Generator{
		Version:     "3.0.3", // Use 3.0.3 for better compatibility with code generators
		Title:       service.Name + " API",
		Description: "Generated API documentation",
	}

	document, err := openapiGenerator.GenerateFromService(service)
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	// Create server generator with default configuration
	serverConfig := server.NewDefaultConfig()
	serverConfig.PackageName = "api"
	serverGenerator := server.NewGenerator(serverConfig)

	// Generate server code
	serverCode, err := serverGenerator.GenerateFromDocument(document)
	if err != nil {
		return fmt.Errorf("failed to generate server code: %w", err)
	}

	// Determine output file path - always use .go for server code
	outputPath := outputFile
	if outputPath == "" {
		outputPath = generateServerOutputPath(inputFile)
	} else {
		// Ensure output path has .go extension
		outputPath = ensureGoExtension(outputPath)
	}

	// Write output file
	if err := os.WriteFile(outputPath, serverCode, 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated Go server code", logKeyFile, outputPath)
	fmt.Printf("Go server code generated: %s\n", outputPath)

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

// generateOpenAPIYAMLOutputPath generates an output file path for OpenAPI YAML documents.
func generateOpenAPIYAMLOutputPath(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + suffixOpenAPI + extYAML
}

// generateServerOutputPath generates an output file path for Go server code (always .go).
func generateServerOutputPath(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + suffixServer + ".go"
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

// ensureYAMLExtension ensures the output path has a .yaml extension.
func ensureYAMLExtension(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext != extYAML && ext != extYML {
		base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
		return base + extYAML
	}
	return outputPath
}

// ensureGoExtension ensures the output path has a .go extension.
func ensureGoExtension(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext != ".go" {
		base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
		return base + ".go"
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
