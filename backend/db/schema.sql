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
	status VARCHAR(50) NOT NULL DEFAULT 'success',
	is_duplicate BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS print_jobs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	label_id UUID NOT NULL REFERENCES labels(id),
	user_id UUID NOT NULL REFERENCES users(id),
	status VARCHAR(50) NOT NULL DEFAULT 'success',
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
