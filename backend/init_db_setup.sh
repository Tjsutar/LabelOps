#!/bin/bash

set -e

echo "Creating folder structure..."
# mkdir -p db

echo "Creating db/init.go..."
cat > db/init.go <<'EOF'
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
EOF

echo "Creating db/schema.sql..."
cat > db/schema.sql <<'EOF'
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	role VARCHAR(50) NOT NULL DEFAULT 'user',
	is_active BOOLEAN NOT NULL DEFAULT true,
	last_login TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS labels (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	label_id VARCHAR(255) UNIQUE NOT NULL,
	location VARCHAR(100),
	bundle_no VARCHAR(255) NOT NULL,
	pqd VARCHAR(255) NOT NULL,
	unit VARCHAR(50) NOT NULL,
	time VARCHAR(10) NOT NULL,
	length INTEGER,
	heat_no VARCHAR(100) NOT NULL,
	product_heading VARCHAR(255) NOT NULL,
	isi_bottom VARCHAR(255) NOT NULL,
	isi_top VARCHAR(255) NOT NULL,
	charge_dtm VARCHAR(255),
	mill VARCHAR(50) NOT NULL,
	grade VARCHAR(100) NOT NULL,
	url_apikey VARCHAR(255) NOT NULL,
	weight VARCHAR(50),
	section VARCHAR(255) NOT NULL,
	date VARCHAR(20) NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id),
	status VARCHAR(50) NOT NULL DEFAULT 'completed',
	is_duplicate BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS print_jobs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	label_id UUID NOT NULL REFERENCES labels(id),
	user_id UUID NOT NULL REFERENCES users(id),
	status VARCHAR(50) NOT NULL DEFAULT 'completed',
	error_message TEXT,
	retry_count INTEGER NOT NULL DEFAULT 0,
	max_retries INTEGER NOT NULL DEFAULT 3,
	zpl_content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id),
	action VARCHAR(100) NOT NULL,
	resource VARCHAR(100) NOT NULL,
	resource_id VARCHAR(255),
	details TEXT,
	ip_address VARCHAR(45),
	user_agent TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_labels_label_id ON labels(label_id);
CREATE INDEX IF NOT EXISTS idx_labels_user_id ON labels(user_id);
CREATE INDEX IF NOT EXISTS idx_labels_status ON labels(status);
CREATE INDEX IF NOT EXISTS idx_labels_created_at ON labels(created_at);
CREATE INDEX IF NOT EXISTS idx_print_jobs_status ON print_jobs(status);
CREATE INDEX IF NOT EXISTS idx_print_jobs_user_id ON print_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
EOF

echo "Creating db/procedures.sql..."
cat > db/procedures.sql <<'EOF'
CREATE OR REPLACE FUNCTION batch_label_process(
	labels_json JSONB,
	user_uuid UUID
) RETURNS JSONB AS $$
DECLARE
	label_record JSONB;
	label_id_val VARCHAR(255);
	existing_label_id VARCHAR(255);
	new_labels JSONB := '[]'::JSONB;
	duplicate_labels JSONB := '[]'::JSONB;
	new_count INTEGER := 0;
	duplicate_count INTEGER := 0;
BEGIN
	FOR label_record IN SELECT * FROM jsonb_array_elements(labels_json)
	LOOP
		label_id_val := label_record->>'PQD';

		SELECT label_id INTO existing_label_id 
		FROM labels 
		WHERE label_id = label_id_val;

		IF existing_label_id IS NULL THEN
			INSERT INTO labels (
				label_id, location, bundle_no, pqd, unit, time, length,
				heat_no, product_heading, isi_bottom, isi_top, charge_dtm,
				mill, grade, url_apikey, weight, section, date, user_id,
				status, is_duplicate
			) VALUES (
				label_id_val,
				label_record->>'LOCATION',
				(label_record->>'BUNDLE_NO')::INTEGER,
				label_record->>'PQD',
				label_record->>'UNIT',
				label_record->>'TIME',
				(label_record->>'LENGTH')::INTEGER,
				label_record->>'HEAT_NO',
				label_record->>'PRODUCT_HEADING',
				label_record->>'ISI_BOTTOM',
				label_record->>'ISI_TOP',
				label_record->>'CHARGE_DTM',
				label_record->>'MILL',
				label_record->>'GRADE',
				label_record->>'URL_APIKEY',
				label_record->>'WEIGHT',
				label_record->>'SECTION',
				label_record->>'DATE',
				user_uuid,
				'pending',
				false
			);

			new_labels := new_labels || label_record;
			new_count := new_count + 1;
		ELSE
			duplicate_labels := duplicate_labels || label_record;
			duplicate_count := duplicate_count + 1;
		END IF;
	END LOOP;

	RETURN jsonb_build_object(
		'new_labels', new_labels,
		'duplicate_labels', duplicate_labels,
		'total_processed', jsonb_array_length(labels_json),
		'new_count', new_count,
		'duplicate_count', duplicate_count
	);
END;
$$ LANGUAGE plpgsql;
EOF

echo "Creating db/seed.sql..."
cat > db/seed.sql <<'EOF'
INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
VALUES (
	'admin@gmail.com',
	crypt('Admin@123', gen_salt('bf')),
	'admin',
	'admin',
	'admin',
	true
)
ON CONFLICT (email) DO NOTHING;
EOF

echo "Creating db/flush.sql..."
cat > db/flush.sql <<'EOF'
-- Truncate tables with cascade for FK relations
TRUNCATE audit_logs, print_jobs, labels, users RESTART IDENTITY CASCADE;
EOF

echo "Creating .env.example..."
cat > .env.example <<'EOF'
# Database Config
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=labelops

# JWT Config
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Config
PORT=8080
GIN_MODE=debug

# CORS Config
CORS_ORIGIN=http://localhost:4200

# Logging Config
LOG_LEVEL=info

# Database Options
FLUSH_DB=false
SEED_DB=true
EOF

echo "All files created successfully."
