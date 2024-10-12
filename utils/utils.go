// Packages
package utils

// Imports
import (
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

// GenerateOTP generates a random 6-digit numeric OTP
func GenerateOTP() string {
	if os.Getenv("SERVER_MODE") != "PROD" {
		return "1111"
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	otp := rand.Intn(9999)           // Generate a random number between 0 and 999999
	return strconv.Itoa(otp)         // Convert to string
}

// SendEmail sends an email with the provided OTP to the specified email address
func SendEmail(to string, otp string) error {
	from := "your-email@example.com"  // Replace with your email
	password := "your-email-password" // Replace with your email password
	smtpHost := "smtp.example.com"    // Replace with your SMTP server
	smtpPort := "587"                 // SMTP port

	// Set up authentication information
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Create the email message
	message := []byte("Subject: Your OTP\n" +
		"\n" +
		"Your OTP is: " + otp + "\n")

	// Send the email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}
