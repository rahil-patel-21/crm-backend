package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(email string, exp string) (string, error) {
	var expire_time any
	if exp == "min" {
		expire_time = time.Now().Add(time.Minute * 5).Unix()
	} else {
		expire_time = time.Now().Add(time.Hour * 72).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   expire_time,
	})

	return token.SignedString(jwtKey)
}

// ValidateJWT checks if the token is valid
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Return the key for verification
		return jwtKey, nil
	})
}
