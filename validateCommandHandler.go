package main

import (
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
		limitedErrors := limitErrors(err, 10000)
		http.Error(w, fmt.Sprintf("%v", limitedErrors), http.StatusInternalServerError)
		return
	}

	if payload.Size > 0 {
		// Set the response header for JSON format
		w.Header().Set("Content-Type", "application/json")
		
		// Limit the number of records based on payload.Size
		previewSize := min(len(transformedData), payload.Size+1)
		filteredData := transformedData[:previewSize]
	
		// The first row contains the keys (headers)
		headers := filteredData[0]
		
		// Create a slice of maps to hold the transformed data
		var result []map[string]string
	
		// Iterate over the data and construct the result as an array of objects
		for _, record := range filteredData[1:] { // Skip the header row
			recordMap := make(map[string]string)
			for i, value := range record {
				if i < len(headers) {
					recordMap[headers[i]] = value
				}
			}
			result = append(result, recordMap)
		}
	
		// Marshal the result into JSON
		jsonData, err := json.Marshal(result)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
			return
		}
	
		// Write the JSON data to the response
		w.Write(jsonData)
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

func limitErrors(err error, max int) error {
	errStr := err.Error()
	if len(errStr) > max {
		return fmt.Errorf("%s...\n(Truncated, too many errors)", errStr[:max])
	}
	return err
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