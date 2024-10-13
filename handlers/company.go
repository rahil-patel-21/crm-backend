// Packages
package handlers

// Imports
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID   int64  `json:"id"`   // ID of the category
	Name string `json:"name"` // Name of the category
}

// SignUp handles user registration
func CategoryList(c *gin.Context) {

	categoryList := []Category{
		{ID: 1, Name: "IT"},
		{ID: 2, Name: "Finance"},
	}

	c.JSON(http.StatusOK, gin.H{"list": categoryList})
}
