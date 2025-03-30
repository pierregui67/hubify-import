package main

import "strings"

type TransformationAction interface {
	Apply(row []string, index int) string
}

type TrimTransformation struct{}

func (t TrimTransformation) Apply(row []string, index int) string {
	return strings.TrimSpace(row[index])
}

type ConvertTransformation struct {
	Param struct {
		Map             map[string]string `json:"map"`
		IsCaseSensitive bool              `json:"case"`
		Default         string            `json:"default"`
	} `json:"param"`
}

func (c ConvertTransformation) Apply(row []string, index int) string {
	value:= row[index]
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


type BeforeTransformation struct {
	Param struct {
		Delimiter string `json:"delimiter"`
	} `json:"param"`
}


type SubstringTransformation struct {
	Param struct {
		Start int `json:"start"`
		Size  int `json:"size,omitempty"`
	} `json:"param"`
}

func (s SubstringTransformation) Apply(row []string, index int) string {
	value:= row[index]
	if(s.Param.Size == 0){
		s.Param.Size = len(value)
	}
	if s.Param.Start >= len(value) {
		return ""
	}
	end := s.Param.Start + s.Param.Size
	if end > len(value) {
		end = len(value)
	}
	return value[s.Param.Start:end]
}

func (b BeforeTransformation) Apply(row []string, index int) string {
	value:= row[index]
	if b.Param.Delimiter == "" {
		return value
	}

	delimiterIndex := strings.Index(value, b.Param.Delimiter)
	if delimiterIndex == -1 {
		return value
	}

	return value[:delimiterIndex]
}

type AfterTransformation struct {
	Param struct {
		Delimiter string `json:"delimiter"`
	} `json:"param"`
}

func (a AfterTransformation) Apply(row []string, index int) string {
	value:= row[index]
	if a.Param.Delimiter == "" {
		return value // If no delimiter is set, return the full string
	}

	delimiterIndex := strings.Index(value, a.Param.Delimiter)
	if delimiterIndex == -1 {
		return value // If the delimiter is not found, return the full string
	}

	return value[delimiterIndex+len(a.Param.Delimiter):] // Extract everything after the delimiter
}

type ConcatTransformation struct {
	Param struct {
		ColumnId int `json:"columnId"`
	} `json:"param"`
}

func (c ConcatTransformation) Apply(row []string, index int) string {

	if c.Param.ColumnId < 0 || c.Param.ColumnId >= len(row) {
		return row[index]
	}
	return row[index] + row[c.Param.ColumnId]
}

type AddTransformation struct {
	Param struct {
		AddString string `json:"addString"`
	} `json:"param"`
}

func (c AddTransformation) Apply(row []string, index int) string {
	return row[index] + c.Param.AddString
}
