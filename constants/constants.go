// Package constants defines all string constants used throughout the publicapis-gen codebase.
// This centralizes hardcoded strings to improve maintainability and consistency.
package constants

// Error messages
const (
	ErrorNotImplemented    = "not implemented"
	ErrorFailedToRun       = "failed to run"
	ErrorFailedToGenerate  = "failed to generate schema for"
	ErrorValidationFailed  = "validation failed"
	ErrorValidationErrors  = "validation errors"
	ErrorFailedToMarshal   = "failed to marshal schema to JSON"
	ErrorFailedToConvert   = "failed to convert schema to JSON"
	ErrorFailedToUnmarshal = "failed to unmarshal"
	ErrorDataNotValid      = "data is neither valid JSON nor YAML"
	ErrorConversionFailed  = "failed to convert YAML to JSON"
)

// CRUD Operations
const (
	OperationCreate = "Create"
	OperationRead   = "Read"
	OperationUpdate = "Update"
	OperationDelete = "Delete"
)

// Field types
const (
	TypeUUID      = "UUID"
	TypeString    = "String"
	TypeInt       = "Int"
	TypeBool      = "Bool"
	TypeTimestamp = "Timestamp"
	TypeDate      = "Date"
	TypeArray     = "array"
	TypeObject    = "Object"
)

// HTTP Methods
const (
	HTTPMethodGET    = "GET"
	HTTPMethodPOST   = "POST"
	HTTPMethodPUT    = "PUT"
	HTTPMethodDELETE = "DELETE"
	HTTPMethodPATCH  = "PATCH"
)

// Content Types
const (
	ContentTypeJSON       = "application/json"
	ContentTypeFormData   = "multipart/form-data"
	ContentTypeXML        = "application/xml"
	ContentTypeFormURLEnc = "application/x-www-form-urlencoded"
)

// Field modifiers
const (
	ModifierArray    = "array"
	ModifierNullable = "nullable"
	ModifierOptional = "optional"
)

// Common field names
const (
	FieldNameID          = "id"
	FieldNameName        = "name"
	FieldNameEmail       = "email"
	FieldNameDescription = "description"
	FieldNameType        = "type"
	FieldNameDefault     = "default"
	FieldNameExample     = "example"
	FieldNameModifiers   = "modifiers"
	FieldNameOperations  = "operations"
)

// Schema field names
const (
	SchemaFieldName        = "name"
	SchemaFieldEnums       = "enums"
	SchemaFieldObjects     = "objects"
	SchemaFieldResources   = "resources"
	SchemaFieldValues      = "values"
	SchemaFieldFields      = "fields"
	SchemaFieldEndpoints   = "endpoints"
	SchemaFieldHeaders     = "headers"
	SchemaFieldPathParams  = "path_params"
	SchemaFieldQueryParams = "query_params"
	SchemaFieldBodyParams  = "body_params"
	SchemaFieldBodyFields  = "body_fields"
	SchemaFieldBodyObject  = "body_object"
	SchemaFieldContentType = "content_type"
	SchemaFieldStatusCode  = "status_code"
	SchemaFieldTitle       = "title"
	SchemaFieldMethod      = "method"
	SchemaFieldPath        = "path"
	SchemaFieldRequest     = "request"
	SchemaFieldResponse    = "response"
)

// Log message keys
const (
	LogKeyError = "error"
)

// Validation error messages
const (
	ValidationErrorRequired   = "required"
	ValidationErrorAdditional = "Additional property"
	ValidationErrorInteger    = "integer"
	ValidationErrorArray      = "array"
	ValidationErrorString     = "string"
)
