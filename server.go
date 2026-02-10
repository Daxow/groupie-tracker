package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	Artists   []Artist
	Locations []Location
	Dates     []Date
	Relations []Relations
	Templates *template.Template
}

func NewServer() *Server {
	artists := fetchArtists()
	locations := fetchLocations()
	dates := fetchDates()
	relations := fetchRelations()

	templates := template.Must(template.ParseGlob("templates/*.html"))

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
			for _, location := range artist.Location {
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
