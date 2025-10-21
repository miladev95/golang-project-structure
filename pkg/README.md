# Package Directory - Reusable Libraries

The `pkg/` directory contains reusable, production-ready utilities that can be imported and used by other projects.

## 📦 Contents

### `utils/` - Common Utility Functions

Reusable helper functions organized into modules:

#### 1. **String Utilities** (`utils/string.go`)

String manipulation helpers:

```go
import "github.com/miladev95/golang-project-structure/pkg/utils"

// Slugify - Convert to URL-friendly format
slug := utils.Slugify("Hello World") // "hello-world"

// TitleCase - Capitalize each word
title := utils.TitleCase("hello world") // "Hello World"

// Capitalize - Capitalize first character
cap := utils.Capitalize("hello") // "Hello"

// IsEmpty - Check if string is empty/whitespace
empty := utils.IsEmpty("   ") // true

// TruncateString - Truncate with ellipsis
truncated := utils.TruncateString("Hello World", 5) // "Hello..."

// ContainsWord - Check if string contains word
has := utils.ContainsWord("hello world", "world") // true

// ReverseString - Reverse a string
rev := utils.ReverseString("hello") // "olleh"
```

#### 2. **Validation Utilities** (`utils/validation.go`)

Input validation helpers:

```go
// Email validation
valid := utils.IsValidEmail("test@example.com") // true

// Phone validation
valid := utils.IsValidPhoneNumber("123-456-7890") // true

// Username validation (3-20 chars, alphanumeric, -, _)
valid := utils.IsValidUsername("john_doe") // true

// Password validation (8+ chars, uppercase, lowercase, digit)
valid := utils.IsValidPassword("SecurePass123") // true

// URL validation
valid := utils.IsValidURL("https://example.com") // true

// IPv4 validation
valid := utils.IsValidIP("192.168.1.1") // true

// UUID validation
valid := utils.IsValidUUID("550e8400-e29b-41d4-a716-446655440000") // true

// String in slice
exists := utils.IsStringInSlice("apple", []string{"apple", "banana"}) // true

// Number between range
inRange := utils.IsNumberBetween(5, 1, 10) // true
```

**Validator Pattern** (for combining multiple validations):

```go
validator := &utils.Validator{}
validator.AddError("email", "Invalid email format")
validator.AddError("password", "Too short")

if validator.HasErrors() {
    errors := validator.GetErrors() // []ValidationError
}
```

**ValidationErrors Pattern** (API-friendly validation):

```go
ve := utils.NewValidationErrors()
ve.Add("email", "Email already exists")
ve.AddWithValue("age", "Must be 18+", 16)

if ve.HasErrors() {
    // Return to client: ve.Errors
}
```

#### 3. **Pagination Utilities** (`utils/pagination.go`)

Pagination helpers for list endpoints:

```go
// Create pagination
page := 2
pageSize := 10
total := int64(100)
pagination := utils.NewPagination(page, pageSize, total)

// Calculate database query parameters
offset := pagination.GetOffset()    // 10
limit := pagination.GetLimit()      // 10

// Navigation
hasPrev := pagination.HasPreviousPage() // true
hasNext := pagination.HasNextPage()     // true
prevPage := pagination.GetPreviousPage() // 1
nextPage := pagination.GetNextPage()     // 3

// Response with pagination
data := []User{...}
response := utils.NewPaginatedResponse(data, pagination)
// Returns: {data: [...], pagination: {page: 2, page_size: 10, total: 100, total_page: 10}}
```

**Example with GORM:**

```go
pagination := utils.NewPagination(page, pageSize, total)

var users []User
db.Offset(pagination.GetOffset()).
   Limit(pagination.GetLimit()).
   Find(&users)

response := utils.NewPaginatedResponse(users, pagination)
c.JSON(200, response)
```

#### 4. **Error Utilities** (`utils/errors.go`)

Custom error types for consistent error handling:

**Basic AppError:**

```go
err := utils.NewAppError("INVALID_INPUT", "The provided input is invalid")
err = err.SetDetails(map[string]string{"field": "email"})

if err != nil {
    // "INVALID_INPUT: The provided input is invalid"
    log.Println(err.Error())
}
```

**AppError with Cause:**

```go
originalErr := dbQuery.Error
err := utils.NewAppErrorWithCause(
    "DB_QUERY_FAILED",
    "Failed to fetch user",
    originalErr,
)
```

**Not Found Error:**

```go
err := utils.NewNotFoundError("User", 123)
// "User with id 123 not found"
```

**Conflict Error:**

```go
err := utils.NewConflictError("Email already exists")
// "conflict: Email already exists"
```

**Unauthorized Error:**

```go
err := utils.NewUnauthorizedError("Invalid token")
// "unauthorized: Invalid token"
```

**Forbidden Error:**

```go
err := utils.NewForbiddenError("You don't have permission")
// "forbidden: You don't have permission"
```

**Internal Server Error:**

```go
err := utils.NewInternalServerError("Payment processing failed", originalErr)
// "internal server error: Payment processing failed (cause: ...)"
```

**ValidationError (single):**

```go
err := utils.ValidationError{
    Field:   "email",
    Message: "Invalid format",
}
// "validation error on field 'email': Invalid format"
```

**ValidationErrors (collection):**

```go
ve := utils.NewValidationErrors()
ve.Add("email", "Invalid format")
ve.Add("password", "Too short")
ve.AddWithValue("age", "Must be 18+", 16)

if ve.HasErrors() {
    // JSON response with all validation errors
    c.JSON(400, ve)
}
```

## 🧪 Testing

All utilities include comprehensive test coverage. Run tests:

```bash
# Run all tests
go test ./pkg/utils -v

# Run specific test
go test ./pkg/utils -v -run TestSlugify

# With coverage
go test ./pkg/utils -v -cover
```

**Test files:**
- `utils/string_test.go` - 7 test functions, 30+ test cases
- `utils/validation_test.go` - 10 test functions, 50+ test cases
- `utils/pagination_test.go` - 10 test functions, 40+ test cases
- `utils/errors_test.go` - 10 test functions, 30+ test cases

## 📊 Available Test Coverage

### String Utilities (30+ test cases)
- ✅ Slugify (special chars, spaces, case handling)
- ✅ TitleCase (various case combinations)
- ✅ Capitalize (edge cases)
- ✅ IsEmpty (whitespace handling)
- ✅ TruncateString (various lengths)
- ✅ ContainsWord (case sensitivity, partial matches)
- ✅ ReverseString (unicode support)

### Validation Utilities (50+ test cases)
- ✅ Email validation (valid, invalid formats)
- ✅ Phone validation (various formats)
- ✅ Username validation (length, characters)
- ✅ Password validation (strength requirements)
- ✅ URL validation (protocols, ports, paths)
- ✅ IP validation (ranges, invalid values)
- ✅ String in slice (exists, not exists)
- ✅ Number between (ranges)
- ✅ Validator pattern
- ✅ ValidationErrors pattern

### Pagination Utilities (40+ test cases)
- ✅ Pagination creation (defaults, limits)
- ✅ Offset calculation (various pages)
- ✅ Page navigation (previous, next)
- ✅ Validation (valid ranges)
- ✅ PaginatedResponse creation

### Error Utilities (30+ test cases)
- ✅ AppError (creation, with/without cause)
- ✅ NotFoundError
- ✅ ConflictError
- ✅ UnauthorizedError
- ✅ ForbiddenError
- ✅ InternalServerError
- ✅ ValidationError & ValidationErrors

## 🔗 Usage Examples

### In Handlers

```go
package handlers

import (
    "github.com/miladev95/golang-project-structure/pkg/utils"
)

func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    c.BindJSON(&req)

    // Validate input
    ve := utils.NewValidationErrors()
    
    if !utils.IsValidEmail(req.Email) {
        ve.Add("email", "Invalid email format")
    }
    
    if !utils.IsValidPassword(req.Password) {
        ve.Add("password", "Password must be 8+ chars with uppercase, lowercase, digit")
    }

    if ve.HasErrors() {
        c.JSON(400, ve)
        return
    }

    // Business logic...
}
```

### In Services

```go
package services

import "github.com/miladev95/golang-project-structure/pkg/utils"

func (s *UserService) GetUsers(page, pageSize int, total int64) (interface{}, error) {
    pagination := utils.NewPagination(page, pageSize, total)
    
    users, err := s.repo.FindUsers(pagination.GetOffset(), pagination.GetLimit())
    if err != nil {
        return nil, utils.NewInternalServerError("Failed to fetch users", err)
    }

    return utils.NewPaginatedResponse(users, pagination), nil
}
```

### In Repositories

```go
package repositories

import "github.com/miladev95/golang-project-structure/pkg/utils"

func (r *UserRepository) FindByID(id int) (*User, error) {
    var user User
    result := r.db.First(&user, id)
    
    if result.Error != nil {
        return nil, utils.NewNotFoundError("User", id)
    }
    
    return &user, nil
}
```

## 📝 Test Output Example

```bash
$ go test ./pkg/utils -v

=== RUN   TestSlugify
--- PASS: TestSlugify (0.00s)
    === RUN   TestSlugify/simple_string
    --- PASS: TestSlugify/simple_string (0.00s)
    === RUN   TestSlugify/with_special_characters
    --- PASS: TestSlugify/with_special_characters (0.00s)

=== RUN   TestIsValidEmail
--- PASS: TestIsValidEmail (0.00s)
    === RUN   TestIsValidEmail/valid_email
    --- PASS: TestIsValidEmail/valid_email (0.00s)
    === RUN   TestIsValidEmail/invalid_no_@
    --- PASS: TestIsValidEmail/invalid_no_@ (0.00s)

ok  	github.com/miladev95/golang-project-structure/pkg/utils	0.234s
coverage: 95.4% of statements
```

## 🚀 Getting Started

1. **Import in your code:**
   ```go
   import "github.com/miladev95/golang-project-structure/pkg/utils"
   ```

2. **Use utilities:**
   ```go
   if !utils.IsValidEmail(email) {
       return errors.New("invalid email")
   }
   ```

3. **Run tests:**
   ```bash
   go test ./pkg/utils -v -cover
   ```

## 📦 Public API

The following are exported and available for import:

**String Functions:**
- `Slugify(string) string`
- `TitleCase(string) string`
- `Capitalize(string) string`
- `IsEmpty(string) bool`
- `TruncateString(string, int) string`
- `ContainsWord(string, string) bool`
- `ReverseString(string) string`

**Validation Functions:**
- `IsValidEmail(string) bool`
- `IsValidPhoneNumber(string) bool`
- `IsValidUsername(string) bool`
- `IsValidPassword(string) bool`
- `IsValidURL(string) bool`
- `IsValidIP(string) bool`
- `IsValidUUID(string) bool`
- `IsStringInSlice(string, []string) bool`
- `IsNumberBetween(int64, int64, int64) bool`

**Types & Methods:**
- `type Validator struct`
- `type ValidationError struct`
- `type ValidationErrors struct`
- `type Pagination struct`
- `type PaginatedResponse struct`
- `type AppError struct`
- `type NotFoundError struct`
- `type ConflictError struct`
- `type UnauthorizedError struct`
- `type ForbiddenError struct`
- `type InternalServerError struct`

## 🔄 Extending

To add new utilities:

1. Create new file: `pkg/utils/new_feature.go`
2. Create test file: `pkg/utils/new_feature_test.go`
3. Add to this README
4. Run tests: `go test ./pkg/utils`

## 📄 License

Same as parent project