package repository

import (
	"context"
	"task2/internal/domain/entity"

	"github.com/google/uuid"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	// Create a new task
	Create(ctx context.Context, task *entity.Task) error
	
	// Get a task by its UUID
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Task, error)
	
	// Get all tasks
	GetAll(ctx context.Context) ([]*entity.Task, error)
	
	// Update an existing task
	Update(ctx context.Context, task *entity.Task) error
	
	// Delete a task
	Delete(ctx context.Context, uuid uuid.UUID) error
	
	// Get tasks created by a specific user
	GetTasksCreatedByUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error)
	
	// Get tasks assigned to a specific user
	GetTasksAssignedToUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error)
	
	// Assign a task to a user
	AssignTaskToUser(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error
	
	// Complete a task
	CompleteTask(ctx context.Context, taskUUID uuid.UUID) error
	
	// Add a user to a task
	AddUserToTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error
}