package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"labelops-backend/db"
	"labelops-backend/models"
	"labelops-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProcessTMTBarBatch processes a batch of TMT Bar labels
func ProcessTMTBarBatch(c *gin.Context) {
	var req models.TMTBarBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Convert labels to JSON for stored procedure
	labelsJSON, err := json.Marshal(req.Labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal labels"})
		return
	}

	// Call stored procedure
	var resultJSON []byte
	err = db.DB.QueryRow("SELECT process_tmt_bar_batch($1, $2)", labelsJSON, userModel.ID).Scan(&resultJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process batch"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse result"})
		return
	}

	// Log audit
	utils.LogAudit(c, userModel.ID, "process_batch", "tmt_bar_labels", nil, 
		"Processed batch of labels", map[string]interface{}{
			"total_processed": result["total_processed"],
			"new_count":       result["new_count"],
			"duplicate_count": result["duplicate_count"],
		})

	c.JSON(http.StatusOK, result)
}

// GetTMTBarLabels retrieves TMT Bar labels with filtering
func GetTMTBarLabels(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	grade := c.Query("grade")
	section := c.Query("section")

	// Build query
	query := `SELECT id, label_id, location, bundle_nos, pqd, unit, time1, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date1, printed_at, user_id, status, 
			  zpl_content, qr_code, is_duplicate, created_at, updated_at 
			  FROM tmt_bar_labels WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add filters
	if status != "" {
		query += " AND status = $" + strconv.Itoa(argCount)
		args = append(args, status)
		argCount++
	}
	if grade != "" {
		query += " AND grade = $" + strconv.Itoa(argCount)
		args = append(args, grade)
		argCount++
	}
	if section != "" {
		query += " AND section = $" + strconv.Itoa(argCount)
		args = append(args, section)
		argCount++
	}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND user_id = $" + strconv.Itoa(argCount)
		args = append(args, userModel.ID)
		argCount++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argCount) + " OFFSET $" + strconv.Itoa(argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch labels"})
		return
	}
	defer rows.Close()

	var labels []models.TMTBarLabel
	for rows.Next() {
		var label models.TMTBarLabel
		err := rows.Scan(
			&label.ID, &label.LabelID, &label.Location, &label.BundleNos, &label.PQD,
			&label.Unit, &label.Time1, &label.Length, &label.HeatNo, &label.ProductHeading,
			&label.IsiBottom, &label.IsiTop, &label.ChargeDtm, &label.Mill, &label.Grade,
			&label.UrlApikey, &label.Weight, &label.Section, &label.Date1, &label.PrintedAt,
			&label.UserID, &label.Status, &label.ZPLContent, &label.QRCode, &label.IsDuplicate,
			&label.CreatedAt, &label.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan label"})
			return
		}
		labels = append(labels, label)
	}

	c.JSON(http.StatusOK, gin.H{"labels": labels, "count": len(labels)})
}

// PrintTMTBarLabel prints a specific TMT Bar label
func PrintTMTBarLabel(c *gin.Context) {
	labelID := c.Param("id")
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Parse label ID
	labelUUID, err := uuid.Parse(labelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	// Get label
	var label models.TMTBarLabel
	err = db.DB.QueryRow(
		`SELECT id, label_id, location, bundle_nos, pqd, unit, time1, length, 
		 heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
		 url_apikey, weight, section, date1, printed_at, user_id, status, 
		 zpl_content, qr_code, is_duplicate, created_at, updated_at 
		 FROM tmt_bar_labels WHERE id = $1`,
		labelUUID,
	).Scan(
		&label.ID, &label.LabelID, &label.Location, &label.BundleNos, &label.PQD,
		&label.Unit, &label.Time1, &label.Length, &label.HeatNo, &label.ProductHeading,
		&label.IsiBottom, &label.IsiTop, &label.ChargeDtm, &label.Mill, &label.Grade,
		&label.UrlApikey, &label.Weight, &label.Section, &label.Date1, &label.PrintedAt,
		&label.UserID, &label.Status, &label.ZPLContent, &label.QRCode, &label.IsDuplicate,
		&label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch label"})
		return
	}

	// Generate ZPL content if not exists
	if label.ZPLContent == "" {
		label.ZPLContent = utils.GenerateTMTBarZPL(label)
		// Update label with ZPL content
		db.DB.Exec("UPDATE tmt_bar_labels SET zpl_content = $1 WHERE id = $2", label.ZPLContent, label.ID)
	}

	// Create print job
	printJobID := uuid.New()
	_, err = db.DB.Exec(
		`INSERT INTO print_jobs (id, label_id, user_id, status, zpl_content, max_retries) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		printJobID, label.ID, userModel.ID, "pending", label.ZPLContent, 3,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create print job"})
		return
	}

	// Log audit
	utils.LogAudit(c, userModel.ID, "print_label", "tmt_bar_labels", &label.LabelID, 
		"Label print job created", map[string]interface{}{
			"print_job_id": printJobID.String(),
			"label_id":     label.LabelID,
		})

	c.JSON(http.StatusOK, gin.H{
		"message":     "Print job created successfully",
		"print_job_id": printJobID.String(),
		"zpl_content":  label.ZPLContent,
	})
}

// ExportTMTBarLabelsCSV exports TMT Bar labels as CSV
func ExportTMTBarLabelsCSV(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Build query for CSV export
	query := `SELECT label_id, location, bundle_nos, pqd, unit, time1, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date1, printed_at, status, is_duplicate, created_at 
			  FROM tmt_bar_labels WHERE 1=1`
	args := []interface{}{}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND user_id = $1"
		args = append(args, userModel.ID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch labels for export"})
		return
	}
	defer rows.Close()

	// Generate CSV
	csvData := utils.GenerateTMTBarLabelsCSV(rows)

	// Set response headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=tmt_bar_labels.csv")

	// Log audit
	utils.LogAudit(c, userModel.ID, "export_csv", "tmt_bar_labels", nil, "Exported labels to CSV")

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
} 