package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	LogLevel    string
	Debug       bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()
	
	// Get environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Try legacy environment variable name
		databaseURL = os.Getenv("DB_URL")
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Try legacy environment variable name
		jwtSecret = os.Getenv("SECRET")
	}
	
	port := os.Getenv("PORT")
	environment := os.Getenv("ENVIRONMENT")
	
	// Email settings
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default for Gmail
	}
	
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // Default TLS port
	}
	
	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = os.Getenv("EMAIL_SENDER")
	}
	
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		smtpPass = os.Getenv("EMAIL_PASSWORD")
	}
	
	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = smtpUser // Default to the same as SMTP user
	}
	
	logLevel := os.Getenv("LOG_LEVEL")
	debugStr := os.Getenv("DEBUG")
	
	// Set defaults
	if port == "" {
		port = "8080"
	}
	
	if environment == "" {
		environment = "development"
	}
	
	if logLevel == "" {
		logLevel = "info"
	}
	
	// Parse debug flag
	debug := false
	if debugStr != "" {
		debug, _ = strconv.ParseBool(debugStr)
	}
	
	// Create config
	config := &Config{
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		Port:        port,
		Environment: environment,
		SMTPHost:    smtpHost,
		SMTPPort:    smtpPort,
		SMTPUser:    smtpUser,
		SMTPPass:    smtpPass,
		SMTPFrom:    smtpFrom,
		LogLevel:    logLevel,
		Debug:       debug,
	}
	
	return config, nil
}

// IsDevelopment checks if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTest checks if the environment is test
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}