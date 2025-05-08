package service

import (
	"context"
	"errors"
	"task2/internal/domain/entity"
	"task2/internal/domain/repository"

	"github.com/google/uuid"
)

// UserService provides domain logic for users
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *entity.User) error {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(ctx, user.Email)
	if err != nil {
		return err
	}
	
	if exists {
		return errors.New("email already registered")
	}
	
	return s.userRepo.Create(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetUserByUUID retrieves a user by UUID
func (s *UserService) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	return s.userRepo.GetByUUID(ctx, uuid)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return s.userRepo.GetAll(ctx)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) error {
	// Check if user exists
	_, err := s.userRepo.GetByUUID(ctx, user.UUID)
	if err != nil {
		return errors.New("user not found")
	}
	
	// Check if email is already taken by another user
	if user.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && existingUser.UUID != user.UUID {
			return errors.New("email already registered by another user")
		}
	}
	
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return errors.New("user not found")
	}
	
	return s.userRepo.Delete(ctx, uuid)
}