package main

import (
	"log"
	"net/http"
	"fmt"
)

func main() {
	server := NewServer()
	server.RegisterRoutes()

	fmt.Println("Serveur lancé : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
