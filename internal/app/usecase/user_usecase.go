package usecase

import (
	"context"
	"errors"
	"task2/internal/adapter/presenter"
	"task2/internal/app/dto"
	"task2/internal/domain/entity"
	"task2/internal/domain/service"
	"task2/pkg/email"

	"github.com/google/uuid"
)

// UserUseCase handles application logic for users
type UserUseCase struct {
	userService   *service.UserService
	authService   AuthService
	emailService  *email.EmailService
	userPresenter *presenter.UserPresenter
}

// AuthService defines the interface for authentication
type AuthService interface {
	GenerateToken(userUUID uuid.UUID) (string, error)
	ValidatePassword(hashedPassword, password string) bool
	HashPassword(password string) (string, error)
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userService *service.UserService, authService AuthService) *UserUseCase {
	return &UserUseCase{
		userService:   userService,
		authService:   authService,
		userPresenter: presenter.NewUserPresenter(),
	}
}

// SetEmailService sets the email service
func (uc *UserUseCase) SetEmailService(emailService *email.EmailService) {
	uc.emailService = emailService
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Hash password
	hashedPassword, err := uc.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	// Create user entity
	user, err := entity.NewUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}
	
	// Create user
	if err := uc.userService.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	
	// Send welcome email if email service is available
	if uc.emailService != nil {
		go func() {
			_ = uc.emailService.SendRegistrationEmail(user.Email, user.Name)
		}()
	}
	
	// Convert to DTO
	return uc.userPresenter.ToDTO(user), nil
}

// Login authenticates a user
func (uc *UserUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := uc.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	
	// Validate password
	if !uc.authService.ValidatePassword(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}
	
	// Generate token
	token, err := uc.authService.GenerateToken(user.UUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	userDTO := uc.userPresenter.ToDTO(user)
	
	return &dto.LoginResponse{
		User:  *userDTO,
		Token: token,
	}, nil
}

// GetUserByUUID gets a user by UUID
func (uc *UserUseCase) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*dto.UserResponse, error) {
	// Get user
	user, err := uc.userService.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.userPresenter.ToDTO(user), nil
}

// GetAllUsers gets all users
func (uc *UserUseCase) GetAllUsers(ctx context.Context) (*dto.UsersResponse, error) {
	// Get all users
	users, err := uc.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTOs
	return uc.userPresenter.ToDTOList(users), nil
}

// UpdateUser updates a user
func (uc *UserUseCase) UpdateUser(ctx context.Context, userUUID uuid.UUID, req *dto.UserResponse) (*dto.UserResponse, error) {
	// Get user
	user, err := uc.userService.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	
	// Update user fields
	if req.Name != "" {
		if err := user.UpdateName(req.Name); err != nil {
			return nil, err
		}
	}
	
	if req.Email != "" && req.Email != user.Email {
		if err := user.UpdateEmail(req.Email); err != nil {
			return nil, err
		}
	}
	
	// Update user
	if err := uc.userService.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.userPresenter.ToDTO(user), nil
}