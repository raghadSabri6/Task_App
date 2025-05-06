package schemas

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AssignTaskRequest struct {
	TaskID     uuid.UUID `json:"task_id"`
	AssignedTo uuid.UUID `json:"assigned_to"`
}
