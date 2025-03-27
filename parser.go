package main

import (
	"bufio"
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

func ValidateLine(line []string, structure StructureDefinition, lineNumber int, errorsMap map[int][]string) {
	
	for index, fieldDef := range structure {
		i, err := strconv.Atoi(index)
		if err != nil || i >= len(line) {
			continue
		}

		value := line[i]
		if err := ValidateField(value, fieldDef.Type, fieldDef.ExpectedValues ); err != nil {
			errorsMap[i] = append(errorsMap[i], fmt.Sprintf("Erreur ligne %d: %v", lineNumber-1, err))
		}
	}
}

func ValidateCsv(filePath string, structure StructureDefinition) ([][]string, error) {
	startTime := time.Now()
	delimiter, err := DetectDelimiter(filePath)
	if err != nil {
		fmt.Println("Erreur de détection du délimiteur:", err)
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erreur d'ouverture du fichier:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	errorsMap := make(map[int][]string)

	var wg sync.WaitGroup
	var rows [][]string

	for scanner.Scan() {
		
		line := scanner.Text()
		record := strings.Split(line, string(delimiter))

		// header
		if lineNumber == 1 {
			rows = append(rows, record)
			lineNumber++
			continue
		}

		wg.Add(1)
		TransformAndValidateLine(record, structure, lineNumber, errorsMap, &wg)
		
		rows = append(rows, record)
		lineNumber++
	}

	wg.Wait()
	
	elapsedTime := time.Since(startTime)
	fmt.Printf("Validation terminée en %s\n", elapsedTime)
	if len(errorsMap) == 0 {
		fmt.Println("✅ No errors found! Returning transformed data...")
		return rows, nil
	} else {
		errorsString := printColumnErrors(errorsMap)
		return nil, fmt.Errorf(errorsString)
	}
}

func printColumnErrors(errorsMap map[int][]string) (string) {
	errorsString := "";
	for columnIndex, columnErrors := range errorsMap {
		if len(columnErrors) > 0 {
			fmt.Printf("Colonne %d a %d erreur(s):\n", columnIndex+1, len(columnErrors))
			errorsString += fmt.Sprintf("Colonne %d a %d erreur(s):\n", columnIndex+1, len(columnErrors))
			for _, err := range columnErrors {
				fmt.Println("  -", err)
				errorsString += "  - " + err + "\n"
			}
		}
	}
	return errorsString
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

func TransformAndValidateLine(record []string, structure StructureDefinition, lineNumber int, errorsMap map[int][]string, wg *sync.WaitGroup) {
	defer wg.Done() 
	TransformRecord(record, structure)
	ValidateLine(record, structure, lineNumber, errorsMap)
}
