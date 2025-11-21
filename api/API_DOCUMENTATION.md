# Backend API Documentation

## Overview

This document provides comprehensive documentation for the backend API application built with Go and the Gin framework. The API serves as a RESTful service with JWT authentication, PostgreSQL database, and Redis caching.

## Architecture Overview

### Technology Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with Gorp ORM
- **Cache**: Redis for JWT token storage
- **Authentication**: JWT with refresh tokens
- **Validation**: Go Playground Validator
- **Documentation**: Swagger/OpenAPI
- **Compression**: Gzip middleware

### Project Structure

```
api/
├── controllers/     # HTTP request handlers
├── models/         # Data models and business logic
├── forms/          # Request validation structures
├── db/             # Database connection and schemas
├── docs/           # API documentation (Swagger)
├── public/         # Static files
├── tests/          # Unit tests
├── main.go         # Application entry point
└── go.mod          # Dependencies
```

### Design Patterns

- **MVC-like Architecture**: Controllers handle HTTP requests, models manage data and business logic
- **Dependency Injection**: Controllers receive model instances
- **Middleware Chain**: CORS, authentication, logging, compression
- **Repository Pattern**: Models encapsulate database operations
- **Form Validation**: Centralized validation with custom error messages

## API Endpoints

### Authentication & User Management

#### POST `/v1/user/login`

**Description**: Authenticate user and return JWT tokens
**Authentication**: None required
**Request Body**:

```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

**Response**:

```json
{
  "message": "Successfully logged in",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "name": "User Name"
  },
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### POST `/v1/user/register`

**Description**: Register a new user account
**Authentication**: None required
**Request Body**:

```json
{
  "name": "User Name",
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

**Response**:

```json
{
  "message": "Successfully registered",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "name": "User Name"
  }
}
```

#### GET `/v1/user/logout`

**Description**: Logout user by invalidating tokens
**Authentication**: Bearer token required
**Response**:

```json
{
  "message": "Successfully logged out"
}
```

#### GET `/v1/user/profile`

**Description**: Get current user profile information
**Authentication**: Bearer token required
**Response**:

```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "name": "User Name"
}
```

#### POST `/v1/user/forgot-password`

**Description**: Initiate password reset (stub implementation)
**Authentication**: None required
**Request Body**:

```json
{
  "email": "user@example.com"
}
```

**Response**:

```json
{
  "message": "Password reset link sent to your email"
}
```

#### POST `/v1/user/assign-role`

**Description**: Assign a role to a user
**Authentication**: Bearer token + manage_users permission required
**Request Body**:

```json
{
  "user_id": 1,
  "role_name": "admin"
}
```

**Response**:

```json
{
  "message": "Role assigned successfully"
}
```

#### POST `/v1/permission/create`

**Description**: Create a new permission
**Authentication**: Bearer token + manage_users permission required
**Request Body**:

```json
"read_article"
```

**Response**:

```json
{
  "message": "Permission created successfully"
}
```

#### POST `/v1/token/refresh`

**Description**: Refresh access token using refresh token
**Authentication**: None required
**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response**:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Article Management

#### POST `/v1/article`

**Description**: Create a new article
**Authentication**: Bearer token + write_article permission required
**Request Body**:

```json
{
  "title": "Article Title",
  "content": "Article content here..."
}
```

**Response**:

```json
{
  "message": "Article created",
  "id": 1
}
```

#### GET `/v1/articles`

**Description**: Get all articles for the authenticated user
**Authentication**: Bearer token + read_article permission required
**Response**:

```json
{
  "results": [
    {
      "data": [
        {
          "id": 1,
          "title": "Article Title",
          "content": "Article content...",
          "updated_at": 1640995200,
          "created_at": 1640995200,
          "user": {
            "id": 1,
            "name": "User Name",
            "email": "user@example.com"
          }
        }
      ],
      "meta": {
        "total": 1
      }
    }
  ]
}
```

#### GET `/v1/article/{id}`

**Description**: Get a specific article by ID
**Authentication**: Bearer token + read_article permission required
**Parameters**: `id` (integer) - Article ID
**Response**:

```json
{
  "data": {
    "id": 1,
    "title": "Article Title",
    "content": "Article content...",
    "updated_at": 1640995200,
    "created_at": 1640995200,
    "user": {
      "id": 1,
      "name": "User Name",
      "email": "user@example.com"
    }
  }
}
```

#### PUT `/v1/article/{id}`

**Description**: Update an existing article
**Authentication**: Bearer token + write_article permission required
**Parameters**: `id` (integer) - Article ID
**Request Body**:

```json
{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

**Response**:

```json
{
  "message": "Article updated"
}
```

#### DELETE `/v1/article/{id}`

**Description**: Delete an article
**Authentication**: Bearer token + write_article permission required
**Parameters**: `id` (integer) - Article ID
**Response**:

```json
{
  "message": "Article deleted"
}
```

### Sijagur Data Management

#### GET `/v1/realisasi-bulan`

**Description**: Get monthly realization data
**Authentication**: Bearer token required
**Query Parameters**:

- `tahun` (int): Year (default: current year)
- `bulan` (int): Month (default: current month)
- `idsatker` (int): Satker ID (default: 0 for all)
  **Response**:

```json
{
  "results": [
    {
      "data": [
        {
          "category": "barjas",
          "progress": 85.5,
          "progress_formatted": "85.5",
          "items": [
            {
              "type": "perencanaan",
              "value": 10,
              "detail": {
                "selesai": 8,
                "target": 10,
                "terlambat": 2
              }
            }
          ]
        }
      ],
      "meta": {
        "year": 2024,
        "month": 11,
        "month_name": "November",
        "idsatker": 0,
        "type": "bulan"
      }
    }
  ]
}
```

#### GET `/v1/realisasi-tahun`

**Description**: Get yearly realization data
**Authentication**: Bearer token required
**Query Parameters**: Same as `/realisasi-bulan`
**Response**: Similar structure with yearly aggregated data

#### GET `/v1/realisasi-perbulan`

**Description**: Get monthly breakdown data for the year
**Authentication**: Bearer token required
**Query Parameters**:

- `tahun` (int): Year (default: current year)
- `idsatker` (int): Satker ID (default: 0 for all)
  **Response**:

```json
{
  "results": [
    {
      "data": [
        {
          "category": "barjas",
          "monthly": [
            {
              "month": "Januari",
              "value": 75.0,
              "value_formatted": "75",
              "realisasi": 15,
              "target": 20,
              "realisasi_formatted": "15",
              "target_formatted": "20"
            }
          ],
          "current_month": {
            "month": "November",
            "value": 85.0,
            "value_formatted": "85",
            "realisasi": 17,
            "target": 20
          }
        }
      ],
      "meta": {
        "year": 2024,
        "month": 0,
        "idsatker": 0,
        "type": "perbulan"
      }
    }
  ]
}
```

#### GET `/v1/sijagur/peringkat-kinerja`

**Description**: Get performance rankings
**Authentication**: Bearer token required
**Query Parameters**:

- `year` (int): Required - Year
- `month` (int): Optional - Month
- `idsatker` (int): Optional - Satker ID
- `category` (string): Filter category (all/barjas/fisik/anggaran/kinerja)
- `dimension` (string): kumulatif/capaian/periodik
- `scope` (string): skpd/kecamatan
- `sortBy` (string): Sort field
- `sortDir` (string): asc/desc
  **Response**:

```json
{
  "status": "success",
  "scope": "skpd",
  "category": "all",
  "dimension": "kumulatif",
  "year": 2024,
  "month": 11,
  "page": 1,
  "page_size": 50,
  "total": 50,
  "sort_by": "kumulatif_opd",
  "sort_dir": "desc",
  "data": [
    {
      "id": 1,
      "idsatker": 123,
      "nama_opd": "Dinas Kesehatan",
      "rank_number": 1,
      "score_total": 95.5,
      "score_barjas": 92.0,
      "score_fisik": 98.0,
      "score_anggaran": 94.0,
      "score_kinerja": 96.0,
      "score_status": "Melesat",
      "score_total_formatted": "95.5",
      "year": 2024,
      "month": 11
    }
  ]
}
```

## Authentication & Authorization

### JWT Token Flow

1. **Login**: User provides credentials → Server validates → Returns access + refresh tokens
2. **API Requests**: Client sends Bearer token in Authorization header
3. **Token Refresh**: When access token expires, use refresh token to get new pair
4. **Logout**: Server invalidates tokens in Redis

### Permission System

- **Role-Based Access Control (RBAC)**: Users have roles, roles have permissions
- **Permission Checking**: Middleware validates user has required permission
- **Admin Bypass**: Admin role bypasses permission checks
- **Database Relations**: user_roles, role_permissions tables

### Security Features

- **Password Hashing**: bcrypt with default cost
- **Account Locking**: After 5 failed attempts, lock for 1 minute
- **CORS Configuration**: Environment-specific origin validation
- **Request ID Middleware**: Unique ID for each request
- **Sliding Expiration**: Access tokens refresh on activity

## Data Models

### Core Entities

#### User

```go
type User struct {
    ID             int64  `db:"id"`
    Email          string `db:"email"`
    Username       string `db:"username"`
    Password       string `db:"password"`
    Name           string `db:"name"`
    FailedAttempts int64  `db:"failed_attempts"`
    LockedUntil    int64  `db:"locked_until"`
}
```

#### Article

```go
type Article struct {
    ID        int64    `db:"id"`
    UserID    int64    `db:"user_id"`
    Title     string   `db:"title"`
    Content   string   `db:"content"`
    UpdatedAt int64    `db:"updated_at"`
    CreatedAt int64    `db:"created_at"`
}
```

#### Sijagur Data Models

- `DeRankingOpd`: Main ranking data with capaian/kumulatif/periodik scores
- `DeDetailBarjas`: Barjas procurement details
- `DeDetailFisik`: Physical progress details
- `DeDetailAnggaran`: Budget realization details
- `DeDetailKinerja`: Performance indicator details

### Response Structures

- **Consistent JSON Format**: All responses follow `{message, data/results, meta}` pattern
- **Error Responses**: Include error type and user-friendly messages
- **Pagination**: Meta object with total counts
- **Formatted Data**: Numbers and currencies formatted for display

## Validation & Error Handling

### Form Validation

- **Go Validator**: Uses `github.com/go-playground/validator/v10`
- **Custom Rules**: `fullName` validation for user names
- **Error Messages**: Localized error messages per field
- **Binding**: JSON/form binding with validation tags

### Error Response Patterns

```json
{
  "message": "Please enter your email",
  "error": "VALIDATION_ERROR"
}
```

### HTTP Status Codes

- `200`: Success
- `400`: Bad Request / Validation Error
- `401`: Unauthorized
- `403`: Forbidden / Insufficient Permissions
- `404`: Not Found
- `406`: Not Acceptable
- `500`: Internal Server Error

## Database Schema

### Main Tables

- `user`: User accounts with security fields
- `roles`: User roles (admin, user, etc.)
- `permissions`: System permissions
- `user_roles`: User-role relationships
- `role_permissions`: Role-permission relationships
- `article`: User articles
- `login_attempts`: Security logging

### Sijagur Tables

- `de_ranking_opd`: Main performance rankings
- `de_detail_barjas`: Procurement details
- `de_detail_fisik`: Physical progress details
- `de_detail_anggaran`: Budget details
- `de_detail_kinerja`: Performance details
- `de_peta_detail`: Map/geospatial data
- `de_peta_kecamatan`: District data

### Migrations

- **Versioned Migrations**: Sequential migration system
- **Up/Down Functions**: Reversible migrations
- **Auto-creation**: Tables created if not exist

## Dependencies & Libraries

### Core Dependencies

- `github.com/gin-gonic/gin`: Web framework
- `github.com/lib/pq`: PostgreSQL driver
- `github.com/go-gorp/gorp`: ORM
- `github.com/go-redis/redis/v7`: Redis client
- `github.com/golang-jwt/jwt/v4`: JWT handling
- `github.com/go-playground/validator/v10`: Validation
- `github.com/swaggo/swag`: API documentation

### Utility Libraries

- `github.com/google/uuid`: UUID generation
- `github.com/joho/godotenv`: Environment variables
- `github.com/gin-contrib/gzip`: Response compression
- `golang.org/x/crypto/bcrypt`: Password hashing

## Code Patterns & Standards

### Controller Patterns

```go
func (ctrl ControllerName) ActionName(c *gin.Context) {
    // 1. Extract user ID from context
    userID := getUserID(c)

    // 2. Bind and validate request
    var form FormStruct
    if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
        message := formValidator.Validate(validationErr)
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message})
        return
    }

    // 3. Call model method
    result, err := model.Method(userID, form)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Error message"})
        return
    }

    // 4. Return success response
    c.JSON(http.StatusOK, gin.H{"message": "Success", "data": result})
}
```

### Model Patterns

```go
func (m ModelName) MethodName(params ...interface{}) (result interface{}, err error) {
    // Database operation
    err = db.GetDB().SelectOne(&result, query, params...)
    return result, err
}
```

### Middleware Usage

- **CORS**: Environment-specific origin validation
- **Authentication**: JWT token validation with Redis
- **Authorization**: Permission-based access control
- **Request ID**: Unique identifier per request
- **Gzip**: Response compression

## Configuration & Environment

### Environment Variables

- `ENV`: Environment (PRODUCTION/STAGING/DEVELOPMENT)
- `PORT`: Server port
- `DB_USER`, `DB_PASS`, `DB_NAME`: Database credentials
- `REDIS_HOST`, `REDIS_PASSWORD`: Redis connection
- `ACCESS_SECRET`, `REFRESH_SECRET`: JWT secrets
- `FRONTEND_DOMAIN`: CORS allowed domain
- `SSL`: Enable HTTPS

### Database Connection

- **PostgreSQL**: Primary data store
- **Redis**: Session/token storage
- **Connection Pooling**: Gorp handles connection management

## Testing

### Test Structure

- **Unit Tests**: Located in `tests/` directory
- **Test Framework**: `github.com/stretchr/testify`
- **Coverage**: Focus on model and controller logic

### Example Test

```go
func TestArticleCreate(t *testing.T) {
    // Setup
    form := forms.CreateArticleForm{
        Title: "Test Article",
        Content: "Test content",
    }

    // Execute
    id, err := articleModel.Create(userID, form)

    // Assert
    assert.NoError(t, err)
    assert.Greater(t, id, int64(0))
}
```

## Deployment & Operations

### Build Process

```bash
# Build binary
go build -o main main.go

# Run migrations
./main migrate

# Start server
./main
```

### Docker Support

- **Multi-stage Build**: Optimized for production
- **Environment Configuration**: Container-based config
- **SSL Termination**: Optional SSL support

### Monitoring

- **Request IDs**: Track requests across services
- **Error Logging**: Structured error responses
- **Performance**: Gzip compression, connection pooling

## Areas for Improvement

### Code Quality

1. **Swagger Documentation**: Update docs.go to include Sijagur endpoints
2. **Error Consistency**: Standardize error response formats
3. **Validation**: Add more comprehensive validation rules
4. **Testing**: Increase test coverage, especially for edge cases

### Architecture

1. **Service Layer**: Extract business logic from controllers
2. **Repository Pattern**: More consistent data access patterns
3. **Configuration Management**: Centralized config with validation
4. **Logging**: Structured logging with levels

### Security

1. **Rate Limiting**: Implement request rate limiting
2. **Audit Logging**: Log all security-relevant events
3. **Token Blacklisting**: Implement token revocation
4. **Input Sanitization**: Additional input validation

### Performance

1. **Caching**: Implement Redis caching for frequently accessed data
2. **Database Optimization**: Add indexes, optimize queries
3. **Concurrent Processing**: Use goroutines for independent operations
4. **Response Compression**: Already implemented with gzip

### Maintainability

1. **API Versioning**: Implement proper API versioning
2. **Documentation**: Keep API docs synchronized with code
3. **Code Organization**: Separate concerns more clearly
4. **Dependency Management**: Regular dependency updates

## Conclusion

This API provides a solid foundation for a RESTful backend service with proper authentication, authorization, and data management. The architecture follows Go best practices and provides a scalable base for further development. The Sijagur module adds specialized functionality for performance tracking and ranking systems.

The codebase demonstrates good separation of concerns, consistent error handling, and proper use of Go idioms. Areas for improvement focus on documentation completeness, testing coverage, and architectural enhancements for better maintainability and scalability.
