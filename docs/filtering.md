# Build Advanced Filtering

**Task**: Create powerful search capabilities with generated filter objects

## Generate filter objects

### Input: Basic resource definition

```yaml
name: "E-commerce API"
version: "1.0.0"

enums:
  - name: "ProductCategory"
    description: "Product categories"
    values:
      - name: "Electronics"
        description: "Electronic devices"
      - name: "Clothing"
        description: "Apparel and accessories"
      - name: "Books"
        description: "Books and publications"

objects:
  - name: "Price"
    description: "Price information"
    fields:
      - name: "amount"
        type: "Int"
        description: "Price amount in cents"
      - name: "currency"
        type: "String"
        description: "Currency code"

resources:
  - name: "Products"
    description: "Product catalog"
    operations: ["Create", "Read", "Update", "Delete"]
    fields:
      - field:
          name: "title"
          type: "String"
          description: "Product title"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "description"
          type: "String"
          description: "Product description"
          modifiers: ["Nullable"]
        operations: ["Create", "Read", "Update"]
      - field:
          name: "category"
          type: "ProductCategory"
          description: "Product category"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "price"
          type: "Price"
          description: "Product pricing"
        operations: ["Create", "Read", "Update"]
      - field:
          name: "tags"
          type: "String"
          description: "Product tags"
          modifiers: ["Array", "Nullable"]
        operations: ["Create", "Read", "Update"]
      - field:
          name: "inStock"
          type: "Bool"
          description: "Availability status"
        operations: ["Read", "Update"]
```

### Output: Complete filter system

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/meitner-se/publicapis-gen/specification"
)

func generateFilters() {
    // Parse specification and generate filters
    service, err := specification.ParseServiceFromFile("products-api.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("🔍 Generated filter objects:\n")
    for _, obj := range service.Objects {
        if strings.Contains(obj.Name, "Filter") {
            fmt.Printf("  • %s (%d fields)\n", obj.Name, len(obj.Fields))
        }
    }
    
    // Show the main filter structure  
    mainFilter := service.GetObject("ProductsFilter")
    if mainFilter != nil {
        fmt.Printf("\n📋 ProductsFilter structure:\n")
        for _, field := range mainFilter.Fields {
            fmt.Printf("  • %s: %s\n", field.Name, field.Type)
        }
    }
}

func main() {
    generateFilters()
}
```

**Generated filter objects:**
```
🔍 Generated filter objects:
  • ProductsFilter (14 fields)
  • ProductsFilterEquals (6 fields)  
  • ProductsFilterRange (1 fields)
  • ProductsFilterContains (5 fields)
  • ProductsFilterLike (2 fields)
  • ProductsFilterNull (3 fields)
  • PriceFilter (14 fields)
  • PriceFilterEquals (2 fields)
  • PriceFilterRange (1 fields)
  • PriceFilterContains (2 fields)
  • PriceFilterLike (1 fields)
  • PriceFilterNull (0 fields)

📋 ProductsFilter structure:
  • Equals: ProductsFilterEquals
  • NotEquals: ProductsFilterEquals
  • GreaterThan: ProductsFilterRange
  • SmallerThan: ProductsFilterRange
  • GreaterOrEqual: ProductsFilterRange
  • SmallerOrEqual: ProductsFilterRange
  • Contains: ProductsFilterContains
  • NotContains: ProductsFilterContains
  • Like: ProductsFilterLike
  • NotLike: ProductsFilterLike
  • Null: ProductsFilterNull
  • NotNull: ProductsFilterNull
  • OrCondition: Bool
  • NestedFilters: ProductsFilter (Array)
```

## Create simple filters

### Task: Filter by exact values

**Input**: Search for specific products

```json
POST /products/_search
Content-Type: application/json

{
  "Filter": {
    "Equals": {
      "category": "Electronics",
      "inStock": true
    }
  }
}
```

**Output**: Products matching exact criteria

```json
{
  "data": [
    {
      "ID": "123e4567-e89b-12d3-a456-426614174000",
      "Meta": {
        "CreatedAt": "2024-01-15T10:00:00Z",
        "UpdatedAt": "2024-01-15T10:00:00Z"
      },
      "title": "Wireless Headphones",
      "description": "High-quality wireless headphones",
      "category": "Electronics", 
      "price": {
        "amount": 15999,
        "currency": "USD"
      },
      "tags": ["wireless", "audio", "bluetooth"],
      "inStock": true
    }
  ],
  "Pagination": {
    "Offset": 0,
    "Limit": 50,
    "Total": 1
  }
}
```

## Create range filters

### Task: Filter by price range

**Input**: Products within price range

```json
POST /products/_search  
Content-Type: application/json

{
  "Filter": {
    "GreaterOrEqual": {
      "price": {
        "amount": 5000  // $50.00 in cents
      }
    },
    "SmallerOrEqual": {
      "price": {
        "amount": 20000  // $200.00 in cents  
      }
    }
  }
}
```

**Generated Go code for handling this filter:**

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Generated filter types
type ProductsFilter struct {
    Equals         *ProductsFilterEquals   `json:"Equals,omitempty"`
    NotEquals      *ProductsFilterEquals   `json:"NotEquals,omitempty"`
    GreaterThan    *ProductsFilterRange    `json:"GreaterThan,omitempty"`
    SmallerThan    *ProductsFilterRange    `json:"SmallerThan,omitempty"`
    GreaterOrEqual *ProductsFilterRange    `json:"GreaterOrEqual,omitempty"`
    SmallerOrEqual *ProductsFilterRange    `json:"SmallerOrEqual,omitempty"`
    Contains       *ProductsFilterContains `json:"Contains,omitempty"`
    NotContains    *ProductsFilterContains `json:"NotContains,omitempty"`
    Like           *ProductsFilterLike     `json:"Like,omitempty"`
    NotLike        *ProductsFilterLike     `json:"NotLike,omitempty"`
    Null           *ProductsFilterNull     `json:"Null,omitempty"`
    NotNull        *ProductsFilterNull     `json:"NotNull,omitempty"`
    OrCondition    bool                    `json:"OrCondition,omitempty"`
    NestedFilters  []ProductsFilter        `json:"NestedFilters,omitempty"`
}

type ProductsFilterRange struct {
    Price *PriceFilterRange `json:"price,omitempty"`
}

type PriceFilterRange struct {
    Amount *int `json:"amount,omitempty"`  // Only comparable fields
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Filter ProductsFilter `json:"Filter"`
    }
    
    json.NewDecoder(r.Body).Decode(&request)
    
    // Use the filter to build database query
    query := buildDatabaseQuery(request.Filter)
    products := executeQuery(query)
    
    response := map[string]interface{}{
        "data": products,
        "Pagination": map[string]int{
            "Offset": 0,
            "Limit":  50,
            "Total":  len(products),
        },
    }
    
    json.NewEncoder(w).Encode(response)
}
```

## Create text search filters

### Task: Search by title and description content

**Input**: Text pattern matching

```json
POST /products/_search
Content-Type: application/json

{
  "Filter": {
    "Like": {
      "title": "%wireless%",
      "description": "%headphones%"
    }
  }
}
```

**Generated filter types for text search:**

```go
type ProductsFilterLike struct {
    Title       *string `json:"title,omitempty"`       // String fields only
    Description *string `json:"description,omitempty"` // String fields only
}
```

## Create array filters  

### Task: Filter by multiple values

**Input**: Products with specific tags

```json
POST /products/_search
Content-Type: application/json

{
  "Filter": {
    "Contains": {
      "tags": ["wireless", "bluetooth"],      // Match any of these tags
      "category": ["Electronics", "Books"]     // Match any of these categories  
    }
  }
}
```

**Generated filter types for arrays:**

```go
type ProductsFilterContains struct {
    ID          []string           `json:"ID,omitempty"`
    Title       []string           `json:"title,omitempty"`
    Description []string           `json:"description,omitempty"`
    Category    []string           `json:"category,omitempty"`  // Enums become string arrays
    Tags        []string           `json:"tags,omitempty"`      // Already arrays
    InStock     []bool             `json:"inStock,omitempty"`   // Bool becomes bool array
    Price       *PriceFilterContains `json:"price,omitempty"`   // Nested objects
}
```

## Create complex filters

### Task: Combine multiple filter types with OR logic

**Input**: Complex search criteria

```json
POST /products/_search
Content-Type: application/json

{
  "Filter": {
    "OrCondition": true,
    "NestedFilters": [
      {
        "Equals": {
          "category": "Electronics"  
        },
        "GreaterThan": {
          "price": {
            "amount": 10000  // $100.00
          }
        }
      },
      {
        "Equals": {
          "category": "Books"
        },
        "Like": {
          "title": "%programming%"
        }
      }
    ]
  }
}
```

**Meaning**: Find products that are either:
1. Electronics AND cost more than $100, OR  
2. Books AND have "programming" in the title

## Handle null values

### Task: Filter by presence/absence of data

**Input**: Check for null/non-null fields

```json
POST /products/_search
Content-Type: application/json

{
  "Filter": {
    "NotNull": {
      "description": true,  // Must have description
      "tags": true         // Must have tags array
    },
    "Null": {
      "price": false       // Price must NOT be null
    }
  }
}
```

**Generated null filter types:**

```go
type ProductsFilterNull struct {
    Description *bool `json:"description,omitempty"` // Nullable fields only  
    Tags        *bool `json:"tags,omitempty"`        // Array fields only
}
```

## Implementation example

### Task: Build a complete search API

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
)

// Mock database
var products = []Product{
    {
        ID: "1",
        Title: "Wireless Headphones",
        Description: stringPtr("Premium wireless headphones with noise cancellation"),
        Category: "Electronics",
        Price: Price{Amount: 15999, Currency: "USD"},
        Tags: []string{"wireless", "audio", "bluetooth"},
        InStock: true,
    },
    {
        ID: "2", 
        Title: "Programming Book",
        Description: stringPtr("Learn Go programming"),
        Category: "Books",
        Price: Price{Amount: 3999, Currency: "USD"},
        Tags: []string{"programming", "go", "tutorial"},
        InStock: true,
    },
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    
    limit := 50
    offset := 0
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil {
            offset = o
        }
    }
    
    // Parse filter from request body
    var request struct {
        Filter ProductsFilter `json:"Filter"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Apply filters
    filtered := applyFilters(products, request.Filter)
    
    // Apply pagination  
    total := len(filtered)
    start := offset
    end := offset + limit
    if start >= total {
        filtered = []Product{}
    } else {
        if end > total {
            end = total
        }
        filtered = filtered[start:end]
    }
    
    // Return results
    response := SearchResponse{
        Data: filtered,
        Pagination: Pagination{
            Offset: offset,
            Limit:  limit,  
            Total:  total,
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func applyFilters(products []Product, filter ProductsFilter) []Product {
    var result []Product
    
    for _, product := range products {
        if matchesFilter(product, filter) {
            result = append(result, product)
        }
    }
    
    return result
}

func matchesFilter(product Product, filter ProductsFilter) bool {
    // Handle OR condition with nested filters
    if filter.OrCondition && len(filter.NestedFilters) > 0 {
        for _, nestedFilter := range filter.NestedFilters {
            if matchesFilter(product, nestedFilter) {
                return true // OR - any nested filter matching is enough
            }
        }
        return false
    }
    
    // Handle AND conditions (default)
    if len(filter.NestedFilters) > 0 {
        for _, nestedFilter := range filter.NestedFilters {
            if !matchesFilter(product, nestedFilter) {
                return false // AND - all nested filters must match
            }
        }
    }
    
    // Check equals filters
    if filter.Equals != nil {
        if !matchesEquals(product, *filter.Equals) {
            return false
        }
    }
    
    // Check range filters
    if filter.GreaterThan != nil {
        if !matchesRange(product, *filter.GreaterThan, "gt") {
            return false
        }
    }
    
    // Check LIKE filters
    if filter.Like != nil {
        if !matchesLike(product, *filter.Like) {
            return false
        }
    }
    
    // Check contains filters
    if filter.Contains != nil {
        if !matchesContains(product, *filter.Contains) {
            return false
        }
    }
    
    return true
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/products/_search", searchProducts).Methods("POST")
    
    fmt.Println("🔍 Search API running on :8080")
    fmt.Println("Try: POST /products/_search with filter JSON")
    http.ListenAndServe(":8080", r)
}
```

## Filter usage patterns

### Pattern: Category browsing
```json
{
  "Filter": {
    "Equals": {
      "category": "Electronics",
      "inStock": true
    }
  }
}
```

### Pattern: Price range search  
```json
{
  "Filter": {
    "GreaterOrEqual": {
      "price": {"amount": 1000}
    },
    "SmallerOrEqual": {
      "price": {"amount": 5000}
    }
  }
}
```

### Pattern: Text search
```json
{
  "Filter": {
    "Like": {
      "title": "%phone%",
      "description": "%wireless%"
    }
  }
}
```

### Pattern: Multiple choice
```json
{
  "Filter": {
    "Contains": {
      "category": ["Electronics", "Computers"],
      "tags": ["sale", "featured"]
    }
  }
}
```

### Pattern: Advanced combination
```json
{
  "Filter": {
    "OrCondition": true,
    "NestedFilters": [
      {
        "Equals": {"category": "Electronics"},
        "Contains": {"tags": ["featured"]}
      },
      {
        "GreaterThan": {"price": {"amount": 20000}},
        "Like": {"title": "%premium%"}
      }
    ]
  }
}
```

## Best Practices

### ✅ Do's
- **Use appropriate filter types**: Equals for exact matches, Like for text search, Range for numbers
- **Combine filters logically**: Use nested filters for complex AND/OR logic  
- **Handle pagination**: Always include limit/offset for large result sets
- **Validate input**: Check filter values before applying to database
- **Index filtered fields**: Ensure database indexes support your common filters

### ❌ Don'ts
- **Don't ignore null handling**: Check for null values when filtering
- **Don't skip validation**: Malformed filters can crash your API
- **Don't forget performance**: Complex nested filters can be slow
- **Don't hardcode limits**: Make pagination configurable

## Related Tasks

- [🚀 Getting Started](getting-started.md) - Create your first specification with filters
- [⚙️ Working with Specifications](specifications.md) - Design filterable resources  
- [📋 Generate OpenAPI](openapi.md) - Document your search endpoints
- [✅ JSON Schema Validation](schema-validation.md) - Validate filter requests