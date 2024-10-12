package handlers

import (
	"crm-backend/db"
	"crm-backend/models"
	"crm-backend/utils"
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// hashPasswordMD5 hashes a plain-text password using MD5
func hashPasswordMD5(password string) string {
	hash := md5.Sum([]byte(password))  // Hash the password using MD5
	return hex.EncodeToString(hash[:]) // Convert hash to hex string
}

// SignUp handles user registration
func SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password using MD5
	user.Password = hashPasswordMD5(user.Password)
	user.IsVerified = false

	// Save the user in the database
	_, err := db.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Generate OTP and send it to the user's email
	otp := utils.GenerateOTP()
	user.OTP = otp
	utils.SendEmail(user.Email, otp)

	c.JSON(http.StatusOK, gin.H{"message": "User created, verify your email with the OTP sent"})
}

// SignIn handles user login using email and password
func SignIn(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the user exists in the database
	err := db.FindUserByEmail(user.Email, &foundUser)
	if err != nil || !foundUser.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or account not verified"})
		return
	}

	// Compare the provided password with the stored hashed password
	if hashPasswordMD5(user.Password) != foundUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token on successful authentication
	token, err := utils.GenerateJWT(foundUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SignInWithOTP handles user login using email and OTP
func SignInWithOTP(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the user exists and the OTP matches
	err := db.FindUserByEmail(user.Email, &foundUser)
	if err != nil || user.OTP != foundUser.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or OTP"})
		return
	}

	// Generate JWT token on successful OTP authentication
	token, err := utils.GenerateJWT(foundUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
