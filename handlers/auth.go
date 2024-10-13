// Packages
package handlers

// Imports
import (
	"crm-backend/db"
	"crm-backend/models"
	"crm-backend/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

	// Validate request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid email address"})
		return
	}
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password length should be minimum 6 characters"})
		return
	}

	user.Password = hashPasswordMD5(user.Password)
	user.OTP = utils.GenerateOTP()

	// Save the user in the database
	createdData, err := db.InsertUser(user)
	if err != nil {
		// Check if the error message contains the unique constraint name
		if strings.Contains(err.Error(), "unique constraint \"users_email_key\"") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists!"}) // Changed to bad request
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	fmt.Println(createdData)

	c.JSON(http.StatusCreated, gin.H{"message": "OTP successfully sent on your email"})
}

func ResendOTP(c *gin.Context) {
	var user models.User

	// Validate request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid email address"})
		return
	}

	// Save the user in the database
	userId, err := db.FindUserIdByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find email"})
		return
	}

	user.OTP = utils.GenerateOTP()
	err = db.UpdateUserOTPByEmailId(userId, user.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP Send failed !"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "OTP successfully sent on your email"})
}

func VerifyOTP(c *gin.Context) {
	var user models.User

	// Validate request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid email address"})
		return
	}
	if len(user.OTP) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter valid 4 digit OTP"})
		return
	}

	// Save the user in the database
	userId, otp, err := db.FindOTPByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find email"})
		return
	}
	if user.OTP != otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP, Please try again later"})
		return
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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
	userId, err := db.FindUserIdByEmail(user.Email)
	if err != nil || !foundUser.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or account not verified"})
		return
	}
	fmt.Println(userId)

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
	userId, err := db.FindUserIdByEmail(user.Email)
	if err != nil || user.OTP != foundUser.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or OTP"})
		return
	}
	fmt.Println((userId))

	// Generate JWT token on successful OTP authentication
	token, err := utils.GenerateJWT(foundUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func isValidEmail(email string) bool {
	// Define a regex pattern for a valid email
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
