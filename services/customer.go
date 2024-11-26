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
	"strconv"

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
	pageSizeStr := c.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageSize must be a number more than 0"})
		return
	}
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a number more than 0"})
		return
	}
	searchText := c.Query("searchText")

	count, rows, err := db.GetCustomerList(page, pageSize, searchText)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, gin.H{"count": count, "rows": []map[string]interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "rows": rows})
}

func isValidEmail(email string) bool {
	// Define a regex pattern for a valid email
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
