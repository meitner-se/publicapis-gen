package specification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// Test constants for security tests
const (
	testSecurityGroupName  = "Basic"
	testSecurityGroupName2 = "Secure"
	testClientIDName       = "ClientID"
	testClientSecretName   = "ClientSecret"
	testMTLSName           = "mTLS"
	testBearerName         = "Bearer"
	testAPIKeyType         = "apiKey"
	testMutualTLSType      = "mutualTLS"
	testHTTPType           = "http"
	testHeaderLocation     = "header"
	testBearerScheme       = "bearer"
	testJWTFormat          = "JWT"
)

// contains is a helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// TestProcessSecurity_GroupedFormat tests processing of the new grouped security format
func TestProcessSecurity_GroupedFormat(t *testing.T) {
	yamlData := `
name: Test Service
version: 1.0.0
security:
  Basic:
    - name: ClientID
      type: apiKey
      in: header
    - name: ClientSecret
      type: apiKey
      in: header
  Secure:
    - name: mTLS
      type: mutualTLS
    - name: Bearer
      type: http
      scheme: bearer
      bearerFormat: JWT
resources: []
`

	var service Service
	err := yaml.Unmarshal([]byte(yamlData), &service)
	assert.NoError(t, err, "Should parse YAML without error")

	// Process security
	requirements, err := service.ProcessSecurity()
	assert.NoError(t, err, "Should process security without error")
	assert.NotNil(t, requirements, "Should return security requirements")
	assert.Len(t, requirements, 2, "Should have 2 security requirement groups")

	// Check that security schemes were created
	assert.NotNil(t, service.SecuritySchemes, "Should create security schemes map")
	assert.Len(t, service.SecuritySchemes, 4, "Should have 4 security schemes")

	// Verify Basic group schemes
	basicClientID, exists := service.SecuritySchemes["Basic_ClientID"]
	assert.True(t, exists, "Should have Basic_ClientID scheme")
	assert.Equal(t, testAPIKeyType, basicClientID.Type)
	assert.Equal(t, testClientIDName, basicClientID.Name)
	assert.Equal(t, testHeaderLocation, basicClientID.In)

	basicClientSecret, exists := service.SecuritySchemes["Basic_ClientSecret"]
	assert.True(t, exists, "Should have Basic_ClientSecret scheme")
	assert.Equal(t, testAPIKeyType, basicClientSecret.Type)
	assert.Equal(t, testClientSecretName, basicClientSecret.Name)
	assert.Equal(t, testHeaderLocation, basicClientSecret.In)

	// Verify Secure group schemes
	secureMTLS, exists := service.SecuritySchemes["Secure_mTLS"]
	assert.True(t, exists, "Should have Secure_mTLS scheme")
	assert.Equal(t, testMutualTLSType, secureMTLS.Type)

	secureBearer, exists := service.SecuritySchemes["Secure_Bearer"]
	assert.True(t, exists, "Should have Secure_Bearer scheme")
	assert.Equal(t, testHTTPType, secureBearer.Type)
	assert.Equal(t, testBearerScheme, secureBearer.Scheme)
	assert.Equal(t, testJWTFormat, secureBearer.BearerFormat)

	// Verify security requirements - order is not guaranteed
	var foundBasic, foundSecure bool
	for _, req := range requirements {
		if len(req) == 2 && contains(req, "Basic_ClientID") && contains(req, "Basic_ClientSecret") {
			foundBasic = true
		}
		if len(req) == 2 && contains(req, "Secure_mTLS") && contains(req, "Secure_Bearer") {
			foundSecure = true
		}
	}
	assert.True(t, foundBasic, "Should have Basic security requirement")
	assert.True(t, foundSecure, "Should have Secure security requirement")
}

// TestProcessSecurity_OldFormat tests backward compatibility with the old format
func TestProcessSecurity_OldFormat(t *testing.T) {
	// Create service with old format already populated
	service := Service{
		Name:    "Test Service",
		Version: "1.0.0",
		SecuritySchemes: map[string]SecurityScheme{
			"clientId": {
				Type: "apiKey",
				In:   "header",
				Name: "X-Client-Id",
			},
			"clientSecret": {
				Type: "apiKey",
				In:   "header",
				Name: "X-Client-Secret",
			},
		},
		Security: []SecurityRequirement{
			{"clientId", "clientSecret"},
		},
		Resources: []Resource{},
	}

	// Process security
	requirements, err := service.ProcessSecurity()
	assert.NoError(t, err, "Should process security without error")
	assert.NotNil(t, requirements, "Should return security requirements")
	assert.Len(t, requirements, 1, "Should have 1 security requirement")

	// Verify old format is preserved
	assert.Contains(t, requirements[0], "clientId")
	assert.Contains(t, requirements[0], "clientSecret")
	assert.Len(t, requirements[0], 2)

	// SecuritySchemes should remain unchanged
	assert.Len(t, service.SecuritySchemes, 2, "Should keep original security schemes")
}

// TestProcessSecurity_EmptySecurity tests processing with no security defined
func TestProcessSecurity_EmptySecurity(t *testing.T) {
	service := Service{
		Name:      "Test Service",
		Version:   "1.0.0",
		Resources: []Resource{},
	}

	requirements, err := service.ProcessSecurity()
	assert.NoError(t, err, "Should process without error")
	assert.Nil(t, requirements, "Should return nil for empty security")
}

// TestProcessSecurity_InvalidGroupedFormat tests error handling for invalid format
func TestProcessSecurity_InvalidGroupedFormat(t *testing.T) {
	yamlData := `
name: Test Service
version: 1.0.0
security:
  Basic:
    - type: apiKey
      in: header
resources: []
`

	var service Service
	err := yaml.Unmarshal([]byte(yamlData), &service)
	assert.NoError(t, err, "Should parse YAML without error")

	// Process security - should fail due to missing name
	_, err = service.ProcessSecurity()
	assert.Error(t, err, "Should error when name is missing")
	assert.Contains(t, err.Error(), "must have a 'name' field")
}

// TestParseServiceFromYAML_WithGroupedSecurity tests full parsing with grouped security
func TestParseServiceFromYAML_WithGroupedSecurity(t *testing.T) {
	yamlData := `
name: Test Service
version: 1.0.0
security:
  Basic:
    - name: ClientID
      type: apiKey
      in: header
      description: Client ID for authentication
    - name: ClientSecret
      type: apiKey
      in: header
      description: Client secret for authentication
resources: []
`

	service, err := ParseServiceFromYAML([]byte(yamlData))
	assert.NoError(t, err, "Should parse service without error")
	assert.NotNil(t, service, "Should return service")

	// After parsing with overlays, security should be processed
	requirements, ok := service.Security.([]SecurityRequirement)
	assert.True(t, ok, "Security should be converted to []SecurityRequirement")
	assert.Len(t, requirements, 1, "Should have 1 security requirement")

	// Check security schemes were created
	assert.NotNil(t, service.SecuritySchemes)
	assert.Len(t, service.SecuritySchemes, 2, "Should have 2 security schemes")

	// Verify scheme details
	scheme1, exists := service.SecuritySchemes["Basic_ClientID"]
	assert.True(t, exists)
	assert.Equal(t, "Client ID for authentication", scheme1.Description)

	scheme2, exists := service.SecuritySchemes["Basic_ClientSecret"]
	assert.True(t, exists)
	assert.Equal(t, "Client secret for authentication", scheme2.Description)
}
