package main

import (
	"Groupie-Tracker/Api"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Artists []Api.Artist
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		artists, err := Api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes: %v", err)
			http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists: artists,
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

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
