// Package servergen provides server code generation capabilities from specification types.
//
// This package generates Gin-based server implementations directly from specification.Service
// definitions. It creates a complete server API with type-safe request/response handling,
// automatic parameter parsing, and error management.
//
// # Key Features
//
// - Generates complete Gin server implementation from specification.Service
// - Type-safe request/response structures with generics
// - Automatic path, query, and body parameter parsing
// - Built-in error handling with HTTP status code mapping
// - Session management support with generic session types
// - Enum types based on github.com/meitner-se/go-types
// - Embedded OpenAPI specification support
//
// # Generation Process
//
// The package exports a single function GenerateServer that writes the generated
// code to a bytes.Buffer:
//
//	import (
//	    "bytes"
//	    "github.com/meitner-se/publicapis-gen/specification"
//	    "github.com/meitner-se/publicapis-gen/specification/servergen"
//	)
//
//	// Load specification
//	service, err := specification.ParseServiceFromFile("api-spec.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Generate server code
//	var buf bytes.Buffer
//	err = servergen.GenerateServer(&buf, service)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Write to file
//	err = os.WriteFile("generated/server.go", buf.Bytes(), 0644)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Generated Code Structure
//
// The generated server includes:
//
// 1. **Enum Types**: All enums are generated as typed strings using go-types:
//
//	type ErrorCode types.String
//	var (
//	    ErrorCodeBadRequest = ErrorCode(types.NewString("BadRequest"))
//	    ErrorCodeNotFound   = ErrorCode(types.NewString("NotFound"))
//	)
//
// 2. **Object Types**: Objects are generated as Go structs with JSON tags:
//
//	type User struct {
//	    ID       types.UUID   `json:"id"`
//	    Name     types.String `json:"name"`
//	    Email    types.String `json:"email"`
//	}
//
// 3. **Request Types**: Generic request structure with path, query, and body parameters:
//
//	type Request[sessionType, pathParamsType, queryParamsType, bodyParamsType any] struct {
//	    Session     sessionType
//	    PathParams  pathParamsType
//	    QueryParams queryParamsType
//	    BodyParams  bodyParamsType
//	}
//
// 4. **API Interfaces**: Each resource gets an interface defining its endpoints:
//
//	type UsersAPI[Session any] interface {
//	    Create(ctx context.Context, request Request[Session, struct{}, struct{}, CreateUserRequest]) (*User, error)
//	    GetByID(ctx context.Context, request Request[Session, GetUserPathParams, struct{}, struct{}]) (*User, error)
//	}
//
// 5. **Server Registration**: Main registration function that sets up all routes:
//
//	func RegisterServiceAPI[Session any](router *gin.Engine, api *ServiceAPI[Session])
//
// # Session Management
//
// The generated server supports generic session management. Each endpoint receives
// a session object extracted by a user-provided GetSessionFunc:
//
//	type Server[Session any] struct {
//	    GetSessionFunc    func(ctx context.Context, headers http.Header, requestID string) (Session, error)
//	    ConvertErrorFunc  func(err error, requestID string) *Error
//	}
//
// # Error Handling
//
// Errors are automatically converted to API error responses with appropriate
// HTTP status codes. The Error type implements both error and HTTPStatusCode interfaces.
//
// # Type Safety
//
// The package leverages github.com/meitner-se/go-types for type-safe handling of:
// - UUIDs
// - Nullable values
// - Arrays
// - Standard types (String, Int, Bool, etc.)
//
// Field modifiers from the specification (nullable, array) are automatically
// applied to generate the correct Go types.
package servergen
