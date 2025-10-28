package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Config holds migration configuration
type Config struct {
	AutoMigrate bool   // Auto-run migrations on startup
	Verbose     bool   // Enable verbose logging
	DSN         string // Database connection string
}

// Runner manages database migrations
type Runner struct {
	migrate *migrate.Migrate
	config  Config
}

// NewRunner creates a new migration runner
func NewRunner(config Config) (*Runner, error) {
	// Open database connection using pgx/stdlib for compatibility with migrate
	db, err := sql.Open("pgx", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "chatbot_db",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create source driver from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Enable verbose logging if configured
	if config.Verbose {
		m.Log = &logger{}
	}

	return &Runner{
		migrate: m,
		config:  config,
	}, nil
}

// Up runs all pending migrations
func (r *Runner) Up() error {
	if err := r.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	version, dirty, err := r.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if r.config.Verbose {
		if dirty {
			log.Printf("[MIGRATION] Current version: %d (dirty)", version)
		} else {
			log.Printf("[MIGRATION] Current version: %d", version)
		}
	}

	return nil
}

// Down rolls back one migration
func (r *Runner) Down() error {
	if err := r.migrate.Down(); err != nil {
		return fmt.Errorf("migration down failed: %w", err)
	}
	return nil
}

// Steps runs n migrations (positive for up, negative for down)
func (r *Runner) Steps(n int) error {
	if err := r.migrate.Steps(n); err != nil {
		return fmt.Errorf("migration steps failed: %w", err)
	}
	return nil
}

// Version returns the current migration version
func (r *Runner) Version() (uint, bool, error) {
	return r.migrate.Version()
}

// Force sets the migration version without running migrations
func (r *Runner) Force(version int) error {
	if err := r.migrate.Force(version); err != nil {
		return fmt.Errorf("migration force failed: %w", err)
	}
	return nil
}

// Close closes the migration instance
func (r *Runner) Close() error {
	sourceErr, dbErr := r.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close database: %w", dbErr)
	}
	return nil
}

// RunMigrations is a helper function to run migrations based on config
func RunMigrations(config Config) error {
	runner, err := NewRunner(config)
	if err != nil {
		return fmt.Errorf("failed to create migration runner: %w", err)
	}
	defer runner.Close()

	if config.AutoMigrate {
		log.Println("[MIGRATION] Auto-migration enabled, running pending migrations...")
		if err := runner.Up(); err != nil {
			return fmt.Errorf("auto-migration failed: %w", err)
		}
		log.Println("[MIGRATION] Migrations completed successfully")
	} else {
		// Just check the version
		version, dirty, err := runner.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		if err == migrate.ErrNilVersion {
			log.Println("[MIGRATION] No migrations applied yet. Set AUTO_MIGRATE=true to run migrations.")
		} else if dirty {
			log.Printf("[MIGRATION] WARNING: Database is in dirty state at version %d. Manual intervention required.", version)
		} else {
			log.Printf("[MIGRATION] Current migration version: %d", version)
		}
	}

	return nil
}

// logger implements migrate.Logger interface
type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	log.Printf("[MIGRATION] "+format, v...)
}

func (l *logger) Verbose() bool {
	return true
}
