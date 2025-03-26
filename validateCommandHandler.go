package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"bufio"
	"strings"
)

// RequestPayload représente les paramètres attendus dans la requête.
type RequestPayload struct {
	Structure map[string]FieldDefinition `json:"structure"`
	CSVURL    string                 `json:"csv_url"`
	Size      int                    `json:"size,omitempty"`
}

func HandleCSVValidation(w http.ResponseWriter, r *http.Request) {
	// Lire et décoder le JSON de la requête
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Valider le CSV et générer le fichier de sortie
	transformedData, err := ValidateCsv(payload.CSVURL, payload.Structure)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusInternalServerError)
		return
	}

	if payload.Size > 0 {
		// Lire les premières lignes du fichier de sortie
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=validated_data.csv")

		csvWriter := csv.NewWriter(w)
		defer csvWriter.Flush()

		previewSize := min(len(transformedData), payload.Size+1)
		for _, record := range transformedData[:previewSize] {
			if err := csvWriter.Write(record); err != nil {
				http.Error(w, fmt.Sprintf("Error writing CSV: %v", err), http.StatusInternalServerError)
				return
			}
		}
		return
	}
	
	outputFile, err := saveTransformedCsv(payload.CSVURL, transformedData, payload.Structure)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving transformed CSV: %v", err), http.StatusInternalServerError)
		return
	}
	response := map[string]string{"output_file": outputFile}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReadFirstLines lit les premières lignes d'un fichier CSV.
func ReadFirstLines(filePath string, size int) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([][]string, 0, size)
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan() && (size == 0 || i < size); i++ {
		lines = append(lines, strings.Split(scanner.Text(), ","))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}