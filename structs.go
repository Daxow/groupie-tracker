package main

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location	 []string `json:"-"`
}

type Location struct {
	ID		 int      `json:"id"`
	Location []string `json:"location"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	ID        int              `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
