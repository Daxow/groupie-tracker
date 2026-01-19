package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
	"strings"
)

type Server struct {
	Artists   []Artist
	Templates *template.Template
}

func NewServer() *Server {
	artists := fetchArtists()

	templates := template.Must(template.ParseGlob("templates/*.html"))

	return &Server{
		Artists:   artists,
		Templates: templates,
	}
}

func (s *Server) RegisterRoutes() {
	http.HandleFunc("/", s.handleHome)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s.renderPage(w, "index.html", s.Artists)
}

func (s *Server) renderPage(w http.ResponseWriter, name string, data []Artist) {
	err := s.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
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
