// Packages
package handlers

// Imports
import (
	"crm-backend/db"
	"crm-backend/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SignUp handles user registration
func CreateTicket(c *gin.Context) {
	var ticket models.Ticket

	// Validate request
	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !isValidEmail(ticket.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid email address"})
		return
	}

	// Save the user in the database
	createdData, err := db.InsertTicket(ticket)
	if err != nil {
		// Check if the error message contains the unique constraint name
		if strings.Contains(err.Error(), "unique constraint \"users_email_key\"") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists!"}) // Changed to bad request
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create this ticket"})
		return
	}
	fmt.Println(createdData)

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully !"})
}
