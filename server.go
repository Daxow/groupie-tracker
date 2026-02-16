package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"sort"
	"time"
)

type Server struct {
	Artists   []Artist
	Locations []Location
	Dates     []Date
	Relations []Relation
	Templates *template.Template
}

type PageData struct {
	Query   string
	Artists []Artist
	Locations []string
}

type ArtistPageData struct {
    Artist Artist
}

func NewServer() *Server {
	artists := fetchArtists()
	locations := fetchLocations()
	dates := fetchDates()
	relations := fetchRelations()

	artists = linkArtistsWithLocations(artists, relations)
	artists = preloadCoordinates(artists)


	funcMap := template.FuncMap{
		"json": func(data any) template.JS {
			jsonEncodedData, _ := json.Marshal(data)
			return template.JS(jsonEncodedData)
		},
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	return &Server{
		Artists:   artists,
		Locations: locations,
		Dates:     dates,
		Relations: relations,
		Templates: templates,
	}
}

func (server *Server) RegisterRoutes() {
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/search", server.handleSearch)
	http.HandleFunc("/api/artists", server.handleArtistsAPI)
	http.HandleFunc("/filter", server.handleFilter)
	http.HandleFunc("/artist", server.handleArtistPage)
	http.HandleFunc("/api/artist-locations", server.handleArtistLocationsAPI)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (server *Server) handleHome(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}

	data := PageData{
		Query:   "",
		Artists: server.Artists,
		Locations: GetUniqueLocations(server.Artists),
	}

	server.renderPage(response, "index.html", data)
}

func (server *Server) renderPage(response http.ResponseWriter, name string, data any) {
	err := server.Templates.ExecuteTemplate(response, name, data)
	if err != nil {
		http.Error(response, "Erreur serveur", http.StatusInternalServerError)
		log.Println("Erreur template :", err)
	}
}

func SearchArtists(artists []Artist, search string) []Artist {
	var result []Artist

	search = strings.ToLower(search)

	for _, artist := range artists {
		found := false

		if strings.Contains(strings.ToLower(artist.Name), search) {
			found = true
		}

		if !found {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), search) {
					found = true
					break
				}
			}
		}

		if !found {
			creation := fmt.Sprintf("%d", artist.CreationDate)
			if strings.Contains(creation, search) {
				found = true
			}
		}

		if !found {
			if strings.Contains(strings.ToLower(artist.FirstAlbum), search) {
				found = true
			}
		}

		if !found {
			for _, location := range artist.Locations {
				if strings.Contains(strings.ToLower(location), search) {
					found = true
					break
				}
			}
		}

		if found {
			result = append(result, artist)
		}
	}

	return result
}

func (server *Server) handleSearch(response http.ResponseWriter, request *http.Request) {
	query := strings.TrimSpace(request.URL.Query().Get("query"))

	artists := server.Artists
	if query != "" {
		artists = SearchArtists(server.Artists, query)
	}

	data := PageData{
		Query:   query,
		Artists: artists,
		Locations: GetUniqueLocations(server.Artists),
	}

	server.renderPage(response, "index.html", data)
}

func linkArtistsWithLocations(artists []Artist, relations []Relation) []Artist {
	for i, artist := range artists {
		for _, relation := range relations {
			if artist.ID == relation.ID {
				var locations []string

				for location := range relation.DatesLocations {
					locations = append(locations, location)
				}

				artists[i].Locations = locations
				break
			}
		}
	}
	return artists
}

func (server *Server) handleArtistsAPI(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var artistsWithLocations []map[string]any

	for _, artist := range server.Artists {

		artistMap := map[string]any{
			"id":           artist.ID,
			"image":        artist.Image,
			"name":         artist.Name,
			"members":      artist.Members,
			"creationDate": artist.CreationDate,
			"firstAlbum":   artist.FirstAlbum,
			"locations":    artist.Locations,
		}

		artistsWithLocations = append(artistsWithLocations, artistMap)
	}

	json.NewEncoder(response).Encode(artistsWithLocations)
}

func FilterArtists(artists []Artist, creationMin, creationMax int, albumMin, albumMax int, membersMin, membersMax int, selectedLocations []string) []Artist {
    var result []Artist

    for _, artist := range artists {

        if creationMin != 0 && artist.CreationDate < creationMin {
            continue
        }
        if creationMax != 0 && artist.CreationDate > creationMax {
            continue
        }

        albumYear := 0
        if len(artist.FirstAlbum) >= 4 {
    		yearString := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
    		fmt.Sscanf(yearString, "%d", &albumYear)
        }

        if albumMin != 0 && albumYear < albumMin {
            continue
        }
        if albumMax != 0 && albumYear > albumMax {
            continue
        }

        membersCount := len(artist.Members)

        if membersMin != 0 && membersCount < membersMin {
            continue
        }
        if membersMax != 0 && membersCount > membersMax {
            continue
        }

        if len(selectedLocations) > 0 {
            match := false
            for _, location := range artist.Locations {
                for _, selected := range selectedLocations {
                    if strings.Contains(strings.ToLower(location), strings.ToLower(selected)) {
                        match = true
                        break
                    }
                }
            }
            if !match {
                continue
            }
        }

        result = append(result, artist)
    }

    return result
}

func (server *Server) handleFilter(response http.ResponseWriter, request *http.Request) {
    request.ParseForm()

    creationMin, _ := strconv.Atoi(request.FormValue("creationMin"))
    creationMax, _ := strconv.Atoi(request.FormValue("creationMax"))

    albumMin, _ := strconv.Atoi(request.FormValue("albumMin"))
    albumMax, _ := strconv.Atoi(request.FormValue("albumMax"))

    membersMin, _ := strconv.Atoi(request.FormValue("membersMin"))
    membersMax, _ := strconv.Atoi(request.FormValue("membersMax"))

    selectedLocations := request.Form["locations"]

    filtered := FilterArtists(
        server.Artists,
        creationMin, creationMax,
        albumMin, albumMax,
        membersMin, membersMax,
        selectedLocations,
    )

    data := PageData{
        Query:   "",
        Artists: filtered,
		Locations: GetUniqueLocations(server.Artists),
    }

    server.renderPage(response, "index.html", data)
}

func GetUniqueLocations(artists []Artist) []string {
    locationMap := make(map[string]bool)

    for _, artist := range artists {
        for _, location := range artist.Locations {
            locationMap[location] = true
        }
    }

    var uniqueLocations []string
    for location := range locationMap {
        uniqueLocations = append(uniqueLocations, location)
    }

	sort.Strings(uniqueLocations)
    return uniqueLocations
}

func (server *Server) handleArtistPage(response http.ResponseWriter, request *http.Request) {

    idString := request.URL.Query().Get("id")
    id, err := strconv.Atoi(idString)
    if err != nil {
        http.NotFound(response, request)
        return
    }

    var selected Artist
    for _, artist := range server.Artists {
        if artist.ID == id {
            selected = artist
            break
        }
    }

    if selected.ID == 0 {
        http.NotFound(response, request)
        return
    }

    data := ArtistPageData{
        Artist: selected,
    }

    server.renderPage(response, "artist.html", data)
}

func (server *Server) handleArtistLocationsAPI(response http.ResponseWriter, request *http.Request) {

    response.Header().Set("Content-Type", "application/json")

    idStr := request.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(response, "ID invalide", http.StatusBadRequest)
        return
    }

    var selected Artist
    for _, artist := range server.Artists {
        if artist.ID == id {
            selected = artist
            break
        }
    }

    if selected.ID == 0 {
        http.Error(response, "Artiste introuvable", http.StatusNotFound)
        return
    }

    json.NewEncoder(response).Encode(selected.Coordinates)
}

func getCoordinatesFromLocation(location string) (string, string, error) {

    formatted := strings.ReplaceAll(location, " ", "+")
    url := "https://nominatim.openstreetmap.org/search?q=" + formatted + "&format=json&limit=1"

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", "", err
    }

    request.Header.Set("User-Agent", "groupie-tracker")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return "", "", err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return "", "", fmt.Errorf("Erreur API")
    }

    var result []struct {
        Latitude string `json:"lat"`
        Longitude string `json:"lon"`
    }

    err = json.NewDecoder(response.Body).Decode(&result)
    if err != nil || len(result) == 0 {
        return "", "", fmt.Errorf("Pas de résultat")
    }

    return result[0].Latitude, result[0].Longitude, nil
}

func preloadCoordinates(artists []Artist) []Artist {

    locationCache := make(map[string]Coordinate)
	fmt.Println("Chargement des coordonnées...")

    for i, artist := range artists {

        var coords []Coordinate

        for _, location := range artist.Locations {

            if cached, exists := locationCache[location]; exists {
                coords = append(coords, cached)
                continue
            }

            latitude, longitude, err := getCoordinatesFromLocation(location)
            if err == nil {

                coordinate := Coordinate{
                    Latitude: latitude,
                    Longitude: longitude,
                }

                locationCache[location] = coordinate
                coords = append(coords, coordinate)
                time.Sleep(500 * time.Millisecond)
            }
        }

        artists[i].Coordinates = coords
    }

    return artists
}