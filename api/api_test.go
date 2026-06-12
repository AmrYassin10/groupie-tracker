package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"groupie-tracker/models"
)

// ── Mock server helpers ───────────────────────────────────────────

// mockServer creates a fake HTTP server that responds with the given
// status code and body. Tests point the api functions at this server
// instead of the real API — so tests never need a network connection.
func mockServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)

		if _, err := w.Write([]byte(body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
}

// fakeArtistsJSON is valid JSON matching the real API format
const fakeArtistsJSON = `[
	{
		"id": 1,
		"name": "Queen",
		"image": "https://example.com/queen.jpg",
		"members": ["Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"],
		"creationDate": 1970,
		"firstAlbum": "1973-07-13",
		"locations": "https://example.com/locations/1",
		"concertDates": "https://example.com/dates/1",
		"relations": "https://example.com/relation/1"
	}
]`

// fakeRelationsJSON is valid JSON matching the real API format
const fakeRelationsJSON = `{
	"index": [
		{
			"id": 1,
			"datesLocations": {
				"london-uk": ["01-01-2019", "12-05-2020"],
				"paris-france": ["15-03-2019"]
			}
		}
	]
}`

// ── FetchArtists tests ────────────────────────────────────────────

// Happy path — server returns valid JSON
func TestFetchArtistsSuccess(t *testing.T) {
	server := mockServer(http.StatusOK, fakeArtistsJSON)
	defer server.Close()

	// Point the package's URL constant at our mock server
	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	artists, err := FetchArtists()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(artists) != 1 {
		t.Fatalf("expected 1 artist, got %d", len(artists))
	}
	if artists[0].Name != "Queen" {
		t.Errorf("expected name 'Queen', got '%s'", artists[0].Name)
	}
	if artists[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", artists[0].ID)
	}
	if len(artists[0].Members) != 4 {
		t.Errorf("expected 4 members, got %d", len(artists[0].Members))
	}
}

// Server returns non-200 status — should return an error
func TestFetchArtistsNon200Status(t *testing.T) {
	statuses := []int{
		http.StatusInternalServerError, // 500
		http.StatusNotFound,            // 404
		http.StatusServiceUnavailable,  // 503
	}

	for _, status := range statuses {
		server := mockServer(status, "error")
		originalURL := URL
		URL = server.URL + "/"

		_, err := FetchArtists()

		if err == nil {
			t.Errorf("status %d: expected an error, got nil", status)
		}

		URL = originalURL
		server.Close()
	}
}

// Server returns invalid JSON — should return a decode error
func TestFetchArtistsInvalidJSON(t *testing.T) {
	server := mockServer(http.StatusOK, `not valid json {{{`)
	defer server.Close()

	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	_, err := FetchArtists()

	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// Server returns empty array — should succeed with zero artists
func TestFetchArtistsEmptyArray(t *testing.T) {
	server := mockServer(http.StatusOK, `[]`)
	defer server.Close()

	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	artists, err := FetchArtists()

	if err != nil {
		t.Fatalf("expected no error for empty array, got: %v", err)
	}
	if len(artists) != 0 {
		t.Errorf("expected 0 artists, got %d", len(artists))
	}
}

// ── FetchRelations tests ──────────────────────────────────────────

// Happy path — server returns valid relations JSON
func TestFetchRelationsSuccess(t *testing.T) {
	server := mockServer(http.StatusOK, fakeRelationsJSON)
	defer server.Close()

	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	relations, err := FetchRelations()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(relations.Index) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(relations.Index))
	}

	rel := relations.Index[0]

	if rel.ID != 1 {
		t.Errorf("expected relation ID 1, got %d", rel.ID)
	}

	londonDates, ok := rel.DatesLocations["london-uk"]
	if !ok {
		t.Fatal("expected 'london-uk' in datesLocations")
	}
	if len(londonDates) != 2 {
		t.Errorf("expected 2 dates for london-uk, got %d", len(londonDates))
	}
}

// Server returns non-200 — should return an error
func TestFetchRelationsNon200Status(t *testing.T) {
	server := mockServer(http.StatusInternalServerError, "server error")
	defer server.Close()

	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	_, err := FetchRelations()

	if err == nil {
		t.Error("expected error for non-200 status, got nil")
	}
}

// Server returns invalid JSON — should return a decode error
func TestFetchRelationsInvalidJSON(t *testing.T) {
	server := mockServer(http.StatusOK, `{invalid}`)
	defer server.Close()

	originalURL := URL
	URL = server.URL + "/"
	defer func() { URL = originalURL }()

	_, err := FetchRelations()

	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// ── JSON decoding tests ───────────────────────────────────────────

// Artist struct correctly decodes all fields from JSON
func TestArtistStructDecoding(t *testing.T) {
	raw := `{
		"id": 5,
		"name": "Metallica",
		"image": "https://example.com/img.jpg",
		"members": ["James Hetfield", "Lars Ulrich"],
		"creationDate": 1981,
		"firstAlbum": "1983-07-25"
	}`

	var artist models.Artist
	if err := json.Unmarshal([]byte(raw), &artist); err != nil {
		t.Fatalf("failed to decode artist JSON: %v", err)
	}

	if artist.ID != 5 {
		t.Errorf("expected ID 5, got %d", artist.ID)
	}
	if artist.Name != "Metallica" {
		t.Errorf("expected 'Metallica', got '%s'", artist.Name)
	}
	if artist.CreationDate != 1981 {
		t.Errorf("expected 1981, got %d", artist.CreationDate)
	}
	if len(artist.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(artist.Members))
	}
	if artist.Members[0] != "James Hetfield" {
		t.Errorf("expected 'James Hetfield', got '%s'", artist.Members[0])
	}
}

// Relation struct correctly decodes the nested map structure
func TestRelationStructDecoding(t *testing.T) {
	raw := `{
		"id": 3,
		"datesLocations": {
			"new_york-usa": ["25-12-2019"],
			"tokyo-japan":  ["10-01-2020", "11-01-2020"]
		}
	}`

	var rel models.RelationItem
	if err := json.Unmarshal([]byte(raw), &rel); err != nil {
		t.Fatalf("failed to decode relation JSON: %v", err)
	}

	if rel.ID != 3 {
		t.Errorf("expected ID 3, got %d", rel.ID)
	}

	tokyo, ok := rel.DatesLocations["tokyo-japan"]
	if !ok {
		t.Fatal("expected 'tokyo-japan' key in datesLocations")
	}
	if len(tokyo) != 2 {
		t.Errorf("expected 2 dates for tokyo-japan, got %d", len(tokyo))
	}
}