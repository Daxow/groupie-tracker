package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewServer()
	server.RegisterRoutes()

	log.Println("Serveur lancé : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}