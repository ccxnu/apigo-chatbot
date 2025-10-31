package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

// Simple script to create an admin user
// Usage: go run scripts/create_admin.go -username=admin -email=admin@example.com -password=yourpassword -name="Admin User"

func main() {
	// Parse command line flags
	username := flag.String("username", "", "Admin username (required)")
	email := flag.String("email", "", "Admin email (required)")
	password := flag.String("password", "", "Admin password (required)")
	name := flag.String("name", "", "Admin full name (required)")
	role := flag.String("role", "ROLE_ADMIN", "Admin role (default: ROLE_ADMIN)")

	flag.Parse()

	// Validate required fields
	if *username == "" || *email == "" || *password == "" || *name == "" {
		log.Fatal("Error: All fields are required (username, email, password, name)")
	}

	// Connect to database
	connString := "host=localhost port=5432 user=postgres password=lo0G4Rfaw7gtHw0wvpm4aqi4 dbname=chatbot_db sslmode=disable"
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database")

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	fmt.Println("✓ Password hashed successfully")

	// Call stored procedure to create admin user
	var success bool
	var code string
	var adminID int

	ctx := context.Background()
	row := db.QueryRowContext(ctx,
		"CALL sp_create_admin_user($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		nil, // OUT success
		nil, // OUT code
		nil, // OUT admin_id
		*username,
		*email,
		string(hashedPassword),
		*name,
		*role,
		[]byte("[]"), // permissions (empty array)
		[]byte("{}"), // claims (empty object)
	)

	// For PostgreSQL procedures, we need to use a different approach
	// Let's use a direct query instead
	err = db.QueryRowContext(ctx, `
		DO $$
		DECLARE
			v_success BOOLEAN;
			v_code VARCHAR;
			v_admin_id INT;
		BEGIN
			CALL sp_create_admin_user(
				v_success, v_code, v_admin_id,
				$1, $2, $3, $4, $5, $6::jsonb, $7::jsonb
			);

			IF NOT v_success THEN
				RAISE EXCEPTION 'Failed to create admin: %', v_code;
			END IF;

			RAISE NOTICE 'Admin created with ID: %', v_admin_id;
		END $$;
	`, *username, *email, string(hashedPassword), *name, *role, "[]", "{}").Scan()

	if err != nil && err != sql.ErrNoRows {
		// Check if it's a specific error
		if contains(err.Error(), "ERR_USERNAME_EXISTS") {
			log.Fatalf("✗ Error: Username '%s' already exists", *username)
		} else if contains(err.Error(), "ERR_EMAIL_EXISTS") {
			log.Fatalf("✗ Error: Email '%s' already exists", *email)
		} else {
			log.Fatalf("✗ Failed to create admin user: %v", err)
		}
	}

	fmt.Printf("\n✓ Admin user created successfully!\n\n")
	fmt.Printf("Username: %s\n", *username)
	fmt.Printf("Email:    %s\n", *email)
	fmt.Printf("Name:     %s\n", *name)
	fmt.Printf("Role:     %s\n", *role)
	fmt.Printf("\nYou can now login at: http://localhost:8080/admin/login\n")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
