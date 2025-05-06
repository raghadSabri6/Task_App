package schemas

import "github.com/google/uuid"

type TaskIdRequest struct {
	TaskID uuid.UUID `json:"task_id"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description"`
	Users       []UserAssign `json:"users,omitempty"`
}

type UserAssign struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
