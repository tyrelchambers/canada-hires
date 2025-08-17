package db

import (
	"canada-hires/utils"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database interface defines the contract for database operations
type Database interface {
	GetDB() *sqlx.DB
	Close() error
	Ping() error
}

// Global database instance
var instance Database

// PostgresDB implements the Database interface for PostgreSQL
type PostgresDB struct {
	db *sqlx.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config Config) (Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresDB{db: db}, nil
}

// GetDB returns the underlying sqlx.DB instance
func (p *PostgresDB) GetDB() *sqlx.DB {
	return p.db
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// Ping checks if the database connection is alive
func (p *PostgresDB) Ping() error {
	return p.db.Ping()
}

// NewConfigFromEnv loads database configuration from environment variables
func NewConfigFromEnv() Config {
	// Load .env file if it exists

	return Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "canada-hires"),
		SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}
}

// InitDB initializes the database connection
func InitDB() Database {
	config := NewConfigFromEnv()

	db, err := NewPostgresDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Debug("Successfully connected to database")
	// Set the global instance
	instance = db

	// Start connection pool monitoring goroutine

	return db
}

// GetInstance returns the global database instance
// This can be called from any package to get access to the database
func GetInstance() Database {
	if instance == nil {
		log.Fatal("Database not initialized. Call InitDB() first")
	}
	return instance
}

// RunMigrations executes database migrations using the migrate command
func RunMigrations(db *sql.DB) error {
	log.Debug("Running migrations")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("Failed to create postgres driver", "error", err)
		return err
	}

	// Check if we're running in a Docker container (migrations at /app/migrations)
	// or locally (migrations at ./migrations)
	migrationPath := "./migrations"
	if _, err := os.Stat("/app/migrations"); err == nil {
		migrationPath = "/app/migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres", driver)
	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return err
	}
	m.Up() // or m.Steps(2) if you want to explicitly set the number of migrations to run

	log.Debug("Migrations completed successfully")
	return nil
}
