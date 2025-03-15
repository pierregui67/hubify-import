package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Transformation struct {
	Action string `json:"action"`
	Start  int    `json:"start,omitempty"`
	Size   int    `json:"size,omitempty"`
}

type FieldDefinition struct {
	Type           string        `json:"type"`
	Target         string        `json:"target"`
	Transformation Transformation `json:"transformation,omitempty"`
}

type StructureDefinition map[string]FieldDefinition

func LoadValidationStructure(filePath string) (StructureDefinition, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var structure StructureDefinition
	if err := json.Unmarshal(file, &structure); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return structure, nil
}
