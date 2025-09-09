# Blog API Example

Complete example of a blog management API using publicapis-gen.

## What this example shows

- 📝 **Complete blog management** with posts, authors, and categories
- 🔍 **Advanced filtering** with full-text search and category filtering
- 📋 **OpenAPI 3.1 generation** with complete documentation
- ✅ **Request validation** with JSON schemas
- 🚀 **HTTP server implementation** with generated endpoints

## Files

- `blog-api.yaml` - Main API specification
- `main.go` - Complete HTTP server implementation  
- `generate.go` - Script to generate OpenAPI and schemas
- `generated/` - Generated OpenAPI and JSON schemas
- `README.md` - This file

## Quick Start

1. **Install dependencies:**
   ```bash
   go mod init blog-api-example
   go get github.com/meitner-se/publicapis-gen
   go get github.com/gorilla/mux
   ```

2. **Generate documentation:**
   ```bash
   go run generate.go
   ```

3. **Start the API server:**
   ```bash
   go run main.go
   ```

4. **View API documentation:**
   - OpenAPI: `http://localhost:8080/docs`
   - JSON Schema: See `generated/schemas/` directory

## API Endpoints

### Generated CRUD Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/posts` | List all posts with pagination |
| `POST` | `/posts` | Create a new post |  
| `GET` | `/posts/{id}` | Get a specific post |
| `PATCH` | `/posts/{id}` | Update a post |
| `DELETE` | `/posts/{id}` | Delete a post |
| `POST` | `/posts/_search` | Search posts with filters |

### Search Examples

**Find published posts:**
```bash
curl -X POST http://localhost:8080/posts/_search \
  -H "Content-Type: application/json" \
  -d '{
    "Filter": {
      "Equals": {
        "status": "Published"
      }
    }
  }'
```

**Full-text search in titles:**
```bash
curl -X POST http://localhost:8080/posts/_search \
  -H "Content-Type: application/json" \
  -d '{
    "Filter": {
      "Like": {
        "title": "%golang%"
      }
    }
  }'
```

**Complex filter - published posts about programming:**
```bash
curl -X POST http://localhost:8080/posts/_search \
  -H "Content-Type: application/json" \
  -d '{
    "Filter": {
      "Equals": {
        "status": "Published"
      },
      "Contains": {
        "tags": ["programming", "tutorial"]
      }
    }
  }'
```

## Learning Objectives

After studying this example, you'll understand:

- ✅ How to structure a complete API specification
- ✅ How to use enums for controlled values
- ✅ How to create reusable objects  
- ✅ How to leverage auto-generated CRUD endpoints
- ✅ How to implement advanced search with filters
- ✅ How to validate requests with generated schemas
- ✅ How to serve OpenAPI documentation
- ✅ How to handle errors consistently

## Extension Ideas

Try extending this example:

- 🔐 **Add authentication** - JWT tokens, API keys
- 📊 **Add analytics** - View counts, popular posts
- 🏷️ **Tag management** - CRUD operations for tags
- 💬 **Add comments** - Nested resources with moderation
- 🖼️ **File uploads** - Image attachments for posts
- 📧 **Notifications** - Email alerts for new posts