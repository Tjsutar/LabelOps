package controllers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"labelops-backend/db"
	"labelops-backend/models"
	"labelops-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user authentication
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, email, password_hash, first_name, last_name, role, is_active FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update last login
	db.DB.Exec("UPDATE users SET last_login = NOW() WHERE id = $1", user.ID)

	// Log audit
	utils.LogAudit(c, user.ID, "login", "user", nil, "User logged in successfully")

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     tokenString,
		User:      user,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingID uuid.UUID
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
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
idStr := user.ID.String()
utils.LogAudit(c, user.ID, "register", "user", &idStr, "User registered successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUserProfile returns the current user's profile
func GetUserProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserProfile updates the current user's profile
func UpdateUserProfile(c *gin.Context) {
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	userModel := user.(models.User)

	// Update user
	_, err := db.DB.Exec(
		"UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = NOW() WHERE id = $4",
		req.FirstName, req.LastName, req.Email, userModel.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Log audit
idStr := userModel.ID.String()
utils.LogAudit(c, userModel.ID, "update_profile", "user", &idStr, "User profile updated")

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
} 