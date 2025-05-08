package presenter

import (
	"task2/internal/app/dto"
	"task2/internal/domain/entity"
)

// UserPresenter converts between domain entities and DTOs
type UserPresenter struct{}

// NewUserPresenter creates a new user presenter
func NewUserPresenter() *UserPresenter {
	return &UserPresenter{}
}

// ToDTO converts a user entity to a DTO
func (p *UserPresenter) ToDTO(user *entity.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	
	return &dto.UserResponse{
		ID:        user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

// ToSummary converts a user entity to a summary DTO
func (p *UserPresenter) ToSummary(user *entity.User) *dto.UserSummary {
	if user == nil {
		return nil
	}
	
	return &dto.UserSummary{
		ID:    user.UUID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// ToDTOList converts a list of user entities to DTOs
func (p *UserPresenter) ToDTOList(users []*entity.User) *dto.UsersResponse {
	if users == nil {
		return &dto.UsersResponse{
			Users: []dto.UserResponse{},
		}
	}
	
	// Create user responses
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponse := p.ToDTO(user)
		if userResponse != nil {
			userResponses[i] = *userResponse
		}
	}
	
	return &dto.UsersResponse{
		Users: userResponses,
	}
}