package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// User represents a user record from the database.
type User struct {
	ID   int
	Name string
	Age  int
}

// initDB establishes a connection to the PostgreSQL database.
func initDB() (*pgx.Conn, error) {
	// IMPORTANT: You must replace these values with your own PostgreSQL database credentials.
	// This connection string uses the "database/sql" format, which pgx can handle.
	connStr := "user=youruser password=yourpassword dbname=yourdbname host=localhost port=5432 sslmode=disable"
	
	// Establish a connection to the database.
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Ping the database to verify the connection is alive.
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return conn, nil
}

// getUsers retrieves a slice of users from the database.
func getUsers(ctx context.Context, db *pgx.Conn) ([]User, error) {
	// Execute the query to get all users.
	rows, err := db.Query(ctx, "SELECT id, name, age FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	// It's crucial to defer closing the rows to free up resources.
	defer rows.Close()

	var users []User
	// Iterate through the rows and scan the data into User structs.
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	// Check for any errors that may have occurred during row iteration.
	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return users, nil
}

func main() {
	// First, connect to the database.
	conn, err := initDB()
	if err != nil {
		log.Fatalf("could not initialize database connection: %v", err)
	}
	// Ensure the connection is closed when the main function exits.
	defer conn.Close(context.Background())

	// Use a context for the query.
	ctx := context.Background()

	// Example: Add a new user (assuming the 'users' table exists with 'id', 'name', 'age' columns).
	// This is just to ensure the table has some data for the select query to work.
	// You might need to create the table first in your PostgreSQL instance.
	// CREATE TABLE users (id serial PRIMARY KEY, name VARCHAR(50), age INT);
	insertSQL := `INSERT INTO users (name, age) VALUES ($1, $2)`
	_, err = conn.Exec(ctx, insertSQL, "Alice", 30)
	if err != nil {
		// Log the error but don't exit, as the user might already exist.
		log.Printf("could not insert user: %v", err)
	}

	// Now, retrieve all the users.
	users, err := getUsers(ctx, conn)
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}

	// Print the retrieved users.
	fmt.Println("\nUsers found in the database:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

// To run this code, you'll need a PostgreSQL database running and accessible.
// Make sure you have Go installed and the pgx driver.
// 1. Initialize a Go module: go mod init mydriver
// 2. Install the pgx driver: go get github.com/jackc/pgx/v5
// 3. Update the connection string with your credentials.
// 4. Run the code: go run .
