package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

var validate = validator.New()

func ValidateField(value string, expectedType string, expectedValues []string) error {

	if value == "" {
		return fmt.Errorf("expected type '%s' but received an empty value", expectedType)
	}

	switch expectedType {
	case "email":
		if err := validate.Var(value, "email"); err != nil {
			return fmt.Errorf("invalid email format: '%s'", value)
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("expected a float, but got '%s'", value)
		}
	case "date":
		if err := validate.Var(value, "datetime=2006-01-02"); err != nil {
			return fmt.Errorf("invalid date format (expected YYYY-MM-DD) but got '%s'", value)
		}
	case "required":
		if err := validate.Var(value, "required"); err != nil {
			return fmt.Errorf("this field is required but was empty")
		}
	case "int":
		if err := validate.Var(value, "number"); err != nil {
			return fmt.Errorf("expected an integer but got '%s'", value)
		}
	case "equals":
		if len(expectedValues) > 0 {
			validatorParam := "eq=" + strings.Join(expectedValues, "|eq=")
			if err := validate.Var(value, validatorParam); err != nil {
				return fmt.Errorf("value '%s' does not match allowed options: %s", value, "M")
			}
		}
		
	default:
		return fmt.Errorf("unsupported validation type: '%s'", expectedType)
	}

	return nil
}
