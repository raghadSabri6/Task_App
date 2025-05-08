package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity
type User struct {
	ID        int64
	UUID      uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// References to other entities - initialized as empty slice to avoid nil issues
	Tasks []*Task
}

// NewUser creates a new user with the given parameters
func NewUser(name, email, hashedPassword string) (*User, error) {
	if name == "" || email == "" || hashedPassword == "" {
		return nil, errors.New("name, email, and password are required")
	}

	return &User{
		UUID:      uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tasks:     make([]*Task, 0), // Initialize with empty slice
	}, nil
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(hashedPassword string) error {
	if hashedPassword == "" {
		return errors.New("password cannot be empty")
	}
	
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}