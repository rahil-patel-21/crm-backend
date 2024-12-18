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
func InsertUser(user models.User) error {
	var userId int64

	// Construct the raw SQL query string with user inputs directly embedded (be cautious!)
	query := fmt.Sprintf(`
        INSERT INTO users (email, password, otp, is_verified)
        VALUES ('%s', '%s', '%s', %t)
        RETURNING id
    `, user.Email, user.Password, user.OTP, false)

	// Execute the SQL query directly using Exec
	err := db.QueryRow(context.Background(), query).Scan(&userId)
	if err != nil {
		log.Printf("Error inserting user with email %s into database: %v", user.Email, err)
		return err
	}

	return nil
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

func FindUserByEmail(email string) (models.User, error) {
	// Define the query to get the user by email
	query := `SELECT "id", "password" FROM users WHERE email = $1`

	// Create an empty User model to hold the results
	var user models.User

	// QueryRow is used for fetching a single row result
	err := db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}

func InsertTicket(ticket models.Ticket) error {
	var id int64

	query := fmt.Sprintf(`
    INSERT INTO tickets (
        customer_name, contact, email, address, city, district, state, pincode, GST, brand, 
		model_no, serial_no, issue_description, created_by, due_date, "status_id", "emp_id" ) 
    VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%d', '%d') 
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
		ticket.Status_Id,
		ticket.Emp_Id,
	)

	// Execute the SQL query directly using Exec
	err := db.QueryRow(context.Background(), query).Scan(&id)
	return err
}

func GetTicketList(page int, pageSize int, created_by int, status_id int) (int64, []models.Ticket, error) {
	var count int64
	var tickets []models.Ticket
	queryParams := []interface{}{}
	queryParams = append(queryParams, created_by)
	whereStr := `WHERE created_by = $1`

	// Query to count total tickets
	countQuery := `SELECT COUNT(*) FROM "tickets" WHERE created_by = $1`
	if status_id != 0 {
		countQuery += ` AND "status_id" = $2`
		queryParams = append(queryParams, status_id)
		whereStr += ` AND "status_id" = $2`
	}
	err := db.QueryRow(context.Background(), countQuery, queryParams...).Scan(&count)
	if err != nil {
		fmt.Println("Error counting tickets:", err)
		return 0, nil, err
	}
	if count == 0 {
		return count, tickets, nil // Avoid checking rows as count is already zero
	}

	// Query to select tickets with pagination
	query := `SELECT "tickets"."id", "customer_name", "contact", "email", "address", "city", "district",
	"state", "pincode", "gst", "brand", "model_no", "serial_no", "issue_description", "due_date",
	"created_at", "updated_at", "status_id", "TicketStatus"."name" AS "status", 
	CONCAT("Employee"."firstName", ' ', "Employee"."lastName") AS "assignee", "emp_id"
	FROM "tickets" 

	INNER JOIN "TicketStatus" ON "TicketStatus"."id" = "tickets"."status_id"
	INNER JOIN "Employee" ON "Employee"."id" = "tickets"."emp_id"`
	if status_id != 0 {
		query += whereStr
		query += ` LIMIT $3 OFFSET $4`
	} else {
		query += whereStr
		query += ` LIMIT $2 OFFSET $3`
	}

	queryParams = append(queryParams, pageSize)
	queryParams = append(queryParams, (page-1)*pageSize)
	rows, err := db.Query(context.Background(), query, queryParams...)
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
			&ticket.Created_At, &ticket.Updated_At, &ticket.Status_Id, &ticket.Status, &ticket.Assignee, &ticket.Emp_Id)
		if err != nil {
			fmt.Println("Error scanning ticket:", err)
			return 0, nil, err
		}
		tickets = append(tickets, ticket)
	}

	return count, tickets, nil
}

// Create a new customer into the database
func CreateCustomer(customer models.Customer) error {

	// Construct the raw SQL query string with user inputs directly embedded (be cautious!)
	query := fmt.Sprintf(`
        INSERT INTO "Customer" 
		("mobileNumber", "fullName", "email", "companyName", "address_1", "address_2", "city", "district", "state", "pincode",
		"type", "reference", "gstNumber")
        VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')
    `, customer.MobileNumber, customer.FullName, customer.Email, customer.CompanyName, customer.Address_1, customer.Address_2,
		customer.City, customer.District, customer.State, customer.Pincode, customer.Type, customer.Reference, customer.GstNumber)

	// Execute the SQL query directly using Exec
	_, err := db.Exec(context.Background(), query)
	return err
}

func GetCustomerList(page int, pageSize int, searchText string) (int64, []models.Customer, error) {
	var count int64
	var customers []models.Customer
	queryParams := []interface{}{}

	countQuery := `SELECT COUNT(*) FROM "Customer"`
	listQuery := `SELECT "mobileNumber", "fullName", "email", "companyName", "address_1", "address_2", "city", "district", "state", 
	"pincode", "type", "reference", "gstNumber"
	FROM "Customer"`
	// Search via number
	if searchText != "" {
		countQuery += ` WHERE "mobileNumber" ILIKE $1 OR "email" ILIKE $1`
		listQuery += ` WHERE "mobileNumber" ILIKE $1 OR "email" ILIKE $1`
		queryParams = append(queryParams, "%"+searchText+"%")
		listQuery += ` LIMIT $2 OFFSET $3`
	} else {
		listQuery += ` LIMIT $1 OFFSET $2`
	}
	err := db.QueryRow(context.Background(), countQuery, queryParams...).Scan(&count)
	if err != nil {
		fmt.Println("Error counting tickets:", err)
		return 0, nil, err
	}
	if count == 0 {
		return count, customers, nil // Avoid checking rows as count is already zero
	}

	// Query to select tickets with pagination
	queryParams = append(queryParams, pageSize, (page-1)*pageSize)
	rows, err := db.Query(context.Background(), listQuery, queryParams...)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}
	defer rows.Close()

	// Iterate over the result set and scan into tickets slice
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.MobileNumber, &customer.FullName, &customer.Email, &customer.CompanyName, &customer.Address_1,
			&customer.Address_2, &customer.City, &customer.District, &customer.State, &customer.Pincode, &customer.Type,
			&customer.Reference, &customer.GstNumber)
		if err != nil {
			return 0, nil, err
		}
		customers = append(customers, customer)
	}

	return count, customers, nil
}

// Create a new customer into the database
func CreateEmployee(employee models.Employee) error {

	// Construct the raw SQL query string with user inputs directly embedded (be cautious!)
	query := fmt.Sprintf(`
        INSERT INTO "Employee" 
		("mobileNumber", "firstName", "lastName", "org_id")
        VALUES ('%s', '%s', '%s', '%d')
    `, employee.MobileNumber, employee.FirstName, employee.LastName, employee.OrgId)

	// Execute the SQL query directly using Exec
	_, err := db.Exec(context.Background(), query)
	return err
}

func GetEmployeeList(page int, pageSize int, searchText string) (int64, []models.Employee, error) {
	var count int64
	var employees []models.Employee
	queryParams := []interface{}{}

	countQuery := `SELECT COUNT(*) FROM "Employee"`
	listQuery := `SELECT "mobileNumber", "firstName", "lastName", "id" 
	FROM "Employee"`
	// Search via number
	if searchText != "" {
		countQuery += ` WHERE "mobileNumber" ILIKE $1 OR "firstName" ILIKE $1 OR "lastName" ILIKE $1`
		listQuery += ` WHERE "mobileNumber" ILIKE $1 OR "firstName" ILIKE $1 OR "lastName" ILIKE $1`
		queryParams = append(queryParams, "%"+searchText+"%")
		listQuery += ` LIMIT $2 OFFSET $3`
	} else {
		listQuery += ` LIMIT $1 OFFSET $2`
	}
	err := db.QueryRow(context.Background(), countQuery, queryParams...).Scan(&count)
	if err != nil {
		fmt.Println("Error counting tickets:", err)
		return 0, nil, err
	}
	if count == 0 {
		return count, employees, nil // Avoid checking rows as count is already zero
	}

	// Query to select tickets with pagination
	queryParams = append(queryParams, pageSize, (page-1)*pageSize)
	rows, err := db.Query(context.Background(), listQuery, queryParams...)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}
	defer rows.Close()

	// Iterate over the result set and scan into tickets slice
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(&employee.MobileNumber, &employee.FirstName, &employee.LastName, &employee.ID)
		if err != nil {
			return 0, nil, err
		}
		employees = append(employees, employee)
	}

	return count, employees, nil
}
