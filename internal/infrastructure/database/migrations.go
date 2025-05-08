package database

import (
	"context"
	"fmt"
	"task2/internal/infrastructure/persistence"

	"github.com/uptrace/bun"
)

// RunMigrations runs database migrations
func RunMigrations(db *bun.DB) error {
	ctx := context.Background()
	
	// Create users table
	_, err := db.NewCreateTable().
		Model((*persistence.User)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	
	// Create tasks table
	_, err = db.NewCreateTable().
		Model((*persistence.Task)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}
	
	// Create user_tasks table
	_, err = db.NewCreateTable().
		Model((*persistence.UserTask)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user_tasks table: %w", err)
	}
	
	return nil
}

// AddIndexes adds database indexes
func AddIndexes(db *bun.DB) error {
	ctx := context.Background()
	
	// Add index on users.email
	_, err := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on users.email: %w", err)
	}
	
	// Add index on users.uuid
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_users_uuid ON users (uuid);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on users.uuid: %w", err)
	}
	
	// Add index on tasks.uuid
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_tasks_uuid ON tasks (uuid);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on tasks.uuid: %w", err)
	}
	
	// Add index on tasks.created_by_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_tasks_created_by_id ON tasks (created_by_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on tasks.created_by_id: %w", err)
	}
	
	// Add index on tasks.assigned_to_id
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to_id ON tasks (assigned_to_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on tasks.assigned_to_id: %w", err)
	}
	
	return nil
}