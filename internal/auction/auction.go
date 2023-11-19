// internal/auction/auction.go
package auction

import (
	"database/sql"
	"fmt"
	"time"
)

// Auction represents an auction.
type Auction struct {
	AdID      int       `json:"ad_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	State     string    `json:"state"`
}

// AuctionRepository represents the repository for managing auctions.
type AuctionRepository struct {
	DB *sql.DB
}

// NewAuctionRepository creates a new AuctionRepository instance.
func NewAuctionRepository(db *sql.DB) *AuctionRepository {
	return &AuctionRepository{DB: db}
}

func (repo *AuctionRepository) StartAuction(a *Auction) error {
	query := "INSERT INTO auction (ad_id, start_time, end_time, state) VALUES (?, ?, ?, ?)"

	// Prepare the SQL statement
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the auction's data
	_, err = stmt.Exec(a.AdID, a.StartTime, a.EndTime, "active")
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}

func (repo *AuctionRepository) EndAuction(adID int) error {
	query := "UPDATE auction SET state = 'completed' WHERE ad_id = ?"

	// Prepare the SQL statement
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement to end the auction
	_, err = stmt.Exec(adID)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}

// GetActiveAuctions retrieves a list of active auctions from the database.
func (repo *AuctionRepository) GetActiveAuctions() ([]Auction, error) {
	query := "SELECT ad_id, start_time, end_time, state FROM auction WHERE state = 'active'"
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeAuctions []Auction
	for rows.Next() {
		var auction Auction
		err := rows.Scan(&auction.AdID, &auction.StartTime, &auction.EndTime, &auction.State)
		if err != nil {
			return nil, err
		}
		activeAuctions = append(activeAuctions, auction)
	}

	return activeAuctions, nil
}

// CloseAuction closes the auction for a specific ad.
func (repo *AuctionRepository) CloseAuction(adID int) error {
	query := "UPDATE auction SET state = 'completed' WHERE ad_id = ?"
	_, err := repo.DB.Exec(query, adID)
	return err
}
