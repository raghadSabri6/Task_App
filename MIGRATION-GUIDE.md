# Migration Guide

This document provides guidance on migrating from the old project structure to the new Clean Architecture structure.

## Overview

The project has been refactored to follow Clean Architecture principles and implement the Repository Design Pattern. This guide will help you understand the changes and how to migrate your code.

## Directory Structure Changes

### Old Structure

```
├── controllers/
├── database/
├── helperFunc/
├── images/
├── initializers/
├── middlewares/
├── migrations/
├── models/
├── routes/
├── schemas/
├── templates/
└── main.go
```

### New Structure

```
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   ├── repository/
│   │   └── service/
│   ├── app/
│   │   ├── dto/
│   │   └── usecase/
│   ├── adapter/
│   │   ├── controller/
│   │   ├── presenter/
│   │   └── repository/
│   └── infrastructure/
│       ├── auth/
│       ├── config/
│       ├── database/
│       ├── middleware/
│       ├── persistence/
│       ├── router/
│       └── validator/
├── pkg/
│   ├── email/
│   └── utils/
├── migrations/
└── templates/
```

## Component Mapping

| Old Component | New Component |
|---------------|--------------|
| `controllers/` | `internal/adapter/controller/` |
| `database/` | `internal/infrastructure/database/` |
| `helperFunc/` | `pkg/utils/` |
| `initializers/` | `internal/infrastructure/config/` |
| `middlewares/` | `internal/infrastructure/middleware/` |
| `models/` | `internal/domain/entity/` and `internal/infrastructure/persistence/` |
| `routes/` | `internal/infrastructure/router/` |
| `schemas/` | `internal/app/dto/` |
| `main.go` | `cmd/api/main.go` |

## Migration Steps

1. **Domain Layer**:
   - Move business entities from `models/` to `internal/domain/entity/`
   - Create repository interfaces in `internal/domain/repository/`
   - Implement domain services in `internal/domain/service/`

2. **Application Layer**:
   - Move request/response schemas from `schemas/` to `internal/app/dto/`
   - Implement use cases in `internal/app/usecase/`

3. **Interface Adapters Layer**:
   - Move controllers from `controllers/` to `internal/adapter/controller/`
   - Implement repository implementations in `internal/adapter/repository/`
   - Create presenters in `internal/adapter/presenter/`

4. **Infrastructure Layer**:
   - Move database connection from `database/` to `internal/infrastructure/database/`
   - Move middleware from `middlewares/` to `internal/infrastructure/middleware/`
   - Move initialization code from `initializers/` to `internal/infrastructure/config/`
   - Create persistence models in `internal/infrastructure/persistence/`
   - Implement router in `internal/infrastructure/router/`

5. **Shared Utilities**:
   - Move helper functions from `helperFunc/` to `pkg/utils/`
   - Create email service in `pkg/email/`

## Database Changes

No database schema changes are required. The new persistence models map to the same database tables as the old models.

## API Changes

The API endpoints remain the same, but the implementation has been refactored to follow Clean Architecture principles.

## Testing

After migration, run the following tests:

1. Unit tests for domain services and use cases
2. Integration tests for repositories and controllers
3. End-to-end tests for API endpoints