package models

// ARTIST — matches /api/artists
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// LOCATIONS — matches /api/locations
type LocationItem struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Locations struct {
	Index []LocationItem `json:"index"`
}

// DATES — matches /api/dates
type DateItem struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Dates struct {
	Index []DateItem `json:"index"`
}

// RELATION — matches /api/relation
type RelationItem struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relation struct {
	Index []RelationItem `json:"index"`
}

// TEMPLATE DATA — bundles everything artist.html needs
type ArtistData struct {
	Artist   Artist
	Relation map[string][]string
}