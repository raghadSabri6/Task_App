package database

import (
	"database/sql"
	"fmt"
	"os"
	"task2/internal/infrastructure/persistence"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// DB is the database connection
var DB *bun.DB

// ConnectToDB connects to the database
func ConnectToDB() (*bun.DB, error) {
	// Get database connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Try legacy environment variable name
		dsn = os.Getenv("DB_URL")
		if dsn == "" {
			return nil, fmt.Errorf("DATABASE_URL or DB_URL environment variable not set")
		}
	}
	
	// Connect to database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	
	// Set connection pool settings
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)
	
	// Create bun DB
	db := bun.NewDB(sqldb, pgdialect.New())
	
	// Set global DB
	DB = db
	
	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

// RegisterModels registers database models
func RegisterModels(db *bun.DB) {
	// Register models in the correct order
	// Register the join table (UserTask) first before the models that use it in m2m relationships
	db.RegisterModel((*persistence.UserTask)(nil))
	db.RegisterModel((*persistence.User)(nil))
	db.RegisterModel((*persistence.Task)(nil))
}

// SyncDatabase synchronizes the database schema
func SyncDatabase(db *bun.DB) error {
	// Run migrations
	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	// Add indexes
	if err := AddIndexes(db); err != nil {
		return fmt.Errorf("failed to add indexes: %w", err)
	}
	
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}