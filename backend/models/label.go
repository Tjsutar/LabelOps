package models

import (
	"time"
	"github.com/google/uuid"
)

// Label represents a label entry in the system
type Label struct {
	ID          uuid.UUID `json:"id" db:"id"`
	LabelID     string    `json:"label_id" db:"label_id" binding:"required"`
	ProductName string    `json:"product_name" db:"product_name"`
	SKU         string    `json:"sku" db:"sku"`
	Quantity    int       `json:"quantity" db:"quantity"`
	BatchNumber string    `json:"batch_number" db:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date" db:"expiry_date"`
	PrintedAt   *time.Time `json:"printed_at" db:"printed_at"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Status      string    `json:"status" db:"status"` // "pending", "printed", "failed"
	ZPLContent  string    `json:"zpl_content" db:"zpl_content"`
	QRCode      string    `json:"qr_code" db:"qr_code"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// LabelBatchRequest represents a batch of labels to be processed
type LabelBatchRequest struct {
	Labels []LabelData `json:"labels" binding:"required"`
}

// LabelData represents individual label data in a batch
type LabelData struct {
	LabelID     string     `json:"label_id" binding:"required"`
	ProductName string     `json:"product_name"`
	SKU         string     `json:"sku"`
	Quantity    int        `json:"quantity"`
	BatchNumber string     `json:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date"`
}

// LabelBatchResponse represents the response after processing a batch
type LabelBatchResponse struct {
	NewLabels     []Label `json:"new_labels"`
	DuplicateLabels []Label `json:"duplicate_labels"`
	TotalProcessed int     `json:"total_processed"`
	NewCount      int     `json:"new_count"`
	DuplicateCount int    `json:"duplicate_count"`
}

// PrintJob represents a print job in the system
type PrintJob struct {
	ID         uuid.UUID `json:"id" db:"id"`
	LabelID    uuid.UUID `json:"label_id" db:"label_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Status     string    `json:"status" db:"status"` // "pending", "printing", "completed", "failed"
	ErrorMsg   *string   `json:"error_msg" db:"error_msg"`
	Retries    int       `json:"retries" db:"retries"`
	MaxRetries int       `json:"max_retries" db:"max_retries"`
	ZPLContent string    `json:"zpl_content" db:"zpl_content"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// LabelFilter represents filters for querying labels
type LabelFilter struct {
	Status      *string    `json:"status"`
	UserID      *uuid.UUID `json:"user_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	ProductName *string    `json:"product_name"`
	SKU         *string    `json:"sku"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// LabelStats represents statistics about labels
type LabelStats struct {
	TotalLabels     int `json:"total_labels"`
	PrintedLabels   int `json:"printed_labels"`
	PendingLabels   int `json:"pending_labels"`
	FailedLabels    int `json:"failed_labels"`
	TotalPrintJobs  int `json:"total_print_jobs"`
	FailedPrintJobs int `json:"failed_print_jobs"`
} 