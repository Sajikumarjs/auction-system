package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sajikumarjs/auction-system/dbutil"
	"github.com/Sajikumarjs/auction-system/internal/ad"
)

func AddAdHandler(adRepo *ad.AdRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newAd ad.Ad
		err := json.NewDecoder(r.Body).Decode(&newAd)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		db, err := dbutil.OpenDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Initialize AdRepository with a database connection
		repo := &ad.AdRepository{DB: db}

		// Add ad to the database
		err = repo.CreateAd(&newAd)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to add ad", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
func ListAdsHandler(adRepo *ad.AdRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize AdRepository with a database connection
		db, err := dbutil.OpenDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		repo := &ad.AdRepository{DB: db}

		// Get all ads from the database
		ads, err := repo.GetAds()
		if err != nil {
			http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
			return
		}

		// Return ads as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ads)
	}
}
