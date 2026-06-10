package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"groupie-tracker/models"
)

// URL is a var (not const) so tests can point it at a mock server
var URL = "https://groupietrackers.herokuapp.com/api/"

// client has a 10s timeout — prevents server hanging if API is slow
var client = &http.Client{Timeout: 10 * time.Second}

func FetchArtists() ([]models.Artist, error) {

	resp, err := client.Get(URL + "artists")
	if err != nil {
		return nil, fmt.Errorf("can't get artists: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var artists []models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("error decoding artists JSON: %v", err)
	}

	return artists, nil
}

func FetchRelations() (models.Relation, error) {

	resp, err := client.Get(URL + "relation")
	if err != nil {
		return models.Relation{}, fmt.Errorf("can't get relations: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Relation{}, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var relations models.Relation
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return models.Relation{}, fmt.Errorf("error decoding relations JSON: %v", err)
	}

	return relations, nil
}