package middleware

import (
	"crm-backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var exceptionUrls = []string{
	"/auth/signup",
	"/auth/resendOTP",
	"/auth/verifyOTP",
	"/auth/signIn",
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// If the requested URL is in the exception list, skip authentication
		requestedURL := c.FullPath() // Get the requested URL
		if isExceptionURL(requestedURL) {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // Removing "Bearer " as the token is Bearer Token

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized"})
			c.Abort()
			return
		}

		validatedToken, err := utils.ValidateJWT(tokenString)

		if err != nil || !validatedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized"})
			c.Abort()
			return
		}
		claims := validatedToken.Claims.(jwt.MapClaims)
		retrievedEmail, ok := claims["email"].(string)
		if !ok {
			log.Fatalf("Email not found or not a string")
		}
		c.Set("email", retrievedEmail)
		c.Next()
	}
}

// Check if the current URL is in the exceptionUrls list
func isExceptionURL(requestedURL string) bool {
	for _, url := range exceptionUrls {
		if url == requestedURL {
			return true
		}
	}
	return false
}
