# Go Web MVC Architecture Documentation

## Overview
This is a comprehensive Go web application built using the Echo framework with a clean MVC (Model-View-Controller) architecture. The application implements a RESTful API with JWT authentication, role-based access control, domain-based authorization, and comprehensive documentation.

## Architecture Layers

### 1. Presentation Layer
- **Echo Framework (v4.13.4)**: High-performance HTTP web framework
  - Middleware support for authentication, logging, CORS, and compression
  - Static file serving for documentation
  - Request/response handling with JSON support
  - Gzip compression with configurable levels
- **Swagger Documentation**: Auto-generated API documentation
  - Interactive testing interface
  - OpenAPI 3.0 specification
  - Embedded Swagger UI
- **Health Check Endpoints**: Application monitoring and status reporting

### 2. Application Layer
- **Handlers**: HTTP request handlers implementing CRUD operations
  - `BaseHandler`: Abstract base class for all handlers
  - `AuthHandler`: Authentication and token management
  - `UserHandler`: User management operations
  - `PostHandler`: Post/content management
  - `RoleHandler`: Role and permission management
  - `DomainHandler`: Domain-specific operations
  - `RegisterHandler`: User registration functionality
- **Routes**: RESTful API route configuration
  - Public routes (login, register, health check, refresh token)
  - Protected routes with middleware chaining
  - Resource-based routing with standard CRUD endpoints
  - API versioning under `/api` prefix
- **Interceptors**: Middleware components
  - JWT authentication middleware with HS512 signing
  - Claims authorization middleware with Redis token validation
  - Casbin RBAC authorization with domain support
  - Resource-level authorization middleware
  - Performance monitoring and request logging
  - Request recovery and error handling

### 3. Business Logic Layer
- **Services**: Business logic implementation
  - `TokenService`: JWT token generation, validation, and refresh management
  - `UserService`: User business logic with domain association
  - `PostService`: Post business logic
  - `RoleService`: Role management logic
  - `DomainService`: Domain-specific business logic
- **Validation**: Request validation using ozzo-validation v4.3.0
- **Request/Response Models**: Structured data transfer objects
- **CLI Commands**: Cobra-based command-line interface

### 4. Data Access Layer
- **GORM (v1.30.1)**: Object-Relational Mapping
  - Support for MySQL, PostgreSQL, and SQLite
  - Automatic migrations and schema management
  - Soft deletes and audit fields
  - Connection pooling and query optimization
- **Database Migrations**: Schema version control
  - Table creation and modification
  - Data seeding and initial setup
  - Multi-database support
- **Redis Cache (v9.12.0)**: In-memory caching
  - Session management with TTL
  - Token storage and invalidation
  - Performance optimization
  - Asynchronous TTL updates

## Security Framework

### Authentication
- **JWT (JSON Web Tokens v5.3.0)**: Token-based authentication
  - Access tokens (30 minutes expiration)
  - Refresh tokens (2 hours expiration)
  - HS512 signing algorithm
  - Token invalidation on logout
  - Redis-based token caching with automatic TTL extension

### Authorization
- **Casbin RBAC (v2.115.0)**: Role-based access control
  - Domain-specific permissions
  - Resource-level authorization
  - Policy-based access control
  - Dynamic permission management
  - Multi-domain user support

### Security Features
- **Password Hashing**: bcrypt for secure password storage
- **Input Validation**: Request sanitization and validation
- **CORS Protection**: Cross-origin resource sharing configuration
- **Rate Limiting**: Request throttling capabilities
- **Domain-based Access Control**: Multi-tenant architecture support

## Request Processing Workflow

### Authentication Flow
1. Client sends credentials to `/login` endpoint
2. Service validates credentials against database
3. JWT access and refresh tokens are generated with HS512 signing
4. Tokens are cached in Redis with expiration
5. Response includes tokens and user information

### Protected Resource Access
1. Client includes JWT token in Authorization header
2. JWT middleware validates token signature and expiration
3. Claims authorization middleware extracts user information and validates Redis token
4. Domain validation middleware checks user-domain association
5. Casbin middleware checks resource permissions for specific domain
6. Resource authorization middleware validates action permissions
7. Handler processes request and returns response

### Token Refresh Flow
1. Client sends refresh token to `/refresh` endpoint
2. Service validates refresh token against Redis
3. New access and refresh tokens are generated
4. Old tokens are invalidated
5. New tokens are cached in Redis

### Data Flow
```
Client Request → Echo Router → Middleware Chain → Handler → Service → Database
                ↓
            Middleware Chain:
            - Recovery
            - Gzip Compression
            - Logger
            - Request ID
            - Performance Monitoring
            - JWT Authentication
            - Claims Authorization
            - Domain Authorization
            - Resource Authorization
```

## Core Libraries & Dependencies

### Web Framework
- **Echo v4.13.4**: High-performance HTTP web framework
- **Echo-JWT v4.3.1**: JWT middleware for Echo framework

### Database & ORM
- **GORM v1.30.1**: Object-Relational Mapping library
- **MySQL Driver v1.6.0**: MySQL database connectivity
- **PostgreSQL Driver v1.6.0**: PostgreSQL database connectivity
- **SQLite Support**: Embedded database support
- **Redis v9.12.0**: In-memory data structure store

### Authentication & Security
- **JWT v5.3.0**: JSON Web Token implementation
- **Casbin v2.115.0**: Authorization library with RBAC support
- **Casbin GORM Adapter v3.36.0**: Database adapter for Casbin
- **Ozzo-Validation v4.3.0**: Request validation library
- **Golang Crypto v0.40.0**: Cryptographic functions

### CLI & Configuration
- **Cobra v1.9.1**: CLI application framework
- **Godotenv v1.5.1**: Environment variable management
- **Zerolog v1.34.0**: Structured logging library

### Documentation
- **Swaggo/Swag v1.16.6**: Swagger documentation generator
- **Swagger UI**: Interactive API documentation

### Utilities
- **UUID v1.6.0**: UUID generation and management
- **Gofakeit v7.3.0**: Data generation for testing
- **Golang Sync v0.16.0**: Concurrency utilities

## Project Structure

```
gowebmvc/
├── api/                    # API response utilities
├── casbin/                 # RBAC policy configuration
├── cmd/                    # CLI command implementations
│   ├── migrate.go         # Database migration commands
│   ├── root.go            # Main CLI entry point
│   └── version.go         # Version information
├── config/                 # Configuration management
│   ├── auth.go            # Authentication configuration
│   ├── config.go          # Main configuration structure
│   ├── db.go              # Database configuration
│   ├── http.go            # HTTP server configuration
│   └── redis.go           # Redis configuration
├── db/                     # Database connection and migrations
│   ├── connection.go      # Database connection setup
│   ├── migrations/        # Database migration files
│   └── migrator.go        # Migration utilities
├── docs/                   # Swagger documentation files
├── handlers/               # HTTP request handlers
│   ├── auth_handler.go    # Authentication handlers
│   ├── base_handler.go    # Base handler interface
│   ├── domain_handler.go  # Domain management
│   ├── post_handler.go    # Post management
│   ├── register_handler.go # User registration
│   ├── role_handler.go    # Role management
│   └── user_handler.go    # User management
├── interceptor/            # Middleware implementations
│   └── middlewares.go     # All middleware functions
├── models/                 # Data models and entities
│   ├── base.go            # Base model structure
│   ├── domain.go          # Domain model
│   ├── post.go            # Post model
│   └── user.go            # User model
├── requests/               # Request DTOs
├── responses/              # Response DTOs
├── routes/                 # Route configuration
│   └── routes.go          # Route setup and middleware
├── server/                 # Server initialization
│   └── server.go          # Server configuration
├── services/               # Business logic services
│   ├── domain_service.go  # Domain business logic
│   ├── post_service.go    # Post business logic
│   ├── role_service.go    # Role business logic
│   ├── token_service.go   # Token management
│   └── user_service.go    # User business logic
└── util/                   # Utility functions
    └── utilities.go       # Common utilities
```

## Key Features

### RESTful API Design
- Standard CRUD operations for all resources
- Consistent response format with structured error handling
- Proper HTTP status codes and error messages
- Resource-based URL structure with UUID identifiers
- API versioning and grouping

### Multi-Database Support
- MySQL, PostgreSQL, and SQLite support
- Database-agnostic ORM layer
- Migration system for schema management
- Connection pooling and optimization

### Domain-Based Architecture
- Multi-tenant support with domain isolation
- Domain-specific user associations
- Cross-domain authorization controls
- Domain header-based routing

### Comprehensive Logging
- Structured logging with Zerolog
- Request/response logging with performance metrics
- Error tracking and debugging
- Slow request detection and alerting

### Development Tools
- Hot reloading capabilities
- Comprehensive testing framework
- Code generation utilities
- Development and production configurations
- CLI-based database migrations

### Deployment Ready
- Docker containerization with multi-stage builds
- Environment-based configuration
- Health check endpoints with database and Redis monitoring
- Graceful shutdown handling
- Performance monitoring and metrics

## Performance Optimizations

### Caching Strategy
- Redis-based session caching with TTL management
- Token caching with automatic expiration extension
- Database query optimization with GORM
- Response compression (Gzip) with configurable levels

### Database Optimization
- Connection pooling with configurable limits
- Query optimization and indexing
- Migration efficiency
- Multi-database support

### Server Optimization
- Request timeout configuration
- Memory management and garbage collection
- Concurrent request handling
- Resource cleanup and monitoring

## Monitoring & Observability

### Health Checks
- Database connectivity monitoring with ping tests
- Redis connectivity monitoring
- Application status endpoints with timestamps
- Performance metrics collection

### Logging
- Structured JSON logging with Zerolog
- Request tracing with unique request IDs
- Error tracking and debugging
- Performance monitoring with duration tracking

### Metrics
- Request/response times with percentile tracking
- Error rates and status code distribution
- Resource utilization monitoring
- Custom business metrics
- Slow request detection (>1 second threshold)

## Security Enhancements

### Token Management
- Automatic token TTL extension in Redis
- Token invalidation on logout
- Refresh token rotation
- Domain-specific token validation

### Authorization Layers
- Multi-level authorization checks
- Domain-based access control
- Resource-level permissions
- Action-based authorization

### Input Validation
- Request sanitization
- Structured validation with ozzo-validation
- Type-safe request/response models
- Error message standardization

This architecture provides a solid foundation for building scalable, secure, and maintainable web applications with Go, featuring modern security practices, comprehensive monitoring, and multi-tenant support.
