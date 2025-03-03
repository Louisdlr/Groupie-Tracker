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
	Locations    []string
	ConcertDates []string
	Relations    map[string][]string
	ArtistsJSON  string
}

func FilterArtistsByMembers(artists []Api.Artist, minMembers, maxMembers int) []Api.Artist {
	var filtered []Api.Artist
	for _, artist := range artists {
		membersCount := len(artist.Members)
		if membersCount >= minMembers && membersCount <= maxMembers {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

func main() {
	// Serveur pour les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Route : page d'accueil (index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		// Récupération des valeurs de filtrage à partir des paramètres de requête
		minStr := r.URL.Query().Get("min")
		maxStr := r.URL.Query().Get("max")

		var minMembers, maxMembers int
		if minStr != "" {
			minMembers, err = strconv.Atoi(minStr)
			if err != nil {
				minMembers = 0 // Si l'entrée est invalide, prendre 0 comme valeur par défaut
			}
		}
		if maxStr != "" {
			maxMembers, err = strconv.Atoi(maxStr)
			if err != nil {
				maxMembers = 1000 // Valeur par défaut pour max
			}
		}

		// Filtrage des artistes par nombre de membres
		artists = FilterArtistsByMembers(artists, minMembers, maxMembers)

		pageData := PageData{
			Artists: artists,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template : %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template : %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/filtered_catalogue", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		// Récupération des valeurs depuis l'URL
		minStr := r.URL.Query().Get("min")
		maxStr := r.URL.Query().Get("max")

		minMembers, errMin := strconv.Atoi(minStr)
		maxMembers, errMax := strconv.Atoi(maxStr)

		if errMin == nil && errMax == nil {
			artists = FilterArtistsByMembers(artists, minMembers, maxMembers)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(artists)
	})

	// Route : cartes des artistes (cards.html)
	http.HandleFunc("/cards.html", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists: artists,
		}

		tmpl, err := template.ParseFiles("templates/cards.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template : %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template : %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Route : détail d'un artiste (artist_detail.html)
	http.HandleFunc("/artist_detail.html", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

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

		locations, err := Api.GetLocations(artist.LocationsURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des localisations : %v", err)
			http.Error(w, "Erreur de récupération des localisations", http.StatusInternalServerError)
			return
		}

		concertDates, err := Api.GetConcertDates(artist.ConcertDatesURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des dates de concert : %v", err)
			http.Error(w, "Erreur de récupération des dates de concert", http.StatusInternalServerError)
			return
		}

		relations, err := Api.GetRelations(artist.RelationsURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des relations : %v", err)
			http.Error(w, "Erreur de récupération des relations", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artist:       artist,
			Locations:    locations,
			ConcertDates: concertDates,
			Relations:    relations,
		}

		tmpl, err := template.ParseFiles("templates/artist_detail.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template : %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template : %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Route : catalogue des artistes (catalog.html)
	http.HandleFunc("/catalogue.html", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists: artists,
		}

		tmpl, err := template.ParseFiles("templates/catalogue.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template : %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template : %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Route : détail d'un artiste (artist_detail_catalogue.html)
	http.HandleFunc("/artist_detail_catalogue.html", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes : %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

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

		locations, err := Api.GetLocations(artist.LocationsURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des localisations : %v", err)
			http.Error(w, "Erreur de récupération des localisations", http.StatusInternalServerError)
			return
		}

		concertDates, err := Api.GetConcertDates(artist.ConcertDatesURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des dates de concert : %v", err)
			http.Error(w, "Erreur de récupération des dates de concert", http.StatusInternalServerError)
			return
		}

		relations, err := Api.GetRelations(artist.RelationsURL)
		if err != nil {
			log.Printf("Erreur lors de la récupération des relations : %v", err)
			http.Error(w, "Erreur de récupération des relations", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artist:       artist,
			Locations:    locations,
			ConcertDates: concertDates,
			Relations:    relations,
		}

		tmpl, err := template.ParseFiles("templates/artist_detail_catalogue.html")
		if err != nil {
			log.Printf("Erreur lors du chargement du template : %v", err)
			http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template : %v", err)
			http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		}
	})

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
