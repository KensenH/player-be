package valmsg

import "github.com/go-playground/validator/v10"

// error message for golang playground validate

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return fe.Error()
}
