package main

type Artist struct {
	ID           int          `json:"id"`
	Image        string       `json:"image"`
	Name         string       `json:"name"`
	Members      []string     `json:"members"`
	CreationDate int          `json:"creationDate"`
	FirstAlbum   string       `json:"firstAlbum"`
	Locations    []string     `json:"-"`
	Coordinates  []Coordinate `json:"-"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationResponse struct {
	Index []Location `json:"index"`
}

type DateResponse struct {
	Index []Date `json:"index"`
}

type RelationResponse struct {
	Index []Relation `json:"index"`
}

type Coordinate struct {
    Latitude  string `json:"lat"`
    Longitude string `json:"lon"`
}