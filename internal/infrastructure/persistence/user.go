package persistence

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64      `bun:",pk,autoincrement"`
	UUID      uuid.UUID  `bun:",type:uuid,default:uuid_generate_v4()" json:"id"`
	Name      string     `bun:",notnull" json:"name" validate:"required"`
	Email     string     `bun:",unique,notnull" json:"email" validate:"required,email"`
	Password  string     `bun:",notnull" json:"password,omitempty" validate:"required,min=6"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:",soft_delete" json:"deleted_at,omitempty"`

	Tasks []*Task `bun:"m2m:user_tasks" json:"tasks,omitempty"`
}
