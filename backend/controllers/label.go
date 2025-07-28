package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"labelops-backend/db"
	"labelops-backend/models"
	"labelops-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// func ProcessLabelBatch(c *gin.Context)    { c.JSON(501, gin.H{"error": "Not implemented"}) }
// func GetLabels(c *gin.Context)            { c.JSON(501, gin.H{"error": "Not implemented"}) }
// GetLabelByID retrieves a specific label
func GetLabelByID(c *gin.Context) {
	labelID := c.Param("id")

	// Parse label ID
	labelUUID, err := uuid.Parse(labelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	// Get label
	var label models.Label
	err = db.DB.QueryRow(
		`SELECT id, label_id, location, bundle_nos, pqd, unit, time1, length, 
		 heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
		 url_apikey, weight, section, date1, printed_at, user_id, status, 
		 zpl_content, qr_code, is_duplicate, created_at, updated_at 
		 FROM labels WHERE id = $1`,
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

	c.JSON(http.StatusOK, label)
}

// PrintLabel prints a specific label
func PrintLabel(c *gin.Context) {
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
	var label models.Label
	err = db.DB.QueryRow(
		`SELECT id, label_id, location, bundle_nos, pqd, unit, time1, length, 
		 heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
		 url_apikey, weight, section, date1, printed_at, user_id, status, 
		 zpl_content, qr_code, is_duplicate, created_at, updated_at 
		 FROM labels WHERE id = $1`,
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
		label.ZPLContent = utils.GenerateLabelZPL(label)
		// Update label with ZPL content
		db.DB.Exec("UPDATE labels SET zpl_content = $1 WHERE id = $2", label.ZPLContent, label.ID)
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
	utils.LogAudit(c, userModel.ID, "print_label", "labels", &label.LabelID,
		"Label print job created", map[string]interface{}{
			"print_job_id": printJobID.String(),
			"label_id":     label.LabelID,
		})

	c.JSON(http.StatusOK, gin.H{
		"message":      "Print job created successfully",
		"print_job_id": printJobID.String(),
		"zpl_content":  label.ZPLContent,
	})
}

// ExportLabelsCSV exports labels as CSV
func ExportLabelsCSV(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Build query for CSV export
	query := `SELECT label_id, location, bundle_nos, pqd, unit, time1, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date1, printed_at, status, is_duplicate, created_at 
			  FROM labels WHERE 1=1`
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
	csvData := utils.GenerateLabelsCSV(rows)

	// Set response headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=labels.csv")

	// Log audit
	utils.LogAudit(c, userModel.ID, "export_csv", "labels", nil, "Exported labels to CSV")

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}

// GetPrintJobs retrieves print jobs
func GetPrintJobs(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	query := `SELECT id, label_id, user_id, status, zpl_content, max_retries, 
			  retry_count, error_message, created_at, updated_at 
			  FROM print_jobs WHERE 1=1`
	args := []interface{}{}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND user_id = $1"
		args = append(args, userModel.ID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print jobs"})
		return
	}
	defer rows.Close()

	var printJobs []map[string]interface{}
	for rows.Next() {
		var (
			id, labelID, userID, status, zplContent, errorMessage string
			maxRetries, retryCount                                int
			createdAt, updatedAt                                  sql.NullTime
		)
		err := rows.Scan(
			&id, &labelID, &userID, &status, &zplContent, &maxRetries, &retryCount,
			&errorMessage, &createdAt, &updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan print job"})
			return
		}

		job := map[string]interface{}{
			"id":            id,
			"label_id":      labelID,
			"user_id":       userID,
			"status":        status,
			"zpl_content":   zplContent,
			"max_retries":   maxRetries,
			"retry_count":   retryCount,
			"error_message": errorMessage,
			"created_at":    createdAt.Time,
			"updated_at":    updatedAt.Time,
		}
		printJobs = append(printJobs, job)
	}

	c.JSON(http.StatusOK, gin.H{"print_jobs": printJobs, "count": len(printJobs)})
}

// GetPrintJobByID retrieves a specific print job
func GetPrintJobByID(c *gin.Context) {
	jobID := c.Param("id")

	// Parse job ID
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid print job ID"})
		return
	}

	// Get print job
	var (
		id, labelID, userID, status, zplContent, errorMessage string
		maxRetries, retryCount                                int
		createdAt, updatedAt                                  sql.NullTime
	)
	err = db.DB.QueryRow(
		`SELECT id, label_id, user_id, status, zpl_content, max_retries, 
		 retry_count, error_message, created_at, updated_at 
		 FROM print_jobs WHERE id = $1`,
		jobUUID,
	).Scan(
		&id, &labelID, &userID, &status, &zplContent, &maxRetries, &retryCount,
		&errorMessage, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Print job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print job"})
		return
	}

	job := map[string]interface{}{
		"id":            id,
		"label_id":      labelID,
		"user_id":       userID,
		"status":        status,
		"zpl_content":   zplContent,
		"max_retries":   maxRetries,
		"retry_count":   retryCount,
		"error_message": errorMessage,
		"created_at":    createdAt.Time,
		"updated_at":    updatedAt.Time,
	}

	c.JSON(http.StatusOK, job)
}

// RetryPrintJob retries a failed print job
func RetryPrintJob(c *gin.Context) {
	jobID := c.Param("id")
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Parse job ID
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid print job ID"})
		return
	}

	// Update print job status
	_, err = db.DB.Exec(
		`UPDATE print_jobs SET status = 'pending', retry_count = retry_count + 1, 
		 updated_at = NOW() WHERE id = $1 AND user_id = $2`,
		jobUUID, userModel.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry print job"})
		return
	}

	// Log audit
	utils.LogAudit(c, userModel.ID, "retry_print_job", "print_jobs", &jobID, "Print job retry initiated")

	c.JSON(http.StatusOK, gin.H{"message": "Print job retry initiated"})
}

// GetAllUsers retrieves all users (admin only)
func GetAllUsers(c *gin.Context) {
	rows, err := db.DB.Query(
		`SELECT id, email, first_name, last_name, role, is_active, 
		 last_login, created_at, updated_at 
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role,
			&user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "count": len(users)})
}

// CreateUser creates a new user (admin only)
func CreateUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert new user
	var user models.User
	err = db.DB.QueryRow(
		`INSERT INTO users (email, password_hash, first_name, last_name, role) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, email, first_name, last_name, role, is_active, created_at`,
		req.Email, string(hashedPassword), req.FirstName, req.LastName, req.Role,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.IsActive, &user.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Log audit
	currentUser, _ := c.Get("user")
	adminUser := currentUser.(models.User)
	idStr := user.ID.String()
	utils.LogAudit(c, adminUser.ID, "create_user", "users", &idStr, "User created by admin")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// UpdateUser updates a user (admin only)
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Update user
	_, err = db.DB.Exec(
		`UPDATE users SET first_name = $1, last_name = $2, email = $3, 
		 updated_at = NOW() WHERE id = $4`,
		req.FirstName, req.LastName, req.Email, userUUID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Log audit
	user, _ := c.Get("user")
	adminUser := user.(models.User)
	utils.LogAudit(c, adminUser.ID, "update_user", "users", &userID, "User updated by admin")

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete user
	_, err = db.DB.Exec("DELETE FROM users WHERE id = $1", userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Log audit
	user, _ := c.Get("user")
	adminUser := user.(models.User)
	utils.LogAudit(c, adminUser.ID, "delete_user", "users", &userID, "User deleted by admin")

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetSystemStats retrieves system statistics (admin only)
func GetSystemStats(c *gin.Context) {
	// Get total users
	var totalUsers int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}

	// Get total labels
	var totalLabels int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels").Scan(&totalLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get label count"})
		return
	}

	// Get total print jobs
	var totalPrintJobs int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM print_jobs").Scan(&totalPrintJobs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get print job count"})
		return
	}

	// Get recent activity
	var recentLabels int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE created_at >= NOW() - INTERVAL '24 hours'").Scan(&recentLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent label count"})
		return
	}

	stats := gin.H{
		"total_users":      totalUsers,
		"total_labels":     totalLabels,
		"total_print_jobs": totalPrintJobs,
		"recent_labels":    recentLabels,
	}

	c.JSON(http.StatusOK, stats)
}

// GetAuditLogs retrieves audit logs with filtering
func GetAuditLogs(c *gin.Context) {
	utils.GetAuditLogs(c)
}

// ExportAuditLogsCSV exports audit logs as CSV
func ExportAuditLogsCSV(c *gin.Context) {
	utils.ExportAuditLogsCSV(c)
}

// BatchLabelProcess processes a batch of labels
func BatchLabelProcess(c *gin.Context) {
	var req models.LabelBatchRequest
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
	err = db.DB.QueryRow("SELECT batch_label_process($1, $2)", labelsJSON, userModel.ID).Scan(&resultJSON)
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
	utils.LogAudit(c, userModel.ID, "process_batch", "labels", nil,
		"Processed batch of labels", map[string]interface{}{
			"total_processed": result["total_processed"],
			"new_count":       result["new_count"],
			"duplicate_count": result["duplicate_count"],
		})

	c.JSON(http.StatusOK, result)
}

// GetLabels retrieves labels with filtering
func GetLabels(c *gin.Context) {
	// Safely get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userModel, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	// Parse query parameters with defaults
	limit := 50
	if limitStr := c.DefaultQuery("limit", "50"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.DefaultQuery("offset", "0"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	status := c.Query("status")
	grade := c.Query("grade")
	section := c.Query("section")

	// Build base query
	query := `SELECT id, label_id, location, bundle_nos, pqd, unit, time1, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date1, printed_at, user_id, status, 
			  zpl_content, qr_code, is_duplicate, created_at, updated_at 
			  FROM labels WHERE 1=1`

	var args []interface{}
	argCount := 1

	// Add filters safely
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}
	if grade != "" {
		query += fmt.Sprintf(" AND grade = $%d", argCount)
		args = append(args, grade)
		argCount++
	}
	if section != "" {
		query += fmt.Sprintf(" AND section = $%d", argCount)
		args = append(args, section)
		argCount++
	}

	// Add user filter for non-admin users
	if userModel.Role != "admin" && userModel.Role != "user" {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, userModel.ID)
		argCount++
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Execute query
	log.Printf("Executing query: %s", query)
	log.Printf("Query args: %v", args)
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch labels", "details": err.Error()})
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Failed to get columns: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get columns", "details": err.Error()})
		return
	}

	var labels []map[string]interface{}

	// Process each row
	rowCount := 0
	for rows.Next() {
		rowCount++
		log.Printf("Processing row %d", rowCount)

		// Create slices to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Row scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan label row", "details": err.Error()})
			return
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Handle different types properly
			switch v := val.(type) {
			case []byte:
				// Convert UUID and other binary types to string
				row[col] = string(v)
			case time.Time:
				// Format time properly
				row[col] = v.Format(time.RFC3339)
			case nil:
				row[col] = nil
			default:
				row[col] = v
			}
		}
		log.Printf("Row %d data: %v", rowCount, row)
		labels = append(labels, row)
	}
	log.Printf("Total rows processed: %d", rowCount)

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over labels", "details": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"labels":    labels,
		"count":     len(labels),
		"limit":     limit,
		"offset":    offset,
		"user_id":   userModel.ID,
		"user_role": userModel.Role,
		"debug":     "GetLabels function executed at " + time.Now().Format("2006-01-02 15:04:05"),
	})
}
