# Clean Architecture Task Management API

This project has been refactored to follow Clean Architecture principles and the Repository Design Pattern.

## Project Structure

```
├── cmd/                  # Application entry points
│   └── api/              # API server
│       └── main.go       # Main application
├── internal/             # Private application code
│   ├── domain/           # Enterprise business rules (entities)
│   │   ├── entity/       # Domain entities
│   │   ├── repository/   # Repository interfaces
│   │   └── service/      # Domain services
│   ├── app/              # Application business rules
│   │   ├── dto/          # Data Transfer Objects
│   │   └── usecase/      # Use cases (application logic)
│   ├── adapter/          # Interface adapters
│   │   ├── controller/   # HTTP controllers
│   │   ├── presenter/    # Response formatters
│   │   └── repository/   # Repository implementations
│   └── infrastructure/   # Frameworks & drivers
│       ├── auth/         # Authentication
│       ├── config/       # Configuration
│       ├── database/     # Database connection
│       ├── dependencies/ # Application dependencies
│       ├── middleware/   # HTTP middleware
│       ├── persistence/  # Database models
│       ├── router/       # HTTP router
│       └── validator/    # Validation
├── pkg/                  # Public libraries
│   ├── email/            # Email utilities
│   └── utils/            # Utility functions
├── migrations/           # Database migrations
└── templates/            # Email templates
```

## Architecture Overview

1. **Domain Layer**: Contains business entities and repository interfaces
2. **Application Layer**: Contains use cases and DTOs
3. **Interface Adapters Layer**: Contains controllers, presenters, and repository implementations
4. **Infrastructure Layer**: Contains frameworks and drivers

## Repository Pattern

The Repository Pattern is implemented to abstract data access:
- Repository interfaces are defined in the domain layer
- Repository implementations are in the adapter layer
- Use cases depend on repository interfaces, not implementations

## Running the Application

1. Set up environment variables in `.env` file:
```
DATABASE_URL=postgres://username:password@localhost:5432/dbname
JWT_SECRET=your_jwt_secret
PORT=8080
```

2. Build and run the application:
```
go build ./cmd/api
./api
```

## API Endpoints

### User Endpoints
- `POST /register` - Register a new user
- `POST /login` - Login a user
- `GET /profile` - Get user profile
- `GET /users` - Get all users

### Task Endpoints
- `POST /tasks` - Create a new task
- `GET /tasks` - Get all tasks
- `GET /tasks/{id}` - Get a task by ID
- `DELETE /tasks/{id}` - Delete a task
- `PUT /tasks/{id}/complete` - Complete a task
- `PUT /tasks/{id}/assign/{userId}` - Assign a task to a user
- `GET /tasks/created` - Get tasks created by the current user
- `GET /tasks/assigned` - Get tasks assigned to the current user
