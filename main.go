// Packages
package main

// Imports
import (
	"log"
	"os"

	"crm-backend/db"
	"crm-backend/handlers"
	"crm-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db.ConnectDB()     // Initialize the database connection
	defer db.CloseDB() // Ensures the DB connection is closed when the application shuts down

	r := gin.Default()       // Initialize the Gin router
	r.Use(middleware.CORS()) // Apply the CORS middleware

	r.Use(middleware.AuthRequired())

	// Routes for authentication
	r.POST("/auth/signup", handlers.SignUp)
	r.POST("/auth/resendOTP", handlers.ResendOTP)
	r.POST("/auth/verifyOTP", handlers.VerifyOTP)
	r.POST("/auth/signIn", handlers.SignIn)

	r.GET("/company/categoryList", handlers.CategoryList)

	r.POST("/ticket/create", handlers.CreateTicket)
	r.GET("/ticket/list", handlers.GetTickets)

	// Get the port from environment variable or use 8080 by default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start the server on the specified port
	if err := r.Run(":" + port); err != nil { // Listen and serve on the specified port
		log.Fatalf("Error starting server: %v", err)
	}
}
