package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(newUser interface{}) []string {
	v := validator.New()
	if err := v.Struct(newUser); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			if fieldName == "" {
				fieldName = err.StructField()
			}
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required", fieldName))
			case "email":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be a valid email", fieldName))
			case "min":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param()))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be at most %s characters long", fieldName, err.Param()))
			}
		}
		return validationErrors
	}
	return nil
}
