package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"groupie-tracker/models"
)

// TestMain runs before all tests in this package.
// It moves up to the project root so InitTemplates can find templates/,
// then seeds the global data with fake artists so handler tests work
// without hitting the real API.
func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}


    // test code

	// Initialize templates from the project root
	InitTemplates("templates")

	// Seed fake data — tests never call the real API
	fakeArtists := []models.Artist{
		{
			ID:           1,
			Name:         "Queen",
			Image:        "https://example.com/queen.jpg",
			Members:      []string{"Freddie Mercury", "Brian May"},
			CreationDate: 1970,
			FirstAlbum:   "1973-07-13",
		},
		{
			ID:           2,
			Name:         "Pantera",
			Image:        "https://example.com/pantera.jpg",
			Members:      []string{"Phil Anselmo", "Dimebag Darrell"},
			CreationDate: 1981,
			FirstAlbum:   "1983-01-01",
		},
	}

	fakeArtistsMap := map[int]models.Artist{
		1: fakeArtists[0],
		2: fakeArtists[1],
	}

	fakeRelationsMap := map[int]map[string][]string{
		1: {
			"london-uk":      {"01-01-2019", "12-05-2020"},
			"paris-france":   {"15-03-2019"},
		},
		// artist 2 has no relations — tests the {{else}} branch in template
	}

	SetData(fakeArtists, fakeArtistsMap, fakeRelationsMap)

	// Run all tests
	os.Exit(m.Run())
}

// ── HomeHandler tests ─────────────────────────────────────────────

// GET / → 200 and contains artist names
func TestHomeHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	HomeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	// Template should have rendered both artist names
	if !strings.Contains(body, "Queen") {
		t.Error("expected response to contain 'Queen'")
	}
	if !strings.Contains(body, "Pantera") {
		t.Error("expected response to contain 'Pantera'")
	}
}
func BenchmarkHomeHandler(b *testing.B) {
    req := httptest.NewRequest("GET", "/", nil)

    for i := 0; i < b.N; i++ {
        rr := httptest.NewRecorder()
        HomeHandler(rr, req)
    }
}
// POST / → 405 Method Not Allowed
func TestHomeHandlerRejectsNonGET(t *testing.T) {

	// Test multiple non-GET methods — all should be rejected
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/", nil)
		rr := httptest.NewRecorder()

		HomeHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, rr.Code)
		}
	}
}

// GET /unknown → 404
func TestHomeHandlerRejectsUnknownPath(t *testing.T) {
	paths := []string{"/unknown", "/about", "/foo/bar", "/artist/1"}

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		HomeHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("path %s: expected 404, got %d", path, rr.Code)
		}
	}
}

// ── ArtistHandler tests ───────────────────────────────────────────

// GET /artist?id=1 → 200 and shows correct artist
func TestArtistHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist?id=1", nil)
	rr := httptest.NewRecorder()

	ArtistHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	if !strings.Contains(body, "Queen") {
		t.Error("expected response to contain 'Queen'")
	}
	if !strings.Contains(body, "Freddie Mercury") {
		t.Error("expected response to contain member 'Freddie Mercury'")
	}
	if !strings.Contains(body, "london-uk") {
		t.Error("expected response to contain concert location 'london-uk'")
	}
}

// GET /artist?id=2 → 200, artist with no relations shows fallback message
func TestArtistHandlerNoRelations(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist?id=2", nil)
	rr := httptest.NewRecorder()

	ArtistHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	if !strings.Contains(body, "No concert data available") {
		t.Error("expected 'No concert data available' for artist with no relations")
	}
}

// POST /artist?id=1 → 405
func TestArtistHandlerRejectsNonGET(t *testing.T) {
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/artist?id=1", nil)
		rr := httptest.NewRecorder()

		ArtistHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, rr.Code)
		}
	}
}

// GET /artist?id=abc → 400 Bad Request (non-numeric ID)
func TestArtistHandlerInvalidID(t *testing.T) {
	invalidIDs := []string{"abc", "1.5", "!!", "null", "xyz"}

	for _, id := range invalidIDs {
		req := httptest.NewRequest(http.MethodGet, "/artist?id="+id, nil)
		rr := httptest.NewRecorder()

		ArtistHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("id=%q: expected 400, got %d", id, rr.Code)
		}
	}
}

// GET /artist → 400 (missing id param — Atoi("") fails)
func TestArtistHandlerMissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist", nil)
	rr := httptest.NewRecorder()

	ArtistHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing id, got %d", rr.Code)
	}
}

// GET /artist?id=999 → 404 (ID doesn't exist in data)
func TestArtistHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/artist?id=999", nil)
	rr := httptest.NewRecorder()

	ArtistHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown artist, got %d", rr.Code)
	}
}

// ── SetData tests ─────────────────────────────────────────────────

// SetData correctly replaces all three data structures
func TestSetDataUpdatesAllStructures(t *testing.T) {
	newArtists := []models.Artist{
		{ID: 99, Name: "Test Band"},
	}
	newArtistsMap := map[int]models.Artist{
		99: {ID: 99, Name: "Test Band"},
	}
	newRelationsMap := map[int]map[string][]string{
		99: {"berlin-germany": {"01-01-2025"}},
	}

	SetData(newArtists, newArtistsMap, newRelationsMap)

	mu.RLock()
	defer mu.RUnlock()

	if len(Artists) != 1 || Artists[0].Name != "Test Band" {
		t.Errorf("Artists not updated correctly, got: %v", Artists)
	}

	artist, ok := ArtistsMap[99]
	if !ok || artist.Name != "Test Band" {
		t.Error("ArtistsMap not updated correctly")
	}

	rel, ok := RelationsMap[99]
	if !ok || len(rel["berlin-germany"]) != 1 {
		t.Error("RelationsMap not updated correctly")
	}
}