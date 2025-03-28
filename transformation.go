package main

import "strings"

type TransformationAction interface {
	Apply(value string) string
}

type TrimTransformation struct{}

func (t TrimTransformation) Apply(value string) string {
	return strings.TrimSpace(value)
}

type ConvertTransformation struct {
	Param struct {
		Map             map[string]string `json:"map"`
		IsCaseSensitive bool              `json:"case"`
		Default         string            `json:"default"`
	} `json:"param"`
}

func (c ConvertTransformation) Apply(value string) string {
	currentValue := value
	if !c.Param.IsCaseSensitive {
		currentValue = strings.ToLower(value)
	}

	for key, newValue := range c.Param.Map {
		proposedKey := key
		if !c.Param.IsCaseSensitive {
			proposedKey = strings.ToLower(key)
		}

		if proposedKey == currentValue {
			return newValue
		}
	}

	if c.Param.Default != "" {
		return c.Param.Default
	}
	return value
}


type SubstringTransformation struct {
	Param struct {
		Start int `json:"start"`
		Size  int `json:"size"`
	} `json:"param"`
}

func (s SubstringTransformation) Apply(value string) string {
	if s.Param.Start >= len(value) {
		return ""
	}
	end := s.Param.Start + s.Param.Size
	if end > len(value) {
		end = len(value)
	}
	return value[s.Param.Start:end]
}

