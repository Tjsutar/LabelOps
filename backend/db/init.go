package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes DB connection, creates tables, procedures, and seeds data if configured
func InitDB() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "root")
	dbname := getEnv("DB_NAME", "labelops")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL database")

	if err = executeSQLFile("db/schema.sql"); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err = executeSQLFile("db/procedures.sql"); err != nil {
		return fmt.Errorf("failed to create stored procedures: %w", err)
	}

	// If FLUSH_DB env var is "true", flush DB
	if getEnv("FLUSH_DB", "false") == "true" {
		log.Println("⚠️  FLUSH_DB=true, truncating tables...")
		if err = executeSQLFile("db/flush.sql"); err != nil {
			return fmt.Errorf("failed to flush database: %w", err)
		}
	}

	// Seed initial data
	if getEnv("SEED_DB", "true") == "true" {
		log.Println("🌱 Seeding initial data...")
		if err = executeSQLFile("db/seed.sql"); err != nil {
			return fmt.Errorf("failed to seed database: %w", err)
		}
	}

	log.Println("🎉 Database initialized successfully")
	return nil
}

// executeSQLFile reads and executes SQL file content
func executeSQLFile(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filepath, err)
	}

	_, err = DB.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute SQL from file %s: %w", filepath, err)
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
