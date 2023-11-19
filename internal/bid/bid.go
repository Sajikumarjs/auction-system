// internal/bid/bid.go
package bid

import (
	"database/sql"
	"fmt"
)

// Bid represents a bid.
type Bid struct {
	BidderID int `json:"bidder_id"`
	AdID     int `json:"ad_id"`
	Price    int `json:"price"`
}

// BidRepository represents the repository for managing bids.
type BidRepository struct {
	DB *sql.DB
}

// NewBidRepository creates a new BidRepository instance.
func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{DB: db}
}

func (repo *BidRepository) PlaceBid(b *Bid) error {
	query := "INSERT INTO bid (bidder_id, ad_id, price) VALUES (?, ?, ?)"

	// Prepare the SQL statement
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the bid's data
	_, err = stmt.Exec(b.BidderID, b.AdID, b.Price)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}

// GetBidsForAd retrieves a list of bids for a specific ad from the database.
func (repo *BidRepository) GetBidsForAd(adID int) ([]Bid, error) {
	query := "SELECT bidder_id, ad_id, price FROM bid WHERE ad_id = ?"
	rows, err := repo.DB.Query(query, adID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []Bid
	for rows.Next() {
		var currentBid Bid
		err := rows.Scan(&currentBid.BidderID, &currentBid.AdID, &currentBid.Price)
		if err != nil {
			return nil, err
		}
		bids = append(bids, currentBid)
	}

	return bids, nil
}
