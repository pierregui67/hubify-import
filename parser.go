package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func DetectDelimiter(filePath string) (rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counts := map[rune]int{',': 0, ';': 0}

	// Read first few lines to detect the most common delimiter
	for i := 0; i < 5 && scanner.Scan(); i++ {
		line := scanner.Text()
		counts[','] += strings.Count(line, ",")
		counts[';'] += strings.Count(line, ";")
	}

	// Choose the most frequent delimiter
	if counts[';'] > counts[','] {
		return ';', nil
	}
	return ',', nil
}

func ValidateLine(line []string, structure StructureDefinition, lineNumber int, errorsMap map[int][]string, wg *sync.WaitGroup) {
	defer wg.Done() 
	
	for index, fieldDef := range structure {
		i, err := strconv.Atoi(index)
		if err != nil || i >= len(line) {
			continue
		}

		value := line[i]
		if err := ValidateField(value, fieldDef.Type, fieldDef.ExpectedValues ); err != nil {
			errorsMap[i] = append(errorsMap[i], fmt.Sprintf("Erreur ligne %d: %v", lineNumber, err))
		}
	}
}

func ValidateCsv(filePath string, structure StructureDefinition) {
	startTime := time.Now()
	delimiter, err := DetectDelimiter(filePath)
	if err != nil {
		fmt.Println("Erreur de détection du délimiteur:", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erreur d'ouverture du fichier:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	errorsMap := make(map[int][]string)

	var wg sync.WaitGroup
	var originalHeader []string
	var rows [][]string

	for scanner.Scan() {
		
		line := scanner.Text()
		record := strings.Split(line, string(delimiter))

		// header
		if lineNumber == 1 {
			originalHeader = record
			rows = append(rows, record)
			lineNumber++
			continue
		}

		TransformRecord(record, structure)

		wg.Add(1)
		go ValidateLine(record, structure, lineNumber, errorsMap, &wg)
		rows = append(rows, record)
		lineNumber++
	}

	wg.Wait()
	
	elapsedTime := time.Since(startTime)
	fmt.Printf("Validation terminée en %s\n", elapsedTime)
	if len(errorsMap) == 0 {
		fmt.Println("✅ No errors found! Saving transformed file...")
		saveTransformedCsv(filePath, originalHeader, rows, structure)
	} else {
		printColumnErrors(errorsMap)
	}
}

func printColumnErrors(errorsMap map[int][]string) {
	for columnIndex, columnErrors := range errorsMap {
		if len(columnErrors) > 0 {
			fmt.Printf("Colonne %d a %d erreur(s):\n", columnIndex+1, len(columnErrors))
			for _, err := range columnErrors {
				fmt.Println("  -", err)
			}
		}
	}
}

func saveTransformedCsv(originalFile string, oldHeader []string, rows [][]string, structure StructureDefinition) {
	// Generate new file name with timestamp
	startTime := time.Now()

	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("%s_%s.csv", strings.TrimSuffix(originalFile, ".csv"), timestamp)

	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Replace column headers with "target" values
	newHeader := make([]string, len(oldHeader))
	for i, header := range oldHeader {
		if fieldDef, exists := structure[strconv.Itoa(i)]; exists {
			newHeader[i] = fieldDef.Target // Use target name
		} else {
			newHeader[i] = header // Keep original
		}
	}

	// Write new header and rows
	writer.Write(newHeader)
	writer.WriteAll(rows[1:]) // Skip old header
	writer.Flush()

	elapsedTime := time.Since(startTime)
	fmt.Printf("✅ CSV transformed and saved as %s in %s ms", newFileName, elapsedTime)
}

func TransformRecord(record []string, structure StructureDefinition) {

	for index, fieldDef := range structure {
		i, err := strconv.Atoi(index)
		if err != nil || i >= len(record) {
			continue
		}

		for _, transformation := range fieldDef.Transformations {
			record[i] = transformation.Apply(record[i])
		}
	}
}
