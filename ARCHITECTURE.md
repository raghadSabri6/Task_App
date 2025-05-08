# Clean Architecture Implementation

This document describes the Clean Architecture implementation of the Task Management API.

## Overview

The application follows the principles of Clean Architecture as described by Robert C. Martin. The architecture is divided into the following layers:

1. **Domain Layer** (innermost layer)
2. **Application Layer**
3. **Interface Adapters Layer**
4. **Infrastructure Layer** (outermost layer)

## Layer Dependencies

Dependencies flow from the outer layers to the inner layers. Inner layers have no knowledge of outer layers.

```
Infrastructure → Interface Adapters → Application → Domain
```

## Layer Descriptions

### Domain Layer (`internal/domain/`)

The Domain Layer contains enterprise business rules and entities. It is independent of any external frameworks or libraries.

- **Entities**: Core business objects (e.g., User, Task)
- **Repository Interfaces**: Define data access contracts
- **Domain Services**: Implement core business logic

### Application Layer (`internal/app/`)

The Application Layer contains application-specific business rules. It orchestrates the flow of data and business rules.

- **Use Cases**: Implement application-specific business rules
- **DTOs**: Data Transfer Objects for input/output data transformation

### Interface Adapters Layer (`internal/adapter/`)

The Interface Adapters Layer converts data between the Application Layer and external frameworks or tools.

- **Controllers**: Handle HTTP requests and responses
- **Presenters**: Format data for presentation
- **Repository Implementations**: Implement repository interfaces

### Infrastructure Layer (`internal/infrastructure/`)

The Infrastructure Layer contains frameworks, tools, and drivers.

- **Database**: Database connection and models
- **Authentication**: Authentication middleware and services
- **Configuration**: Application configuration
- **Router**: HTTP routing
- **Middleware**: HTTP middleware

## Repository Pattern

The Repository Pattern is implemented to abstract data access:

1. **Repository Interfaces** (`internal/domain/repository/`): Define methods for data access
2. **Repository Implementations** (`internal/adapter/repository/`): Implement repository interfaces
3. **Persistence Models** (`internal/infrastructure/persistence/`): Database models

## Dependency Injection

Dependencies are injected from the main function to ensure loose coupling and testability.

## Error Handling

Errors are propagated up the call stack and handled at the appropriate level:

1. **Domain Layer**: Returns domain-specific errors
2. **Application Layer**: Translates domain errors to application errors
3. **Interface Adapters Layer**: Translates application errors to HTTP responses
4. **Infrastructure Layer**: Handles infrastructure-specific errors

## Testing

Each layer can be tested independently:

1. **Domain Layer**: Unit tests for business rules
2. **Application Layer**: Unit tests for use cases
3. **Interface Adapters Layer**: Integration tests for controllers and repositories
4. **Infrastructure Layer**: Integration tests for infrastructure components