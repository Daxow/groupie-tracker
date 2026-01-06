package main

import (
	"html/template"
	"log"
	"net/http"
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