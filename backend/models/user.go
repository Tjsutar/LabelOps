package models

import (
	"time"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" binding:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash" binding:"required"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Role         string    `json:"role" db:"role"` // "admin", "user", "operator"
	IsActive     bool      `json:"is_active" db:"is_active"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresAt int64  `json:"expires_at"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role"`
}

// UserUpdateRequest represents a user update request
type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"email"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	Resource  string    `json:"resource" db:"resource"`
	ResourceID *string  `json:"resource_id" db:"resource_id"`
	Details   string    `json:"details" db:"details"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SystemStats represents system statistics
type SystemStats struct {
	TotalUsers      int `json:"total_users"`
	ActiveUsers     int `json:"active_users"`
	TotalLabels     int `json:"total_labels"`
	PrintedLabels   int `json:"printed_labels"`
	PendingLabels   int `json:"pending_labels"`
	FailedLabels    int `json:"failed_labels"`
	TotalPrintJobs  int `json:"total_print_jobs"`
	FailedPrintJobs int `json:"failed_print_jobs"`
} 