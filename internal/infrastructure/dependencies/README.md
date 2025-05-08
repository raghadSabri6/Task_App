# Dependencies Package

This package centralizes all external dependencies used by the application, such as:

- Database connections
- Email services
- Redis clients (future)
- Message queues (future)
- External API clients (future)
- Cache services (future)

## Purpose

The dependencies package serves several important purposes:

1. **Single Responsibility**: Centralizes the initialization and management of all external services
2. **Dependency Injection**: Provides a clean way to inject dependencies into application components
3. **Testing**: Makes it easier to mock external dependencies for testing
4. **Configuration**: Handles configuration of all external services in one place
5. **Lifecycle Management**: Manages the lifecycle (connect, disconnect) of all external services

## Usage

To use the dependencies package:

1. Initialize dependencies in your application entry point:

```go
// Load configuration
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}

// Initialize dependencies
deps, err := dependencies.NewDependencies(cfg)
if err != nil {
    log.Fatalf("Failed to initialize dependencies: %v", err)
}
defer deps.Close()

// Use dependencies
repo := repository.NewRepository(deps.DB)
```

2. Add new dependencies by extending the `Dependencies` struct and adding initialization in the `NewDependencies` function.

## Adding New Dependencies

To add a new external dependency:

1. Add the dependency client to the `Dependencies` struct
2. Create an initialization function for the dependency
3. Initialize the dependency in the `NewDependencies` function
4. Add cleanup code in the `Close` method

Example for adding Redis:

```go
// In dependencies.go
type Dependencies struct {
    DB          *bun.DB
    EmailClient *email.EmailService
    RedisClient *redis.Client  // New dependency
}

// Initialize Redis
func initRedis(cfg *config.Config) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisAddr,
        Password: cfg.RedisPassword,
        DB:       0,
    })
    
    // Test connection
    if _, err := client.Ping(context.Background()).Result(); err != nil {
        return nil, err
    }
    
    return client, nil
}

// In NewDependencies
redisClient, err := initRedis(cfg)
if err != nil {
    return nil, fmt.Errorf("failed to initialize Redis: %w", err)
}

// In Close method
if d.RedisClient != nil {
    if err := d.RedisClient.Close(); err != nil {
        errs = append(errs, fmt.Errorf("error closing Redis: %w", err))
    }
}
```