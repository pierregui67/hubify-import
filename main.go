package main

import (
	"fmt"
)

func main() {
	// Load validation structure from JSON file
	structure, err := LoadValidationStructure("structure.json")
	if err != nil {
		fmt.Println("Error loading structure:", err)
		return
	}

	// Valider le fichier CSV avec la structure définie
	ValidateCsv("bigdata.csv", structure)
}
