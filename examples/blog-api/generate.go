package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/meitner-se/publicapis-gen/specification"
	"github.com/meitner-se/publicapis-gen/specification/openapi"
	"github.com/meitner-se/publicapis-gen/specification/schema"
)

func main() {
	fmt.Println("🚀 Generating Blog API documentation and schemas...")

	// Create output directories
	os.MkdirAll("generated/schemas", 0755)

	// Parse the specification
	service, err := specification.ParseServiceFromFile("blog-api.yaml")
	if err != nil {
		log.Fatal("❌ Failed to parse specification:", err)
	}

	fmt.Printf("✅ Parsed specification: %s v%s\n", service.Name, service.Version)
	fmt.Printf("📊 Generated: %d enums, %d objects, %d resources\n",
		len(service.Enums), len(service.Objects), len(service.Resources))

	// Generate OpenAPI 3.1 specification
	generator := openapi.NewGenerator()
	generator.SetContactInfo("API Team", "api@myblog.com", "https://myblog.com/contact")
	generator.SetLicenseInfo("MIT", "https://opensource.org/licenses/MIT")
	generator.AddTag("posts", "Blog post management")
	generator.AddTag("authors", "Author management")
	generator.AddTag("categories", "Category management")

	document, err := generator.GenerateFromService(service)
	if err != nil {
		log.Fatal("❌ Failed to generate OpenAPI:", err)
	}

	// Save OpenAPI as YAML
	yamlBytes, err := generator.ToYAML(document)
	if err != nil {
		log.Fatal("❌ Failed to convert to YAML:", err)
	}

	err = os.WriteFile("generated/openapi.yaml", yamlBytes, 0644)
	if err != nil {
		log.Fatal("❌ Failed to save OpenAPI YAML:", err)
	}

	// Save OpenAPI as JSON
	jsonBytes, err := generator.ToJSON(document)
	if err != nil {
		log.Fatal("❌ Failed to convert to JSON:", err)
	}

	err = os.WriteFile("generated/openapi.json", jsonBytes, 0644)
	if err != nil {
		log.Fatal("❌ Failed to save OpenAPI JSON:", err)
	}

	fmt.Println("📋 Generated OpenAPI 3.1 specification:")
	fmt.Println("   • generated/openapi.yaml")
	fmt.Println("   • generated/openapi.json")

	// Generate JSON schemas
	schemaGenerator := schema.NewSchemaGenerator()
	schemas, err := schemaGenerator.GenerateAllSchemas()
	if err != nil {
		log.Fatal("❌ Failed to generate schemas:", err)
	}

	fmt.Printf("📄 Generated %d JSON schemas:\n", len(schemas))

	// Save each schema to a file
	for name, schema := range schemas {
		schemaJSON, err := schemaGenerator.SchemaToJSON(schema)
		if err != nil {
			log.Printf("⚠️  Failed to convert schema %s to JSON: %v", name, err)
			continue
		}

		filename := fmt.Sprintf("generated/schemas/%s.json", name)
		err = os.WriteFile(filename, []byte(schemaJSON), 0644)
		if err != nil {
			log.Printf("⚠️  Failed to save schema %s: %v", name, err)
			continue
		}

		fmt.Printf("   • %s\n", filename)
	}

	// Generate summary statistics
	stats := generateStats(service)
	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	os.WriteFile("generated/stats.json", statsJSON, 0644)

	fmt.Println("\n📈 API Statistics:")
	fmt.Printf("   • Total endpoints: %d\n", stats.TotalEndpoints)
	fmt.Printf("   • Filter objects: %d\n", stats.FilterObjects)
	fmt.Printf("   • Request validation objects: %d\n", stats.ValidationObjects)
	
	printEndpointSummary(service)

	fmt.Println("\n✅ Generation complete! Next steps:")
	fmt.Println("   • Run 'go run main.go' to start the API server")
	fmt.Println("   • Visit http://localhost:8080/docs for interactive documentation")
	fmt.Println("   • Test endpoints with curl or your favorite HTTP client")
}

func generateStats(service *specification.Service) map[string]interface{} {
	totalEndpoints := 0
	filterObjects := 0
	validationObjects := 0

	for _, resource := range service.Resources {
		totalEndpoints += len(resource.Endpoints)
	}

	for _, obj := range service.Objects {
		if contains(obj.Name, "Filter") {
			filterObjects++
		}
		if contains(obj.Name, "RequestError") {
			validationObjects++
		}
	}

	return map[string]interface{}{
		"TotalEndpoints":     totalEndpoints,
		"FilterObjects":      filterObjects,
		"ValidationObjects":  validationObjects,
		"ResourceCount":      len(service.Resources),
		"ObjectCount":        len(service.Objects),
		"EnumCount":         len(service.Enums),
	}
}

func printEndpointSummary(service *specification.Service) {
	fmt.Println("\n🔗 Generated Endpoints:")
	
	for _, resource := range service.Resources {
		fmt.Printf("\n   %s Resource:\n", resource.Name)
		
		for _, endpoint := range resource.Endpoints {
			path := endpoint.Path
			if path == "" {
				path = "/"
			}
			fmt.Printf("      %s /%s%s - %s\n", 
				endpoint.Method, 
				toKebabCase(resource.Name),
				path,
				endpoint.Description)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func toKebabCase(s string) string {
	// Simple conversion - in real code, use the one from specification package
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "-"
		}
		result += string(r + 32) // Convert to lowercase
	}
	return result
}