package models

import (
	"time"

	"github.com/google/uuid"
)

// LabelData represents the label data structure from the API
type LabelData struct {
	LOCATION        *string `json:"LOCATION"`
	BUNDLE_NOS      int     `json:"BUNDLE_NOS"`
	PQD             string  `json:"PQD"`
	UNIT            string  `json:"UNIT"`
	TIME1           string  `json:"TIME1"`
	LENGTH          string  `json:"LENGTH"`
	HEAT_NO         string  `json:"HEAT_NO"`
	PRODUCT_HEADING string  `json:"PRODUCT_HEADING"`
	ISI_BOTTOM      string  `json:"ISI_BOTTOM"`
	ISI_TOP         string  `json:"ISI_TOP"`
	CHARGE_DTM      string  `json:"CHARGE_DTM"`
	MILL            string  `json:"MILL"`
	GRADE           string  `json:"GRADE"`
	URL_APIKEY      string  `json:"URL_APIKEY"`
	ID              *string `json:"ID"`
	WEIGHT          *string `json:"WEIGHT"`
	SECTION         string  `json:"SECTION"`
	DATE1           string  `json:"DATE1"`
}

// Label represents a label in the system
type Label struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	LabelID        string     `json:"label_id" db:"label_id" binding:"required"`
	Location       *string    `json:"location" db:"location"`
	BundleNos      int        `json:"bundle_nos" db:"bundle_nos"`
	PQD            string     `json:"pqd" db:"pqd"`
	Unit           string     `json:"unit" db:"unit"`
	Time1          string     `json:"time1" db:"time1"`
	Length         string     `json:"length" db:"length"`
	HeatNo         string     `json:"heat_no" db:"heat_no"`
	ProductHeading string     `json:"product_heading" db:"product_heading"`
	IsiBottom      string     `json:"isi_bottom" db:"isi_bottom"`
	IsiTop         string     `json:"isi_top" db:"isi_top"`
	ChargeDtm      string     `json:"charge_dtm" db:"charge_dtm"`
	Mill           string     `json:"mill" db:"mill"`
	Grade          string     `json:"grade" db:"grade"`
	UrlApikey      string     `json:"url_apikey" db:"url_apikey"`
	Weight         *string    `json:"weight" db:"weight"`
	Section        string     `json:"section" db:"section"`
	Date1          string     `json:"date1" db:"date1"`
	PrintedAt      *time.Time `json:"printed_at" db:"printed_at"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	Status         string     `json:"status" db:"status"` // "pending", "printed", "failed"
	ZPLContent     string     `json:"zpl_content" db:"zpl_content"`
	QRCode         string     `json:"qr_code" db:"qr_code"`
	IsDuplicate    bool       `json:"is_duplicate" db:"is_duplicate"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// LabelBatchRequest represents a batch of labels to be processed
type LabelBatchRequest struct {
	Labels []LabelData `json:"labels" binding:"required"`
}

// LabelBatchResponse represents the response after processing a label batch
type LabelBatchResponse struct {
	NewLabels       []Label `json:"new_labels"`
	DuplicateLabels []Label `json:"duplicate_labels"`
	TotalProcessed  int     `json:"total_processed"`
	NewCount        int     `json:"new_count"`
	DuplicateCount  int     `json:"duplicate_count"`
}

// LabelFilter represents filters for querying labels
type LabelFilter struct {
	Status      *string    `json:"status"`
	UserID      *uuid.UUID `json:"user_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Grade       *string    `json:"grade"`
	Section     *string    `json:"section"`
	HeatNo      *string    `json:"heat_no"`
	IsDuplicate *bool      `json:"is_duplicate"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// LabelStats represents statistics about labels
type LabelStats struct {
	TotalLabels     int            `json:"total_labels"`
	PrintedLabels   int            `json:"printed_labels"`
	PendingLabels   int            `json:"pending_labels"`
	FailedLabels    int            `json:"failed_labels"`
	DuplicateLabels int            `json:"duplicate_labels"`
	TotalPrintJobs  int            `json:"total_print_jobs"`
	FailedPrintJobs int            `json:"failed_print_jobs"`
	ByGrade         map[string]int `json:"by_grade"`
	BySection       map[string]int `json:"by_section"`
}
