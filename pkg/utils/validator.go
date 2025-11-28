package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationErrorResponse represents the structure of validation errors
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns friendly error messages
func ValidateStruct(s interface{}) []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationErrorResponse
			element.Field = strings.ToLower(err.Field())
			element.Message = msgForTag(err)
			errors = append(errors, &element)
		}
	}
	return errors
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	}
	return fe.Error() // Default error
}