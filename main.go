package main

import (
	"bytes"
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
	"github.com/meitner-se/publicapis-gen/specification/openapigen"
	"github.com/meitner-se/publicapis-gen/specification/schemagen"
	"github.com/meitner-se/publicapis-gen/specification/servergen"
	"github.com/meitner-se/publicapis-gen/specification/testgen"
)

// Error messages and log keys
const (
	errorNotImplemented = "not implemented"
	errorFailedToRun    = "failed to run"
	errorInvalidFile    = "invalid input file"
	errorInvalidMode    = "invalid operation mode"
	errorFileWrite      = "failed to write file"
	errorFileRead       = "failed to read file"
	errorFilesDiffer    = "files differ from generated content"
	logKeyError         = "error"
	logKeyFile          = "file"
	logKeyMode          = "mode"
)

// Diff output constants
const (
	diffFileNotExist    = "file does not exist on disk"
	diffContentDiffers  = "content differs"
	diffSeparator       = "---"
	diffGeneratedHeader = "Generated content:"
	diffDiskHeader      = "File on disk:"
	diffLinePrefix      = "  "
	diffMaxLinesToShow  = 10
	diffContextLines    = 3
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
	usageExample     = "\nExamples:\n  # Using config file\n  publicapis-gen generate -config=build-config.yaml\n  publicapis-gen generate -config=build-config.yaml -log-level=info\n\n  # Using default config file (automatically detects publicapis.yaml or publicapis.yml)\n  publicapis-gen generate\n  publicapis-gen generate -log-level=info"
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
	commandDiff         = "diff"
	commandHelp         = "help"
	errorInvalidCommand = "invalid command"
	errorMissingCommand = "missing command"
)

// Command usage messages
const (
	mainUsageDescription     = "publicapis-gen - Generate API specifications and OpenAPI documents"
	mainUsageCommands        = "\nAvailable Commands:\n  generate    Generate API specifications and OpenAPI documents\n  diff        Check for differences between generated files and files on disk\n  help        Show help for commands\n\nUse \"publicapis-gen [command] --help\" for more information about a command."
	generateUsageDescription = "Generate API specifications and OpenAPI documents from specification files"
	diffUsageDescription     = "Check for differences between generated files and files on disk"
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
	ServerPackage string `yaml:"server_package,omitempty" json:"server_package,omitempty"`
	TestGo        string `yaml:"test_go,omitempty" json:"test_go,omitempty"`
	TestPackage   string `yaml:"test_package,omitempty" json:"test_package,omitempty"`
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
	case commandDiff:
		return runDiffCommand(ctx, os.Args[2:])
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
	case commandDiff:
		showDiffUsage()
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
	fmt.Fprintf(os.Stderr, "  -config string\n        Path to YAML config file containing multiple jobs\n")
	fmt.Fprintf(os.Stderr, "  -log-level string\n        Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)\n")
	fmt.Fprintf(os.Stderr, "  -help\n        Show this help message\n")
	fmt.Fprintf(os.Stderr, "%s\n", usageExample)
}

func showDiffUsage() {
	fmt.Fprintf(os.Stderr, "%s\n\n", diffUsageDescription)
	fmt.Fprintf(os.Stderr, "Usage: %s diff [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -config string\n        Path to YAML config file containing multiple jobs\n")
	fmt.Fprintf(os.Stderr, "  -log-level string\n        Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)\n")
	fmt.Fprintf(os.Stderr, "  -help\n        Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  # Using config file\n")
	fmt.Fprintf(os.Stderr, "  publicapis-gen diff -config=build-config.yaml\n")
	fmt.Fprintf(os.Stderr, "  publicapis-gen diff -config=build-config.yaml -log-level=info\n\n")
	fmt.Fprintf(os.Stderr, "  # Using default config file (automatically detects publicapis.yaml or publicapis.yml)\n")
	fmt.Fprintf(os.Stderr, "  publicapis-gen diff\n")
	fmt.Fprintf(os.Stderr, "  publicapis-gen diff -log-level=info\n")
}

func runGenerateCommand(ctx context.Context, args []string) error {
	// Create a new FlagSet for the generate command
	generateFlags := flag.NewFlagSet(commandGenerate, flag.ContinueOnError)
	generateFlags.Usage = showGenerateUsage

	// Parse command line flags for generate command
	var (
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

	// Determine config file path
	configPath := *configFlag
	if configPath == "" {
		// Try to find default config file
		defaultConfigPath := findDefaultConfigFile()
		if defaultConfigPath != "" {
			slog.InfoContext(ctx, "Using default config file", logKeyFile, defaultConfigPath)
			configPath = defaultConfigPath
		} else {
			// No default config file found, require explicit configuration
			showGenerateUsage()
			return fmt.Errorf("%s: config file is required", errorInvalidConfig)
		}
	}

	return runConfigMode(ctx, configPath)
}

func runDiffCommand(ctx context.Context, args []string) error {
	// Create a new FlagSet for the diff command
	diffFlags := flag.NewFlagSet(commandDiff, flag.ContinueOnError)
	diffFlags.Usage = showDiffUsage

	// Parse command line flags for diff command
	var (
		configFlag   = diffFlags.String(configFileFlag, "", "Path to YAML config file containing multiple jobs")
		logLevelFlag = diffFlags.String("log-level", logLevelOff, "Log level: 'debug', 'info', 'warn', 'error', or 'off' (default: off)")
		helpFlag     = diffFlags.Bool("help", false, "Show help message")
	)

	if err := diffFlags.Parse(args); err != nil {
		return err
	}

	// Configure logging
	if err := configureLogging(*logLevelFlag); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	// Show help if requested
	if *helpFlag {
		showDiffUsage()
		return nil
	}

	// Determine config file path
	configPath := *configFlag
	if configPath == "" {
		// Try to find default config file
		defaultConfigPath := findDefaultConfigFile()
		if defaultConfigPath != "" {
			slog.InfoContext(ctx, "Using default config file", logKeyFile, defaultConfigPath)
			configPath = defaultConfigPath
		} else {
			// No default config file found, require explicit configuration
			showDiffUsage()
			return fmt.Errorf("%s: config file is required", errorInvalidConfig)
		}
	}

	return runDiffMode(ctx, configPath)
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
		if job.OpenAPIJSON == "" && job.OpenAPIYAML == "" && job.SchemaJSON == "" && job.OverlayYAML == "" && job.OverlayJSON == "" && job.ServerGo == "" && job.TestGo == "" {
			return nil, fmt.Errorf("%s: job %d must specify at least one output format (openapi_json, openapi_yaml, schema_json, overlay_yaml, overlay_json, server_go, or test_go)", errorInvalidConfig, i+1)
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
		// Generate server code using servergen from the specification
		if err := generateServerFromSpecification(ctx, service, job.Specification, job.ServerGo, job.ServerPackage); err != nil {
			return fmt.Errorf("failed to generate Go server to '%s': %w", job.ServerGo, err)
		}
	}

	if job.TestGo != "" {
		// Generate test code using testgen from the specification
		if err := generateTestsFromSpecification(ctx, service, job.Specification, job.TestGo, job.TestPackage); err != nil {
			return fmt.Errorf("failed to generate Go tests to '%s': %w", job.TestGo, err)
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

// generateOpenAPIBytes generates an OpenAPI document from the specification and returns it as bytes.
func generateOpenAPIBytes(ctx context.Context, service *specification.Service) ([]byte, error) {
	slog.InfoContext(ctx, "Generating OpenAPI document bytes", logKeyMode, modeOpenAPI)

	// Generate OpenAPI document as JSON using the new API pattern
	var buf bytes.Buffer
	err := openapigen.GenerateOpenAPI(&buf, service)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OpenAPI document: %w", err)
	}

	return buf.Bytes(), nil
}

// generateOpenAPI generates an OpenAPI document from the specification.
func generateOpenAPI(ctx context.Context, service *specification.Service, inputFile, outputFile string) error {
	slog.InfoContext(ctx, "Generating OpenAPI document", logKeyMode, modeOpenAPI)

	// Generate OpenAPI document as bytes
	outputData, err := generateOpenAPIBytes(ctx, service)
	if err != nil {
		return err
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

	// Generate OpenAPI document as JSON bytes first
	outputData, err := generateOpenAPIBytes(ctx, service)
	if err != nil {
		return err
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

	// Create buffer for schema generation
	var buf bytes.Buffer

	// Generate schemas using the new API
	if err := schemagen.GenerateSchemas(&buf); err != nil {
		return fmt.Errorf("failed to generate schemas: %w", err)
	}

	// The buffer now contains the JSON data
	outputData := buf.Bytes()

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

// generateOpenAPIYAMLOutputPath generates an output file path for OpenAPI YAML documents.
func generateOpenAPIYAMLOutputPath(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + suffixOpenAPI + extYAML
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

// generateServerFromSpecification generates Go server code from a specification (for config mode).
// It uses servergen to generate the server code directly from the specification.
func generateServerFromSpecification(ctx context.Context, service *specification.Service, specPath, outputPath, packageName string) error {
	slog.InfoContext(ctx, "Generating Go server code from specification using servergen", logKeyMode, modeServer)

	// Note: packageName parameter is currently not used as servergen hardcodes the package to "api"
	// This is kept for future compatibility when servergen supports custom package names

	// Generate server code using servergen
	var buf bytes.Buffer
	if err := servergen.GenerateServer(&buf, service); err != nil {
		return fmt.Errorf("failed to generate server code: %w", err)
	}

	// Write the generated code to file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated Go server code", logKeyFile, outputPath)
	fmt.Printf("Go server code generated: %s\n", outputPath)

	return nil
}

// generateTestsFromSpecification generates HTTP API tests from a service specification using testgen.
func generateTestsFromSpecification(ctx context.Context, service *specification.Service, specPath, outputPath, packageName string) error {
	slog.InfoContext(ctx, "Generating Go test code from specification using testgen", logKeyMode, "test")

	// Default package name if not provided
	if packageName == "" {
		packageName = "main"
	}

	// Generate test code using testgen
	var buf bytes.Buffer
	if err := testgen.GenerateTests(&buf, service, packageName); err != nil {
		return fmt.Errorf("failed to generate test code: %w", err)
	}

	// Write the generated code to file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("%s: %w", errorFileWrite, err)
	}

	slog.InfoContext(ctx, "Successfully generated Go test code", logKeyFile, outputPath)
	fmt.Printf("Go test code generated: %s\n", outputPath)

	return nil
}

// runDiffMode processes jobs from a config file and checks for differences
func runDiffMode(ctx context.Context, configPath string) error {
	// Parse config file
	config, err := parseConfigFile(configPath)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully parsed config file", logKeyFile, configPath)

	var hasDifferences bool
	var diffResults []string

	// Check each job in the config
	for i, job := range config {
		slog.InfoContext(ctx, "Checking job", "job_index", i+1, "specification", job.Specification)

		jobDiffs, err := checkJobDifferences(ctx, job)
		if err != nil {
			return fmt.Errorf("failed to check job %d (spec: %s): %w", i+1, job.Specification, err)
		}

		if len(jobDiffs) > 0 {
			hasDifferences = true
			diffResults = append(diffResults, fmt.Sprintf("Job %d (spec: %s):", i+1, job.Specification))
			diffResults = append(diffResults, jobDiffs...)
			diffResults = append(diffResults, "")
		}
	}

	slog.InfoContext(ctx, "Successfully checked all jobs", "total_jobs", len(config))

	if hasDifferences {
		fmt.Printf("Differences found:\n\n")
		for _, result := range diffResults {
			fmt.Println(result)
		}
		return fmt.Errorf("%s: generated content differs from files on disk", errorFilesDiffer)
	}

	fmt.Printf("No differences found. All generated files match the files on disk.\n")
	return nil
}

// checkJobDifferences checks a single job for differences between generated content and disk files
func checkJobDifferences(ctx context.Context, job Job) ([]string, error) {
	var differences []string

	// Read and parse the specification file
	service, err := readSpecificationFile(job.Specification)
	if err != nil {
		return nil, fmt.Errorf("failed to read specification file '%s': %w", job.Specification, err)
	}

	slog.InfoContext(ctx, "Successfully parsed specification file", logKeyFile, job.Specification)

	// Check OpenAPI JSON output
	if job.OpenAPIJSON != "" {
		if diff, err := checkOpenAPIJSONDifference(ctx, service, job.OpenAPIJSON); err != nil {
			return nil, fmt.Errorf("failed to check OpenAPI JSON '%s': %w", job.OpenAPIJSON, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  OpenAPI JSON (%s): %s", job.OpenAPIJSON, diff))
		}
	}

	// Check OpenAPI YAML output
	if job.OpenAPIYAML != "" {
		if diff, err := checkOpenAPIYAMLDifference(ctx, service, job.OpenAPIYAML); err != nil {
			return nil, fmt.Errorf("failed to check OpenAPI YAML '%s': %w", job.OpenAPIYAML, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  OpenAPI YAML (%s): %s", job.OpenAPIYAML, diff))
		}
	}

	// Check Schema JSON output
	if job.SchemaJSON != "" {
		if diff, err := checkSchemaJSONDifference(ctx, service, job.SchemaJSON); err != nil {
			return nil, fmt.Errorf("failed to check Schema JSON '%s': %w", job.SchemaJSON, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  Schema JSON (%s): %s", job.SchemaJSON, diff))
		}
	}

	// Check Overlay YAML output
	if job.OverlayYAML != "" {
		if diff, err := checkOverlayDifference(ctx, service, job.OverlayYAML); err != nil {
			return nil, fmt.Errorf("failed to check Overlay YAML '%s': %w", job.OverlayYAML, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  Overlay YAML (%s): %s", job.OverlayYAML, diff))
		}
	}

	// Check Overlay JSON output
	if job.OverlayJSON != "" {
		if diff, err := checkOverlayDifference(ctx, service, job.OverlayJSON); err != nil {
			return nil, fmt.Errorf("failed to check Overlay JSON '%s': %w", job.OverlayJSON, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  Overlay JSON (%s): %s", job.OverlayJSON, diff))
		}
	}

	// Check Server Go output
	if job.ServerGo != "" {
		if diff, err := checkServerGoDifference(ctx, service, job.ServerGo); err != nil {
			return nil, fmt.Errorf("failed to check Server Go '%s': %w", job.ServerGo, err)
		} else if diff != "" {
			differences = append(differences, fmt.Sprintf("  Server Go (%s): %s", job.ServerGo, diff))
		}
	}

	return differences, nil
}

// checkOpenAPIJSONDifference checks if the generated OpenAPI JSON differs from the file on disk
func checkOpenAPIJSONDifference(ctx context.Context, service *specification.Service, filePath string) (string, error) {
	// Generate OpenAPI content in memory
	generatedData, err := generateOpenAPIBytes(ctx, service)
	if err != nil {
		return "", err
	}

	return compareWithDiskFile(filePath, generatedData)
}

// checkOpenAPIYAMLDifference checks if the generated OpenAPI YAML differs from the file on disk
func checkOpenAPIYAMLDifference(ctx context.Context, service *specification.Service, filePath string) (string, error) {
	// Generate OpenAPI JSON bytes first
	jsonData, err := generateOpenAPIBytes(ctx, service)
	if err != nil {
		return "", err
	}

	// Parse JSON to interface{} so we can convert to YAML
	var openAPIDoc interface{}
	if err := json.Unmarshal(jsonData, &openAPIDoc); err != nil {
		return "", fmt.Errorf("failed to parse generated OpenAPI JSON: %w", err)
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(openAPIDoc)
	if err != nil {
		return "", fmt.Errorf("failed to convert OpenAPI document to YAML: %w", err)
	}

	return compareWithDiskFile(filePath, yamlData)
}

// checkSchemaJSONDifference checks if the generated Schema JSON differs from the file on disk
func checkSchemaJSONDifference(ctx context.Context, service *specification.Service, filePath string) (string, error) {
	// Generate schemas in memory
	var buf bytes.Buffer
	if err := schemagen.GenerateSchemas(&buf); err != nil {
		return "", fmt.Errorf("failed to generate schemas: %w", err)
	}

	return compareWithDiskFile(filePath, buf.Bytes())
}

// checkOverlayDifference checks if the generated overlay differs from the file on disk
func checkOverlayDifference(ctx context.Context, service *specification.Service, filePath string) (string, error) {
	// Determine output format based on extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var generatedData []byte
	var err error

	switch ext {
	case extYAML, extYML:
		generatedData, err = yaml.Marshal(service)
		if err != nil {
			return "", fmt.Errorf("failed to marshal specification to YAML: %w", err)
		}
	case extJSON:
		generatedData, err = json.MarshalIndent(service, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal specification to JSON: %w", err)
		}
	default:
		// Default to YAML if extension is not recognized
		generatedData, err = yaml.Marshal(service)
		if err != nil {
			return "", fmt.Errorf("failed to marshal specification to YAML: %w", err)
		}
	}

	return compareWithDiskFile(filePath, generatedData)
}

// checkServerGoDifference checks if the generated Server Go code differs from the file on disk
func checkServerGoDifference(ctx context.Context, service *specification.Service, filePath string) (string, error) {
	// Generate server code in memory
	var buf bytes.Buffer
	if err := servergen.GenerateServer(&buf, service); err != nil {
		return "", fmt.Errorf("failed to generate server code: %w", err)
	}

	return compareWithDiskFile(filePath, buf.Bytes())
}

// compareWithDiskFile compares generated content with the content of a file on disk
func compareWithDiskFile(filePath string, generatedData []byte) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return diffFileNotExist, nil
		}
		return "", fmt.Errorf("%s: cannot access file: %w", errorFileRead, err)
	}

	// Read file content from disk
	diskData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errorFileRead, err)
	}

	// Compare content
	if !bytes.Equal(generatedData, diskData) {
		diffOutput := generateDetailedDiff(generatedData, diskData)
		return fmt.Sprintf("%s\n%s", diffContentDiffers, diffOutput), nil
	}

	return "", nil
}

// generateDetailedDiff creates a detailed diff output showing line-by-line differences
func generateDetailedDiff(generatedData, diskData []byte) string {
	generatedLines := strings.Split(string(generatedData), "\n")
	diskLines := strings.Split(string(diskData), "\n")

	var diffOutput strings.Builder

	// Find the first difference
	firstDiffLine := findFirstDifference(generatedLines, diskLines)
	if firstDiffLine == -1 {
		// This shouldn't happen since we already know they differ
		return diffContentDiffers
	}

	// Show context around the first difference
	startLine := max(0, firstDiffLine-diffContextLines)
	endLine := min(len(generatedLines), len(diskLines))
	endLine = min(endLine, firstDiffLine+diffContextLines+1)

	diffOutput.WriteString(fmt.Sprintf("\nFirst difference found at line %d:\n\n", firstDiffLine+1))

	// Show generated content section
	diffOutput.WriteString(diffGeneratedHeader)
	diffOutput.WriteString("\n")
	for i := startLine; i < min(endLine, len(generatedLines)); i++ {
		marker := diffLinePrefix
		if i == firstDiffLine {
			marker = "→ "
		}
		diffOutput.WriteString(fmt.Sprintf("%s%d: %s\n", marker, i+1, generatedLines[i]))
	}

	diffOutput.WriteString("\n")
	diffOutput.WriteString(diffSeparator)
	diffOutput.WriteString("\n\n")

	// Show disk file content section
	diffOutput.WriteString(diffDiskHeader)
	diffOutput.WriteString("\n")
	for i := startLine; i < min(endLine, len(diskLines)); i++ {
		marker := diffLinePrefix
		if i == firstDiffLine {
			marker = "→ "
		}
		diffOutput.WriteString(fmt.Sprintf("%s%d: %s\n", marker, i+1, diskLines[i]))
	}

	// Show additional differences if there are more
	additionalDiffs := countAdditionalDifferences(generatedLines, diskLines, firstDiffLine+1)
	if additionalDiffs > 0 {
		diffOutput.WriteString(fmt.Sprintf("\n... and %d more line(s) differ\n", additionalDiffs))
	}

	return diffOutput.String()
}

// findFirstDifference finds the line number where files first differ
func findFirstDifference(lines1, lines2 []string) int {
	minLen := min(len(lines1), len(lines2))

	for i := 0; i < minLen; i++ {
		if lines1[i] != lines2[i] {
			return i
		}
	}

	// If one file is longer than the other
	if len(lines1) != len(lines2) {
		return minLen
	}

	return -1
}

// countAdditionalDifferences counts how many more lines differ after the first difference
func countAdditionalDifferences(lines1, lines2 []string, startFrom int) int {
	count := 0
	maxLen := max(len(lines1), len(lines2))
	minLen := min(len(lines1), len(lines2))

	// Count different lines in the overlapping part
	for i := startFrom; i < minLen; i++ {
		if lines1[i] != lines2[i] {
			count++
		}
	}

	// Count extra lines if one file is longer
	if len(lines1) != len(lines2) {
		count += maxLen - minLen
	}

	return count
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
