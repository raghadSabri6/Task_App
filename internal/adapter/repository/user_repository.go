package repository

import (
	"context"
	"task2/internal/domain/entity"
	"task2/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// UserRepository implements the domain.UserRepository interface
type UserRepository struct {
	db *bun.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	// Convert domain entity to persistence model
	dbUser := &persistence.User{
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		// Tasks field is explicitly excluded with bun:"-" tag
	}
	
	// Insert user - explicitly specify columns to avoid tasks field
	_, err := r.db.NewInsert().
		Model(dbUser).
		Column("uuid", "name", "email", "password").
		Returning("id").
		Exec(ctx)
	
	if err != nil {
		return err
	}
	
	// Update user ID
	user.ID = dbUser.ID
	
	return nil
}

// GetByUUID gets a user by UUID
func (r *UserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	dbUser := new(persistence.User)
	
	// Get user
	err := r.db.NewSelect().
		Model(dbUser).
		Where("uuid = ?", uuid).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	
	// Convert to domain entity
	user := &entity.User{
		ID:        dbUser.ID,
		UUID:      dbUser.UUID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		DeletedAt: dbUser.DeletedAt,
	}
	
	return user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser := new(persistence.User)
	
	// Get user
	err := r.db.NewSelect().
		Model(dbUser).
		Where("email = ?", email).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	
	// Convert to domain entity
	user := &entity.User{
		ID:        dbUser.ID,
		UUID:      dbUser.UUID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		DeletedAt: dbUser.DeletedAt,
	}
	
	return user, nil
}

// GetAll gets all users
func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var dbUsers []persistence.User
	
	// Get all users
	err := r.db.NewSelect().
		Model(&dbUsers).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	
	// Convert to domain entities
	users := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = &entity.User{
			ID:        dbUser.ID,
			UUID:      dbUser.UUID,
			Name:      dbUser.Name,
			Email:     dbUser.Email,
			Password:  dbUser.Password,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			DeletedAt: dbUser.DeletedAt,
		}
	}
	
	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	// Convert domain entity to persistence model
	dbUser := &persistence.User{
		ID:       user.ID,
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UpdatedAt: user.UpdatedAt,
	}
	
	// Update user
	_, err := r.db.NewUpdate().
		Model(dbUser).
		Column("name", "email", "password", "updated_at").
		WherePK().
		Exec(ctx)
	
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	// Get user
	dbUser := new(persistence.User)
	err := r.db.NewSelect().
		Model(dbUser).
		Where("uuid = ?", uuid).
		Scan(ctx)
	
	if err != nil {
		return err
	}
	
	// Delete user
	_, err = r.db.NewDelete().
		Model(dbUser).
		WherePK().
		Exec(ctx)
	
	return err
}

// EmailExists checks if an email exists
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return r.db.NewSelect().
		Model((*persistence.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
}