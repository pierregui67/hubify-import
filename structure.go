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
		Type           string            `json:"type"`
		Target         string            `json:"target"`
		Transformations []json.RawMessage `json:"transformations,omitempty"`
		ExpectedValues []string          `json:"equalValues,omitempty"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	f.Type = raw.Type
	f.Target = raw.Target
	f.ExpectedValues = raw.ExpectedValues

	for _, tData := range raw.Transformations {
		var tType struct {
			Action string `json:"action"`
		}
		if err := json.Unmarshal(tData, &tType); err != nil {
			return err
		}

		transformation, err := doUnmarshal(tType.Action, tData)
		if err != nil {
			return err
		}

		f.Transformations = append(f.Transformations, transformation)
	}

	return nil
}

func doUnmarshal(action string, data json.RawMessage) (TransformationAction, error) {
	var transformation TransformationAction

	switch action {
	case "trim":
		transformation = TrimTransformation{}
	case "convert":
		transformation = &ConvertTransformation{}
	case "before":
		transformation = &BeforeTransformation{}
	case "after":
		transformation = &AfterTransformation{}
	case "substring":
		transformation = &SubstringTransformation{}
	case "concat":
		transformation = &ConcatTransformation{}
	case "addString":
		transformation = &AddTransformation{}
	default:
		return nil, fmt.Errorf("unknown transformation: %s", action)
	}

	if err := json.Unmarshal(data, transformation); err != nil {
		return nil, err
	}

	return transformation, nil
}