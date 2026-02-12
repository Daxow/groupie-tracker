package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"encoding/json"
)

type Server struct {
	Artists   []Artist
	Locations []Location
	Dates     []Date
	Relations []Relation
	Templates *template.Template
}

func NewServer() *Server {
	artists := fetchArtists()
	locations := fetchLocations()
	dates := fetchDates()
	relations := fetchRelations()

	artists = linkArtistsWithLocations(artists, relations)

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

type PageData struct {
	Query   string
	Artists []Artist
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
	json.NewEncoder(response).Encode(server.Artists)
}