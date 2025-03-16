package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FieldDefinition struct {
	Type           string        `json:"type"`
	Target         string        `json:"target"`
    Transformations []TransformationAction `json:"transformations,omitempty"`
	ExpectedValues []string			`json:"equalValues,omitempty"`}

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

func (f *FieldDefinition) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type           string          `json:"type"`
		Target         string          `json:"target"`
		Transformations []json.RawMessage `json:"transformations,omitempty"`
		ExpectedValues []string			`json:"equalValues,omitempty"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	f.Type = raw.Type
	f.Target = raw.Target
	f.ExpectedValues = raw.ExpectedValues

	for _, tData := range raw.Transformations {
		// Detect transformation type
		var tType struct {
			Action string `json:"action"`
		}
		if err := json.Unmarshal(tData, &tType); err != nil {
			return err
		}

		var transformation TransformationAction
		switch tType.Action {
		case "trim":
			transformation = TrimTransformation{}
		case "convert":
			var convert ConvertTransformation
			if err := json.Unmarshal(tData, &convert); err != nil {
				return err
			}
			transformation = convert
		case "substring":
			var substring SubstringTransformation
			if err := json.Unmarshal(tData, &substring); err != nil {
				return err
			}
			transformation = substring
		default:
			return fmt.Errorf("unknown transformation: %s", tType.Action)
		}

		f.Transformations = append(f.Transformations, transformation)
	}

	return nil
}

