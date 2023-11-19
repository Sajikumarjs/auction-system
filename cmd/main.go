package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/Sajikumarjs/auction-system/dbutil"
	"github.com/Sajikumarjs/auction-system/internal/ad"
	"github.com/Sajikumarjs/auction-system/internal/auction"
	"github.com/Sajikumarjs/auction-system/internal/bid"
	"github.com/Sajikumarjs/auction-system/web/handlers"
)

func main() {
	// Initialize database connection
	db, err := dbutil.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	adRepo := ad.NewAdRepository(db)
	bidRepo := bid.NewBidRepository(db)
	auctionRepo := auction.NewAuctionRepository(db)

	// Initialize router
	r := mux.NewRouter()

	// Ad endpoints
	r.HandleFunc("/add-ad", handlers.AddAdHandler(adRepo)).Methods("POST")
	r.HandleFunc("/list-ads", handlers.ListAdsHandler(adRepo)).Methods("GET")

	// Bid endpoints
	r.HandleFunc("/place-bid", handlers.PlaceBidHandler(bidRepo)).Methods("POST")
	r.HandleFunc("/list-bids", handlers.ListBidsHandler(bidRepo)).Methods("GET")

	// Auction endpoints
	r.HandleFunc("/run-auction", handlers.RunAuctionHandler(auctionRepo, bidRepo)).Methods("GET")
	r.HandleFunc("/start-auction", handlers.StartAuctionHandler(auctionRepo)).Methods("POST")

	port := 8080
	log.Printf("Auction service is running on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
