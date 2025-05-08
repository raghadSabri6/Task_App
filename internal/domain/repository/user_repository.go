package repository

import (
	"context"
	"task2/internal/domain/entity"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create a new user
	Create(ctx context.Context, user *entity.User) error
	
	// Get a user by UUID
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	
	// Get a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// Get all users
	GetAll(ctx context.Context) ([]*entity.User, error)
	
	// Update an existing user
	Update(ctx context.Context, user *entity.User) error
	
	// Delete a user
	Delete(ctx context.Context, uuid uuid.UUID) error
	
	// Check if email exists
	EmailExists(ctx context.Context, email string) (bool, error)
}