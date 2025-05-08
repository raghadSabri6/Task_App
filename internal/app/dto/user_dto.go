package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the response for a user
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// UserSummary represents a simplified user response
type UserSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// LoginResponse represents the response for a login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UsersResponse represents the response for multiple users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}