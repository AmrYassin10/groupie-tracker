package models

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

// RELATION (INNER ITEM)
type RelationItem struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// RELATION RESPONSE (WRAPPER)
type Relation struct {
	Index []RelationItem `json:"index"`
}

type ArtistData struct {
	Artist   Artist
	Relation map[string][]string
}
