package main

import (
	"Groupie-Tracker/Api"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	Artists      []Api.Artist
	Artist       Api.Artist
	ConcertDates []Api.ConcertDate // Ajouter les dates de concerts
	Locations    []Api.Location    // Ajouter les localisations
	Relations    []Api.Relation    // Ajouter les relations
	ArtistsJSON  string            // Pour passer les artistes en JSON pour le JavaScript
}

func main() {

	// Gérer les fichiers statiques (CSS, JS, images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve la page d'accueil (index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes: %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		// Convertir les artistes en JSON
		artistsJSON, err := json.Marshal(artists)
		if err != nil {
			log.Printf("Erreur lors de la conversion des artistes en JSON: %v", err)
			http.Error(w, "Erreur de conversion des artistes", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists:     artists,
			ArtistsJSON: string(artistsJSON), // Passer les artistes JSON à la page
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template: %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template: %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Serve la page des cartes d'artistes (cards.html)
	http.HandleFunc("/cards.html", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes: %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		// Convertir les artistes en JSON
		artistsJSON, err := json.Marshal(artists)
		if err != nil {
			log.Printf("Erreur lors de la conversion des artistes en JSON: %v", err)
			http.Error(w, "Erreur de conversion des artistes", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists:     artists,
			ArtistsJSON: string(artistsJSON), // Passer les artistes JSON à la page
		}

		tmpl, err := template.ParseFiles("templates/cards.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template: %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template: %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/artist_detail.html", func(w http.ResponseWriter, r *http.Request) {
		// Récupère l'ID de l'artiste passé dans l'URL (par exemple, /artist_detail.html?id=3)
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		// Récupère les artistes depuis l'API
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes: %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		// Cherche l'artiste correspondant à l'ID
		var artist Api.Artist
		for _, a := range artists {
			if a.ID == id {
				artist = a
				break
			}
		}

		if artist.ID == 0 {
			http.Error(w, "Artiste non trouvé", http.StatusNotFound)
			return
		}

		// Récupérer les autres données (dates de concert, localisations, relations)
		concertDates, err := Api.GetConcertDates()
		if err != nil {
			log.Printf("Erreur lors de la récupération des dates de concerts: %v", err)
			http.Error(w, "Erreur de récupération des dates de concerts", http.StatusInternalServerError)
			return
		}

		locations, err := Api.GetLocations()
		if err != nil {
			log.Printf("Erreur lors de la récupération des localisations: %v", err)
			http.Error(w, "Erreur de récupération des localisations", http.StatusInternalServerError)
			return
		}

		relations, err := Api.GetRelations()
		if err != nil {
			log.Printf("Erreur lors de la récupération des relations: %v", err)
			http.Error(w, "Erreur de récupération des relations", http.StatusInternalServerError)
			return
		}

		// Préparer les données pour le template
		pageData := PageData{
			Artist:       artist,
			ConcertDates: concertDates,
			Locations:    locations,
			Relations:    relations,
		}

		// Charger et exécuter le template
		tmpl, err := template.ParseFiles("templates/artist_detail.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template: %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template: %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Démarrer le serveur sur le port 8080
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
