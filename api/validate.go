package api

import "github.com/go-playground/validator"

// ValidationError
// Represents detailed validation error message
type ValidationError struct {
	FailedField string
	Tag         string
	Value       string
}

// ValidateStruct
// Validates the struct values based on struct tags.
func ValidateStruct(v interface{}) []*ValidationError {
	var errors []*ValidationError
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
