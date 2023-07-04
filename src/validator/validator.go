package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname= regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func validateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := validateString(value, 3, 100); err != nil {
		return err 
	}

	if !isValidUsername(value){
		return fmt.Errorf("Must contain only lowercaseletters, digits, or underscore")
	}
	return nil 
}

func ValidatePassword(value string) error {
	return validateString(value, 3, 100)
}


func ValidateEmail(value string) error {
	if err := validateString(value, 3, 200); err != nil {
		return err 
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil 
}

func ValidateFullName(value string) error {
	if err := validateString(value, 3, 100); err != nil {
		return err 
	}

	if !isValidFullname(value){
		return fmt.Errorf("Must contain only letters or spaces")
	}
	return nil 
}