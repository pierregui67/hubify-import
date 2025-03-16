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
	Map              map[string]string `json:"map"`
	IsCaseSensitive  bool              `json:"casse"`
	Default  		 string            `json:"default"`
}

func (c ConvertTransformation) Apply(value string) string {
	currentValue := value
	if !c.IsCaseSensitive {
		currentValue = strings.ToLower(value)
	}

	for key, newValue := range c.Map {
		proposedKey := key
		if !c.IsCaseSensitive {
			proposedKey = strings.ToLower(key)
		}

		if proposedKey == currentValue {
			return newValue
		}
	}

	if c.Default != "" {
		return c.Default
	}
	return value
}

type SubstringTransformation struct {
	Start int `json:"start"`
	Size  int `json:"size"`
}

func (s SubstringTransformation) Apply(value string) string {
	if s.Start >= len(value) {
		return ""
	}
	end := s.Start + s.Size
	if end > len(value) {
		end = len(value)
	}
	return value[s.Start:end]
}
