package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"labelops-backend/db"
	"labelops-backend/internal/printer"
	"labelops-backend/models"
	"labelops-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUserFromContext extracts the authenticated user from Gin context
// Returns the user model and a boolean indicating success; on failure, writes the appropriate response
func getUserFromContext(c *gin.Context) (models.User, bool) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return models.User{}, false
	}
	userModel, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
		return models.User{}, false
	}
	return userModel, true
}

// ensurePrinterDirectories creates required folders used by the printer package
// printers/bat -> where print batch scripts live
func ensurePrinterDirectories() error {
    if err := os.MkdirAll("printers/bat", 0755); err != nil {
        return err
    }
    return nil
}

// convertLabelDataToLabelWithID maps incoming LabelData to models.Label using an existing DB UUID
// This is used for ZPL generation where a full models.Label is expected
func convertLabelDataToLabelWithID(data models.LabelData, userID uuid.UUID, id uuid.UUID) models.Label {
    return models.Label{
        ID:             id,
        LabelID:        data.ID,
        Location:       data.LOCATION,
        BundleNo:       data.BUNDLE_NO,
        BundleType:     data.BUNDLE_TYPE,
        PQD:            data.PQD,
        Unit:           data.UNIT,
        Time:           data.TIME,
        Length:         data.LENGTH,
        HeatNo:         data.HEAT_NO,
        ProductHeading: data.PRODUCT_HEADING,
        IsiBottom:      data.ISI_BOTTOM,
        IsiTop:         data.ISI_TOP,
        ChargeDtm:      "",
        Mill:           data.MILL,
        Grade:          data.GRADE,
        UrlApikey:      data.URL_APIKEY,
        Weight:         data.WEIGHT,
        Section:        data.SECTION,
        Date:           data.DATE,
        UserID:         userID,
        Status:         "success",
        IsDuplicate:    false,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }
}

// createPrintJob inserts a new print job using the actual DB label UUID, user ID, heat number,
// and ZPL content read from the provided file path. Returns the new job ID as string.
func createPrintJob(labelID uuid.UUID, userID uuid.UUID, heatNo string, zplFilePath string, actualLabelID string) (string, error) {
    // Read ZPL content from file path generated earlier
    content, err := os.ReadFile(zplFilePath)
    if err != nil {
        return "", fmt.Errorf("failed to read ZPL file: %w", err)
    }

    jobID := uuid.New()
    _, err = db.DB.Exec(`
        INSERT INTO print_jobs (id, label_id, heat_no, actual_label_id, user_id, status, zpl_content, max_retries)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, jobID, labelID, heatNo, actualLabelID, userID, "pending", string(content), 3)
    if err != nil {
        return "", fmt.Errorf("failed to insert print job: %w", err)
    }
    return jobID.String(), nil
}

// updatePrintJobStatus updates the status and optional error message for a job
func updatePrintJobStatus(jobID string, status string, errorMessage string) {
    // Best-effort update; log on failure but do not interrupt caller flow
    id, err := uuid.Parse(jobID)
    if err != nil {
        log.Printf("updatePrintJobStatus: invalid job id %s: %v", jobID, err)
        return
    }
    _, err = db.DB.Exec(`
        UPDATE print_jobs
        SET status = $1, error_message = NULLIF($2, ''), updated_at = NOW()
        WHERE id = $3
    `, status, errorMessage, id)
    if err != nil {
        log.Printf("updatePrintJobStatus: failed updating job %s: %v", jobID, err)
    }
}

// nilIfInvalidString converts sql.NullString to either its string or nil for JSON
func nilIfInvalidString(ns sql.NullString) interface{} {
    if ns.Valid {
        return ns.String
    }
    return nil
}

// nilIfInvalidTime converts sql.NullTime to RFC3339 string or nil for JSON
func nilIfInvalidTime(nt sql.NullTime) interface{} {
    if nt.Valid {
        return nt.Time.Format(time.RFC3339)
    }
    return nil
}

// BatchLabelProcess processes a batch of labels and sends new labels to printer
func BatchLabelProcess(c *gin.Context) {
	var req models.LabelBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate request
	if len(req.Labels) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No labels provided"})
		return
	}

	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	// Ensure printer directories exist
	if err := ensurePrinterDirectories(); err != nil {
		log.Printf("Failed to create printer directories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize printer system", "details": err.Error()})
		return
	}

	// Convert labels to JSON for DB stored procedure
	labelsJSON, err := json.Marshal(req.Labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal labels", "details": err.Error()})
		return
	}

	// Process batch in database
	var resultStr string
	err = db.DB.QueryRow("SELECT batch_label_process($1, $2)", labelsJSON, userModel.ID).Scan(&resultStr)
	if err != nil {
		log.Printf("Database batch processing failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process batch", "details": err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse batch result", "details": err.Error()})
		return
	}

	// Extract new labels for printing - these should contain the DB-generated IDs (business IDs provided)
	newLabelsInterface, exists := result["new_labels"]
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid batch result: missing new_labels"})
		return
	}

	// Convert new labels to slice of maps to access both label data and business IDs
	newLabelsJSON, err := json.Marshal(newLabelsInterface)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process new labels", "details": err.Error()})
		return
	}

	var newLabelsWithIDs []map[string]interface{}
	if err := json.Unmarshal(newLabelsJSON, &newLabelsWithIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse new labels", "details": err.Error()})
		return
	}

	// Generate ZPL files and create print jobs only for NEW labels
	var zplPaths []string
	var printJobIDs []string

	for _, labelMap := range newLabelsWithIDs {
		// Extract the business ID (bundle number / external label ID)
		businessID, ok := labelMap["ID"].(string)
		if !ok {
			log.Printf("Invalid business ID in new_labels result: %v", labelMap["ID"])
			continue
		}

		// Query the database to get the actual UUID for this label
		var labelUUID uuid.UUID
		err := db.DB.QueryRow(`
			SELECT id FROM labels 
			WHERE label_id = $1 AND user_id = $2 
			ORDER BY created_at DESC 
			LIMIT 1
		`, businessID, userModel.ID).Scan(&labelUUID)
		if err != nil {
			log.Printf("Failed to find database UUID for label %s: %v", businessID, err)
			continue
		}

		// Convert the label map back to LabelData for ZPL generation
		labelDataJSON, err := json.Marshal(labelMap)
		if err != nil {
			log.Printf("Failed to marshal label data: %v", err)
			continue
		}

		var labelData models.LabelData
		if err := json.Unmarshal(labelDataJSON, &labelData); err != nil {
			log.Printf("Failed to unmarshal label data: %v", err)
			continue
		}

		// Convert LabelData to Label for ZPL generation, using the DB ID
		label := convertLabelDataToLabelWithID(labelData, userModel.ID, labelUUID)

		// Generate ZPL file
		zplPath, err := printer.GenerateAndSaveZPL(label)
		if err != nil {
			log.Printf("Failed to generate ZPL for label %s: %v", labelData.PQD, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate ZPL file",
				"label":   labelData.PQD,
				"details": err.Error(),
			})
			return
		}
		zplPaths = append(zplPaths, zplPath)

		// Create print job record in database using the actual DB label ID and store business ID as actual_label_id
		// Pass heat number for the NOT NULL heat_no column
		printJobID, err := createPrintJob(labelUUID, userModel.ID, label.HeatNo, zplPath, businessID)
		if err != nil {
			log.Printf("Failed to create print job for label %s: %v", businessID, err)
			// Continue processing other labels, but log the error
		} else {
			printJobIDs = append(printJobIDs, printJobID)
		}
	}

	// Print all ZPL files in batch if any new labels exist
	var printError error
	if len(zplPaths) > 0 {
		if err := printer.PrintZPLBatch(zplPaths); err != nil {
			printError = err
			log.Printf("Batch printing failed: %v", err)

			// Update print job statuses to failed
			for _, jobID := range printJobIDs {
				updatePrintJobStatus(jobID, "failed", err.Error())
			}
		} else {
			// Update print job statuses to completed
			for _, jobID := range printJobIDs {
				updatePrintJobStatus(jobID, "completed", "")
			}
		}
	}

	// Audit logging
	utils.LogAudit(c, userModel.ID, "process_batch", "labels", nil,
		"Processed batch of labels", map[string]interface{}{
			"total_processed":    result["total_processed"],
			"new_count":          result["new_count"],
			"duplicate_count":    result["duplicate_count"],
			"print_jobs_created": len(printJobIDs),
		})

	// Prepare response
	response := gin.H{
		"message":            "Batch processed successfully",
		"total_processed":    result["total_processed"],
		"new_count":          result["new_count"],
		"duplicate_count":    result["duplicate_count"],
		"print_jobs_created": len(printJobIDs),
	}

	if printError != nil {
		response["print_warning"] = "Labels processed but printing failed: " + printError.Error()
		response["message"] = "Batch processed with print errors"
	} else if len(zplPaths) > 0 {
		response["message"] = "Batch processed and sent to printer"
	}

	c.JSON(http.StatusOK, response)
}

// GetPrintJobs retrieves print jobs (includes actual_label_id for UI)
func GetPrintJobs(c *gin.Context) {
	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	query := `SELECT id, label_id, actual_label_id, user_id, status, heat_no, error_message,
       zpl_content, max_retries, retry_count, created_at, updated_at
FROM print_jobs WHERE 1=1
`
	args := []interface{}{}
	if userModel.Role != "admin" {
		query += " AND user_id = $1"
		args = append(args, userModel.ID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print jobs", "details": err.Error()})
		return
	}
	defer rows.Close()

	var printJobs []map[string]interface{}
	for rows.Next() {
		var (
			id, labelID, actualLabelID, userID, status, heatNo string
			errorMessage                                       sql.NullString
			zplContent                                         string
			maxRetries                                         int
			retryCount                                         int
			createdAt, updatedAt                               sql.NullTime
		)
		err := rows.Scan(
			&id, &labelID, &actualLabelID, &userID, &status, &heatNo, &errorMessage, &zplContent, &maxRetries,
			&retryCount, &createdAt, &updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan print job", "details": err.Error()})
			return
		}

		job := map[string]interface{}{
			"id":              id,
			"label_id":        labelID,
			"actual_label_id": actualLabelID,
			"user_id":         userID,
			"status":          status,
			"heat_no":         heatNo,
			"error_message":   nilIfInvalidString(errorMessage),
			"zpl_content":     zplContent,
			"max_retries":     maxRetries,
			"retry_count":     retryCount,
			"created_at":      nilIfInvalidTime(createdAt),
			"updated_at":      nilIfInvalidTime(updatedAt),
		}
		printJobs = append(printJobs, job)
	}

	c.JSON(http.StatusOK, gin.H{"print_jobs": printJobs, "count": len(printJobs)})
}

// GetLabels retrieves labels with filtering
func GetLabels(c *gin.Context) {
	// Safely get user from context
	userModel, ok := getUserFromContext(c)
	if !ok {
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
	query := `SELECT id, label_id, location, bundle_no, pqd, unit, time, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date, user_id, status, 
			  is_duplicate, created_at, updated_at 
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

	// Restrict non-admin users to their own records
	if userModel.Role != "admin" {
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

// GetPrintJobByID retrieves a specific print job

// GetPrintJobByID retrieves a specific print job with debug logs
func GetPrintJobByID(c *gin.Context) {
	jobID := c.Param("id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid print job ID"})
		return
	}

	query := `
        SELECT id, label_id, user_id, status, zpl_content, max_retries, 
           retry_count, error_message, actual_label_id, heat_no, created_at, updated_at 
    FROM print_jobs WHERE id = $1
    `
	log.Printf("Executing SQL Query: %s", query)
	log.Printf("With Parameter jobID: %s", jobID)

	var (
		id, labelID, userID, status, zplContent string
		maxRetries, retryCount                  int
		errorMessage                            sql.NullString
		actualLabelID                           sql.NullString
		heatNoCol                               string
		createdAt, updatedAt                    sql.NullTime
	)

	err := db.DB.QueryRow(query, jobID).Scan(
		&id, &labelID, &userID, &status, &zplContent, &maxRetries, &retryCount,
		&errorMessage, &actualLabelID, &heatNoCol, &createdAt, &updatedAt,
	)

	if err != nil {
		log.Printf("Error querying print job: %v", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Print job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print job"})
		return
	}

	job := map[string]interface{}{
		"id":              id,
		"label_id":        labelID,
		"user_id":         userID,
		"status":          status,
		"zpl_content":     zplContent,
		"max_retries":     maxRetries,
		"retry_count":     retryCount,
		"error_message":   nilIfInvalidString(errorMessage),
		"heat_no":         heatNoCol,
		"actual_label_id": nilIfInvalidString(actualLabelID),
		"created_at":      nilIfInvalidTime(createdAt),
		"updated_at":      nilIfInvalidTime(updatedAt),
	}

	c.JSON(http.StatusOK, job)
}

// GetPrintJobByHeatNo retrieves a specific print job by heat number
func GetPrintJobsByHeatNo(c *gin.Context) {
	heatNo := c.Param("heatno")

	if heatNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid heat number"})
		return
	}

	query := `
        SELECT id, label_id, user_id, status, zpl_content, max_retries,
               retry_count, error_message, actual_label_id, heat_no, created_at, updated_at
        FROM print_jobs
        WHERE heat_no = $1
          AND created_at >= NOW() - INTERVAL '10 days'
    `
	log.Printf("Executing SQL Query: %s", query)
	log.Printf("With Parameter heatNo: %s", heatNo)

	rows, err := db.DB.Query(query, heatNo)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print jobs"})
		return
	}
	defer rows.Close()

	var jobs []map[string]interface{}

	for rows.Next() {
		var (
			id, labelID, userID, status, zplContent string
			maxRetries, retryCount                  int
			errorMessage                            sql.NullString
			actualLabelID                           sql.NullString
			heatNoCol                               string
			createdAt, updatedAt                    sql.NullTime
		)

		err := rows.Scan(
			&id, &labelID, &userID, &status, &zplContent, &maxRetries,
			&retryCount, &errorMessage, &actualLabelID, &heatNoCol, &createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		job := map[string]interface{}{
			"id":              id,
			"label_id":        labelID,
			"user_id":         userID,
			"status":          status,
			"zpl_content":     zplContent,
			"max_retries":     maxRetries,
			"retry_count":     retryCount,
			"error_message":   nilIfInvalidString(errorMessage),
			"heat_no":         heatNoCol,
			"actual_label_id": nilIfInvalidString(actualLabelID),
			"created_at":      nilIfInvalidTime(createdAt),
			"updated_at":      nilIfInvalidTime(updatedAt),
		}

		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No print jobs found in past 10 days"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// PrintLabel prints a specific label
func PrintLabel(c *gin.Context) {
	var request struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("PrintLabel: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	log.Printf("PrintLabel: Received request with label_id: %s", request.ID)

	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	// Fetch label using the label_id (string), but retrieve its UUID `id`
	var label models.Label
	err := db.DB.QueryRow(`
		SELECT id, label_id, location, bundle_no, pqd, unit, time, length,
		       heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade,
		       url_apikey, weight, section, date, user_id, status, is_duplicate,
		       created_at, updated_at
		FROM labels
		WHERE id = $1
	`, request.ID).Scan(
		&label.ID, &label.LabelID, &label.Location, &label.BundleNo, &label.PQD,
		&label.Unit, &label.Time, &label.Length, &label.HeatNo, &label.ProductHeading,
		&label.IsiBottom, &label.IsiTop, &label.ChargeDtm, &label.Mill, &label.Grade,
		&label.UrlApikey, &label.Weight, &label.Section, &label.Date,
		&label.UserID, &label.Status, &label.IsDuplicate,
		&label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("PrintLabel: No label found for label_id: %s", request.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Label not found",
				"details": fmt.Sprintf("No label found with label_id: %s", request.ID),
			})
			return
		}
		log.Printf("PrintLabel: Database error fetching label: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch label",
			"details": err.Error(),
		})
		return
	}

	log.Printf("PrintLabel: Found label UUID: %s", label.ID.String())

	// Generate ZPL content
	zplContent := utils.GenerateLabelZPL(label)
	log.Printf("PrintLabel: ZPL content generated (length: %d)", len(zplContent))
	log.Printf("PrintLabel: ZPL content generated: %s", zplContent)

	// Create print job
	printJobID := uuid.New()
	_, err = db.DB.Exec(`
		INSERT INTO print_jobs (id, label_id, heat_no, user_id, status, zpl_content, max_retries)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, printJobID, label.ID, label.HeatNo, userModel.ID, "pending", zplContent, 3)

	if err != nil {
		log.Printf("PrintLabel: Failed to insert print job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create print job",
			"details": err.Error(),
		})
		return
	}

	log.Printf("PrintLabel: Print job created with ID: %s", printJobID.String())

	// Audit log
	utils.LogAudit(c, userModel.ID, "print_label", "labels", &label.LabelID,
		"Label print job created", map[string]interface{}{
			"print_job_id": printJobID.String(),
			"label_id":     label.LabelID,
		})

	c.JSON(http.StatusOK, gin.H{
		"message":      "Print job created successfully",
		"print_job_id": printJobID.String(),
		"zpl_content":  zplContent,
	})
}

// ExportLabelsCSV exports labels as CSV
func ExportLabelsCSV(c *gin.Context) {
	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	query := `SELECT label_id, location, bundle_no, pqd, unit, time, length, 
			  heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
			  url_apikey, weight, section, date, status, is_duplicate, created_at 
			  FROM labels WHERE 1=1`
	args := []interface{}{}

	if userModel.Role != "admin" {
		query += " AND user_id = $1"
		args = append(args, userModel.ID)
	}

	query += " ORDER BY created_at DESC"

	log.Printf("Running Query: %s", query)
	log.Printf("With Args: %v", args)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch labels for export",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	csvData := utils.GenerateLabelsCSV(rows)

	if len(csvData) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No data to export"})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=labels.csv")

	utils.LogAudit(c, userModel.ID, "export_csv", "labels", nil, "Exported labels to CSV")

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}

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
		`SELECT id, label_id, location, bundle_no, pqd, unit, time, length, 
		 heat_no, product_heading, isi_bottom, isi_top, charge_dtm, mill, grade, 
		 url_apikey, weight, section, date, user_id, status, 
		 is_duplicate, created_at, updated_at 
		 FROM labels WHERE id = $1`,
		labelUUID,
	).Scan(
		&label.ID, &label.LabelID, &label.Location, &label.BundleNo, &label.PQD,
		&label.Unit, &label.Time, &label.Length, &label.HeatNo, &label.ProductHeading,
		&label.IsiBottom, &label.IsiTop, &label.ChargeDtm, &label.Mill, &label.Grade,
		&label.UrlApikey, &label.Weight, &label.Section, &label.Date,
		&label.UserID, &label.Status, &label.IsDuplicate,
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

// RetryPrintJob retries a failed print job
func RetryPrintJob(c *gin.Context) {
	var request struct {
		JobID string `json:"job_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	// Parse job ID
	jobUUID, err := uuid.Parse(request.JobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid print job ID"})
		return
	}

	// Update print job status and increment retry count
	_, err = db.DB.Exec(
		`UPDATE print_jobs
         SET retry_count = retry_count + 1, status = 'retrying', updated_at = NOW()
         WHERE id = $1 AND user_id = $2`,
		jobUUID, userModel.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry print job"})
		return
	}

	// Log audit
	utils.LogAudit(c, userModel.ID, "retry_print_job", "print_jobs", &request.JobID, "Print job retry initiated")

	// Return updated job
	c.JSON(http.StatusOK, gin.H{"message": "Print job retry initiated"})
}

// ExportPrintJobsCSV exports print jobs as CSV
func ExportPrintJobsCSV(c *gin.Context) {
	userModel, ok := getUserFromContext(c)
	if !ok {
		return
	}

	// Build query for print jobs
	query := `SELECT id, label_id, user_id,actual_label_id,heat_no,zpl_content, status, max_retries, retry_count,
		 created_at, updated_at
			  FROM print_jobs WHERE 1=1`
	args := []interface{}{}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND user_id = $1"
		args = append(args, userModel.ID)
	}

	// Add optional status filter
	if status := c.Query("status"); status != "" {
		if userModel.Role != "admin" {
			query += " AND status = $2"
		} else {
			query += " AND status = $1"
		}
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	log.Printf("Running Print Jobs Export Query: %s", query)
	log.Printf("With Args: %v", args)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch print jobs for export",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	csvData := utils.GeneratePrintJobsCSV(rows)

	if len(csvData) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No print jobs data to export"})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=print_jobs.csv")

	utils.LogAudit(c, userModel.ID, "export_csv", "print_jobs", nil, "Exported print jobs to CSV")

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}
