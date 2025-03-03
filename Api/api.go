package Api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Structures des données

type Artist struct {
	ID              int                 `json:"id"`
	Image           string              `json:"image"`
	Name            string              `json:"name"`
	Members         []string            `json:"members"`
	CreationDate    int                 `json:"creationDate"`
	FirstAlbum      string              `json:"firstAlbum"`
	LocationsURL    string              `json:"locations"`
	ConcertDatesURL string              `json:"concertDates"`
	RelationsURL    string              `json:"relations"`
	Relations       map[string][]string `json:"-"`
	Dates           []string            `json:"-"`
	Locations       []string            `json:"-"`
}

type Location struct {
	Locations []string `json:"locations"`
}

type ConcertDate struct {
	Dates []string `json:"dates"`
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Fonctions pour récupérer les données des APIs

// Récupère les artistes
func GetArtists() ([]Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := fetchData[[]Artist](url)
	if err != nil {
		return nil, err
	}

	for i, artist := range artists {
		relations, err := GetRelations(artist.RelationsURL)
		if err != nil {
			return nil, err
		}
		artists[i].Relations = relations

		// Extract dates and locations from relations
		var dates []string
		var locations []string
		for location, dateList := range relations {
			locations = append(locations, location)
			dates = append(dates, dateList...)
		}
		artists[i].Dates = dates
		artists[i].Locations = locations
	}

	return artists, nil
}

// Récupère les localisations
func GetLocations(url string) ([]string, error) {
	var locData Location
	err := fetchDataFromURL(url, &locData)
	if err != nil {
		return nil, err
	}
	return locData.Locations, nil
}

// Récupère les dates de concert
func GetConcertDates(url string) ([]string, error) {
	var dateData ConcertDate
	err := fetchDataFromURL(url, &dateData)
	if err != nil {
		return nil, err
	}
	return dateData.Dates, nil
}

// Récupère les relations
func GetRelations(url string) (map[string][]string, error) {
	var relData Relation
	err := fetchDataFromURL(url, &relData)
	if err != nil {
		return nil, err
	}
	return relData.DatesLocations, nil
}

// Fonction générique pour récupérer des données JSON
func fetchData[T any](url string) (T, error) {
	resp, err := http.Get(url)
	if err != nil {
		return *new(T), fmt.Errorf("erreur lors de la requête : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *new(T), fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
	}

	var data T
	err = json.Unmarshal(body, &data)
	if err != nil {
		return *new(T), fmt.Errorf("erreur lors du décodage JSON : %w", err)
	}

	return data, nil
}

// Fonction pour récupérer et décoder des données à partir d'un URL spécifique
func fetchDataFromURL(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erreur lors de la requête : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("erreur lors du décodage JSON : %w", err)
	}

	return nil
}
