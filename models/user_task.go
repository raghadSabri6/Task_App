package models

import (
	"github.com/uptrace/bun"
)

type UserTask struct {
	bun.BaseModel `bun:"table:user_tasks,alias:ut"`

	TaskID int64 `bun:",pk"`
	UserID int64 `bun:",pk"`

	Task *Task `bun:"rel:belongs-to,join:task_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
