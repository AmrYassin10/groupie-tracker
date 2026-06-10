package main

import (
	"log"
	"net/http"
	"time"

	"groupie-tracker/api"
	"groupie-tracker/handlers"
	"groupie-tracker/models"
)

func main() {

	// Initialize templates before anything else
	handlers.InitTemplates("templates")

	// 1. Load data once at startup
	if err := loadAndSetData(); err != nil {
		log.Fatal("Initial data load failed: ", err)
	}

	// 2. Background refresh every 2 hours
	go func() {
		for {
			time.Sleep(2 * time.Hour)
			log.Println("Refreshing data from API...")
			if err := loadAndSetData(); err != nil {
				log.Printf("Refresh failed, keeping old data: %v", err)
			} else {
				log.Println("Data refreshed successfully")
			}
		}
	}()

	// 3. Register routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}

func loadAndSetData() error {

	artists, err := api.FetchArtists()
	if err != nil {
		return err
	}

	for i := range artists {
		if artists[i].ID == 21 {
			artists[i].Image = "/static/mamonasassassinas.jpeg"
			break
		}
	}

	artistsMap := make(map[int]models.Artist, len(artists))
	for _, artist := range artists {
		artistsMap[artist.ID] = artist
	}

	relationsResp, err := api.FetchRelations()
	if err != nil {
		return err
	}

	relationsMap := make(map[int]map[string][]string, len(relationsResp.Index))
	for _, rel := range relationsResp.Index {
		relationsMap[rel.ID] = rel.DatesLocations
	}

	handlers.SetData(artists, artistsMap, relationsMap)
	return nil
}