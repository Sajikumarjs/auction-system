package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sajikumarjs/auction-system/internal/bid"
	"github.com/gorilla/mux"
)

// PlaceBidHandler handles the /place-bid endpoint.
func PlaceBidHandler(bidRepo *bid.BidRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request parameters
		var request struct {
			BidderID int `json:"bidder_id"`
			AdID     int `json:"ad_id"`
			Price    int `json:"price"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate request parameters (add your validation logic here)
		if request.BidderID <= 0 || request.AdID <= 0 || request.Price <= 0 {
			http.Error(w, "Invalid bid parameters", http.StatusBadRequest)
			return
		}

		// Create a new Bid instance
		newBid := bid.Bid{
			BidderID: request.BidderID,
			AdID:     request.AdID,
			Price:    request.Price,
		}

		// Place the bid in the database
		err = bidRepo.PlaceBid(&newBid)
		if err != nil {
			http.Error(w, "Failed to place the bid", http.StatusInternalServerError)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusCreated)
	}
}

// ListBidsHandler handles the /list-bids endpoint.
func ListBidsHandler(bidRepo *bid.BidRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ad ID from the URL parameters
		adIDStr := mux.Vars(r)["adID"]
		adID, err := strconv.Atoi(adIDStr)
		if err != nil {
			http.Error(w, "Invalid ad ID", http.StatusBadRequest)
			return
		}

		// Retrieve bids for the specified ad from the database
		bids, err := bidRepo.GetBidsForAd(adID)
		if err != nil {
			http.Error(w, "Failed to retrieve bids for ad", http.StatusInternalServerError)
			return
		}

		// Return the list of bids as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bids)
	}
}
