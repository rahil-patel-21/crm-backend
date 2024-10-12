package db

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
	query := `
        INSERT INTO users (email, password, otp, is_verified)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := db.QueryRow(context.Background(), query, user.Email, user.Password, user.OTP, user.IsVerified).Scan(&userID)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return 0, err
	}

	return userID, nil
}

// FindUserByEmail finds a user by email
func FindUserByEmail(email string, user *models.User) error {
	query := `
        SELECT id, email, password, otp, is_verified
        FROM users
        WHERE email = $1
    `

	row := db.QueryRow(context.Background(), query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.OTP, &user.IsVerified)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return err
	}

	return nil
}
