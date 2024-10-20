// Packages
package db

// Imports
import (
	"context"
	"crm-backend/models"
	"fmt"
	"log"
	"os"
	"time"

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
        INSERT INTO users (email, password, otp,token, is_verified)
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

func UpdateUserVerifiedByID(userID int64, is_verified bool) error {
	query := fmt.Sprintf(`UPDATE users SET is_verified='%t' WHERE id = %d`, is_verified, userID)
	fmt.Print(query)
	// Execute the query
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Println("Error updating OTP for user:", err)
		return err
	}

	log.Println("OTP updated successfully for user ID:", userID)
	return nil
}

func FindOTPByEmail(email string) (int64, string, bool, error) {
	query := `SELECT "id", "otp","is_verified" FROM users WHERE email = $1`

	var userID int64
	var otp string
	var is_verified bool
	err := db.QueryRow(context.Background(), query, email).Scan(&userID, &otp, &is_verified) // Scan both values into variables
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", false, nil // User not found, return 0 and an empty OTP
		}
		log.Println("Error fetching user by email:", err)
		return 0, "", false, err // Return error with empty OTP
	}

	return userID, otp, is_verified, nil // Return both userID and otp
}

func InsertTicket(ticket models.Ticket) (int64, error) {
	var id int64

	query := fmt.Sprintf(`
    INSERT INTO tickets (
        customer_name, contact, email, address, city, district, state, pincode, GST, brand, 
		model_no, serial_no, issue_description, created_by, due_date ) 
    VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s') 
    RETURNING id`,
		ticket.Customer_Name,
		ticket.Contact_Number,
		ticket.Email,
		ticket.Address,
		ticket.City,
		ticket.District,
		ticket.State,
		ticket.Pincode,
		ticket.Gst_Number,
		ticket.Brand,
		ticket.Model_Number,
		ticket.Serial_Number,
		ticket.Issue_Description,
		ticket.Created_By,                    // '%d' for integer type
		ticket.Due_Date.Format(time.RFC3339), // '%s' for date in string format (ISO8601 format)
	)

	// Execute the SQL query directly using Exec
	err := db.QueryRow(context.Background(), query).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func GetTicketList(page int, pageSize int, created_by int) (int64, []models.Ticket, error) {
	var count int64
	var tickets []models.Ticket

	// Query to count total tickets
	countQuery := `SELECT COUNT(*) FROM "tickets" WHERE created_by = $1`
	err := db.QueryRow(context.Background(), countQuery, created_by).Scan(&count)
	if err != nil {
		fmt.Println("Error counting tickets:", err)
		return 0, nil, err
	}
	if count == 0 {
		return count, tickets, nil // Avoid checking rows as count is already zero
	}

	// Query to select tickets with pagination
	query := `SELECT "id", "customer_name", "contact", "email", "address", "city", "district",
	"state", "pincode", "gst", "brand", "model_no", "serial_no", "issue_description", "due_date",
	"created_at", "updated_at"
	FROM "tickets" 
	WHERE created_by = $1 LIMIT $2 OFFSET $3`
	rows, err := db.Query(context.Background(), query, created_by, pageSize, (page-1)*pageSize)
	if err != nil {
		fmt.Println("Error querying tickets:", err)
		return 0, nil, err
	}
	defer rows.Close()

	// Iterate over the result set and scan into tickets slice
	for rows.Next() {
		var ticket models.Ticket
		err := rows.Scan(&ticket.ID, &ticket.Customer_Name, &ticket.Contact_Number, &ticket.Email,
			&ticket.Address, &ticket.City, &ticket.District, &ticket.State, &ticket.Pincode, &ticket.Gst_Number,
			&ticket.Brand, &ticket.Model_Number, &ticket.Serial_Number, &ticket.Issue_Description, &ticket.Due_Date,
			&ticket.Created_At, &ticket.Updated_At)
		if err != nil {
			fmt.Println("Error scanning ticket:", err)
			return 0, nil, err
		}
		tickets = append(tickets, ticket)
	}

	return count, tickets, nil
}
