package middleware

import (
	"crm-backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication is required"})
			c.Abort()
			return
		}

		validatedToken, err := utils.ValidateJWT(tokenString)

		if err != nil || !validatedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
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
