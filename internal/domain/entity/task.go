package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Task represents the core task entity
type Task struct {
	ID          int64
	UUID        uuid.UUID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	CreatedByID uuid.UUID
	AssignedToID *uuid.UUID

	// References to other entities
	CreatedBy  *User
	AssignedTo *User
	Users      []*User
}

// NewTask creates a new task with the given parameters
func NewTask(title, description string, createdByID uuid.UUID) (*Task, error) {
	if title == "" {
		return nil, errors.New("task title is required")
	}

	return &Task{
		UUID:        uuid.New(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedByID: createdByID,
	}, nil
}

// Complete marks a task as completed
func (t *Task) Complete() error {
	if t.Completed {
		return errors.New("task is already completed")
	}
	
	t.Completed = true
	t.UpdatedAt = time.Now()
	return nil
}

// CanBeModifiedBy checks if a user can modify this task
func (t *Task) CanBeModifiedBy(userID uuid.UUID) bool {
	// Task creator can always modify
	if t.CreatedByID == userID {
		return true
	}
	
	// Check if user is assigned to this task
	for _, user := range t.Users {
		if user.UUID == userID {
			return true
		}
	}
	
	return false
}

// AssignTo assigns the task to a user
func (t *Task) AssignTo(userID uuid.UUID) {
	t.AssignedToID = &userID
	t.UpdatedAt = time.Now()
}

// AddUser adds a user to the task
func (t *Task) AddUser(user *User) {
	// Check if user is already assigned
	for _, u := range t.Users {
		if u.UUID == user.UUID {
			return // User already assigned
		}
	}
	
	t.Users = append(t.Users, user)
	t.UpdatedAt = time.Now()
}