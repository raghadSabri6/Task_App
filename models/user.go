package models

import (
	"context"
	"errors"
	"fmt"
	"task2/database"
	"task2/helperFunc"
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

func (u *User) CreateUser(db bun.IDB, name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	hashedPassword, err := helperFunc.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := &User{
		UUID:     uuid.New(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	ctx := context.Background()

	existing := new(User)
	err = db.NewSelect().Model(existing).Where("email = ?", email).Scan(ctx)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	_, err = db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	user := new(User)

	err := database.DB.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {
	ctx := context.Background()
	var users []User

	err := database.DB.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
