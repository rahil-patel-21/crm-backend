// Packages
package db

// Imports
import (
	"context"
	"crm-backend/models"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("Unable to connect to database")
	}

	db, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
}

// CloseDB closes the PostgreSQL database connection
func CloseDB() {
	if db != nil {
		db.Close(context.Background())
	}
}

// InsertUser inserts a new user into the database
func InsertUser(user models.User) (int64, error) {
	var userID int64

	// Construct the raw SQL query string with user inputs directly embedded (be cautious!)
	query := fmt.Sprintf(`
        INSERT INTO users (email, password, otp, is_verified)
        VALUES ('%s', '%s', '%s', %t)
        RETURNING id
    `, user.Email, user.Password, user.OTP, false)

	// Execute the SQL query directly using Exec
	err := db.QueryRow(context.Background(), query).Scan(&userID)
	if err != nil {
		log.Printf("Error inserting user with email %s into database: %v", user.Email, err)
		return 0, err
	}

	return userID, nil
}

// FindUserByEmail finds a user by email
func FindUserIdByEmail(email string) (int64, error) {
	query := `SELECT id FROM users WHERE email = $1`

	var userID int64
	err := db.QueryRow(context.Background(), query, email).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil // User not found, return 0
		}
		log.Println("Error fetching user by email:", err)
		return 0, err
	}

	return userID, nil
}

func UpdateUserOTPByEmailId(userID int64, otp string) error {
	// Raw SQL query to update the otp field for a specific user by ID
	query := fmt.Sprintf(`UPDATE users SET otp = '%s' WHERE id = %d`, otp, userID)

	// Execute the query
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Println("Error updating OTP for user:", err)
		return err
	}

	log.Println("OTP updated successfully for user ID:", userID)
	return nil
}

func FindOTPByEmail(email string) (int64, string, error) {
	query := `SELECT "id", "otp" FROM users WHERE email = $1`

	var userID int64
	var otp string
	err := db.QueryRow(context.Background(), query, email).Scan(&userID, &otp) // Scan both values into variables
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", nil // User not found, return 0 and an empty OTP
		}
		log.Println("Error fetching user by email:", err)
		return 0, "", err // Return error with empty OTP
	}

	return userID, otp, nil // Return both userID and otp
}
