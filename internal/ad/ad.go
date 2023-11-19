package ad

import (
	"database/sql"
	"fmt"
	"time"
)

type Ad struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	BasePrice int       `json:"base_price"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type AdRepository struct {
	DB *sql.DB
}

// NewAdRepository creates a new AdRepository instance.
func NewAdRepository(db *sql.DB) *AdRepository {
	return &AdRepository{DB: db}
}

func (repo *AdRepository) CreateAd(a *Ad) error {
	query := "INSERT INTO ad (text, base_price, start_time, end_time) VALUES (?, ?, ?, ?)"

	// Prepare the SQL statement
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the ad's data
	_, err = stmt.Exec(a.Text, a.BasePrice, a.StartTime, a.EndTime)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}

func (repo *AdRepository) GetAds() ([]Ad, error) {
	query := "SELECT * FROM ad"

	// Execute the SQL query to fetch all ads
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", err)
	}
	defer rows.Close()

	var ads []Ad

	// Iterate through the result set and scan each row into an Ad struct
	for rows.Next() {
		var a Ad
		err := rows.Scan(&a.ID, &a.Text, &a.BasePrice, &a.StartTime, &a.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		ads = append(ads, a)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return ads, nil
}
