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
		log.Println("Erreur réseau :", err)
        return []Artist{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Erreur API : statut", resp.StatusCode)
        return []Artist{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lecture réponse :", err)
        return []Artist{}
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Println("Erreur parsing JSON :", err)
        return []Artist{}
	}

	return artists
}

func fetchLocations() []Location {
	url := "https://groupietrackers.herokuapp.com/api/locations"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur réseau :", err)
		return []Location{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lecture réponse :", err)
		return []Location{}
	}

	var response LocationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Erreur parsing JSON :", err)
		return []Location{}
	}

	return response.Index
}

func fetchDates() []Date {
	url := "https://groupietrackers.herokuapp.com/api/dates"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur réseau :", err)
		return []Date{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lecture réponse :", err)
		return []Date{}
	}

	var response DateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Erreur parsing JSON :", err)
		return []Date{}
	}

	return response.Index
}

func fetchRelations() []Relation {
	url := "https://groupietrackers.herokuapp.com/api/relation"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur réseau :", err)
		return []Relation{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lecture réponse :", err)
		return []Relation{}
	}

	var response RelationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Erreur parsing JSON :", err)
		return []Relation{}
	}

	return response.Index
}