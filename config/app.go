package config

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
	"api-chatbot/internal/cache"
	"api-chatbot/internal/migration"
	"api-chatbot/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Env          *Env
	Db           *pgxpool.Pool
	Cache        domain.ParameterCache // Expose parameter cache for direct access
	Logger       *slog.Logger
	LoggerCloser func() error // Call on shutdown to close log files
}

func App() Application {
	app := &Application{}

	app.Env = NewEnv()

	app.Db = NewPostgresDatabase(app.Env)

	// Run database migrations
	if err := runMigrations(app.Env); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	app.Cache = cache.NewParameterCache()
	paramRepo := repository.NewParameterRepository(dal.NewDAL(app.Db))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params, err := paramRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Warning: Could not load parameters from database: %v", err)
	} else {
		app.Cache.LoadAll(params)
		log.Printf("Loaded %d parameters from database into cache", len(params))

		if param, exists := app.Cache.Get("APP_CONFIG"); exists {
			if data, err := param.GetDataAsMap(); err == nil {
				if appEnv, ok := data["appEnv"].(string); ok && appEnv == "development" {
					log.Println("La aplicación está corriendo en modo de desarrollo.")
				}
			}
		}
	}

	app.Logger, app.LoggerCloser = SetupLogger(app.Cache)

	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.Db)
}

// Shutdown gracefully closes all application resources
func (app *Application) Shutdown() {
	if app.LoggerCloser != nil {
		if err := app.LoggerCloser(); err != nil {
			slog.Error("Failed to close logger", "error", err)
		}
	}
	app.CloseDBConnection()
}

// runMigrations executes database migrations based on configuration
func runMigrations(env *Env) error {
	// Build DSN for migration runner
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.Database.User,
		env.Database.Password,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
	)

	migrationConfig := migration.Config{
		AutoMigrate: env.Migration.AutoMigrate,
		Verbose:     env.Migration.Verbose,
		DSN:         dsn,
	}

	return migration.RunMigrations(migrationConfig)
}
