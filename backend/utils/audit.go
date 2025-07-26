package utils

import (
	"encoding/json"
	"net/http"

	"labelops-backend/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogAudit logs an audit entry to the database
func LogAudit(c *gin.Context, userID uuid.UUID, action, resource string, resourceID *string, details string, metadata ...map[string]interface{}) {
	// Get client IP
	ipAddress := c.ClientIP()
	if ipAddress == "" {
		ipAddress = c.GetHeader("X-Forwarded-For")
	}

	// Get user agent
	userAgent := c.GetHeader("User-Agent")

	// Prepare details with metadata
	finalDetails := details
	if len(metadata) > 0 && metadata[0] != nil {
		if metadataJSON, err := json.Marshal(metadata[0]); err == nil {
			finalDetails = details + " | Metadata: " + string(metadataJSON)
		}
	}

	// Insert audit log
	_, err := db.DB.Exec(
		`INSERT INTO audit_logs (user_id, action, resource, resource_id, details, ip_address, user_agent) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, action, resource, resourceID, finalDetails, ipAddress, userAgent,
	)

	if err != nil {
		// Log error but don't fail the request
		// In production, you might want to use a proper logging framework
		// log.Printf("Failed to log audit: %v", err)
	}
}

// GetAuditLogs retrieves audit logs with filtering
func GetAuditLogs(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	action := c.Query("action")
	resource := c.Query("resource")
	userID := c.Query("user_id")

	// Build query
	query := `SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, al.details, 
			  al.ip_address, al.user_agent, al.created_at, u.email, u.first_name, u.last_name 
			  FROM audit_logs al 
			  LEFT JOIN users u ON al.user_id = u.id 
			  WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add filters
	if action != "" {
		query += " AND al.action = $" + strconv.Itoa(argCount)
		args = append(args, action)
		argCount++
	}
	if resource != "" {
		query += " AND al.resource = $" + strconv.Itoa(argCount)
		args = append(args, resource)
		argCount++
	}
	if userID != "" {
		query += " AND al.user_id = $" + strconv.Itoa(argCount)
		args = append(args, userID)
		argCount++
	}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND al.user_id = $" + strconv.Itoa(argCount)
		args = append(args, userModel.ID)
		argCount++
	}

	query += " ORDER BY al.created_at DESC LIMIT $" + strconv.Itoa(argCount) + " OFFSET $" + strconv.Itoa(argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}
	defer rows.Close()

	var auditLogs []map[string]interface{}
	for rows.Next() {
		var log map[string]interface{}
		var userEmail, firstName, lastName sql.NullString
		
		err := rows.Scan(
			&log["id"], &log["user_id"], &log["action"], &log["resource"], &log["resource_id"],
			&log["details"], &log["ip_address"], &log["user_agent"], &log["created_at"],
			&userEmail, &firstName, &lastName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan audit log"})
			return
		}

		// Add user info
		if userEmail.Valid {
			log["user_email"] = userEmail.String
			log["user_name"] = firstName.String + " " + lastName.String
		}

		auditLogs = append(auditLogs, log)
	}

	c.JSON(http.StatusOK, gin.H{"audit_logs": auditLogs, "count": len(auditLogs)})
}

// ExportAuditLogsCSV exports audit logs as CSV
func ExportAuditLogsCSV(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Build query for CSV export
	query := `SELECT al.action, al.resource, al.resource_id, al.details, al.ip_address, 
			  al.created_at, u.email, u.first_name, u.last_name 
			  FROM audit_logs al 
			  LEFT JOIN users u ON al.user_id = u.id 
			  WHERE 1=1`
	args := []interface{}{}

	// Add user filter for non-admin users
	if userModel.Role != "admin" {
		query += " AND al.user_id = $1"
		args = append(args, userModel.ID)
	}

	query += " ORDER BY al.created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs for export"})
		return
	}
	defer rows.Close()

	// Generate CSV
	csvData := generateAuditLogsCSV(rows)

	// Set response headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")

	// Log audit
	LogAudit(c, userModel.ID, "export_csv", "audit_logs", nil, "Exported audit logs to CSV")

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}

// generateAuditLogsCSV generates CSV data for audit logs
func generateAuditLogsCSV(rows interface{}) string {
	// This would be implemented to convert database rows to CSV format
	// For now, returning a placeholder
	return "Action,Resource,Resource ID,Details,IP Address,Created At,User Email,User Name\n"
} 