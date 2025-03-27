package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func saveTransformedCsv(originalFile string, rows [][]string, structure StructureDefinition) (string, error) {
	// Generate new file name with timestamp
	startTime := time.Now()

	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("%s_%s.csv", strings.TrimSuffix(originalFile, ".csv"), timestamp)

	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if len(rows) == 0 {
		return "", fmt.Errorf("no data available to save")
	}

	// Replace column headers with "target" values
	newHeader := make([]string, len(rows[0]))
	for i, header := range rows[0] {
		if fieldDef, exists := structure[strconv.Itoa(i)]; exists {
			newHeader[i] = fieldDef.Target // Use target name
		} else {
			newHeader[i] = header // Keep original
		}
	}

	// Write new header and rows
	writer.Write(newHeader)
	writer.WriteAll(rows[1:])
	writer.Flush()

	elapsedTime := time.Since(startTime)
	fmt.Printf("✅ CSV saved as %s in %s ms", newFileName, elapsedTime)
	return newFileName, nil
}