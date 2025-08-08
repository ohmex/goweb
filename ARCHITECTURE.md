# Go Web MVC Architecture Documentation

## Overview
This is a comprehensive Go web application built using the Echo framework with a clean MVC (Model-View-Controller) architecture. The application implements a RESTful API with JWT authentication, role-based access control, and comprehensive documentation.

## Architecture Layers

### 1. Presentation Layer
- **Echo Framework (v4)**: High-performance HTTP web framework
  - Middleware support for authentication, logging, and CORS
  - Static file serving for documentation
  - Request/response handling with JSON support
- **Swagger Documentation**: Auto-generated API documentation
  - Interactive testing interface
  - OpenAPI 3.0 specification
  - Embedded Swagger UI

### 2. Application Layer
- **Handlers**: HTTP request handlers implementing CRUD operations
  - `BaseHandler`: Abstract base class for all handlers
  - `AuthHandler`: Authentication and token management
  - `UserHandler`: User management operations
  - `PostHandler`: Post/content management
  - `RoleHandler`: Role and permission management
  - `DomainHandler`: Domain-specific operations
- **Routes**: RESTful API route configuration
  - Public routes (login, register, health check)
  - Protected routes with middleware chaining
  - Resource-based routing with standard CRUD endpoints
- **Interceptors**: Middleware components
  - JWT authentication middleware
  - Claims authorization middleware
  - Casbin RBAC authorization
  - Performance monitoring
  - Request logging and recovery

### 3. Business Logic Layer
- **Services**: Business logic implementation
  - `TokenService`: JWT token generation and validation
  - `UserService`: User business logic
  - `PostService`: Post business logic
  - `RoleService`: Role management logic
  - `DomainService`: Domain-specific business logic
- **Validation**: Request validation using ozzo-validation
- **Request/Response Models**: Structured data transfer objects

### 4. Data Access Layer
- **GORM**: Object-Relational Mapping
  - Support for MySQL, PostgreSQL, and SQLite
  - Automatic migrations and schema management
  - Soft deletes and audit fields
- **Database Migrations**: Schema version control
  - Table creation and modification
  - Data seeding and initial setup
- **Redis Cache**: In-memory caching
  - Session management
  - Token storage and invalidation
  - Performance optimization

## Security Framework

### Authentication
- **JWT (JSON Web Tokens)**: Token-based authentication
  - Access tokens (30 minutes expiration)
  - Refresh tokens (2 hours expiration)
  - Token invalidation on logout
  - Redis-based token caching

### Authorization
- **Casbin RBAC**: Role-based access control
  - Domain-specific permissions
  - Resource-level authorization
  - Policy-based access control
  - Dynamic permission management

### Security Features
- **Password Hashing**: bcrypt for secure password storage
- **Input Validation**: Request sanitization and validation
- **CORS Protection**: Cross-origin resource sharing configuration
- **Rate Limiting**: Request throttling capabilities

## Request Processing Workflow

### Authentication Flow
1. Client sends credentials to `/login` endpoint
2. Service validates credentials against database
3. JWT access and refresh tokens are generated
4. Tokens are cached in Redis with expiration
5. Response includes tokens and user information

### Protected Resource Access
1. Client includes JWT token in Authorization header
2. JWT middleware validates token signature and expiration
3. Claims authorization middleware extracts user information
4. Casbin middleware checks resource permissions
5. Handler processes request and returns response

### Data Flow
```
Client Request → Echo Router → Middleware Chain → Handler → Service → Database
```

## Core Libraries & Dependencies

### Web Framework
- **Echo v4.13.4**: High-performance HTTP web framework
- **Echo-JWT v4.3.1**: JWT middleware for Echo framework

### Database & ORM
- **GORM v1.30.1**: Object-Relational Mapping library
- **MySQL Driver v1.6.0**: MySQL database connectivity
- **PostgreSQL Driver v1.6.0**: PostgreSQL database connectivity
- **Redis v9.12.0**: In-memory data structure store

### Authentication & Security
- **JWT v5.3.0**: JSON Web Token implementation
- **Casbin v2.115.0**: Authorization library with RBAC support
- **Ozzo-Validation v4.3.0**: Request validation library
- **Golang Crypto**: Cryptographic functions

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
- **Golang Sync**: Concurrency utilities

## Project Structure

```
gowebmvc/
├── api/                    # API response utilities
├── casbin/                 # RBAC policy configuration
├── cmd/                    # CLI command implementations
├── config/                 # Configuration management
├── db/                     # Database connection and migrations
├── docs/                   # Swagger documentation files
├── handlers/               # HTTP request handlers
├── interceptor/            # Middleware implementations
├── models/                 # Data models and entities
├── requests/               # Request DTOs
├── responses/              # Response DTOs
├── routes/                 # Route configuration
├── server/                 # Server initialization
├── services/               # Business logic services
└── util/                   # Utility functions
```

## Key Features

### RESTful API Design
- Standard CRUD operations for all resources
- Consistent response format
- Proper HTTP status codes
- Resource-based URL structure

### Multi-Database Support
- MySQL, PostgreSQL, and SQLite support
- Database-agnostic ORM layer
- Migration system for schema management

### Comprehensive Logging
- Structured logging with Zerolog
- Request/response logging
- Performance monitoring
- Error tracking and debugging

### Development Tools
- Hot reloading capabilities
- Comprehensive testing framework
- Code generation utilities
- Development and production configurations

### Deployment Ready
- Docker containerization
- Environment-based configuration
- Health check endpoints
- Graceful shutdown handling

## Performance Optimizations

### Caching Strategy
- Redis-based session caching
- Token caching with expiration
- Database query optimization
- Response compression (Gzip)

### Database Optimization
- Connection pooling
- Query optimization
- Index management
- Migration efficiency

### Server Optimization
- Request timeout configuration
- Memory management
- Concurrent request handling
- Resource cleanup

## Monitoring & Observability

### Health Checks
- Database connectivity monitoring
- Redis connectivity monitoring
- Application status endpoints
- Performance metrics

### Logging
- Structured JSON logging
- Request tracing
- Error tracking
- Performance monitoring

### Metrics
- Request/response times
- Error rates
- Resource utilization
- Custom business metrics

This architecture provides a solid foundation for building scalable, secure, and maintainable web applications with Go.
