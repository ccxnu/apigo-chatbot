package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"api-chatbot/config"
	"api-chatbot/internal/migration"
)

func main() {
	// Define CLI flags
	upCmd := flag.Bool("up", false, "Run all pending migrations")
	downCmd := flag.Bool("down", false, "Rollback one migration")
	versionCmd := flag.Bool("version", false, "Show current migration version")
	forceCmd := flag.Int("force", -1, "Force set migration version (use with caution!)")
	stepsCmd := flag.Int("steps", 0, "Run N migrations (positive for up, negative for down)")

	flag.Parse()

	// Load environment configuration
	env := config.NewEnv()

	// Build DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.Database.User,
		env.Database.Password,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
	)

	// Create migration runner
	runner, err := migration.NewRunner(migration.Config{
		DSN:     dsn,
		Verbose: true,
	})
	if err != nil {
		log.Fatalf("Failed to create migration runner: %v", err)
	}
	defer runner.Close()

	// Execute commands
	switch {
	case *versionCmd:
		version, dirty, err := runner.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		if dirty {
			fmt.Printf("Current version: %d (DIRTY - needs manual intervention)\n", version)
			os.Exit(1)
		}
		fmt.Printf("Current version: %d\n", version)

	case *forceCmd >= 0:
		if err := runner.Force(*forceCmd); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		fmt.Printf("Forced migration version to: %d\n", *forceCmd)

	case *upCmd:
		if err := runner.Up(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		version, _, _ := runner.Version()
		fmt.Printf("Migrations completed successfully. Current version: %d\n", version)

	case *downCmd:
		if err := runner.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		version, _, _ := runner.Version()
		fmt.Printf("Rollback completed successfully. Current version: %d\n", version)

	case *stepsCmd != 0:
		if err := runner.Steps(*stepsCmd); err != nil {
			log.Fatalf("Migration steps failed: %v", err)
		}
		version, _, _ := runner.Version()
		fmt.Printf("Migrations completed successfully. Current version: %d\n", version)

	default:
		flag.Usage()
		fmt.Println("\nExamples:")
		fmt.Println("  migrate -up              Run all pending migrations")
		fmt.Println("  migrate -down            Rollback one migration")
		fmt.Println("  migrate -version         Show current version")
		fmt.Println("  migrate -force 7         Force version to 7 (fixes dirty state)")
		fmt.Println("  migrate -steps 2         Run next 2 migrations")
		fmt.Println("  migrate -steps -1        Rollback 1 migration")
	}
}
