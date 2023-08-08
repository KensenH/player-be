package valmsg

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// error message for golang playground validate

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "alpha":
		return fmt.Sprintf("%s can only be alpha characters", fe.Field())
	case "e164":
		return "Invalid phone number, example: '+62813XXXXXXXX'"
	}
	return fe.Error()
}
