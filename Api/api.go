package Api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Structure pour stocker les informations d'un artiste
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Structure pour les localisations retournées par l'API
type LocationsResponse struct {
	Locations []Location `json:"locations"`
}

// Structure pour une localisation
type Location struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
}

// Structure pour les dates de concerts retournées par l'API
type ConcertDatesResponse struct {
	ConcertDates []ConcertDate `json:"concertDates"`
}

// Structure pour une date de concert
type ConcertDate struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Location string `json:"location"`
}

// Structure pour les relations retournées par l'API
type RelationsResponse struct {
	Relations []Relation `json:"relations"`
}

// Structure pour une relation
type Relation struct {
	ID       int    `json:"id"`
	Relation string `json:"relation"`
}

// Fonction pour récupérer les artistes depuis l'API
func GetArtists() ([]Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des artistes: %v", err)
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("Erreur de décodage JSON des artistes: %v", err)
	}

	return artists, nil
}

// Fonction pour récupérer les localisations depuis l'API
func GetLocations() ([]Location, error) {
	url := "https://groupietrackers.herokuapp.com/api/locations"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des localisations: %v", err)
	}
	defer resp.Body.Close()

	var locationsResponse LocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&locationsResponse); err != nil {
		return nil, fmt.Errorf("Erreur de décodage JSON des localisations: %v", err)
	}

	return locationsResponse.Locations, nil
}

// Fonction pour récupérer les dates de concerts depuis l'API
func GetConcertDates() ([]ConcertDate, error) {
	url := "https://groupietrackers.herokuapp.com/api/dates"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des dates de concerts: %v", err)
	}
	defer resp.Body.Close()

	var concertDatesResponse ConcertDatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&concertDatesResponse); err != nil {
		return nil, fmt.Errorf("Erreur de décodage JSON des dates de concerts: %v", err)
	}

	return concertDatesResponse.ConcertDates, nil
}

// Fonction pour récupérer les relations depuis l'API
func GetRelations() ([]Relation, error) {
	url := "https://groupietrackers.herokuapp.com/api/relation"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des relations: %v", err)
	}
	defer resp.Body.Close()

	var relationsResponse RelationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&relationsResponse); err != nil {
		return nil, fmt.Errorf("Erreur de décodage JSON des relations: %v", err)
	}

	return relationsResponse.Relations, nil
}
