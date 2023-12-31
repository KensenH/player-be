package valmsg

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// error message for golang playground validate

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "Invalid email"
	case "alpha":
		return fmt.Sprintf("%s can only be alpha characters", fe.Field())
	case "e164":
		return "Invalid phone number (ex: '+62813XXXXXXXX')"
	case "datetime":
		return fmt.Sprintf("%s date is invalid (ex: '02-01-2006', or 'DD-MM-YYYY')", fe.Field())
	}
	return fe.Error()
}
