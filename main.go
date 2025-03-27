package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Définir les routes HTTP
	http.HandleFunc("/validate-csv", HandleCSVValidation)

	// Démarrer le serveur HTTP sur le port 8085
	port := 8085
	fmt.Printf("Server started on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
