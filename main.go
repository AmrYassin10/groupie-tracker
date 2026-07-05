package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"groupie-tracker/api"
	"groupie-tracker/handlers"
	"groupie-tracker/models"
)

func main() {

	// parse templates once at startup ✅
	handlers.InitTemplates("templates")

	// retry with exponential backoff + jitter ✅
	// solves thundering herd problem
	if err := fetchWithBackoff(5); err != nil {
		log.Fatal("All retries failed: ", err)
	}

	

	// register routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	// serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}

// fetchWithBackoff retries loadAndSetData with exponential backoff + jitter
// exponential: wait doubles each attempt 2s → 4s → 8s → 16s
// jitter: random extra wait so 1000 servers don't retry at same time
// solves thundering herd problem ✅
func fetchWithBackoff(maxRetries int) error {
	for attempt := 0; attempt < maxRetries; attempt++ {

		if attempt > 0 {
			// double the wait each attempt
			base := time.Duration(1<<uint(attempt)) * time.Second

			// random jitter between 0 and base
			jitter := time.Duration(rand.Int63n(int64(base)))

			wait := base + jitter
			log.Printf("Attempt %d failed. Waiting %v...", attempt, wait)
			time.Sleep(wait)
		}

		log.Printf("Attempt %d of %d...", attempt+1, maxRetries)

		if err := loadAndSetData(); err != nil {
			log.Printf("Attempt %d failed: %v", attempt+1, err)
			continue
		}

		log.Printf("Success on attempt %d ✅", attempt+1)
		return nil
	}

	return fmt.Errorf("all %d attempts failed", maxRetries)
}

// loadAndSetData fetches all 4 API endpoints concurrently
// uses goroutines + WaitGroup — cuts fetch time in half ✅
func loadAndSetData() error {

	// local variables — stored on STACK ✅
	var (
		artists      []models.Artist
		relationsResp models.Relation
		locationsResp models.Locations
		datesResp     models.Dates
	)

	var (
		errArtists   error
		errRelations error
		errLocations error
		errDates     error
	)

	// fetch all 4 endpoints concurrently ✅
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		artists, errArtists = api.FetchArtists()
	}()

	go func() {
		defer wg.Done()
		relationsResp, errRelations = api.FetchRelations()
	}()

	go func() {
		defer wg.Done()
		locationsResp, errLocations = api.FetchLocations()
	}()

	go func() {
		defer wg.Done()
		datesResp, errDates = api.FetchDates()
	}()

	// wait for ALL goroutines to finish
	wg.Wait()

	// check all errors — all must succeed
	if errArtists != nil {
		return errArtists
	}
	if errRelations != nil {
		return errRelations
	}
	if errLocations != nil {
		return errLocations
	}
	if errDates != nil {
		return errDates
	}

	// fix broken image for artist ID 21
	for i := range artists {
		if artists[i].ID == 21 {
			artists[i].Image = "/static/mamonasassassinas.jpeg"
			break
		}
	}

	// build maps with exact capacity — avoids rehashing ✅
	// O(1) lookup instead of O(n) loop ✅
	artistsMap := make(map[int]models.Artist, len(artists))
	for _, a := range artists {
		artistsMap[a.ID] = a
	}

	relationsMap := make(map[int]map[string][]string, len(relationsResp.Index))
	for _, rel := range relationsResp.Index {
		relationsMap[rel.ID] = rel.DatesLocations
	}

	// locations and dates stored separately if needed later
	_ = locationsResp
	_ = datesResp

	// atomic update — all three updated under ONE lock ✅
	handlers.SetData(artists, artistsMap, relationsMap)

	return nil
}