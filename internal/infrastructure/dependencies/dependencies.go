package dependencies

import (
	"fmt"
	"strconv"
	"task2/internal/infrastructure/config"
	"task2/internal/infrastructure/database"
	"task2/pkg/email"

	"github.com/uptrace/bun"
)

// Dependencies holds all external service dependencies
type Dependencies struct {
	DB          *bun.DB
	EmailClient *email.EmailService
	// Add other external dependencies here as needed:
	// RedisClient  *redis.Client
	// CacheClient  cache.Client
	// MessageQueue mq.Client
	// ExternalAPI  api.Client
}

// NewDependencies initializes all external dependencies
func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize email client
	emailClient, err := initEmailClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email client: %w", err)
	}

	// Create dependencies container
	deps := &Dependencies{
		DB:          db,
		EmailClient: emailClient,
		// Initialize other dependencies here
	}

	return deps, nil
}

// initDatabase initializes the database connection
func initDatabase() (*bun.DB, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}

	// Register models
	database.RegisterModels(db)

	// Sync database schema
	if err := database.SyncDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}

// initEmailClient initializes the email client
func initEmailClient(cfg *config.Config) (*email.EmailService, error) {
	// Parse SMTP port
	smtpPort, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Create email service
	emailClient := email.NewEmailService(
		cfg.SMTPHost,
		smtpPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.SMTPFrom,
	)

	return emailClient, nil
}

// Close closes all connections
func (d *Dependencies) Close() error {
	var errs []error

	// Close database connection
	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing database: %w", err))
		}
	}

	// Add code to close other connections here
	// For example:
	// if d.RedisClient != nil {
	//     if err := d.RedisClient.Close(); err != nil {
	//         errs = append(errs, fmt.Errorf("error closing Redis: %w", err))
	//     }
	// }

	// Return combined errors if any
	if len(errs) > 0 {
		return fmt.Errorf("errors closing dependencies: %v", errs)
	}

	return nil
}