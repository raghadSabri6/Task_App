package persistence

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          int64      `bun:",pk,autoincrement"`
	UUID        uuid.UUID  `bun:",type:uuid,default:uuid_generate_v4()" json:"id"`
	Title       string     `bun:",notnull" json:"title"`
	Description string     `json:"description"`
	Completed   bool       `bun:",default:false"`
	CreatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt   *time.Time `bun:",soft_delete" json:"deleted_at,omitempty"`

	CreatedByID uuid.UUID `bun:",type:uuid,notnull"`
	CreatedBy   *User     `bun:"rel:belongs-to,join:created_by_id=uuid"`

	AssignedToID *uuid.UUID `bun:",type:uuid"`
	AssignedTo   *User      `bun:"rel:belongs-to,join:assigned_to_id=uuid"`

	Users []*User `bun:"m2m:user_tasks" json:"users,omitempty"`
}
