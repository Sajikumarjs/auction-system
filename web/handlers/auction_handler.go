// web/handlers/auction_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sajikumarjs/auction-system/internal/auction"
	"github.com/Sajikumarjs/auction-system/internal/bid"
)

// RunAuctionHandler handles the /run-auction endpoint.
func RunAuctionHandler(auctionRepo *auction.AuctionRepository, bidRepo *bid.BidRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve active auctions from the database
		activeAuctions, err := auctionRepo.GetActiveAuctions()
		if err != nil {
			http.Error(w, "Failed to retrieve active auctions", http.StatusInternalServerError)
			return
		}

		// Create a map to store winning bids for each ad
		winningBids := make(map[int]bid.Bid)

		// Iterate through active auctions and determine winners
		for _, activeAuction := range activeAuctions {
			// Retrieve bids for the current ad
			adID := activeAuction.AdID
			bids, err := bidRepo.GetBidsForAd(adID)
			if err != nil {
				http.Error(w, "Failed to retrieve bids for ad", http.StatusInternalServerError)
				return
			}

			// Find the highest bid for the current ad
			var highestBid bid.Bid
			for _, currentBid := range bids {
				if currentBid.Price > highestBid.Price {
					highestBid = currentBid
				}
			}

			// Store the winning bid for the current ad
			winningBids[adID] = highestBid
		}

		// Close the auctions
		for _, activeAuction := range activeAuctions {
			err := auctionRepo.CloseAuction(activeAuction.AdID)
			if err != nil {
				http.Error(w, "Failed to close auction", http.StatusInternalServerError)
				return
			}
		}

		// Return the winning bids as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(winningBids)
	}
}

// StartAuctionHandler handles the /start-auction endpoint.
func StartAuctionHandler(auctionRepo *auction.AuctionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request parameters
		var request struct {
			AdID int `json:"ad_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Check if the auction for the given ad is already active
		activeAuctions, err := auctionRepo.GetActiveAuctions()
		if err != nil {
			http.Error(w, "Failed to check for active auctions", http.StatusInternalServerError)
			return
		}

		for _, activeAuction := range activeAuctions {
			if activeAuction.AdID == request.AdID {
				http.Error(w, "Auction for this ad is already active", http.StatusBadRequest)
				return
			}
		}

		// Set the start and end times for the auction
		startTime := time.Now()
		endTime := startTime.Add(24 * time.Hour) // Auction lasts for 24 hours (adjust as needed)

		// Create a new auction
		newAuction := auction.Auction{
			AdID:      request.AdID,
			StartTime: startTime,
			EndTime:   endTime,
			State:     "active",
		}

		// Start the auction in the database
		err = auctionRepo.StartAuction(&newAuction)
		if err != nil {
			http.Error(w, "Failed to start the auction", http.StatusInternalServerError)
			return
		}

		// Return the auction details as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newAuction)
	}
}
