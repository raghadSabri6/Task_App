package router

import (
	"log"
	"net/http"
	"strings"
	"task2/internal/adapter/controller"
	"task2/internal/app/dto"
	"task2/internal/infrastructure/middleware"
)

// Router handles HTTP routing
type Router struct {
	mux               *http.ServeMux
	authMiddleware    *middleware.AuthMiddleware
	loggingMiddleware *middleware.LoggingMiddleware
	corsMiddleware    *middleware.CorsMiddleware
	logger            *log.Logger
}

// NewRouter creates a new router
func NewRouter(authMiddleware *middleware.AuthMiddleware) *Router {
	return &Router{
		mux:            http.NewServeMux(),
		authMiddleware: authMiddleware,
		logger:         log.New(log.Writer(), "[ROUTER] ", log.LstdFlags),
	}
}

// SetLoggingMiddleware sets the logging middleware
func (r *Router) SetLoggingMiddleware(loggingMiddleware *middleware.LoggingMiddleware) {
	r.loggingMiddleware = loggingMiddleware
}

// SetCorsMiddleware sets the CORS middleware
func (r *Router) SetCorsMiddleware(corsMiddleware *middleware.CorsMiddleware) {
	r.corsMiddleware = corsMiddleware
}

// RegisterUserRoutes registers user routes
func (r *Router) RegisterUserRoutes(userController *controller.UserController) {
	r.logger.Println("Registering user routes")

	// Register handler with debug logging
	r.mux.Handle("/api/v1/register", r.wrapHandler(
		middleware.MethodCheck("POST")(
			middleware.BindAndValidate(&dto.CreateUserRequest{})(
				http.HandlerFunc(userController.Register)))))

	// Login handler with debug logging
	r.mux.Handle("/api/v1/login", r.wrapHandler(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Set CORS headers for login specifically
			origin := req.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if req.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.WriteHeader(http.StatusOK)
				return
			}

			// Use the middleware chain for POST requests
			if req.Method == "POST" {
				middleware.BindAndValidate(&dto.LoginRequest{})(
					http.HandlerFunc(userController.Login)).ServeHTTP(w, req)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	// Profile handler
	r.mux.Handle("/api/v1/profile", r.wrapHandler(
		r.authMiddleware.Middleware(
			middleware.MethodCheck("GET")(
				http.HandlerFunc(userController.GetProfile)))))

	// Users handler
	r.mux.Handle("/api/v1/users", r.wrapHandler(
		r.authMiddleware.Middleware(
			middleware.MethodCheck("GET")(
				http.HandlerFunc(userController.GetAllUsers)))))

	// Add a debug cookie endpoint
	r.mux.Handle("/api/v1/debug/cookie", r.wrapHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Debug cookie endpoint called: Method=%s", r.Method)

			// Set CORS headers
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Set a test cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "DebugCookie",
				Value:    "test-value",
				Path:     "/",
				HttpOnly: false,
				MaxAge:   3600,
			})

			// Log all cookies in the request
			cookies := r.Cookies()
			log.Printf("Found %d cookies in request", len(cookies))
			for _, cookie := range cookies {
				log.Printf("Cookie: %s=%s", cookie.Name, cookie.Value)
			}

			// Return response with cookies
			http.HandlerFunc(userController.DebugCookie).ServeHTTP(w, r)
		})))
}

// RegisterTaskRoutes registers task routes
func (r *Router) RegisterTaskRoutes(taskController *controller.TaskController) {
	r.logger.Println("Registering task routes")

	// Create task handler and Get all tasks handler
	r.mux.Handle("/api/v1/tasks", r.wrapHandler(
		r.authMiddleware.Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" {
					// Create task using BindAndValidate middleware
					middleware.BindAndValidate(&dto.CreateTaskRequest{})(
						http.HandlerFunc(taskController.CreateTask)).ServeHTTP(w, r)
				} else if r.Method == "GET" {
					// Get all tasks
					taskController.GetAllTasks(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}))))

	// Get tasks created by user handler
	r.mux.Handle("/api/v1/tasks/created", r.wrapHandler(
		r.authMiddleware.Middleware(
			middleware.MethodCheck("GET")(
				http.HandlerFunc(taskController.GetTasksCreatedByUser)))))

	// Get tasks assigned to user handler
	r.mux.Handle("/api/v1/tasks/assigned", r.wrapHandler(
		r.authMiddleware.Middleware(
			middleware.MethodCheck("GET")(
				http.HandlerFunc(taskController.GetTasksAssignedToUser)))))

	// Get task by ID, Delete task, Complete task, and Assign task handlers
	r.mux.Handle("/api/v1/tasks/", r.wrapHandler(
		r.authMiddleware.Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					taskController.GetTaskByID(w, r)
				case "DELETE":
					taskController.DeleteTask(w, r)
				case "PUT":
					if len(r.URL.Path) > 16 && r.URL.Path[len(r.URL.Path)-9:] == "/complete" {
						taskController.CompleteTask(w, r)
					} else if strings.Contains(r.URL.Path, "/assign/") {
						taskController.AssignTask(w, r)
					} else {
						http.NotFound(w, r)
					}
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}))))
}

// wrapHandler wraps a handler with the logging middleware if available
func (r *Router) wrapHandler(handler http.Handler) http.Handler {
	// Apply CORS middleware if available
	if r.corsMiddleware != nil {
		handler = r.corsMiddleware.Middleware(handler)
	}

	// Apply logging middleware if available
	if r.loggingMiddleware != nil {
		handler = r.loggingMiddleware.Middleware(handler)
	}

	return handler
}

// GetHandler returns the HTTP handler
func (r *Router) GetHandler() http.Handler {
	return r.mux
}