// dbutil/db.go
package dbutil

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// OpenDB opens a new database connection.
func OpenDB() (*sql.DB, error) {
	// Replace the connection string with your actual database connection details
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/auction_db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}
