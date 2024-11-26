// Packages
package customer

// Imports
import (
	"crm-backend/db"
	"crm-backend/models"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func Create(c *gin.Context) {
	var customer models.Customer

	// Validate request
	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !isValidEmail(customer.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid email address"})
		return
	}
	if len(customer.Pincode) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid 6 digit pincode"})
		return
	}
	if len(customer.MobileNumber) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid 10 digit mobile number"})
		return
	}

	err := db.CreateCustomer(customer)
	if err != nil {
		// Check if it's a PostgreSQL error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique constraint violation
				c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number already exists !"})
				return
			}
		}

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return

	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Customer created successfully !"})
}

func List(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}

func isValidEmail(email string) bool {
	// Define a regex pattern for a valid email
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
