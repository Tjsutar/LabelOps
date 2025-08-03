package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "root")
	dbname := getEnv("DB_NAME", "labelops")

	// Create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Create stored procedures
	if err = createStoredProcedures(); err != nil {
		return fmt.Errorf("failed to create stored procedures: %v", err)
	}

	return nil
}

// createTables creates all necessary tables
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
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
		)`,
		`CREATE TABLE IF NOT EXISTS labels (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			label_id VARCHAR(255) UNIQUE NOT NULL,
			location VARCHAR(100),
			bundle_nos INTEGER NOT NULL,
			pqd VARCHAR(255) NOT NULL,
			unit VARCHAR(50) NOT NULL,
			time1 VARCHAR(10) NOT NULL,
			length VARCHAR(50) NOT NULL,
			heat_no VARCHAR(100) NOT NULL,
			product_heading VARCHAR(255) NOT NULL,
			isi_bottom VARCHAR(255) NOT NULL,
			isi_top VARCHAR(255) NOT NULL,
			charge_dtm VARCHAR(255) NOT NULL,
			mill VARCHAR(50) NOT NULL,
			grade VARCHAR(100) NOT NULL,
			url_apikey VARCHAR(255) NOT NULL,
			weight VARCHAR(50),
			section VARCHAR(255) NOT NULL,
			date1 VARCHAR(20) NOT NULL,
			user_id UUID NOT NULL REFERENCES users(id),
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			is_duplicate BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS print_jobs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			label_id UUID NOT NULL REFERENCES labels(id),
			user_id UUID NOT NULL REFERENCES users(id),
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			error_message TEXT,
			retry_count INTEGER NOT NULL DEFAULT 0,
			max_retries INTEGER NOT NULL DEFAULT 3,
			zpl_content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id),
			action VARCHAR(100) NOT NULL,
			resource VARCHAR(100) NOT NULL,
			resource_id VARCHAR(255),
			details TEXT,
			ip_address VARCHAR(45),
			user_agent TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_labels_label_id ON labels(label_id)`,
		`CREATE INDEX IF NOT EXISTS idx_labels_user_id ON labels(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_labels_status ON labels(status)`,
		`CREATE INDEX IF NOT EXISTS idx_labels_created_at ON labels(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_print_jobs_status ON print_jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_print_jobs_user_id ON print_jobs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// createStoredProcedures creates the stored procedures for batch processing
func createStoredProcedures() error {
	queries := []string{
		`CREATE OR REPLACE FUNCTION batch_label_process(
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
			-- Loop through each label in the batch
			FOR label_record IN SELECT * FROM jsonb_array_elements(labels_json)
			LOOP
				label_id_val := label_record->>'PQD';
				
				-- Check if label already exists
				SELECT label_id INTO existing_label_id 
				FROM labels 
				WHERE label_id = label_id_val;
				
				IF existing_label_id IS NULL THEN
					-- Insert new label
					INSERT INTO labels (
						label_id, location, bundle_nos, pqd, unit, time1, length,
						heat_no, product_heading, isi_bottom, isi_top, charge_dtm,
						mill, grade, url_apikey, weight, section, date1, user_id,
						status, is_duplicate
					) VALUES (
						label_id_val,
						label_record->>'LOCATION',
						(label_record->>'BUNDLE_NOS')::INTEGER,
						label_record->>'PQD',
						label_record->>'UNIT',
						label_record->>'TIME1',
						label_record->>'LENGTH',
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
						label_record->>'DATE1',
						user_uuid,
						'pending',
						false
					);
					
					new_labels := new_labels || label_record;
					new_count := new_count + 1;
				ELSE
					-- Add to duplicates
					duplicate_labels := duplicate_labels || label_record;
					duplicate_count := duplicate_count + 1;
				END IF;
			END LOOP;
			
			-- Return result
			RETURN jsonb_build_object(
				'new_labels', new_labels,
				'duplicate_labels', duplicate_labels,
				'total_processed', jsonb_array_length(labels_json),
				'new_count', new_count,
				'duplicate_count', duplicate_count
			);
		END;
		$$ LANGUAGE plpgsql;`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to create stored procedure: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Stored procedures created successfully")
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
