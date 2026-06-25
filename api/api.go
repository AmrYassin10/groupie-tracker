package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"groupie-tracker/models"
)

// URL is var not const so tests can point it at a mock server
var URL = "https://groupietrackers.herokuapp.com/api/"

// client has 10s timeout — prevents hanging if API is slow
// transport layer protection ✅
var client = &http.Client{Timeout: 10 * time.Second}

// fetchAndDecode is a shared helper — avoids repeating the same
// http get + status check + json decode pattern in every function
func fetchAndDecode(endpoint string, target interface{}) error {

	// transport layer — did the request reach the server?
	resp, err := client.Get(URL + endpoint)
	if err != nil {
		return fmt.Errorf("can't reach %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// application layer — did the server respond correctly?
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned status: %d", endpoint, resp.StatusCode)
	}

	// decode JSON into the target struct
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding %s: %v", endpoint, err)
	}

	return nil
}

func FetchArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if err := fetchAndDecode("artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func FetchRelations() (models.Relation, error) {
	var relations models.Relation
	if err := fetchAndDecode("relation", &relations); err != nil {
		return models.Relation{}, err
	}
	return relations, nil
}

func FetchLocations() (models.Locations, error) {
	var locations models.Locations
	if err := fetchAndDecode("locations", &locations); err != nil {
		return models.Locations{}, err
	}
	return locations, nil
}

func FetchDates() (models.Dates, error) {
	var dates models.Dates
	if err := fetchAndDecode("dates", &dates); err != nil {
		return models.Dates{}, err
	}
	return dates, nil
}