// Packages
package main

// Imports
import (
	"log"
	"os"

	"crm-backend/db"
	"crm-backend/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database connection
	db.ConnectDB()
	defer db.CloseDB() // Ensure the DB connection is closed when the application shuts down

	// Initialize the Gin router
	r := gin.Default()

	// Routes for authentication
	r.POST("/signup", handlers.SignUp)            // Endpoint for user signup
	r.POST("/signin", handlers.SignIn)            // Endpoint for user login via email/password
	r.POST("/signin-otp", handlers.SignInWithOTP) // Endpoint for user login via OTP

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
