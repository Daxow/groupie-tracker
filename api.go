package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func fetchArtists() []Artist {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur réseau :", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Erreur API : statut", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erreur lecture réponse :", err)
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal("Erreur parsing JSON :", err)
	}

	return artists
}

func fetchLocations() []Location {
	url := "https://groupietrackers.herokuapp.com/api/locations"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur réseau :", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erreur lecture réponse :", err)
	}

	var locations []Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal("Erreur parsing JSON :", err)
	}

	return locations
}

func fetchDates() []Date {
	url := "https://groupietrackers.herokuapp.com/api/dates"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur réseau :", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erreur lecture réponse :", err)
	}

	var dates []Date
	err = json.Unmarshal(body, &dates)
	if err != nil {
		log.Fatal("Erreur parsing JSON :", err)
	}

	return dates
}

func fetchRelations() []Relations {
	url := "https://groupietrackers.herokuapp.com/api/relation"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur réseau :", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erreur lecture réponse :", err)
	}

	var relations []Relations
	err = json.Unmarshal(body, &relations)
	if err != nil {
		log.Fatal("Erreur parsing JSON :", err)
	}

	return relations
}