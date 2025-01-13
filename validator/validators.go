package validator

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
)

// IsNotEmpty checks if a value is not empty.
func IsNotEmpty(value interface{}) error {
	if value == nil {
		return errors.New("value is nil")
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return errors.New("value is empty")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return errors.New("value is zero")
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return errors.New("value is zero")
		}
	case reflect.Bool:
		if !v.Bool() {
			return errors.New("value is false")
		}
	default:
		// For unsupported types, assume the value is not empty
		return nil
	}

	return nil
}

// IsAlphanumeric checks if a string contains only alphanumeric characters.
func IsAlphanumeric(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}
	for _, char := range str {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return errors.New("value contains invalid characters")
		}
	}
	return nil
}

// MinLength checks if a string meets a minimum length requirement.
func MinLength(min int) ValidatorFunc {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		if len(str) < min {
			return fmt.Errorf("value must be at least %d characters long", min)
		}
		return nil
	}
}

// MaxLength checks if a string meets a maximum length requirement.
func MaxLength(max int) ValidatorFunc {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		if len(str) > max {
			return fmt.Errorf("value must be at most %d characters long", max)
		}
		return nil
	}
}

// IsEmail checks if a string is a valid email address.
func IsEmail(value interface{}) error {
	// Check if the input is a string
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	// Use net/mail to validate the email
	_, err := mail.ParseAddress(str)
	if err != nil {
		return errors.New("value is not a valid email address")
	}

	return nil
}

// Regex validates a string against a regular expression.
func Regex(pattern string) ValidatorFunc {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("Invalid regex pattern: %s", err))
	}
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value must be a string")
		}

		if !re.MatchString(str) {
			return errors.New("value does not match the required pattern")
		}

		return nil
	}
}
