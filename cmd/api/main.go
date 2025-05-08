package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"task2/internal/adapter/controller"
	"task2/internal/adapter/repository"
	"task2/internal/app/usecase"
	"task2/internal/domain/service"
	"task2/internal/infrastructure/auth"
	"task2/internal/infrastructure/config"
	"task2/internal/infrastructure/dependencies"
	"task2/internal/infrastructure/middleware"
	"task2/internal/infrastructure/router"
)

func main() {
	// Create logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)
	
	// Load configuration
	logger.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Log configuration for debugging
	logger.Printf("Database URL: %s", maskConnectionString(cfg.DatabaseURL))
	logger.Printf("JWT Secret: %s", maskString(cfg.JWTSecret))
	logger.Printf("Port: %s", cfg.Port)
	
	// Initialize dependencies
	logger.Println("Initializing dependencies...")
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer deps.Close()
	
	// Create repositories
	logger.Println("Creating repositories...")
	userRepo := repository.NewUserRepository(deps.DB)
	taskRepo := repository.NewTaskRepository(deps.DB)
	
	// Create domain services
	logger.Println("Creating domain services...")
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	
	// Create auth service
	logger.Println("Creating auth service...")
	authService := auth.NewAuthService(cfg.JWTSecret)
	
	// Create use cases
	logger.Println("Creating use cases...")
	userUseCase := usecase.NewUserUseCase(userService, authService)
	userUseCase.SetEmailService(deps.EmailClient)
	taskUseCase := usecase.NewTaskUseCase(taskService, userService)
	
	// Create controllers
	logger.Println("Creating controllers...")
	userController := controller.NewUserController(userUseCase)
	taskController := controller.NewTaskController(taskUseCase)
	
	// Create middleware
	logger.Println("Creating middleware...")
	authMiddleware := middleware.NewAuthMiddleware(authService)
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	corsMiddleware := middleware.NewCorsMiddleware(logger)
	
	// Create router
	logger.Println("Setting up router...")
	r := router.NewRouter(authMiddleware)
	r.SetLoggingMiddleware(loggingMiddleware)
	r.SetCorsMiddleware(corsMiddleware)
	
	// Register routes
	logger.Println("Registering routes...")
	r.RegisterUserRoutes(userController)
	r.RegisterTaskRoutes(taskController)
	
	// Create server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r.GetHandler(),
	}
	
	// Start server in a goroutine
	go func() {
		logger.Printf("Server starting on http://localhost:%s", port)
		logger.Printf("CORS is enabled with credentials support")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Println("Server shutting down...")
	
	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}
	
	logger.Println("Server exited properly")
}

// maskConnectionString masks a database connection string for logging
func maskConnectionString(connStr string) string {
	if len(connStr) < 20 {
		return "***"
	}
	return connStr[:10] + "..." + connStr[len(connStr)-10:]
}

// maskString masks a string for logging
func maskString(s string) string {
	if len(s) < 8 {
		return "***"
	}
	return s[:2] + "..." + s[len(s)-2:]
}