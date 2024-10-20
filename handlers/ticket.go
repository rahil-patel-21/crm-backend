// Packages
package handlers

// Imports
import (
	"crm-backend/db"
	"crm-backend/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Create new ticket
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

// Get list of ticket
func GetTickets(c *gin.Context) {
	pageSizeStr := c.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize >= 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageSize must be a number between 1 and 100"})
		return
	}
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a number more than 0"})
		return
	}
	created_by_str := c.Query("created_by")
	created_by, err := strconv.Atoi(created_by_str)
	if err != nil || created_by <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "created_by should be number more than 0"})
		return
	}

	// Get list of tickets with pagination
	count, rows, err := db.GetTicketList(page, pageSize, created_by)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusOK, gin.H{"count": count, "rows": []map[string]interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "rows": rows})

}
