package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description"`
	Users       []UserAssign `json:"users,omitempty"`
}

// UserAssign represents a user to be assigned to a task
type UserAssign struct {
	ID string `json:"id" validate:"required"`
}

// TaskResponse represents the response for a task
type TaskResponse struct {
	ID          uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Completed   bool        `json:"completed"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `json:"deleted_at,omitempty"`
	CreatedBy   UserSummary `json:"created_by"`
	AssignedTo  *UserSummary `json:"assigned_to,omitempty"`
	Users       []UserSummary `json:"users,omitempty"`
}

// TasksResponse represents the response for multiple tasks
type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

// AssignTaskRequest represents the request to assign a task
type AssignTaskRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// CompleteTaskRequest represents the request to complete a task
type CompleteTaskRequest struct {
	// Empty as it's just a status change
}