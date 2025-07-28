package main

import (
	"log"
	"os"

	"labelops-backend/controllers"
	"labelops-backend/db"
	"labelops-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.Default())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "LabelOps Backend is running",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/login", controllers.Login)
		api.POST("/auth/register", controllers.Register)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Label routes
			protected.POST("/labels/batch", controllers.BatchLabelProcess)
			protected.GET("/labels", controllers.GetLabels)
			protected.GET("/labels/:id", controllers.GetLabelByID)
			protected.POST("/labels/:id/print", controllers.PrintLabel)
			protected.GET("/labels/export/csv", controllers.ExportLabelsCSV)

			// Print job routes
			protected.GET("/print-jobs", controllers.GetPrintJobs)
			protected.GET("/print-jobs/:id", controllers.GetPrintJobByID)
			protected.POST("/print-jobs/:id/retry", controllers.RetryPrintJob)

			// User routes
			protected.GET("/users/profile", controllers.GetUserProfile)
			protected.PUT("/users/profile", controllers.UpdateUserProfile)

			// Audit log routes (fixed)
			protected.GET("/audit-logs", controllers.GetAuditLogs)
			protected.GET("/audit-logs/export/csv", controllers.ExportAuditLogsCSV)

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", controllers.GetAllUsers)
				admin.POST("/users", controllers.CreateUser)
				admin.PUT("/users/:id", controllers.UpdateUser)
				admin.DELETE("/users/:id", controllers.DeleteUser)
				admin.GET("/stats", controllers.GetSystemStats)
			}
		}

		// Get port from environment or use default
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Starting LabelOps Backend on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
}
