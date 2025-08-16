package controllers

import (
	"fmt"
	"net/http"
	"time"

	"labelops-backend/db"
	"labelops-backend/models"
	"labelops-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

// GetDashboardStats retrieves comprehensive dashboard statistics
func GetDashboardStats(c *gin.Context) {
	// Get basic label counts
	var totalLabels, printedLabels, pendingLabels, failedLabels, duplicateLabels int

	// Total labels count
	err := db.DB.QueryRow("SELECT COUNT(*) FROM labels").Scan(&totalLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total labels count"})
		return
	}

	// Printed labels count (status = 'printed')
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE status = 'success'").Scan(&printedLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get printed labels count"})
		return
	}

	// Pending labels count (status = 'pending')
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE status = 'pending'").Scan(&pendingLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending labels count"})
		return
	}

	// Failed labels count (status = 'failed')
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE status = 'failed'").Scan(&failedLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get failed labels count"})
		return
	}

	// Duplicate labels count
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE is_duplicate = true").Scan(&duplicateLabels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get duplicate labels count"})
		return
	}

	// Get labels by grade
	gradeRows, err := db.DB.Query("SELECT grade, COUNT(*) FROM labels GROUP BY grade ORDER BY COUNT(*) DESC LIMIT 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get labels by grade"})
		return
	}
	defer gradeRows.Close()

	byGrade := make(map[string]int)
	for gradeRows.Next() {
		var grade string
		var count int
		if err := gradeRows.Scan(&grade, &count); err != nil {
			continue
		}
		byGrade[grade] = count
	}

	// Get labels by section
	sectionRows, err := db.DB.Query("SELECT section, COUNT(*) FROM labels GROUP BY section ORDER BY COUNT(*) DESC LIMIT 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get labels by section"})
		return
	}
	defer sectionRows.Close()

	bySection := make(map[string]int)
	for sectionRows.Next() {
		var section string
		var count int
		if err := sectionRows.Scan(&section, &count); err != nil {
			continue
		}
		bySection[section] = count
	}

	// Get recent activity (labels created in last 24 hours)
	var recentLabels int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM labels WHERE created_at >= NOW() - INTERVAL '24 hours'").Scan(&recentLabels)
	if err != nil {
		recentLabels = 0 // Default to 0 if query fails
	}

	// Get print success rate
	var totalPrintJobs, successfulPrintJobs int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM print_jobs").Scan(&totalPrintJobs)
	if err != nil {
		totalPrintJobs = 0
	}

	err = db.DB.QueryRow("SELECT COUNT(*) FROM print_jobs WHERE status = 'completed'").Scan(&successfulPrintJobs)
	if err != nil {
		successfulPrintJobs = 0
	}

	printSuccessRate := float64(0)
	if totalPrintJobs > 0 {
		printSuccessRate = float64(successfulPrintJobs) / float64(totalPrintJobs) * 100
	}

	// Get active users count
	var activeUsers int
	err = db.DB.QueryRow("SELECT COUNT(DISTINCT user_id) FROM labels WHERE created_at >= NOW() - INTERVAL '7 days'").Scan(&activeUsers)
	if err != nil {
		activeUsers = 0
	}

	// Create comprehensive dashboard response
	dashboardStats := gin.H{
		"overview": gin.H{
			"total_labels":     totalLabels,
			"printed_labels":   printedLabels,
			"pending_labels":   pendingLabels,
			"failed_labels":    failedLabels,
			"duplicate_labels": duplicateLabels,
		},
		"breakdown": gin.H{
			"by_grade":   byGrade,
			"by_section": bySection,
		},
		"activity": gin.H{
			"recent_labels_24h": recentLabels,
			"active_users_7d":   activeUsers,
		},
		"performance": gin.H{
			"print_success_rate": fmt.Sprintf("%.1f%%", printSuccessRate),
			"total_print_jobs":   totalPrintJobs,
		},
		"timestamp": time.Now().UTC(),
	}

	c.JSON(http.StatusOK, dashboardStats)
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
