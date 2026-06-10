package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"groupie-tracker/models"
)

// mu protects all three data structures from race conditions.
var mu sync.RWMutex

var Artists      []models.Artist
var ArtistsMap   = make(map[int]models.Artist)
var RelationsMap = make(map[int]map[string][]string)

// tmpl starts nil — initialized explicitly by InitTemplates.
// This allows tests to call InitTemplates with the correct path
// before any handler runs, instead of failing at package init time.
var tmpl *template.Template

// InitTemplates parses both HTML templates from the given directory.
// Must be called once before any handler is used — in main() and in TestMain.
func InitTemplates(dir string) {
	tmpl = template.Must(template.ParseFiles(
		dir+"/index.html",
		dir+"/artist.html",
	))
}

// SetData replaces all three data structures atomically.
func SetData(artists []models.Artist, artistsMap map[int]models.Artist, relations map[int]map[string][]string) {
	mu.Lock()
	defer mu.Unlock()

	Artists      = artists
	ArtistsMap   = artistsMap
	RelationsMap = relations
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl == nil {
        http.Error(w, "Server not initialised", http.StatusInternalServerError)
        return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	mu.RLock()
	data := Artists
	mu.RUnlock()

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error rendering index: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.RLock()
	artist, ok := ArtistsMap[id]
	relation   := RelationsMap[id]
	mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	data := models.ArtistData{
		Artist:   artist,
		Relation: relation,
	}

	if err := tmpl.ExecuteTemplate(w, "artist.html", data); err != nil {
		log.Printf("Error rendering artist: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}