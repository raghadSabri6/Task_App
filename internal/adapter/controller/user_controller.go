package controller

import (
	"log"
	"net/http"
	"time"
	"task2/internal/app/dto"
	"task2/internal/app/usecase"
	"task2/internal/infrastructure/middleware"
	"task2/pkg/utils"
)

// UserController handles HTTP requests for users
type UserController struct {
	userUseCase *usecase.UserUseCase
}

// NewUserController creates a new user controller
func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// Register handles user registration
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	// Get request body from context
	ctx := r.Context()
	val := ctx.Value(middleware.BindKey)
	
	// Log the type of the value for debugging
	log.Printf("Register: Value type from context: %T", val)
	
	// Try to cast to CreateUserRequest
	userReq, ok := val.(*dto.CreateUserRequest)
	if !ok {
		log.Printf("Register: Type assertion failed, value is not *dto.CreateUserRequest")
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	
	log.Printf("Register: User request: %+v", userReq)
	
	// Create user
	user, err := c.userUseCase.CreateUser(ctx, userReq)
	if err != nil {
		log.Printf("Register: Failed to create user: %v", err)
		utils.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusCreated, "", map[string]interface{}{"user": user})
}

// Login handles user login
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// Log request details
	log.Printf("Login request received: Method=%s, ContentType=%s", r.Method, r.Header.Get("Content-Type"))
	
	// Get request body from context
	ctx := r.Context()
	val := ctx.Value(middleware.BindKey)
	
	// Log the type of the value for debugging
	log.Printf("Login: Value type from context: %T", val)
	
	// Try to cast to LoginRequest
	loginReq, ok := val.(*dto.LoginRequest)
	if !ok {
		log.Printf("Login: Type assertion failed, value is not *dto.LoginRequest")
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	
	log.Printf("Login: Login request: %+v", loginReq)
	
	// Login user
	loginResp, err := c.userUseCase.Login(ctx, loginReq)
	if err != nil {
		log.Printf("Login: Failed to login: %v", err)
		utils.RespondJSON(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	
	// Set token in cookie
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    loginResp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	http.SetCookie(w, cookie)
	
	log.Printf("Login: Set Authorization cookie with token")
	
	// Also set a test cookie to verify cookie functionality
	testCookie := &http.Cookie{
		Name:     "TestCookie",
		Value:    "test-value",
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, testCookie)
	log.Printf("Login: Set test cookie")
	
	// Set CORS headers to allow credentials
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	
	// Send response
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{
		"user":  loginResp.User,
		"token": loginResp.Token,
	})
}

// GetProfile handles getting the user's profile
func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Get user
	user, err := c.userUseCase.GetUserByUUID(r.Context(), userUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusNotFound, "User not found", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"user": user})
}

// GetAllUsers handles getting all users
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get all users
	usersResp, err := c.userUseCase.GetAllUsers(r.Context())
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch users", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"users": usersResp.Users})
}

// DebugCookie handles the debug cookie endpoint
func (c *UserController) DebugCookie(w http.ResponseWriter, r *http.Request) {
	// Return response
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{
		"message": "Debug cookie set",
		"cookies": r.Cookies(),
	})
}