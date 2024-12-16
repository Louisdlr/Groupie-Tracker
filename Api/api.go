package Api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Structures pour stocker les données de l'API
type Artist struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Started    int      `json:"start_year"`
	FirstAlbum string   `json:"first_album_date"`
	Members    []string `json:"members"`
}

// Récupérer les artistes
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
